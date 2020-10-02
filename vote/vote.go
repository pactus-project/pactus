package vote

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/errors"
)

// Vote represents a prevote, precommit, or commit vote from validators for
// consensus.
type Vote struct {
	data voteData
}

type voteData struct {
	Type      VoteType         `cbor:"1,keyasint"`
	Height    int              `cbor:"2,keyasint"`
	Round     int              `cbor:"3,keyasint"`
	BlockHash crypto.Hash      `cbor:"4,keyasint"`
	Signer    crypto.Address   `cbor:"5,keyasint"`
	Signature crypto.Signature `cbor:"6,keyasint"`
}

func NewPrevoteVote(height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
	return NewVote(VoteTypePrevote, height, round, blockHash, signer)
}

func NewPrecommitVote(height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
	return NewVote(VoteTypePrecommit, height, round, blockHash, signer)

}
func NewVote(voteType VoteType, height int, round int, blockHash crypto.Hash, signer crypto.Address) *Vote {
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

func (vote *Vote) Type() VoteType              { return vote.data.Type }
func (vote *Vote) Height() int                 { return vote.data.Height }
func (vote *Vote) Round() int                  { return vote.data.Round }
func (vote *Vote) BlockHash() crypto.Hash      { return vote.data.BlockHash }
func (vote *Vote) Signer() crypto.Address      { return vote.data.Signer }
func (vote *Vote) Signature() crypto.Signature { return vote.data.Signature }

func (vote *Vote) SetSignature(signature crypto.Signature) { vote.data.Signature = signature }

func (vote Vote) SignBytes() []byte {
	type signVote struct {
		Type      VoteType       `cbor:"1,keyasint"`
		Height    int            `cbor:"2,keyasint"`
		Round     int            `cbor:"3,keyasint"`
		BlockHash crypto.Hash    `cbor:"4,keyasint"`
		Signer    crypto.Address `cbor:"5,keyasint"`
	}

	bz, _ := json.Marshal(signVote{
		Type:      vote.data.Type,
		Height:    vote.data.Height,
		Round:     vote.data.Round,
		BlockHash: vote.data.BlockHash,
		Signer:    vote.data.Signer,
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
	return crypto.HashH(vote.SignBytes())
}

func (vote *Vote) Verify(pubKey crypto.PublicKey) error {
	if !pubKey.Address().EqualsTo(vote.data.Signer) {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signer")
	}
	if !pubKey.Verify(vote.SignBytes(), vote.data.Signature) {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid signature")
	}
	return nil
}

func (vote *Vote) SanityCheck() error {
	if !vote.data.Type.IsValid() {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid type")
	}
	if vote.data.Height < 0 {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid height")
	}
	if vote.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid round")
	}
	if err := vote.data.BlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidVote, err.Error())
	}
	if vote.data.Signer.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signer")
	}
	if vote.data.Signature.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidVote, "Invalid signature")
	}
	return nil
}

func (vote Vote) Fingerprint() string {
	return fmt.Sprintf("{%v/%d/%s S:%s B:%s}",
		vote.data.Height,
		vote.data.Round,
		vote.data.Type,
		vote.data.Signer.Fingerprint(),
		vote.data.BlockHash.Fingerprint(),
	)
}
