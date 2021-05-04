package serialization

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	Name            string
	Value           interface{}
	HexEncodedValue string
}{
	{
		"Bool",
		true,
		"01",
	},
	{
		"I32",
		int32(7),
		"07000000",
	},
	{
		"I64",
		int64(7),
		"0700000000000000",
	},
	{
		"U8",
		byte(7),
		"07",
	},
	{
		"U32",
		uint32(7),
		"07000000",
	},
	{
		"U64",
		uint64(1024),
		"0004000000000000",
	},
	{
		"U128",
		U128{*big.NewInt(7)},
		"0107",
	},
	{
		"U256",
		U256{*big.NewInt(7)},
		"0107",
	},
	{
		"U512",
		U512{*big.NewInt(7)},
		"0107",
	},
	{
		"String",
		"Hello, World!",
		"0d00000048656c6c6f2c20576f726c6421",
	},
	{
		"Optional nil",
		optional,
		"00",
	},
	{
		"Optional big.Int(7)",
		big.NewInt(7),
		"010107",
	},
	{
		"Empty list",
		[]uint32{},
		"00000000",
	},
	{
		"List",
		[]uint32{1, 2, 3},
		"03000000010000000200000003000000",
	},
	{
		"Fixed size list",
		[3]uint32{1, 2, 3},
		"010000000200000003000000",
	},
	{
		"Result success",
		result{
			IsSuccess: true,
			Success:   toPtrU64(314),
		},
		"013a01000000000000",
	},
	{
		"Result error",
		result{
			IsSuccess: false,
			Error:     toPtrString("Uh oh"),
		},
		"00050000005568206f68",
	},
	{
		"map[string]string",
		map[string]string{
			"test": "test",
		},
		"0100000004000000746573740400000074657374",
	},
	{
		"Tuple(u32)",
		tuple1{
			First: 1,
		},
		"01000000",
	},
	{
		"Tuple(u32,string)",
		tuple2{
			First:  1,
			Second: "Hello, World!",
		},
		"010000000d00000048656c6c6f2c20576f726c6421",
	},
	{
		"Tuple(u32,string,bool)",
		tuple3{
			First:  1,
			Second: "Hello, World!",
			Third:  true,
		},
		"010000000d00000048656c6c6f2c20576f726c642101",
	},
	{
		"Struct",
		struct {
			First  uint32
			Second string
			Third  bool
		}{
			1,
			"Hello, World!",
			true,
		},
		"010000000d00000048656c6c6f2c20576f726c642101",
	},
	{
		"Union",
		union{
			Discriminant: 0,
			FirstElem:    toPtrString("Hello, World!"),
		},
		"000d00000048656c6c6f2c20576f726c6421",
	},
}

func toPtrU64(val uint64) *uint64 {
	return &val
}

func toPtrString(val string) *string {
	return &val
}

var optional *big.Int

type result struct {
	IsSuccess bool
	Success   *uint64
	Error     *string
}

func (r result) ResultFieldName() string {
	return "IsSuccess"
}

func (r result) SuccessFieldName() string {
	return "Success"
}

func (r result) ErrorFieldName() string {
	return "Error"
}

type tuple struct {
	First  uint32
	Second string
	Third  bool
}

type tuple1 tuple

func (t tuple1) TupleFields() []string {
	return []string{"First"}
}

type tuple2 tuple

func (t tuple2) TupleFields() []string {
	return []string{"First", "Second"}
}

type tuple3 tuple

func (t tuple3) TupleFields() []string {
	return []string{"First", "Second", "Third"}
}

type union struct {
	Discriminant byte
	FirstElem    *string
	SecondElem   *uint32
}

func (u union) SwitchFieldName() string {
	return "Discriminant"
}

func (u union) ArmForSwitch(sw byte) (string, bool) {
	switch sw {
	case 0:
		return "FirstElem", true
	case 1:
		return "SecondElem", true
	}
	return "-", false
}

func TestEncoding(t *testing.T) {
	for _, c := range cases {
		result, err := Marshal(c.Value)
		if err != nil {
			t.Errorf("failed to marshal %s: %v", c.Name, err)
		}

		hexResult := hex.EncodeToString(result)

		assert.Equal(t, c.HexEncodedValue, hexResult, "failed to marshal %s", c.Name)
	}
}

func TestDecoding(t *testing.T) {

}
