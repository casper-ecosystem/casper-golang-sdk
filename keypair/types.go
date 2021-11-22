package keypair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
)

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

func (key PublicKey) MarshalJSON() ([]byte, error) {
	var w bytes.Buffer
	_, err := key.Marshal(&w)
	if err != nil {
		return nil, nil
	}
	return json.Marshal(hex.EncodeToString(w.Bytes()))
}

func (key *PublicKey) UnmarshalJSON(data []byte) error {
	var result string

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	resultHex, err := hex.DecodeString(result)
	if err != nil {
		return err
	}

	key.Tag = KeyTag(resultHex[0])
	key.PubKeyData = resultHex[1:]
	return nil
}

func (key PublicKey) ToBytes() ([]byte, error) {
	var w bytes.Buffer
	_, err := key.Marshal(&w)
	return w.Bytes(), err
}

type Signature struct {
	Tag           KeyTag
	SignatureData []byte
}

func (signature Signature) Marshal(w io.Writer) (int, error) {
	n, err := w.Write([]byte{byte(signature.Tag)})
	if err != nil {
		return n, err
	}

	n2, err := w.Write(signature.SignatureData)
	n += n2

	return n, err
}

func (signature Signature) MarshalJSON() ([]byte, error) {
	var w bytes.Buffer
	_, err := signature.Marshal(&w)
	if err != nil {
		return nil, nil
	}
	return json.Marshal(hex.EncodeToString(w.Bytes()))
}
