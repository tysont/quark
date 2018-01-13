package quark

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain{
	return &BlockChain{}
}

func (bc *BlockChain) Length() int {
	return len(bc.Blocks)
}

func (bc *BlockChain) Last() *Block {
	return bc.Blocks[bc.Length() - 1]
}