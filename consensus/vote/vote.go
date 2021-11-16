package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
)

// Vote represents a prepare, precommit, or commit vote from validators for
// consensus.
type Vote struct {
	data voteData
}

type voteData struct {
	Type      Type             `cbor:"1,keyasint"`
	Height    int              `cbor:"2,keyasint"`
	Round     int              `cbor:"3,keyasint"`
	BlockHash hash.Hash        `cbor:"4,keyasint"`
	Signer    crypto.Address   `cbor:"5,keyasint"`
	Signature crypto.Signature `cbor:"6,keyasint"`
}

type signVote struct {
	BlockHash hash.Hash `cbor:"1,keyasint"`
	Round     int       `cbor:"2,keyasint"`
	Tail      string    `cbor:"3,keyasint,omitempty"`
}

func (v *Vote) SignBytes() []byte {
	tail := ""
	if v.Type() == VoteTypePrepare {
		tail = "prepare"
	} else if v.Type() == VoteTypeChangeProposer {
		tail = "change-proposer"
	}
	// Note:
	// We omit block height, because finally block height is not matter, block hash is matter
	bz, _ := cbor.Marshal(signVote{
		Round:     v.data.Round,
		BlockHash: v.data.BlockHash,
		Tail:      tail,
	})

	return bz
}

func NewVote(voteType Type, height int, round int, blockHash hash.Hash, signer crypto.Address) *Vote {
	return &Vote{
		data: voteData{
			Type:      voteType,
			Height:    height,
			Round:     round,
			BlockHash: blockHash,
			Signer:    signer,
		},
	}
}

func (v *Vote) Type() Type                  { return v.data.Type }
func (v *Vote) Height() int                 { return v.data.Height }
func (v *Vote) Round() int                  { return v.data.Round }
func (v *Vote) BlockHash() hash.Hash        { return v.data.BlockHash }
func (v *Vote) Signer() crypto.Address      { return v.data.Signer }
func (v *Vote) Signature() crypto.Signature { return v.data.Signature }

func (v *Vote) SetSignature(sig crypto.Signature) {
	v.data.Signature = sig
}

// SetPublicKey is doing nothing and just satisfies SignableMsg interface
func (v *Vote) SetPublicKey(crypto.PublicKey) {}

func (v *Vote) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(v.data)
}

func (v *Vote) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &v.data)
}

func (v *Vote) Hash() hash.Hash {
	bz, _ := cbor.Marshal(v.data)
	return hash.HashH(bz)
}

func (v *Vote) Verify(pubKey crypto.PublicKey) error {
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidVote, "no signature")
	}
	if !pubKey.Address().EqualsTo(v.Signer()) {
		return errors.Errorf(errors.ErrInvalidVote, "invalid signer")
	}
	if !pubKey.Verify(v.SignBytes(), v.Signature()) {
		return errors.Errorf(errors.ErrInvalidProposal, "invalid signature")
	}
	return nil
}

func (v *Vote) SanityCheck() error {
	if !v.data.Type.IsValid() {
		return errors.Errorf(errors.ErrInvalidVote, "invalid vote type")
	}
	if v.data.Height <= 0 {
		return errors.Errorf(errors.ErrInvalidVote, "invalid height")
	}
	if v.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidVote, "invalid round")
	}
	if v.data.Signer.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "invalid signer")
	}
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidVote, "no signature")
	}
	if v.Signature().SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "invalid signature")
	}
	return nil
}

func (v *Vote) Fingerprint() string {
	return fmt.Sprintf("{%v/%d/%s ⌘ %v 👤 %s}",
		v.Height(),
		v.Round(),
		v.Type(),
		v.BlockHash().Fingerprint(),
		v.Signer().Fingerprint(),
	)
}

// ---------
// For tests
func GenerateTestPrecommitVote(height, round int) (*Vote, crypto.Signer) {
	s := bls.GenerateTestSigner()
	v := NewVote(
		VoteTypePrecommit,
		height,
		round,
		hash.GenerateTestHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}

func GenerateTestPrepareVote(height, round int) (*Vote, crypto.Signer) {
	s := bls.GenerateTestSigner()
	v := NewVote(
		VoteTypePrepare,
		height,
		round,
		hash.GenerateTestHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}

func GenerateTestChangeProposerVote(height, round int) (*Vote, crypto.Signer) {
	s := bls.GenerateTestSigner()
	v := NewVote(
		VoteTypeChangeProposer,
		height,
		round,
		hash.GenerateTestHash(),
		s.Address())
	s.SignMsg(v)

	return v, s
}
