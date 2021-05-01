package serialization

import (
	"io"
)

type Marshaler interface {
	Marshal(w io.Writer) (int, error)
}

type Unmarshaler interface {
	Unmarshal(r io.Reader) error
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
