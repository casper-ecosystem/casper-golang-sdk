package types

import (
	"encoding/hex"
	"io"

	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
)

type CLMap struct {
	KeyType   CLType
	ValueType CLType
	// key of the map is encoded CLValue of the KeyType type
	Raw map[string]CLValue
}

func (clmap CLMap) Marshal(w io.Writer) (int, error) {
	size := len(clmap.Raw)

	bytes := make([]byte, 0)

	sizeBytes, err := serialization.Marshal(int32(size))
	if err != nil {
		return 0, err
	}

	bytes = append(bytes, sizeBytes...)

	for k, v := range clmap.Raw {
		var keysBytes []byte
		if clmap.KeyType == CLTypeString {
			keysBytes, err = serialization.Marshal(k)
			if err != nil {
				return 0, err
			}
		} else {
			keysBytes, err = hex.DecodeString(k)
			if err != nil {
				return 0, err
			}
		}
		valueBytes, err := serialization.Marshal(v)
		if err != nil {
			return 0, err
		}

		bytes = append(bytes, keysBytes...)
		bytes = append(bytes, valueBytes...)
	}

	return w.Write(bytes)
}
