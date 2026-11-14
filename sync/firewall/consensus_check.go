package firewall

import (
	"bytes"
	gosync "sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
)

// maxSignatureChecksPerPeer bounds the forged consensus messages per second that
// the firewall checks for one peer. An honest peer relays no forged message,
// because it checks the messages before relaying them. Past this bound, the
// messages of a flooding peer are no longer relayed, and checking them costs no more.
//
// Only the forged messages count. The peer is the one that relays the message,
// not its author: counting the valid messages, or the messages that can't be checked,
// would let anyone use up the budget of an honest relay, whose authentic votes
// would then no longer be relayed.
const maxSignatureChecksPerPeer = 100

// committeeKeys keeps the public keys of the committee of the current height,
// so a gossip message can be checked without reading the whole committee from
// the state each time. It is refreshed once per height.
type committeeKeys struct {
	lk     gosync.Mutex
	height types.Height
	keys   map[crypto.Address]*bls.PublicKey
}

// signatureCheck is the result of checking the signature of a consensus message.
type signatureCheck int

const (
	// signatureNotChecked is for a message that the firewall can't check:
	// it is not a vote or a proposal, or it is for another height.
	signatureNotChecked signatureCheck = iota
	signatureValid
	signatureForged
)

// checkConsensusSignature checks whether a vote or a proposal of the current height
// carries the signature of its committee signer, or of the proposer of its round.
// The firewall does not relay a forged message, so a flood of forged votes stops
// at the first node instead of reaching every validator.
//
// Messages of other heights, and other messages, are left to the consensus
// as before: the committee of another height may differ.
//
// This runs in the gossip validator of every node, for every message, so it only
// does what the signature check needs. The block of a proposal is not checked here:
// checking it verifies the signature of each of its transactions, before knowing
// that the proposal comes from the proposer. The consensus checks the block.
func (f *Firewall) checkConsensusSignature(data []byte) signatureCheck {
	bdl, _, err := f.decodeBundle(bytes.NewReader(data), maxGossipBundleSize)
	if err != nil {
		// The bundle is checked again, and rejected, when it is opened.
		return signatureNotChecked
	}

	height := f.state.LastBlockHeight() + 1

	switch msg := bdl.Message.(type) {
	case *message.VoteMessage:
		vte := msg.Vote
		if vte == nil || vte.Height() != height {
			return signatureNotChecked
		}
		if vte.BasicCheck() != nil {
			return signatureForged
		}
		pub := f.committeePublicKey(height, vte.Signer())
		if pub == nil {
			// The signer is not in the committee of this height.
			return signatureForged
		}
		if vte.Verify(pub) != nil {
			return signatureForged
		}

		return signatureValid

	case *message.ProposalMessage:
		prop := msg.Proposal
		if prop == nil || prop.Height() != height {
			return signatureNotChecked
		}
		if prop.Block() == nil || prop.Signature() == nil || prop.Round() < 0 {
			return signatureForged
		}
		proposer := f.state.Proposer(prop.Round())
		if proposer == nil {
			return signatureNotChecked
		}
		if prop.Verify(proposer.PublicKey()) != nil {
			return signatureForged
		}

		return signatureValid

	default:
		return signatureNotChecked
	}
}

// committeePublicKey returns the public key of a member of the committee of
// the given height, or nil if the address is not a member.
func (f *Firewall) committeePublicKey(height types.Height, addr crypto.Address) *bls.PublicKey {
	f.committee.lk.Lock()
	defer f.committee.lk.Unlock()

	if f.committee.keys == nil || f.committee.height != height {
		keys := make(map[crypto.Address]*bls.PublicKey)
		for _, val := range f.state.CommitteeValidators() {
			keys[val.Address()] = val.PublicKey()
		}
		f.committee.keys = keys
		f.committee.height = height
	}

	return f.committee.keys[addr]
}
