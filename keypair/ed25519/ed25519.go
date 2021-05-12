package ed25519

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"io"
)

const ED25519_PEM_SECRET_KEY_TAG = "PRIVATE KEY"
const ED25519_PEM_PUBLIC_KEY_TAG = "PUBLIC KEY"

// Ed25519Random creates a random Ed25519 keypair
func Ed25519Random() (keypair.KeyPair, error) {
	var rawSeed [32]byte

	_, err := io.ReadFull(rand.Reader, rawSeed[:])
	if err != nil {
		return nil, err
	}

	return Ed25519FromSeed(rawSeed[:]), nil
}

// Ed25519FromSeed creates a new keypair from the provided Ed25519 seed
func Ed25519FromSeed(seed []byte) keypair.KeyPair {
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

func (key *ed25519KeyPair) KeyTag() keypair.KeyTag {
	return keypair.KeyTagEd25519
}

// PublicKey returns Ed25519 public key
func (key *ed25519KeyPair) PublicKey() keypair.PublicKey {
	pubKey, _ := key.keys()

	return keypair.PublicKey {
		Tag:        key.KeyTag(),
		PubKeyData: pubKey,
	}
}

// Sign signs the message by using the keyPair
func (key *ed25519KeyPair) Sign(mes []byte) keypair.Signature {
	_, privKey := key.keys()
	sign := ed25519.Sign(privKey, mes)
	return keypair.Signature{
		Tag:           key.KeyTag(),
		SignatureData: sign,
	}
}

// Verify verifies the signature along with the raw message
func (key *ed25519KeyPair) Verify(mes []byte, sign []byte) bool {
	pubKey, _ := key.keys()
	return ed25519.Verify(pubKey, mes, sign)
}

// AccountHash generates the accountHash for the Ed25519 public key
func (key *ed25519KeyPair) AccountHash() string {
	pubKey, _ := key.keys()
	buffer := append([]byte(keypair.StrKeyTagEd25519), keypair.Separator)
	buffer = append(buffer, pubKey...)

	hash := blake2b.Sum256(buffer)

	return fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))
}

// keys generates a public and private key
func (key *ed25519KeyPair) keys() (ed25519.PublicKey, ed25519.PrivateKey) {
	reader := bytes.NewReader(key.seed)
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	return pub, priv
}

// ParseKeyPair constructs keyPair from a public key and private key
func ParseKeyPair(publicKey []byte, privateKey []byte) keypair.KeyPair {
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

func ParsePublicKeyFile(path string) ([]byte, error) {
	key, _ := keypair.ReadBase64File(path)
	return ParsePublicKey(key)
}

func ParsePrivateKeyFile(path string) ([]byte, error) {
	key, _ := keypair.ReadBase64File(path)
	return ParsePrivateKey(key)
}

// AccountHex generates the accountHex for the Ed25519 public key
func AccountHex(publicKey []byte) string {
	dst := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(dst, publicKey)
	return "01"+string(dst)
}

// ParseKeyFiles parses the key pair from publicKey file and privateKey file
func ParseKeyFiles(pubKeyPath, privKeyPath string) keypair.KeyPair {
	pub, _ := ParsePublicKeyFile(pubKeyPath)
	priv, _ := ParsePublicKeyFile(privKeyPath)

	keyPair := ed25519KeyPair{
		PublKey: pub,
		PrivateKey: priv,
	}
	return &keyPair
}

// ExportPrivateKeyInPem expects the private key encoded in pem
func ExportPrivateKeyInPem(privateKey string) []byte {
	block := &pem.Block{
		Type:  ED25519_PEM_SECRET_KEY_TAG,
		Bytes: []byte(privateKey),
	}
	return pem.EncodeToMemory(block)
}

// ExportPublicKeyInPem exports the public key encoded in pem
func ExportPublicKeyInPem(publicKey string) []byte {
	block := &pem.Block{
		Type:  ED25519_PEM_PUBLIC_KEY_TAG,
		Bytes: []byte(publicKey),
	}
	return pem.EncodeToMemory(block)
}