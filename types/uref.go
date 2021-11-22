package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type AccessRight byte

const (
	AccessRightRead AccessRight = 1 << iota
	AccessRightWrite
	AccessRightAdd
	AccessRightReadWrite                = AccessRightRead | AccessRightWrite
	AccessRightReadAdd                  = AccessRightRead | AccessRightAdd
	AccessRightAddWrite                 = AccessRightAdd | AccessRightWrite
	AccessRightReadAddWrite             = AccessRightRead | AccessRightAdd | AccessRightWrite
	AccessRightNone         AccessRight = 0

	URefPrefix = "uref-"
)

type URef struct {
	AccessRight AccessRight
	Address     [32]byte
}

func (uRef URef) Marshal(w io.Writer) (int, error) {
	n, err := w.Write(append(uRef.Address[:], byte(uRef.AccessRight)))
	if err != nil {
		return n, err
	}

	return n, err
}

func (uRef *URef) Unmarshal(data []byte) (int, error) {
	var temp [32]byte
	if len(data) < 33 {
		return 0, errors.New("inproper source, not enough bytes")
	}
	for i := 0; i < 32; i++ {
		temp[i] = data[i]
	}

	uRef.Address = temp
	uRef.AccessRight = AccessRight(data[32])

	return 33, nil
}

func (uRef URef) ToFormattedString() string {
	return fmt.Sprintf("%s%s-%03d", URefPrefix, hex.EncodeToString(uRef.Address[:]), uRef.AccessRight)
}

func URefFromFormattedString(str string) (*URef, error) {
	if str[:len(URefPrefix)] != URefPrefix {
		return nil, errors.New("invalid prefix (not 'uref-')")
	}

	splitRes := strings.Split(str[len(URefPrefix):], "-")

	if len(splitRes) < 2 {
		return nil, errors.New("no access rights as suffix")
	}

	decodedAddress, err := hex.DecodeString(splitRes[0])
	if err != nil {
		return nil, err
	}

	if len(decodedAddress) != 32 {
		return nil, errors.New("invalid address length")
	}

	var result URef

	for i := 0; i < 32; i++ {
		result.Address[i] = decodedAddress[i]
	}

	parsedAccessRights, err := strconv.ParseInt(splitRes[1], 8, 8)
	if err != nil {
		return nil, err
	}

	result.AccessRight = AccessRight(byte(parsedAccessRights))

	return &result, nil
}
