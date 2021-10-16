package sdk

import (
	"bytes"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestNamedArg_Marshal(t *testing.T) {
	serializedArg := []byte{
		6, 0, 0, 0, 97, 109, 111,
		117, 110, 116, 4, 0, 0, 0,
		3, 128, 150, 152, 8}
	arg := NamedArg{
		Name: "amount",
		Value: types.CLValue{
			Type: types.CLTypeU512,
			U512: big.NewInt(10000000),
		},
	}

	AssertEqualSerialization(t, arg, serializedArg)
}

func TestModuleBytes_Marshal(t *testing.T) {
	serializedDeployItem := []byte{0, 3, 0, 0, 0, 10, 20, 30, 1,
		0, 0, 0, 6, 0, 0, 0, 97, 109,
		111, 117, 110, 116, 4, 0, 0, 0, 3,
		128, 150, 152, 8}
	deployItem := ExecutableDeployItem{
		Type: ExecutableDeployItemTypeModuleBytes,
		ModuleBytes: &ModuleBytes{
			ModuleBytes: []byte{10, 20, 30},
			Args: []NamedArg{{
				Name: "amount",
				Value: types.CLValue{
					Type: types.CLTypeU512,
					U512: big.NewInt(10000000),
				},
			}},
		},
	}

	AssertEqualSerialization(t, deployItem, serializedDeployItem)
}

func TestStoredContractByHash_Marshal(t *testing.T) {
	serializedDeployItem := []byte{
		1,
		113, 29, 198, 74, 172, 204, 98, 45, 244, 151, 41, 226, 67, 58, 230, 46, 223, 254, 224, 124, 110, 137, 119, 203, 109, 96, 91, 27, 120, 151, 46, 113,
		4, 0, 0, 0, 116, 101, 115, 116,
		1, 0, 0, 0, 6, 0, 0, 0, 97, 109, 111, 117, 110, 116,
		4, 0, 0, 0, 3, 128, 150, 152, 8}
	hash, err := keypair.HashFromHex("711dc64aaccc622df49729e2433ae62edffee07c6e8977cb6d605b1b78972e71")
	assert.NoError(t, err)

	deployItem := ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByHash,
		StoredContractByHash: &StoredContractByHash{
			Hash:       hash,
			Entrypoint: "test",
			Args: []NamedArg{{
				Name: "amount",
				Value: types.CLValue{
					Type: types.CLTypeU512,
					U512: big.NewInt(10000000),
				},
			}},
		},
	}

	AssertEqualSerialization(t, deployItem, serializedDeployItem)
}

func TestStoredContractByName_Marshal(t *testing.T) {
	serializedDeployItem := []byte{
		2, 7, 0, 0, 0, 101, 120, 97, 109, 112,
		108, 101, 4, 0, 0, 0, 116, 101, 115, 116,
		1, 0, 0, 0, 6, 0, 0, 0, 97, 109,
		111, 117, 110, 116, 4, 0, 0, 0, 3, 128,
		150, 152, 8}

	deployItem := ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredContractByName,
		StoredContractByName: &StoredContractByName{
			Name:       "example",
			Entrypoint: "test",
			Args: []NamedArg{{
				Name: "amount",
				Value: types.CLValue{
					Type: types.CLTypeU512,
					U512: big.NewInt(10000000),
				},
			}},
		},
	}

	AssertEqualSerialization(t, deployItem, serializedDeployItem)
}

func TestStoredVersionedContractByHash_Marshal(t *testing.T) {
	serializedDeployItem := []byte{
		3, 113, 29, 198, 74, 172, 204, 98, 45, 244, 151, 41,
		226, 67, 58, 230, 46, 223, 254, 224, 124, 110, 137, 119,
		203, 109, 96, 91, 27, 120, 151, 46, 113, 1, 1, 0,
		0, 0, 4, 0, 0, 0, 116, 101, 115, 116, 1, 0,
		0, 0, 6, 0, 0, 0, 97, 109, 111, 117, 110, 116,
		4, 0, 0, 0, 3, 128, 150, 152, 8}
	version := uint32(1)
	hash, err := keypair.HashFromHex("711dc64aaccc622df49729e2433ae62edffee07c6e8977cb6d605b1b78972e71")
	assert.NoError(t, err)

	deployItem := ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByHash,
		StoredVersionedContractByHash: &StoredVersionedContractByHash{
			Hash:       hash,
			Version:    &version,
			Entrypoint: "test",
			Args: []NamedArg{{
				Name: "amount",
				Value: types.CLValue{
					Type: types.CLTypeU512,
					U512: big.NewInt(10000000),
				},
			}},
		},
	}

	AssertEqualSerialization(t, deployItem, serializedDeployItem)
}

