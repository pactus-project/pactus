package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

//-----------------------------------------------------------------------------

var (
	msgQueueSize = 1000
)

type Consensus struct {
	lk deadlock.RWMutex

	config        *config.Config
	hrs           hrs.HRS
	votes         *HeightVoteSet
	valset        *validator.ValidatorSet
	privValidator *validator.PrivValidator
	lastCommit    *block.Commit
	commitRound   int
	state         *state.State
	store         *store.Store
	syncer        *synchronizer
	logger        *logger.Logger
}

func NewConsensus(
	conf *config.Config,
	state *state.State,
	net *network.Network,
	store *store.Store,
	privValidator *validator.PrivValidator,
) (*Consensus, error) {
	cs := &Consensus{
		config:        conf,
		state:         state,
		store:         store,
		valset:        state.ValidatorSet(),
		privValidator: privValidator,
	}

	// See enterNewHeight.
	cs.hrs = hrs.NewHRS(state.LastBlockHeight(), 0, hrs.StepTypeNewHeight)
	cs.logger = logger.NewLogger("consensus", cs)

	syncer, err := newSynchronizer(conf, cs, net, cs.logger)
	if err != nil {
		return nil, err
	}
	cs.syncer = syncer

	return cs, nil
}

func (cs *Consensus) Start() error {
	cs.scheduleNewHeight()
	cs.stateListener()

	if err := cs.syncer.Start(); err != nil {
		return err
	}

	return nil
}

func (cs *Consensus) Stop() {
	cs.syncer.Stop()
}

func (cs *Consensus) stateListener() {

	ch := make(chan int, 10)
	cs.state.SetNewHeightListener(ch)
	for {
		select {
		case height := <-ch:
			cs.logger.Info("New height", "h", height)
		}

	}
}
func (cs *Consensus) Fingerprint() string {
	return fmt.Sprintf("{%v}",
		cs.hrs.Fingerprint())
}

func (cs *Consensus) HeightRoundStep() hrs.HRS {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.hrs
}

func (cs *Consensus) updateRoundStep(round int, step hrs.StepType) {
	cs.hrs.UpdateRoundStep(round, step)

	go cs.syncer.BroadcastNewStep(cs.hrs)
}

func (cs *Consensus) updateHeight(height int) {
	cs.hrs.UpdateHeight(height)
}

func (cs *Consensus) Proposal(round int) *vote.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.votes.RoundProposal(round)
}

func (cs *Consensus) AllVotes() []*vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	votes := cs.votes.votes
	slice := make([]*vote.Vote, len(votes))
	i := 0
	for _, v := range votes {
		slice[i] = v
		i++
	}
	return slice
}

func (cs *Consensus) Vote(h crypto.Hash) *vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.votes.votes[h]
}

func (cs *Consensus) scheduleTimeout(duration time.Duration, height int, round int, step hrs.StepType) {
	to := timeout{duration, height, round, step}
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		cs.handleTimeout(to)
	}()
	logger.Debug("Scheduled timeout", "dur", duration, "height", height, "round", round, "step", step)
}

func (cs *Consensus) AddVote(v *vote.Vote) error {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.addVote(v)
}

func (cs *Consensus) SetProposal(proposal *vote.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.setProposal(proposal)
}

func (cs *Consensus) handleTimeout(ti timeout) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Debug("Handle timeout", "timeout", ti)

	// timeouts must be for current height, round, step
	if ti.Height != cs.hrs.Height() || ti.Round < cs.hrs.Round() {
		cs.logger.Debug("Ignoring timeout", "timeout", ti)
		return
	}

	switch ti.Step {
	case hrs.StepTypeNewHeight:
		cs.enterNewHeight(ti.Height + 1)
	case hrs.StepTypeNewRound:
		cs.enterNewRound(ti.Height, ti.Round+1)
	case hrs.StepTypePrevote:
		cs.enterPrevote(ti.Height, ti.Round)
	case hrs.StepTypePrecommit:
		cs.enterPrecommit(ti.Height, ti.Round)
	default:
		panic(fmt.Sprintf("Invalid timeout step: %v", ti.Step))
	}

}

//-----------------------------------------------------------------------------

func (cs *Consensus) addVote(v *vote.Vote) error {
	// Height mismatch is ignored.
	if cs.hrs.InvalidHeight(v.Height()) {
		return errors.Errorf(errors.ErrInvalidVote, "Vote ignored, height mismatch: %v", v.Height())
	}

	added, err := cs.votes.AddVote(v)
	if err != nil {
		return err
	}
	if !added {
		// we probably had this vote before
		return nil
	}

	height := v.Height()
	round := v.Round()
	switch v.Type() {
	case vote.VoteTypePrevote:
		prevotes := cs.votes.Prevotes(round)
		cs.logger.Debug("Vote added to prevote", "vote", v, "voteset", prevotes)
		// current round
		if cs.hrs.Round() == round {
			if ok := prevotes.HasQuorum(); ok {
				blockHash := prevotes.QuorumBlock()
				if blockHash == nil {
					cs.enterPrevoteWait(height, round)
				} else if blockHash.IsUndef() {
					cs.enterPrecommit(height, round)
				} else {
					cs.enterPrecommit(height, round)
				}
			}
		}

	case vote.VoteTypePrecommit:
		precommits := cs.votes.Precommits(round)
		cs.logger.Debug("Vote added to precommit", "vote", v, "voteset", precommits)
		// current round
		if cs.hrs.Round() == round {
			if ok := precommits.HasQuorum(); ok {
				blockHash := precommits.QuorumBlock()
				if blockHash == nil {
					cs.enterPrecommitWait(height, round)
				} else if blockHash.IsUndef() {
					cs.enterNewRound(height, round+1)
				} else {
					cs.enterCommit(height, round)
				}
			}
		}

	default:
		cs.logger.Panic("Unexpected vote type %X", v.Type)
	}

	return err
}

func (cs *Consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	if cs.privValidator == nil {
		cs.logger.Error("This node is not a validator")
		return
	}

	address := cs.privValidator.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Error("This node is not in validator set", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.hrs.Height(), cs.hrs.Round(), hash, address)
	cs.privValidator.SignMsg(v)
	cs.addVote(v)

	// Broadcast our vote
	go cs.syncer.BroadcastVote(v)
}
