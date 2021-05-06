package secp256k1

import (
	"crypto/ecdsa"
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

func (key *secp256k1KeyPair) RawSeed() []byte {
	return key.seed
}

func (key *secp256k1KeyPair) KeyTag() keypair.KeyTag {
	return keypair.KeyTagSecp256k1
}

func (key *secp256k1KeyPair) PublicKey() keypair.PublicKey {
	pubKey, _ := key.keys()

	return keypair.PublicKey{
		Tag: key.KeyTag(),
		PubKeyData: pubKey,
	}
}

func (key *secp256k1KeyPair) Sign(mes []byte) keypair.Signature {
	_, privKey := key.keys()
	sign, _ := privKey.Sign(mes)
	return keypair.Signature{
		Tag: key.KeyTag(),
		SignatureData: sign,
	}
}

func (key *secp256k1KeyPair) Verify(mes []byte, sign []byte) bool {
	pubKey, _ := key.keys()
	return pubKey.VerifySignature(mes, sign)
}

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

func AccountHex(publicKey []byte) string {
	dst := make([]byte, hex.EncodedLen(len(publicKey)))
	hex.Encode(dst, publicKey)
	return "02"+string(dst)
}

func (key *secp256k1KeyPair) ExportPublicKeyInPem() ([]byte, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(key.PublKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(block), nil
}

func (key *secp256k1KeyPair) ExportPrivateKeyInPem() ([]byte, error) {
	derKey, err := x509.MarshalECPrivateKey(key.PrivateKey)
	if err != nil {
		return nil, err
	}

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derKey,
	}

	return pem.EncodeToMemory(keyBlock), nil
}

func ParsePublicKey(bytes []byte) ([]byte, error) {

}

func ParsePrivateKey(bytes []byte) ([]byte, error) {

}

func ParseKeyPair(pub, priv []byte) keypair.KeyPair {
	public := ParsePublicKey(pub)
	private := ParsePrivateKey(priv)

	keyPair := secp256k1KeyPair{PublKey: public, PrivateKey: private}
	return &keyPair
}

func ParsePublicKeyFile(path string) ([]byte, error) {
	key, _ := keypair.ReadBase64File(path)
	return ParsePublicKey(key)
}

func ParsePrivateKeyFile(path string) ([]byte, error) {
	key, _ := keypair.ReadBase64File(path)
	return ParsePrivateKey(key)
}

func ParseKeyFiles(pubKeyPath, privKeyPath string) keypair.KeyPair {
	pub, _ := ParsePublicKeyFile(pubKeyPath)
	priv, _ := ParsePublicKeyFile(privKeyPath)

	keyPair := secp256k1KeyPair{
		PublKey: pub,
		PrivateKey: priv,
	}
	return &keyPair
}