package tx_test

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeSignedUnbondBytes returns the encoding of a signed, non-delegated unbond
// transaction and the byte offset of the end of its payload.
func makeSignedUnbondBytes(t *testing.T, ts *testsuite.TestSuite) ([]byte, int) {
	t.Helper()

	valPub, valPrv := ts.RandBLSKeyPair()
	trx := tx.NewUnbondTx(ts.RandHeight(), valPub.ValidatorAddress())
	ts.HelperSignTransaction(valPrv, trx)
	require.NoError(t, trx.BasicCheck())

	data, err := trx.Bytes()
	require.NoError(t, err)

	// flags(1) + version(1) + lock_time(4) + fee(1) + memo(1) + payload type(1) + validator(21)
	return data, 1 + 1 + 4 + 1 + 1 + 1 + 21
}

// A transaction has a single canonical encoding: either the decoder rejects
// bytes that are not the canonical encoding, or re-encoding the decoded
// transaction gives back exactly the same bytes. Flag bits that are not
// defined must not be silently ignored, because the flags are not part of
// SignBytes: anybody can flip them without invalidating the signature, and
// the transaction ID stays the same.
func TestDecodeRejectsUnknownFlagBits(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	data, _ := makeSignedUnbondBytes(t, ts)

	mutated := append([]byte(nil), data...)
	mutated[0] |= 0x80 // undefined flag bit

	decoded, err := tx.FromBytes(mutated)
	if err != nil {
		return
	}

	reencoded, err := decoded.Bytes()
	require.NoError(t, err)
	assert.True(t, bytes.Equal(mutated, reencoded),
		"bytes with an undefined flag bit are accepted and re-encoded differently "+
			"(same tx ID, different wire bytes)")
}

// The delegation flag with a zero (treasury) owner decodes to a non-delegated
// payload. The decoded transaction has the same ID as the original, but a
// different encoding.
func TestDecodeRejectsDelegationFlagWithTreasuryOwner(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	data, end := makeSignedUnbondBytes(t, ts)

	mutated := make([]byte, 0, len(data)+1)
	mutated = append(mutated, data[:end]...)
	mutated = append(mutated, 0x00) // delegate owner = treasury address (1 byte)
	mutated = append(mutated, data[end:]...)
	mutated[0] |= 0x04 // flagWithDelegation

	decoded, err := tx.FromBytes(mutated)
	if err != nil {
		return
	}

	original, err := tx.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, original.ID(), decoded.ID())

	reencoded, err := decoded.Bytes()
	require.NoError(t, err)
	assert.True(t, bytes.Equal(mutated, reencoded),
		"a non-canonical delegated encoding of a non-delegated unbond is accepted "+
			"with the same tx ID as the canonical one")
}

// FromBytes must consume the whole input.
func TestFromBytesRejectsTrailingBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	data, _ := makeSignedUnbondBytes(t, ts)

	withTrailing := append(append([]byte(nil), data...), 0xde, 0xad, 0xbe, 0xef)

	_, err := tx.FromBytes(withTrailing)
	assert.Error(t, err, "FromBytes accepts trailing garbage after a valid transaction")
}

// Decoding into a reused Tx must not keep the state of the previous content.
func TestDecodeIntoReusedTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signed, _ := makeSignedUnbondBytes(t, ts)
	first, err := tx.FromBytes(signed)
	require.NoError(t, err)

	// Another transaction, unsigned.
	other := tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAmount(), ts.RandFee())
	otherBytes, err := other.Bytes()
	require.NoError(t, err)

	// Warm the caches of the reused object.
	reused := new(tx.Tx)
	require.NoError(t, reused.Decode(bytesReader(signed)))
	require.NoError(t, reused.BasicCheck())
	require.Equal(t, first.ID(), reused.ID())

	require.NoError(t, reused.Decode(bytesReader(otherBytes)))

	assert.Equal(t, other.ID(), reused.ID(),
		"ID() returns the memorized ID of the previously decoded transaction")
	assert.False(t, reused.IsSigned(),
		"signature of the previous content is kept after decoding an unsigned transaction")
	assert.Error(t, reused.BasicCheck(),
		"BasicCheck result of the previous content is reused (unsigned tx passes)")
}

func bytesReader(bs []byte) *bytes.Reader {
	return bytes.NewReader(bs)
}
