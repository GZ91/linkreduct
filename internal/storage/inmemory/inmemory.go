package inmemory

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ConfigerStorage interface {
	GetMaxIterLen() int
}

type GeneratorRunes interface {
	RandStringRunes(int) string
}

func New(ctx context.Context, conf ConfigerStorage, genrun GeneratorRunes) (*DB, error) {
	return &DB{data: make(map[string]*models.StructURL, 1), config: conf, genrun: genrun, chURLsForDel: make(chan models.StructDelURLs)}, nil
}

type DB struct {
	data          map[string]*models.StructURL
	config        ConfigerStorage
	genrun        GeneratorRunes
	mutex         sync.Mutex
	chsURLsForDel chan []models.StructDelURLs
	chURLsForDel  chan models.StructDelURLs
}

func (r *DB) setDB(ctx context.Context, key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}
	StructURL := &models.StructURL{
		OriginalURL: value,
		ShortURL:    key,
		UserID:      UserID,
		ID:          uuid.New().String(),
	}
	r.data[key] = StructURL
	return true
}

func (r *DB) GetURL(ctx context.Context, key string) (val string, ok bool, errs error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	valueStruct, found := r.data[key]
	if found {
		if valueStruct.DeletedFlag {
			return "", false, errorsapp.ErrLineURLDeleted
		}
		ok = found
		val = valueStruct.OriginalURL
	}
	return
}

func (r *DB) AddURL(ctx context.Context, url string) (string, error) {

	lenID := 5
	iterLen := 0
	MaxIterLen := r.config.GetMaxIterLen()

	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := r.genrun.RandStringRunes(lenID)
		if _, found, err := r.GetURL(ctx, idString); found {
			if err != nil {
				return "", err
			}
			iterLen++
			continue
		}
		r.setDB(ctx, idString, url)
		return idString, nil
	}
}

func (r *DB) Close() error {
	return nil
}

func (r *DB) Ping(ctx context.Context) error {
	return nil
}

func (r *DB) FindLongURL(ctx context.Context, OriginalURL string) (string, bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, val := range r.data {
		if val.OriginalURL == OriginalURL {
			return key, true, nil
		}
	}
	return "", false, nil
}

func (r *DB) AddBatchLink(ctx context.Context, batchLink []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {
	for _, data := range batchLink {
		link := data.OriginalURL
		var shortURL string
		shortURL, ok, err := r.FindLongURL(ctx, link)
		if err != nil {
			return nil, err
		}
		if ok {
			errs = errorsapp.ErrLinkAlreadyExists
		} else {
			var err error
			shortURL, err = r.AddURL(ctx, link)
			if err != nil {
				logger.Log.Error("error when writing an add link to the inmemory storage", zap.Error(err), zap.String("unadded value", link))
				errs = err
				return
			}
		}
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{CorrelationID: data.CorrelationID, ShortURL: shortURL})
	}
	return
}

func (r *DB) GetLinksUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	returnData := make([]models.ReturnedStructURL, 0)
	for _, val := range r.data {
		if val.UserID == userID {
			returnData = append(returnData, models.ReturnedStructURL{OriginalURL: val.OriginalURL, ShortURL: val.ShortURL})
		}
	}
	return returnData, nil
}

func (r *DB) InitializingRemovalChannel(ctx context.Context, chsURLs chan []models.StructDelURLs) error {
	r.chsURLsForDel = chsURLs
	go r.GroupingDataForDeleted(ctx)
	go r.FillBufferDelete(ctx)
	return nil
}

func (d *DB) GroupingDataForDeleted(ctx context.Context) {
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			close(d.chURLsForDel)
			return
		case sliceVal := <-d.chsURLsForDel:
			wg.Add(1)
			go func(*sync.WaitGroup) {
				for _, val := range sliceVal {
					d.chURLsForDel <- val
				}
				wg.Done()
			}(&wg)
		}
	}
}

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
				r.deletedURLs(listForDel)
			}
		}

	}
}

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
