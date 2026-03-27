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

func (bc *BlockChain) Balance(address string) int64 {
	var balance int64
	for _, block := range bc.Blocks {
		for _, tx := range block.Data {
			if tx.Recipient == address {
				balance += tx.Amount
			}
			if tx.Sender == address {
				balance -= tx.Amount
			}
		}
	}
	return balance
}

func (bc *BlockChain) IsValid() bool {
	for i, block := range bc.Blocks {
		if !block.Header.IsValid(block.Data) {
			return false
		}
		if i == 0 {
			if block.Header.PreviousHash != "" {
				return false
			}
		} else {
			if block.Header.PreviousHash != bc.Blocks[i-1].Header.Hash {
				return false
			}
		}
	}
	return true
}