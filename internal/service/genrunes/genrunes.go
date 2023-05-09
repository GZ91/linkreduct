package genrunes

import (
	"math/rand"
	"time"
)

type Genrun struct {
	letterRunes []rune
	rander      *rand.Rand
}

func New() *Genrun {
	return &Genrun{
		letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		rander:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g Genrun) RandStringRunes(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = g.letterRunes[g.rander.Intn(len(g.letterRunes))]
	}
	return string(b)
}
