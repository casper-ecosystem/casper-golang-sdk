package sdk

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"time"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
	"golang.org/x/crypto/blake2b"
)

type UniqAddress struct {
	PublicKey  *keypair.PublicKey
	TransferId uint64
}

type Deploy struct {
	Hash      Hash                  `json:"hash"`
	Header    *DeployHeader         `json:"header"`
	Payment   *ExecutableDeployItem `json:"payment"`
	Session   *ExecutableDeployItem `json:"session"`
	Approvals []Approval            `json:"approvals"`
}

func (d *Deploy) IsStandardPayment() bool {
	if d.Session != nil && d.Session.IsTransfer() {
		return true
	}
	return false
}

func (d *Deploy) IsTransfer() bool {
	if d.Payment != nil && d.Payment.IsModuleBytes() {
		if d.Payment.ModuleBytes.ModuleBytes == nil || len(d.Payment.ModuleBytes.ModuleBytes) == 0 {
			return true
		}
	}
	return false
}

func (d *Deploy) ValidateDeploy() bool {
	serializedBody := SerializeBody(d.Payment, d.Session)
	bodyHash := blake2b.Sum256(serializedBody)

	for i := range bodyHash {
		if d.Header.BodyHash[i] != bodyHash[i] {
			return false
		}
	}

	serializedHeader := SerializeHeader(d.Header)
	deployHash := blake2b.Sum256(serializedHeader)

	for i := range deployHash {
		if d.Hash[i] != deployHash[i] {
			return false
		}
	}

	// TODO: Verify included signatures

	return true
}

func (d *Deploy) SignDeploy(keys keypair.KeyPair) {
	signature := keys.Sign(d.Hash)
	approval := Approval{
		Signer:    keys.PublicKey(),
		Signature: signature,
	}

	if d.Approvals == nil {
		d.Approvals = make([]Approval, 0)
	}

	d.Approvals = append(d.Approvals, approval)
}

func NewDeploy(hash []byte, header *DeployHeader, payment *ExecutableDeployItem, sessions *ExecutableDeployItem,
	approvals []Approval) *Deploy {
	d := new(Deploy)
	d.Hash = hash
	d.Header = header
	d.Payment = payment
	d.Session = sessions
	d.Approvals = approvals
	return d
}

type DeployParams struct {
	AccountPublicKey keypair.PublicKey
	ChainName        string
	GasPrice         uint64
	Ttl              int64
	Dependencies     [][]uint8
	Timestamp        int64
}

func NewDeployParams(accountPublicKey keypair.PublicKey, chainName string, dependencies [][]uint8, timestamp int64) *DeployParams {
	d := new(DeployParams)
	d.AccountPublicKey = accountPublicKey
	d.ChainName = chainName
	d.GasPrice = 1
	d.Ttl = 1800000

	if dependencies == nil {
		d.Dependencies = [][]uint8{}
	} else {
		d.Dependencies = dependencies
	}

	if timestamp == 0 {
		d.Timestamp = time.Now().UnixNano() / 1000000
	} else {
		d.Timestamp = timestamp
	}

	return d
}

type Approval struct {
	Signer    keypair.PublicKey `json:"signer"`
	Signature keypair.Signature `json:"signature"`
}

type DeployHeader struct {
	Account      keypair.PublicKey `json:"account"`
	Timestamp    Timestamp         `json:"timestamp"`
	TTL          Duration          `json:"ttl"`
	GasPrice     uint64            `json:"gas_price"`
	BodyHash     Hash              `json:"body_hash"`
	Dependencies [][]byte          `json:"dependencies"`
	ChainName    string            `json:"chain_name"`
}

func NewDeployHeader(account keypair.PublicKey, timeStamp int64, ttl int64, gasPrice uint64,
	bodyHash []byte, dependencies [][]byte, chainName string) *DeployHeader {
	d := new(DeployHeader)
	d.Account = account
	d.Timestamp = Timestamp(timeStamp)
	d.TTL = Duration(ttl)
	d.GasPrice = gasPrice
	d.BodyHash = bodyHash
	d.Dependencies = dependencies
	d.ChainName = chainName
	return d
}

