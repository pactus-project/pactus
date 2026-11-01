package bls_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandBLSKeyPair()
	valKey := bls.NewValidatorKey(prv)
	sig := valKey.Sign([]byte("foo"))

	// Compare by raw key material, not by deep struct equality: PublicKey's
	// pointG2 is a lazy, independently-computed cache (atomic.Pointer since
	// M-L2-19), so two PublicKey values for the same key can legitimately
	// differ in it (e.g. one populated, one not, or two separate but equal
	// points) without the keys themselves being different.
	assert.Equal(t, prv.Bytes(), valKey.PrivateKey().Bytes())
	assert.Equal(t, pub.Bytes(), valKey.PublicKey().Bytes())
	assert.Equal(t, pub.ValidatorAddress(), valKey.Address())
	require.NoError(t, pub.Verify([]byte("foo"), sig))
}
