package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockHeaderHashIsDeterministic(t *testing.T) {
	bh := &BlockHeader{
		PreviousHash: "abc",
		MerkleRoot:   "def",
		Timestamp:    1,
		Nonce:        2,
		Difficulty:   3,
	}
	assert.Equal(t, bh.computeHash(), bh.computeHash())
}

func TestBlockHeaderIsValidRequiresMatchingHash(t *testing.T) {
	bh := mineHeader("", nil, 4, 1)
	assert.True(t, bh.IsValid())

	bh.Hash = "0000000000000000000000000000000000000000000000000000000000000000"
	assert.False(t, bh.IsValid())
}

func TestBlockHeaderIsValidRequiresDifficulty(t *testing.T) {
	bh := mineHeader("", nil, 4, 1)
	bh.Difficulty = 200
	bh.Hash = bh.computeHash()
	assert.False(t, bh.IsValid())
}
