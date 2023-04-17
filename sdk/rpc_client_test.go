package sdk

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = NewRpcClient("http://65.21.238.180:7777/rpc")

func TestRpcClient_GetLatestBlock(t *testing.T) {
	_, err := client.GetLatestBlock()

	if err != nil {
		t.Errorf("can't get latest block")
	}
}

func TestRpcClient_GetDeploy(t *testing.T) {
	hash := "1dfdf144eb0422eae3076cd8a17e55089010a133c6555c881ac0b9e2714a1605"
	_, err := client.GetDeploy(hash)

	if err != nil {
		t.Errorf("can't get deploy info")
	}
}

func TestRpcClient_GetBlockState(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	path := []string{"special_value"}
	_, err := client.GetStateItem(stateRootHash, key, path)

	if err != nil {
		t.Errorf("can't get block state")
	}
}

func TestRpcClient_GetAccountBalance(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"

	balanceUref := client.GetAccountMainPurseURef(key)

	_, err := client.GetAccountBalance(stateRootHash, balanceUref)

	if err != nil {
		t.Errorf("can't get account balance")
	}
}

func TestRpcClient_GetAccountInfo(t *testing.T) {
	pubkey := "020237037ff4845669e59d3e7698e7d58eb97ca378960ac57478a86a6a3535460292"
	bad_pubkey := "020237037ff4845669e59d3e7698e7d58ee97ca378960ac57478a86a6a3535460292"

	block := "a705d5cf0bca0ec2cc0ffceb2913900669b1234c230340086304989d67dde7d7"

	_, err := client.GetAccountInfo(bad_pubkey, block)
	var rpcerror *RpcError

	if err != nil {
		if errors.As(err, &rpcerror) {
			if rpcerror.Code != -32003 {
				t.Errorf("bad error code")
			}
		}
	}
	_, err = client.GetAccountInfo(pubkey, block)
	if err != nil {
		t.Errorf("can't get account info")
	}
	//t.Errorf("can't get account info")
	//log.Println("Running test...")
	//log.Println(err)
}

func TestRpcClient_GetAccountBalanceByKeypair(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	stateRootHashNew, err := client.GetStateRootHash(stateRootHash)
	if err != nil {
		return
	}
	path := []string{"special_value"}
	_, err = client.GetStateItem(stateRootHashNew.StateRootHash, sourceKeyPair.AccountHash(), path)
	if err != nil {
		_, err := client.GetAccountBalanceByKeypair(stateRootHashNew.StateRootHash, sourceKeyPair)

		if err != nil {
			t.Errorf("can't get account balance")
		}
	}
}

func TestRpcClient_GetBlockByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block by height")
	}
}

func TestRpcClient_GetBlockTransfersByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block transfers by height")
	}
}

func TestRpcClient_GetBlockByHash(t *testing.T) {
	_, err := client.GetBlockByHash("")

	if err != nil {
		t.Errorf("can't get block by hash")
	}
}

func TestRpcClient_GetBlockTransfersByHash(t *testing.T) {
	_, err := client.GetBlockTransfersByHash("")

	if err != nil {
		t.Errorf("can't get block transfers by hash")
	}
}

func TestRpcClient_GetLatestBlockTransfers(t *testing.T) {
	_, err := client.GetLatestBlockTransfers()

	if err != nil {
		t.Errorf("can't get latest block transfers")
	}
}

func TestRpcClient_GetValidator(t *testing.T) {
	_, err := client.GetValidator()

	if err != nil {
		t.Errorf("can't get validator")
	}
}

func TestRpcClient_GetStatus(t *testing.T) {
	_, err := client.GetStatus()

	if err != nil {
		t.Errorf("can't get status")
	}
}

func TestRpcClient_GetPeers(t *testing.T) {
	_, err := client.GetPeers()

	if err != nil {
		t.Errorf("can't get peers")
	}
}

// make sure your account has balance
func TestRpcClient_PutDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3000000000), big.NewInt(10000), "casper-test", "")

	assert.True(t, deploy.ValidateDeploy())
	deploy.SignDeploy(sourceKeyPair)

	result, err := client.PutDeploy(*deploy)

	if !assert.NoError(t, err) {
		t.Errorf("error : %v", err)
		return
	}

	assert.Equal(t, hex.EncodeToString(deploy.Hash), result.Hash)
}
