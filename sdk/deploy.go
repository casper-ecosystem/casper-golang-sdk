package sdk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
	"io"
	"math/big"
	"time"
)

type AccessRight byte

const (
	AccessRightRead AccessRight = 1 << iota
	AccessRightWrite
	AccessRightAdd
	AccessRightNone AccessRight = 0
)

type URef struct {
	AccessRight AccessRight
	Address     [32]byte
}

type KeyType byte

const (
	KeyTypeAccount KeyType = iota
	KeyTypeHash
	KeyTypeURef
)

type Key struct {
	Type    KeyType
	Account [32]byte
	Hash    [32]byte
	URef    *URef
}

func (u Key) SwitchFieldName() string {
	return "Type"
}

func (u Key) ArmForSwitch(sw byte) (string, bool) {
	switch KeyType(sw) {
	case KeyTypeAccount:
		return "Account", true
	case KeyTypeHash:
		return "Hash", true
	case KeyTypeURef:
		return "URef", true
	}
	return "-", false
}

// Time is custom struct type containing time.Time but it implements its own marshaling.
type Time struct {
	time.Time
}

func (d Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(time.RFC3339Nano))
}

func (d Time) Marshal(w io.Writer) (int, error) {
	encoder := serialization.NewEncoder(w)
	n, err := encoder.Encode(uint64(d.UnixMilli()))
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Duration is custom struct type containing time.Duration but it implements its own marshaling.
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d Duration) Marshal(w io.Writer) (int, error) {
	encoder := serialization.NewEncoder(w)
	n, err := encoder.Encode(uint64(d.Milliseconds()))
	if err != nil {
		return 0, err
	}
	return n, nil
}

type NamedArg struct {
	Name  string
	Value types.CLValue
}

func (r NamedArg) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{r.Name, r.Value})
}

func (r NamedArg) Marshal(w io.Writer) (int, error) {
	size := 0
	encoder := serialization.NewEncoder(w)
	n, err := encoder.Encode(r.Name)
	if err != nil {
		return 0, err
	}
	size += n
	n2, err := r.Value.MarshalWithType(w)
	if err != nil {
		return 0, err
	}
	return n + n2, nil
}

type Deploy struct {
	Hash      keypair.Hash         `json:"hash"`
	Header    DeployHeader         `json:"header"`
	Payment   ExecutableDeployItem `json:"payment"`
	Session   ExecutableDeployItem `json:"session"`
	Approvals []Approval           `json:"approvals"`
}

// MakeDeploy constructs a Deploy from DeployParams, payment and session ExecutableDeployItem
func MakeDeploy(deployParams DeployParams, payment, session ExecutableDeployItem) (Deploy, error) {
	if deployParams.Timestamp.IsZero() {
		deployParams.Timestamp = time.Now()
	}
	deployParams.Timestamp = deployParams.Timestamp.UTC().Truncate(time.Millisecond)
	if deployParams.TTL == 0 {
		deployParams.TTL = time.Minute * 30
	}
	if deployParams.Dependencies == nil {
		deployParams.Dependencies = []keypair.Hash{}
	}
	header := DeployHeader{
		Account:      deployParams.Account,
		Timestamp:    Time{deployParams.Timestamp},
		TTL:          Duration{deployParams.TTL},
		GasPrice:     deployParams.GasPrice,
		Dependencies: deployParams.Dependencies,
		ChainName:    deployParams.ChainName,
	}
	d := Deploy{
		Header:  header,
		Payment: payment,
		Session: session,
	}
	bodyHash, err := d.BodyHash()
	if err != nil {
		return Deploy{}, err
	}
	d.Header.BodyHash = bodyHash
	headerHash, err := d.Header.Hash()
	if err != nil {
		return Deploy{}, err
	}
	d.Hash = headerHash
	return d, nil
}

// BodyHash serializes payment and session and calculates a hash
func (d *Deploy) BodyHash() (keypair.Hash, error) {
	var w bytes.Buffer
	encoder := serialization.NewEncoder(&w)
	_, err := encoder.Encode(d.Payment)
	if err != nil {
		return keypair.Hash{}, err
	}
	_, err = encoder.Encode(d.Session)
	if err != nil {
		return keypair.Hash{}, err
	}
	return keypair.CalculateHash(w.Bytes()), nil
}

// Sign signs the Deploy with given keypair.KeyPair
func (d *Deploy) Sign(pair keypair.KeyPair) error {
	var empty [32]byte
	if d.Hash == empty {
		return errors.New("can't sign because hash is missing (call makeDeploy)")
	}
	if d.Approvals == nil {
		d.Approvals = make([]Approval, 0, 1)
	}
	d.Approvals = append(d.Approvals, Approval{
		Signer:    pair.PublicKey(),
		Signature: pair.Sign(d.Hash[:]),
	})
	return nil
}

