// ABOUTME: Mining utilities for SHA-256 block hashing, proof-of-work
// ABOUTME: difficulty checking, and the core mining loop.
package quark

import (
	"encoding/hex"
	"encoding/binary"
	"bytes"
	"crypto/sha256"
	"math/big"
	"math/rand"
)

func hash(previousHash string, nonce int64, difficulty int32, timestamp int64, data []*Transaction) (string, error) {
	p, err := hex.DecodeString(previousHash)
	if err != nil {
		return "", err
	}

	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(timestamp))

	n := make([]byte, 8)
	binary.LittleEndian.PutUint64(n, uint64(nonce))

	d := make([]byte, 8)
	binary.LittleEndian.PutUint64(d, uint64(difficulty))

	ts, err := encode(data)
	if err != nil {
		return "", err
	}

	b := bytes.Join([][]byte{p, t, n, d, ts}, []byte{})
	h := sha256.Sum256(b)
	hash := hex.EncodeToString(h[:])
	return hash, nil
}

func check(hash string, difficulty int32) (bool, error) {
	h, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	var x big.Int
	x.SetBytes(h)
	y := big.NewInt(1)
	y.Lsh(y, uint(256 - difficulty))

	if x.Cmp(y) >= 0 {
		return false, nil
	}

	return true, nil
}

func mine(bc *BlockChain, difficulty int32, data []*Transaction) *Block {
	b := &Block{
		Data: data,
	}

	for true {
		nonce := rand.Int63()
		var bh *BlockHeader
		if bc.Length() == 0 {
			bh, _ = NewGenesisBlockHeader(nonce, difficulty, data)
		} else {
			bh, _ = NewBlockHeader(bc.Last().Header.Hash, nonce, difficulty, data)
		}

		if bh.IsValid(data) {
			b.Header = bh
			bc.Blocks = append(bc.Blocks, b)
			break
		}
	}

	return b
}
