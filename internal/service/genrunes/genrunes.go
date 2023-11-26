package genrunes

import (
	"math/rand"
	"time"
)

// Genrun представляет генератор случайных строк на основе заданных символов.
type Genrun struct {
	letterRunes []rune
	rander      *rand.Rand
}

// New создает и возвращает новый экземпляр Genrun.
func New() *Genrun {
	return &Genrun{
		letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		rander:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RandStringRunes генерирует случайную строку длины l из символов, заданных в letterRunes.
func (g Genrun) RandStringRunes(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = g.letterRunes[g.rander.Intn(len(g.letterRunes))]
	}
	return string(b)
}
