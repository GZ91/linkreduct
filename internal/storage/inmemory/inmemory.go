package inmemory

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/google/uuid"
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
	return &db{data: make(map[string]models.StructURL, 1), config: conf, genrun: genrun}
}

type db struct {
	data   map[string]models.StructURL
	config ConfigerStorage
	genrun GeneratorRunes
	mutex  sync.Mutex
}

func (r *db) setDB(ctx context.Context, key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}
	StructURL := models.StructURL{
		OriginalURL: value,
		ShortURL:    key,
		UserID:      UserID,
		ID:          uuid.New().String(),
	}
	r.data[key] = StructURL
	return true
}

func (r *db) GetURL(ctx context.Context, key string) (val string, ok bool, errs error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	valueStruct, found := r.data[key]
	if valueStruct.DeletedFlag {
		return "", false, errorsapp.ErrLineURLDeleted
	}
	if found {
		ok = found
		val = valueStruct.OriginalURL
	}
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
		r.setDB(ctx, idString, url)
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
				logger.Log.Error("error when writing an add link to the inmemory storage", zap.Error(err), zap.String("unadded value", link))
				errs = err
				return
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
