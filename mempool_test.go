package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMempoolAddAndPending(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 10)
	err = tx.Sign(w)
	assert.NoError(t, err)

	mp := NewMempool()
	assert.NoError(t, mp.Add(tx))
	assert.Equal(t, 1, mp.Len())
	assert.True(t, mp.Has(tx.Hash()))
	assert.Equal(t, []*Transaction{tx}, mp.Pending())
}

func TestMempoolRejectsDuplicate(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 10)
	err = tx.Sign(w)
	assert.NoError(t, err)

	mp := NewMempool()
	assert.NoError(t, mp.Add(tx))
	assert.Error(t, mp.Add(tx))
	assert.Equal(t, 1, mp.Len())
}

func TestMempoolRemove(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx1 := NewTransaction(w.Address(), "a", 10)
	tx1.Nonce = 1
	assert.NoError(t, tx1.Sign(w))
	tx2 := NewTransaction(w.Address(), "b", 10)
	tx2.Nonce = 2
	assert.NoError(t, tx2.Sign(w))

	mp := NewMempool()
	assert.NoError(t, mp.Add(tx1))
	assert.NoError(t, mp.Add(tx2))

	mp.Remove(tx1.Hash())
	assert.Equal(t, 1, mp.Len())
	assert.False(t, mp.Has(tx1.Hash()))
	assert.True(t, mp.Has(tx2.Hash()))
}

func TestMempoolPendingPreservesOrder(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	mp := NewMempool()
	hashes := []string{}
	for i := 0; i < 5; i++ {
		tx := NewTransaction(w.Address(), "r", 10)
		tx.Nonce = int64(i)
		assert.NoError(t, tx.Sign(w))
		assert.NoError(t, mp.Add(tx))
		hashes = append(hashes, tx.Hash())
	}

	pending := mp.Pending()
	for i, tx := range pending {
		assert.Equal(t, hashes[i], tx.Hash())
	}
}
