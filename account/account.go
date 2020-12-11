package account

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

// Account structure
type Account struct {
	data accountData
}

type accountData struct {
	Address  crypto.Address `cbor:"1,keyasint"`
	Number   int            `cbor:"2,keyasint"`
	Sequence int            `cbor:"3,keyasint"`
	Balance  int64          `cbor:"4,keyasint"`
}

///---- Constructors
func NewAccount(addr crypto.Address) *Account {
	return &Account{
		data: accountData{
			Address: addr,
		},
	}
}

func (acc Account) Address() crypto.Address { return acc.data.Address }
func (acc Account) Number() int             { return acc.data.Number }
func (acc Account) Sequence() int           { return acc.data.Sequence }
func (acc Account) Balance() int64          { return acc.data.Balance }

func (acc *Account) SetBalance(bal int64) error {
	acc.data.Balance = bal
	return nil
}

func (acc *Account) SubtractFromBalance(amt int64) error {
	if amt < 0 {
		return errors.Errorf(errors.ErrInvalidAmount, "amount is negative: %v", amt)
	}
	if amt > acc.Balance() {
		return errors.Errorf(errors.ErrInsufficientFunds, "Attempt to subtract %v from the balance of %s", amt, acc.Address())
	}
	acc.data.Balance -= amt
	return nil
}

func (acc *Account) AddToBalance(amt int64) error {
	if amt < 0 {
		return errors.Errorf(errors.ErrInvalidAmount, "amount is negative: %v", amt)
	}
	acc.data.Balance += amt
	return nil
}

func (acc *Account) IncSequence() {
	acc.data.Sequence++
}

func (acc *Account) Hash() crypto.Hash {
	bs, err := acc.Encode()
	if err != nil {
		panic(err)
	}
	return crypto.HashH(bs)
}

func (acc *Account) Encode() ([]byte, error) {
	return cbor.Marshal(acc.data)
}

func (acc *Account) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &acc.data)
}

func (acc Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(acc.data)
}

func (acc *Account) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &acc.data)
}

func (acc Account) Fingerprint() string {
	return fmt.Sprintf("{ %s %v}",
		acc.Address().Fingerprint(),
		acc.Balance())
}

// ---------
// For tests
func GenerateTestAccount() (*Account, crypto.PrivateKey) {
	a, _, priv := crypto.GenerateTestKeyPair()
	acc := NewAccount(a)
	acc.data.Number = util.RandInt(10000)
	acc.data.Balance = util.RandInt64(10000000)
	acc.data.Sequence = util.RandInt(100)
	return acc, priv
}
