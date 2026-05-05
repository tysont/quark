// ABOUTME: Miner performs proof-of-work mining and injects coinbase
// ABOUTME: reward transactions into each mined block.
package quark

import (
	"errors"
	"time"
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

func (m *Miner) Mine(bc *BlockChain, transactions []*Transaction) (*Block, error) {
	return m.mineAt(bc, transactions, time.Now().Unix())
}

func (m *Miner) mineAt(bc *BlockChain, transactions []*Transaction, timestamp int64) (*Block, error) {
	coinbase := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	coinbase.Nonce = int64(bc.Length())
	all := append([]*Transaction{coinbase}, transactions...)

	difficulty := bc.NextDifficulty()
	prev := bc.Last()
	if timestamp < prev.Header.Timestamp {
		timestamp = prev.Header.Timestamp
	}
	header := mineHeader(prev.Header.Hash, all, difficulty, timestamp)
	block := &Block{Header: header, Data: all}

	if err := bc.Append(block); err != nil {
		return nil, errors.Join(errors.New("mined block failed validation"), err)
	}
	return block, nil
}
