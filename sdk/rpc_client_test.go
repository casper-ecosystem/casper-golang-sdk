package sdk

import "testing"

var client = NewRpcClient("http://157.245.67.232:7777/rpc")

func TestRpcClient_GetLatestBlock(t *testing.T) {
	_, err:= client.GetLatestBlock()

	if err != nil {
		t.Errorf("can't get latest block")
	}
}

func TestRpcClient_GetDeploy(t *testing.T) {
	_, err := client.GetDeploy("")

	if err != nil {
		t.Errorf("can't get latest block")
	}
}

func TestRpcClient_GetBlockState(t *testing.T) {
	path := []string{}
	_, err := client.GetBlockState("", "", path)

	if err != nil {
		t.Errorf("can't get block state")
	}
}

func TestRpcClient_GetAccountBalance(t *testing.T) {
	_, err := client.GetAccountBalance("", "")

	if err != nil {
		t.Errorf("can't get account balance")
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





