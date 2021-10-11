package keypair

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/blake2b"
	"io"
	"strconv"
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

func FromPublicKeyHex(publicKeyHex string) (PublicKey, error) {
	tag, err := strconv.Atoi(publicKeyHex[:2])
	if err != nil {
		return PublicKey{}, nil
	}
	publicKeyBytes, err := hex.DecodeString(publicKeyHex[2:])
	if err != nil {
		return PublicKey{}, nil
	}
	return PublicKey{
		Tag:        KeyTag(tag),
		PubKeyData: publicKeyBytes,
	}, nil
}

func (key *PublicKey) Hex() string {
	preimage := make([]byte, 0, len(key.PubKeyData)+len([]byte{byte(key.Tag)}))
	preimage = append(preimage, byte(key.Tag))
	preimage = append(preimage, key.PubKeyData...)
	return hex.EncodeToString(preimage)
}

func (key *PublicKey) AccountHash() []byte {
	var algorithmName []byte
	switch key.Tag {
	case KeyTagEd25519:
		algorithmName = []byte(StrKeyTagEd25519)
	case KeyTagSecp256k1:
		algorithmName = []byte(StrKeyTagSecp256k1)
	}
	preimage := make([]byte, 0, len(key.PubKeyData)+len(algorithmName)+1)
	preimage = append(preimage, algorithmName...)
	preimage = append(preimage, Separator)
	preimage = append(preimage, key.PubKeyData...)
	var accountHashBytes [32]byte
	accountHashBytes = blake2b.Sum256(preimage)
	return accountHashBytes[:]
}

func (key *PublicKey) AccountHashHex() string {
	accountHash := key.AccountHash()
	return hex.EncodeToString(accountHash)
}

func (key PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(key.Hex())
}

type Signature struct {
	Tag           KeyTag
	SignatureData []byte
}

func (key Signature) Marshal(w io.Writer) (int, error) {
	n, err := w.Write([]byte{byte(key.Tag)})
	if err != nil {
		return n, err
	}

	n2, err := w.Write(key.SignatureData)
	n += n2

	return n, err
}

func (key *Signature) Hex() string {
	preimage := make([]byte, 0, len(key.SignatureData)+len([]byte{byte(key.Tag)}))
	preimage = append(preimage, byte(key.Tag))
	preimage = append(preimage, key.SignatureData...)
	return hex.EncodeToString(preimage)
}

func (key Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(key.Hex())
}

type Hash [32]byte

func (h *Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Hex())
}

// HashFromSlice converts a hash in slice format to type Hash
func HashFromSlice(hashSlice []byte) Hash {
	hash := Hash{}
	copy(hash[:], hashSlice)
	return hash
}

// HashFromHex converts a hash in hex format to type Hash
func HashFromHex(hexHash string) (Hash, error) {
	hashSlice, err := hex.DecodeString(hexHash)
	if err != nil {
		return Hash{}, err
	}
	return HashFromSlice(hashSlice), nil
}

// CalculateHash calculates a Hash from given byte slice
func CalculateHash(data []byte) Hash {
	hash := blake2b.Sum256(data)
	return hash
}
