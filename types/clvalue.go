package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"io"
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
	Optional  *Optional
	List      []CLValue
	ByteArray FixedByteArray
	Result    *CLValueResult
	Map       *CLMap
	Tuple1    *[1]CLValue
	Tuple2    *[2]CLValue
	Tuple3    *[3]CLValue
	PublicKey *PublicKey
	Any       *interface{}
}

type Optional struct {
	*CLValue
}

func (o Optional) Marshal(w io.Writer) (int, error) {
	size := 0
	encoder := serialization.NewEncoder(w)

	if o.CLValue == nil {
		n, err := encoder.Encode(byte(0))
		if err != nil {
			return 0, err
		}
		size += n
	} else {
		n, err := encoder.Encode(byte(1))
		if err != nil {
			return 0, err
		}
		size += n
		n, err = o.CLValue.Marshal(w)
		if err != nil {
			return 0, err
		}
		size += n
	}
	return size, nil
}

func (o Optional) MarshalJSON() ([]byte, error) {
	_, typeValue := o.GetValue()
	return json.Marshal(typeValue)
}

func (c *CLValue) GetValue() (string, interface{}) {
	switch c.Type {
	case CLTypeBool:
		return "Bool", *c.Bool
	case CLTypeI32:
		return "I32", *c.I32
	case CLTypeI64:
		return "I64", *c.I64
	case CLTypeU8:
		return "U8", *c.U8
	case CLTypeU32:
		return "U32", *c.U32
	case CLTypeU64:
		return "U64", *c.U64
	case CLTypeU128:
		return "U128", *c.U128
	case CLTypeU256:
		return "U256", *c.U256
	case CLTypeU512:
		return "U512", *c.U512
	case CLTypeUnit:
		return "", nil
	case CLTypeString:
		return "String", *c.String
	case CLTypeKey:
		return "Key", *c.Key
	case CLTypeURef:
		return "URef", *c.URef
	case CLTypeOption:
		return "Optional", c.Optional
	case CLTypeList:
		return "List", c.List
	case CLTypeByteArray:
		return "ByteArray", c.ByteArray
	case CLTypeResult:
		return "Result", *c.Result
	case CLTypeMap:
		return "Map", *c.Map
	case CLTypeTuple1:
		return "Tuple1", *c.Tuple1
	case CLTypeTuple2:
		return "Tuple2", *c.Tuple2
	case CLTypeTuple3:
		return "Tuple3", *c.Tuple3
	case CLTypeAny:
		return "Any", *c.Any
	case CLTypePublicKey:
		return "PublicKey", *c.PublicKey
	}
	return "", nil
}

// TODO: Union methods

func (c CLValue) MarshalJSON() ([]byte, error) {
	typeName, typeValue := c.GetValue()
	bytes, err := serialization.Marshal(typeValue)
	if err != nil {
		return nil, err
	}
	var typeJSON interface{}
	if typeName == "Optional" {
		optionalType, _ := c.Optional.GetValue()
		typeJSON = struct {
			Option string `json:"Option"`
		}{
			optionalType,
		}
	} else if typeName == "ByteArray" {
		typeJSON = struct {
			ByteArray int `json:"ByteArray"`
		}{
			len(c.ByteArray),
		}
	} else {
		typeJSON = typeName
	}
	return json.Marshal(struct {
		Type   interface{} `json:"cl_type"`
		Bytes  string      `json:"bytes"`
		Parsed interface{} `json:"parsed"`
	}{
		Type:   typeJSON,
		Bytes:  hex.EncodeToString(bytes),
		Parsed: typeValue,
	})
}

func (c CLValue) Marshal(w io.Writer) (int, error) {
	_, typeValue := c.GetValue()
	encoder := serialization.NewEncoder(w)
	n, err := encoder.Encode(typeValue)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c CLValue) MarshalWithType(w io.Writer) (int, error) {
	size := 0
	typeName, typeValue := c.GetValue()
	encoder := serialization.NewEncoder(w)
	valueBytes, err := serialization.Marshal(typeValue)
	if err != nil {
		return 0, err
	}
	n, err := encoder.Encode(valueBytes)
	if err != nil {
		return 0, err
	}
	size += n

	n, err = encoder.Encode(c.Type)
	if err != nil {
		return 0, err
	}
	size += n

	if typeName == "Optional" {
		n, err = encoder.Encode(c.Optional.Type)
		if err != nil {
			return 0, err
		}
		size += n
	} else if typeName == "ByteArray" {
		n, err = encoder.Encode(uint32(len(c.ByteArray)))
		if err != nil {
			return 0, err
		}
		size += n
	}
	return size, nil
}
