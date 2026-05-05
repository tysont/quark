// ABOUTME: Mempool holds pending transactions in submission order
// ABOUTME: until a miner includes them in a block.
package quark

import (
	"errors"
	"sync"
)

type Mempool struct {
	mu      sync.Mutex
	order   []string
	byHash  map[string]*Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		byHash: map[string]*Transaction{},
	}
}

func (m *Mempool) Add(tx *Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hash := tx.Hash()
	if _, exists := m.byHash[hash]; exists {
		return errors.New("transaction already in mempool")
	}
	m.byHash[hash] = tx
	m.order = append(m.order, hash)
	return nil
}

func (m *Mempool) Has(hash string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.byHash[hash]
	return ok
}

func (m *Mempool) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.byHash)
}

func (m *Mempool) Pending() []*Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Transaction, 0, len(m.order))
	for _, h := range m.order {
		out = append(out, m.byHash[h])
	}
	return out
}

func (m *Mempool) Remove(hashes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	dropped := map[string]bool{}
	for _, h := range hashes {
		if _, ok := m.byHash[h]; ok {
			delete(m.byHash, h)
			dropped[h] = true
		}
	}
	if len(dropped) == 0 {
		return
	}
	next := make([]string, 0, len(m.order)-len(dropped))
	for _, h := range m.order {
		if !dropped[h] {
			next = append(next, h)
		}
	}
	m.order = next
}
