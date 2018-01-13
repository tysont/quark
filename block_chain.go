package quark

import (
	"math/rand"
)

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain{
	return &BlockChain{}
}

func (bc *BlockChain) Mine(difficulty int32, data []byte) *Block {
	b := &Block{
		Data: data,
	}

	for true {
		nonce := rand.Int63()
		var bh *BlockHeader
		if bc.Length() == 0 {
			bh = NewGenesisBlockHeader(nonce, difficulty, data)
		} else {
			bh = NewBlockHeader(bc.Last().Header.Hash, nonce, difficulty, data)
		}

		if bh.IsValid(data) {
			b.Header = bh
			bc.Blocks = append(bc.Blocks, b)
			break
		}
	}

	return b
}

func (bc *BlockChain) Length() int {
	return len(bc.Blocks)
}

func (bc *BlockChain) Last() *Block {
	return bc.Blocks[bc.Length() - 1]
}