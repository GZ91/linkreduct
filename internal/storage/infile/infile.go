package infile

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

type GeneratorRunes interface {
	RandStringRunes(int) string
}

type ConfigerStorage interface {
	GetMaxIterLen() int
	GetNameFileStorage() string
	GetStartLenShortLink() int
}

func New(ctx context.Context, conf ConfigerStorage, gen GeneratorRunes) *db {
	DB := &db{
		generatorRunes: gen,
		conf:           conf,
		data:           make(map[string]*models.StructURL),
	}
	DB.open()
	return DB
}

type db struct {
	generatorRunes GeneratorRunes
	conf           ConfigerStorage
	mutex          sync.Mutex
	data           map[string]*models.StructURL
	newdata        []string
	chsURLsForDel  chan []models.StructDelURLs
	chURLsForDel   chan models.StructDelURLs
}

func (r *db) GetURL(ctx context.Context, key string) (string, bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	datainModel, ok := r.data[key]
	if datainModel.DeletedFlag {
		return "", false, errorsapp.ErrLineURLDeleted
	}
	var retval string
	if ok {
		retval = datainModel.OriginalURL
	}
	return retval, ok, nil
}

func (r *db) AddURL(ctx context.Context, url string) (string, error) {
	shortURL, err := r.getShortURL(ctx)
	if err != nil {
		return "", err
	}
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}
	model := &models.StructURL{
		ID:          uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: url,
		UserID:      UserID,
	}
	r.newdata = append(r.newdata, shortURL)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[shortURL] = model
	return shortURL, nil
}

func (r *db) save() (errs error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	nameFile := r.conf.GetNameFileStorage()
	if nameFile == "" {
		return nil
	}
	if len(r.newdata) == 0 {
		return nil
	}

	file, err := os.OpenFile(nameFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("an error occurred when opening the file", zap.String("error", err.Error()))
		return err
	}
	defer file.Close()

	var ovLen int
	var count int

	for _, shorURL := range r.newdata {
		val, ok := r.data[shorURL]
		if !ok {
			logger.Log.Error("record was added to the array of new, but was not found in the storage map")
			continue
		}
		data, err := json.Marshal(val)
		if err != nil {
			logger.Log.Error("when serializing data to json", zap.String("error", err.Error()))
			errs = err
			continue
		}
		data = append(data, '\n')
		len, err := file.Write(data)
		if err != nil {
			logger.Log.Error("when writing a json string to a file", zap.String("error", err.Error()))
			errs = err
			continue
		}
		ovLen += len
		count++
	}

	logger.Log.Info("the file has been successfully written to",
		zap.String("len", strconv.Itoa(ovLen)),
		zap.String("count", strconv.Itoa(count)))

	return
}

func (r *db) open() (errs error) {
	nameFile := r.conf.GetNameFileStorage()
	if nameFile == "" {
		return nil
	}
	file, err := os.OpenFile(r.conf.GetNameFileStorage(), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("error when opening a file with saved data", zap.String("error", err.Error()))
		return err
	}
	defer file.Close()
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		data := scaner.Bytes()
		modelData := &models.StructURL{}
		err := json.Unmarshal(data, modelData)
		if err != nil {
			logger.Log.Error("error when trying to decode a string", zap.String("error", err.Error()))
			errs = err
			continue
		}
		r.data[modelData.ShortURL] = modelData
	}
	return
}

func (r *db) Close() error {
	return r.save()
}

func (r *db) getShortURL(ctx context.Context) (string, error) {
	lenID := r.conf.GetStartLenShortLink()
	iterLen := 0
	MaxIterLen := r.conf.GetMaxIterLen()

	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := r.generatorRunes.RandStringRunes(lenID)
		if _, found, err := r.GetURL(ctx, idString); found {
			if err != nil {
				return "", err
			}
			iterLen++
			continue
		}
		return idString, nil
	}
}

func (r *db) Ping(ctx context.Context) error {
	return nil
}

func (r *db) FindLongURL(ctx context.Context, OriginalURL string) (string, bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, val := range r.data {
		if val.OriginalURL == OriginalURL {
			return key, true, nil
		}
	}
	return "", false, nil
}

func (r *db) AddBatchLink(ctx context.Context, batchLink []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {
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
				logger.Log.Error("error when writing an add link to the file storage", zap.Error(err), zap.String("unadded value", link))
				return nil, err
			}
		}
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{CorrelationID: data.CorrelationID, ShortURL: shortURL})
	}
	return
}

func (r *db) GetLinksUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	returnData := make([]models.ReturnedStructURL, 0)
	for _, val := range r.data {
		if val.UserID == userID {
			returnData = append(returnData, models.ReturnedStructURL{OriginalURL: val.OriginalURL, ShortURL: val.ShortURL})
		}
	}
	return returnData, nil
}

func (r *db) InitializingRemovalChannel(chsURLs chan []models.StructDelURLs) error {
	r.chsURLsForDel = chsURLs
	go r.GroupingDataForDeleted()
	go r.FillBufferDelete()
	return nil
}

func (r *db) GroupingDataForDeleted() {
	for sliceVal := range r.chsURLsForDel {
		sliceVal := sliceVal
		go func() {
			for _, val := range sliceVal {
				r.chURLsForDel <- val
			}
		}()
	}
}

func (r *db) FillBufferDelete() {
	t := time.Tick(time.Second * 10)
	var listForDel []models.StructDelURLs
	for {
		select {
		case val := <-r.chURLsForDel:
			listForDel = append(listForDel, val)
		case <-t:
			if len(listForDel) > 0 {
				r.deletedURLs(listForDel)
			}
		}

	}
}

func (r *db) deletedURLs(listForDel []models.StructDelURLs) {
	for _, val := range listForDel {
		for index, _ := range r.data {
			if r.data[index].ShortURL == val.URL {
				r.data[index].DeletedFlag = true
				break
			}
		}
	}
}
