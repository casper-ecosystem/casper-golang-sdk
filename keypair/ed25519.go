package keypair

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"

	"golang.org/x/crypto/blake2b"
)

func Ed25519Random() (KeyPair, error) {
	var rawSeed [32]byte

	_, err := io.ReadFull(rand.Reader, rawSeed[:])
	if err != nil {
		return nil, err
	}

	return Ed25519FromSeed(rawSeed[:]), nil
}

func Ed25519FromSeed(seed []byte) KeyPair {
	return &ed25519KeyPair{
		seed: seed,
	}
}

type ed25519KeyPair struct {
	seed []byte
}

func (key *ed25519KeyPair) RawSeed() []byte {
	return key.seed
}

func (key *ed25519KeyPair) KeyTag() KeyTag {
	return KeyTagEd25519
}

func (key *ed25519KeyPair) PublicKey() PublicKey {
	pubKey, _ := key.keys()

	return PublicKey{
		Tag:        key.KeyTag(),
		PubKeyData: pubKey,
	}
}

func (key *ed25519KeyPair) Sign(mes []byte) Signature {
	_, privKey := key.keys()
	sign := ed25519.Sign(privKey, mes)
	return Signature{
		Tag:           key.KeyTag(),
		SignatureData: sign,
	}
}

func (key *ed25519KeyPair) Verify(mes []byte, sign []byte) bool {
	pubKey, _ := key.keys()
	return ed25519.Verify(pubKey, mes, sign)
}

func (key *ed25519KeyPair) AccountHash() string {
	pubKey, _ := key.keys()
	buffer := append([]byte(strKeyTagEd25519), separator)
	buffer = append(buffer, pubKey...)

	hash := blake2b.Sum256(buffer)

	return fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))
}

func (key *ed25519KeyPair) keys() (ed25519.PublicKey, ed25519.PrivateKey) {
	reader := bytes.NewReader(key.seed)
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	return pub, priv
}

func ParseKeyPair(publicKey []uint8, privateKey []uint8) {
	pub, _ := ParsePublicKey(publicKey)
	priv, _ := ParsePrivateKey(privateKey)
}

func ParseKey(bytes []uint8, from int, to int) ([]byte, error) {
	length := len(bytes)
	var key []uint8
	if  length == 32 {
		key = bytes;
	}
	if  length == 64 {
		key = bytes[from:to]
	}
	if  length > 32 && length < 64 {
		key = bytes[length %32:]
	}
	if key == nil || len(key) != 32 {
		return nil, errors.New("Unexpected key length")
	}
	return key, nil
}

func ParsePublicKey(bytes []uint8) ([]byte, error) {
	return ParseKey(bytes, 32, 64)
}

func ParsePrivateKey(bytes []uint8) ([]byte, error) {
	return ParseKey(bytes, 0, 32)
}


