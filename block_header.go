package quark

import (
	"time"
	"strings"
)

type BlockHeader struct {
	PreviousHash string
	Nonce int64
	Difficulty int32
	Timestamp int64
	Hash string
}

func NewGenesisBlockHeader(nonce int64, difficulty int32, data []*Transaction) (*BlockHeader, error) {
	return NewBlockHeader("", nonce, difficulty, data)
}

func NewBlockHeader(previousHash string, nonce int64, difficulty int32, data []*Transaction) (*BlockHeader, error) {
	timestamp := time.Now().Unix()
	hash, err := hash(previousHash, nonce, difficulty, timestamp, data)
	if err != nil {
		return nil, err
	}

	bh := &BlockHeader{
		PreviousHash: previousHash,
		Nonce: nonce,
		Difficulty: difficulty,
		Timestamp: timestamp,
		Hash: hash,
	}
	return bh, nil
}

func (bh *BlockHeader) IsValid(data []*Transaction) bool {
	hash, err := hash(bh.PreviousHash, bh.Nonce, bh.Difficulty, bh.Timestamp, data)
	if err != nil {
		return false
	}
	if strings.Compare(bh.Hash, hash) != 0 {
		return false
	}
	c, err := check(bh.Hash, bh.Difficulty)
	if err != nil {
		return false
	}
	if !c {
		return false
	}
	return true
}