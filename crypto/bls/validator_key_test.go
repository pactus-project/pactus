package bls_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestValidatorKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandBLSKeyPair()
	valKey := bls.NewValidatorKey(prv)
	sig := valKey.Sign([]byte("foo"))

	assert.Equal(t, prv, valKey.PrivateKey())
	assert.Equal(t, pub, valKey.PublicKey())
	assert.Equal(t, pub.ValidatorAddress(), valKey.Address())
	assert.NoError(t, pub.Verify([]byte("foo"), sig))
}
