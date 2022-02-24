package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingProposalMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing proposal message", func(t *testing.T) {
		consensusHeight := tState.LastBlockHeight() + 1
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
		msg := message.NewProposalMessage(prop)

		assert.NoError(t, testReceiveingNewMessage(tSync, msg, util.RandomPeerID()))
		assert.NotNil(t, tSync.cache.GetProposal(consensusHeight, 0))
		assert.NotNil(t, tConsensus.RoundProposal(0))
	})
}
