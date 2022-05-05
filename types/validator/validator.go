package validator

import (
	"bytes"

	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/encoding"
)

type Validator struct {
	data validatorData
}

type validatorData struct {
	PublicKey         *bls.PublicKey
	Number            int32
	Sequence          int32
	Stake             int64
	LastBondingHeight int32
	UnbondingHeight   int32
	LastJoinedHeight  int32
}

// NewValidator constructs a new validator object.
func NewValidator(publicKey *bls.PublicKey, number int32) *Validator {
	val := &Validator{
		data: validatorData{
			PublicKey: publicKey,
			Number:    number,
		},
	}
	return val
}

// FromBytes constructs a new validator from byte array.
func FromBytes(data []byte) (*Validator, error) {
	acc := new(Validator)
	r := bytes.NewReader(data)

	acc.data.PublicKey = new(bls.PublicKey)
	if err := acc.data.PublicKey.Decode(r); err != nil {
		return nil, err
	}

	err := encoding.ReadElements(r,
		&acc.data.Number,
		&acc.data.Sequence,
		&acc.data.Stake,
		&acc.data.LastBondingHeight,
		&acc.data.UnbondingHeight,
		&acc.data.LastJoinedHeight,
	)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (val *Validator) PublicKey() *bls.PublicKey { return val.data.PublicKey }
func (val *Validator) Address() crypto.Address   { return val.data.PublicKey.Address() }
func (val *Validator) Number() int32             { return val.data.Number }
func (val *Validator) Sequence() int32           { return val.data.Sequence }
func (val *Validator) Stake() int64              { return val.data.Stake }
func (val *Validator) LastBondingHeight() int32  { return val.data.LastBondingHeight }
func (val *Validator) UnbondingHeight() int32    { return val.data.UnbondingHeight }
func (val *Validator) LastJoinedHeight() int32   { return val.data.LastJoinedHeight }

func (val Validator) Power() int64 {
	//if the validator requested to unbond ignore stake
	if val.data.UnbondingHeight > 0 {
		return 0
	} else if val.data.Stake == 0 { // Only bootstrap validators at genesis block has no stake
		return 1
	}
	return val.data.Stake
}

func (val *Validator) SubtractFromStake(amt int64) {
	val.data.Stake -= amt
}

// AddToStake increases the stake by bonding transaction.
func (val *Validator) AddToStake(amt int64) {
	val.data.Stake += amt
}

// IncSequence increases the sequence anytime this validator signs a transaction.
func (val *Validator) IncSequence() {
	val.data.Sequence++
}

// UpdateLastJoinedHeight updates the last height that this validator joined the committee.
func (val *Validator) UpdateLastJoinedHeight(height int32) {
	val.data.LastJoinedHeight = height
}

// UpdateLastBondingHeight updates the last height that this validator bonded some stakes.
func (val *Validator) UpdateLastBondingHeight(height int32) {
	val.data.LastBondingHeight = height
}

// UpdateUnbondingHeight updates the unbonding height for the validator.
func (val *Validator) UpdateUnbondingHeight(height int32) {
	val.data.UnbondingHeight = height
}

// Hash return the hash of this validator.
func (val *Validator) Hash() hash.Hash {
	bs, err := val.Bytes()
	if err != nil {
		panic(err)
	}
	return hash.CalcHash(bs)
}

func (val *Validator) SerializeSize() int {
	return 124 // 96+4+4+8+4+4+4
}

func (val *Validator) Bytes() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, val.SerializeSize()))

	if err := val.data.PublicKey.Encode(w); err != nil {
		return nil, err
	}

	err := encoding.WriteElements(w,
		val.data.Number,
		val.data.Sequence,
		val.data.Stake,
		val.data.LastBondingHeight,
		val.data.UnbondingHeight,
		val.data.LastJoinedHeight)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// GenerateTestValidator generates a validator for testing purpose.
func GenerateTestValidator(number int32) (*Validator, crypto.Signer) {
	pub, pv := bls.GenerateTestKeyPair()
	val := NewValidator(pub, number)
	val.data.Stake = util.RandInt64(100 * 1e8)
	val.data.Sequence = util.RandInt32(100)
	return val, crypto.NewSigner(pv)
}
