package consensus

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/consensus/hrs"
	"gitlab.com/zarb-chain/zarb-go/consensus/message"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/network"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/vote"
)

type synchronizer struct {
	ctx     context.Context
	config  *config.Config
	cstate  *Consensus
	store   *store.Store
	topic   *pubsub.Topic
	sub     *pubsub.Subscription
	self    peer.ID
	syncing bool
	logger  *logger.Logger
}

func newSynchronizer(conf *config.Config, cstate *Consensus, net *network.Network, logger *logger.Logger) (*synchronizer, error) {
	syncer := &synchronizer{
		ctx:     context.Background(),
		config:  conf,
		cstate:  cstate,
		syncing: true,
		logger:  logger,
	}

	topic, err := net.JoinTopic("consensus")
	if err != nil {
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	syncer.self = net.ID()
	syncer.topic = topic
	syncer.sub = sub

	return syncer, nil
}

func (syncer *synchronizer) Start() error {
	go syncer.readLoop()

	return nil
}

func (syncer *synchronizer) Stop() error {
	syncer.ctx.Done()
	syncer.sub.Cancel()
	syncer.topic.Close()
	return nil
}

func (syncer *synchronizer) readLoop() {

	// Let's other peers know we don't have any votes
	hrs := syncer.cstate.HeightRoundStep()
	syncer.BroadcastNewStep(hrs)

	for {
		m, err := syncer.sub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}
		// only forward messages delivered by others
		if m.ReceivedFrom == syncer.self {
			continue
		}

		msg := new(message.Message)
		err = msg.UnmarshalCBOR(m.Data)
		if err != nil {
			syncer.logger.Error("Error decoding message", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
			continue
		}
		syncer.logger.Trace("Received a message", "from", m.ReceivedFrom.Pretty(), "message", msg)

		if err = msg.SanityCheck(); err != nil {
			syncer.logger.Error("Peer sent us invalid msg", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
			continue
		}

		switch msg.PayloadType() {
		case message.PayloadTypeStep:
			pld := msg.Payload.(*message.StepPayload)
			hrs := syncer.cstate.HeightRoundStep()
			if msg.Height == hrs.Height() {
				if pld.Round > hrs.Round() || // Peer is in further round
					(pld.Round == hrs.Round() && pld.Step > hrs.Step()) { // Peer is in further step
					// We are behind of the peer, ask for more votes
					votes := syncer.cstate.AllVotes()
					hashes := make([]crypto.Hash, len(votes))
					for i, v := range votes {
						hashes[i] = v.Hash()
					}
					msg := message.NewVoteSetMessage(hrs.Height(), hashes)
					syncer.publishMessage(msg)
				}
			}

		case message.PayloadTypeProposal:
			pld := msg.Payload.(*message.ProposalPayload)

			syncer.cstate.SetProposal(pld.Proposal)
		case message.PayloadTypeBlock:
			//pld := msg.Payload.(*message.BlockPayload)

		case message.PayloadTypeVote:
			pld := msg.Payload.(*message.VotePayload)

			syncer.cstate.AddVote(pld.Vote)

		case message.PayloadTypeVoteSet:
			pld := msg.Payload.(*message.VoteSetPayload)
			hrs := syncer.cstate.HeightRoundStep()
			if msg.Height == hrs.Height() {
				// Sending votes to peer
				ourVotes := syncer.cstate.AllVotes()
				peerVotes := pld.Votes

				for _, v1 := range ourVotes {
					hasVote := false
					for _, v2 := range peerVotes {
						if v1.Hash() == v2 {
							hasVote = true
							break
						}
					}

					if !hasVote {
						msg := message.NewVoteMessage(v1)
						syncer.publishMessage(msg)
					}
				}
			}

		default:
			syncer.logger.Error("Unknown message type", "msg", msg)
		}
	}
}

func (syncer *synchronizer) BroadcastNewStep(hrs hrs.HRS) {
	msg := message.NewStepMessage(hrs)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) BroadcastProposal(p *vote.Proposal) {
	msg := message.NewProposalMessage(p)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) BroadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) publishMessage(msg *message.Message) {
	bs, _ := msg.MarshalCBOR()
	if err := syncer.topic.Publish(syncer.ctx, bs); err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Trace("Publishing new message", "message", msg)
	}
}
