package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"reflect"

	"github.com/Simplewallethq/casper-golang-sdk/serialization"
)

func UnmarshalCLValue(src []byte, dest *CLValue) (int, error) {
	r := bytes.NewReader(src)
	dec := serialization.NewDecoder(r)
	switch dest.Type {
	case CLTypeBool:
		//var temp bool
		temp, n, err := dec.DecodeBool()
		//err = serialization.Unmarshal(src, &temp)
		if err != nil {
			return 0, err
		}
		dest.Bool = &temp
		return n, nil
	case CLTypeI32:
		temp, n, err := dec.DecodeInt32()
		if err != nil {
			return 0, err
		}
		dest.I32 = &temp
		return n, nil
	case CLTypeI64:
		temp, n, err := dec.DecodeInt64()
		if err != nil {
			return 0, err
		}
		dest.I64 = &temp
		return n, nil
	case CLTypeU8:
		temp, n, err := dec.DecodeByte()
		if err != nil {
			return 0, err
		}
		dest.U8 = &temp
		return n, nil
	case CLTypeU32:
		temp, n, err := dec.DecodeUInt32()
		if err != nil {
			return 0, err
		}
		dest.U32 = &temp
		return n, nil
	case CLTypeU64:
		temp, n, err := dec.DecodeUInt64()
		if err != nil {
			return 0, err
		}
		dest.U64 = &temp
		return n, nil
	case CLTypeU128:
		temp := reflect.New(reflect.TypeOf(serialization.U128{})).Interface()
		n, err := dec.DecodeU128(reflect.ValueOf(temp).Elem())
		if err != nil {
			return 0, err
		}
		dest.U128 = &temp.(*serialization.U128).Int
		return n, nil
	case CLTypeU256:
		temp := reflect.New(reflect.TypeOf(serialization.U256{})).Interface()
		n, err := dec.DecodeU256(reflect.ValueOf(temp).Elem())
		if err != nil {
			return 0, err
		}
		dest.U256 = &temp.(*serialization.U256).Int
		return n, nil
	case CLTypeU512:
		temp := reflect.New(reflect.TypeOf(serialization.U512{})).Interface()
		n, err := dec.DecodeU512(reflect.ValueOf(temp).Elem())
		if err != nil {
			return 0, err
		}
		dest.U512 = &temp.(*serialization.U512).Int
		return n, nil
	case CLTypeUnit:
		return 0, nil
	case CLTypeString:
		temp, n, err := dec.DecodeString()
		if err != nil {
			return 0, err
		}
		dest.String = &temp
		return n, nil
	case CLTypeKey:
		if dest.Key == nil {
			return 0, errors.New("inproper destination, mission key clvalue type")
		}

		return dest.Key.Unmarshal(src[1:])
	case CLTypeURef:
		dest.URef = &URef{}
		return dest.URef.Unmarshal(src)
	case CLTypeOption:
		var temp CLValue
		var n int

		if src[0] == byte(0) {
			dest.Option = nil
			return 0, nil
		} else if src[0] == byte(1) {
			if dest.Option == nil {
				return 0, errors.New("inproper destination, mission option clvalue type")
			}
			temp.Type = dest.Option.Type
			n1, err := UnmarshalCLValue(src[1:], &temp)
			if err != nil {
				return 0, err
			}
			n += n1
		} else {
			return 0, errors.New("inproper bytes or dest type")
		}

		dest.Option = &temp
		return n, nil
	case CLTypeList:
		size, n, err := dec.DecodeUInt32()
		if err != nil {
			return 0, err
		}

		if size == 0 {
			list := make([]CLValue, 0)
			dest.List = &list
			return n, nil
		} else {
			if dest.List == nil || len(*dest.List) == 0 {
				return n, errors.New("inproper destination, missing first element of list with inner type")
			}

			list := make([]CLValue, 0)
			for i := uint32(0); i < size; i++ {
				var temp CLValue
				temp.Type = (*dest.List)[0].Type
				n1, err := UnmarshalCLValue(src[n:], &temp)
				if err != nil {
					return n, err
				}
				n += n1
				list = append(list, temp)
			}

			dest.List = &list
		}

		return n, nil
	case CLTypeByteArray:
		var temp FixedByteArray
		temp = src
		dest.ByteArray = &temp
		return len(src), nil
	case CLTypeResult:
		var temp CLValueResult
		var n int
		if src[0] == byte(0) {
			temp.IsSuccess = false
			resultError := CLValue{
				Type: CLTypeString,
			}
			n1, err := UnmarshalCLValue(src[1:], &resultError)
			if err != nil {
				return n, err
			}
			n += n1
			temp.Error = &resultError
		} else if src[0] == byte(1) {
			temp.IsSuccess = true
			if dest.Result == nil || dest.Result.Success == nil {
				return n, errors.New("inproper destination, missing result success inner type")
			}
			result := CLValue{Type: dest.Result.Success.Type}
			n1, err := UnmarshalCLValue(src[1:], &result)
			if err != nil {
				return n, err
			}
			n += n1
			temp.Success = &result
		} else {
			return n, errors.New("inproper bytes or dest type")
		}

		dest.Result = &temp
		return n, nil
	case CLTypeMap:
		if dest.Map == nil {
			return 0, errors.New("inproper destination, missing map inner types")
		}

		size, n, err := dec.DecodeUInt32()
		if err != nil {
			return n, err
		}

		raw := make(map[string]CLValue)

		for i := uint32(0); i < size; i++ {
			var key CLValue
			key.Type = dest.Map.KeyType
			n1, err := UnmarshalCLValue(src[n:], &key)
			if err != nil {
				return n, err
			}
			n += n1
			var value CLValue
			value.Type = dest.Map.ValueType
			n1, err = UnmarshalCLValue(src[n:], &value)
			if err != nil {
				return n, err
			}
			n += n1

			if key.Type != CLTypeString {
				marshaledKey, err := serialization.Marshal(key)
				if err != nil {
					return n, err
				}

				raw[hex.EncodeToString(marshaledKey)] = value
			} else {
				raw[*key.String] = value
			}
		}

		dest.Map.Raw = raw

		return n, nil
	case CLTypeTuple1:
		if dest.Tuple1 == nil {
			return 0, errors.New("inproper destination, missing tuple inner type")
		}
		temp := CLValue{Type: dest.Tuple1[0].Type}
		n, err := UnmarshalCLValue(src, &temp)
		if err != nil {
			return n, err
		}
		dest.Tuple1[0] = temp
		return n, nil
	case CLTypeTuple2:
		if dest.Tuple2 == nil {
			return 0, errors.New("inproper destination, missing tuple inner type")
		}
		temp := CLValue{Type: dest.Tuple2[0].Type}
		n, err := UnmarshalCLValue(src, &temp)
		if err != nil {
			return n, err
		}
		dest.Tuple2[0] = temp
		temp = CLValue{Type: dest.Tuple2[1].Type}
		n1, err := UnmarshalCLValue(src[n:], &temp)
		if err != nil {
			return n, err
		}
		n += n1
		dest.Tuple2[1] = temp
		return n, nil
	case CLTypeTuple3:
		if dest.Tuple3 == nil {
			return 0, errors.New("inproper destination, missing tuple inner type")
		}
		temp := CLValue{Type: dest.Tuple3[0].Type}
		n, err := UnmarshalCLValue(src, &temp)
		if err != nil {
			return n, err
		}
		dest.Tuple3[0] = temp

		temp = CLValue{Type: dest.Tuple3[1].Type}
		n1, err := UnmarshalCLValue(src[n:], &temp)
		if err != nil {
			return n, err
		}
		n += n1
		dest.Tuple3[1] = temp

		temp = CLValue{Type: dest.Tuple3[2].Type}
		n1, err = UnmarshalCLValue(src[n:], &temp)
		if err != nil {
			return n, err
		}
		n += n1
		dest.Tuple3[2] = temp
		return n, nil
	}

	return 0, errors.New("Invalid CLtype provided")
}
