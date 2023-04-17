package types

import (
	"math/big"

	"github.com/Simplewallethq/casper-golang-sdk/keypair"
)

type CLType byte

const (
	CLTypeBool CLType = iota
	CLTypeI32
	CLTypeI64
	CLTypeU8
	CLTypeU32
	CLTypeU64
	CLTypeU128
	CLTypeU256
	CLTypeU512
	CLTypeUnit
	CLTypeString
	CLTypeKey
	CLTypeURef
	CLTypeOption
	CLTypeList
	CLTypeByteArray
	CLTypeResult
	CLTypeMap
	CLTypeTuple1
	CLTypeTuple2
	CLTypeTuple3
	CLTypeAny
	CLTypePublicKey
)

func (t CLType) ToString() string {
	switch t {
	case CLTypeBool:
		return "Bool"
	case CLTypeI32:
		return "I32"
	case CLTypeI64:
		return "I64"
	case CLTypeU8:
		return "U8"
	case CLTypeU32:
		return "U32"
	case CLTypeU64:
		return "U64"
	case CLTypeU128:
		return "U128"
	case CLTypeU256:
		return "U256"
	case CLTypeU512:
		return "U512"
	case CLTypeUnit:
		return "Unit"
	case CLTypeString:
		return "String"
	case CLTypeKey:
		return "Key"
	case CLTypeURef:
		return "URef"
	case CLTypeOption:
		return "Option"
	case CLTypeList:
		return "List"
	case CLTypeByteArray:
		return "ByteArray"
	case CLTypeResult:
		return "Result"
	case CLTypeMap:
		return "Map"
	case CLTypeTuple1:
		return "Tuple1"
	case CLTypeTuple2:
		return "Tuple2"
	case CLTypeTuple3:
		return "Tuple3"
	case CLTypePublicKey:
		return "PublicKey"
	}

	return "Any"
}

func FromString(clType string) CLType {
	switch clType {
	case "Bool":
		return CLTypeBool
	case "I32":
		return CLTypeI32
	case "I64":
		return CLTypeI64
	case "U8":
		return CLTypeU8
	case "U32":
		return CLTypeU32
	case "U64":
		return CLTypeU64
	case "U128":
		return CLTypeU128
	case "U256":
		return CLTypeU256
	case "U512":
		return CLTypeU512
	case "Unit":
		return CLTypeUnit
	case "String":
		return CLTypeString
	case "Key":
		return CLTypeKey
	case "URef":
		return CLTypeURef
	case "Option":
		return CLTypeOption
	case "List":
		return CLTypeList
	case "ByteArray":
		return CLTypeByteArray
	case "Result":
		return CLTypeResult
	case "Map":
		return CLTypeMap
	case "Tuple1":
		return CLTypeTuple1
	case "Tuple2":
		return CLTypeTuple2
	case "Tuple3":
		return CLTypeTuple3
	case "PublicKey":
		return CLTypePublicKey
	}

	return CLTypeAny
}

type CLValue struct {
	Type      CLType
	Bool      *bool
	I32       *int32
	I64       *int64
	U8        *byte
	U32       *uint32
	U64       *uint64
	U128      *big.Int
	U256      *big.Int
	U512      *big.Int
	String    *string
	Key       *Key
	URef      *URef
	Option    *CLValue
	List      *[]CLValue
	ByteArray *FixedByteArray
	Result    *CLValueResult
	Map       *CLMap
	Tuple1    *[1]CLValue
	Tuple2    *[2]CLValue
	Tuple3    *[3]CLValue
	PublicKey *keypair.PublicKey
}

func (v CLValue) SwitchFieldName() string {
	return "Type"
}

func (v CLValue) ArmForSwitch(sw byte) (string, bool) {
	switch CLType(sw) {
	case CLTypeBool:
		return "Bool", true
	case CLTypeI32:
		return "I32", true
	case CLTypeI64:
		return "I64", true
	case CLTypeU8:
		return "U8", true
	case CLTypeU32:
		return "U32", true
	case CLTypeU64:
		return "U64", true
	case CLTypeU128:
		return "U128", true
	case CLTypeU256:
		return "U256", true
	case CLTypeU512:
		return "U512", true
	case CLTypeUnit:
		return "", false
	case CLTypeString:
		return "String", true
	case CLTypeKey:
		return "Key", true
	case CLTypeURef:
		return "URef", true
	case CLTypeOption:
		return "Option", true
	case CLTypeList:
		return "List", true
	case CLTypeByteArray:
		return "ByteArray", true
	case CLTypeResult:
		return "Result", true
	case CLTypeMap:
		return "Map", true
	case CLTypeTuple1:
		return "Tuple1", true
	case CLTypeTuple2:
		return "Tuple2", true
	case CLTypeTuple3:
		return "Tuple3", true
	case CLTypeAny:
		return "-", false
	case CLTypePublicKey:
		return "PublicKey", false
	}

	return "-", false
}
