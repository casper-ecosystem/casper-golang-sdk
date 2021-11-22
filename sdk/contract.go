package sdk

import (
	"encoding/hex"
	"io/ioutil"
	"math/big"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
)

type Contract struct {
	SessionWasm []byte
	PaymentWasm []byte
}

func NewContract(sessionContractPath, paymentContractPath string) Contract {
	var result Contract
	sessionWasm, err := ioutil.ReadFile(sessionContractPath)
	if err != nil {
		return Contract{}
	}

	result.SessionWasm = sessionWasm

	if paymentContractPath != "" {
		payment, err := ioutil.ReadFile(paymentContractPath)
		if err != nil {
			return Contract{}
		}
		result.PaymentWasm = payment
	}

	return result
}

func (c Contract) Deploy(args RuntimeArgs, paymentAmount big.Int, pubKey keypair.PublicKey, keyPair keypair.KeyPair, chainName string) *Deploy {
	session := NewModuleBytes(c.SessionWasm, args)
	marshaledPaymentAmount, err := serialization.Marshal(paymentAmount)

	if err != nil {
		return nil
	}

	payment := NewModuleBytes(c.PaymentWasm, RuntimeArgs{Args: map[string]Value{"amount": {
		Tag:         types.CLTypeU512,
		StringBytes: hex.EncodeToString(marshaledPaymentAmount),
	}}})

	deploy := MakeDeploy(NewDeployParams(pubKey, chainName, nil, 0), payment, session)

	deploy.SignDeploy(keyPair)
	return deploy
}

type BoundContract struct {
	ContractStruct Contract
	KeyPair        keypair.KeyPair
}

func (b BoundContract) Deploy(args RuntimeArgs, paymentAmount big.Int, chainName string) *Deploy {
	return b.ContractStruct.Deploy(args, paymentAmount, b.KeyPair.PublicKey(), b.KeyPair, chainName)
}

type FaucetContract struct{}

func (f FaucetContract) MakeArgs(accountHash string) RuntimeArgs {
	decodedHash, err := hex.DecodeString(accountHash)
	if err != nil || len(decodedHash) != 32 {
		return RuntimeArgs{}
	}

	key := types.Key{
		Type:    types.KeyTypeAccount,
		Account: [32]byte{},
	}

	for i := 0; i < 32; i++ {
		key.Account[i] = decodedHash[i]
	}

	marshaledKey, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &key})

	if err != nil {
		return RuntimeArgs{}
	}

	return RuntimeArgs{
		Args: map[string]Value{"account": {
			Tag:         types.CLTypeKey,
			StringBytes: hex.EncodeToString(marshaledKey),
		}},
		KeyOrder: []string{"account"},
	}
}

type TransferContract struct{}

func (t TransferContract) MakeArgs(accountHash string, amount big.Int) RuntimeArgs {
	decodedHash, err := hex.DecodeString(accountHash)
	if err != nil || len(decodedHash) != 32 {
		return RuntimeArgs{}
	}

	key := types.Key{
		Type:    types.KeyTypeAccount,
		Account: [32]byte{},
	}

	for i := 0; i < 32; i++ {
		key.Account[i] = decodedHash[i]
	}

	marshaledKey, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &key})

	if err != nil {
		return RuntimeArgs{}
	}

	marshaledAmount, err := serialization.Marshal(types.CLValue{Type: types.CLTypeU512, U512: &amount})

	if err != nil {
		return RuntimeArgs{}
	}

	return RuntimeArgs{
		Args: map[string]Value{
			"account": {
				Tag:         types.CLTypeKey,
				StringBytes: hex.EncodeToString(marshaledKey),
			},
			"amount": {
				Tag:         types.CLTypeU512,
				StringBytes: hex.EncodeToString(marshaledAmount),
			},
		},
		KeyOrder: []string{"account", "amount"},
	}
}
