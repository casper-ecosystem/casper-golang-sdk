package keypair

import (
	"log"
	"testing"
)

func TestPublicKey_ToAccountHash(t *testing.T) {
	publicKeyHex := "0172a54c123b336fb1d386bbdff450623d1b5da904f5e2523b3e347b6d7573ae80"
	publicKey, err := FromPublicKeyHex(publicKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	accountHashHex := publicKey.AccountHashHex()

	want := "273fdfb1f9dac9b3a8ef104982d253c44a2c915416dbbe9cbb8d3d31647b4a10"

	if accountHashHex != want {
		t.Fatalf("Got account hash: \n%s\nbut want: \n%s\n", accountHashHex, want)
	}
}
