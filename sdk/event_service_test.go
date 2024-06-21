package sdk

import "testing"

var eventService = NewEventService("http://3.136.227.9:7777/rpc")

func TestEventService_GetAccountDeploy(t *testing.T) {
	res, err := eventService.GetAccountDeploy("0169ce4172b5d8f58d1ee9e0f5d24e8210cdad1265d159dd7cdd2aa8beb4ab8ad6", 1, 10)
	t.Log(res)
	if err != nil {
		t.Errorf("can't get account deploy")
	}
}

func TestEventService_GetTransfersByAccountHash(t *testing.T) {
	_, err := eventService.GetTransfersByAccountHash("f57ff0bdcf33a86c34204d46b267d5e8dafae3b3eae376bc5229143c6e6c7897")
	if err != nil {
		t.Errorf("can't get transfers by account hash")
	}
}