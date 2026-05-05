// ABOUTME: Wallet provides keypair-based identity for blockchain
// ABOUTME: participants, deriving addresses from public key hashes.
package quark

import (
	"crypto/rsa"
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
	addr, _ := addressFromPublicKey(w.publicKey)
	return addr
}

func (w *Wallet) PublicKey() *rsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PrivateKey() *rsa.PrivateKey {
	return w.privateKey
}