func (d DeployHeader) ToBytes() []uint8 {
	acc, err := d.Account.ToBytes()
	if err != nil {
		return nil
	}
	result := append(acc, acc...)

	timeSt, err := serialization.Marshal(d.Timestamp)
	if err != nil {
		return nil
	}
	result = append(acc, timeSt...)

	ttl, err := serialization.Marshal(d.TTL)
	if err != nil {
		return nil
	}
	result = append(result, ttl...)

	gasPr, err := serialization.Marshal(d.GasPrice)
	if err != nil {
		return nil
	}
	result = append(result, gasPr...)

	result = append(result, d.BodyHash...)

	dependencies, err := serialization.Marshal(d.Dependencies)
	if err != nil {
		return nil
	}
	result = append(result, dependencies...)

	chainN, err := serialization.Marshal(d.ChainName)
	if err != nil {
		return nil
	}
	result = append(result, chainN...)

	return result
}

func MakeDeploy(deployParam *DeployParams, payment *ExecutableDeployItem, session *ExecutableDeployItem) *Deploy {
	serializedBody := SerializeBody(payment, session)
	bodyHash := blake2b.Sum256(serializedBody)
	resBodyHash := bodyHash[:]

	header := NewDeployHeader(
		deployParam.AccountPublicKey,
		deployParam.Timestamp,
		deployParam.Ttl,
		deployParam.GasPrice,
		resBodyHash,
		deployParam.Dependencies,
		deployParam.ChainName,
	)

	serializedHeader := SerializeHeader(header)
	deployHash := blake2b.Sum256(serializedHeader)
	resDeployHash := make([]uint8, 32)

	for i, v := range deployHash {
		resDeployHash[i] = v
	}

	approvals := make([]Approval, 0)
	return NewDeploy(resDeployHash, header, payment, session, approvals)
}

func NewTransferToUniqAddress(source keypair.PublicKey, target UniqAddress, amount *big.Int, paymentAmount *big.Int, chainName string, sourcePurse string) *Deploy {
	deployParams := NewDeployParams(source, chainName, nil, 0)
	payment := StandardPayment(paymentAmount)
	session := NewTransfer(amount, target.PublicKey, sourcePurse, target.TransferId)

	return MakeDeploy(deployParams, payment, session)
}

func SerializeBody(payment *ExecutableDeployItem, session *ExecutableDeployItem) []byte {
	return append(payment.ToBytes(), session.ToBytes()...)
}

func SerializeHeader(deployHeader *DeployHeader) []byte {
	return deployHeader.ToBytes()
}

func StandardPayment(amount *big.Int) *ExecutableDeployItem {
	amountBytes, err := serialization.Marshal(serialization.U512{Int: *amount})

	if err != nil {
		return nil
	}

	return NewModuleBytes(nil, *NewRunTimeArgs(map[string]Value{"amount": {
		Tag:         types.CLTypeU512,
		StringBytes: hex.EncodeToString(amountBytes),
	}}, append(make([]string, 0), "amount")))
}

func (deploy *Deploy) AddArgToDeploy(key string, value types.CLValue) error {
	if len(deploy.Approvals) != 0 {
		return errors.New("can not add argument to already signed deploy")
	}

	err := deploy.Session.SetArg(key, value)
	if err != nil {
		return err
	}

	return nil
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
	Type                          ExecutableDeployItemType
	ModuleBytes                   *ModuleBytes
	StoredContractByHash          *StoredContractByHash
	StoredContractByName          *StoredContractByName
	StoredVersionedContractByHash *StoredVersionedContractByHash
	StoredVersionedContractByName *StoredVersionedContractByName
	Transfer                      *Transfer
}

