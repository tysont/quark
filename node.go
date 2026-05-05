// ABOUTME: Node owns the blockchain, mempool, and miner identity, and
// ABOUTME: exposes the operations a peer performs on the network.
package quark

import (
	"errors"
	"sync"
)

type Node struct {
	mu       sync.Mutex
	Chain    *BlockChain
	Mempool  *Mempool
	Miner    *Miner
}

func NewNode() (*Node, error) {
	miner, err := NewMiner()
	if err != nil {
		return nil, err
	}
	return &Node{
		Chain:   NewBlockChain(),
		Mempool: NewMempool(),
		Miner:   miner,
	}, nil
}

func (n *Node) Address() string {
	return n.Miner.Wallet.Address()
}

func (n *Node) SubmitTransaction(tx *Transaction) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if tx.IsCoinbase() {
		return errors.New("coinbase transactions cannot be submitted directly")
	}
	if tx.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}
	if !tx.Verify() {
		return errors.New("transaction signature is invalid")
	}

	hash := tx.Hash()
	if n.Mempool.Has(hash) {
		return errors.New("transaction already in mempool")
	}
	if n.chainContainsTransaction(hash) {
		return errors.New("transaction already in chain")
	}

	pendingDebit := int64(0)
	for _, p := range n.Mempool.Pending() {
		if p.Sender == tx.Sender {
			pendingDebit += p.Amount
		}
	}
	available := n.Chain.Balance(tx.Sender) - pendingDebit
	if available < tx.Amount {
		return errors.New("sender has insufficient balance")
	}

	return n.Mempool.Add(tx)
}

func (n *Node) Mine() (*Block, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	pending := n.Mempool.Pending()
	block, err := n.Miner.Mine(n.Chain, pending)
	if err != nil {
		return nil, err
	}
	hashes := make([]string, 0, len(block.Data))
	for _, tx := range block.Data {
		hashes = append(hashes, tx.Hash())
	}
	n.Mempool.Remove(hashes...)
	return block, nil
}

func (n *Node) ReceiveBlock(block *Block) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if err := n.Chain.Append(block); err != nil {
		return err
	}
	hashes := make([]string, 0, len(block.Data))
	for _, tx := range block.Data {
		hashes = append(hashes, tx.Hash())
	}
	n.Mempool.Remove(hashes...)
	return nil
}

func (n *Node) Balance(address string) int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.Chain.Balance(address)
}

func (n *Node) chainContainsTransaction(hash string) bool {
	for _, block := range n.Chain.Blocks {
		for _, tx := range block.Data {
			if tx.Hash() == hash {
				return true
			}
		}
	}
	return false
}
