# Casper Golang SDK

## Rpc client usage

```
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

client := sdk.NewRpcClient("http://157.245.67.232:7777/rpc")

// Getting info about last block
resp, err := client.GetLatestBlock()

// Getting info about last block by hash
resp, err := client.GetBlockByHash(hash)

// Getting info about last block by height
resp, err := client.GetBlockByHeight(height)

// Getting info about a single deploy by hash
resp, err := client.GetDeploy(hash)

// Getting global state item
resp, err := client.GetStateItem(stateRootHash, key, path)

// Geting info about status 
resp, err := client.GetStatus()

// Getting info about auction
resp, err := client.GetValidator()

// Getting info about peers 
resp, err := client.GetPeers()

// Getting info abot latest block transfers
resp, err := client.GetLatestBlockTransfers()

// Getting info about block transfers by hash
resp, err := client.GetBlockTransfersByHash(blockHash)

// Getting info about block transfers by height
resp, err := client.GetBlockTransfersByHeight(blockHeight)

// Getting info about account balance
resp, err := client.GetAccountBalance(stateRootHash, balanceUref)
```

## KeyPair package usage

```
import "github.com/casper-ecosystem/casper-golang-sdk/keypair"

// Creating a random Ed25519 keypair
edKeyPair := keypair.ed25519.Ed25519Random()

// Getting Ed25519 public key
edPubKey := keypair.ed22519.PublicKey()
```
