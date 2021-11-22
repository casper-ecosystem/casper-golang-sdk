package ed25519

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
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
	newKeyPair := ed25519KeyPair{
		seed: seed,
	}

	publKey, privKey := newKeyPair.keys()

	newKeyPair.publicKey = publKey
	newKeyPair.privateKey = privKey

	return &newKeyPair
}

// in private key stored only part of the private key
// ed25519 private key for signature is 64 bytes, first 32 bytes is a private key like here and other 32 bytes are public key
type ed25519KeyPair struct {
	seed       []byte
	publicKey  []byte
	privateKey []byte
}

func (key *ed25519KeyPair) RawSeed() []byte {
	return key.seed
}

func (key *ed25519KeyPair) KeyTag() keypair.KeyTag {
	return keypair.KeyTagEd25519
}

// PublicKey returns Ed25519 public key
func (key *ed25519KeyPair) PublicKey() keypair.PublicKey {
	return keypair.PublicKey{
		Tag:        key.KeyTag(),
		PubKeyData: key.publicKey,
	}
}

// Sign signs the message by using the keyPair
func (key *ed25519KeyPair) Sign(mes []byte) keypair.Signature {
	sign := ed25519.Sign(append(key.privateKey, key.publicKey...), mes)

	return keypair.Signature{
		Tag:           key.KeyTag(),
		SignatureData: sign,
	}
}

// Verify verifies the signature along with the raw message
func (key *ed25519KeyPair) Verify(sign []byte, mes []byte) bool {
	return ed25519.Verify(key.publicKey, mes, sign)
}

// AccountHash generates the accountHash for the Ed25519 public key
func (key *ed25519KeyPair) AccountHash() string {
	buffer := append([]byte(keypair.StrKeyTagEd25519), keypair.Separator)
	buffer = append(buffer, key.publicKey...)

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

	return pub, append(make([]byte, 0), priv[:32]...)
}

// ParseKeyPair constructs keyPair from a public key and private key
func ParseKeyPair(publicKey []byte, privateKey []byte) keypair.KeyPair {
	pub, _ := ParsePublicKey(publicKey)
	priv, _ := ParsePrivateKey(privateKey)
	keyPair := ed25519KeyPair{
		publicKey:  pub,
		privateKey: priv,
	}
	return &keyPair
}

func ParseKey(bytes []byte, from int, to int) ([]byte, error) {
	length := len(bytes)
	var key []byte
	if length == 32 {
		key = bytes
	}
	if length == 64 {
		key = bytes[from:to]
	}
	if length > 32 && length < 64 {
		key = bytes[length%32:]
	}
	if key == nil || len(key) != 32 {
		return nil, errors.New("Unexpected key length")
	}
	return key, nil
}

func ParsePublicKey(bytes []byte) ([]byte, error) {
	return ParseKey(bytes, 12, 32)
}

func ParsePrivateKey(bytes []byte) ([]byte, error) {
	return ParseKey(bytes, 16, 32)
}

func ParsePublicKeyFile(path string) ([]byte, error) {
	key, err := keypair.ReadBase64File(path)

	if err != nil {
		return nil, err
	}

	return ParsePublicKey(key)
}

func ParsePrivateKeyFile(path string) ([]byte, error) {
	key, err := keypair.ReadBase64File(path)

	if err != nil {
		return nil, err
	}

	return ParsePrivateKey(key)
}

// AccountHex generates the accountHex for the Ed25519 public key
func AccountHex(publicKey []byte) string {
	dst := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(dst, publicKey)
	return "01" + string(dst)
}

// AccountHex generates the accountHex for the Ed25519 public key
func AccountHash(pubKey []byte) string {
	buffer := append([]byte(keypair.StrKeyTagEd25519), keypair.Separator)
	buffer = append(buffer, pubKey...)

	hash := blake2b.Sum256(buffer)

	return hex.EncodeToString(hash[:])
}

// ParseKeyFiles parses the key pair from publicKey file and privateKey file
func ParseKeyFiles(pubKeyPath, privKeyPath string) (keypair.KeyPair, error) {
	pub, err := ParsePublicKeyFile(pubKeyPath)

	if err != nil {
		return nil, err
	}

	priv, err := ParsePrivateKeyFile(privKeyPath)

	if err != nil {
		return nil, err
	}

	keyPair := ed25519KeyPair{
		publicKey:  pub,
		privateKey: priv,
	}
	return &keyPair, nil
}

// ExportPrivateKeyInPem expects the private key encoded in pem
func ExportPrivateKeyInPem(privateKey string) []byte {
	derPrefix := make([]byte, 16)
	derPrefix[0] = 48
	derPrefix[1] = 46
	derPrefix[2] = 2
	derPrefix[3] = 1
	derPrefix[4] = 0
	derPrefix[5] = 48
	derPrefix[6] = 5
	derPrefix[7] = 6
	derPrefix[8] = 3
	derPrefix[9] = 43
	derPrefix[10] = 101
	derPrefix[11] = 112
	derPrefix[12] = 4
	derPrefix[13] = 34
	derPrefix[14] = 4
	derPrefix[15] = 32

	privateKeyBytes, err := hex.DecodeString(privateKey)

	if err != nil {
		return nil
	}

	block := &pem.Block{
		Type:  ED25519_PEM_SECRET_KEY_TAG,
		Bytes: append(derPrefix, privateKeyBytes...),
	}

	return pem.EncodeToMemory(block)
}

// ExportPublicKeyInPem exports the public key encoded in pem
func ExportPublicKeyInPem(publicKey string) []byte {
	derPrefix := make([]byte, 12)
	derPrefix[0] = 48
	derPrefix[1] = 42
	derPrefix[2] = 48
	derPrefix[3] = 5
	derPrefix[4] = 6
	derPrefix[5] = 3
	derPrefix[6] = 43
	derPrefix[7] = 101
	derPrefix[8] = 112
	derPrefix[9] = 3
	derPrefix[10] = 33
	derPrefix[11] = 0

	publicKeyBytes, err := hex.DecodeString(publicKey)

	if err != nil {
		return nil
	}

	block := &pem.Block{
		Type:  ED25519_PEM_PUBLIC_KEY_TAG,
		Bytes: append(derPrefix, publicKeyBytes...),
	}

	return pem.EncodeToMemory(block)
}
