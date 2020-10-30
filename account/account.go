package account

import (
	"encoding/json"
	"math/rand"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	e "github.com/zarbchain/zarb-go/errors"
)

// Account structure
type Account struct {
	data accountData
}

type accountData struct {
	Address  crypto.Address `cbor:"1,keyasint"`
	Sequence int            `cbor:"2,keyasint"`
	Balance  int64          `cbor:"3,keyasint"`
	Code     []byte         `cbor:"4,keyasint"`
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
func (acc Account) Sequence() int           { return acc.data.Sequence }
func (acc Account) Balance() int64          { return acc.data.Balance }
func (acc Account) Code() []byte            { return acc.data.Code }

func (acc *Account) SetBalance(bal int64) error {
	acc.data.Balance = bal
	return nil
}

func (acc *Account) SubtractFromBalance(amt int64) error {
	if amt > acc.Balance() {
		return errors.Errorf(e.ErrInsufficientFunds, "Attempt to subtract %v from the balance of %s", amt, acc.Address())
	}
	acc.data.Balance -= amt
	return nil
}

func (acc *Account) AddToBalance(amt int64) error {
	acc.data.Balance += amt
	return nil
}

func (acc *Account) SetCode(code []byte) error {
	acc.data.Code = code
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

func (acc Account) String() string {
	b, _ := acc.MarshalJSON()
	return string(b)
}

// ---------
// For tests
func GenerateTestAccount() *Account {
	a, _, _ := crypto.GenerateTestKeyPair()
	acc := NewAccount(a)
	acc.data.Balance = rand.Int63n(100000)
	acc.data.Sequence = rand.Intn(100)
	return acc
}
