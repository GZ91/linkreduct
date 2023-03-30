package database

import (
	"github.com/GZ91/linkreduct/internal/config"
	genesis_runes "github.com/GZ91/linkreduct/internal/genesis-runes"
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

func (r db) SetDB(key, value string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[key] = value
	return true
}

func (r db) GetDB(key string) (val string, ok bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	val, ok = r.data[key]
	return
}

func GetUrl(id string) (val string, found bool) {
	val, found = DB.GetDB(id)
	return
}

func AddUrl(url string, config config.Config) string {
	lenId := 5
	iterLen := 0
	MaxIterLen := config.GetMaxIterLen()
	for {
		if iterLen == MaxIterLen {
			lenId++
		}
		idString := genesis_runes.RandStringRunes(lenId)
		if _, found := DB.GetDB(idString); found {
			iterLen++
			continue
		}
		DB.SetDB(idString, url)
		return idString
	}
}
