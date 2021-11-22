package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingURefString(t *testing.T) {
	testString := "uref-6ad5075addcdef0308bf9100a88292fd16e49edeb724dea2cc9f6f3730352d97-007"

	uref, err := URefFromFormattedString(testString)
	if !assert.NoError(t, err, "failed decode"){
		return
	}


	addressBytes, err := hex.DecodeString("6ad5075addcdef0308bf9100a88292fd16e49edeb724dea2cc9f6f3730352d97")

	if !assert.NoError(t, err, "failed decode") || !assert.Equal(t, len(addressBytes), 32, "failed decoded len"){
		return
	}

	var addressBytesToCompare [32]byte

	for i := 0; i < 32; i++{
		addressBytesToCompare[i] = addressBytes[i]
	}


	assert.Equal(t, uref.Address,addressBytesToCompare, "invalid address")
	assert.Equal(t, uref.AccessRight, AccessRightReadAddWrite, "invalid access right")
}

func TestToURefString(t *testing.T) {
	addressBytes, err := hex.DecodeString("6ad5075addcdef0308bf9100a88292fd16e49edeb724dea2cc9f6f3730352d97")

	if !assert.NoError(t, err, "failed decode") || !assert.Equal(t, len(addressBytes), 32, "failed decoded len"){
		return
	}

	var addressBytes32 [32]byte

	for i := 0; i < 32; i++{
		addressBytes32[i] = addressBytes[i]
	}

	uRef := URef{
		AccessRight: AccessRightReadAddWrite,
		Address:     addressBytes32,
	}

	testString := "uref-6ad5075addcdef0308bf9100a88292fd16e49edeb724dea2cc9f6f3730352d97-007"

	assert.Equal(t, uRef.ToFormattedString(),testString, "invalid formatted string")
}
