package sdk

import "github.com/casper-ecosystem/casper-golang-sdk/keypair"

type Deploy struct {
	Hash      [32]byte
	Header    DeployHeader
	Payment   ExecutableDeployItem
	Session   ExecutableDeployItem
	Approvals []Approval
}

type Approval struct {
	Signer    keypair.PublicKey
	Signature keypair.Signature
}

type DeployHeader struct {
	Account      keypair.PublicKey
	Timestamp    uint64
	TTL          uint64
	GasPrice     uint64
	BodyHash     [32]byte
	Dependencies [][32]byte
	ChainName    string
}

type ModuleBytes struct {
	ModuleBytes []byte
	Args        []byte
}

type StoredContractByHash struct {
	Hash       [32]byte
	Entrypoint string
	Args       []byte
}

type StoredContractByName struct {
	Name       string
	Entrypoint string
	Args       []byte
}

type StoredVersionedContractByHash struct {
	Hash       [32]byte
	Version    *uint32
	Entrypoint string
	Args       []byte
}

type StoredVersionedContractByName struct {
	Name       string
	Version    *uint32
	Entrypoint string
	Args       []byte
}

type Transfer struct {
	Args []byte
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
