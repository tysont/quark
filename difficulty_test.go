package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifficultyConstantWhenRetargetDisabled(t *testing.T) {
	bc := NewBlockChain()
	for i := 1; i <= 5; i++ {
		assert.Equal(t, bc.Config.InitialDifficulty, bc.DifficultyAt(i))
	}
}

func TestAdjustDifficultyTooFastIncreases(t *testing.T) {
	target := int64(100)
	actual := int64(25) // 4x faster than target
	got := adjustDifficulty(8, target, actual, 4)
	assert.Equal(t, int32(10), got) // ratio 4 → +2 bits
}

func TestAdjustDifficultyTooSlowDecreases(t *testing.T) {
	target := int64(100)
	actual := int64(400) // 4x slower
	got := adjustDifficulty(8, target, actual, 4)
	assert.Equal(t, int32(6), got)
}

func TestAdjustDifficultyClampsAtMaxFactor(t *testing.T) {
	target := int64(100)
	actual := int64(10000) // way slower, should clamp
	got := adjustDifficulty(8, target, actual, 4)
	assert.Equal(t, int32(6), got)
}

func TestAdjustDifficultyFloorAtOne(t *testing.T) {
	got := adjustDifficulty(1, 100, 1000, 4)
	assert.Equal(t, int32(1), got)
}

func TestRetargetIncreasesDifficultyOnFastBlocks(t *testing.T) {
	cfg := &DifficultyConfig{
		InitialDifficulty: 8,
		RetargetInterval:  5,
		TargetBlockTime:   60,
		MaxAdjustFactor:   4,
	}
	bc := NewBlockChainWithConfig(cfg)
	m, err := NewMiner()
	assert.NoError(t, err)

	for i := 1; i <= 5; i++ {
		_, err := m.mineAt(bc, nil, int64(i))
		assert.NoError(t, err)
	}
	assert.Greater(t, bc.DifficultyAt(6), int32(8))
}

func TestRetargetDecreasesDifficultyOnSlowBlocks(t *testing.T) {
	cfg := &DifficultyConfig{
		InitialDifficulty: 8,
		RetargetInterval:  5,
		TargetBlockTime:   1,
		MaxAdjustFactor:   4,
	}
	bc := NewBlockChainWithConfig(cfg)
	m, err := NewMiner()
	assert.NoError(t, err)

	for i := 1; i <= 5; i++ {
		_, err := m.mineAt(bc, nil, int64(i)*100)
		assert.NoError(t, err)
	}
	assert.Less(t, bc.DifficultyAt(6), int32(8))
}

func TestRejectBlockWithWrongDifficulty(t *testing.T) {
	bc := NewBlockChain()
	m, err := NewMiner()
	assert.NoError(t, err)

	coinbase := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	coinbase.Nonce = int64(bc.Length())
	txs := []*Transaction{coinbase}

	header := mineHeader(bc.Last().Header.Hash, txs, 12, 1) // wrong: should be 8
	block := &Block{Header: header, Data: txs}
	assert.Error(t, bc.Append(block))
}
