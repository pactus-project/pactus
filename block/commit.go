package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

const (
	CommitNotSigned = 0
	CommitSigned    = 1
)

type Committer struct {
	Address crypto.Address `cbor:"1,keyasint"`
	Status  int            `cbor:"2,keyasint"`
}

func (committer *Committer) HasSigned() bool {
	return committer.Status == CommitSigned
}

type Commit struct {
	data commitData
}
type commitData struct {
	Round      int              `cbor:"1,keyasint"`
	Signature  crypto.Signature `cbor:"2,keyasint"`
	Committers []Committer      `cbor:"3,keyasint"`
}

func NewCommit(round int, committers []Committer, signature crypto.Signature) *Commit {
	return &Commit{
		data: commitData{
			Round:      round,
			Committers: committers,
			Signature:  signature,
		},
	}
}

func (commit *Commit) Round() int                  { return commit.data.Round }
func (commit *Commit) Committers() []Committer     { return commit.data.Committers }
func (commit *Commit) Signature() crypto.Signature { return commit.data.Signature }

func (commit *Commit) SanityCheck() error {
	if commit.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}
	if err := commit.data.Signature.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	signed := 0
	for _, c := range commit.data.Committers {
		if c.Status > 1 {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid commit status")
		}
		if c.Status == CommitSigned {
			signed++
		}
	}

	if signed <= (len(commit.data.Committers) * 2 / 3) {
		return errors.Errorf(errors.ErrInvalidBlock, "Not enough committers")
	}

	return nil
}

func (commit *Commit) Size() int {
	if commit == nil {
		return 0
	}
	return len(commit.data.Committers)
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

func (commit *Commit) CommittersHash() crypto.Hash {
	data := make([][]byte, len(commit.data.Committers))

	for i, c := range commit.data.Committers {
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
