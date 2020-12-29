package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

type consensus struct {
	lk deadlock.RWMutex

	config         *Config
	hrs            hrs.HRS
	pendingVotes   *PendingVotes
	valset         validator.ValidatorSetReader
	signer         crypto.Signer
	isProposed     bool
	isPrepared     bool
	isPreCommitted bool
	isCommitted    bool
	state          state.State
	broadcastCh    chan *message.Message
	logger         *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.State,
	signer crypto.Signer,
	broadcastCh chan *message.Message) (Consensus, error) {
	cs := &consensus{
		config:      conf,
		state:       state,
		valset:      state.ValidatorSet(),
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.pendingVotes = NewPendingVotes(0, cs.valset)
	cs.hrs = hrs.NewHRS(0, 0, hrs.StepTypeUnknown)
	cs.logger = logger.NewLogger("_consensus", cs)

	return cs, nil
}

func (cs *consensus) Stop() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	cs.hrs.UpdateHeight(-1)
}

func (cs *consensus) Fingerprint() string {
	isProposed := "-"
	if cs.isProposed {
		isProposed = "X"
	}
	isPrepared := "-"
	if cs.isPrepared {
		isPrepared = "X"
	}
	isPreCommitted := "-"
	if cs.isPreCommitted {
		isPreCommitted = "X"
	}
	isCommitted := "-"
	if cs.isCommitted {
		isCommitted = "X"
	}
	status := fmt.Sprintf("%s%s%s%s", isProposed, isPrepared, isPreCommitted, isCommitted)

	return fmt.Sprintf("{%v %s}",
		cs.hrs.String(), status)
}

func (cs *consensus) HRS() hrs.HRS {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.hrs
}

func (cs *consensus) updateRound(round int) {
	cs.hrs.UpdateRound(round)
}

func (cs *consensus) updateStep(step hrs.StepType) {
	cs.hrs.UpdateStep(step)
}

func (cs *consensus) updateHeight(height int) {
	if cs.hrs.Height() != height {
		cs.pendingVotes.Reset(height)
		cs.hrs.UpdateHeight(height)
	}
}

func (cs *consensus) LastProposal() *vote.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.pendingVotes.RoundProposal(cs.hrs.Round())
}

func (cs *consensus) RoundVotes(round int) []*vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	rv := cs.pendingVotes.MustGetRoundVotes(round)
	return rv.votes
}

func (cs *consensus) RoundVotesHash(round int) []crypto.Hash {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	rv := cs.pendingVotes.MustGetRoundVotes(round)
	hashes := make([]crypto.Hash, len(rv.votes))

	for i, v := range rv.votes {
		hashes[i] = v.Hash()
	}

	return hashes
}

func (cs *consensus) HasVote(hash crypto.Hash) bool {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.pendingVotes.HasVote(hash)
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height int, round int, step hrs.StepType) {
	to := timeout{duration, height, round, step}
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		cs.handleTimeout(to)
	}()
	logger.Debug("Scheduled timeout", "duration", duration, "height", height, "round", round, "step", step)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if err := cs.addVote(v); err != nil {
		cs.logger.Error("Error on adding a vote", "vote", v, "err", err)
	}
}

func (cs *consensus) SetProposal(proposal *vote.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.setProposal(proposal)
}

func (cs *consensus) handleTimeout(ti timeout) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Debug("Handle timeout", "timeout", ti)

	// timeouts must be for current height
	if ti.Height != cs.hrs.Height() {
		cs.logger.Debug("Ignoring timeout", "timeout", ti)
		return
	}

	switch ti.Step {
	case hrs.StepTypeNewHeight:
		cs.enterNewHeight()
	case hrs.StepTypePrepare:
		cs.enterPrepare(ti.Round)
	case hrs.StepTypePrecommit:
		cs.enterPrecommit(ti.Round)
	default:
		panic(fmt.Sprintf("Invalid timeout step: %v", ti.Step))
	}

}

//-----------------------------------------------------------------------------

func (cs *consensus) addVote(v *vote.Vote) error {
	// Height mismatch is ignored.
	if cs.hrs.Height() != v.Height() {
		return nil
	}

	added, err := cs.pendingVotes.AddVote(v)
	if err != nil {
		if v.Signer().EqualsTo(cs.signer.Address()) {
			cs.logger.Error("Detecting a duplicated vote from ourself. Did you restart the node?")
		}

		return err
	}
	if !added {
		// we probably had this vote before
		return nil
	}

	round := v.Round()
	switch v.VoteType() {
	case vote.VoteTypePrepare:
		prepares := cs.pendingVotes.PrepareVoteSet(round)
		cs.logger.Debug("Vote added to prepare", "vote", v, "voteset", prepares)

		if ok := prepares.HasQuorum(); ok {
			blockHash := prepares.QuorumBlock()
			if blockHash == nil {
				cs.enterPrepareWait(round)
			} else {
				cs.enterPrecommit(round)
			}
		}

	case vote.VoteTypePrecommit:
		precommits := cs.pendingVotes.PrecommitVoteSet(round)
		cs.logger.Debug("Vote added to precommit", "vote", v, "voteset", precommits)

		if ok := precommits.HasQuorum(); ok {
			blockHash := precommits.QuorumBlock()
			if blockHash != nil {
				if blockHash.IsUndef() {
					cs.enterNewRound(round + 1)
				} else {
					cs.enterCommit(round)
				}
			}
		}

	default:
		cs.logger.Panic("Unexpected vote type %X", v.VoteType)
	}

	return err
}

func (cs *consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	address := cs.signer.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Trace("This node is not in validator set", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.hrs.Height(), cs.hrs.Round(), hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("Our vote signed and broadcasted", "vote", v)

	err := cs.addVote(v)
	if err != nil {
		cs.logger.Error("Error on adding our vote!", "err", err, "vote", v)
		return
	}

	// Broadcast our vote
	if cs.config.FuzzTesting {
		rand := util.RandInt(3)
		go func() {
			time.Sleep(time.Duration(rand) * time.Second)
			msg := message.NewVoteMessage(v)
			cs.broadcastCh <- msg
		}()
	} else {
		msg := message.NewVoteMessage(v)
		cs.broadcastCh <- msg
	}
}

func (cs *consensus) requestForProposal() {
	msg := message.NewProposalReqMessage(cs.hrs.Height(), cs.hrs.Round())
	cs.broadcastCh <- msg
}
