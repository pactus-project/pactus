package message

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HelloMessage.BasicCheck aggregates every key in PublicKeys
// (bls.PublicKeyAggregate, which calls PointG2() -> Unmarshal + a subgroup
// check for each one) before any cheap bound is checked. An honest node has
// at most 32 validator keys (cmd/cmd.go caps wallets at 32), yet the CBOR
// decoder accepts arrays up to 65_536 elements (sync/bundle/bundle.go
// MaxArrayElements). This test expresses the desired behavior — rejecting a
// Hello with far more keys than any honest node would ever send — and fails
// on the current code, which accepts it as long as the (cheaply forgeable,
// see TestHelloRogueKeyWithoutProofOfPossession and the repeated-key
// case below) aggregate signature verifies.
func TestHelloRejectsExcessivePublicKeyCount(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	const honestMaxValidatorKeys = 32
	const forgedKeyCount = honestMaxValidatorKeys * 10

	valKey := ts.RandValKey()
	msg := NewHelloMessage(ts.RandPeerID(), "Oscar", service.New(service.FullNode),
		ts.RandHeight(), ts.RandHash(), ts.RandHash())

	// A single real key, repeated: computing a valid aggregate signature for
	// n copies of the same signer costs the ATTACKER one real signature plus
	// cheap point additions, yet costs the RECEIVER one expensive subgroup
	// check per listed key in PublicKeyAggregate.
	repeated := make([]*bls.ValidatorKey, forgedKeyCount)
	for i := range repeated {
		repeated[i] = valKey
	}
	require.NoError(t, msg.Sign(repeated))

	require.Len(t, msg.PublicKeys, forgedKeyCount)

	err := msg.BasicCheck()
	assert.Error(t, err,
		"a Hello listing far more keys than any honest node would send must be rejected before the expensive aggregation")
}
