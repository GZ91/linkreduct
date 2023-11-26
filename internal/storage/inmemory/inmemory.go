package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ConfigerStorage определяет интерфейс для получения конфигурации, связанной с хранилищем.
type ConfigerStorage interface {
	GetMaxIterLen() int
}

// GeneratorRunes определяет интерфейс для генерации случайных строк.
type GeneratorRunes interface {
	RandStringRunes(int) string
}

// New инициализирует новый экземпляр типа DB, представляющий хранилище данных в памяти для URL.
func New(ctx context.Context, conf ConfigerStorage, genrun GeneratorRunes) (*DB, error) {
	return &DB{
		data:         make(map[string]*models.StructURL, 1),
		config:       conf,
		genrun:       genrun,
		chURLsForDel: make(chan models.StructDelURLs),
	}, nil
}

// DB представляет собой хранилище данных в памяти для URL.
type DB struct {
	data          map[string]*models.StructURL
	config        ConfigerStorage
	genrun        GeneratorRunes
	mutex         sync.Mutex
	chsURLsForDel chan []models.StructDelURLs
	chURLsForDel  chan models.StructDelURLs
}

// setDB устанавливает данные для заданного ключа в хранилище в памяти.
func (r *DB) setDB(ctx context.Context, key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Извлечение идентификатора пользователя из контекста, если доступно.
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}

	// Создание новой структуры StructURL и сохранение ее в карте данных.
	StructURL := &models.StructURL{
		OriginalURL: value,
		ShortURL:    key,
		UserID:      UserID,
		ID:          uuid.New().String(),
	}
	r.data[key] = StructURL

	return true
}

// GetURL извлекает оригинальный URL, связанный с заданным коротким ключом.
func (r *DB) GetURL(ctx context.Context, key string) (val string, ok bool, errs error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Извлечение структуры StructURL, связанной с ключом.
	valueStruct, found := r.data[key]
	if found {
		// Проверка, помечен ли URL как удаленный.
		if valueStruct.DeletedFlag {
			return "", false, errorsapp.ErrLineURLDeleted
		}
		ok = found
		val = valueStruct.OriginalURL
	}

	return
}

// AddURL добавляет новый URL в хранилище в памяти, генерируя уникальный короткий ключ для него.
func (r *DB) AddURL(ctx context.Context, url string) (string, error) {
	// Установка начальных параметров для генерации короткого ключа.
	lenID := 5
	iterLen := 0
	MaxIterLen := r.config.GetMaxIterLen()

	// Генерация уникального короткого ключа для нового URL.
	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := r.genrun.RandStringRunes(lenID)

		// Проверка, существует ли уже сгенерированный ключ.
		if _, found, err := r.GetURL(ctx, idString); found {
			if err != nil {
				return "", err
			}
			iterLen++
			continue
		}

		// Установка нового URL в хранилище в памяти.
		r.setDB(ctx, idString, url)
		return idString, nil
	}
}

// Close выполняет необходимую очистку или освобождение ресурсов.
func (r *DB) Close() error {
	return nil
}

// Ping проверяет доступность или состояние хранилища.
func (r *DB) Ping(ctx context.Context) error {
	return nil
}

// FindLongURL находит короткий ключ, связанный с заданным оригинальным URL.
func (r *DB) FindLongURL(ctx context.Context, OriginalURL string) (string, bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Итерация по карте данных для поиска ключа, связанного с оригинальным URL.
	for key, val := range r.data {
		if val.OriginalURL == OriginalURL {
			return key, true, nil
		}
	}

	return "", false, nil
}

// AddBatchLink добавляет пакет URL в хранилище в памяти.
func (r *DB) AddBatchLink(ctx context.Context, batchLink []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {
	for _, data := range batchLink {
		link := data.OriginalURL
		var shortURL string

		// Проверка существования URL.
		shortURL, ok, err := r.FindLongURL(ctx, link)
		if err != nil {
			return nil, err
		}

		if ok {
			// URL уже существует, возврат ошибки.
			errs = errorsapp.ErrLinkAlreadyExists
		} else {
			// Добавление URL в хранилище в памяти.
			var err error
			shortURL, err = r.AddURL(ctx, link)
			if err != nil {
				logger.Log.Error("error when writing an add link to the inmemory storage", zap.Error(err), zap.String("unadded value", link))
				errs = err
				return
			}
		}

		// Добавление освобожденного URL в срез результатов.
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{CorrelationID: data.CorrelationID, ShortURL: shortURL})
	}

	return
}

// GetLinksUser извлекает URL, связанные с конкретным пользователем из хранилища в памяти.
func (r *DB) GetLinksUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	returnData := make([]models.ReturnedStructURL, 0)

	// Итерация по данным хранилища для поиска URL, принадлежащих конкретному пользователю.
	for _, val := range r.data {
		if val.UserID == userID {
			returnData = append(returnData, models.ReturnedStructURL{OriginalURL: val.OriginalURL, ShortURL: val.ShortURL})
		}
	}
	return returnData, nil
}

// InitializingRemovalChannel инициализирует канал для удаления URL.
func (r *DB) InitializingRemovalChannel(ctx context.Context, chsURLs chan []models.StructDelURLs) error {
	r.chsURLsForDel = chsURLs

	// Запуск горутин для обработки данных о удалении и очистки буфера.
	go r.GroupingDataForDeleted(ctx)
	go r.FillBufferDelete(ctx)

	return nil
}

// GroupingDataForDeleted группирует данные для удаления в буфере.
func (r *DB) GroupingDataForDeleted(ctx context.Context) {
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			close(r.chURLsForDel)
			return
		case sliceVal := <-r.chsURLsForDel:
			wg.Add(1)
			go func(*sync.WaitGroup) {
				// Пересылка данных о удалении из среза в канал.
				for _, val := range sliceVal {
					r.chURLsForDel <- val
				}
				wg.Done()
			}(&wg)
		}
	}
}

// FillBufferDelete заполняет буфер данными для удаления.
func (r *DB) FillBufferDelete(ctx context.Context) {
	t := time.NewTicker(time.Second * 10)
	var listForDel []models.StructDelURLs
	for {
		select {
		case <-ctx.Done():
			return
		case val := <-r.chURLsForDel:
			listForDel = append(listForDel, val)
		case <-t.C:
			if len(listForDel) > 0 {
				// Удаление URL, указанных в буфере.
				r.deletedURLs(listForDel)
			}
		}
	}
}

// deletedURLs устанавливает флаг DeletedFlag для URL, указанных в списке для удаления.
func (r *DB) deletedURLs(listForDel []models.StructDelURLs) {
	for _, val := range listForDel {
		for index := range r.data {
			if r.data[index].ShortURL == val.URL && r.data[index].UserID == val.UserID {
				r.data[index].DeletedFlag = true
				break
			}
		}
	}
}
