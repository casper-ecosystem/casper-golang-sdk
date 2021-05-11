package ed25519

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"strings"
	"testing"
)

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
	const publKey = "48656c6c6f20476f706865722134563f"
	publKeyInPem := ExportPublicKeyInPem(publKey)

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/public.pem"
	err := ioutil.WriteFile(fileName, publKeyInPem, 0644)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	signKeyPair2, _ := ParsePublicKeyFile(fileName)
	encKeyPair := base64.StdEncoding.EncodeToString([]byte(publKey))
	encKeyPair2 := base64.StdEncoding.EncodeToString(signKeyPair2)
	assert.Equal(t, encKeyPair, encKeyPair2)
}

func TestEd25519KeyPair_ExportPrivateKeyInPem(t *testing.T) {
	const privKey = "48656c6c6f20476f706865722134563f"
	privKeyInPem := ExportPrivateKeyInPem(privKey)

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/private.pem"
	err := ioutil.WriteFile(fileName, privKeyInPem, 0644)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	signKeyPair2, _ := ParsePrivateKeyFile(fileName)
	encKeyPair := base64.StdEncoding.EncodeToString([]byte(privKey))
	encKeyPair2 := base64.StdEncoding.EncodeToString(signKeyPair2)
	assert.Equal(t, encKeyPair, encKeyPair2)
}

func TestEd25519KeyPair_Sign(t *testing.T) {
	kp,_ := Ed25519Random()
	message := []byte("hello world")
	sign := kp.Sign(message).SignatureData
	dst := make([]byte, hex.EncodedLen(len(sign)))
	hex.Encode(dst, sign)
	assert.Equal(t, 128, len(dst))
}

func TestEd25519KeyPair_Verify(t *testing.T) {
	kp,_ := Ed25519Random()
	message := []byte("hello world")
	sign := kp.Sign(message).SignatureData
	verify := kp.Verify(sign, message)
	assert.Equal(t, false, verify)
}