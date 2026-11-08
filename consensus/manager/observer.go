package manager

import (
	"bytes"
	"slices"
	"sync"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/ratelimit"
)

// The observer keeps, for the current height, the votes and proposals a node
// receives through gossip while none of its validators is in the committee.
// It only reads what the node already relays: it never signs, votes or writes
// to the state, so it cannot change consensus.
//
// It is built to stay cheap on a node that does no consensus work, even under
// a flood of forged messages or during a long stall:
//   - every cheap check (height, signer, slot, cache of bad messages) runs
//     before any signature verification, and a slot that is already filled is
//     never verified again;
//   - signature verifications are rate limited for the whole node;
//   - only the two highest rounds keep their votes; older rounds are folded into
//     tallies, and only a bounded number of rounds is kept;
//   - everything is dropped when the node moves to a new height.
const (
	// observerFullRounds is how many of the highest rounds keep each vote.
	observerFullRounds = 2
	// observerMaxRounds bounds the rounds kept, with votes or as tallies.
	observerMaxRounds = 16
	// observerMaxCPRound bounds the change-proposer round of a vote.
	observerMaxCPRound = 16
	// observerMaxVerifications bounds the signatures verified per second.
	// A committee produces a few dozen votes per second at most.
	observerMaxVerifications = 100
	// observerMaxBadMessages bounds the cache of messages that failed verification.
	observerMaxBadMessages = 1024
)

// ObservedVote is a vote the observer verified, without its signature and
// change-proposer justification.
type ObservedVote struct {
	Type      vote.Type
	Signer    crypto.Address
	BlockHash hash.Hash
	Round     types.Round
	CPRound   int16
	CPValue   vote.CPValue
}

// ObservedTally sums the votes of an older round that share the same type,
// change-proposer round, value and block hash.
type ObservedTally struct {
	Round     types.Round
	Type      vote.Type
	CPRound   int16
	CPValue   vote.CPValue
	BlockHash hash.Hash
	Votes     int
	Power     int64
}

// ObservedProposal is a proposal whose signature the observer verified against
// the proposer of its round. Its block was not validated.
type ObservedProposal struct {
	Round     types.Round
	BlockHash hash.Hash
	Proposer  crypto.Address
}

// Observation is what the observer saw for one height.
type Observation struct {
	Height types.Height
	// TotalPower is the power of the committee of this height.
	TotalPower int64
	// Votes are the votes of the highest rounds, ordered by round.
	Votes []ObservedVote
	// Tallies sum the votes of older rounds, ordered by round.
	Tallies []ObservedTally
	// Proposals are ordered by round.
	Proposals []ObservedProposal
	// Dropped counts the messages the observer skipped because of its bounds
	// or its rate limit, not the duplicates or the invalid ones.
	Dropped uint64
}

// voteSlot identifies the one vote a committee member can cast for a type and round.
type voteSlot struct {
	signer  crypto.Address
	typ     vote.Type
	round   types.Round
	cpRound int16
}

type tallyKey struct {
	typ       vote.Type
	cpRound   int16
	cpValue   vote.CPValue
	blockHash hash.Hash
}

type observedRound struct {
	votes   map[voteSlot]ObservedVote
	tallies map[tallyKey]*ObservedTally
	power   int64
	folded  bool
}

// addToTally counts a vote in the tally of its type, round, value and block hash.
func (r *observedRound) addToTally(vte *ObservedVote, power int64) {
	key := tallyKey{typ: vte.Type, cpRound: vte.CPRound, cpValue: vte.CPValue, blockHash: vte.BlockHash}
	tally, found := r.tallies[key]
	if !found {
		tally = &ObservedTally{
			Round: vte.Round, Type: vte.Type, CPRound: vte.CPRound,
			CPValue: vte.CPValue, BlockHash: vte.BlockHash,
		}
		r.tallies[key] = tally
	}
	tally.Votes++
	tally.Power += power
}

func (r *observedRound) fold(powers map[crypto.Address]int64) {
	if r.folded {
		return
	}
	r.tallies = make(map[tallyKey]*ObservedTally)
	for _, vte := range r.votes {
		r.addToTally(&vte, powers[vte.Signer])
	}
	// The slots stay, without their votes, so a folded round still ignores
	// duplicates before verifying them.
	for slot := range r.votes {
		r.votes[slot] = ObservedVote{}
	}
	r.folded = true
}

type observer struct {
	lk sync.RWMutex

	// enabled is set by the configuration and never changes.
	enabled bool
	// watching is true for a height none of the local validators takes part in.
	watching bool

	height     types.Height
	committee  map[crypto.Address]*validator.Validator
	powers     map[crypto.Address]int64
	totalPower int64
	proposerOf func(types.Round) *validator.Validator

	rounds      map[types.Round]*observedRound
	proposals   map[types.Round]ObservedProposal
	badMessages map[hash.Hash]struct{}
	dropped     uint64

	verifyLimit *ratelimit.RateLimit
}

