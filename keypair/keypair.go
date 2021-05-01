package keypair

type KeyPair interface {
	RawSeed() []byte
	KeyTag() KeyTag
	PublicKey() PublicKey
	Sign(mes []byte) Signature
	Verify(sign []byte, mes []byte) bool
}
