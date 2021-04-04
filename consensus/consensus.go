package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/consensus/pending_votes"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

const (
	newHeightName = "new-height"
	newRoundName  = "new-round"
	proposeName   = "propose"
	prepareName   = "prepare"
	precommitName = "precommit"
	commitName    = "commit"
)

type consensus struct {
	lk deadlock.RWMutex

	config         *Config
	pendingVotes   *pending_votes.PendingVotes
	signer         crypto.Signer
	state          state.StateFacade
	height         int
	round          int
	newHeightState consState
	newRoundState  consState
	proposeState   consState
	prepareState   consState
	precommitState consState
	commitState    consState
	currentState   consState
	broadcastCh    chan *message.Message
	logger         *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.StateFacade,
	signer crypto.Signer,
	broadcastCh chan *message.Message) (Consensus, error) {
	cs := &consensus{
		config:      conf,
		state:       state,
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.pendingVotes = pending_votes.NewPendingVotes()
	cs.logger = logger.NewLogger("_consensus", cs)

	cs.newHeightState = &newHeightState{cs}
	cs.newRoundState = &newRoundState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}

	cs.height = -1
	cs.round = -1
	cs.currentState = &initState{}

	return cs, nil
}

func (cs *consensus) Stop() {

}

func (cs *consensus) Fingerprint() string {
	return fmt.Sprintf("{%d/%d/%s}", cs.height, cs.round, cs.currentState.name())
}

func (cs *consensus) HRS() hrs.HRS {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return hrs.NewHRS(cs.height, cs.round, hrs.StepTypePropose)
}

func (cs *consensus) RoundProposal(round int) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.pendingVotes.RoundProposal(round)
}

func (cs *consensus) RoundVotes(round int) []*vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	rv := cs.pendingVotes.MustGetRoundVotes(round)
	return rv.AllVotes()
}

func (cs *consensus) HasVote(hash crypto.Hash) bool {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.pendingVotes.HasVote(hash)
}

func (cs *consensus) enterNewState(s consState) {
	cs.currentState = s
	cs.currentState.enter()
}

func (cs *consensus) MoveToNewHeight() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	if cs.state.LastBlockHeight()+1 > cs.height {
		cs.enterNewState(cs.newHeightState)
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height int, round int, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		cs.handleTimeout(ti)
	}()
	logger.Trace("Scheduled timeout", "duration", duration, "height", height, "round", round, "target", target)
}

func (cs *consensus) SetProposal(p *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.setProposal(p)
}

func (cs *consensus) handleTimeout(t *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("Handle ticker", "ticker", t)

	// Old tickers might trigged now. Ignore them
	if cs.height != t.Height || cs.round != t.Round {
		cs.logger.Trace("Stale ticker", "ticker", t)
		return
	}

	cs.currentState.timedout(t)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if cs.height != v.Height() {
		return
	}

	if err := cs.addVote(v); err != nil {
		cs.logger.Error("Error on adding a vote", "vote", v, "err", err)
	}
}

func (cs *consensus) addVote(v *vote.Vote) error {
	added, err := cs.pendingVotes.AddVote(v)
	if !added {
		// we probably have this vote
		return err
	}

	cs.logger.Debug("New vote added", "vote", v)
	cs.currentState.voteAdded(v)

	return err
}

func (cs *consensus) proposer(round int) *validator.Validator {
	return cs.state.Proposer(round)
}

func (cs *consensus) setProposal(p *proposal.Proposal) {
	if p.Height() != cs.height {
		cs.logger.Trace("Propose: Invalid height", "proposal", p)
		return
	}

	roundProposal := cs.pendingVotes.RoundProposal(p.Round())
	if roundProposal != nil {
		cs.logger.Trace("propose: This round has proposal", "proposal", p)
		return
	}

	proposer := cs.proposer(p.Round())
	if err := p.Verify(proposer.PublicKey()); err != nil {
		cs.logger.Error("propose: Proposal has invalid signature", "proposal", p, "err", err)
		return
	}

	if err := cs.state.ValidateBlock(p.Block()); err != nil {
		cs.logger.Warn("propose: Invalid block", "proposal", p, "err", err)
		return
	}

	cs.logger.Info("propose: Proposal set", "proposal", p)
	cs.pendingVotes.SetRoundProposal(p.Round(), p)

	if p.Round() == cs.round {
		cs.currentState.execute()
	}
}

func (cs *consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	address := cs.signer.Address()
	if !cs.pendingVotes.CanVote(address) {
		cs.logger.Trace("This node is not in committee", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.height, cs.round, hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("Our vote signed and broadcasted", "vote", v)

	err := cs.addVote(v)
	if err != nil {
		cs.logger.Error("Error on adding our vote!", "err", err, "vote", v)
		return
	}

	cs.broadcastVote(v)
}

func (cs *consensus) requestForProposal() {
	msg := message.NewOpaqueQueryProposalMessage(cs.height, cs.round)
	cs.broadcastCh <- msg
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	msg := message.NewProposalMessage(p)
	cs.broadcastCh <- msg
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	cs.broadcastCh <- msg
}

func (cs *consensus) broadcastBlock(h int, b *block.Block, c *block.Certificate) {
	msg := message.NewOpaqueBlockAnnounceMessage(h, b, c)
	cs.broadcastCh <- msg
}

func (cs *consensus) PickRandomVote() *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	round := util.RandInt(cs.round + 1)
	rv := cs.pendingVotes.MustGetRoundVotes(round)
	votes := rv.AllVotes()
	if len(votes) == 0 {
		return nil
	}
	r := util.RandInt(len(votes))
	return votes[r]
}
