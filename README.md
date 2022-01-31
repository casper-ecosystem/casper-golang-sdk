# Casper Go SDK

### Install

`go get -u github.com/casper-ecosystem/casper-golang-sdk`

### Running tests

`go test ./...`

### Keypair package usage

#### ED25519

```go
import "github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"

...

// Creating a random Ed25519 keypair
randomKeyPair, err := ed25519.Ed25519Random()

// Parsing keypair from pem files
keyPairFromFileED25519, err := ed25519.ParseKeyFiles(
"keypair/test_account_keys/account1/public_key.pem",
"keypair/test_account_keys/account1/secret_key.pem")

...
```

Alternatively you can parse keys via various functions

ED25519 specification can be found here: \
https://datatracker.ietf.org/doc/html/rfc8410

```go
// Parse public key from file
func ParsePublicKeyFile(path string) ([]byte, error)
// Parse private key from file
func ParsePrivateKeyFile(path string) ([]byte, error)

// Parse private key from bytes, according to ED25519 specification
// it is encoded with extra bytes before the actual key
// in case of private key it will omit first 16 bytes of the array
func ParsePrivateKey(bytes []byte) ([]byte, error)
// in case of public key it will omit first 12 bytes of the array
func ParsePublicKey(bytes []byte) ([]byte, error)

// Utility function for custom key parsing
func ParseKey(bytes []byte, from int, to int) ([]byte, error)
// Creating ed25519 keypair from bytes
func ParseKeyPair(publicKey []byte, privateKey []byte) keypair.KeyPair
// Get ed25519 public key in the form of 010..0
func AccountHex(publicKey []byte) string
// Get account hash by public key
func AccountHash(pubKey []byte) string
// Export private key in pem
func ExportPrivateKeyInPem(privateKey string) []byte
// Export public key in pem
func ExportPublicKeyInPem(publicKey string) []byte
```

#### KeyPair methods

```go

// Getting public key
publicKey := keyPairFromFile.PublicKey()

// Get account hash in the form 'account-hash-0..0'
accountHash := keyPairFromFile.AccountHash()

// Get seed from which keypair was generated
seedBytes := keyPairFromFile.RawSeed()

// Get signature
signature := keyPairFromFile.Sign([]byte("Hello Casper!"))

// Verify signature
verified := keyPairFromFile.Verify(signature.SignatureData, []byte("Hello Casper!")))
```

### Serialization

Information about serialization standarts can be found here \
https://docs.casperlabs.io/en/latest/implementation/serialization-standard.html

#### Simple types

More simple types examples can be found at _serialization/coding_test.go_

###### Encoding

```go
import "github.com/casper-ecosystem/casper-golang-sdk/serialization"

...

res, err := serialization.Marshal(int64(1024)) // res == "0004000000000000"

...
```

###### Decoding

```go
import (
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"encoding/hex"
	)

...

toDecode, err := hex.DecodeString("0004000000000000")
if err != nil{return}

var dest int64

err = serialization.Unmarshal(toDecode, &dest) // dest == 1024 
if err != nil{return}

...
```

#### CLValue

More CLValue types examples can be found at _types/coding_test.go_

###### Encoding

```go
import "github.com/casper-ecosystem/casper-golang-sdk/serialization"

...

helloString := "Hello Casper!"

res, err := serialization.Marshal(
    types.CLValue{
        Type: types.CLTypeString,
        String: &helloString,
    }) // res == "0d00000048656c6c6f2043617370657221"

if err != nil {return}

...
```

###### Decoding

When decoding CLValue type of the CLValue destination struct has to be specified \
(for each subtype too) \
E.g. When decoding CLValue of the type Option, you have to specify that you are expecting Option type to be parsed and the inner type.

```go
// CLValue of the type Option, containing String
dest := types.CLValue{
	Type: types.CLTypeOption,
	Option: &types.CLValue{
		Type: types.CLTypeString,
	},
}
```

```go
import (
    "github.com/casper-ecosystem/casper-golang-sdk/serialization"
    "encoding/hex"
    )

...

toDecode, err := hex.DecodeString("0d00000048656c6c6f2043617370657221")
if err != nil {return}

dest := types.CLValue{Type: types.CLTypeString}
_, err = types.UnmarshalCLValue(toDecode, &dest)
if err != nil {return}
// *dest.String == "Hello Casper!" 

...
```

## SDK

### Deploy

For deploy following structs are defined

#### ExecutableDeployItem structs

