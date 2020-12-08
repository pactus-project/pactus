package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

type Commiter struct {
	Address crypto.Address `cbor:"1,keyasint"`
	Signed  bool           `cbor:"2,keyasint"`
}

type Commit struct {
	data commitData
}
type commitData struct {
	Round     int              `cbor:"1,keyasint"`
	Signature crypto.Signature `cbor:"2,keyasint"`
	Commiters []Commiter       `cbor:"3,keyasint"`
}

func NewCommit(round int, commiters []Commiter, signature crypto.Signature) *Commit {
	return &Commit{
		data: commitData{
			Round:     round,
			Commiters: commiters,
			Signature: signature,
		},
	}
}

func (commit *Commit) Round() int                  { return commit.data.Round }
func (commit *Commit) Commiters() []Commiter       { return commit.data.Commiters }
func (commit *Commit) Signature() crypto.Signature { return commit.data.Signature }

func (commit *Commit) SanityCheck() error {
	if commit.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}
	if err := commit.data.Signature.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	signed := 0
	for _, c := range commit.data.Commiters {
		if c.Signed {
			signed++
		}
	}

	if signed <= (len(commit.data.Commiters) * 2 / 3) {
		return errors.Errorf(errors.ErrInvalidBlock, "Not enough commiters")
	}

	return nil
}

func (commit *Commit) Size() int {
	if commit == nil {
		return 0
	}
	return len(commit.data.Commiters)
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

func (commit *Commit) CommitersHash() crypto.Hash {
	data := make([][]byte, len(commit.data.Commiters))

	for i, c := range commit.data.Commiters {
		data[i] = make([]byte, 20)
		copy(data[i], c.Address.RawBytes())
	}
	merkle := simpleMerkle.NewTreeFromSlices(data)

	return merkle.Root()
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
