package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("sender", "recipient", 100)
	assert.Equal(t, "sender", tx.Sender)
	assert.Equal(t, "recipient", tx.Recipient)
	assert.Equal(t, int64(100), tx.Amount)
	assert.Nil(t, tx.Signature)
	assert.Nil(t, tx.SenderPublicKey)
}

func TestSignTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w)
	assert.NoError(t, err)
	assert.NotEmpty(t, tx.Signature)
	assert.NotEmpty(t, tx.SenderPublicKey)
}

func TestVerifySignedTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w)
	assert.NoError(t, err)
	assert.True(t, tx.Verify())
}

func TestVerifyTamperedTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w)
	assert.NoError(t, err)

	tx.Amount = 999
	assert.False(t, tx.Verify())
}

func TestVerifySwappedPublicKey(t *testing.T) {
	w1, err := NewWallet()
	assert.NoError(t, err)
	w2, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w1.Address(), "recipient", 50)
	err = tx.Sign(w1)
	assert.NoError(t, err)

	otherPub, err := marshalPublicKey(w2.PublicKey())
	assert.NoError(t, err)
	tx.SenderPublicKey = otherPub

	assert.False(t, tx.Verify())
}

func TestSignWithMismatchedAddress(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction("wrong-address", "recipient", 50)
	err = tx.Sign(w)
	assert.Error(t, err)
}

func TestNewCoinbaseTransaction(t *testing.T) {
	tx := NewCoinbaseTransaction("recipient", 50)
	assert.Equal(t, "", tx.Sender)
	assert.Equal(t, "recipient", tx.Recipient)
	assert.Equal(t, int64(50), tx.Amount)
	assert.Nil(t, tx.Signature)
}

func TestCoinbaseIsCoinbase(t *testing.T) {
	tx := NewCoinbaseTransaction("recipient", 50)
	assert.True(t, tx.IsCoinbase())

	tx2 := NewTransaction("sender", "recipient", 50)
	assert.False(t, tx2.IsCoinbase())
}

func TestCoinbaseDoesNotRequireSignature(t *testing.T) {
	tx := NewCoinbaseTransaction("recipient", 50)
	assert.True(t, tx.Verify())
}

func TestTransactionHashChangesWithFields(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w)
	assert.NoError(t, err)
	h1 := tx.Hash()

	tx.Amount = 60
	assert.NotEqual(t, h1, tx.Hash())
}