```go
type ExecutableDeployItemType byte

const (
	ExecutableDeployItemTypeModuleBytes ExecutableDeployItemType = iota
	ExecutableDeployItemTypeStoredContractByHash
	ExecutableDeployItemTypeStoredContractByName
	ExecutableDeployItemTypeStoredVersionedContractByHash
	ExecutableDeployItemTypeStoredVersionedContractByName
	ExecutableDeployItemTypeTransfer
)

type ExecutableDeployItem struct {
	Type                          ExecutableDeployItemType
	ModuleBytes                   *ModuleBytes
	StoredContractByHash          *StoredContractByHash
	StoredContractByName          *StoredContractByName
	StoredVersionedContractByHash *StoredVersionedContractByHash
	StoredVersionedContractByName *StoredVersionedContractByName
	Transfer                      *Transfer
}

type ModuleBytes struct {
	Tag         ExecutableDeployItemType
	ModuleBytes []byte
	Args        RuntimeArgs
}

type StoredContractByHash struct {
	Tag        ExecutableDeployItemType
	Hash       [32]byte
	Entrypoint string
	Args       RuntimeArgs
}

type StoredContractByName struct {
	Tag        ExecutableDeployItemType
	Name       string
	Entrypoint string
	Args       RuntimeArgs
}

type StoredVersionedContractByHash struct {
	Tag        ExecutableDeployItemType
	Hash       [32]byte
	Version    *types.CLValue
	Entrypoint string
	Args       RuntimeArgs
}

type StoredVersionedContractByName struct {
	Tag ExecutableDeployItemType
	Name       string
	Version    *types.CLValue
	Entrypoint string
	Args       RuntimeArgs
}

type Transfer struct {
	Tag  ExecutableDeployItemType
	Args RuntimeArgs
}
```

#### Runtime arguments

Value struct differs from the CLValue, as they are serialized differently for the deploy bytes. But CLValue bytes used in StringBytes field.

Optional value can be enforced with IsOptional flag. If Tag of the Value is Map, then Map field should be defined with types of the keys and values.
```go
type ValueMap struct {
	KeyType types.CLType
	ValueType types.CLType
}

type Value struct {
	Tag         types.CLType
	IsOptional  bool
	StringBytes string
	Optional    *Value
	Map 	    *ValueMap
}

type RuntimeArgs struct {
	KeyOrder []string
	Args     map[string]Value
}

```

#### Deploy structs

```go
type Hash []byte

type Timestamp int64

type Duration time.Duration

type DeployHeader struct {
	Account      keypair.PublicKey `json:"account"`
	Timestamp    Timestamp         `json:"timestamp"`
	TTL          Duration          `json:"ttl"`
	GasPrice     uint64            `json:"gas_price"`
	BodyHash     Hash              `json:"body_hash"`
	Dependencies [][]byte          `json:"dependencies"`
	ChainName    string            `json:"chain_name"`
}

type UniqAddress struct {
	PublicKey  *keypair.PublicKey
	TransferId uint64
}

type DeployParams struct {
	AccountPublicKey keypair.PublicKey
	ChainName        string
	GasPrice         uint64
	Ttl              int64
	Dependencies     [][]uint8
	Timestamp        int64
}

type Approval struct {
	Signer    keypair.PublicKey `json:"signer"`
	Signature keypair.Signature `json:"signature"`
}

type Deploy struct {
	Hash      Hash                  `json:"hash"`
	Header    *DeployHeader         `json:"header"`
	Payment   *ExecutableDeployItem `json:"payment"`
	Session   *ExecutableDeployItem `json:"session"`
	Approvals []Approval            `json:"approvals"`
}
```

##### Create deploy

```go
import (
	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	)

...

srcKP, err := ed25519.ParseKeyFiles(
	"keypair/test_account_keys/account1/public_key.pem", 
	"keypair/test_account_keys/account1/secret_key.pem")
if err != nil {return}

// Amount of fees paid in motes (10e9)
payment := sdk.StandardPayment(big.NewInt(10000))

// Amount of to transfer in motes (10e9)
session := sdk.NewTransfer(big.NewInt(2500000000), &destPub, "", uint64(1))

// Creating deploy on testnet
deploy := sdk.MakeDeploy(sdk.NewDeployParams(srcKP.PublicKey(), "casper-test", nil, 0), payment, session)

// Signing the deploy
deploy.SignDeploy(srcKP)

// Validating it
deploy.ValidateDeploy() //returns true if deploy is correct

// We can validate our signature
srcKP.Verify(deploy.Approvals[0].Signature.SignatureData,deploy.Hash) //returns true if the signature was correct

...
```

