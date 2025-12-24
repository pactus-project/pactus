package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
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
	vote := newVote(VoteTypeCPPreVote, blockHash, height, round, signer)
	vote.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return vote
}

// NewCPMainVote creates a new cp:MAIN-VOTE with the specified parameters.
func NewCPMainVote(blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpValue CPValue, just Just, signer crypto.Address,
) *Vote {
	vote := newVote(VoteTypeCPMainVote, blockHash, height, round, signer)
	vote.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return vote
}

// NewCPDecidedVote creates a new cp:Decided with the specified parameters.
func NewCPDecidedVote(blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpValue CPValue, just Just, signer crypto.Address,
) *Vote {
	vote := newVote(VoteTypeCPDecided, blockHash, height, round, signer)
	vote.data.CPVote = &cpVote{
		Round: cpRound,
		Value: cpValue,
		Just:  just,
	}

	return vote
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
	cert := certificate.NewCertificate(v.data.Height, v.data.Round)

	switch typ := v.Type(); typ {
	case VoteTypePrecommit:
		return cert.SignBytesPrecommit(v.data.BlockHash)

	case VoteTypePrepare:
		return cert.SignBytesPrepare(v.data.BlockHash)

	case VoteTypeCPPreVote:
		return cert.SignBytesCPPreVote(v.data.BlockHash, v.data.CPVote.Round, byte(v.data.CPVote.Value))

	case VoteTypeCPMainVote:
		return cert.SignBytesCPMainVote(v.data.BlockHash, v.data.CPVote.Round, byte(v.data.CPVote.Value))

	case VoteTypeCPDecided:
		return cert.SignBytesCPDecided(v.data.BlockHash, v.data.CPVote.Round, byte(v.data.CPVote.Value))

	default:
		return nil
	}
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

// CPJust returns the change proposer justification for the vote.
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

// MarshalCBOR encodes the vote into CBOR format.
func (v *Vote) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(v.data)
}

// UnmarshalCBOR decodes the vote from CBOR format.
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
		return BasicCheckError{
			Reason: "invalid vote type",
		}
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
			return BasicCheckError{
				Reason: "should have CP data",
			}
		}
		if err := v.data.CPVote.BasicCheck(); err != nil {
			return err
		}
	} else if v.data.CPVote != nil {
		return BasicCheckError{
			Reason: "should not have CP data",
		}
	}
	if v.Signature() == nil {
		return BasicCheckError{
			Reason: "no signature",
		}
	}

	return nil
}

// LogString returns a concise string representation intended for use in logs.
func (v *Vote) LogString() string {
	switch v.Type() {
	case VoteTypePrepare, VoteTypePrecommit:
		return fmt.Sprintf("{%d/%d/%s ⌘ %v 👤 %s}",
			v.Height(),
			v.Round(),
			v.Type(),
			v.BlockHash().LogString(),
			v.Signer().LogString(),
		)
	case VoteTypeCPPreVote, VoteTypeCPMainVote, VoteTypeCPDecided:
		return fmt.Sprintf("{%d/%d/%s/%d/%s ⌘ %v 👤 %s}",
			v.Height(),
			v.Round(),
			v.Type(),
			v.CPRound(),
			v.CPValue(),
			v.BlockHash().LogString(),
			v.Signer().LogString(),
		)

	default:
		return "unknown type"
	}
}
