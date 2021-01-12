package block

import (
	"encoding/json"
	"sort"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Commit struct {
	data commitData
}
type commitData struct {
	BlockHash crypto.Hash      `cbor:"1,keyasint"`
	Round     int              `cbor:"2,keyasint"`
	Signed    []int            `cbor:"3,keyasint"` // validator numbers that signed the commit
	Missed    []int            `cbor:"4,keyasint"` // validator numbers that missed the commit
	Signature crypto.Signature `cbor:"5,keyasint"`
}

func NewCommit(blockHash crypto.Hash, round int, signed, missed []int, signature crypto.Signature) *Commit {
	return &Commit{
		data: commitData{
			BlockHash: blockHash,
			Round:     round,
			Signed:    signed,
			Missed:    missed,
			Signature: signature,
		},
	}
}

func (c *Commit) BlockHash() crypto.Hash      { return c.data.BlockHash }
func (c *Commit) Round() int                  { return c.data.Round }
func (c *Commit) Signed() []int               { return c.data.Signed }
func (c *Commit) Missed() []int               { return c.data.Missed }
func (c *Commit) Signature() crypto.Signature { return c.data.Signature }

func (c *Commit) Committers() []int {
	nums := make([]int, len(c.data.Signed))
	copy(nums, c.data.Signed)
	nums = append(nums, c.data.Missed...)
	sort.Ints(nums)

	return nums
}

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
	return len(c.data.Signed) * 100 / (len(c.data.Missed) + len(c.data.Signed))
}

func (c *Commit) HasTwoThirdThreshold() bool {
	return c.Threshold() > (2 * 100 / 3)
}

func (c *Commit) CommittersHash() crypto.Hash {
	numbers := c.Committers()
	data := make([]byte, 0)
	for _, n := range numbers {
		data = append(data, util.IntToSlice(n)...)
	}
	return crypto.HashH(data)
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
