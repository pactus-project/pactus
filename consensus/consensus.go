package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/log"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

type consensus struct {
	lk sync.RWMutex

	config              *Config
	log                 *log.Log
	signer              crypto.Signer
	state               state.Facade
	height              int
	round               int
	newHeightState      consState
	proposeState        consState
	prepareState        consState
	precommitState      consState
	commitState         consState
	currentState        consState
	changeProposerState consState
	broadcastCh         chan payload.Payload
	logger              *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.Facade,
	signer crypto.Signer,
	broadcastCh chan payload.Payload) (Consensus, error) {
	cs := &consensus{
		config:      conf,
		state:       state,
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.log = log.NewLog()
	cs.logger = logger.NewLogger("_consensus", cs)

	cs.newHeightState = &newHeightState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}
	cs.changeProposerState = &changeProposerState{cs}

	cs.height = -1
	cs.round = -1

	return cs, nil
}

func (cs *consensus) Start() error {
	return nil
}

func (cs *consensus) Stop() {
}

func (cs *consensus) Fingerprint() string {
	return fmt.Sprintf("{%d/%d/%s}", cs.height, cs.round, cs.currentState.name())
}

func (cs *consensus) HeightRound() (int, int) {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height, cs.round
}

func (cs *consensus) Height() int {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height
}

func (cs *consensus) Round() int {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.round
}

func (cs *consensus) RoundProposal(round int) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.RoundProposal(round)
}
func (cs *consensus) AllVotes() []*vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	for r := 0; r <= cs.round; r++ {
		m := cs.log.MustGetRoundMessages(r)
		votes = append(votes, m.AllVotes()...)
	}
	return votes
}
func (cs *consensus) RoundVotes(round int) []*vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	rm := cs.log.RoundMessages(round)
	if rm != nil {
		return rm.AllVotes()
	}
	return nil
}

func (cs *consensus) HasVote(hash hash.Hash) bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.HasVote(hash)
}

func (cs *consensus) enterNewState(s consState) {
	cs.currentState = s
	cs.currentState.enter()
}

func (cs *consensus) MoveToNewHeight() {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if cs.state.LastBlockHeight()+1 > cs.height {
		if cs.currentState != cs.newHeightState {
			cs.enterNewState(cs.newHeightState)
		}
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height int, round int, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	cs.logger.Debug("new timer scheduled ⏱️", "duration", duration, "height", height, "round", round, "target", target)

	go func() {
		<-timer.C
		cs.handleTimeout(ti)
	}()
}

func (cs *consensus) SetProposal(p *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if p.Height() != cs.height {
		cs.logger.Trace("invalid height", "proposal", p)
		return
	}

	roundProposal := cs.log.RoundProposal(p.Round())
	if roundProposal != nil {
		cs.logger.Trace("this round has proposal", "proposal", p)
		return
	}

	proposer := cs.proposer(p.Round())
	if err := p.Verify(proposer.PublicKey()); err != nil {
		cs.logger.Error("proposal has invalid signature", "proposal", p, "err", err)
		return
	}

	if err := cs.state.ValidateBlock(p.Block()); err != nil {
		cs.logger.Warn("invalid block", "proposal", p, "err", err)
		return
	}

	cs.currentState.onSetProposal(p)
}

func (cs *consensus) doSetProposal(p *proposal.Proposal) {
	cs.logger.Info("proposal set", "proposal", p)
	cs.log.SetRoundProposal(p.Round(), p)
}

func (cs *consensus) handleTimeout(t *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("handle ticker", "ticker", t)

	// Old tickers might trigged now. Ignore them
	if cs.height != t.Height || cs.round != t.Round {
		cs.logger.Trace("stale ticker", "ticker", t)
		return
	}

	cs.currentState.onTimedout(t)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if v.Height() != cs.height {
		cs.logger.Trace("vote has invalid height", "vote", v)
		return
	}

	if cs.log.HasVote(v.Hash()) {
		cs.logger.Trace("vote exists", "vote", v)
		return
	}

	cs.currentState.onAddVote(v)
}

func (cs *consensus) doAddVote(v *vote.Vote) {
	err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding a vote", "vote", v, "err", err)
	}

	cs.logger.Debug("new vote added", "vote", v)
}

func (cs *consensus) proposer(round int) *validator.Validator {
	return cs.state.Proposer(round)
}

func (cs *consensus) signAddVote(msgType vote.Type, hash hash.Hash) {
	address := cs.signer.Address()
	if !cs.log.CanVote(address) {
		cs.logger.Trace("this node is not in committee", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.height, cs.round, hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("our vote signed and broadcasted", "vote", v)

	err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding our vote!", "err", err, "vote", v)
	} else {
		cs.broadcastVote(v)
	}
}

func (cs *consensus) queryProposal() {
	pld := payload.NewQueryProposalPayload(cs.height, cs.round)
	cs.broadcastCh <- pld
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	pld := payload.NewProposalPayload(p)
	cs.broadcastCh <- pld
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	pld := payload.NewVotePayload(v)
	cs.broadcastCh <- pld
}

func (cs *consensus) announceNewBlock(h int, b *block.Block, c *block.Certificate) {
	pld := payload.NewBlockAnnouncePayload(h, b, c)
	cs.broadcastCh <- pld
}

// TODO: Improve the performance?
func (cs *consensus) PickRandomVote() *vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	rndRound := util.RandInt(cs.round + 1)
	votes := []*vote.Vote{}
	if rndRound == cs.round {
		m := cs.log.MustGetRoundMessages(rndRound)
		votes = append(votes, m.AllVotes()...)
	} else {
		// Don't broadcast prepare and precommit votes for previous rounds
		vs := cs.log.ChangeProposerVoteSet(rndRound)
		votes = append(votes, vs.AllVotes()...)
	}
	if len(votes) == 0 {
		return nil
	}
	return votes[util.RandInt(len(votes))]
}
