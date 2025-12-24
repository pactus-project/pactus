// Package validator provides functionality for managing validator information.
package validator

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util/encoding"
)

// The Validator struct represents a validator object.
type Validator struct {
	data validatorData
}

// validatorData contains the data associated with a validator.
type validatorData struct {
	PublicKey           *bls.PublicKey
	Number              int32
	Stake               amount.Amount
	LastBondingHeight   uint32
	UnbondingHeight     uint32
	LastSortitionHeight uint32

	// The protocol version of the validator.
	// This is in memory and not saved to the blockchain.
	ProtocolVersion protocol.Version
}

// NewValidator constructs a new validator from the given public key and number.
func NewValidator(publicKey *bls.PublicKey, number int32) *Validator {
	val := &Validator{
		data: validatorData{
			PublicKey: publicKey,
			Number:    number,
		},
	}

	return val
}

// FromBytes constructs a new validator from raw byte data.
func FromBytes(data []byte) (*Validator, error) {
	acc := new(Validator)
	reader := bytes.NewReader(data)

	acc.data.PublicKey = new(bls.PublicKey)
	if err := acc.data.PublicKey.Decode(reader); err != nil {
		return nil, err
	}

	err := encoding.ReadElements(reader,
		&acc.data.Number,
		&acc.data.Stake,
		&acc.data.LastBondingHeight,
		&acc.data.UnbondingHeight,
		&acc.data.LastSortitionHeight,
	)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// PublicKey returns the public key of the validator.
func (val *Validator) PublicKey() *bls.PublicKey {
	return val.data.PublicKey
}

// Address returns the address of the validator.
func (val *Validator) Address() crypto.Address {
	return val.data.PublicKey.ValidatorAddress()
}

// Number returns the number of the validator.
func (val *Validator) Number() int32 {
	return val.data.Number
}

// Stake returns the stake of the validator.
func (val *Validator) Stake() amount.Amount {
	return val.data.Stake
}

// LastBondingHeight returns the last height in which the validator bonded stake.
func (val *Validator) LastBondingHeight() uint32 {
	return val.data.LastBondingHeight
}

// UnbondingHeight returns the last height in which the validator unbonded stake.
func (val *Validator) UnbondingHeight() uint32 {
	return val.data.UnbondingHeight
}

// IsUnbonded returns true if the validator is unbonded.
func (val *Validator) IsUnbonded() bool {
	return val.data.UnbondingHeight > 0
}

// LastSortitionHeight returns the last height in which the validator evaluated sortition.
func (val *Validator) LastSortitionHeight() uint32 {
	return val.data.LastSortitionHeight
}

// Power returns the power of the validator.
func (val Validator) Power() int64 {
	if val.data.UnbondingHeight > 0 {
		// Power for unbonded validators is set to zero.
		return 0
	} else if val.data.Stake == 0 {
		// Only bootstrap validators at the genesis block have no stake.
		return 1
	}

	return int64(val.data.Stake)
}

// SubtractFromStake subtracts the given amount from the validator's stake.
func (val *Validator) SubtractFromStake(amt amount.Amount) {
	val.data.Stake -= amt
}

// AddToStake adds the given amount to the validator's stake.
func (val *Validator) AddToStake(amt amount.Amount) {
	val.data.Stake += amt
}

// UpdateLastSortitionHeight updates the last height at which the validator performed a valid sortition.
func (val *Validator) UpdateLastSortitionHeight(height uint32) {
	val.data.LastSortitionHeight = height
}

// UpdateLastBondingHeight updates the last height at which the validator bonded some stakes.
func (val *Validator) UpdateLastBondingHeight(height uint32) {
	val.data.LastBondingHeight = height
}

// UpdateUnbondingHeight updates the unbonding height for the validator.
func (val *Validator) UpdateUnbondingHeight(height uint32) {
	val.data.UnbondingHeight = height
}

// Hash calculates and returns the hash of the validator.
func (val *Validator) Hash() hash.Hash {
	bs, err := val.Bytes()
	if err != nil {
		panic(err)
	}

	return hash.CalcHash(bs)
}

// SerializeSize returns the size in bytes required to serialize the validator.
func (*Validator) SerializeSize() int {
	return 120 // 96+4+4+8+4+4
}

// Bytes returns the serialized byte representation of the validator.
func (val *Validator) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, val.SerializeSize()))

	if err := val.data.PublicKey.Encode(buf); err != nil {
		return nil, err
	}

	err := encoding.WriteElements(buf,
		val.data.Number,
		val.data.Stake,
		val.data.LastBondingHeight,
		val.data.UnbondingHeight,
		val.data.LastSortitionHeight)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Clone creates a deep copy of the validator.
func (val *Validator) Clone() *Validator {
	cloned := new(Validator)
	*cloned = *val

	return cloned
}

// UpdateProtocolVersion updates the protocol version of the validator.
func (val *Validator) UpdateProtocolVersion(ver protocol.Version) {
	val.data.ProtocolVersion = ver
}

// ProtocolVersion returns the protocol version of the validator.
func (val *Validator) ProtocolVersion() protocol.Version {
	return val.data.ProtocolVersion
}
