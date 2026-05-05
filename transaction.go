// ABOUTME: Transaction represents a signed transfer of value between
// ABOUTME: wallet addresses, with coinbase support for mining rewards.
package quark

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Transaction struct {
	Sender          string
	Recipient       string
	Amount          int64
	Nonce           int64
	SenderPublicKey []byte
	Signature       []byte
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

func (tx *Transaction) Sign(w *Wallet) error {
	pubBytes, err := marshalPublicKey(w.PublicKey())
	if err != nil {
		return err
	}
	tx.SenderPublicKey = pubBytes

	addr, err := addressFromPublicKey(w.PublicKey())
	if err != nil {
		return err
	}
	if addr != tx.Sender {
		return errors.New("wallet does not match transaction sender")
	}

	digest := tx.digest()
	sig, err := rsa.SignPKCS1v15(rand.Reader, w.PrivateKey(), crypto.SHA256, digest)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.IsCoinbase() {
		return tx.SenderPublicKey == nil && tx.Signature == nil
	}
	if len(tx.SenderPublicKey) == 0 || len(tx.Signature) == 0 {
		return false
	}
	pub, err := unmarshalPublicKey(tx.SenderPublicKey)
	if err != nil {
		return false
	}
	addr, err := addressFromPublicKey(pub)
	if err != nil {
		return false
	}
	if addr != tx.Sender {
		return false
	}
	digest := tx.digest()
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest, tx.Signature) == nil
}

func (tx *Transaction) Hash() string {
	h := sha256.New()
	writeString(h, tx.Sender)
	writeString(h, tx.Recipient)
	writeInt64(h, tx.Amount)
	writeInt64(h, tx.Nonce)
	writeBytes(h, tx.SenderPublicKey)
	writeBytes(h, tx.Signature)
	return hex.EncodeToString(h.Sum(nil))
}

func (tx *Transaction) digest() []byte {
	h := sha256.New()
	writeString(h, tx.Sender)
	writeString(h, tx.Recipient)
	writeInt64(h, tx.Amount)
	writeInt64(h, tx.Nonce)
	return h.Sum(nil)
}
