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
}

func TestSignTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w.PrivateKey())
	assert.NoError(t, err)
	assert.NotNil(t, tx.Signature)
	assert.NotEmpty(t, tx.Signature)
}

func TestVerifySignedTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w.PrivateKey())
	assert.NoError(t, err)
	assert.True(t, tx.Verify(w.PublicKey()))
}

func TestVerifyTamperedTransaction(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 50)
	err = tx.Sign(w.PrivateKey())
	assert.NoError(t, err)

	tx.Amount = 999
	assert.False(t, tx.Verify(w.PublicKey()))
}

func TestVerifyWrongKeyTransaction(t *testing.T) {
	w1, err := NewWallet()
	assert.NoError(t, err)
	w2, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w1.Address(), "recipient", 50)
	err = tx.Sign(w1.PrivateKey())
	assert.NoError(t, err)
	assert.False(t, tx.Verify(w2.PublicKey()))
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
	assert.True(t, tx.Verify(nil))
}
