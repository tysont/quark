package quark

import (
	"crypto/rsa"
	"crypto/rand"
)

func generateKeypair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	reader := rand.Reader
	privateKey, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}
