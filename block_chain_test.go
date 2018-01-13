package quark

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"strings"
)

func TestMineBlock(t *testing.T) {
	bc := NewBlockChain()
	data := []byte("foobar")

	d := int32(7)
	m := ""
	for i := 0; i < int(d); i++ {
		m = m + "0"
	}

	for i := 0; i < 3; i++ {
		block := bc.Mine(d, data)
		assert.NotNil(t, block)
		s := string(block.Header.Hash)
		fmt.Println(s)
		assert.True(t, strings.HasPrefix(s, m))
	}
}