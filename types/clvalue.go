package types

import (
	"math/big"
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
	Optional  *CLValue
	List      []CLValue
	ByteArray FixedByteArray
	Result    *CLValueResult
	Map       *CLMap
	Tuple1    *[1]CLValue
	Tuple2    *[2]CLValue
	Tuple3    *[3]CLValue
	PublicKey *PublicKey
}

// TODO: Union methods
