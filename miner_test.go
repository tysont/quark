package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMiner(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	assert.NotNil(t, m.Wallet)
	assert.NotEmpty(t, m.Wallet.Address())
}

func TestMinerMineBlock(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	block, err := m.Mine(bc, nil)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, 2, bc.Length())
}

func TestMinerCoinbaseReward(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	block, err := m.Mine(bc, nil)
	assert.NoError(t, err)
	assert.True(t, len(block.Data) >= 1)

	coinbase := block.Data[0]
	assert.True(t, coinbase.IsCoinbase())
	assert.Equal(t, m.Wallet.Address(), coinbase.Recipient)
	assert.Equal(t, MiningReward, coinbase.Amount)
}

func TestMinerIncludesTransactions(t *testing.T) {
	sender, err := NewMiner()
	assert.NoError(t, err)
	miner, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	_, err = sender.Mine(bc, nil) // fund sender
	assert.NoError(t, err)

	tx := NewTransaction(sender.Wallet.Address(), "recipient", 10)
	err = tx.Sign(sender.Wallet)
	assert.NoError(t, err)

	block, err := miner.Mine(bc, []*Transaction{tx})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(block.Data))
	assert.True(t, block.Data[0].IsCoinbase())
	assert.Equal(t, sender.Wallet.Address(), block.Data[1].Sender)
}