func (e *ExecutableDeployItem) IsModuleBytes() bool {
	if e.Type == ExecutableDeployItemTypeModuleBytes {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) IsStoredContractByHash() bool {
	if e.Type == ExecutableDeployItemTypeStoredContractByHash {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) IsStoredContractByName() bool {
	if e.Type == ExecutableDeployItemTypeStoredContractByName {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) IsStoredVersionedContractByHash() bool {
	if e.Type == ExecutableDeployItemTypeStoredVersionedContractByHash {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) IsStoredVersionedContractByName() bool {
	if e.Type == ExecutableDeployItemTypeStoredVersionedContractByName {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) IsTransfer() bool {
	if e.Type == ExecutableDeployItemTypeTransfer {
		return true
	} else {
		return false
	}
}

func (e *ExecutableDeployItem) SetArg(key string, value types.CLValue) error {

	marshaledValue, err := serialization.Marshal(value)
	if err != nil {
		return err
	}

	valueToAdd := Value{
		Tag: value.Type,
	}

	if value.Type == types.CLTypeOption {
		valueToAdd.IsOptional = true
		valueToAdd.Optional = &Value{Tag: value.Option.Type, StringBytes: hex.EncodeToString(marshaledValue)}
	} else {
		valueToAdd.StringBytes = hex.EncodeToString(marshaledValue)
	}

	switch e.Type {
	case ExecutableDeployItemTypeModuleBytes:
		e.ModuleBytes.Args.Insert(key, valueToAdd)
		break
	case ExecutableDeployItemTypeStoredContractByHash:
		e.StoredContractByHash.Args.Insert(key, valueToAdd)
	case ExecutableDeployItemTypeStoredContractByName:
		e.StoredContractByName.Args.Insert(key, valueToAdd)
	case ExecutableDeployItemTypeStoredVersionedContractByHash:
		e.StoredVersionedContractByHash.Args.Insert(key, valueToAdd)
	case ExecutableDeployItemTypeStoredVersionedContractByName:
		e.StoredVersionedContractByName.Args.Insert(key, valueToAdd)
	case ExecutableDeployItemTypeTransfer:
		e.Transfer.Args.Insert(key, valueToAdd)
		break
	}
	return nil
}

func (e *ExecutableDeployItem) ToBytes() []byte {
	switch e.Type {
	case ExecutableDeployItemTypeModuleBytes:
		return e.ModuleBytes.ToBytes()
	case ExecutableDeployItemTypeStoredContractByHash:
		return e.StoredContractByHash.ToBytes()
	case ExecutableDeployItemTypeStoredContractByName:
		return e.StoredContractByName.ToBytes()
	case ExecutableDeployItemTypeStoredVersionedContractByHash:
		return e.StoredVersionedContractByHash.ToBytes()
	case ExecutableDeployItemTypeStoredVersionedContractByName:
		return e.StoredVersionedContractByName.ToBytes()
	case ExecutableDeployItemTypeTransfer:
		return e.Transfer.ToBytes()
	}
	return nil
}

func (e ExecutableDeployItem) MarshalJSON() ([]byte, error) {
	switch e.Type {
	case ExecutableDeployItemTypeModuleBytes:
		return e.ModuleBytes.MarshalJSON()
	case ExecutableDeployItemTypeStoredContractByHash:
		return e.StoredContractByHash.MarshalJSON()
	case ExecutableDeployItemTypeStoredContractByName:
		return e.StoredContractByName.MarshalJSON()
	case ExecutableDeployItemTypeStoredVersionedContractByHash:
		return e.StoredVersionedContractByHash.MarshalJSON()
	case ExecutableDeployItemTypeStoredVersionedContractByName:
		return e.StoredVersionedContractByName.MarshalJSON()
	case ExecutableDeployItemTypeTransfer:
		return e.Transfer.MarshalJSON()
	}
	return nil, nil
}

func (e *ExecutableDeployItem) UnmarshalJSON(data []byte) error {
	var tempMap map[string]interface{}

	err := json.Unmarshal(data, &tempMap)
	if err != nil {
		return err
	}

	if tempMap["ModuleBytes"] != nil {
		tempMap = tempMap["ModuleBytes"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeModuleBytes
		e.ModuleBytes = &ModuleBytes{Tag: e.Type}
		return e.ModuleBytes.UnmarshalModuleBytes(tempMap)
	} else if tempMap["StoredContractByHash"] != nil {
		tempMap = tempMap["StoredContractByHash"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeStoredContractByHash
		e.StoredContractByHash = &StoredContractByHash{Tag: e.Type}
		return e.StoredContractByHash.UnmarshalStoredContractByHash(tempMap)
	} else if tempMap["StoredContractByName"] != nil {
		tempMap = tempMap["StoredContractByName"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeStoredContractByName
		e.StoredContractByName = &StoredContractByName{Tag: e.Type}
		return e.StoredContractByName.UnmarshalStoredContractByName(tempMap)
	} else if tempMap["StoredVersionedContractByHash"] != nil {
		tempMap = tempMap["StoredVersionedContractByHash"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeStoredVersionedContractByHash
		e.StoredVersionedContractByHash = &StoredVersionedContractByHash{Tag: e.Type}
		return e.StoredVersionedContractByHash.UnmarshalStoredVersionedContractByHash(tempMap)
	} else if tempMap["StoredVersionedContractByName"] != nil {
		tempMap = tempMap["StoredVersionedContractByName"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeStoredVersionedContractByName
		e.StoredVersionedContractByName = &StoredVersionedContractByName{Tag: e.Type}
		return e.StoredVersionedContractByName.UnmarshalStoredVersionedContractByName(tempMap)
	} else if tempMap["Transfer"] != nil {
		tempMap = tempMap["Transfer"].(map[string]interface{})
		e.Type = ExecutableDeployItemTypeTransfer
		e.Transfer = &Transfer{Tag: e.Type}
		return e.Transfer.UnmarshalTransfer(tempMap)
	}

	return nil
}

func (e ExecutableDeployItem) SwitchFieldName() string {
	return "Type"
}

func (e ExecutableDeployItem) ArmForSwitch(sw byte) (string, bool) {
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

type ModuleBytes struct {
	Tag         ExecutableDeployItemType
	ModuleBytes []byte
	Args        RuntimeArgs
}

func NewModuleBytes(moduleBytes []byte, args RuntimeArgs) *ExecutableDeployItem {
	if moduleBytes == nil {
		moduleBytes = make([]byte, 0)
	}

	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeModuleBytes,
		ModuleBytes: &ModuleBytes{
			Tag:         ExecutableDeployItemTypeModuleBytes,
			ModuleBytes: moduleBytes,
			Args:        args,
		},
	}
}

func (m ModuleBytes) ToBytes() []byte {
	res := make([]byte, 0)

	moduleBytes, err := serialization.Marshal(m.ModuleBytes)
	if err != nil {
		return nil
	}

	res = append([]byte{byte(m.Tag)}, moduleBytes...)
	res = append(res, m.Args.ToBytes()...)
	return res
}

func (m ModuleBytes) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"ModuleBytes": map[string]interface{}{
			"module_bytes": hex.EncodeToString(m.ModuleBytes),
			"args":         m.Args.ToJSONInterface(),
		},
	}
	return json.Marshal(data)
}

func (m *ModuleBytes) UnmarshalModuleBytes(tempMap map[string]interface{}) error {
	if tempMap["module_bytes"] == nil {
		m.ModuleBytes = make([]byte, 0)
	} else {
		m.ModuleBytes = []byte(tempMap["module_bytes"].(string))
	}

	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})

	var err error
	m.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}

type StoredContractByHash struct {
	Tag        ExecutableDeployItemType
	Hash       [32]byte
	Entrypoint string
	Args       RuntimeArgs
}

func NewStoredContractByHash(hash [32]byte, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByHash,
		StoredContractByHash: &StoredContractByHash{
			Tag:        ExecutableDeployItemTypeStoredContractByHash,
			Hash:       hash,
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func (s StoredContractByHash) ToBytes() []byte {
	res := make([]byte, 0)

	entrypoint, err := serialization.Marshal(s.Entrypoint)
	if err != nil {
		return nil
	}

	res = append([]byte{byte(s.Tag)}, s.Hash[:]...)
	res = append(res, entrypoint...)
	res = append(res, s.Args.ToBytes()...)
	return res
}

func (s StoredContractByHash) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"StoredContractByHash": map[string]interface{}{
			"hash":        hex.EncodeToString(s.Hash[:]),
			"entry_point": s.Entrypoint,
			"args":        s.Args.ToJSONInterface(),
		},
	}
	return json.Marshal(data)
}

func (s *StoredContractByHash) UnmarshalStoredContractByHash(tempMap map[string]interface{}) error {
	var ok bool
	tempHash, ok := tempMap["hash"].(string)

	decodedHash, err := hex.DecodeString(tempHash)
	if err != nil {
		return err
	}

	if len(decodedHash) != 32 {
		return errors.New("invalid hash length")
	}

	for i := 0; i < 32; i++ {
		s.Hash[i] = decodedHash[i]
	}

	s.Entrypoint, ok = tempMap["entry_point"].(string)

	if !ok {
		return errors.New("invalid json, no entry_point")
	}

	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})

	s.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}

type StoredContractByName struct {
	Tag        ExecutableDeployItemType
	Name       string
	Entrypoint string
	Args       RuntimeArgs
}

func NewStoredContractByName(name, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByName,
		StoredContractByName: &StoredContractByName{
			Tag:        ExecutableDeployItemTypeStoredContractByName,
			Name:       name,
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func (s StoredContractByName) ToBytes() []byte {
	res := make([]byte, 0)

	name, err := serialization.Marshal(s.Name)
	if err != nil {
		return nil
	}

	entrypoint, err := serialization.Marshal(s.Entrypoint)
	if err != nil {
		return nil
	}

	res = append([]byte{byte(s.Tag)}, name...)
	res = append(res, entrypoint...)
	res = append(res, s.Args.ToBytes()...)
	return res
}

func (s StoredContractByName) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"StoredContractByName": map[string]interface{}{
			"name":        s.Name,
			"entry_point": s.Entrypoint,
			"args":        s.Args.ToJSONInterface(),
		},
	}
	return json.Marshal(data)
}

func (s *StoredContractByName) UnmarshalStoredContractByName(tempMap map[string]interface{}) error {
	var ok bool
	s.Name, ok = tempMap["name"].(string)

	if !ok {
		return errors.New("invalid json, no name")
	}

	s.Entrypoint, ok = tempMap["entry_point"].(string)

	if !ok {
		return errors.New("invalid json, no entry_point")
	}

	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})

	var err error
	s.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}

type StoredVersionedContractByHash struct {
	Tag        ExecutableDeployItemType
	Hash       [32]byte
	Version    *types.CLValue
	Entrypoint string
	Args       RuntimeArgs
}

func NewStoredVersionedContractByHashWithoutVersion(hash [32]byte, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByHash,
		StoredVersionedContractByHash: &StoredVersionedContractByHash{
			Tag:        ExecutableDeployItemTypeStoredVersionedContractByHash,
			Hash:       hash,
			Version:    &types.CLValue{Type: types.CLTypeOption, Option: nil},
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func NewStoredVersionedContractByHash(hash [32]byte, version uint32, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByHash,
		StoredVersionedContractByHash: &StoredVersionedContractByHash{
			Tag:  ExecutableDeployItemTypeStoredVersionedContractByHash,
			Hash: hash,
			Version: &types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{
				Type: types.CLTypeU32,
				U32:  &version,
			}},
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func (s StoredVersionedContractByHash) ToBytes() []byte {
	res := make([]byte, 0)

	entrypoint, err := serialization.Marshal(s.Entrypoint)
	if err != nil {
		return nil
	}

	version, err := serialization.Marshal(*s.Version)

	res = append([]byte{byte(s.Tag)}, s.Hash[:]...)
	res = append(res, version...)
	res = append(res, entrypoint...)
	res = append(res, s.Args.ToBytes()...)
	return res
}

func (s StoredVersionedContractByHash) MarshalJSON() ([]byte, error) {
	var version string
	if s.Version.Option != nil {
		version = strconv.Itoa(int(*s.Version.Option.U32))
	} else {
		version = "None"
	}

	data := map[string]interface{}{
		"StoredVersionedContractByHash": map[string]interface{}{
			"hash":        hex.EncodeToString(s.Hash[:]),
			"entry_point": s.Entrypoint,
			"args":        s.Args.ToJSONInterface(),
			"version":     version,
		},
	}

	return json.Marshal(data)
}

func (s *StoredVersionedContractByHash) UnmarshalStoredVersionedContractByHash(tempMap map[string]interface{}) error {
	var ok bool
	tempHash, ok := tempMap["hash"].(string)

	decodedHash, err := hex.DecodeString(tempHash)
	if err != nil {
		return err
	}

	if len(decodedHash) != 32 {
		return errors.New("invalid hash length")
	}

	for i := 0; i < 32; i++ {
		s.Hash[i] = decodedHash[i]
	}

	version, ok := tempMap["version"].(string)

	if version == "None" {
		s.Version = &types.CLValue{Type: types.CLTypeOption, Option: nil}
	} else {
		parsedVersion, err := strconv.ParseInt(version, 10, 32)
		if err != nil {
			return err
		}
		parsedVersion32 := uint32(parsedVersion)
		s.Version = &types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU32, U32: &parsedVersion32}}
	}

	s.Entrypoint, ok = tempMap["entry_point"].(string)

	if !ok {
		return errors.New("invalid json, no entry_point")
	}

	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})

	s.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}

type StoredVersionedContractByName struct {
	Tag        ExecutableDeployItemType
	Name       string
	Version    *types.CLValue
	Entrypoint string
	Args       RuntimeArgs
}

func NewStoredVersionedContractByNameWithoutVersion(name string, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByName,
		StoredVersionedContractByName: &StoredVersionedContractByName{
			Tag:        ExecutableDeployItemTypeStoredVersionedContractByName,
			Name:       name,
			Version:    &types.CLValue{Type: types.CLTypeOption, Option: nil},
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func NewStoredVersionedContractByName(name string, version uint32, entryPoint string, args RuntimeArgs) *ExecutableDeployItem {
	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByName,
		StoredVersionedContractByName: &StoredVersionedContractByName{
			Tag:  ExecutableDeployItemTypeStoredVersionedContractByName,
			Name: name,
			Version: &types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{
				Type: types.CLTypeU32,
				U32:  &version,
			}},
			Entrypoint: entryPoint,
			Args:       args,
		},
	}
}

func (s StoredVersionedContractByName) ToBytes() []byte {
	res := make([]byte, 0)

	name, err := serialization.Marshal(s.Name)
	if err != nil {
		return nil
	}

	entrypoint, err := serialization.Marshal(s.Entrypoint)
	if err != nil {
		return nil
	}

	version, err := serialization.Marshal(*s.Version)

	res = append([]byte{byte(s.Tag)}, name...)
	res = append(res, version...)
	res = append(res, entrypoint...)
	res = append(res, s.Args.ToBytes()...)
	return res
}

func (s StoredVersionedContractByName) MarshalJSON() ([]byte, error) {
	var version string
	if s.Version.Option != nil {
		version = strconv.Itoa(int(*s.Version.Option.U32))
	} else {
		version = "None"
	}

	data := map[string]interface{}{
		"StoredVersionedContractByName": map[string]interface{}{
			"name":        s.Name,
			"entry_point": s.Entrypoint,
			"args":        s.Args.ToJSONInterface(),
			"version":     version,
		},
	}
	return json.Marshal(data)
}

func (s *StoredVersionedContractByName) UnmarshalStoredVersionedContractByName(tempMap map[string]interface{}) error {
	var ok bool
	s.Name, ok = tempMap["name"].(string)

	if !ok {
		return errors.New("invalid json, no name")
	}

	s.Entrypoint, ok = tempMap["entry_point"].(string)

	if !ok {
		return errors.New("invalid json, no entry_point")
	}

	version, ok := tempMap["version"].(string)

	if version == "None" {
		s.Version = &types.CLValue{Type: types.CLTypeOption, Option: nil}
	} else {
		parsedVersion, err := strconv.ParseInt(version, 10, 32)
		if err != nil {
			return err
		}
		parsedVersion32 := uint32(parsedVersion)
		s.Version = &types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU32, U32: &parsedVersion32}}
	}

	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})

	var err error
	s.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}

