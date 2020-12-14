package validator

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

type Validator struct {
	data validatorData
}

type validatorData struct {
	PublicKey     crypto.PublicKey `cbor:"1,keyasint"`
	Number        int              `cbor:"2,keyasint"`
	Sequence      int              `cbor:"3,keyasint"`
	Stake         int64            `cbor:"4,keyasint"`
	BondingHeight int              `cbor:"5,keyasint"`
}

func NewValidator(publicKey crypto.PublicKey, number, bondingHeight int) *Validator {
	val := &Validator{
		data: validatorData{
			PublicKey:     publicKey,
			Number:        number,
			BondingHeight: bondingHeight,
		},
	}
	return val
}

func (val *Validator) PublicKey() crypto.PublicKey { return val.data.PublicKey }
func (val *Validator) Address() crypto.Address     { return val.data.PublicKey.Address() }
func (val *Validator) Number() int                 { return val.data.Number }
func (val *Validator) Sequence() int               { return val.data.Sequence }
func (val *Validator) Stake() int64                { return val.data.Stake }
func (val *Validator) BondingHeight() int          { return val.data.BondingHeight }

func (val Validator) Power() int64 {
	// Viva democracy, everybody should be treated equally
	return 1
}

func (val *Validator) SubtractFromStake(amt int64) {
	val.data.Stake -= amt
}

func (val *Validator) AddToStake(amt int64) {
	val.data.Stake += amt
}

func (val *Validator) IncSequence() {
	val.data.Sequence++
}

func (val *Validator) Hash() crypto.Hash {
	bs, err := val.Encode()
	if err != nil {
		panic(err)
	}
	return crypto.HashH(bs)
}

///---- Serialization methods
func (val Validator) Encode() ([]byte, error) {
	return cbor.Marshal(val.data)
}

func (val *Validator) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &val.data)
}

func (val Validator) MarshalJSON() ([]byte, error) {
	return json.Marshal(val.data)
}

func (val *Validator) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &val.data)
}

func (val Validator) Fingerprint() string {
	return fmt.Sprintf("{%s %v}",
		val.Address().Fingerprint(),
		val.Stake())
}

// GenerateTestValidator generates a validator for testing purpose
func GenerateTestValidator(number int) (*Validator, crypto.PrivateKey) {
	_, pub, priv := crypto.GenerateTestKeyPair()
	val := NewValidator(pub, number, util.RandInt(100))
	val.data.Stake = util.RandInt64(1000000000)
	val.data.Sequence = util.RandInt(1000)
	return val, priv
}
