package sync

import (
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

type ConsensusSync struct {
	config    *Config
	consensus consensus.Consensus
	logger    *logger.Logger
	publishFn PublishMessageFn
}

func NewConsensusSync(
	conf *Config,
	consensus consensus.Consensus,
	logger *logger.Logger,
	publishFn PublishMessageFn) *ConsensusSync {
	return &ConsensusSync{
		config:    conf,
		consensus: consensus,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (cs *ConsensusSync) BroadcastQueryProposal() {
	hrs := cs.consensus.HRS()
	msg := message.NewQueryProposalMessage(hrs.Height(), hrs.Round())
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastProposal(p *vote.Proposal) {
	msg := message.NewProposalMessage(p)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastVoteSet() {
	hrs := cs.consensus.HRS()
	hashes := cs.consensus.RoundVotesHash(hrs.Round())

	msg := message.NewVoteSetMessage(hrs.Height(), hrs.Round(), hashes)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) ProcessVotePayload(pld *payload.VotePayload) {
	cs.logger.Trace("Process vote payload", "pld", pld)

	cs.consensus.AddVote(pld.Vote)
}

func (cs *ConsensusSync) ProcessVoteSetPayload(pld *payload.VoteSetPayload) {
	cs.logger.Trace("Process vote-set payload", "pld", pld)

	hrs := cs.consensus.HRS()
	if pld.Height == hrs.Height() {

		// Check peers vote and send the votes he doesn't have
		ourVotes := cs.consensus.RoundVotes(pld.Round)
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
				cs.BroadcastVote(v1)
			}
		}
	}
}
func (cs *ConsensusSync) ProcessQueryProposalPayload(pld *payload.QueryProposalPayload) {
	cs.logger.Trace("Process proposal request payload", "pld", pld)

	hrs := cs.consensus.HRS()
	if pld.Height == hrs.Height() {
		p := cs.consensus.LastProposal()
		if p != nil {
			if p.Round() >= pld.Round {
				cs.BroadcastProposal(p)
			}
		}
	}
}

func (cs *ConsensusSync) ProcessProposalPayload(pld *payload.ProposalPayload) {
	cs.logger.Trace("Process proposal payload", "pld", pld)

	cs.consensus.SetProposal(pld.Proposal)
}