type Approval struct {
	Signer    keypair.PublicKey `json:"signer"`
	Signature keypair.Signature `json:"signature"`
}

type DeployHeader struct {
	Account      keypair.PublicKey `json:"account"`
	Timestamp    Time              `json:"timestamp"`
	TTL          Duration          `json:"ttl"`
	GasPrice     uint64            `json:"gas_price"`
	BodyHash     keypair.Hash      `json:"body_hash"`
	Dependencies []keypair.Hash    `json:"dependencies"`
	ChainName    string            `json:"chain_name"`
}

// Hash serializes the DeployHeader and calculates the hash.
func (d *DeployHeader) Hash() (keypair.Hash, error) {
	var w bytes.Buffer
	encoder := serialization.NewEncoder(&w)
	_, err := encoder.Encode(*d)
	if err != nil {
		return keypair.Hash{}, err
	}
	return keypair.CalculateHash(w.Bytes()), nil
}

type DeployParams struct {
	Account      keypair.PublicKey
	ChainName    string
	GasPrice     uint64
	TTL          time.Duration
	Dependencies []keypair.Hash
	Timestamp    time.Time
}

type ModuleBytes struct {
	ModuleBytes []byte     `json:"module_bytes"`
	Args        []NamedArg `json:"args"`
}

func (m ModuleBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ModuleBytes string     `json:"module_bytes"`
		Args        []NamedArg `json:"args"`
	}{
		hex.EncodeToString(m.ModuleBytes),
		m.Args,
	})
}

type StoredContractByHash struct {
	Hash       keypair.Hash `json:"hash"`
	Entrypoint string       `json:"entry_point"`
	Args       []NamedArg   `json:"args"`
}

type StoredContractByName struct {
	Name       string     `json:"name"`
	Entrypoint string     `json:"entry_point"`
	Args       []NamedArg `json:"args"`
}

type StoredVersionedContractByHash struct {
	Hash       keypair.Hash `json:"hash"`
	Version    *uint32      `json:"version"`
	Entrypoint string       `json:"entry_point"`
	Args       []NamedArg   `json:"args"`
}

type StoredVersionedContractByName struct {
	Name       string     `json:"name"`
	Version    *uint32    `json:"version"`
	Entrypoint string     `json:"entry_point"`
	Args       []NamedArg `json:"args"`
}

type Transfer struct {
	Args []NamedArg `json:"args"`
}

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
	Type                          ExecutableDeployItemType       `json:"-"`
	ModuleBytes                   *ModuleBytes                   `json:"ModuleBytes,omitempty"`
	StoredContractByHash          *StoredContractByHash          `json:"StoredContractByHash,omitempty"`
	StoredContractByName          *StoredContractByName          `json:"StoredContractByName,omitempty"`
	StoredVersionedContractByHash *StoredVersionedContractByHash `json:"StoredVersionedContractByHash,omitempty"`
	StoredVersionedContractByName *StoredVersionedContractByName `json:"StoredVersionedContractByName,omitempty"`
	Transfer                      *Transfer                      `json:"Transfer,omitempty"`
}

func (e ExecutableDeployItem) Marshal(w io.Writer) (int, error) {
	size := 0
	encoder := serialization.NewEncoder(w)
	n, err := encoder.Encode(e.Type)
	if err != nil {
		return 0, err
	}
	size += n
	_, typeValue := e.GetValue()
	n, err = encoder.Encode(typeValue)
	if err != nil {
		return 0, err
	}
	size += n
	return size, nil
}

func (u ExecutableDeployItem) SwitchFieldName() string {
	return "Type"
}

func (u ExecutableDeployItem) ArmForSwitch(sw byte) (string, bool) {
	switch ExecutableDeployItemType(sw) {
	case ExecutableDeployItemTypeModuleBytes:
		return "ModuleBytes", true
	case ExecutableDeployItemTypeStoredContractByHash:
		return "StoredContractByHash", true
	case ExecutableDeployItemTypeStoredContractByName:
		return "StoredContractByName", true
	case ExecutableDeployItemTypeStoredVersionedContractByHash:
		return "StoredVersionedContractByHash", true
	case ExecutableDeployItemTypeStoredVersionedContractByName:
		return "StoredVersionedContractByName", true
	case ExecutableDeployItemTypeTransfer:
		return "Transfer", true
	}
	return "-", false
}

