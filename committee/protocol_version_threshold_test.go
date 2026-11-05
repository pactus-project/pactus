package committee_test

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// "more than 75%" is a strict threshold: a committee where exactly 75% of the
// power supports the version (18 of 24 here) must never report the support.
// The percentages are floats summed in map iteration order, so the sum of
// 12.5 + 58.33.. + 4.16.. can end above or below 75.0 from one call to another.
func TestSupportProtocolVersionExactlyThreeQuarters(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	newVal := func(stake int64, ver protocol.Version) *validator.Validator {
		val := ts.GenerateTestValidator(testsuite.ValidatorWithStake(amount.Amount(stake)))
		val.UpdateProtocolVersion(ver)

		return val
	}

	vals := []*validator.Validator{
		newVal(6, protocol.Version(1)),  // does not support version 2
		newVal(1, protocol.Version(4)),  // supports
		newVal(3, protocol.Version(3)),  // supports
		newVal(14, protocol.Version(2)), // supports
	}

	cmt, err := committee.NewCommittee(vals, 4, vals[0].Address())
	require.NoError(t, err)
	require.Equal(t, int64(24), cmt.Power())

	for i := 0; i < 300; i++ {
		require.False(t, cmt.SupportProtocolVersion(protocol.Version(2)),
			"exactly 75%% of the power is not more than 75%% (call %d)", i)
	}
}
