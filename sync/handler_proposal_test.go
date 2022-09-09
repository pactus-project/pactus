package sync

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/stretchr/testify/assert"
)

func TestParsingProposalMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing proposal message", func(t *testing.T) {
		consensusHeight := tState.LastBlockHeight() + 1
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
		msg := message.NewProposalMessage(prop)

		assert.NoError(t, testReceiveingNewMessage(tSync, msg, network.TestRandomPeerID()))
		assert.NotNil(t, tConsensus.RoundProposal(0))
	})
}
