# quark

A minimal, feature-complete blockchain implementation in Go.

## Features

- **Wallets** -- RSA keypair identity with SHA-256 address derivation
- **Transactions** -- signed transfers of value between addresses with tamper detection
- **Coinbase rewards** -- miners receive a block reward for each mined block
- **Proof-of-work mining** -- SHA-256 hash-based mining with configurable difficulty
- **Balance tracking** -- account-based balance computation by walking the chain
- **Chain validation** -- verifies block header integrity and hash-chain linkage

## Usage

```go
// Create a blockchain and two miners
bc := quark.NewBlockChain()
minerA, _ := quark.NewMiner()
minerB, _ := quark.NewMiner()

// Mine a block (miner receives a coinbase reward)
minerA.Mine(bc, 12, nil)

// Create and sign a transaction
tx := quark.NewTransaction(minerA.Wallet.Address(), minerB.Wallet.Address(), 30)
tx.Sign(minerA.Wallet.PrivateKey())

// Mine a block containing the transaction
minerB.Mine(bc, 12, []*quark.Transaction{tx})

// Query balances
bc.Balance(minerA.Wallet.Address()) // 20 (50 reward - 30 sent)
bc.Balance(minerB.Wallet.Address()) // 80 (50 reward + 30 received)

// Validate the chain
bc.IsValid() // true
```

## Testing

```sh
go test ./...
```
