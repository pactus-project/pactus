package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type Commit struct {
	data commitData
}
type commitData struct {
	Round      int                `cbor:"1,keyasint"`
	Commiters  []crypto.Address   `cbor:"2,keyasint"`
	Signatures []crypto.Signature `cbor:"3,keyasint"`
}

func NewCommit(round int, commiters []crypto.Address, signatures []crypto.Signature) *Commit {
	return &Commit{
		data: commitData{
			Round:      round,
			Commiters:  commiters,
			Signatures: signatures,
		},
	}
}

func (commit *Commit) Round() int                     { return commit.data.Round }
func (commit *Commit) Commiters() []crypto.Address    { return commit.data.Commiters }
func (commit *Commit) Signatures() []crypto.Signature { return commit.data.Signatures }

func (commit *Commit) SanityCheck() error {
	if commit.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}

	if len(commit.data.Commiters) != len(commit.data.Signatures) {
		return errors.Errorf(errors.ErrInvalidBlock, "Not enough signatures")
	}

	return nil
}

func (commit *Commit) Size() int {
	if commit == nil {
		return 0
	}
	return len(commit.data.Signatures)
}

func (commit *Commit) Hash() crypto.Hash {
	if commit == nil {
		return crypto.UndefHash
	}
	bs, err := commit.MarshalCBOR()
	if err != nil {
		return crypto.UndefHash
	}
	return crypto.HashH(bs)
}

func (commit *Commit) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(commit.data)
}

func (commit *Commit) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &commit.data)
}

func (commit Commit) MarshalJSON() ([]byte, error) {
	return json.Marshal(commit.data)
}

func (commit *Commit) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &commit.data)
}
