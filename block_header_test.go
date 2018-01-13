package quark

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateBlockHeader(t *testing.T) {
	tx := &Transaction{}
	data := make([]*Transaction, 0)
	data = append(data, tx)
	bh, err := NewGenesisBlockHeader(0, 0, data)
	assert.NoError(t, err)
	assert.True(t, bh.IsValid(data))
	assert.False(t, bh.IsValid(make([]*Transaction, 0)))
}