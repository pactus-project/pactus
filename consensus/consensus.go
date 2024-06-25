package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/consensus/log"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

type broadcaster func(crypto.Address, message.Message)

type consensus struct {
	lk sync.RWMutex

	config          *Config
	logger          *logger.SubLogger
	log             *log.Log
	validators      []*validator.Validator
	cpWeakValidity  *hash.Hash // The change proposer's weak validity that is a prepared block hash
	cpDecided       int
	height          uint32
	round           int16
	cpRound         int16
	valKey          *bls.ValidatorKey
	rewardAddr      crypto.Address
	bcState         state.Facade // Blockchain state
	changeProposer  *changeProposer
	newHeightState  consState
	proposeState    consState
	prepareState    consState
	precommitState  consState
	commitState     consState
	cpPreVoteState  consState
	cpMainVoteState consState
	cpDecideState   consState
	currentState    consState
	broadcaster     broadcaster
	mediator        mediator
	active          bool
}

func NewConsensus(
	conf *Config,
	bcState state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcastCh chan message.Message,
	mediator mediator,
) Consensus {
	broadcaster := func(_ crypto.Address, msg message.Message) {
		broadcastCh <- msg
	}

	return makeConsensus(conf, bcState,
		valKey, rewardAddr, broadcaster, mediator)
}

func makeConsensus(
	conf *Config,
	bcState state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcaster broadcaster,
	mediator mediator,
) *consensus {
	cs := &consensus{
		config:      conf,
		bcState:     bcState,
		broadcaster: broadcaster,
		valKey:      valKey,
	}

	// Update height later, See enterNewHeight.
	cs.log = log.NewLog()
	cs.logger = logger.NewSubLogger("_consensus", cs)
	cs.rewardAddr = rewardAddr

	cs.changeProposer = &changeProposer{cs}
	cs.newHeightState = &newHeightState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}
	cs.cpPreVoteState = &cpPreVoteState{cs.changeProposer}
	cs.cpMainVoteState = &cpMainVoteState{cs.changeProposer}
	cs.cpDecideState = &cpDecideState{cs.changeProposer}
	cs.currentState = cs.newHeightState
	cs.mediator = mediator

	cs.height = 0
	cs.round = 0
	cs.active = false
	cs.mediator = mediator

	mediator.Register(cs)

	logger.Info("consensus instance created",
		"validator address", valKey.Address().String(),
		"reward address", rewardAddr.String())

	return cs
}

func (cs *consensus) Start() {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.moveToNewHeight()
}

func (cs *consensus) String() string {
	return fmt.Sprintf("{%s %d/%d/%s/%d}",
		cs.valKey.Address().ShortString(),
		cs.height, cs.round, cs.currentState.name(), cs.cpRound)
}

func (cs *consensus) ConsensusKey() *bls.PublicKey {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.valKey.PublicKey()
}

func (cs *consensus) HeightRound() (uint32, int16) {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height, cs.round
}

func (cs *consensus) Proposal() *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.RoundProposal(cs.round)
}

func (cs *consensus) HasVote(h hash.Hash) bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.HasVote(h)
}

// AllVotes returns all valid votes inside the consensus log up to and including
// the current consensus round.
// Valid votes from subsequent rounds are not included.
func (cs *consensus) AllVotes() []*vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	for r := int16(0); r <= cs.round; r++ {
		m := cs.log.RoundMessages(r)
		votes = append(votes, m.AllVotes()...)
	}

	return votes
}

func (cs *consensus) enterNewState(s consState) {
	cs.currentState = s
	cs.currentState.enter()
}

func (cs *consensus) MoveToNewHeight() {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.moveToNewHeight()
}

