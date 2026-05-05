package quark

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoadEmptyNode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	n, err := NewNode()
	assert.NoError(t, err)
	assert.NoError(t, n.Save(path))

	loaded, err := LoadNode(path)
	assert.NoError(t, err)
	assert.Equal(t, n.Address(), loaded.Address())
	assert.Equal(t, n.Chain.Length(), loaded.Chain.Length())
	assert.True(t, loaded.Chain.IsValid())
}

func TestSaveAndLoadAfterMining(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	n, err := NewNode()
	assert.NoError(t, err)
	for i := 0; i < 3; i++ {
		_, err := n.Mine()
		assert.NoError(t, err)
	}
	assert.NoError(t, n.Save(path))

	loaded, err := LoadNode(path)
	assert.NoError(t, err)
	assert.Equal(t, 4, loaded.Chain.Length())
	assert.Equal(t, n.Balance(n.Address()), loaded.Balance(n.Address()))
	assert.True(t, loaded.Chain.IsValid())
}

func TestSaveAndLoadPreservesMempool(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine()
	assert.NoError(t, err)

	tx := NewTransaction(n.Address(), "recipient", 10)
	tx.Nonce = 1
	assert.NoError(t, tx.Sign(n.Miner.Wallet))
	assert.NoError(t, n.SubmitTransaction(tx))
	assert.NoError(t, n.Save(path))

	loaded, err := LoadNode(path)
	assert.NoError(t, err)
	assert.Equal(t, 1, loaded.Mempool.Len())
	assert.True(t, loaded.Mempool.Has(tx.Hash()))
}

func TestLoadedWalletStillSigns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine()
	assert.NoError(t, err)
	assert.NoError(t, n.Save(path))

	loaded, err := LoadNode(path)
	assert.NoError(t, err)

	tx := NewTransaction(loaded.Address(), "recipient", 5)
	assert.NoError(t, tx.Sign(loaded.Miner.Wallet))
	assert.True(t, tx.Verify())
}

func TestLoadOrCreateCreatesWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	n, err := LoadOrCreateNode(path)
	assert.NoError(t, err)
	assert.NotNil(t, n)
	assert.Equal(t, 1, n.Chain.Length())
}

func TestLoadOrCreateReadsExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	original, err := NewNode()
	assert.NoError(t, err)
	_, err = original.Mine()
	assert.NoError(t, err)
	assert.NoError(t, original.Save(path))

	loaded, err := LoadOrCreateNode(path)
	assert.NoError(t, err)
	assert.Equal(t, original.Address(), loaded.Address())
}
