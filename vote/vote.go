package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

// Vote represents a prevote, precommit, or commit vote from validators for
// consensus.
type Vote struct {
	data voteData
}

type voteData struct {
	VoteType  VoteType          `cbor:"1,keyasint"`
	Height    int               `cbor:"2,keyasint"`
	Round     int               `cbor:"3,keyasint"`
	BlockHash crypto.Hash       `cbor:"4,keyasint"`
	Signer    crypto.Address    `cbor:"5,keyasint"`
	Signature *crypto.Signature `cbor:"6,keyasint"`
}

func NewPrevote(height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
	return NewVote(VoteTypePrevote, height, round, blockHash, signer)
}

func NewPrecommit(height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
	return NewVote(VoteTypePrecommit, height, round, blockHash, signer)

}

func NewVote(voteType VoteType, height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
	return &Vote{
		data: voteData{
			VoteType:  voteType,
			Height:    height,
			Round:     round,
			BlockHash: blockHash,
			Signer:    signer,
		},
	}
}

func (vote *Vote) VoteType() VoteType           { return vote.data.VoteType }
func (vote *Vote) Height() int                  { return vote.data.Height }
func (vote *Vote) Round() int                   { return vote.data.Round }
func (vote *Vote) BlockHash() crypto.Hash       { return vote.data.BlockHash }
func (vote *Vote) Signer() crypto.Address       { return vote.data.Signer }
func (vote *Vote) Signature() *crypto.Signature { return vote.data.Signature }

func (vote *Vote) SetSignature(signature *crypto.Signature) {
	vote.data.Signature = signature
}

type signVote struct {
	VoteType  VoteType    `cbor:"1,keyasint"`
	BlockHash crypto.Hash `cbor:"2,keyasint"`
	Round     int         `cbor:"3,keyasint"`
}

func CommitSignBytes(blockHash crypto.Hash, round int) []byte {
	bz, _ := cbor.Marshal(signVote{
		VoteType:  VoteTypePrecommit,
		Round:     round,
		BlockHash: blockHash,
	})

	return bz
}

func (vote Vote) SignBytes() []byte {
	// Note:
	// We omit block height, because finally block height is not matter, block hash is matter
	bz, _ := cbor.Marshal(signVote{
		VoteType:  vote.data.VoteType,
		Round:     vote.data.Round,
		BlockHash: vote.data.BlockHash,
	})

	return bz
}

func (vote *Vote) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(vote.data)
}

func (vote *Vote) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &vote.data)
}

func (vote *Vote) Hash() crypto.Hash {
	bz, _ := cbor.Marshal(vote.data)
	return crypto.HashH(bz)
}

func (vote *Vote) Verify(pubKey crypto.PublicKey) error {
	if vote.data.Signature == nil {
		return errors.Errorf(errors.ErrInvalidVote, "No signature")
	}
	if !pubKey.Address().EqualsTo(vote.data.Signer) {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signer")
	}
	if !pubKey.Verify(vote.SignBytes(), vote.data.Signature) {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid signature")
	}
	return nil
}

func (vote *Vote) SanityCheck() error {
	if !vote.data.VoteType.IsValid() {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid vote type")
	}
	if vote.data.Height < 0 {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid height")
	}
	if vote.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid round")
	}
	if vote.data.Signer.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signer")
	}
	if vote.data.Signature == nil {
		return errors.Errorf(errors.ErrInvalidVote, "No signature")
	}
	if vote.data.Signature.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signature")
	}
	return nil
}

func (vote Vote) Fingerprint() string {
	return fmt.Sprintf("{%v/%d/%s ⌘ %v 👤 %s 🖊 %s}",
		vote.data.Height,
		vote.data.Round,
		vote.data.VoteType,
		vote.data.BlockHash.Fingerprint(),
		vote.data.Signer.Fingerprint(),
		vote.data.Signature.Fingerprint(),
	)
}

// ---------
// For tests
func GenerateTestPrecommitVote(height, round int) (*Vote, crypto.PrivateKey) {
	addr, _, pv := crypto.GenerateTestKeyPair()
	v := NewPrecommit(
		height,
		round,
		crypto.GenerateTestHash(),
		addr)
	sig := pv.Sign(v.SignBytes())
	v.SetSignature(sig)

	return v, pv
}

func GenerateTestPrevoteVote(height, round int) (*Vote, crypto.PrivateKey) {
	addr, _, pv := crypto.GenerateTestKeyPair()
	v := NewPrevote(
		height,
		round,
		crypto.GenerateTestHash(),
		addr)
	sig := pv.Sign(v.SignBytes())
	v.SetSignature(sig)

	return v, pv
}