func TestStoredVersionedContractByName_Marshal(t *testing.T) {
	serializedDeployItem := []byte{
		4, 7, 0, 0, 0, 101, 120, 97, 109, 112, 108,
		101, 1, 1, 0, 0, 0, 4, 0, 0, 0, 116,
		101, 115, 116, 1, 0, 0, 0, 6, 0, 0, 0,
		97, 109, 111, 117, 110, 116, 4, 0, 0, 0, 3,
		128, 150, 152, 8}
	version := uint32(1)
	deployItem := ExecutableDeployItem{
		Type: ExecutableDeployItemTypeStoredVersionedContractByName,
		StoredVersionedContractByName: &StoredVersionedContractByName{
			Name:       "example",
			Version:    &version,
			Entrypoint: "test",
			Args: []NamedArg{{
				Name: "amount",
				Value: types.CLValue{
					Type: types.CLTypeU512,
					U512: big.NewInt(10000000),
				},
			}},
		},
	}

	AssertEqualSerialization(t, deployItem, serializedDeployItem)
}

func TestDeploy(t *testing.T) {
	account, err := keypair.FromPublicKeyHex("01e456c3779510fd14e83fa3be84ff4b2a22de76ef6be677ed7936f37f7712a0a4")
	assert.NoError(t, err)
	param := "test"

	timestamp, err := time.Parse(time.RFC3339Nano, "2021-09-21T14:58:41.048Z")
	assert.NoError(t, err)
	deployParams := DeployParams{
		Account:   account,
		Timestamp: timestamp,
		TTL:       time.Minute * 30,
		GasPrice:  1,
		ChainName: "casper-test",
	}
	payment := StandardPayment(big.NewInt(10000000))
	session := NewModuleBytes([]byte{1, 2, 3}, []NamedArg{{
		Name: "message",
		Value: types.CLValue{
			Type:   types.CLTypeString,
			String: &param,
		},
	}})

	expectedHash, err := keypair.HashFromHex("6463f022a7114a2ee92cbeefd563431d8d6e2f4efedb73f12e8d069c83777b25")
	assert.NoError(t, err)
	deploy, err := MakeDeploy(deployParams, payment, session)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, deploy.Hash)
}

func TestTransfer(t *testing.T) {
	account, err := keypair.FromPublicKeyHex("01e456c3779510fd14e83fa3be84ff4b2a22de76ef6be677ed7936f37f7712a0a4")
	assert.NoError(t, err)
	timestamp, err := time.Parse(time.RFC3339Nano, "2021-09-21T14:58:41.048Z")
	assert.NoError(t, err)
	target, err := keypair.FromPublicKeyHex("0172a54c123b336fb1d386bbdff450623d1b5da904f5e2523b3e347b6d7573ae80")
	assert.NoError(t, err)
	deployParams := DeployParams{
		Account:   account,
		Timestamp: timestamp,
		TTL:       time.Second,
		GasPrice:  1,
		ChainName: "casper-test",
	}
	payment := StandardPayment(big.NewInt(10000000000))
	session := NewTransfer(big.NewInt(25000000000), target, uint64(5))

	expectedHash, err := keypair.HashFromHex("c54103f6b97eb999c7b92ed80d681020dc1c506052eab23f5bac5fe65532e489")
	assert.NoError(t, err)
	deploy, err := MakeDeploy(deployParams, payment, session)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, deploy.Hash)
}

func AssertEqualSerialization(t *testing.T, item serialization.Marshaler, actual []byte) {
	t.Helper()
	var w bytes.Buffer
	encoder := serialization.NewEncoder(&w)
	_, err := encoder.Encode(item)
	assert.NoError(t, err)
	assert.Equal(t, actual, w.Bytes())
}
