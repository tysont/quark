// ABOUTME: BlockChain is the ordered ledger of blocks with balance
// ABOUTME: computation and full chain integrity validation.
package quark

import (
	"errors"
)

const GenesisPreviousHash = ""

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*Block{newGenesisBlock()},
	}
}

func newGenesisBlock() *Block {
	header := &BlockHeader{
		PreviousHash: GenesisPreviousHash,
		MerkleRoot:   merkleRoot(nil),
		Timestamp:    0,
		Nonce:        0,
		Difficulty:   0,
	}
	header.Hash = header.computeHash()
	return &Block{Header: header, Data: nil}
}

func (bc *BlockChain) Length() int {
	return len(bc.Blocks)
}

func (bc *BlockChain) Last() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
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
	return bc.Validate() == nil
}

func (bc *BlockChain) Validate() error {
	if len(bc.Blocks) == 0 {
		return errors.New("chain has no genesis block")
	}

	expectedGenesis := newGenesisBlock()
	if bc.Blocks[0].Header.Hash != expectedGenesis.Header.Hash {
		return errors.New("genesis block does not match expected genesis")
	}

	balances := map[string]int64{}
	seenTxHashes := map[string]bool{}

	for i := 1; i < len(bc.Blocks); i++ {
		block := bc.Blocks[i]
		prev := bc.Blocks[i-1]

		if err := validateBlockStructure(block, prev); err != nil {
			return err
		}
		if err := applyBlockTransactions(block, balances, seenTxHashes); err != nil {
			return err
		}
	}
	return nil
}

func (bc *BlockChain) Append(block *Block) error {
	prev := bc.Last()
	if err := validateBlockStructure(block, prev); err != nil {
		return err
	}

	balances := map[string]int64{}
	seen := map[string]bool{}
	for i := 1; i < len(bc.Blocks); i++ {
		if err := applyBlockTransactions(bc.Blocks[i], balances, seen); err != nil {
			return err
		}
	}
	if err := applyBlockTransactions(block, balances, seen); err != nil {
		return err
	}

	bc.Blocks = append(bc.Blocks, block)
	return nil
}

func validateBlockStructure(block, prev *Block) error {
	if block.Header.PreviousHash != prev.Header.Hash {
		return errors.New("block previous hash does not link to prior block")
	}
	if block.Header.Timestamp < prev.Header.Timestamp {
		return errors.New("block timestamp is older than previous block")
	}
	if block.Header.MerkleRoot != merkleRoot(block.Data) {
		return errors.New("block merkle root does not match transactions")
	}
	if !block.Header.IsValid() {
		return errors.New("block header hash is invalid or fails difficulty")
	}
	if len(block.Data) == 0 {
		return errors.New("block has no transactions")
	}
	if !block.Data[0].IsCoinbase() {
		return errors.New("first transaction in block must be coinbase")
	}
	if block.Data[0].Amount != MiningReward {
		return errors.New("coinbase amount does not equal mining reward")
	}
	for _, tx := range block.Data[1:] {
		if tx.IsCoinbase() {
			return errors.New("block contains more than one coinbase transaction")
		}
	}
	return nil
}

func applyBlockTransactions(block *Block, balances map[string]int64, seen map[string]bool) error {
	for _, tx := range block.Data {
		if tx.Amount <= 0 {
			return errors.New("transaction amount must be positive")
		}
		hash := tx.Hash()
		if seen[hash] {
			return errors.New("duplicate transaction in chain")
		}
		seen[hash] = true

		if !tx.IsCoinbase() {
			if !tx.Verify() {
				return errors.New("transaction signature is invalid")
			}
			if balances[tx.Sender] < tx.Amount {
				return errors.New("sender has insufficient balance")
			}
			balances[tx.Sender] -= tx.Amount
		}
		balances[tx.Recipient] += tx.Amount
	}
	return nil
}