func newObserver(enabled bool) *observer {
	return &observer{
		enabled:     enabled,
		verifyLimit: ratelimit.NewRateLimit(observerMaxVerifications, time.Second),
	}
}

// moveToHeight drops what was seen and starts watching a new height, unless a
// local validator is in its committee: its consensus instance keeps the votes then.
func (o *observer) moveToHeight(height types.Height, committee []*validator.Validator,
	proposerOf func(types.Round) *validator.Validator, hasActiveInstance bool,
) {
	if !o.enabled {
		return
	}

	o.lk.Lock()
	defer o.lk.Unlock()

	o.height = height
	o.watching = !hasActiveInstance
	o.proposerOf = proposerOf
	o.rounds = make(map[types.Round]*observedRound)
	o.proposals = make(map[types.Round]ObservedProposal)
	o.badMessages = make(map[hash.Hash]struct{})
	o.dropped = 0
	o.committee = make(map[crypto.Address]*validator.Validator, len(committee))
	o.powers = make(map[crypto.Address]int64, len(committee))
	o.totalPower = 0
	for _, val := range committee {
		o.committee[val.Address()] = val
		o.powers[val.Address()] = val.Power()
		o.totalPower += val.Power()
	}
}

// addVote keeps a vote of the watched height if its signer is a committee member
// and its signature is valid. Cheap checks run first; see the package notes.
func (o *observer) addVote(vte *vote.Vote) {
	if !o.enabled {
		return
	}

	o.lk.Lock()
	defer o.lk.Unlock()

	if !o.watching || vte.Height() != o.height {
		return
	}
	if vte.BasicCheck() != nil {
		return
	}
	val, ok := o.committee[vte.Signer()]
	if !ok {
		return
	}
	cpRound, cpValue := voteCPFields(vte)
	if vte.Round() < 0 || cpRound < 0 || cpRound >= observerMaxCPRound {
		o.dropped++

		return
	}

	slot := voteSlot{signer: vte.Signer(), typ: vte.Type(), round: vte.Round(), cpRound: cpRound}
	if rnd, known := o.rounds[vte.Round()]; known {
		if _, filled := rnd.votes[slot]; filled {
			// A duplicate, or a second vote of the same member for the same slot.
			return
		}
	}

	voteHash := vte.Hash()
	if _, bad := o.badMessages[voteHash]; bad {
		return
	}
	if !o.verifyLimit.AllowRequest() {
		o.dropped++

		return
	}
	if err := vte.Verify(val.PublicKey()); err != nil {
		o.rememberBad(voteHash)

		return
	}

	o.storeVote(slot, ObservedVote{
		Type:      vte.Type(),
		Signer:    vte.Signer(),
		BlockHash: vte.BlockHash(),
		Round:     vte.Round(),
		CPRound:   cpRound,
		CPValue:   cpValue,
	})
}

// voteCPFields returns the change-proposer round and value of a vote, or zeros
// for a vote that has none: Vote.CPRound and Vote.CPValue panic on such a vote.
func voteCPFields(vte *vote.Vote) (int16, vote.CPValue) {
	if !vte.IsCPVote() {
		return 0, 0
	}

	return vte.CPRound(), vte.CPValue()
}

// setProposal keeps the first proposal of a round that the proposer of the
// round signed. The block itself is not validated.
func (o *observer) setProposal(prop *proposal.Proposal) {
	if !o.enabled {
		return
	}

	o.lk.Lock()
	defer o.lk.Unlock()

	if !o.watching || prop.Height() != o.height {
		return
	}
	if _, ok := o.proposals[prop.Round()]; ok {
		return
	}
	if prop.BasicCheck() != nil {
		return
	}
	if len(o.proposals) >= observerMaxRounds {
		o.pruneProposals()
		if len(o.proposals) >= observerMaxRounds {
			o.dropped++

			return
		}
	}
	proposer := o.proposerOf(prop.Round())
	if proposer == nil || prop.Block().Header().ProposerAddress() != proposer.Address() {
		return
	}

	propHash := prop.Hash()
	if _, bad := o.badMessages[propHash]; bad {
		return
	}
	if !o.verifyLimit.AllowRequest() {
		o.dropped++

		return
	}
	if err := prop.Verify(proposer.PublicKey()); err != nil {
		o.rememberBad(propHash)

		return
	}

	o.proposals[prop.Round()] = ObservedProposal{
		Round:     prop.Round(),
		BlockHash: prop.Block().Hash(),
		Proposer:  proposer.Address(),
	}
}

