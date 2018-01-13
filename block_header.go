package quark

import (
	"time"
	"encoding/binary"
	"bytes"
	"crypto/sha256"
	"math/big"
)

type BlockHeader struct {
	PreviousHash []byte
	Nonce int64
	Difficulty int32
	Timestamp int64
	Hash []byte
}

func NewGenesisBlockHeader(nonce int64, difficulty int32, data []byte) *BlockHeader {
	return NewBlockHeader(make([]byte, 0), nonce, difficulty, data)
}

func NewBlockHeader(previousHash []byte, nonce int64, difficulty int32, data []byte) *BlockHeader {
	timestamp := time.Now().Unix()
	hash := hash(previousHash, nonce, difficulty, timestamp, data)
	return &BlockHeader{
		PreviousHash: previousHash,
		Nonce: nonce,
		Difficulty: difficulty,
		Timestamp: timestamp,
		Hash: hash,
	}
}

func (bh *BlockHeader) IsValid(data []byte) bool {
	hash := hash(bh.PreviousHash, bh.Nonce, bh.Difficulty, bh.Timestamp, data)
	if bytes.Compare(bh.Hash, hash) != 0 {
		return false
	}

	if !check(bh.Hash, bh.Difficulty) {
		return false
	}

	return true
}

func hash(previousHash []byte, nonce int64, difficulty int32, timestamp int64, data []byte) []byte {
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(timestamp))

	n := make([]byte, 8)
	binary.LittleEndian.PutUint64(n, uint64(nonce))

	d := make([]byte, 8)
	binary.LittleEndian.PutUint64(d, uint64(difficulty))

	b := bytes.Join([][]byte{previousHash, t, n, d, data}, []byte{})
	h := sha256.Sum256(b)
	return h[:]
}

func check(hash []byte, difficulty int32) bool {
	var h big.Int
	h.SetBytes(hash[:])
	m := big.NewInt(1)
	m.Lsh(m, uint(256 - difficulty))

	if h.Cmp(m) >= 0 {
		return false
	}

	return true
}