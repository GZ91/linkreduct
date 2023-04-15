package storage

import (
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"sync"
)

type ConfigerStorage interface {
	GetMaxIterLen() int
}

func New(conf ConfigerStorage) *db {
	return &db{data: make(map[string]string, 1), config: conf}
}

type db struct {
	data   map[string]string
	config ConfigerStorage
	mutex  sync.Mutex
}

func (r *db) setDB(key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[key] = value
	return true
}

func (r *db) GetURL(key string) (val string, ok bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	val, ok = r.data[key]
	return
}

func (r *db) AddURL(url string) string {
	lenID := 5
	iterLen := 0
	MaxIterLen := r.config.GetMaxIterLen()
	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := genrunes.RandStringRunes(lenID)
		if _, found := r.GetURL(idString); found {
			iterLen++
			continue
		}
		r.setDB(idString, url)
		return idString
	}
}
