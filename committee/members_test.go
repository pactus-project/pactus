package committee_test

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// The committee finds its members by their address. Over a long sequence of random updates,
// with validators that join, join again and leave, a validator is found if and only if
// it is in the list of the committee, and the proposer is found as the proposer.
func TestMembersFollowTheList(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	const committeeSize = 7

	pool := make([]*validator.Validator, 0, 30)
	for i := range 30 {
		pool = append(pool, ts.GenerateTestValidator(testsuite.ValidatorWithNumber(int32(i))))
	}

	cmt, err := committee.NewCommittee(pool[:committeeSize], committeeSize, pool[0].Address())
	require.NoError(t, err)

	for step := range 500 {
		height := types.Height(step + 1)

		joined := make([]*validator.Validator, 0, 3)
		for range ts.RandIntMax(4) {
			val := pool[ts.RandIntMax(len(pool))]
			val.UpdateLastSortitionHeight(height)
			joined = append(joined, val)
		}
		cmt.Update(types.Round(ts.RandIntMax(3)), joined)

		inList := make(map[crypto.Address]bool, committeeSize)
		for _, val := range cmt.Validators() {
			inList[val.Address()] = true
		}
		require.Len(t, inList, cmt.Size(), "step %d", step)
		require.LessOrEqual(t, cmt.Size(), committeeSize, "step %d", step)

		for _, val := range pool {
			require.Equal(t, inList[val.Address()], cmt.Contains(val.Address()),
				"step %d: validator %d", step, val.Number())
		}

		for round := types.Round(0); round < 3; round++ {
			proposer := cmt.Proposer(round)
			require.True(t, cmt.IsProposer(proposer.Address(), round), "step %d: round %d", step, round)
			require.False(t, cmt.IsProposer(ts.RandValAddress(), round), "step %d: round %d", step, round)
		}
	}
}

func BenchmarkContains(b *testing.B) {
	ts := testsuite.NewTestSuite(&testing.T{})
	const committeeSize = 51

	vals := make([]*validator.Validator, 0, committeeSize)
	for i := range committeeSize {
		vals = append(vals, ts.GenerateTestValidator(testsuite.ValidatorWithNumber(int32(i))))
	}
	cmt, _ := committee.NewCommittee(vals, committeeSize, vals[0].Address())
	last := vals[committeeSize-1].Address()
	unknown := ts.RandValAddress()

	b.ResetTimer()
	for range b.N {
		cmt.Contains(last)
		cmt.Contains(unknown)
	}
}
