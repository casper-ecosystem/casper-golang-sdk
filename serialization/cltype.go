package serialization

import "io"

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
	Type  CLType
	Value interface{}
}

func (v CLValue) Marshal(w io.Writer) (int, error) {
	enc := Encoder{
		w: w,
	}

	n, err := enc.EncodeByte(byte(v.Type))
	if err != nil {
		return n, err
	}

	n2, err := enc.Encode(v.Value)
	n += n2

	return n, err
}
