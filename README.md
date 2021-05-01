# Casper Golang SDK

## Rpc client usage

```
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

client := sdk.NewRpcClient("http://157.245.67.232:7777/rpc")

// Getting info about last block
resp, err := client.GetLatestBlock()
```

