package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	addr := ts.RandAccAddress()

	assert.Equal(t, "stake amount can't be less than 1 PAC", SmallStakeError{Minimum: 1e9}.Error())
	assert.Equal(t, "validator's stake amount can't be more than 1,000 PAC", MaximumStakeError{Maximum: 1000e9}.Error())
	assert.Equal(t, "unknown payload type: transfer", InvalidPayloadTypeError{PayloadType: payload.TypeTransfer}.Error())
	assert.Equal(t, "no account found for address: "+addr.String(), AccountNotFoundError{Address: addr}.Error())
	assert.Equal(t, "no validator found for address: "+addr.String(), ValidatorNotFoundError{Address: addr}.Error())
}
