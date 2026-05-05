// ABOUTME: Merkle root computation over a list of transactions, used
// ABOUTME: as the block header's commitment to its transaction set.
package quark

import (
	"crypto/sha256"
	"encoding/hex"
)

func merkleRoot(txs []*Transaction) string {
	if len(txs) == 0 {
		var zero [32]byte
		return hex.EncodeToString(zero[:])
	}

	layer := make([][]byte, len(txs))
	for i, tx := range txs {
		b, _ := hex.DecodeString(tx.Hash())
		layer[i] = b
	}

	for len(layer) > 1 {
		if len(layer)%2 == 1 {
			layer = append(layer, layer[len(layer)-1])
		}
		next := make([][]byte, 0, len(layer)/2)
		for i := 0; i < len(layer); i += 2 {
			h := sha256.New()
			h.Write(layer[i])
			h.Write(layer[i+1])
			next = append(next, h.Sum(nil))
		}
		layer = next
	}
	return hex.EncodeToString(layer[0])
}
