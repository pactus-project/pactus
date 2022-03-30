package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/block"
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
	Type      Type           `cbor:"1,keyasint"`
	Height    int32          `cbor:"2,keyasint"`
	Round     int16          `cbor:"3,keyasint"`
	BlockHash hash.Hash      `cbor:"4,keyasint"`
	Signer    crypto.Address `cbor:"5,keyasint"`
	Signature *bls.Signature `cbor:"6,keyasint"`
}

func (v *Vote) SignBytes() []byte {
	// Note:
	// We omit block height, because finally block height is not matter, block hash is matter
	sb := block.CertificateSignBytes(v.data.BlockHash, v.data.Round)
	if v.Type() == VoteTypePrepare {
		sb = append(sb, []byte("prepare")...)
	} else if v.Type() == VoteTypeChangeProposer {
		sb = append(sb, []byte("change-proposer")...)
	}

	return sb
}

func NewVote(voteType Type, height int32, round int16, blockHash hash.Hash, signer crypto.Address) *Vote {
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

func (v *Vote) Type() Type                { return v.data.Type }
func (v *Vote) Height() int32             { return v.data.Height }
func (v *Vote) Round() int16              { return v.data.Round }
func (v *Vote) BlockHash() hash.Hash      { return v.data.BlockHash }
func (v *Vote) Signer() crypto.Address    { return v.data.Signer }
func (v *Vote) Signature() *bls.Signature { return v.data.Signature }

func (v *Vote) SetSignature(sig crypto.Signature) {
	v.data.Signature = sig.(*bls.Signature)
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
	return hash.CalcHash(bz)
}

func (v *Vote) Verify(pubKey *bls.PublicKey) error {
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidVote, "no signature")
	}
	if err := pubKey.VerifyAddress(v.Signer()); err != nil {
		return err
	}
	return pubKey.Verify(v.SignBytes(), v.Signature())
}

func (v *Vote) SanityCheck() error {
	if !v.data.Type.IsValid() {
		return errors.Errorf(errors.ErrInvalidVote, "invalid vote type")
	}
	if v.data.Height <= 0 {
		return errors.Error(errors.ErrInvalidHeight)
	}
	if v.data.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}
	if err := v.data.Signer.SanityCheck(); err != nil {
		return err
	}
	if v.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidVote, "no signature")
	}
	if err := v.Signature().SanityCheck(); err != nil {
		return err
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
func GenerateTestPrecommitVote(height int32, round int16) (*Vote, crypto.Signer) {
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

func GenerateTestPrepareVote(height int32, round int16) (*Vote, crypto.Signer) {
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

func GenerateTestChangeProposerVote(height int32, round int16) (*Vote, crypto.Signer) {
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
