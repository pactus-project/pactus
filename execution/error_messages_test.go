package execution

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	txID := ts.RandHash()
	addr := ts.RandAccAddress()

	assert.Equal(t, "the transaction submitted before: "+txID.String(), TransactionCommittedError{ID: txID}.Error())
	assert.Equal(t, "lock time expired: 10", LockTimeExpiredError{LockTime: 10}.Error())
	assert.Equal(t, "lock time is in the future: 20", LockTimeInFutureError{LockTime: 20}.Error())
	assert.Equal(t, "the signer is banned: "+addr.String(), SignerBannedError{Address: addr}.Error())
}
