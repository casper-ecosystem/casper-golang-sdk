package ed25519

import (
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
	resHash := fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))

	assert.Equal(t, signKeyPair.AccountHash(), resHash)
}

func TestEd25519KeyPair_AccountHex(t *testing.T) {
	const publKey = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	publKeyBytes, err := hex.DecodeString(publKey)
	if err != nil {
		return
	}
	assert.Equal(t, "01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061", AccountHex(publKeyBytes))
}

func TestEd25519KeyPair_ExportPublicKeyInPem(t *testing.T) {
	const publKey = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	publKeyInPem := ExportPublicKeyInPem(publKey)

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/public.pem"
	err := ioutil.WriteFile(fileName, publKeyInPem, 0644)

	if err != nil {
		assert.NoError(t, err)
		return
	}

	publKeyFromPem, err := ParsePublicKeyFile(fileName)

	if err != nil {
		assert.NoError(t, err)
		return
	}

	publKeyFromPemString := hex.EncodeToString(publKeyFromPem)
	assert.Equal(t, publKey, publKeyFromPemString)
}

func TestEd25519KeyPair_ParseKeyFiles(t *testing.T) {
	const key = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	publKeyInPem := ExportPublicKeyInPem(key)
	privKeyInPem := ExportPrivateKeyInPem(key)

	dir, _ := ioutil.TempDir("", "test")
	fileNamePubl := dir + "/public.pem"
	fileNamePriv := dir + "/private.pem"
	err := ioutil.WriteFile(fileNamePubl, publKeyInPem, 0644)

	if err != nil {
		assert.NoError(t, err)
		return
	}
	err = ioutil.WriteFile(fileNamePriv, privKeyInPem, 0644)

	if err != nil {
		assert.NoError(t, err)
		return
	}

	keypair, err := ParseKeyFiles(fileNamePubl, fileNamePriv)
	if err != nil {
		return
	}

	assert.Equal(t, key, hex.EncodeToString(keypair.PublicKey().PubKeyData))
}

func TestEd25519KeyPair_ParseKeyPair(t *testing.T) {
	const key = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	keyBytes, err := hex.DecodeString(key)

	if err != nil {
		return
	}

	keypair := ParseKeyPair(keyBytes, keyBytes)

	assert.Equal(t, key, hex.EncodeToString(keypair.PublicKey().PubKeyData))
}

func TestEd25519KeyPair_ParseKey32(t *testing.T) {
	const key = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	keyBytes, err := hex.DecodeString(key)

	assert.Equal(t, 32, len(keyBytes))

	parseKey, err := ParseKey(keyBytes, 0, 32)

	if !assert.NoError(t, err, "failed to parse") {
		return
	}

	assert.Equal(t, key, hex.EncodeToString(parseKey))
}

func TestEd25519KeyPair_ParseKey64(t *testing.T) {
	const key = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	keyBytes, err := hex.DecodeString(key)

	assert.Equal(t, 64, len(keyBytes))

	parseKey, err := ParseKey(keyBytes, 0, 32)

	if !assert.NoError(t, err, "failed to parse") {
		return
	}

	assert.Equal(t, key[:64], hex.EncodeToString(parseKey))
}

func TestEd25519KeyPair_ParseError(t *testing.T) {
	const key = "d995c93ac47e76"
	keyBytes, err := hex.DecodeString(key)

	assert.Greater(t, 32, len(keyBytes))

	_, err = ParseKey(keyBytes, 0, 32)

	assert.Error(t, err, "failed to parse")
}

func TestEd25519KeyPair_Random(t *testing.T) {
	keypair, err := Ed25519Random()

	assert.NoError(t, err, "failed to get random")

	assert.Equal(t, 32, len(keypair.PublicKey().PubKeyData))
	assert.Equal(t, 32, len(keypair.RawSeed()))
}

func TestEd25519KeyPair_ExportPrivateKeyInPem(t *testing.T) {
	const privKey = "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061"
	privKeyInPem := ExportPrivateKeyInPem(privKey)

	dir, _ := ioutil.TempDir("", "test")
	fileName := dir + "/public.pem"
	err := ioutil.WriteFile(fileName, privKeyInPem, 0644)

	if err != nil {
		assert.NoError(t, err)
		return
	}

	privKeyFromPem, err := ParsePrivateKeyFile(fileName)

	if err != nil {
		assert.NoError(t, err)
		return
	}

	privKeyFromPemString := hex.EncodeToString(privKeyFromPem)
	assert.Equal(t, privKey, privKeyFromPemString)
}

func TestEd25519KeyPair_Sign(t *testing.T) {
	kp, _ := Ed25519Random()
	message := []byte("hello world")
	sign := kp.Sign(message).SignatureData
	dst := make([]byte, hex.EncodedLen(len(sign)))
	hex.Encode(dst, sign)
	assert.Equal(t, 128, len(dst))
}

func TestEd25519KeyPair_Verify(t *testing.T) {
	kp, _ := ParseKeyFiles("../test_account_keys/account1/public_key.pem", "../test_account_keys/account1/secret_key.pem")
	message := []byte("hello world")
	sign := kp.Sign(message).SignatureData
	verify := kp.Verify(sign, message)
	assert.Equal(t, true, verify)
}

func TestEd25519KeyPair_VerifyFalse(t *testing.T) {
	kp, _ := ParseKeyFiles("../test_account_keys/account1/public_key.pem", "../test_account_keys/account2/secret_key.pem")
	message := []byte("hello world")
	sign := kp.Sign(message).SignatureData
	verify := kp.Verify(sign, message)
	assert.Equal(t, false, verify)
}