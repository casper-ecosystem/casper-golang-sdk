package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Testnet node
var client, _ = NewRpcClient("http://159.65.118.250:7777/rpc")

func TestRpcClient_GetLatestBlock(t *testing.T) {
	_, err := client.GetLatestBlock()

	assert.NoError(t, err)
}

func TestRpcClient_GetDeploy(t *testing.T) {
	hash := "1dfdf144eb0422eae3076cd8a17e55089010a133c6555c881ac0b9e2714a1605"
	_, err := client.GetDeploy(hash)

	assert.NoError(t, err)
}

func TestRpcClient_GetBlockState(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	path := []string{"special_value"}
	_, err := client.GetStateItem(stateRootHash, key, path)

	assert.NoError(t, err)
}

func TestRpcClient_GetAccountBalanceUrefByPublicKeyHash(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	accountHash := "a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	want := "uref-15f95f6acf6eea00353318bd76648f9c1a319020d295479164db18019be49e55-007"

	uref, err := client.GetAccountBalanceUrefByPublicKeyHash(stateRootHash, accountHash)
	assert.NoError(t, err)
	assert.Equal(t, want, uref)
}

func TestRpcClient_GetAccountBalance(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	//path := []string{"special_value"}
	item, err := client.GetStateItem(stateRootHash, key, nil)

	assert.NoError(t, err)
	assert.NotNil(t, item.Account)

	balanceUref := item.Account.MainPurse
	_, err = client.GetAccountBalance(stateRootHash, balanceUref)

	assert.NoError(t, err)
}

func TestRpcClient_GetBlockByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	assert.NoError(t, err)
}

func TestRpcClient_GetBlockTransfersByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	assert.NoError(t, err)
}

func TestRpcClient_GetBlockByHash(t *testing.T) {
	_, err := client.GetBlockByHash("")

	assert.NoError(t, err)
}

func TestRpcClient_GetBlockTransfersByHash(t *testing.T) {
	_, err := client.GetBlockTransfersByHash("")

	assert.NoError(t, err)
}

func TestRpcClient_GetLatestBlockTransfers(t *testing.T) {
	_, err := client.GetLatestBlockTransfers()

	assert.NoError(t, err)
}

func TestRpcClient_GetValidator(t *testing.T) {
	_, err := client.GetValidator()

	assert.NoError(t, err)
}

func TestRpcClient_GetStatus(t *testing.T) {
	_, err := client.GetStatus()

	assert.NoError(t, err)
}

func TestRpcClient_GetMetrics(t *testing.T) {
	metrics, err := client.GetMetrics()

	assert.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "metrics should be a long string")
}

func TestRpcClient_GetRpcSchema(t *testing.T) {
	schema, err := client.GetRpcSchema()

	assert.NoError(t, err)
	assert.Greater(t, len(schema), 0, "schema should be a long string")
}

func TestRpcClient_GetPeers(t *testing.T) {
	_, err := client.GetPeers()

	assert.NoError(t, err)
}

func TestRpcClient_GetEraBySwitchBlockHeight(t *testing.T) {
	t.Run("should return error because block height 1 is not a switch block", func(t *testing.T) {
		_, err := client.GetEraBySwitchBlockHeight(1)
		assert.Error(t, err)
	})
	t.Run("should get era summary because block height 212742 is a switch block", func(t *testing.T) {
		_, err := client.GetEraBySwitchBlockHeight(212742)
		assert.NoError(t, err)
	})
}

func TestRpcClient_GetLatestStateRootHash(t *testing.T) {
	_, err := client.GetLatestStateRootHash()

	assert.NoError(t, err)
}

func TestRpcClient_GetStateRootHashByHeight(t *testing.T) {
	got, err := client.GetStateRootHashByHeight(1)
	want := StateRootHashResult{"e88b7c061760134ba37ad312c1e2d6373121748e9c61bcea19cc57510829addf"}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRpcClient_GetStateRootHashByHash(t *testing.T) {
	got, err := client.GetStateRootHashByHash("f8ebb3a81c9c70faeaec896c94e7e56c75a5e1548e1b8dfe639c1b31610e5d22")
	want := StateRootHashResult{"e88b7c061760134ba37ad312c1e2d6373121748e9c61bcea19cc57510829addf"}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
