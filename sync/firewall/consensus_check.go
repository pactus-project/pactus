package firewall

import (
	"bytes"
	gosync "sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
)

// maxSignatureChecksPerPeer bounds the consensus messages per second whose
// signature the firewall checks for one peer. A committee sends a few dozen
// votes per second, so an honest peer stays well below it; past it, a flooding
// peer's messages are no longer relayed, and checking them costs no more.
const maxSignatureChecksPerPeer = 100

// committeeKeys keeps the public keys of the committee of the current height,
// so a gossip message can be checked without reading the whole committee from
// the state each time. It is refreshed once per height.
type committeeKeys struct {
	lk     gosync.Mutex
	height types.Height
	keys   map[crypto.Address]*bls.PublicKey
}

// hasForgedSignature reports whether a vote or a proposal of the current height
// carries a signature that its committee signer, or the proposer of its round,
// did not make. The firewall does not relay such a message, so a flood of
// forged votes stops at the first node instead of reaching every validator.
//
// Messages of other heights, and other messages, are left to the consensus
// as before: the committee of another height may differ.
func (f *Firewall) hasForgedSignature(data []byte) bool {
	bdl, _, err := f.decodeBundle(bytes.NewReader(data), maxGossipBundleSize)
	if err != nil {
		// The bundle is checked again, and rejected, when it is opened.
		return false
	}

	height := f.state.LastBlockHeight() + 1

	switch msg := bdl.Message.(type) {
	case *message.VoteMessage:
		vte := msg.Vote
		if vte == nil || vte.Height() != height {
			return false
		}
		if vte.BasicCheck() != nil {
			return true
		}
		pub := f.committeePublicKey(height, vte.Signer())
		if pub == nil {
			// The signer is not in the committee of this height.
			return true
		}

		return vte.Verify(pub) != nil

	case *message.ProposalMessage:
		prop := msg.Proposal
		if prop == nil || prop.Height() != height {
			return false
		}
		if prop.BasicCheck() != nil {
			return true
		}
		proposer := f.state.Proposer(prop.Round())
		if proposer == nil {
			return false
		}

		return prop.Verify(proposer.PublicKey()) != nil

	default:
		return false
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