func (u ExecutableDeployItem) GetValue() (string, interface{}) {
	switch u.Type {
	case ExecutableDeployItemTypeModuleBytes:
		return "ModuleBytes", *u.ModuleBytes
	case ExecutableDeployItemTypeStoredContractByHash:
		return "StoredContractByHash", *u.StoredContractByHash
	case ExecutableDeployItemTypeStoredContractByName:
		return "StoredContractByName", *u.StoredContractByName
	case ExecutableDeployItemTypeStoredVersionedContractByHash:
		return "StoredVersionedContractByHash", *u.StoredVersionedContractByHash
	case ExecutableDeployItemTypeStoredVersionedContractByName:
		return "StoredVersionedContractByName", *u.StoredVersionedContractByName
	case ExecutableDeployItemTypeTransfer:
		return "Transfer", *u.Transfer
	}
	return "-", nil
}

// StandardPayment provides default payment code for execution engine
// paymentAmount is the number of motos
func StandardPayment(paymentAmount *big.Int) ExecutableDeployItem {
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeModuleBytes,
		ModuleBytes: &ModuleBytes{
			ModuleBytes: []byte{},
			Args: []NamedArg{
				{
					Name: "amount",
					Value: types.CLValue{
						Type: types.CLTypeU512,
						U512: paymentAmount,
					},
				},
			},
		},
	}
}

// NewModuleBytes creates ModuleBytes with given byte slice and args and returns ExecutableDeployItem containing it.
func NewModuleBytes(moduleBytes []byte, args []NamedArg) ExecutableDeployItem {
	if args == nil {
		args = []NamedArg{}
	}
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeModuleBytes,
		ModuleBytes: &ModuleBytes{
			ModuleBytes: moduleBytes,
			Args:        args,
		},
	}
}

// NewStoredContractByHash creates StoredContractByHash with hash, entrypoint and args and returns ExecutableDeployItem containing it.
func NewStoredContractByHash(hash keypair.Hash, entrypoint string, args []NamedArg) ExecutableDeployItem {
	if args == nil {
		args = []NamedArg{}
	}
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByHash,
		StoredContractByHash: &StoredContractByHash{
			Hash:       hash,
			Entrypoint: entrypoint,
			Args:       args,
		},
	}
}

// NewStoredContractByName creates StoredContractByName with name, entrypoint and args and returns ExecutableDeployItem containing it.
func NewStoredContractByName(name, entrypoint string, args []NamedArg) ExecutableDeployItem {
	if args == nil {
		args = []NamedArg{}
	}
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByName,
		StoredContractByName: &StoredContractByName{
			Name:       name,
			Entrypoint: entrypoint,
			Args:       args,
		},
	}
}

// NewStoredVersionedContractByHash creates StoredVersionedContractByHash with hash, entrypoint, optional version and args and returns ExecutableDeployItem containing it.
func NewStoredVersionedContractByHash(hash keypair.Hash, entrypoint string, version *uint32, args []NamedArg) ExecutableDeployItem {
	if args == nil {
		args = []NamedArg{}
	}
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByHash,
		StoredVersionedContractByHash: &StoredVersionedContractByHash{
			Hash:       hash,
			Entrypoint: entrypoint,
			Version:    version,
			Args:       args,
		},
	}
}

// NewStoredVersionedContractByName creates StoredVersionedContractByName with name, entrypoint, optional version and args and returns ExecutableDeployItem containing it.
func NewStoredVersionedContractByName(name, entrypoint string, version *uint32, args []NamedArg) ExecutableDeployItem {
	if args == nil {
		args = []NamedArg{}
	}
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByName,
		StoredVersionedContractByName: &StoredVersionedContractByName{
			Name:       name,
			Entrypoint: entrypoint,
			Version:    version,
			Args:       args,
		},
	}
}

// NewTransfer creates a session ExecutableDeployItem for transferring.
// amount is the number of motos to be transferred.
// target is the account receiving motos.
// id is transfer-id
func NewTransfer(amount *big.Int, target keypair.PublicKey, id uint64) ExecutableDeployItem {
	var args []NamedArg
	args = append(args, NamedArg{
		Name: "amount",
		Value: types.CLValue{
			Type: types.CLTypeU512,
			U512: amount,
		},
	})
	args = append(args, NamedArg{
		Name: "target",
		Value: types.CLValue{
			Type:      types.CLTypeByteArray,
			ByteArray: target.AccountHash(),
		},
	})
	args = append(args, NamedArg{
		Name: "id",
		Value: types.CLValue{
			Type: types.CLTypeOption,
			Optional: &types.Optional{CLValue: &types.CLValue{
				Type: types.CLTypeU64,
				U64:  &id,
			}},
		},
	})
	return ExecutableDeployItem{
		Type: ExecutableDeployItemTypeTransfer,
		Transfer: &Transfer{
			Args: args,
		},
	}
}