type Transfer struct {
	Tag  ExecutableDeployItemType
	Args RuntimeArgs
}

func NewTransfer(amount *big.Int, target *keypair.PublicKey, sourcePurse string, id uint64) *ExecutableDeployItem {
	return buildTransfer(amount, target, sourcePurse, true, id)
}

func NewTransferWithoutId(amount *big.Int, target *keypair.PublicKey, sourcePurse string) *ExecutableDeployItem {
	return buildTransfer(amount, target, sourcePurse, false, 0)
}

func buildTransfer(amount *big.Int, target *keypair.PublicKey, sourcePurse string, idPresent bool, id uint64) *ExecutableDeployItem {
	amountBytes, err := serialization.Marshal(serialization.U512{Int: *amount})
	if err != nil {
		return nil
	}

	var accountHex string

	if target.Tag == keypair.KeyTagEd25519 {
		accountHex = ed25519.AccountHash(target.PubKeyData)
	} 
	if target.Tag == keypair.KeyTagSecp256k1 {
		accountHex = append(0x02,target.PubKeyData...)
	}

	idValue := Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &Value{
			Tag: types.CLTypeU64,
		},
	}

	if idPresent {
		idBytes, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU64, U64: &id}})

		if err != nil {
			return nil
		}

		idValue.Optional.StringBytes = hex.EncodeToString(idBytes)
	} else {
		idBytes, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption})

		if err != nil {
			return nil
		}

		idValue.Optional.StringBytes = hex.EncodeToString(idBytes)
	}

	var argsOrder []string

	if sourcePurse == "" {
		argsOrder = append(make([]string, 0), "amount", "target", "id")
	} else {
		argsOrder = append(make([]string, 0), "amount", "target", "sourcePurse", "id")
	}


	if target.Tag == keypair.KeyTagEd25519 {
		args := NewRunTimeArgs(map[string]Value{
		"amount": {
			Tag:         types.CLTypeU512,
			StringBytes: hex.EncodeToString(amountBytes),
		},
		"target": {
			Tag:         types.CLTypeByteArray,
			StringBytes: accountHex,
		},
		"id": idValue,
		}, argsOrder)
	}

	if target.Tag == keypair.KeyTagSecp256k1 {
		args := NewRunTimeArgs(map[string]Value{
		"amount": {
			Tag:         types.CLTypeU512,
			StringBytes: hex.EncodeToString(amountBytes),
		},
		"target": {
			Tag:         types.CLTypePublicKey,
			StringBytes: accountHex,
		},
		"id": idValue,
		}, argsOrder)
	}



	if sourcePurse != "" {
		args.Args["source"] = Value{Tag: types.CLTypeURef, StringBytes: sourcePurse}
	}

	return &ExecutableDeployItem{
		Type: ExecutableDeployItemTypeTransfer,
		Transfer: &Transfer{
			Tag:  ExecutableDeployItemTypeTransfer,
			Args: *args,
		},
	}
}

func (t Transfer) ToBytes() []byte {
	a := make([]uint8, 1)
	a[0] = uint8(t.Tag)
	return append(a, t.Args.ToBytes()...)
}

func (t Transfer) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"Transfer": map[string]interface{}{
			"args": t.Args.ToJSONInterface(),
		},
	}

	return json.Marshal(data)
}

func (t *Transfer) UnmarshalTransfer(tempMap map[string]interface{}) error {
	if tempMap["args"] == nil {
		return errors.New("'args' key doesn't exist, invalid json")
	}

	args := tempMap["args"].([]interface{})
	var err error
	t.Args, err = ParseRuntimeArgs(args)

	if err != nil {
		return errors.New("error while parsing args")
	}

	return nil
}
