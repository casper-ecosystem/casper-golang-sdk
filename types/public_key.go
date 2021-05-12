package types

import "io"

type KeyTag byte

const (
	KeyTagEd25519 KeyTag = iota + 1
	KeyTagSecp256k1
)

// PublicKey representing public key
type PublicKey struct {
	Tag        KeyTag
	PubKeyData []byte
}

func (key PublicKey) Marshal(w io.Writer) (int, error) {
	n, err := w.Write([]byte{byte(key.Tag)})
	if err != nil {
		return n, err
	}

	n2, err := w.Write(key.PubKeyData)
	n += n2

	return n, err
}
