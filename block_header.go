package quark

import (
	"time"
	"encoding/binary"
	"bytes"
	"crypto/sha256"
)

type BlockHeader struct {
	PreviousHash []byte
	Timestamp []byte
	Nonce []byte
	Hash []byte
}

func NewGenesisBlockHeader(nonce int64, data []byte) *BlockHeader {
	return NewBlockHeader(make([]byte, 0), nonce, data)
}

func NewBlockHeader(previousHash []byte, nonce int64, data []byte) *BlockHeader {
	u := time.Now().Unix()
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(u))

	n := make([]byte, 8)
	binary.LittleEndian.PutUint64(n, uint64(nonce))
	h := Hash(previousHash, t, n, data)

	return &BlockHeader{
		PreviousHash: previousHash,
		Timestamp: t,
		Nonce: n,
		Hash: h,
	}
}

func (bh *BlockHeader) TimestampAsUnixInt() int64 {
	return int64(binary.LittleEndian.Uint64(bh.Timestamp))
}

func (bh *BlockHeader) TimestampAsLocalTime() time.Time {
	return time.Unix(bh.TimestampAsUnixInt(), 0)
}

func (bh *BlockHeader) NonceAsInt() int64 {
	return int64(binary.LittleEndian.Uint64(bh.Nonce))
}

func (bh *BlockHeader) IsValid(data []byte) bool {
	h := Hash(bh.PreviousHash, bh.Timestamp, bh.Nonce, data)
	return bytes.Compare(bh.Hash, h) == 0
}

func Hash(previousHash []byte, timestamp []byte, nonce []byte, data []byte) []byte {
	b := bytes.Join([][]byte{previousHash, timestamp, nonce, data}, []byte{})
	h := sha256.Sum256(b)
	return h[:]
}