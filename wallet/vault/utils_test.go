package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMnemonic(t *testing.T) {
	_, err := GenerateMnemonic(127)
	assert.Error(t, err, "low entropy")

	_, err = GenerateMnemonic(128)
	assert.NoError(t, err)

	_, err = GenerateMnemonic(257)
	assert.Error(t, err, "high entropy")

	_, err = GenerateMnemonic(256)
	assert.NoError(t, err)
}

func TestValidateMnemonic(t *testing.T) {
	tests := []struct {
		mnenomic string
		errStr   string
	}{
		{
			"",
			"invalid mnenomic",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access",
			"invalid mnenomic",
		},
		{
			"bandon ability able about above absent absorb abstract absurd abuse access ability",
			"word `bandon` not found in reverse map",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access accident",
			"checksum incorrect",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access ability",
			"",
		},
	}
	for no, tt := range tests {
		err := CheckMnemonic(tt.mnenomic)
		if err != nil {
			assert.Equal(t, tt.errStr, err.Error(), "test %v failed", no)
		}
	}
}
