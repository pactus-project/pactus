package validator_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Clearing the delegation clears its share and its expiry too: they are not stored without
// a delegation, so the validator in memory is the one that is read back.
func TestClearDelegation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	plain, err := val.Bytes()
	require.NoError(t, err)

	val.SetDelegation(ts.RandAccAddress(), 1e9, 1000)
	val.SetDelegation(crypto.TreasuryAddress, 1e9, 1000)

	assert.Zero(t, val.DelegateShare())
	assert.Zero(t, val.DelegateExpiry())
	cleared, err := val.Bytes()
	require.NoError(t, err)
	assert.Equal(t, plain, cleared)
}

// A delegation section without a delegation owner is not written by the encoder:
// it is a second encoding of a validator without delegation.
func TestEmptyDelegationSectionIsRejected(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	plain, err := val.Bytes()
	require.NoError(t, err)

	val.SetDelegation(ts.RandAccAddress(), 1e9, 1000)
	delegated, err := val.Bytes()
	require.NoError(t, err)
	_, err = validator.FromBytes(delegated)
	require.NoError(t, err)

	// The treasury address is encoded in one byte, and the share and the expiry follow the owner.
	shareAndExpiry := delegated[len(plain)+crypto.AddressSize:]
	forged := append(append([]byte{}, plain...), crypto.TreasuryAddress.Bytes()[:1]...)
	forged = append(forged, shareAndExpiry...)

	_, err = validator.FromBytes(forged)
	require.ErrorIs(t, err, validator.ErrInvalidData)
}
