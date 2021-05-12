package keypair

type KeyPair interface {
	RawSeed() []byte
	KeyTag() KeyTag
	PublicKey() PublicKey
	AccountHash() string
	Sign(mes []byte) Signature
	Verify(sign []byte, mes []byte) bool
}


