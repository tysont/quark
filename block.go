// ABOUTME: Block pairs a proof-of-work header with a list of
// ABOUTME: transactions to form a single unit in the chain.
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