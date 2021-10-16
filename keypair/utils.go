package keypair

import (
	"encoding/asn1"
	"encoding/pem"
	"os"
)

type privateKeyDer struct {
	Integer int
	Object  struct {
		Identifier asn1.ObjectIdentifier
	} `asn1:"omitempty"`
	Key []byte
}

//ReadBase64File reads the Base64 content of a file, get rid of PEM frames
func ReadBase64File(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	if len(block.Bytes) == 32 {
		return block.Bytes, nil
	}
	var privateKey privateKeyDer
	_, err = asn1.Unmarshal(block.Bytes, &privateKey)
	if err != nil {
		return nil, err
	}
	return privateKey.Key, nil
}
