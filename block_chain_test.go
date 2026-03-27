package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceAfterMining(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	m.Mine(bc, 8, nil)
	assert.Equal(t, MiningReward, bc.Balance(m.Wallet.Address()))
}

func TestBalanceAfterMultipleMines(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	n := 3
	for i := 0; i < n; i++ {
		m.Mine(bc, 8, nil)
	}
	assert.Equal(t, MiningReward*int64(n), bc.Balance(m.Wallet.Address()))
}

func TestBalanceAfterTransfer(t *testing.T) {
	sender, err := NewMiner()
	assert.NoError(t, err)
	recipient, err := NewWallet()
	assert.NoError(t, err)
	bc := NewBlockChain()

	sender.Mine(bc, 8, nil) // sender gets 50

	tx := NewTransaction(sender.Wallet.Address(), recipient.Address(), 30)
	err = tx.Sign(sender.Wallet.PrivateKey())
	assert.NoError(t, err)

	sender.Mine(bc, 8, []*Transaction{tx}) // sender gets another 50, sends 30

	assert.Equal(t, int64(70), bc.Balance(sender.Wallet.Address()))
	assert.Equal(t, int64(30), bc.Balance(recipient.Address()))
}

func TestBalanceOfUnknownAddress(t *testing.T) {
	bc := NewBlockChain()
	assert.Equal(t, int64(0), bc.Balance("nonexistent"))
}

func TestValidateEmptyChain(t *testing.T) {
	bc := NewBlockChain()
	assert.True(t, bc.IsValid())
}

func TestValidateChainWithBlocks(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	for i := 0; i < 3; i++ {
		m.Mine(bc, 8, nil)
	}
	assert.True(t, bc.IsValid())
}

func TestValidateChainWithTamperedBlock(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	m.Mine(bc, 8, nil)
	m.Mine(bc, 8, nil)

	bc.Blocks[0].Data = []*Transaction{NewTransaction("fake", "fake", 999)}
	assert.False(t, bc.IsValid())
}

func TestValidateChainBrokenLink(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	m.Mine(bc, 8, nil)
	m.Mine(bc, 8, nil)

	bc.Blocks[1].Header.PreviousHash = "0000000000000000000000000000000000000000000000000000000000000000"
	assert.False(t, bc.IsValid())
}

func TestEndToEnd(t *testing.T) {
	minerA, err := NewMiner()
	assert.NoError(t, err)
	minerB, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	// Miner A mines 2 blocks: balance = 100
	minerA.Mine(bc, 8, nil)
	minerA.Mine(bc, 8, nil)
	assert.Equal(t, MiningReward*2, bc.Balance(minerA.Wallet.Address()))

	// A sends 30 to B, mined by B: A = 70, B = 80
	tx := NewTransaction(minerA.Wallet.Address(), minerB.Wallet.Address(), 30)
	err = tx.Sign(minerA.Wallet.PrivateKey())
	assert.NoError(t, err)
	assert.True(t, tx.Verify(minerA.Wallet.PublicKey()))

	minerB.Mine(bc, 8, []*Transaction{tx})

	assert.Equal(t, int64(70), bc.Balance(minerA.Wallet.Address()))
	assert.Equal(t, int64(80), bc.Balance(minerB.Wallet.Address()))
	assert.Equal(t, 3, bc.Length())
	assert.True(t, bc.IsValid())
}
