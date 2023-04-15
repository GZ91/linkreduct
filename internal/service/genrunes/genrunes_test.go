package genrunes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	id := RandStringRunes(10)
	assert.Equal(t, 10, len(id))
}
