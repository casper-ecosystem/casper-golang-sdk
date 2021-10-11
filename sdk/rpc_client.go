package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type RpcClient struct {
	endpoint *url.URL
}

// NewRpcClient create a new RPC client with specified URL.
// URL should contain RPC path for example http://159.65.118.250:7777/rpc.
func NewRpcClient(endpoint string) (*RpcClient, error) {
	parsedUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return &RpcClient{
		endpoint: parsedUrl,
	}, nil
}

// GetDeploy gets information about a deploy with specified hash.
func (c *RpcClient) GetDeploy(hash string) (DeployResult, error) {
	resp, err := c.rpcCall("info_get_deploy", map[string]string{
		"deploy_hash": hash,
	})
	if err != nil {
		return DeployResult{}, err
	}

	var result DeployResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return DeployResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

// GetStateItem retrieves a stored value from the network.
// key can be: public key in hex format, account hash, contract address hash, uref, transfer hash or deploy-info hash.
// path is optional path to the item.
func (c *RpcClient) GetStateItem(stateRootHash, key string, path []string) (StoredValue, error) {
	params := map[string]interface{}{
		"state_root_hash": stateRootHash,
		"key":             key,
	}
	if len(path) > 0 {
		params["path"] = path
	}
	resp, err := c.rpcCall("state_get_item", params)
	if err != nil {
		return StoredValue{}, err
	}

	var result storedValueResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StoredValue{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result.StoredValue, nil
}

// GetAccountBalance returns balance for specified uref and state root hash.
func (c *RpcClient) GetAccountBalance(stateRootHash, balanceUref string) (big.Int, error) {
	resp, err := c.rpcCall("state_get_balance", map[string]string{
		"state_root_hash": stateRootHash,
		"purse_uref":      balanceUref,
	})
	if err != nil {
		return big.Int{}, err
	}

	var result balanceResponse
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return big.Int{}, fmt.Errorf("failed to get result: %w", err)
	}

	balance := big.Int{}
	balance.SetString(result.BalanceValue, 10)
	return balance, nil
}

// GetLatestAccountBalance returns latest balance for specified uref.
func (c *RpcClient) GetLatestAccountBalance(balanceUref string) (big.Int, error) {
	result, err := c.GetLatestStateRootHash()
	if err != nil {
		return big.Int{}, err
	}
	balance, err := c.GetAccountBalance(result.StateRootHash, balanceUref)
	if err != nil {
		return big.Int{}, err
	}
	return balance, nil
}

// GetAccountBalanceUrefByPublicKeyHex returns main/balance purse for specified public key in hex format.
func (c *RpcClient) GetAccountBalanceUrefByPublicKeyHex(stateRootHash, hex string) (string, error) {
	value, err := c.GetStateItem(stateRootHash, hex, nil)
	if err != nil {
		return "", nil
	}
	if value.Account == nil {
		return "", errors.New("supplied key is not an Account")
	}
	return value.Account.MainPurse, nil
}

// GetAccountBalanceUrefByPublicKeyHash returns main/balance purse for specified public key hash.
func (c *RpcClient) GetAccountBalanceUrefByPublicKeyHash(stateRootHash, hash string) (string, error) {
	value, err := c.GetStateItem(stateRootHash, "account-hash-"+hash, nil)
	if err != nil {
		return "", err
	}
	if value.Account == nil {
		return "", errors.New("supplied key is not an Account")
	}
	return value.Account.MainPurse, nil
}

// GetAccountBalanceUrefByPublicKey returns main/balance purse for specified public key.
func (c *RpcClient) GetAccountBalanceUrefByPublicKey(stateRootHash string, key keypair.PublicKey) (string, error) {
	value, err := c.GetStateItem(stateRootHash, "account-hash-"+key.AccountHashHex(), nil)
	if err != nil {
		return "", nil
	}
	if value.Account == nil {
		return "", errors.New("supplied key is not an Account")
	}
	return value.Account.MainPurse, nil
}

func (c *RpcClient) getBlock(params blockParams) (BlockResponse, error) {
	resp, err := c.rpcCall("chain_get_block", params)
	if err != nil {
		return BlockResponse{}, err
	}

	var result blockResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return BlockResponse{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result.Block, nil
}

// GetLatestBlock returns the latest block information.
func (c *RpcClient) GetLatestBlock() (BlockResponse, error) {
	block, err := c.getBlock(blockParams{})
	if err != nil {
		return BlockResponse{}, err
	}
	return block, nil
}

// GetBlockByHeight returns a block specified by height.
func (c *RpcClient) GetBlockByHeight(height uint64) (BlockResponse, error) {
	block, err := c.getBlock(
		blockParams{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return BlockResponse{}, err
	}
	return block, nil
}

// GetBlockByHash returns a block specified by hash.
func (c *RpcClient) GetBlockByHash(hash string) (BlockResponse, error) {
	block, err := c.getBlock(
		blockParams{blockIdentifier{
			Hash: hash,
		}})
	if err != nil {
		return BlockResponse{}, err
	}
	return block, nil
}

func (c *RpcClient) getBlockTransfers(params blockParams) ([]TransferResponse, error) {
	resp, err := c.rpcCall("chain_get_block_transfers", params)
	if err != nil {
		return nil, err
	}

	var result transferResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get result: %w", err)
	}

	return result.Transfers, nil
}

// GetLatestBlockTransfers returns information about all transfers in the latest block.
func (c *RpcClient) GetLatestBlockTransfers() ([]TransferResponse, error) {
	transfers, err := c.getBlockTransfers(blockParams{})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetBlockTransfersByHeight returns information about all transfers in a block specified by height.
func (c *RpcClient) GetBlockTransfersByHeight(height uint64) ([]TransferResponse, error) {
	transfers, err := c.getBlockTransfers(
		blockParams{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetBlockTransfersByHash returns information about all transfers in a block specified by hash.
func (c *RpcClient) GetBlockTransfersByHash(blockHash string) ([]TransferResponse, error) {
	transfers, err := c.getBlockTransfers(
		blockParams{blockIdentifier{
			Hash: blockHash,
		}})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// GetValidator returns information about era validators of the latest block.
func (c *RpcClient) GetValidator() (ValidatorPesponse, error) {
	resp, err := c.rpcCall("state_get_auction_info", nil)
	if err != nil {
		return ValidatorPesponse{}, err
	}

	var result validatorResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return ValidatorPesponse{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result.Validator, nil
}

// GetStatus retrieves node information.
func (c *RpcClient) GetStatus() (StatusResult, error) {
	resp, err := c.rpcCall("info_get_status", nil)
	if err != nil {
		return StatusResult{}, err
	}

	var result StatusResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StatusResult{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result, nil
}

// GetMetrics retrieves node metrics.
func (c *RpcClient) GetMetrics() (string, error) {
	resp, err := c.httpCall("/metrics")
	if err != nil {
		return "", err
	}
	return resp, nil
}

// GetRpcSchema retrieves node's RPC schema.
func (c *RpcClient) GetRpcSchema() (string, error) {
	resp, err := c.rpcCall("rpc.discover", nil)
	if err != nil {
		return "", err
	}
	return string(resp.Result), nil
}

// GetPeers returns all peers connected to the node.
func (c *RpcClient) GetPeers() (PeerResult, error) {
	resp, err := c.rpcCall("info_get_peers", nil)
	if err != nil {
		return PeerResult{}, err
	}

	var result PeerResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return PeerResult{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result, nil
}

func (c *RpcClient) getEraBySwitchBlock(params blockParams) (EraSummary, error) {
	resp, err := c.rpcCall("chain_get_era_info_by_switch_block", params)
	if err != nil {
		return EraSummary{}, err
	}

	var result eraResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return EraSummary{}, err
	}

	if result.EraSummary == nil {
		return EraSummary{}, fmt.Errorf("provided block is not a switch block")
	}

	return *result.EraSummary, nil
}

// GetLatestEraBySwitchBlock gets era by latest block. If latest block is not a switch block, it returns an error.
func (c *RpcClient) GetLatestEraBySwitchBlock() (EraSummary, error) {
	eraSummary, err := c.getEraBySwitchBlock(blockParams{})
	if err != nil {
		return EraSummary{}, err
	}
	return eraSummary, nil
}

// GetEraBySwitchBlockHeight gets era by switch block height. If provided block is not a switch block, it returns an error.
func (c *RpcClient) GetEraBySwitchBlockHeight(height uint64) (EraSummary, error) {
	eraSummary, err := c.getEraBySwitchBlock(blockParams{blockIdentifier{Height: height}})
	if err != nil {
		return EraSummary{}, err
	}
	return eraSummary, nil
}

// GetEraBySwitchBlockHash gets era by switch block hash. If provided block is not a switch block, it returns an error.
func (c *RpcClient) GetEraBySwitchBlockHash(hash string) (EraSummary, error) {
	eraSummary, err := c.getEraBySwitchBlock(blockParams{blockIdentifier{Hash: hash}})
	if err != nil {
		return EraSummary{}, err
	}
	return eraSummary, nil
}

func (c *RpcClient) getStateRootHash(params blockParams) (StateRootHashResult, error) {
	resp, err := c.rpcCall("chain_get_state_root_hash", params)
	if err != nil {
		return StateRootHashResult{}, err
	}

	var result StateRootHashResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StateRootHashResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

// GetLatestStateRootHash returns state root hash for the latest block.
func (c *RpcClient) GetLatestStateRootHash() (StateRootHashResult, error) {
	result, err := c.getStateRootHash(blockParams{})
	if err != nil {
		return StateRootHashResult{}, err
	}
	return result, nil
}

// GetStateRootHashByHeight returns state root hash for a block specified by height.
func (c *RpcClient) GetStateRootHashByHeight(height uint64) (StateRootHashResult, error) {
	result, err := c.getStateRootHash(
		blockParams{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return StateRootHashResult{}, err
	}
	return result, nil
}

// GetStateRootHashByHash returns state root hash for a block specified by hash.
func (c *RpcClient) GetStateRootHashByHash(hash string) (StateRootHashResult, error) {
	result, err := c.getStateRootHash(
		blockParams{blockIdentifier{
			Hash: hash,
		}})
	if err != nil {
		return StateRootHashResult{}, err
	}
	return result, nil
}

// PutDeploy sends deploy to the node.
func (c *RpcClient) PutDeploy(deploy Deploy) (DeployResponse, error) {
	resp, err := c.rpcCall("account_put_deploy", map[string]interface{}{
		"deploy": deploy,
	})
	if err != nil {
		return DeployResponse{}, err
	}

	var result DeployResponse
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return DeployResponse{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *RpcClient) httpCall(uri string) (string, error) {
	host, _, _ := net.SplitHostPort(c.endpoint.Host)
	endpointUrl := fmt.Sprintf("%s://%s:%s%s", c.endpoint.Scheme, host, "8888", uri)
	resp, err := http.Get(endpointUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	return string(b), nil
}

func (c *RpcClient) rpcCall(method string, params interface{}) (RpcResponse, error) {
	body, err := json.Marshal(RpcRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return RpcResponse{}, errors.Wrap(err, "failed to marshal json")
	}

	resp, err := http.Post(c.endpoint.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to get response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return RpcResponse{}, fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	var rpcResponse RpcResponse
	err = json.Unmarshal(b, &rpcResponse)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to parse response body: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse, fmt.Errorf("rpc call failed, code - %d, message - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse, nil
}

type RpcRequest struct {
	Version string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type RpcResponse struct {
	Version string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RpcError       `json:"error,omitempty"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type transferResult struct {
	Transfers []TransferResponse `json:"transfers"`
}

type TransferResponse struct {
	ID         *string `json:"id,omitempty"`
	DeployHash string  `json:"deploy_hash"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Amount     string  `json:"amount"`
	Gas        string  `json:"gas"`
}

type blockResult struct {
	Block BlockResponse `json:"block"`
}

type BlockResponse struct {
	Hash   string      `json:"hash"`
	Header BlockHeader `json:"header"`
	Body   BlockBody   `json:"body"`
	Proofs []Proof     `json:"proofs"`
}

type BlockHeader struct {
	ParentHash      string    `json:"parent_hash"`
	StateRootHash   string    `json:"state_root_hash"`
	BodyHash        string    `json:"body_hash"`
	RandomBit       bool      `json:"random_bit"`
	AccumulatedSeed string    `json:"accumulated_seed"`
	Timestamp       time.Time `json:"timestamp"`
	EraID           int       `json:"era_id"`
	Height          int       `json:"height"`
	ProtocolVersion string    `json:"protocol_version"`
}

type BlockBody struct {
	Proposer       string   `json:"proposer"`
	DeployHashes   []string `json:"deploy_hashes"`
	TransferHashes []string `json:"transfer_hashes"`
}

type Proof struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type DeployResponse struct {
	DeployHash string `json:"deploy_hash"`
}

type DeployResult struct {
	Deploy           JsonDeploy            `json:"deploy"`
	ExecutionResults []JsonExecutionResult `json:"execution_results"`
}

type JsonDeploy struct {
	Hash      string           `json:"hash"`
	Header    JsonDeployHeader `json:"header"`
	Approvals []JsonApproval   `json:"approvals"`
}

type JsonDeployHeader struct {
	Account      string    `json:"account"`
	Timestamp    time.Time `json:"timestamp"`
	TTL          string    `json:"ttl"`
	GasPrice     int       `json:"gas_price"`
	BodyHash     string    `json:"body_hash"`
	Dependencies []string  `json:"dependencies"`
	ChainName    string    `json:"chain_name"`
}

type JsonApproval struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}

type JsonExecutionResult struct {
	BlockHash string          `json:"block_hash"`
	Result    ExecutionResult `json:"result"`
}

type ExecutionResult struct {
	Success      SuccessExecutionResult `json:"success"`
	ErrorMessage *string                `json:"error_message,omitempty"`
}

type SuccessExecutionResult struct {
	Transfers []string `json:"transfers"`
	Cost      string   `json:"cost"`
}

type storedValueResult struct {
	StoredValue StoredValue `json:"stored_value"`
}

type StoredValue struct {
	CLValue         *JsonCLValue          `json:"CLValue,omitempty"`
	Account         *JsonAccount          `json:"Account,omitempty"`
	Contract        *JsonContractMetadata `json:"Contract,omitempty"`
	ContractWASM    *string               `json:"ContractWASM,omitempty"`
	ContractPackage *string               `json:"ContractPackage,omitempty"`
	Transfer        *TransferResponse     `json:"Transfer,omitempty"`
	DeployInfo      *JsonDeployInfo       `json:"DeployInfo,omitempty"`
	EraInfo         *EraInfo              `json:"EraInfo,omitempty"`
}

type EraInfo struct {
	SeigniorageAllocations []SeigniorageAllocation `json:"seigniorage_allocations"`
}

type SeigniorageAllocation struct {
	Validator *Validator `json:"Validator,omitempty"`
	Delegator *Delegator `json:"Delegator,omitempty"`
}

type Delegator struct {
	Amount                string `json:"amount"`
	ValidatorPublicKeyHex string `json:"validator_public_key"`
	DelegatorPublicKeyHex string `json:"delegator_public_key"`
}

type Validator struct {
	Amount                string `json:"amount"`
	ValidatorPublicKeyHex string `json:"validator_public_key"`
}

type JsonCLValue struct {
	Bytes  string      `json:"bytes"`
	CLType string      `json:"cl_type"`
	Parsed interface{} `json:"parsed"`
}

type JsonAccount struct {
	AccountHash      string           `json:"account_hash"`
	NamedKeys        []NamedKey       `json:"named_keys"`
	MainPurse        string           `json:"main_purse"`
	AssociatedKeys   []AssociatedKey  `json:"associated_keys"`
	ActionThresholds ActionThresholds `json:"action_thresholds"`
}

type NamedKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type AssociatedKey struct {
	AccountHash string `json:"account_hash"`
	Weight      uint64 `json:"weight"`
}

type ActionThresholds struct {
	Deployment    uint64 `json:"deployment"`
	KeyManagement uint64 `json:"key_management"`
}

type JsonContractMetadata struct {
	ContractPackageHash string `json:"contract_package_hash"`
	ContractWasmHash    string `json:"contract_wasm_hash"`
	ProtocolVersion     string `json:"protocol_version"`
}

type JsonDeployInfo struct {
	DeployHash string   `json:"deploy_hash"`
	Transfers  []string `json:"transfers"`
	From       string   `json:"from"`
	Source     string   `json:"source"`
	Gas        string   `json:"gas"`
}

type blockParams struct {
	BlockIdentifier blockIdentifier `json:"block_identifier"`
}

type blockIdentifier struct {
	Hash   string `json:"Hash,omitempty"`
	Height uint64 `json:"Height,omitempty"`
}

type balanceResponse struct {
	BalanceValue string `json:"balance_value"`
}

type ValidatorWeight struct {
	PublicKey string `json:"public_key"`
	Weight    string `json:"weight"`
}

type EraValidators struct {
	EraId            int               `json:"era_id"`
	ValidatorWeights []ValidatorWeight `json:"validator_weights"`
}

type AuctionState struct {
	StateRootHash string          `json:"state_root_hash"`
	BlockHeight   uint64          `json:"block_height"`
	EraValidators []EraValidators `json:"era_validators"`
}

type ValidatorPesponse struct {
	Version      string `json:"jsonrpc"`
	AuctionState `json:"auction_state"`
}

type validatorResult struct {
	Validator ValidatorPesponse `json:"validator"`
}

type StatusResult struct {
	ApiVersion            string         `json:"api_version"`
	ChainspecName         string         `json:"chainspec_name"`
	LastAddedBlock        *BlockResponse `json:"last_added_block,omitempty"`
	NextUpgrade           *Upgrade       `json:"next_upgrade,omitempty"`
	PublicSigningKey      string         `json:"our_public_signing_key"`
	Peers                 []Peer         `json:"peers"`
	RoundLength           *string        `json:"round_length,omitempty"`
	StartingStateRootHash *string        `json:"starting_state_root_hash,omitempty"`
	BuildVersion          string         `json:"build_version"`
}

type Upgrade struct {
	ActivationPoint int    `json:"activation_point"`
	ProtocolVersion string `json:"protocol_version"`
}

type Peer struct {
	NodeId  string `json:"node_id"`
	Address string `json:"address"`
}

type PeerResult struct {
	Peers []Peer `json:"peers"`
}

type StateRootHashResult struct {
	StateRootHash string `json:"state_root_hash"`
}

type eraResult struct {
	EraSummary *EraSummary `json:"era_summary,omitempty"`
}

type EraSummary struct {
	BlockHash     string      `json:"block_hash"`
	EraId         int         `json:"era_id"`
	MerkleProof   string      `json:"merkle_proof"`
	StateRootHash string      `json:"state_root_hash"`
	StoredValue   StoredValue `json:"stored_value"`
}
