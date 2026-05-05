// ABOUTME: Miner performs proof-of-work mining and injects coinbase
// ABOUTME: reward transactions into each mined block.
package quark

import (
	"errors"
)

const MiningReward int64 = 50

type Miner struct {
	Wallet *Wallet
}

func NewMiner() (*Miner, error) {
	w, err := NewWallet()
	if err != nil {
		return nil, err
	}
	return &Miner{Wallet: w}, nil
}

func (m *Miner) Mine(bc *BlockChain, difficulty int32, transactions []*Transaction) (*Block, error) {
	coinbase := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	coinbase.Nonce = int64(bc.Length())
	all := append([]*Transaction{coinbase}, transactions...)

	header := mineHeader(bc.Last().Header.Hash, all, difficulty)
	block := &Block{Header: header, Data: all}

	if err := bc.Append(block); err != nil {
		return nil, errors.Join(errors.New("mined block failed validation"), err)
	}
	return block, nil
}
