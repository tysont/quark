package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockChainHasGenesis(t *testing.T) {
	bc := NewBlockChain()
	assert.Equal(t, 1, bc.Length())
	assert.Equal(t, GenesisPreviousHash, bc.Blocks[0].Header.PreviousHash)
	assert.True(t, bc.IsValid())
}

func TestBalanceAfterMining(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)
	assert.Equal(t, MiningReward, bc.Balance(m.Wallet.Address()))
}

func TestBalanceAfterMultipleMines(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	n := 3
	for i := 0; i < n; i++ {
		_, err = m.Mine(bc, 8, nil)
		assert.NoError(t, err)
	}
	assert.Equal(t, MiningReward*int64(n), bc.Balance(m.Wallet.Address()))
}

func TestBalanceAfterTransfer(t *testing.T) {
	sender, err := NewMiner()
	assert.NoError(t, err)
	recipient, err := NewWallet()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = sender.Mine(bc, 8, nil)
	assert.NoError(t, err)

	tx := NewTransaction(sender.Wallet.Address(), recipient.Address(), 30)
	err = tx.Sign(sender.Wallet)
	assert.NoError(t, err)

	_, err = sender.Mine(bc, 8, []*Transaction{tx})
	assert.NoError(t, err)

	assert.Equal(t, int64(70), bc.Balance(sender.Wallet.Address()))
	assert.Equal(t, int64(30), bc.Balance(recipient.Address()))
}

func TestBalanceOfUnknownAddress(t *testing.T) {
	bc := NewBlockChain()
	assert.Equal(t, int64(0), bc.Balance("nonexistent"))
}

func TestValidateChainWithBlocks(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	for i := 0; i < 3; i++ {
		_, err = m.Mine(bc, 8, nil)
		assert.NoError(t, err)
	}
	assert.True(t, bc.IsValid())
}

func TestValidateChainWithTamperedTransaction(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)
	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)

	bc.Blocks[1].Data[0].Amount = 999
	assert.False(t, bc.IsValid())
}

func TestValidateChainBrokenLink(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)
	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)

	bc.Blocks[2].Header.PreviousHash = "0000000000000000000000000000000000000000000000000000000000000000"
	assert.False(t, bc.IsValid())
}

func TestRejectBlockWithoutCoinbase(t *testing.T) {
	bc := NewBlockChain()
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 10)
	err = tx.Sign(w)
	assert.NoError(t, err)

	header := mineHeader(bc.Last().Header.Hash, []*Transaction{tx}, 8)
	block := &Block{Header: header, Data: []*Transaction{tx}}
	assert.Error(t, bc.Append(block))
}

func TestRejectBlockWithTwoCoinbases(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	c1 := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	c2 := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	txs := []*Transaction{c1, c2}
	header := mineHeader(bc.Last().Header.Hash, txs, 8)
	block := &Block{Header: header, Data: txs}
	assert.Error(t, bc.Append(block))
}

func TestRejectBlockWithOversizedCoinbase(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	cb := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward+1)
	txs := []*Transaction{cb}
	header := mineHeader(bc.Last().Header.Hash, txs, 8)
	block := &Block{Header: header, Data: txs}
	assert.Error(t, bc.Append(block))
}

func TestRejectInsufficientBalance(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	tx := NewTransaction(m.Wallet.Address(), "recipient", 100)
	err = tx.Sign(m.Wallet)
	assert.NoError(t, err)

	_, err = m.Mine(bc, 8, []*Transaction{tx})
	assert.Error(t, err)
}

func TestRejectDuplicateTransaction(t *testing.T) {
	sender, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()
	_, err = sender.Mine(bc, 8, nil)
	assert.NoError(t, err)

	tx := NewTransaction(sender.Wallet.Address(), "recipient", 10)
	err = tx.Sign(sender.Wallet)
	assert.NoError(t, err)

	_, err = sender.Mine(bc, 8, []*Transaction{tx})
	assert.NoError(t, err)

	_, err = sender.Mine(bc, 8, []*Transaction{tx})
	assert.Error(t, err)
}

func TestRejectInvalidSignature(t *testing.T) {
	sender, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()
	_, err = sender.Mine(bc, 8, nil)
	assert.NoError(t, err)

	tx := NewTransaction(sender.Wallet.Address(), "recipient", 10)
	err = tx.Sign(sender.Wallet)
	assert.NoError(t, err)
	tx.Signature[0] ^= 0xFF

	_, err = sender.Mine(bc, 8, []*Transaction{tx})
	assert.Error(t, err)
}

func TestRejectTamperedMerkleRoot(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()
	_, err = m.Mine(bc, 8, nil)
	assert.NoError(t, err)

	bc.Blocks[1].Header.MerkleRoot = "0000000000000000000000000000000000000000000000000000000000000000"
	assert.False(t, bc.IsValid())
}

func TestEndToEnd(t *testing.T) {
	minerA, err := NewMiner()
	assert.NoError(t, err)
	minerB, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = minerA.Mine(bc, 8, nil)
	assert.NoError(t, err)
	_, err = minerA.Mine(bc, 8, nil)
	assert.NoError(t, err)
	assert.Equal(t, MiningReward*2, bc.Balance(minerA.Wallet.Address()))

	tx := NewTransaction(minerA.Wallet.Address(), minerB.Wallet.Address(), 30)
	err = tx.Sign(minerA.Wallet)
	assert.NoError(t, err)
	assert.True(t, tx.Verify())

	_, err = minerB.Mine(bc, 8, []*Transaction{tx})
	assert.NoError(t, err)

	assert.Equal(t, int64(70), bc.Balance(minerA.Wallet.Address()))
	assert.Equal(t, int64(80), bc.Balance(minerB.Wallet.Address()))
	assert.Equal(t, 4, bc.Length()) // genesis + 3 mined
	assert.True(t, bc.IsValid())
}
