package quark

type Block struct {
	Header *BlockHeader
	Data []*Transaction
}

func NewBlock(header *BlockHeader, data []*Transaction) *Block {
	return &Block{
		Header: header,
		Data: data,
	}
}