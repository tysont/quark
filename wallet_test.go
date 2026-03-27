package quark

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.NotNil(t, w.PrivateKey())
	assert.NotNil(t, w.PublicKey())
	assert.NotEmpty(t, w.Address())
}

func TestWalletAddressDeterministic(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)
	assert.Equal(t, w.Address(), w.Address())
}

func TestWalletAddressUniqueness(t *testing.T) {
	w1, err := NewWallet()
	assert.NoError(t, err)
	w2, err := NewWallet()
	assert.NoError(t, err)
	assert.NotEqual(t, w1.Address(), w2.Address())
}

func TestWalletAddressIsValidHex(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)
	_, err = hex.DecodeString(w.Address())
	assert.NoError(t, err)
	assert.Len(t, w.Address(), 64) // SHA-256 produces 32 bytes = 64 hex chars
}