func (cs *consensus) moveToNewHeight() {
	stateHeight := cs.bcState.LastBlockHeight()
	if cs.height != stateHeight+1 {
		cs.enterNewState(cs.newHeightState)
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height uint32, round int16, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	cs.logger.Trace("new timer scheduled ⏱️", "duration", duration, "height", height, "round", round, "target", target)

	go func() {
		<-timer.C
		cs.handleTimeout(ti)
	}()
}

func (cs *consensus) handleTimeout(t *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("handle ticker", "ticker", t)

	// Old tickers might be triggered now. Ignore them.
	if cs.height != t.Height || cs.round != t.Round {
		cs.logger.Trace("stale ticker", "ticker", t)

		return
	}

	cs.logger.Debug("timer expired", "ticker", t)
	cs.currentState.onTimeout(t)
}

func (cs *consensus) SetProposal(p *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		cs.logger.Trace("we are not in the committee")

		return
	}

	if p.Height() != cs.height {
		cs.logger.Trace("invalid height", "proposal", p)

		return
	}

	if p.Round() < cs.round {
		cs.logger.Trace("proposal for expired round", "proposal", p)

		return
	}

	if err := p.BasicCheck(); err != nil {
		cs.logger.Warn("invalid proposal", "proposal", p, "error", err)

		return
	}

	roundProposal := cs.log.RoundProposal(p.Round())
	if roundProposal != nil {
		cs.logger.Trace("this round has proposal", "proposal", p)

		return
	}

	if p.Height() == cs.bcState.LastBlockHeight() {
		// A slow node might receive a proposal after committing the proposed block.
		// In this case, we accept the proposal and allow nodes to continue.
		// By doing so, we enable the validator to broadcast its votes and
		// prevent it from being marked as absent in the block certificate.
		cs.logger.Trace("block is committed for this height", "proposal", p)
		if p.Block().Hash() != cs.bcState.LastBlockHash() {
			cs.logger.Warn("proposal is not for the committed block", "proposal", p)

			return
		}
	} else {
		proposer := cs.proposer(p.Round())
		if err := p.Verify(proposer.PublicKey()); err != nil {
			cs.logger.Warn("proposal is invalid", "proposal", p, "error", err)

			return
		}

		if err := cs.bcState.ValidateBlock(p.Block(), p.Round()); err != nil {
			cs.logger.Warn("invalid block", "proposal", p, "error", err)

			return
		}
	}

	cs.logger.Info("proposal set", "proposal", p)
	cs.log.SetRoundProposal(p.Round(), p)

	cs.currentState.onSetProposal(p)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		cs.logger.Trace("we are not in the committee")

		return
	}

	if v.Height() != cs.height {
		cs.logger.Trace("vote has invalid height", "vote", v)

		return
	}

	if v.Round() < cs.round {
		cs.logger.Trace("vote for expired round", "vote", v)

		return
	}

	if v.Type() == vote.VoteTypeCPPreVote ||
		v.Type() == vote.VoteTypeCPMainVote ||
		v.Type() == vote.VoteTypeCPDecided {
		err := cs.changeProposer.checkJust(v)
		if err != nil {
			cs.logger.Error("error on adding a cp vote", "vote", v, "error", err)

			return
		}
	}

	added, err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding a vote", "vote", v, "error", err)
	}
	if added {
		cs.logger.Info("new vote added", "vote", v)

		cs.currentState.onAddVote(v)

		if v.Type() == vote.VoteTypeCPDecided {
			if v.Round() > cs.round {
				cs.changeProposer.cpDecide(v.Round(), v.CPValue())
			}
		}
	}
}

func (cs *consensus) proposer(round int16) *validator.Validator {
	return cs.bcState.Proposer(round)
}

func (cs *consensus) IsProposer() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.isProposer()
}

func (cs *consensus) isProposer() bool {
	return cs.proposer(cs.round).Address() == cs.valKey.Address()
}

