package quark

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
	"fmt"
	"strings"
)

func TestEncodeDecode(t *testing.T) {
	s := "siamese dreams"
	o := struct {
		Payload string
	}{
		Payload: s,
	}

	b, err := encode(o)
	assert.NoError(t, err)
	d, err := decode(reflect.TypeOf(o), b)
	assert.NoError(t, err)
	assert.Equal(t, o, d)
}

func TestMineBlock(t *testing.T) {
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
		s := string(block.Header.Hash)
		fmt.Println(s)
		assert.True(t, strings.HasPrefix(s, m))
	}
}