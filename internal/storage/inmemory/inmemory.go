package inmemory

import (
	"sync"
)

type ConfigerStorage interface {
	GetMaxIterLen() int
}

type GeneratorRunes interface {
	RandStringRunes(int) string
}

func New(conf ConfigerStorage, genrun GeneratorRunes) *db {
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
		idString := r.genrun.RandStringRunes(lenID)
		if _, found := r.GetURL(idString); found {
			iterLen++
			continue
		}
		r.setDB(idString, url)
		return idString
	}
}

func (r *db) Close() error {
	return nil
}

func (r *db) Ping() error {
	return nil
}
