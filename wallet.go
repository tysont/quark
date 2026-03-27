package quark

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

type Wallet struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := generateKeypair()
	if err != nil {
		return nil, err
	}
	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (w *Wallet) Address() string {
	b, _ := encode(w.publicKey)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func (w *Wallet) PublicKey() *rsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PrivateKey() *rsa.PrivateKey {
	return w.privateKey
}