#### All methods that are available

###### ExecutableDeployItem

The methods are either getters of the type or marshaling to bytes or unmarshaling to instance 

```go
func (e *ExecutableDeployItem) IsModuleBytes() bool
func (e *ExecutableDeployItem) IsStoredContractByHash()
func (e *ExecutableDeployItem) IsStoredContractByName() bool
func (e *ExecutableDeployItem) IsStoredVersionedContractByHash() bool
func (e *ExecutableDeployItem) IsStoredVersionedContractByName() bool
func (e *ExecutableDeployItem) IsTransfer() bool
func (e *ExecutableDeployItem) SetArg(key string, value types.CLValue) error
func (e *ExecutableDeployItem) ToBytes() []byte
func (e ExecutableDeployItem)  MarshalJSON() ([]byte, error)
func (e *ExecutableDeployItem) UnmarshalJSON(data []byte) error
func (e ExecutableDeployItem)  SwitchFieldName() string
func (e ExecutableDeployItem)  ArmForSwitch(sw byte) (string, bool)

func StandardPayment(amount *big.Int) *ExecutableDeployItem
```

For each ExecutableDeployItem subtype (as defined above) the following methods implemented

```go
func New<SUBTYPE_NAME>(moduleBytes []byte, args RuntimeArgs) *ExecutableDeployItem
func (m <SUBTYPE_NAME>)  ToBytes() []byte
func (m <SUBTYPE_NAME>)  MarshalJSON() ([]byte, error)
func (m *<SUBTYPE_NAME>) Unmarshal<SUBTYPE_NAME>(tempMap map[string]interface{}) error
```

Additionally, there exist modifications of the methods when some arguments are not needed

###### Runtime arguments

```go
// Marshaling Value struct
func (v Value) MarshalJSON() ([]byte, error)
// Getting bytes according to serialization of the deploy fields
func (v Value) ToBytes() []uint8

// Creating new RuntimeArgs from map, with provided keyOrdering (required for correct serialization)
func NewRunTimeArgs(args map[string]Value, keyOrder []string) *RuntimeArgs
// Getting RuntimeArgs from map of value, setting it's own order
func ParseRuntimeArgs(args []interface{}) (RuntimeArgs, error)
// Getting RuntimeArgs from map of value with key ordering
func (r RuntimeArgs) FromMap(args map[string]Value, keyOrder []string) *RuntimeArgs
// Insert Value to the map of args
func (r RuntimeArgs) Insert(key string, value Value)
// Serializing args according to the rules of serialization
func (r RuntimeArgs) ToBytes() []byte

// Returns array of arrays of length 2 of interfaces ([2][]interface{})
// runtime args in deploy are encoded as an array of pairs "key" : "value"
func (r RuntimeArgs) ToJSONInterface() []interface{}
```

###### Deploy

```go
// Check if deploy is standart payment of CSPR
func (d *Deploy) IsStandardPayment() bool
// Check if deploy is transfer of CSPR (not contract interaction)
func (d *Deploy) IsTransfer() bool
// Check if deploy is valid (correct hash and signatures)
func (d *Deploy) ValidateDeploy() bool
// Add signature to the deploy
func (d *Deploy) SignDeploy(keys keypair.KeyPair)
// Add arguments to the session of the constructed deploy
func (d *Deploy) AddArgToDeploy(key string, value types.CLValue) error

// Marshal deploy header to bytes array
func (d DeployHeader) ToBytes() []uint8

// Create new deploy from the known elements
func NewDeploy(hash []byte, header *DeployHeader, payment *ExecutableDeployItem, sessions *ExecutableDeployItem, approvals []Approval) *Deploy
// Create new default transfer to unique address (public key with corresponding transfer id
func NewTransferToUniqAddress(source keypair.PublicKey, target UniqAddress, amount *big.Int, paymentAmount *big.Int, chainName string, sourcePurse string) *Deploy
// Create deploy from params, payment and session
func MakeDeploy(deployParam *DeployParams, payment *ExecutableDeployItem, session *ExecutableDeployItem) *Deploy

// Create deploy header from the parameters
func NewDeployHeader(account keypair.PublicKey, timeStamp int64, ttl int64, gasPrice uint64, bodyHash []byte, dependencies [][]byte, chainName string) *DeployHeader
// Create correct byte representation of the body of the deploy
func SerializeBody(payment *ExecutableDeployItem, session *ExecutableDeployItem) []byte
// Create correct byte representation of the header of the deploy
func SerializeHeader(deployHeader *DeployHeader) []byte

// Create deploy params, last two arguments are optional
func NewDeployParams(accountPublicKey keypair.PublicKey, chainName string, dependencies [][]uint8, timestamp int64) *DeployParams
```

