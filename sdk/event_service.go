package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type BlockResult struct {
	BlockHash	string
	ParentHash	string
	TimeStamp	string
	Eraid		int
	Proposer	string
	State		string
	DeployCount	int
	Height		uint64
	Deploys		[]string
}

type Page struct {
	Number	int
	Url		string
}

type BlocksResult struct {
	Data		[]BlockResult
	PageCount	int
	ItemCount	int
	Pages 		[]Page
}

type DeployRes struct {
	DeployHash		string
	State			string
	Cost 			int
	ErrorMessage	string
	Account			string
	BlockHash		string
}

type DeployHash struct {
	BlockHash		string
	DeployHash		string
	State 			string
	Cost 			int
	ErrorMessage	string
}

type AccountDeploy struct {
	DeployHash		string
	Account 		string
	State			string
	Cost			int
	ErrorMessage	string
	BlockHash		string
}

type AccountDeployResult struct {
	Data		[]AccountDeploy
	PageCount	int
	ItemCount	int
	Pages		[]Page
}

type TransferResult struct {
	DeployHash	string
	SourcePurse	string
	TargetPurse	string
	Amount 		string
	Id			string
	FromAccount	string
	ToAccount	string
}

type EventService struct {
	Url	string
}

func NewEventService(url string) *EventService {
	return &EventService{
		Url: url,
	}
}

func (e EventService) GetResponseData(endpoint string) string {
	get := fmt.Sprintf("%s%s", e.Url, endpoint)

	resp, err := http.Get(get)
	if err != nil {
		fmt.Errorf("failed to get data: %w", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("failed to get response body: %w", err)
	}
	return string(b)
}

func (e EventService) GetBlocks(page int, count int) string {
	endpoint := fmt.Sprintf("/blocks?page=$%d&limit=%d", page, count)
	return e.GetResponseData(endpoint)
}

func (e EventService) GetDeployByHash(deployHash string) string {
	endpoint := fmt.Sprintf("/deploy/%s", deployHash)
	return e.GetResponseData(endpoint)
}

func (e EventService) GetAccountDeploy(accountHex string, page int, limit int) string {
	endpoint := fmt.Sprintf("/accountDeploys/%s?page=%d&limit=%d", accountHex, page, limit)
	return e.GetResponseData(endpoint)
}

func (e EventService) GetTransfersByAccountHash(accountHash string) string {
	endpoint := fmt.Sprintf("/transfers/%s", accountHash)
	return e.GetResponseData(endpoint)
}