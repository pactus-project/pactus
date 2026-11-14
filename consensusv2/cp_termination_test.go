package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A `yes` decision in two change-proposer rounds moves the validator to the next round once.
func TestStrongTerminationMovesOneRound(t *testing.T) {
	td := setup(t)
	td.enterNewHeight(td.consP)

	cons := td.consP
	cons.lk.Lock()
	defer cons.lk.Unlock()

	height, round := cons.height, cons.round
	for _, cpRound := range []int16{0, 1} {
		just := &vote.JustDecided{QCert: td.GenerateTestCertificate(height)}
		vte := vote.NewCPDecidedVote(td.RandHash(), height, round, cpRound, vote.CPValueYes, just,
			td.valKeys[tIndexX].Address())
		td.HelperSignVote(td.valKeys[tIndexX], vte)
		_, err := cons.log.AddVote(vte)
		require.NoError(t, err)
	}

	cons.cpRound = 1
	cons.changeProposer.cpStrongTermination()

	assert.Equal(t, round+1, cons.round)
}