## RPC client

#### Usage example

```go
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

...

var client = sdk.NewRpcClient("http://0.000.000.0:7777/rpc")

// Get info about latest block
latestBlock, err := client.GetLatestBlock()

...
```


#### RPC client structs

```go
type RpcClient struct {
	endpoint string
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
	ID         int64 `json:"id,omitempty"`
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

type DeployResult struct {
	Deploy           JsonDeploy            `json:"deploy"`
	ExecutionResults []JsonExecutionResult `json:"execution_results"`
}

type JsonDeploy struct {
	Hash      string           `json:"hash"`
	Header    JsonDeployHeader `json:"header"`
	Approvals []JsonApproval   `json:"approvals"`
}

type JsonPutDeployRes struct{
	Hash string `json:"deploy_hash"`
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
	PublicKey	string	`json:"public_key"`
	Weight 		string	`json:"weight"`
}

type EraValidators struct {
	EraId             int                   `json:"era_id"`
	ValidatorWeights  []ValidatorWeight 	`json:"validator_weights"`
}

type AuctionState struct {
	StateRootHash	string	`json:"state_root_hash"`
	BlockHeight 	uint64	`json:"block_height"`
	EraValidators 	[]EraValidators `json:"era_validators"`
}

type ValidatorPesponse struct {
	Version	string	`json:"jsonrpc"`
	AuctionState `json:"auction_state"`
}

type validatorResult struct {
	Validator ValidatorPesponse `json:"validator"`
}

type StatusResult struct {
	LastAddedBlock	BlockResponse `json:"last_added_block"`
	BuildVersion	string `json:"build_version"`
}

type Peer struct {
	NodeId	string	`json:"node_id"`
	Address	string	`json:"address"`
}

type PeerResult struct {
	Peers	[]Peer	`json:"peers"`
}

type StateRootHashResult struct {
	StateRootHash	string `json:"state_root_hash"`
}
```

#### All methods that are available

```go
func NewRpcClient(endpoint string) *RpcClient
func (c *RpcClient) GetDeploy(hash string) (DeployResult, error)
func (c *RpcClient) GetStateItem(stateRootHash, key string, path []string) (StoredValue, error)
func (c *RpcClient) GetAccountBalance(stateRootHash, balanceUref string) (big.Int, error)
func (c *RpcClient) GetAccountMainPurseURef(accountHash string) string
func (c *RpcClient) GetAccountBalanceByKeypair(stateRootHash string, key keypair.KeyPair) (big.Int, error)
func (c *RpcClient) GetLatestBlock() (BlockResponse, error)
func (c *RpcClient) GetBlockByHeight(height uint64) (BlockResponse, error)
func (c *RpcClient) GetBlockByHash(hash string) (BlockResponse, error)
func (c *RpcClient) GetLatestBlockTransfers() ([]TransferResponse, error)
func (c *RpcClient) GetBlockTransfersByHeight(height uint64) ([]TransferResponse, error)
func (c *RpcClient) GetBlockTransfersByHash(blockHash string) ([]TransferResponse, error)
func (c *RpcClient) GetValidator() (ValidatorPesponse, error)
func (c *RpcClient) GetStatus() (StatusResult, error)
func (c *RpcClient) GetPeers() (PeerResult, error)
func (c *RpcClient) GetStateRootHash(stateRootHash string) (StateRootHashResult, error)
func (c *RpcClient) PutDeploy(deploy Deploy) (JsonPutDeployRes, error)
```

### HTTP Event service

#### Usage example

```go
import "github.com/casper-ecosystem/casper-golang-sdk/sdk"

...

var eventService = sdk.NewEventService("http://0.000.000.0:7777/api")

// Get 2 blocks on page 0
blocks, err := eventService.GetBlocks(0, 2)

...
```

#### EventService structs

```go
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
```

#### All methods that are available

```go
func NewEventService(url string) *EventService
func (e EventService) GetBlocks(page int, count int) (BlocksResult, error)
func (e EventService) GetDeployByHash(deployHash string) (DeployResult, error)
func (e EventService) GetBlockByHash(blockHash string) (BlockResult, error)
func (e EventService) GetAccountDeploy(accountHex string, page int, limit int) (AccountDeploysResult, error)
func (e EventService) GetTransfersByAccountHash(accountHash string) ([]TransferResult, error)
```
