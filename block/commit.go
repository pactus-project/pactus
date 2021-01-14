package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

const (
	CommitNotSigned = 0
	CommitSigned    = 1
)

type Committer struct {
	Number int `cbor:"1,keyasint"`
	Status int `cbor:"2,keyasint"`
}

func (committer *Committer) HasSigned() bool {
	return committer.Status == CommitSigned
}

type Commit struct {
	data commitData
}
type commitData struct {
	BlockHash  crypto.Hash      `cbor:"1,keyasint"`
	Round      int              `cbor:"2,keyasint"`
	Committers []Committer      `cbor:"3,keyasint"`
	Signature  crypto.Signature `cbor:"4,keyasint"`
}

func NewCommit(blockHash crypto.Hash, round int, committers []Committer, signature crypto.Signature) *Commit {
	return &Commit{
		data: commitData{
			BlockHash:  blockHash,
			Round:      round,
			Committers: committers,
			Signature:  signature,
		},
	}
}

func (c *Commit) BlockHash() crypto.Hash      { return c.data.BlockHash }
func (c *Commit) Round() int                  { return c.data.Round }
func (c *Commit) Committers() []Committer     { return c.data.Committers }
func (c *Commit) Signature() crypto.Signature { return c.data.Signature }

func (c *Commit) SanityCheck() error {
	if err := c.data.BlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if c.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}
	if err := c.data.Signature.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	for _, c := range c.data.Committers {
		if c.Status > 1 {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid status")
		}
		if c.Number < 0 {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid number")
		}
	}
	if !c.HasTwoThirdThreshold() {
		return errors.Errorf(errors.ErrInvalidBlock, "Not enough signatures")
	}

	return nil
}

func (c *Commit) Hash() crypto.Hash {
	if c == nil {
		return crypto.UndefHash
	}
	bs, err := c.MarshalCBOR()
	if err != nil {
		return crypto.UndefHash
	}
	return crypto.HashH(bs)
}

func (c *Commit) Threshold() int {
	signed := 0
	for _, c := range c.data.Committers {
		if c.Status == CommitSigned {
			signed++
		}
	}
	return signed * 100 / len(c.data.Committers) // divide in golang is floor division
}

func (c *Commit) HasTwoThirdThreshold() bool {
	return c.Threshold() > (2 * 100 / 3)
}

func (c *Commit) CommittersHash() crypto.Hash {
	nums := []int{}
	for _, c := range c.data.Committers {
		nums = append(nums, c.Number)
	}
	bz, _ := cbor.Marshal(nums)
	return crypto.HashH(bz)
}

func (c *Commit) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(c.data)
}

func (c *Commit) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &c.data)
}

func (c *Commit) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

func (c *Commit) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &c.data)
}

type signVote struct {
	BlockHash crypto.Hash `cbor:"1,keyasint"`
	Round     int         `cbor:"2,keyasint"`
}

func (c *Commit) SignBytes() []byte {
	return CommitSignBytes(c.data.BlockHash, c.data.Round)
}

func CommitSignBytes(blockHash crypto.Hash, round int) []byte {
	bz, _ := cbor.Marshal(signVote{
		Round:     round,
		BlockHash: blockHash,
	})

	return bz
}
