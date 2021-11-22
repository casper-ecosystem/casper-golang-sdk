package types

import (
	"bytes"
	"io"

	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
)

type KeyType byte

const (
	KeyTypeAccount KeyType = iota
	KeyTypeHash
	KeyTypeURef
	KeyTypeTransfer
	KeyTypeDeployInfo
	KeyTypeEraId
	KeyTypeBalance
	KeyTypeBid
	KeyTypeWithdraw
)

// Key represents key structure
type Key struct {
	Type       KeyType
	Account    [32]byte
	Hash       [32]byte
	URef       *URef
	Transfer   [32]byte
	DeployInfo [32]byte
	EraId      *uint64
	Balance    [32]byte
	Bid        [32]byte
	Withdraw   [32]byte
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
	case KeyTypeTransfer:
		return "Transfer", true
	case KeyTypeDeployInfo:
		return "DeployInfo", true
	case KeyTypeEraId:
		return "EraId", true
	case KeyTypeBalance:
		return "Balance", true
	case KeyTypeBid:
		return "Bid", true
	case KeyTypeWithdraw:
		return "Withdraw", true
	}
	return "-", false
}

func (u Key) Marshal(w io.Writer) (int, error) {
	toMarshal := make([]byte, 1)
	toMarshal[0] = byte(u.Type)
	switch u.Type {
	case KeyTypeAccount:
		toMarshal = append(toMarshal, u.Account[:]...)
	case KeyTypeHash:
		toMarshal = append(toMarshal, u.Hash[:]...)
	case KeyTypeURef:
		mashaledURef, err := serialization.Marshal(u.URef)
		if err != nil {
			return 0, err
		}

		toMarshal = append(toMarshal, mashaledURef...)
	case KeyTypeTransfer:
		toMarshal = append(toMarshal, u.Transfer[:]...)
	case KeyTypeDeployInfo:
		toMarshal = append(toMarshal, u.DeployInfo[:]...)
	case KeyTypeEraId:
		mashaledEraId, err := serialization.Marshal(*u.EraId)
		if err != nil {
			return 0, err
		}

		toMarshal = append(toMarshal, mashaledEraId...)
	case KeyTypeBalance:
		toMarshal = append(toMarshal, u.Balance[:]...)
	case KeyTypeBid:
		toMarshal = append(toMarshal, u.Bid[:]...)
	case KeyTypeWithdraw:
		toMarshal = append(toMarshal, u.Withdraw[:]...)
	}

	return w.Write(toMarshal)
}

func (u *Key) Unmarshal(data []byte) (int, error) {
	r := bytes.NewReader(data)
	dec := serialization.NewDecoder(r)

	switch u.Type {
	case KeyTypeAccount:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Account = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeHash:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Hash = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeURef:
		u.URef = &URef{}
		return u.URef.Unmarshal(data)
	case KeyTypeTransfer:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Transfer = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeDeployInfo:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.DeployInfo = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeEraId:
		tempInt, n, err := dec.DecodeUInt64()
		if err != nil {
			return n, err
		}

		u.EraId = &tempInt
		return n, nil
	case KeyTypeBalance:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Balance = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeBid:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Bid = bytesTo32byte(destByteArray)

		return i, nil
	case KeyTypeWithdraw:
		destByteArray, i, err := dec.DecodeFixedByteArray(32)
		if err != nil {
			return 0, err
		}

		u.Withdraw = bytesTo32byte(destByteArray)

		return i, nil
	}

	return 0, nil
}

func bytesTo32byte(input []byte) [32]byte {
	if len(input) < 32 {
		return [32]byte{}
	}

	var result [32]byte

	for i := 0; i < 32; i++ {
		result[i] = input[i]
	}

	return result
}
