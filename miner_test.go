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

	block := m.Mine(bc, 8, nil)
	assert.NotNil(t, block)
	assert.Equal(t, 1, bc.Length())
}

func TestMinerCoinbaseReward(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	block := m.Mine(bc, 8, nil)
	assert.True(t, len(block.Data) >= 1)

	coinbase := block.Data[0]
	assert.True(t, coinbase.IsCoinbase())
	assert.Equal(t, m.Wallet.Address(), coinbase.Recipient)
	assert.Equal(t, MiningReward, coinbase.Amount)
}

func TestMinerIncludesTransactions(t *testing.T) {
	m, err := NewMiner()
	assert.NoError(t, err)
	bc := NewBlockChain()

	tx := NewTransaction("sender", "recipient", 10)
	block := m.Mine(bc, 8, []*Transaction{tx})

	assert.Equal(t, 2, len(block.Data)) // coinbase + user tx
	assert.True(t, block.Data[0].IsCoinbase())
	assert.Equal(t, "sender", block.Data[1].Sender)
}
