// ABOUTME: Mining utilities for SHA-256 proof-of-work header hashing,
// ABOUTME: difficulty checking, and the core mining loop.
package quark

import (
	"encoding/hex"
	"math/big"
	"time"
)

func meetsDifficulty(headerHash string, difficulty int32) (bool, error) {
	h, err := hex.DecodeString(headerHash)
	if err != nil {
		return false, err
	}
	var x big.Int
	x.SetBytes(h)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	return x.Cmp(target) < 0, nil
}

func mineHeader(previousHash string, txs []*Transaction, difficulty int32) *BlockHeader {
	bh := &BlockHeader{
		PreviousHash: previousHash,
		MerkleRoot:   merkleRoot(txs),
		Timestamp:    time.Now().Unix(),
		Difficulty:   difficulty,
	}
	for nonce := int64(0); ; nonce++ {
		bh.Nonce = nonce
		bh.Hash = bh.computeHash()
		ok, err := meetsDifficulty(bh.Hash, difficulty)
		if err == nil && ok {
			return bh
		}
	}
}
