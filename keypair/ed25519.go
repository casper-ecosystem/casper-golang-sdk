package keypair

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	bs "encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robpike/filter"
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
	"strings"
)

const ED25519_PEM_SECRET_KEY_TAG = "PRIVATE KEY"
const ED25519_PEM_PUBLIC_KEY_TAG = "PUBLIC KEY"

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
	seed 		[]byte
	PublKey 	[]byte
	PrivateKey 	[]byte
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

func ParseKeyPair(publicKey []byte, privateKey []byte) KeyPair{
	pub, _ := ParsePublicKey(publicKey)
	priv, _ := ParsePrivateKey(privateKey)
	keyPair := ed25519KeyPair{
		PublKey: pub,
		PrivateKey: priv,
	}
	return &keyPair
}

func ParseKey(bytes []byte, from int, to int) ([]byte, error) {
	length := len(bytes)
	var key []byte
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

func ParsePublicKey(bytes []byte) ([]byte, error) {
	return ParseKey(bytes, 32, 64)
}

func ParsePrivateKey(bytes []byte) ([]byte, error) {
	return ParseKey(bytes, 0, 32)
}

func ReadBase64File(path string) ([]byte,error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("can't read file")
	}
	resContent := bytesToString(content)
	return ReadBase64WithPEM(resContent)
}

func bytesToString(data []byte) string {
	return string(data[:])
}

func ReadBase64WithPEM(content string) ([]byte, error) {
	base64 := strings.Split(content, "/\r\n/")
	var selec = filter.Choose(base64, filterFunction).([]string)
	join := strings.Join(selec, "")
	res := strings.Trim(join, "\r\n")
	return bs.StdEncoding.DecodeString(res)
}

func filterFunction(a string) bool {
	return ! strings.HasPrefix(a, "-----")
}

func ParsePublicKeyFile(path string) ([]byte, error) {
	key, _ := ReadBase64File(path)
	return ParsePublicKey(key)
}

func ParsePrivateKeyFile(path string) ([]byte, error) {
	key, _ := ReadBase64File(path)
	return ParsePrivateKey(key)
}

func AccountHex(publicKey []byte) string {
	dst := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(dst, publicKey)
	return "01"+string(dst)
}

func ParseKeyFiles(pubKeyPath, privKeyPath string) KeyPair {
	pub, _ := ParsePublicKeyFile(pubKeyPath)
	priv, _ := ParsePublicKeyFile(privKeyPath)

	keyPair := ed25519KeyPair{
		PublKey: pub,
		PrivateKey: priv,
	}
	return &keyPair
}

func (e *ed25519KeyPair) ExportPrivateKeyInPem() []byte {
	block := &pem.Block{
		Type: ED25519_PEM_SECRET_KEY_TAG,
		Bytes: e.PrivateKey,
	}
	return pem.EncodeToMemory(block)
}

func (e *ed25519KeyPair) ExportPublicKeyInPem() []byte {
	block := &pem.Block{
		Type: ED25519_PEM_PUBLIC_KEY_TAG,
		Bytes: e.PublKey,
	}
	return pem.EncodeToMemory(block)
}