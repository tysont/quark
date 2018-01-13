package quark

type Block struct {
	Header *BlockHeader
	Data []byte
}

func NewBlock(header *BlockHeader, data []byte) *Block {
	return &Block{
		Header: header,
		Data: data,
	}
}