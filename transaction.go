package quark

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int64
	Signature []byte
}

func NewTransaction(sender, recipient string, amount int64) *Transaction {
	return &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}

func NewCoinbaseTransaction(recipient string, amount int64) *Transaction {
	return &Transaction{
		Recipient: recipient,
		Amount:    amount,
	}
}

func (tx *Transaction) IsCoinbase() bool {
	return tx.Sender == ""
}

func (tx *Transaction) Sign(privateKey *rsa.PrivateKey) error {
	h, err := tx.digest()
	if err != nil {
		return err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

func (tx *Transaction) Verify(publicKey *rsa.PublicKey) bool {
	if tx.IsCoinbase() {
		return true
	}
	h, err := tx.digest()
	if err != nil {
		return false
	}
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h, tx.Signature) == nil
}

func (tx *Transaction) digest() ([]byte, error) {
	b, err := encode(&struct {
		Sender    string
		Recipient string
		Amount    int64
	}{
		Sender:    tx.Sender,
		Recipient: tx.Recipient,
		Amount:    tx.Amount,
	})
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(b)
	return h[:], nil
}
