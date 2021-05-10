package types

import "io"

type FixedByteArray []byte

func (c FixedByteArray) Marshal(w io.Writer) (int, error) {
	return w.Write(c)
}
