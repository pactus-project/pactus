package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ProposalMessage.ProtocolVersion is not covered by the proposal signature, and
// the handler applies it to the proposer's validator record without any
// verification (the proposal check is deferred to consensus, which returns
// nothing). The committee uses these versions for the PIP-51 upgrade vote
// (committee.SupportProtocolVersion -> state.proposeBlockVersion).
func TestProposalProtocolVersionFromInvalidProposal(t *testing.T) {
	td := setup(t, nil)

	var updatedAddr crypto.Address
	var updatedVer protocol.Version
	updated := 0
	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).
		DoAndReturn(func(addr crypto.Address, ver protocol.Version) {
			updated++
			updatedAddr = addr
			updatedVer = ver
		}).AnyTimes()
	td.consV1Mgr.EXPECT().SetProposal(gomock.Any()).AnyTimes()

	height := td.state.LastBlockHeight() + 1
	prop := td.GenerateTestProposal(height, 0)
	prop.SetSignature(td.RandBLSSignature()) // not signed by the proposer

	msg := message.NewProposalMessage(prop)
	msg.ProtocolVersion = protocol.Version(200)

	td.receivingNewMessage(td.sync, msg, td.RandPeerID())

	assert.Zero(t, updated,
		"protocol version %v recorded for %v from a proposal with an invalid signature",
		updatedVer, updatedAddr)
}

// A proposal that names a member of the committee is recorded only when this member signed it.
func TestProposalProtocolVersionFromCommitteeMember(t *testing.T) {
	td := setup(t, nil)
	keys := td.setupCommittee()

	updated := 0
	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).
		DoAndReturn(func(crypto.Address, protocol.Version) { updated++ }).AnyTimes()
	td.consV1Mgr.EXPECT().SetProposal(gomock.Any()).AnyTimes()

	height := td.state.LastBlockHeight() + 1
	blk, _ := td.GenerateTestBlock(height, testsuite.BlockWithProposer(keys[0].Address()))

	forged := td.GenerateTestProposal(height, 0, testsuite.ProposalWithBlock(blk), testsuite.ProposalWithKey(keys[1]))
	msg := message.NewProposalMessage(forged)
	msg.ProtocolVersion = protocol.Version(200)
	td.receivingNewMessage(td.sync, msg, td.RandPeerID())
	assert.Zero(t, updated, "a proposal signed by another member is not recorded")

	signed := td.GenerateTestProposal(height, 0, testsuite.ProposalWithBlock(blk), testsuite.ProposalWithKey(keys[0]))
	msg = message.NewProposalMessage(signed)
	msg.ProtocolVersion = protocol.Version(200)
	td.receivingNewMessage(td.sync, msg, td.RandPeerID())
	assert.Equal(t, 1, updated, "a proposal signed by its proposer is recorded")
}
