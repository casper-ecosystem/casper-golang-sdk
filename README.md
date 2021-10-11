# Casper Golang SDK

Golang library for interacting with CSPR node and utilities for Casper ecosystem.

## Features
- RPC Client
- Event Stream Client
- Deploy util
- KeyPair util

## Needs improving
- Casper Types (serialization)
- KeyPair util
- Better testing
- Better serialization

## Documentation
Documentation generated from the code can be found here: https://pkg.go.dev/github.com/casper-ecosystem/casper-golang-sdk

## Examples 
End-to-end examples can be found here: https://docs.casperlabs.io/en/latest/dapp-dev-guide/sdk/go-sdk.html

## Tests

```
go test ./...
```

## RPC client usage

```go
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

client := sdk.NewRpcClient("http://157.245.67.232:7777/rpc")

// Getting info about last block
resp, err := client.GetLatestBlock()
```

## Event Stream client usage

```go
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

client = NewEventService("http://159.65.118.250:9999")

// Waiting for 2 blocks and returning that block
block, err := eventClient.AwaitNBlocks(2)

// Waiting for for a deploy and returning information about the deploy
deploy, err := eventClient.AwaitNBlocks(2)
```

## Deploy usage

```go
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

deployParams := DeployParams{
    Account:   account,
    Timestamp: time.Now(),
    TTL:       30 * time.Minute,
    GasPrice:  1,
    ChainName: "casper-test",
}
payment := StandardPayment(big.NewInt(10000000000))
session := NewTransfer(big.NewInt(25000000000), target, uint64(5))
```

## KeyPair package usage

```go
import "github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"

// Creating a random Ed25519 keypair
edKeyPair := ed25519.Ed25519Random()

// Getting Ed25519 public key
edPubKey := edKeyPair.PublicKey()
```

