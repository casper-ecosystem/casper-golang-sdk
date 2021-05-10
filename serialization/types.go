package serialization

import (
	"io"
	"math/big"
)

type Marshaler interface {
	Marshal(w io.Writer) (int, error)
}

type Unmarshaler interface {
	Unmarshal(r io.Reader) (int, error)
}

type Tuple interface {
	TupleFields() []string
}

type Result interface {
	ResultFieldName() string
	SuccessFieldName() string
	ErrorFieldName() string
}

type Union interface {
	ArmForSwitch(byte) (string, bool)
	SwitchFieldName() string
}

type U128 struct {
	big.Int
}

type U256 struct {
	big.Int
}

type U512 struct {
	big.Int
}
