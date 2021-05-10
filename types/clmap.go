package types

import (
	"io"

	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
)

type CLMap struct {
	raw map[string]CLValue
}

func (m *CLMap) Load(key CLValue) (CLValue, bool) {
	b := string(serialization.MustMarshal(key))
	value, ok := m.raw[b]
	return value, ok
}

func (m *CLMap) Store(key, value CLValue) {
	b := string(serialization.MustMarshal(key))
	m.raw[b] = value
}

func (m *CLMap) Delete(key CLValue) {
	b := string(serialization.MustMarshal(key))
	delete(m.raw, b)
}

func (m *CLMap) Range(f func(key, value CLValue) bool) {
	for key, value := range m.raw {
		var keyCLValue CLValue
		serialization.MustUnmarshal([]byte(key), &keyCLValue)
		if !f(keyCLValue, value) {
			break
		}
	}
}

func (m *CLMap) Marshal(w io.Writer) (int, error) {
	enc := serialization.NewEncoder(w)

	n := 0
	for key, value := range m.raw {
		n2, err := enc.EncodeFixedByteArray([]byte(key))
		n += n2
		if err != nil {
			return n, err
		}
		n2, err = enc.Encode(value)
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (m *CLMap) Unmarshal(r io.Reader) (int, error) {
	// TODO
	return 0, nil
}
