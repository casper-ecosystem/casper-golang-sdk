package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/casper-ecosystem/casper-golang-sdk/serialization"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	Name            string
	Value           CLValue
	HexEncodedValue string
}{
	{
		"Bool",
		CLValue{Type: CLTypeBool, Bool: createPtrBool(true)},
		"01",
	},
	{
		"I32",
		CLValue{Type: CLTypeI32, I32: createPtrI32(7)},
		"07000000",
	},
	{
		"I64",
		CLValue{Type: CLTypeI64, I64: createPtrI64(7)},
		"0700000000000000",
	},
	{
		"U8",
		CLValue{Type: CLTypeU8, U8: createPtrU8(7)},
		"07",
	},
	{
		"U32",
		CLValue{Type: CLTypeU32, U32: createPtrU32(7)},
		"07000000",
	},
	{
		"U64",
		CLValue{Type: CLTypeU64, U64: createPtrU64(1024)},
		"0004000000000000",
	},
	{
		"U128",
		CLValue{Type: CLTypeU128, U128: big.NewInt(7)},
		"0107",
	},
	{
		"U256",
		CLValue{Type: CLTypeU256, U256: big.NewInt(7)},
		"0107",
	},
	{
		"U512",
		CLValue{Type: CLTypeU512, U512: big.NewInt(7)},
		"0107",
	},
	{
		"Unit",
		CLValue{Type: CLTypeUnit},
		"",
	},
	{
		"String",
		CLValue{Type: CLTypeString, String: createPtrString("Hello, World!")},
		"0d00000048656c6c6f2c20576f726c6421",
	},
	{
		"Key, address",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeAccount, Account: getAddress()}},
		"004c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, hash",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeHash, Hash: getAddress()}},
		"014c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, uref",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeURef, URef: &URef{
			AccessRight: AccessRightReadAddWrite,
			Address:     getAddress(),
		}}},
		"024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507",
	},
	{
		"Key, transfer",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeTransfer, Transfer: getAddress()}},
		"034c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, deploy info",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeDeployInfo, DeployInfo: getAddress()}},
		"044c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, eraid",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeEraId, EraId: toPtrU64(1024)}},
		"050004000000000000",
	},
	{
		"Key, balance",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeBalance, Balance: getAddress()}},
		"064c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, bid",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeBid, Bid: getAddress()}},
		"074c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Key, withdraw",
		CLValue{Type: CLTypeKey, Key: &Key{Type: KeyTypeWithdraw, Withdraw: getAddress()}},
		"084c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5",
	},
	{
		"Optional nil",
		CLValue{Type: CLTypeOption, Option: nil},
		"00",
	},
	{
		"Optional U64",
		CLValue{Type: CLTypeOption, Option: &CLValue{Type: CLTypeU64, U64: createPtrU64(7)}},
		"010700000000000000",
	},
	{
		"URef",
		CLValue{Type: CLTypeURef, URef: &URef{
			AccessRight: AccessRightReadAddWrite,
			Address:     getAddress(),
		}},
		"4c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507",
	},
	{
		"Empty list",
		CLValue{Type: CLTypeList, List: createPtrList(make([]CLValue, 0))},
		"00000000",
	},
	{
		"List",
		CLValue{Type: CLTypeList, List: &[]CLValue{
			{Type: CLTypeU32,
				U32: createPtrU32(1)},
			{Type: CLTypeU32,
				U32: createPtrU32(2)},
			{Type: CLTypeU32,
				U32: createPtrU32(3)},
		}},
		"03000000010000000200000003000000",
	},
	{
		"Fixed size list",
		CLValue{
			Type:      CLTypeByteArray,
			ByteArray: getByteArray(),
		},
		"010000000200000003000000",
	},
	{
		"Result success",
		CLValue{
			Type: CLTypeResult,
			Result: &CLValueResult{
				IsSuccess: true,
				Success: &CLValue{
					Type: CLTypeU64,
					U64:  toPtrU64(314),
				},
				Error: nil,
			},
		},
		"013a01000000000000",
	},
	{
		"Result error",
		CLValue{
			Type: CLTypeResult,
			Result: &CLValueResult{
				IsSuccess: false,
				Success:   nil,
				Error: &CLValue{
					Type:   CLTypeString,
					String: toPtrString("Uh oh"),
				},
			},
		},
		"00050000005568206f68",
	},
	{
		"map[string]string",
		CLValue{
			Type: CLTypeMap,
			Map:  getMap(),
		},
		//map[string]string{
		//	"test": "test",
		//},
		"0100000004000000746573740400000074657374",
	},
	{
		"Tuple(u32)",
		CLValue{
			Type: CLTypeTuple1,
			Tuple1: &[1]CLValue{
				{
					Type: CLTypeU32,
					U32:  createPtrU32(1),
				},
			},
		},
		"01000000",
	},
	{
		"Tuple(u32, string)",
		CLValue{
			Type: CLTypeTuple2,
			Tuple2: &[2]CLValue{
				{
					Type: CLTypeU32,
					U32:  createPtrU32(1),
				},
				{
					Type:   CLTypeString,
					String: createPtrString("Hello, World!"),
				},
			},
		},
		"010000000d00000048656c6c6f2c20576f726c6421",
	},
	{
		"Tuple(u32, string,bool)",
		CLValue{
			Type: CLTypeTuple3,
			Tuple3: &[3]CLValue{
				{
					Type: CLTypeU32,
					U32:  createPtrU32(1),
				},
				{
					Type:   CLTypeString,
					String: createPtrString("Hello, World!"),
				},
				{
					Type: CLTypeBool,
					Bool: createPtrBool(true),
				},
			},
		},
		"010000000d00000048656c6c6f2c20576f726c642101",
	},
}

