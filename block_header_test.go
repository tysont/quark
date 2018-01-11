package quark

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateBlockHeader(t *testing.T) {
	data := []byte("hello, world!")
	bh := NewGenesisBlockHeader(0, data)
	assert.True(t, bh.IsValid(data))
	assert.False(t, bh.IsValid([]byte("goodbye, world!")))
}