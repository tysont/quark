// ABOUTME: Node persistence to and from a JSON file, capturing the
// ABOUTME: chain, mempool, and miner identity for restart recovery.
package quark

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"os"
)

type nodeFile struct {
	PrivateKey []byte         `json:"private_key"`
	Chain      *BlockChain    `json:"chain"`
	Mempool    []*Transaction `json:"mempool"`
}

func (n *Node) Save(path string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	keyBytes := x509.MarshalPKCS1PrivateKey(n.Miner.Wallet.privateKey)
	file := &nodeFile{
		PrivateKey: keyBytes,
		Chain:      n.Chain,
		Mempool:    n.Mempool.Pending(),
	}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func LoadNode(path string) (*Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var file nodeFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if file.Chain == nil {
		return nil, errors.New("loaded node has no chain")
	}
	if file.Chain.Config == nil {
		file.Chain.Config = DefaultDifficultyConfig()
	}

	priv, err := x509.ParsePKCS1PrivateKey(file.PrivateKey)
	if err != nil {
		return nil, err
	}
	wallet := &Wallet{
		privateKey: priv,
		publicKey:  &priv.PublicKey,
	}

	mempool := NewMempool()
	for _, tx := range file.Mempool {
		if err := mempool.Add(tx); err != nil {
			return nil, err
		}
	}

	node := &Node{
		Chain:   file.Chain,
		Mempool: mempool,
		Miner:   &Miner{Wallet: wallet},
	}
	return node, nil
}

func LoadOrCreateNode(path string) (*Node, error) {
	_, err := os.Stat(path)
	if err == nil {
		return LoadNode(path)
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return NewNode()
}
