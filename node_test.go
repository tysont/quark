package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeMineMovesCoinbaseIntoChain(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)

	_, err = n.Mine(8)
	assert.NoError(t, err)
	assert.Equal(t, MiningReward, n.Balance(n.Address()))
	assert.Equal(t, 0, n.Mempool.Len())
}

func TestNodeSubmitAndMineFlow(t *testing.T) {
	sender, err := NewNode()
	assert.NoError(t, err)
	receiver, err := NewWallet()
	assert.NoError(t, err)

	_, err = sender.Mine(8)
	assert.NoError(t, err)

	tx := NewTransaction(sender.Address(), receiver.Address(), 20)
	assert.NoError(t, tx.Sign(sender.Miner.Wallet))
	assert.NoError(t, sender.SubmitTransaction(tx))
	assert.Equal(t, 1, sender.Mempool.Len())

	_, err = sender.Mine(8)
	assert.NoError(t, err)
	assert.Equal(t, 0, sender.Mempool.Len())
	assert.Equal(t, int64(80), sender.Balance(sender.Address()))
	assert.Equal(t, int64(20), sender.Balance(receiver.Address()))
}

func TestNodeRejectsBadSignature(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine(8)
	assert.NoError(t, err)

	tx := NewTransaction(n.Address(), "recipient", 5)
	assert.NoError(t, tx.Sign(n.Miner.Wallet))
	tx.Signature[0] ^= 0xFF
	assert.Error(t, n.SubmitTransaction(tx))
}

func TestNodeRejectsInsufficientBalance(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine(8)
	assert.NoError(t, err)

	tx := NewTransaction(n.Address(), "recipient", 999)
	assert.NoError(t, tx.Sign(n.Miner.Wallet))
	assert.Error(t, n.SubmitTransaction(tx))
}

func TestNodeRejectsInsufficientBalanceConsideringMempool(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine(8) // 50
	assert.NoError(t, err)

	tx1 := NewTransaction(n.Address(), "a", 30)
	tx1.Nonce = 1
	assert.NoError(t, tx1.Sign(n.Miner.Wallet))
	assert.NoError(t, n.SubmitTransaction(tx1))

	tx2 := NewTransaction(n.Address(), "b", 30)
	tx2.Nonce = 2
	assert.NoError(t, tx2.Sign(n.Miner.Wallet))
	assert.Error(t, n.SubmitTransaction(tx2))
}

func TestNodeRejectsDuplicateSubmission(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)
	_, err = n.Mine(8)
	assert.NoError(t, err)

	tx := NewTransaction(n.Address(), "recipient", 10)
	assert.NoError(t, tx.Sign(n.Miner.Wallet))
	assert.NoError(t, n.SubmitTransaction(tx))
	assert.Error(t, n.SubmitTransaction(tx))
}

func TestNodeRejectsCoinbaseSubmission(t *testing.T) {
	n, err := NewNode()
	assert.NoError(t, err)
	tx := NewCoinbaseTransaction("recipient", MiningReward)
	assert.Error(t, n.SubmitTransaction(tx))
}

func TestNodeReceiveBlockFromPeer(t *testing.T) {
	a, err := NewNode()
	assert.NoError(t, err)
	b, err := NewNode()
	assert.NoError(t, err)

	block, err := a.Mine(8)
	assert.NoError(t, err)
	assert.NoError(t, b.ReceiveBlock(block))
	assert.Equal(t, 2, b.Chain.Length())
	assert.Equal(t, MiningReward, b.Balance(a.Address()))
}

func TestNodeReceiveBlockRejectsInvalid(t *testing.T) {
	a, err := NewNode()
	assert.NoError(t, err)
	b, err := NewNode()
	assert.NoError(t, err)

	block, err := a.Mine(8)
	assert.NoError(t, err)
	block.Header.Hash = "0000000000000000000000000000000000000000000000000000000000000000"
	assert.Error(t, b.ReceiveBlock(block))
}
