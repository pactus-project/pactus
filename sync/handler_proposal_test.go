package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingProposalMessages(t *testing.T) {
	setup(t)

	t.Run("Alice receives a proposal. she sends it to consensus", func(t *testing.T) {
		consensusHeight := tAliceState.LastBlockHeight() + 1
		p1, _ := proposal.GenerateTestProposal(consensusHeight, 0)
		pld := payload.NewProposalPayload(p1)

		assert.NoError(t, simulatingReceiveingNewMessage(t, tAliceSync, pld, util.RandomPeerID()))
		assert.NotNil(t, tAliceSync.cache.GetProposal(consensusHeight, 0))
		assert.NotNil(t, tAliceConsensus.RoundProposal(0))
	})
}
