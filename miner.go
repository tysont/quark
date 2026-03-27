// ABOUTME: Miner performs proof-of-work mining and injects coinbase
// ABOUTME: reward transactions into each mined block.
package quark

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

func (m *Miner) Mine(bc *BlockChain, difficulty int32, transactions []*Transaction) *Block {
	coinbase := NewCoinbaseTransaction(m.Wallet.Address(), MiningReward)
	all := append([]*Transaction{coinbase}, transactions...)
	return mine(bc, difficulty, all)
}
