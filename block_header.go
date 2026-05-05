// ABOUTME: BlockHeader contains proof-of-work metadata and links
// ABOUTME: blocks together via cryptographic hash chaining.
package quark

import (
	"crypto/sha256"
	"encoding/hex"
)

type BlockHeader struct {
	PreviousHash string
	MerkleRoot   string
	Timestamp    int64
	Nonce        int64
	Difficulty   int32
	Hash         string
}

func (bh *BlockHeader) computeHash() string {
	h := sha256.New()
	writeString(h, bh.PreviousHash)
	writeString(h, bh.MerkleRoot)
	writeInt64(h, bh.Timestamp)
	writeInt64(h, bh.Nonce)
	writeInt32(h, bh.Difficulty)
	return hex.EncodeToString(h.Sum(nil))
}

func (bh *BlockHeader) IsValid() bool {
	if bh.computeHash() != bh.Hash {
		return false
	}
	ok, err := meetsDifficulty(bh.Hash, bh.Difficulty)
	if err != nil || !ok {
		return false
	}
	return true
}