func toPtrU64(val uint64) *uint64 {
	return &val
}

func toPtrString(val string) *string {
	return &val
}

func TestEncoding(t *testing.T) {
	for _, c := range cases {
		result, err := serialization.Marshal(c.Value)
		if err != nil {
			t.Errorf("encoding failed: failed to marshal %s: %v", c.Name, err)
		}

		hexResult := hex.EncodeToString(result)

		assert.Equal(t, c.HexEncodedValue, hexResult, "encoding failed: failed to marshal %s", c.Name)
	}
}

func TestDecoding(t *testing.T) {
	keyTypeCounter := 0
	for _, c := range cases {
		buf, err := hex.DecodeString(c.HexEncodedValue)
		if err != nil {
			t.Errorf("decoding failed: failed to decode hex for %s: %v", c.Name, err)
		}
		result := CLValue{Type: c.Value.Type}

		switch c.Value.Type {
		case CLTypeOption:
			if c.Value.Option != nil {
				result.Option = &CLValue{Type: c.Value.Option.Type}
			}
			break
		case CLTypeList:
			if len(*c.Value.List) != 0 {
				list := make([]CLValue, 1)
				list[0] = CLValue{Type: (*c.Value.List)[0].Type}
				result.List = &list
			}
			break
		case CLTypeResult:
			if c.Value.Result.IsSuccess == true {
				result.Result = &CLValueResult{
					Success: &CLValue{Type: c.Value.Result.Success.Type},
				}
			}
			break
		case CLTypeMap:
			result.Map = &CLMap{
				KeyType:   c.Value.Map.KeyType,
				ValueType: c.Value.Map.ValueType,
			}
			break
		case CLTypeTuple1:
			result.Tuple1 = &[1]CLValue{{Type: CLTypeU32}}
			break
		case CLTypeTuple2:
			result.Tuple2 = &[2]CLValue{{Type: CLTypeU32}, {Type: CLTypeString}}
			break
		case CLTypeTuple3:
			result.Tuple3 = &[3]CLValue{{Type: CLTypeU32}, {Type: CLTypeString}, {Type: CLTypeBool}}
			break
		case CLTypeKey:
			switch keyTypeCounter {
			case 0:
				result.Key = &Key{Type: KeyTypeAccount}
			case 1:
				result.Key = &Key{Type: KeyTypeHash}
			case 2:
				result.Key = &Key{Type: KeyTypeURef}
			case 3:
				result.Key = &Key{Type: KeyTypeTransfer}
			case 4:
				result.Key = &Key{Type: KeyTypeDeployInfo}
			case 5:
				result.Key = &Key{Type: KeyTypeEraId}
			case 6:
				result.Key = &Key{Type: KeyTypeBalance}
			case 7:
				result.Key = &Key{Type: KeyTypeBid}
			case 8:
				result.Key = &Key{Type: KeyTypeWithdraw}
			}

			keyTypeCounter++
		}

		_, err = UnmarshalCLValue(buf, &result)
		if err != nil {
			t.Errorf("decoding failed: failed to unmarshal %s: %v", c.Name, err)
		}

		assert.Equal(t, c.Value, result)

	}
}

func createPtrBool(v bool) *bool {
	return &v
}

func createPtrI32(v int32) *int32 {
	return &v
}

func createPtrI64(v int64) *int64 {
	return &v
}

func createPtrU8(v uint8) *uint8 {
	return &v
}

func createPtrU32(v uint32) *uint32 {
	return &v
}

func createPtrU64(v uint64) *uint64 {
	return &v
}

func createPtrString(v string) *string {
	return &v
}

func createPtrList(v []CLValue) *[]CLValue {
	return &v
}

func getByteArray() *FixedByteArray {
	marshal, err := serialization.Marshal([3]CLValue{
		{Type: CLTypeU32,
			U32: createPtrU32(1)},
		{Type: CLTypeU32,
			U32: createPtrU32(2)},
		{Type: CLTypeU32,
			U32: createPtrU32(3)},
	})
	if err != nil {
		return nil
	}

	res := FixedByteArray(marshal)

	return &res
}

func getMap() *CLMap {
	//keys := make([]CLValue,1)
	//values := make([]CLValue,1)
	value := CLValue{Type: CLTypeString, String: toPtrString("test")}
	raw := make(map[string]CLValue)
	raw["test"] = value
	return &CLMap{
		KeyType:   CLTypeString,
		ValueType: CLTypeString,
		Raw:       raw,
	}
}

func getAddress() [32]byte {
	decodeString, err := hex.DecodeString("4c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5")
	if err != nil {
		return [32]byte{}
	}

	var result [32]byte

	for i := 0; i < 32; i++ {
		result[i] = decodeString[i]
	}

	return result
}
