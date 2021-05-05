package keypair

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/blake2b"
)

type secp256k1KeyPair struct {
	seed []byte
}

func (key *secp256k1KeyPair) RawSeed() []byte {
	return key.seed
}

func (key *secp256k1KeyPair) KeyTag() KeyTag {
	return KeyTagSecp256k1
}

func (key *secp256k1KeyPair) PublicKey() PublicKey {
	pubKey, _ := key.keys()

	return PublicKey{
		Tag: key.KeyTag(),
		PubKeyData: pubKey,
	}
}

func (key *secp256k1KeyPair) Sign(mes []byte) Signature {
	_, privKey := key.keys()
	sign, _ := privKey.Sign(mes)
	return Signature{
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
	buffer := append([]byte(strKeyTagSecp256k1), separator)
	buffer = append(buffer, pubKey...)

	hash := blake2b.Sum256(buffer)

	return fmt.Sprintf("account-hash-%s", hex.EncodeToString(hash[:]))
}

func (key *secp256k1KeyPair) keys() (secp256k1.PubKey, secp256k1.PrivKey){
	priv := secp256k1.GenPrivKey()
	pub := priv.PubKey().Bytes()
	return pub, priv
}