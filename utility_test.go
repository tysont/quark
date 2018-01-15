package quark

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestEncodeDecode(t *testing.T) {
	s := "siamese dreams"
	o := struct {
		Payload string
	}{
		Payload: s,
	}
	p := &o

	b, err := encode(p)
	assert.NoError(t, err)
	d, err := decode(reflect.TypeOf(p), b)
	assert.NoError(t, err)
	assert.Equal(t, p, d)
}

func TestMineBlocks(t *testing.T) {
	bc := NewBlockChain()
	tx := &Transaction{}
	data := make([]*Transaction, 0)
	data = append(data, tx)
	d := int32(12)
	m := ""
	for i := 0; i < int(d / 4); i++ {
		m = m + "0"
	}

	size := 4
	for i := 0; i < size; i++ {
		block := mine(bc, d, data)
		assert.NotNil(t, block)
		assert.True(t, block.Header.IsValid(data))
		assert.True(t, strings.HasPrefix(block.Header.Hash, m))
		if i > 0 {
			assert.Equal(t, block.Header.PreviousHash, bc.Blocks[i - 1].Header.Hash)
		}
	}
}