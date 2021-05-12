package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BlockResult struct {
	BlockHash	string		`json:"block_hash"`
	ParentHash	string		`json:"parent_hash"`
	TimeStamp	string		`json:"time_stamp"`
	Eraid		int			`json:"eraid"`
	Proposer	string		`json:"proposer"`
	State		string		`json:"state"`
	DeployCount	int			`json:"deploy_count"`
	Height		uint64		`json:"height"`
	Deploys		[]string	`json:"deploys"`
}

type Page struct {
	Number	int		`json:"number"`
	Url		string	`json:"url"`
}

type BlocksResult struct {
	Data		[]BlockResult	`json:"data"`
	PageCount	int				`json:"page_count"`
	ItemCount	int				`json:"item_count"`
	Pages 		[]Page			`json:"pages"`
}

type DeployRes struct {
	DeployHash		string	`json:"deploy_hash"`
	State			string	`json:"state"`
	Cost 			int		`json:"cost"`
	ErrorMessage	string	`json:"error_message"`
	Account			string	`json:"account"`
	BlockHash		string	`json:"block_hash"`
}

type DeployHash struct {
	BlockHash		string	`json:"block_hash"`
	DeployHash		string	`json:"deploy_hash"`
	State 			string	`json:"state"`
	Cost 			int		`json:"cost"`
	ErrorMessage	string	`json:"error_message"`
}

type AccountDeploy struct {
	DeployHash		string	`json:"deploy_hash"`
	Account 		string	`json:"account"`
	State			string	`json:"state"`
	Cost			int		`json:"cost"`
	ErrorMessage	string	`json:"error_message"`
	BlockHash		string	`json:"block_hash"`
}

type AccountDeploysResult struct {
	Data		[]AccountDeploy	`json:"data"`
	PageCount	int				`json:"page_count"`
	ItemCount	int				`json:"item_count"`
	Pages		[]Page			`json:"pages"`
}

type TransferResult struct {
	DeployHash	string	`json:"deploy_hash"`
	SourcePurse	string	`json:"source_purse"`
	TargetPurse	string	`json:"target_purse"`
	Amount 		string	`json:"amount"`
	Id			string	`json:"id"`
	FromAccount	string	`json:"from_account"`
	ToAccount	string	`json:"to_account"`
}

type EventService struct {
	Url	string
}

func NewEventService(url string) *EventService {
	return &EventService{
		Url: url,
	}
}

func (e EventService) GetResponseData(endpoint string) ([]byte,error) {
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}
	return b, nil
}

func (e EventService) GetBlocks(page int, count int) (BlocksResult, error) {
	endpoint := fmt.Sprintf("/blocks?page=%d&limit=%d", page, count)
	resp, err := e.GetResponseData(endpoint)

	if err != nil {
		return BlocksResult{}, err
	}

	var blocks BlocksResult
	parseResponseBody(resp, blocks)

	return blocks, nil
}

func (e EventService) GetDeployByHash(deployHash string) (DeployResult, error) {
	endpoint := fmt.Sprintf("/deploy/%s", deployHash)
	resp, err := e.GetResponseData(endpoint)

	if err != nil {
		return DeployResult{}, err
	}

	var deploy DeployResult
	parseResponseBody(resp, deploy)

	return deploy, nil
}

func (e EventService) GetBlockByHash(blockHash string) (BlockResult, error) {
	endpoint := fmt.Sprintf("/block/%s", blockHash)
	resp, err := e.GetResponseData(endpoint)

	if err != nil {
		return BlockResult{}, err
	}

	var block BlockResult
	parseResponseBody(resp, block)

	return block, nil
}

func (e EventService) GetAccountDeploy(accountHex string, page int, limit int) (AccountDeploysResult, error) {
	endpoint := fmt.Sprintf("/accountDeploys/%s?page=%d&limit=%d", accountHex, page, limit)
	resp, err := e.GetResponseData(endpoint)

	if err != nil {
		return AccountDeploysResult{}, err
	}

	var accountDeploys AccountDeploysResult
	parseResponseBody(resp, accountDeploys)

	return accountDeploys, nil
}

func (e EventService) GetTransfersByAccountHash(accountHash string) ([]TransferResult, error) {
	endpoint := fmt.Sprintf("/transfers/%s", accountHash)
	resp, err := e.GetResponseData(endpoint)

	if err != nil {
		return []TransferResult{}, err
	}

	var transfers []TransferResult
	parseResponseBody(resp, transfers)

	return transfers, nil
}

func parseResponseBody(response []byte, dest interface{})  {
	err := json.Unmarshal(response, &dest)
	if err != nil {
		fmt.Errorf("failed to parse response body: %w", err)
	}
}