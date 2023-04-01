package storage

import (
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/genrunes"
	"sync"
)

func init() {
	DB = &db{
		data: make(map[string]string, 1),
	}
}

var DB *db

type db struct {
	data  map[string]string
	mutex sync.Mutex
}

func (r *db) SetDB(key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[key] = value
	return true
}

func (r *db) GetDB(key string) (val string, ok bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	val, ok = r.data[key]
	return
}

func GetURL(id string) (val string, found bool) {
	val, found = DB.GetDB(id)
	return
}

func AddURL(url string, config *config.Config) string {
	lenID := 5
	iterLen := 0
	MaxIterLen := config.GetMaxIterLen()
	for {
		if iterLen == MaxIterLen {
			lenID++
		}
		idString := genrunes.RandStringRunes(lenID)
		if _, found := DB.GetDB(idString); found {
			iterLen++
			continue
		}
		DB.SetDB(idString, url)
		return idString
	}
}
