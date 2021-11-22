package types

import (
	"io"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
)

// Signature representing signature
type Signature struct {
	Tag           keypair.KeyTag
	SignatureData []byte
}

func (key Signature) Marshal(w io.Writer) (int, error) {
	n, err := w.Write([]byte{byte(key.Tag)})
	if err != nil {
		return n, err
	}

	n2, err := w.Write(key.SignatureData)
	n += n2

	return n, err
}
