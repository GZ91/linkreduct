package genrunes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringRunes(t *testing.T) {
	genrun := New()
	id := genrun.RandStringRunes(10)
	assert.Equal(t, 10, len(id))
}
