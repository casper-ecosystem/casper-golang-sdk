package keypair

import "io"

type KeyTag byte

const (
	KeyTagEd25519 KeyTag = iota + 1
	KeyTagSecp256k1
)

const (
	StrKeyTagEd25519        = "ed25519"
	StrKeyTagSecp256k1      = "secp256k1"
	Separator          byte = 0
)

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

type Signature struct {
	Tag           KeyTag
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
