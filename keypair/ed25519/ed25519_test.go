package ed25519

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"strings"
	"testing"
)

type RsaIdentity struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

func TestEd25519KeyPair_AccountHash(t *testing.T) {
	signKeyPair, _ := Ed25519Random()
	publ := signKeyPair.PublicKey()
	name := strings.ToLower("ED25519")
	sep := "00"
	decoded_sep, _ := hex.DecodeString(sep)

	buffer := append([]byte(name), decoded_sep...)
	buffer = append(buffer, publ.PubKeyData...)

	hash := blake2b.Sum256(buffer)
	resHash := 	fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))

	assert.Equal(t, signKeyPair.AccountHash(), resHash)
}

func TestEd25519KeyPair_ExportPublicKeyInPem(t *testing.T) {
	signKeyPair,_ := Ed25519Random()
	publKeyInPem := signKeyPair.ExportPublicKeyInPem()

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/public.pem"
	err := ioutil.WriteFile(fileName, publKeyInPem, 0644)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	signKeyPair2, _ := ParsePublicKeyFile(fileName)
	encKeyPair := base64.StdEncoding.EncodeToString(signKeyPair.PublicKey().PubKeyData)
	encKeyPair2 := base64.StdEncoding.EncodeToString(signKeyPair2)
	assert.Equal(t, encKeyPair, encKeyPair2)
}

func TestEd25519KeyPair_ExportPrivateKeyInPem(t *testing.T) {
	signKeyPair,_ := Ed25519Random()
	privKeyInPem := signKeyPair.ExportPrivateKeyInPem()

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/private.pem"
	err := ioutil.WriteFile(fileName, privKeyInPem, 0644)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	signKeyPair2, _ := ParsePrivateKeyFile(fileName)
	encKeyPair := base64.StdEncoding.EncodeToString(signKeyPair.PublicKey().PubKeyData)
	encKeyPair2 := base64.StdEncoding.EncodeToString(signKeyPair2)
	assert.Equal(t, encKeyPair, encKeyPair2)
}

func TestEd25519KeyPair_Sign(t *testing.T) {
	message := []byte("hello world")
	rsa, _ := NewRsaIdentity()
	signKeyPair,_ := Ed25519Random()
	rsaSign, _ := rsa.Sign(message)
	hexRsaSign := hex.EncodeToString(rsaSign)
	hexSignKeyPairSign := hex.EncodeToString(signKeyPair.Sign(message).SignatureData)
	assert.Equal(t, hexRsaSign, hexSignKeyPairSign)
}

func TestEd25519KeyPair_Verify(t *testing.T) {
	message := []byte("hello world")
	rsa, _ := NewRsaIdentity()
	signKeyPair,_ := Ed25519Random()
	rsaSign,_ := rsa.Sign(message)
	RsaVerify := rsa.Verify(message, rsaSign, rsa.public)
	edSign := signKeyPair.Sign(message)
	edVerify := signKeyPair.Verify(edSign.SignatureData, message)
	assert.Equal(t, RsaVerify, edVerify)
}

func (r *RsaIdentity) Sign(msg []byte) ([]byte, error) {
	hs := r.getHashSum(msg)
	return rsa.SignPKCS1v15(rand.Reader, r.private, crypto.SHA256, hs)
}

func (r *RsaIdentity) getHashSum(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)
	return h.Sum(nil)
}

func NewRsaIdentity() (*RsaIdentity, error) {
	identity := new(RsaIdentity)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return identity, err
	}

	identity.private = priv
	identity.public = &priv.PublicKey
	return identity, nil
}

func (r *RsaIdentity) Verify(msg []byte, sig []byte, pk *rsa.PublicKey) error {
	hs := r.getHashSum(msg)
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hs, sig)
}