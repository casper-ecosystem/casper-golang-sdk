package sdk

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/Simplewallethq/casper-golang-sdk/serialization"
	"github.com/Simplewallethq/casper-golang-sdk/types"
)

type ValueMap struct {
	KeyType   types.CLType
	ValueType types.CLType
}

type Value struct {
	Tag         types.CLType
	IsOptional  bool
	StringBytes string
	Optional    *Value
	Map         *ValueMap
}

func (v Value) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	if v.IsOptional {
		data["cl_type"] = map[string]interface{}{"Option": v.Optional.Tag.ToString()}
		//data["bytes"] = hex.EncodeToString(v.Bytes)

		data["bytes"] = v.Optional.StringBytes
	} else if v.Tag == types.CLTypeMap {
		data["cl_type"] = map[string]interface{}{
			"Map": map[string]interface{}{
				"key":   v.Map.KeyType.ToString(),
				"value": v.Map.ValueType.ToString(),
			},
		}
		data["bytes"] = v.StringBytes
	} else {
		if v.Tag == types.CLTypeByteArray {
			data["cl_type"] = map[string]interface{}{"ByteArray": len(v.StringBytes) / 2}
		} else {
			data["cl_type"] = v.Tag.ToString()
		}

		data["bytes"] = v.StringBytes
	}

	return json.Marshal(data)
}

func (v Value) ToBytes() []uint8 {
	res := make([]uint8, 0)

	var bytesSerialized []byte
	var encodedBytesSerialized []byte
	var err error

	if v.IsOptional {
		decodeString, err := hex.DecodeString(v.Optional.StringBytes)
		if err != nil {
			return nil
		}
		marshaledOptional, err2 := serialization.Marshal(decodeString)
		if err2 != nil {
			return nil
		}

		res = append(res, marshaledOptional...)
		//bytesSerialized = v.Bytes
		//encodedBytesSerialized = v.Bytes
	} else {
		bytesSerialized, err = hex.DecodeString(v.StringBytes)
		if err != nil {
			return nil
		}
		encodedBytesSerialized, err = serialization.Marshal(bytesSerialized)
		if err != nil {
			return nil
		}
	}

	res = append(res, encodedBytesSerialized...)

	res = append(res, uint8(v.Tag))

	if v.Tag == types.CLTypeByteArray {
		lenByteArray, err := serialization.Marshal(int32(len(bytesSerialized)))

		if err != nil {
			return nil
		}

		res = append(res, lenByteArray...)
	} else if v.Tag == types.CLTypeOption {
		res = append(res, uint8(v.Optional.Tag))
	}

	return res
}

type RuntimeArgs struct {
	KeyOrder []string
	Args     map[string]Value
}

func ParseRuntimeArgs(args []interface{}) (RuntimeArgs, error) {
	result := make(map[string]Value)
	order := make([]string, 0)

	for _, val := range args {
		argument, ok := val.([]interface{})

		if !ok || len(argument) != 2 {
			return RuntimeArgs{}, errors.New("invalid json")
		}

		valueParsed, ok := argument[1].(map[string]interface{})

		if !ok {
			return RuntimeArgs{}, errors.New("invalid json")
		}

		if valueParsed["bytes"] == nil {
			return RuntimeArgs{}, errors.New("'bytes' key doesn't exist, invalid json")
		}
		if valueParsed["cl_type"] == nil {
			return RuntimeArgs{}, errors.New("'bytes' key doesn't exist, invalid json")
		}

		cl_type, ok := valueParsed["cl_type"].(string)

		var value Value

		if !ok {
			cl_type, ok := valueParsed["cl_type"].(map[string]interface{})
			if !ok {
				return RuntimeArgs{}, errors.New("failed parse cl_type")
			}

			if cl_type["Option"] != nil {
				value = Value{
					Tag:        types.CLTypeOption,
					IsOptional: true,
					Optional: &Value{
						Tag:         types.FromString(cl_type["Option"].(string)),
						StringBytes: valueParsed["bytes"].(string),
					},
				}
			} else if cl_type["ByteArray"] != nil {
				value = Value{
					Tag:         types.CLTypeByteArray,
					StringBytes: valueParsed["bytes"].(string),
				}
			} else if cl_type["Map"] != nil {
				mapTypes, ok := cl_type["Map"].(map[string]interface{})

				if !ok {
					return RuntimeArgs{}, errors.New("unable to unmarshal, not map types")
				}

				value = Value{
					Tag:         types.CLTypeMap,
					StringBytes: valueParsed["bytes"].(string),
					Map: &ValueMap{
						KeyType:   types.FromString(mapTypes["key"].(string)),
						ValueType: types.FromString(mapTypes["value"].(string)),
					},
				}
			} else {
				return RuntimeArgs{}, errors.New("unknown type")
			}
		} else {
			value = Value{
				Tag:         types.FromString(cl_type),
				StringBytes: valueParsed["bytes"].(string),
			}
		}

		result[argument[0].(string)] = value
		order = append(order, argument[0].(string))
	}

	return *NewRunTimeArgs(result, order), nil
}

func NewRunTimeArgs(args map[string]Value, keyOrder []string) *RuntimeArgs {
	r := new(RuntimeArgs)
	r.Args = args
	r.KeyOrder = keyOrder
	return r
}

func (r RuntimeArgs) FromMap(args map[string]Value, keyOrder []string) *RuntimeArgs {
	return NewRunTimeArgs(args, keyOrder)
}

func (r RuntimeArgs) Insert(key string, value Value) {
	r.Args[key] = value
	r.KeyOrder = append(r.KeyOrder, key)
}

func (r RuntimeArgs) ToBytes() []byte {
	res := make([]byte, 0)

	encodedLen, err := serialization.Marshal(int32(len(r.Args)))
	if err != nil {
		return nil
	}

	res = append(res, encodedLen...)

	for _, k := range r.KeyOrder {
		encodedKey, err := serialization.Marshal(k)
		if err != nil {
			return nil
		}
		res = append(res, encodedKey...)
		res = append(res, r.Args[k].ToBytes()...)
	}

	return res
}

// runtime args in deploy are encoded as an array of pairs "key" : "value"
func (r RuntimeArgs) ToJSONInterface() []interface{} {
	result := make([]interface{}, 0)

	for _, k := range r.KeyOrder {
		temp := make([]interface{}, 2)
		temp[0] = k
		temp[1] = r.Args[k]
		result = append(result, temp)
	}

	return result
}
