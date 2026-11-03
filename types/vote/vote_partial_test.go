package vote_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A vote that is decoded but not checked yet must not panic.
func TestPartialVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valKey := ts.RandValKey()

	t.Run("change-proposer vote without its data", func(t *testing.T) {
		vte := vote.NewPrepareVote(ts.RandHash(), 1, 0, valKey.Address())
		vote.SetTypeForTest(vte, vote.VoteTypeCPMainVote) // The data of the change-proposer is missing.

		require.NotPanics(t, func() {
			assert.Nil(t, vte.SignBytes())
			assert.Zero(t, vte.CPRound())
			assert.Equal(t, vote.CPValueNo, vte.CPValue())
			assert.Nil(t, vte.CPJust())
			assert.True(t, vte.IsCPVote())
		})
	})

	t.Run("vote without signature", func(t *testing.T) {
		vte := vote.NewPrepareVote(ts.RandHash(), 1, 0, valKey.Address())

		require.NotPanics(t, func() {
			err := vte.Verify(valKey.PublicKey())
			require.ErrorIs(t, err, crypto.ErrInvalidSignature)
		})
	})
}
