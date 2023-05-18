package infile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
)

type GeneratorRunes interface {
	RandStringRunes(int) string
}

type ConfigerStorage interface {
	GetMaxIterLen() int
	GetNameFileStorage() string
	GetStartLenShortLink() int
}

func New(conf ConfigerStorage, gen GeneratorRunes) *db {
	DB := &db{
		generatorRunes: gen,
		conf:           conf,
		data:           make(map[string]models.StructURL),
	}
	DB.open()
	return DB
}

type db struct {
	generatorRunes GeneratorRunes
	conf           ConfigerStorage
	mutex          sync.Mutex
	data           map[string]models.StructURL
	newdata        []string
}

func (r *db) GetURL(key string) (string, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	datainModel, ok := r.data[key]
	var retval string
	if ok {
		retval = datainModel.OriginalURL
	}
	return retval, ok
}

func (r *db) AddURL(url string) string {
	shortURL := r.getShortURL()

	model := models.StructURL{
		ID:          uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: url,
	}
	r.newdata = append(r.newdata, shortURL)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[shortURL] = model
	return shortURL
}

func (r *db) save() (errs error) {
	defer func() {
		if rec := recover(); rec != nil {
			if errs == nil {
				errs = errors.New(fmt.Sprint(rec))
			} else {
				errors.Wrap(errs, fmt.Sprint(rec))
			}
		}
	}()
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
			errors.Wrap(errs, err.Error())
			continue
		}
		data = append(data, '\n')
		len, err := file.Write(data)
		if err != nil {
			logger.Log.Error("when writing a json string to a file", zap.String("error", err.Error()))
			errors.Wrap(errs, err.Error())
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
		modelData := models.StructURL{}
		err := json.Unmarshal(data, &modelData)
		if err != nil {
			logger.Log.Error("error when trying to decode a string", zap.String("error", err.Error()))
			errors.Wrap(errs, err.Error())
			continue
		}
		r.data[modelData.ShortURL] = modelData
	}
	return
}

func (r *db) Close() error {
	return r.save()
}

func (r *db) getShortURL() string {
	lenID := r.conf.GetStartLenShortLink()
	iterLen := 0
	MaxIterLen := r.conf.GetMaxIterLen()
	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := r.generatorRunes.RandStringRunes(lenID)
		if _, found := r.GetURL(idString); found {
			iterLen++
			continue
		}
		return idString
	}
}

func (r *db) Ping() error {
	return nil
}

func (r *db) FindLongURL(OriginalURL string) (string, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, val := range r.data {
		if val.OriginalURL == OriginalURL {
			return key, true
		}
	}
	return "", false
}
