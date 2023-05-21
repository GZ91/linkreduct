package inmemory

import (
	"context"
	"fmt"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
	"sync"
)

type ConfigerStorage interface {
	GetMaxIterLen() int
}

type GeneratorRunes interface {
	RandStringRunes(int) string
}

func New(ctx context.Context, conf ConfigerStorage, genrun GeneratorRunes) *db {
	return &db{data: make(map[string]string, 1), config: conf, genrun: genrun}
}

type db struct {
	data   map[string]string
	config ConfigerStorage
	genrun GeneratorRunes
	mutex  sync.Mutex
}

func (r *db) setDB(key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[key] = value
	return true
}

func (r *db) GetURL(ctx context.Context, key string) (val string, ok bool, errs error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	val, ok = r.data[key]
	return
}

func (r *db) AddURL(ctx context.Context, url string) (string, error) {

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
		r.setDB(idString, url)
		return idString, nil
	}
}

func (r *db) Close() error {
	return nil
}

func (r *db) Ping(ctx context.Context) error {
	return nil
}

func (r *db) FindLongURL(ctx context.Context, OriginalURL string) (string, bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, val := range r.data {
		if val == OriginalURL {
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
			errs = fmt.Errorf("%w; %w", errs, errorsapp.ErrLinkAlreadyExists)
		} else {
			var err error
			shortURL, err = r.AddURL(ctx, link)
			if err != nil {
				logger.Log.Error("error when writing an add link to the inmemory storage", zap.Error(err), zap.String("unadded value", link))
				errs = fmt.Errorf("%w; %w", errs, err)
			}
		}
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{CorrelationID: data.CorrelationID, ShortURL: shortURL})
	}
	return
}
