# quark

A minimal, feature-complete reference implementation of a proof-of-work blockchain in Go.

## Features

- **Wallets** — RSA keypair identity with SHA-256 address derivation
- **Transactions** — signed transfers with embedded sender public key for self-contained verification
- **Coinbase rewards** — miners receive a fixed block reward
- **Proof-of-work mining** — SHA-256 header hashing with target-difficulty PoW
- **Merkle root** — block headers commit to their transaction set
- **Difficulty adjustment** — Bitcoin-style retargeting based on observed block times
- **Full block validation** — signatures, balances, double-spend detection, coinbase rules, merkle root, header integrity, and chain linkage
- **Mempool** — pending transaction pool with balance-aware admission
- **Persistence** — JSON snapshot of chain, mempool, and miner identity
- **HTTP node** — peer-to-peer block and transaction propagation with longest-valid-chain sync
- **CLI** — single binary for running a node and interacting with one over HTTP

## Library usage

```go
n, _ := quark.NewNode()
n.Mine() // miner gets 50

tx := quark.NewTransaction(n.Address(), recipient, 30)
tx.Sign(n.Miner.Wallet)
n.SubmitTransaction(tx)
n.Mine() // tx is mined, recipient gets 30

n.Balance(n.Address())   // 70 (50 reward + 50 reward - 30 sent)
n.Balance(recipient)     // 30
n.Chain.IsValid()        // true
```

## Running a node

```sh
go install ./cmd/quark

# Start two peers
quark node --listen :8080 --data a.json &
quark node --listen :8081 --data b.json --peers http://localhost:8080 &

# Connect them mutually
quark peer --node http://localhost:8080 --url http://localhost:8081

# Mine and watch propagation
quark mine --node http://localhost:8080
quark balance --node http://localhost:8081 --address $(quark address --node http://localhost:8080 | jq -r .address)
```

## HTTP API

| Method | Path                      | Description                                  |
|--------|---------------------------|----------------------------------------------|
| POST   | `/tx`                     | Submit a signed transaction (JSON body)      |
| POST   | `/send`                   | Sign and submit a transaction with this node's wallet |
| POST   | `/block`                  | Receive a block from a peer                  |
| GET    | `/chain`                  | Return the full chain                        |
| GET    | `/balance?address=ADDR`   | Return balance for an address                |
| GET    | `/address`                | Return this node's miner address             |
| POST   | `/mine`                   | Mine a block and broadcast it                |
| POST   | `/peers`                  | Add a peer URL                               |
| GET    | `/peers`                  | List configured peers                        |
| POST   | `/sync`                   | Pull peers' chains, adopt the longest valid one |

## Testing

```sh
go test ./...
```

Tests cover unit-level behavior for every component, integration-level mempool and node flows, and end-to-end multi-node scenarios including block propagation, transaction broadcast, longest-chain sync, and persistence-restart.
