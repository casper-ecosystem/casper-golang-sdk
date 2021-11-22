package sdk

import (
	"encoding/hex"
	"encoding/json"
	"time"
)

type Hash []byte

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h))
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}

	decodedString, err := hex.DecodeString(dataString)
	if err != nil {
		return err
	}

	*h = decodedString
	return nil
}

type Timestamp int64

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Unix(0, int64(t)*1000000).UTC().Format("2006-01-02T15:04:05.999Z"))
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02T15:04:05.999Z", dataString)
	if err != nil {
		return err
	}

	*t = Timestamp(parse.UnixNano() / 1000000)
	return nil
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d * 1000000).String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}

	duration, err := time.ParseDuration(dataString)
	if err != nil {
		return err
	}

	*d = Duration(duration / 1000000)

	return nil
}
