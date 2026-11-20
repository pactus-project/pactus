package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type proposalHandler struct {
	*synchronizer
}

func newProposalHandler(sync *synchronizer) messageHandler {
	return &proposalHandler{
		sync,
	}
}

func (handler *proposalHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.ProposalMessage)
	handler.logger.Trace("parsing Proposal message", "msg", msg)

	handler.getConsMgr().SetProposal(msg.Proposal)

	// SetProposal validates the proposal's signature internally but
	// reports nothing back to the caller, so the protocol version it
	// carries -- unrelated data the signature doesn't even cover for its
	// own purpose -- would otherwise get attributed to any address an
	// unsigned or forged proposal names, before anyone checked it was
	// really that validator who sent it.
	proposerAddr := msg.Proposal.Block().Header().ProposerAddress()
	for _, val := range handler.state.CommitteeValidators() {
		if val.Address() == proposerAddr {
			if err := msg.Proposal.Verify(val.PublicKey()); err == nil {
				handler.state.UpdateValidatorProtocolVersion(proposerAddr, msg.ProtocolVersion)
			}

			break
		}
	}
}

func (*proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
