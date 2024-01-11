package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
)

// Vote is a struct that represents a consensus vote.
type Vote struct {
	data voteData
}

type voteData struct {
	Type      Type           `cbor:"1,keyasint"`
	Height    uint32         `cbor:"2,keyasint"`
	Round     int16          `cbor:"3,keyasint"`
	BlockHash hash.Hash      `cbor:"4,keyasint"`
	Signer    crypto.Address `cbor:"5,keyasint"`
	CPVote    *cpVote        `cbor:"6,keyasint"`
	Signature *bls.Signature `cbor:"7,keyasint"`
}

// NewPrepareVote creates a new PREPARE with the specified parameters.
func NewPrepareVote(blockHash hash.Hash, height uint32, round int16, signer crypto.Address) *Vote {
	return newVote(VoteTypePrepare, blockHash, height, round, signer)
}

// NewPrecommitVote creates a new PRECOMMIT with the specified parameters.
func NewPrecommitVote(blockHash hash.Hash, height uint32, round int16, signer crypto.Address) *Vote {
	return newVote(VoteTypePrecommit, blockHash, height, round, signer)
}

// NewCPPreVote creates a new cp:PRE-VOTE with the specified parameters.
func NewCPPreVote(blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpValue CPValue, just Just, signer crypto.Address,
) *Vote {
	v := newVote(VoteTypeCPPreVote, blockHash, height, round, signer)
	v.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return v
}

// NewCPMainVote creates a new cp:MAIN-VOTE with the specified parameters.
func NewCPMainVote(blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpValue CPValue, just Just, signer crypto.Address,
) *Vote {
	v := newVote(VoteTypeCPMainVote, blockHash, height, round, signer)
	v.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return v
}

// NewCPDecidedVote creates a new cp:Decided with the specified parameters.
func NewCPDecidedVote(blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpValue CPValue, just Just, signer crypto.Address,
) *Vote {
	v := newVote(VoteTypeCPDecided, blockHash, height, round, signer)
	v.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return v
}

// newVote creates a new vote with the specified parameters.
func newVote(voteType Type, blockHash hash.Hash, height uint32, round int16,
	signer crypto.Address,
) *Vote {
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

// SignBytes generates the bytes to be signed for the vote.
func (v *Vote) SignBytes() []byte {
	sb := certificate.BlockCertificateSignBytes(v.data.BlockHash, v.data.Height, v.data.Round)
	switch t := v.Type(); t {
	case VoteTypePrecommit:
		// Nothing

	case VoteTypePrepare:
		sb = append(sb, util.StringToBytes(t.String())...)

	case VoteTypeCPPreVote, VoteTypeCPMainVote, VoteTypeCPDecided:
		sb = append(sb, util.StringToBytes(t.String())...)
		sb = append(sb, util.Int16ToSlice(v.data.CPVote.Round)...)
		sb = append(sb, byte(v.data.CPVote.Value))
	}

	return sb
}

// Type returns the type of the vote.
func (v *Vote) Type() Type {
	return v.data.Type
}

// Height returns the height of the block in the vote.
func (v *Vote) Height() uint32 {
	return v.data.Height
}

// Round returns the round the vote.
func (v *Vote) Round() int16 {
	return v.data.Round
}

// CPRound returns the change proposer round the vote.
func (v *Vote) CPRound() int16 {
	return v.data.CPVote.Round
}

// CPValue returns the change proposer value the vote.
func (v *Vote) CPValue() CPValue {
	return v.data.CPVote.Value
}

// CPValue returns the change proposer justification for the vote.
func (v *Vote) CPJust() Just {
	return v.data.CPVote.Just
}

// BlockHash returns the hash of the block in the vote.
func (v *Vote) BlockHash() hash.Hash {
	return v.data.BlockHash
}

// Signer returns the address of the signer of the vote.
func (v *Vote) Signer() crypto.Address {
	return v.data.Signer
}

// Signature returns the signature of the vote.
func (v *Vote) Signature() *bls.Signature {
	return v.data.Signature
}

// SetSignature sets the signature of the vote.
func (v *Vote) SetSignature(sig *bls.Signature) {
	v.data.Signature = sig
}

// MarshalCBOR marshals the vote into CBOR format.
func (v *Vote) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(v.data)
}

// UnmarshalCBOR unmarshals the vote from CBOR format.
func (v *Vote) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &v.data)
}

// Hash calculates the hash of the vote.
func (v *Vote) Hash() hash.Hash {
	bz, _ := cbor.Marshal(v.data)

	return hash.CalcHash(bz)
}

// Verify checks the signature of the vote with the given public key.
func (v *Vote) Verify(pubKey *bls.PublicKey) error {
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidVote, "no signature")
	}

	if v.Signer() != pubKey.ValidatorAddress() {
		return InvalidSignerError{
			Expected: pubKey.ValidatorAddress(),
			Got:      v.Signer(),
		}
	}

	return pubKey.Verify(v.SignBytes(), v.Signature())
}

func (v *Vote) IsCPVote() bool {
	if v.data.Type == VoteTypeCPPreVote ||
		v.data.Type == VoteTypeCPMainVote ||
		v.data.Type == VoteTypeCPDecided {
		return true
	}

	return false
}

// BasicCheck performs a basic check on the vote.
func (v *Vote) BasicCheck() error {
	if !v.data.Type.IsValid() {
		return errors.Errorf(errors.ErrInvalidVote, "invalid vote type")
	}
	if v.data.Height <= 0 {
		return BasicCheckError{
			Reason: "invalid height",
		}
	}
	if v.data.Round < 0 {
		return BasicCheckError{
			Reason: "invalid round",
		}
	}
	if v.IsCPVote() {
		if v.data.CPVote == nil {
			return errors.Errorf(errors.ErrInvalidVote, "should have cp data")
		}
		if err := v.data.CPVote.BasicCheck(); err != nil {
			return err
		}
	} else if v.data.CPVote != nil {
		return errors.Errorf(errors.ErrInvalidVote, "should not have cp data")
	}
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidSignature, "no signature")
	}

	return nil
}

func (v *Vote) String() string {
	switch v.Type() {
	case VoteTypePrepare, VoteTypePrecommit:
		return fmt.Sprintf("{%d/%d/%s ⌘ %v 👤 %s}",
			v.Height(),
			v.Round(),
			v.Type(),
			v.BlockHash().ShortString(),
			v.Signer().ShortString(),
		)
	case VoteTypeCPPreVote, VoteTypeCPMainVote, VoteTypeCPDecided:
		return fmt.Sprintf("{%d/%d/%s/%d/%d ⌘ %v 👤 %s}",
			v.Height(),
			v.Round(),
			v.Type(),
			v.CPRound(),
			v.CPValue(),
			v.BlockHash().ShortString(),
			v.Signer().ShortString(),
		)

	default:
		return "unknown type"
	}
}
