package secp256k1

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/blake2b"
)

type secp256k1KeyPair struct {
	seed 		[]byte
	PublKey 	[]byte
	PrivateKey 	[]byte
}

func Secp256k1Random() keypair.KeyPair {
	priv := secp256k1.GenPrivKey()
	pub := priv.PubKey().Bytes()
	return &secp256k1KeyPair{PublKey: pub, PrivateKey: priv}
}

func (key *secp256k1KeyPair) RawSeed() []byte {
	return key.seed
}

func (key *secp256k1KeyPair) KeyTag() keypair.KeyTag {
	return keypair.KeyTagSecp256k1
}

// PublicKey returns Secp256k1 public key
func (key *secp256k1KeyPair) PublicKey() keypair.PublicKey {
	pubKey, _ := key.keys()

	return keypair.PublicKey{
		Tag: key.KeyTag(),
		PubKeyData: pubKey,
	}
}

// Sign signs the message by using the keyPair
func (key *secp256k1KeyPair) Sign(mes []byte) keypair.Signature {
	_, privKey := key.keys()
	sign, _ := privKey.Sign(mes)
	return keypair.Signature{
		Tag: key.KeyTag(),
		SignatureData: sign,
	}
}

// Verify verifies the signature along with the raw message
func (key *secp256k1KeyPair) Verify(mes []byte, sign []byte) bool {
	pubKey, _ := key.keys()
	return pubKey.VerifySignature(mes, sign)
}

// AccountHash generates the accountHash for the Secp256K1 public key
func (key *secp256k1KeyPair) AccountHash() string {
	pubKey, _ := key.keys()
	var buffer = append([]byte(keypair.StrKeyTagSecp256k1), keypair.Separator)
	buffer = append(buffer, pubKey...)

	hash := blake2b.Sum256(buffer)

	return fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))
}

func (key *secp256k1KeyPair) keys() (secp256k1.PubKey, secp256k1.PrivKey){
	priv := secp256k1.GenPrivKey()
	pub := priv.PubKey().Bytes()
	return pub, priv
}

// AccountHex generates the accountHex for the Secp256K1 public key
func AccountHex(publicKey []byte) string {
	dst := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(dst, publicKey)
	return "02"+string(dst)
}

// ExportPublicKeyInPem exports the public key encoded in pem
func ExportPublicKeyInPem(publicKey []byte) []byte {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncodedPub
}

// ExportPrivateKeyInPem expects the private key encoded in pem
func ExportPrivateKeyInPem(privateKey []byte) []byte {
	x509Encoded, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

	return pemEncoded
}

func ParsePublicKey(pemEncodedPub string) []byte {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.([]byte)

	return publicKey
}

func ParsePrivateKey(pemEncoded string) []byte {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	res,_ := x509.MarshalECPrivateKey(privateKey)
	return res
}

// ParseKeyPair constructs keyPair from public key and private key
func ParseKeyPair(publicKey, privateKey string) keypair.KeyPair {
	pub := ParsePublicKey(publicKey)
	priv := ParsePrivateKey(privateKey)
	keyPair := secp256k1KeyPair{PublKey: pub, PrivateKey: priv}

	return &keyPair
}

func ParsePublicKeyFile(path string) []byte {
	key, _ := keypair.ReadBase64File(path)
	return ParsePublicKey(string(key))
}

func ParsePrivateKeyFile(path string) []byte {
	key, _ := keypair.ReadBase64File(path)
	return ParsePrivateKey(string(key))
}

func ParseKeyFiles(pubKeyPath, privKeyPath string) keypair.KeyPair {
	pub := ParsePublicKeyFile(pubKeyPath)
	priv := ParsePrivateKeyFile(privKeyPath)

	keyPair := secp256k1KeyPair{
		PublKey: pub,
		PrivateKey: priv,
	}
	return &keyPair
}