// pruneProposals drops the proposals of rounds that keep no votes, so a member
// that proposes for far rounds alone cannot take the place of the rounds the
// committee votes on.
func (o *observer) pruneProposals() {
	for round := range o.proposals {
		if _, kept := o.rounds[round]; !kept {
			delete(o.proposals, round)
			o.dropped++
		}
	}
}

// evictWeakestRound drops the kept round with the least voting power, the
// lowest one on a tie. It runs only for a verified vote of a new round, so a
// forged vote cannot push a round out, and the rounds a single member signs
// alone are the first to go.
func (o *observer) evictWeakestRound() {
	weakest, found := types.Round(0), false
	for rnd, data := range o.rounds {
		if !found || data.power < o.rounds[weakest].power ||
			(data.power == o.rounds[weakest].power && rnd < weakest) {
			weakest, found = rnd, true
		}
	}
	if found {
		delete(o.rounds, weakest)
		o.dropped++
	}
}

func (o *observer) storeVote(slot voteSlot, vte ObservedVote) {
	rnd, ok := o.rounds[vte.Round]
	if !ok {
		if len(o.rounds) >= observerMaxRounds {
			o.evictWeakestRound()
		}
		rnd = &observedRound{votes: make(map[voteSlot]ObservedVote)}
		o.rounds[vte.Round] = rnd
	}
	rnd.power += o.powers[vte.Signer]

	if rnd.folded {
		rnd.votes[slot] = ObservedVote{}
		rnd.addToTally(&vte, o.powers[vte.Signer])

		return
	}
	rnd.votes[slot] = vte

	o.foldOldRounds()
}

// foldOldRounds keeps the votes of the highest rounds only.
func (o *observer) foldOldRounds() {
	rounds := make([]types.Round, 0, len(o.rounds))
	for rnd := range o.rounds {
		rounds = append(rounds, rnd)
	}
	slices.Sort(rounds)
	for i := 0; i < len(rounds)-observerFullRounds; i++ {
		o.rounds[rounds[i]].fold(o.powers)
	}
}

func (o *observer) rememberBad(msgHash hash.Hash) {
	if len(o.badMessages) >= observerMaxBadMessages {
		o.badMessages = make(map[hash.Hash]struct{})
	}
	o.badMessages[msgHash] = struct{}{}
}

// observation returns a copy of what was seen for the watched height, or nil
// when the observer is off or a local validator is in the committee.
func (o *observer) observation() *Observation {
	if !o.enabled {
		return nil
	}

	o.lk.RLock()
	defer o.lk.RUnlock()

	if !o.watching {
		return nil
	}

	obs := &Observation{
		Height:     o.height,
		TotalPower: o.totalPower,
		Dropped:    o.dropped,
	}

	rounds := make([]types.Round, 0, len(o.rounds))
	for rnd := range o.rounds {
		rounds = append(rounds, rnd)
	}
	slices.Sort(rounds)
	for _, rnd := range rounds {
		data := o.rounds[rnd]
		if data.folded {
			for _, tally := range data.tallies {
				obs.Tallies = append(obs.Tallies, *tally)
			}

			continue
		}
		for _, vte := range data.votes {
			obs.Votes = append(obs.Votes, vte)
		}
	}
	slices.SortStableFunc(obs.Votes, compareObservedVotes)
	slices.SortStableFunc(obs.Tallies, compareObservedTallies)

	for _, prop := range o.proposals {
		obs.Proposals = append(obs.Proposals, prop)
	}
	slices.SortFunc(obs.Proposals, func(a, b ObservedProposal) int { return int(a.Round) - int(b.Round) })

	return obs
}

func compareObservedVotes(left, right ObservedVote) int {
	if left.Round != right.Round {
		return int(left.Round) - int(right.Round)
	}
	if left.Type != right.Type {
		return int(left.Type) - int(right.Type)
	}
	if left.CPRound != right.CPRound {
		return int(left.CPRound) - int(right.CPRound)
	}

	return bytes.Compare(left.Signer.Bytes(), right.Signer.Bytes())
}

func compareObservedTallies(left, right ObservedTally) int {
	if left.Round != right.Round {
		return int(left.Round) - int(right.Round)
	}
	if left.Type != right.Type {
		return int(left.Type) - int(right.Type)
	}
	if left.CPRound != right.CPRound {
		return int(left.CPRound) - int(right.CPRound)
	}
	if left.CPValue != right.CPValue {
		return int(left.CPValue) - int(right.CPValue)
	}

	return bytes.Compare(left.BlockHash.Bytes(), right.BlockHash.Bytes())
}