func (cs *consensus) signAddCPPreVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPPreVote(h, cs.height,
		cs.round, cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddCPMainVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPMainVote(h, cs.height, cs.round,
		cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddCPDecidedVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPDecidedVote(h, cs.height, cs.round,
		cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrepareVote(h hash.Hash) {
	v := vote.NewPrepareVote(h, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrecommitVote(h hash.Hash) {
	v := vote.NewPrecommitVote(h, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddVote(v *vote.Vote) {
	sig := cs.valKey.Sign(v.SignBytes())
	v.SetSignature(sig)
	cs.logger.Info("our vote signed and broadcasted", "vote", v)

	_, err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding our vote", "error", err, "vote", v)
	}
	cs.broadcastVote(v)
}

func (cs *consensus) queryProposal() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryProposalMessage(cs.height, cs.round, cs.valKey.Address()))
}

// queryVotes is an anti-entropy mechanism to retrieve missed votes
// when a validator falls behind the network.
// However, invoking this method might result in unnecessary bandwidth usage.
func (cs *consensus) queryVotes() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryVotesMessage(cs.height, cs.round, cs.valKey.Address()))
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	go cs.mediator.OnPublishProposal(cs, p)
	cs.broadcaster(cs.valKey.Address(),
		message.NewProposalMessage(p))
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	go cs.mediator.OnPublishVote(cs, v)
	cs.broadcaster(cs.valKey.Address(),
		message.NewVoteMessage(v))
}

func (cs *consensus) announceNewBlock(blk *block.Block, cert *certificate.BlockCertificate) {
	go cs.mediator.OnBlockAnnounce(cs)
	cs.broadcaster(cs.valKey.Address(),
		message.NewBlockAnnounceMessage(blk, cert))
}

func (cs *consensus) makeBlockCertificate(votes map[crypto.Address]*vote.Vote,
) *certificate.BlockCertificate {
	cert := certificate.NewBlockCertificate(cs.height, cs.round, false)
	cert.SetSignature(cs.signersInfo(votes))

	return cert
}

// signersInfo processes a map of votes from validators and provides these information:
// - A list of all validators' numbers eligible to vote in this step.
// - A list of absentee validators' numbers who did not vote in this step.
// - An aggregated signature generated from the signatures of participating validators.
func (cs *consensus) signersInfo(votes map[crypto.Address]*vote.Vote) ([]int32, []int32, *bls.Signature) {
	vals := cs.validators
	committers := make([]int32, len(vals))
	absentees := make([]int32, 0)
	sigs := make([]*bls.Signature, 0)

	for i, val := range vals {
		vte := votes[val.Address()]
		if vte != nil {
			sigs = append(sigs, vte.Signature())
		} else {
			absentees = append(absentees, val.Number())
		}

		committers[i] = val.Number()
	}

	aggSig := bls.SignatureAggregate(sigs...)

	return committers, absentees, aggSig
}

func (cs *consensus) makeVoteCertificate(votes map[crypto.Address]*vote.Vote,
) *certificate.VoteCertificate {
	cert := certificate.NewVoteCertificate(cs.height, cs.round)
	cert.SetSignature(cs.signersInfo(votes))

	return cert
}

// IsActive checks if the consensus is in an active state and participating in the consensus algorithm.
func (cs *consensus) IsActive() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.active
}

// TODO: Improve the performance?
func (cs *consensus) PickRandomVote(round int16) *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	switch {
	case round < cs.round:
		// Past round: Only broadcast cp:decided votes
		vs := cs.log.CPDecidedVoteSet(round)
		votes = append(votes, vs.AllVotes()...)

	case round == cs.round:
		// Current round
		m := cs.log.RoundMessages(round)
		votes = append(votes, m.AllVotes()...)

	case round > cs.round:
		// Future round
	}

	if len(votes) == 0 {
		return nil
	}

	return votes[util.RandInt32(int32(len(votes)))]
}

func (cs *consensus) startChangingProposer() {
	// If it is not decided yet.
	// TODO: can we remove this condition in new consensus model?
	if cs.cpDecided == -1 {
		cs.logger.Info("changing proposer started",
			"cpRound", cs.cpRound, "proposer", cs.proposer(cs.round).Address())
		cs.enterNewState(cs.cpPreVoteState)
	}
}
