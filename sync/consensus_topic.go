package sync

import (
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

type ConsensusTopic struct {
	config    *Config
	publishFn PublishMessageFn
	consensus consensus.Consensus
	logger    *logger.Logger
}

func NewConsensusTopic(
	conf *Config,
	consensus consensus.Consensus,
	logger *logger.Logger,
	publishFn PublishMessageFn) *ConsensusTopic {
	return &ConsensusTopic{
		config:    conf,
		consensus: consensus,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (ct *ConsensusTopic) BroadcastProposal(p *vote.Proposal) {
	msg := message.NewProposalMessage(p)
	ct.publishFn(msg)
}

func (ct *ConsensusTopic) BroadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	ct.publishFn(msg)
}

func (ct *ConsensusTopic) BroadcastVoteSet() {
	hrs := ct.consensus.HRS()
	hashes := ct.consensus.RoundVotesHash(hrs.Round())

	msg := message.NewVoteSetMessage(hrs.Height(), hrs.Round(), hashes)
	ct.publishFn(msg)
}

func (ct *ConsensusTopic) ProcessVotePayload(pld *payload.VotePayload) {
	ct.logger.Trace("Process vote payload", "pld", pld)

	ct.consensus.AddVote(pld.Vote)
}

func (ct *ConsensusTopic) ProcessVoteSetPayload(pld *payload.VoteSetPayload) {
	ct.logger.Trace("Process vote-set payload", "pld", pld)

	hrs := ct.consensus.HRS()
	if pld.Height == hrs.Height() {

		// Check peers vote and send the votes he doesn't have
		ourVotes := ct.consensus.RoundVotes(pld.Round)
		peerVotes := pld.Hashes

		for _, v1 := range ourVotes {
			hasVote := false
			for _, v2 := range peerVotes {
				if v1.Hash().EqualsTo(v2) {
					hasVote = true
					break
				}
			}

			if !hasVote {
				ct.BroadcastVote(v1)
			}
		}
	}
}
func (ct *ConsensusTopic) ProcessProposalRequestPayload(pld *payload.ProposalRequestPayload) {
	ct.logger.Trace("Process proposal request payload", "pld", pld)

	hrs := ct.consensus.HRS()
	if pld.Height == hrs.Height() {
		p := ct.consensus.LastProposal()
		if p != nil {
			if p.Round() >= pld.Round {
				ct.BroadcastProposal(p)
			}
		}
	}
}

func (ct *ConsensusTopic) ProcessProposalPayload(pld *payload.ProposalPayload) {
	ct.logger.Trace("Process proposal payload", "pld", pld)

	ct.consensus.SetProposal(pld.Proposal)
}
