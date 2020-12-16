package account

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
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
func NewAccount(addr crypto.Address, number int) *Account {
	return &Account{
		data: accountData{
			Address: addr,
			Number:  number,
		},
	}
}

func (acc Account) Address() crypto.Address { return acc.data.Address }
func (acc Account) Number() int             { return acc.data.Number }
func (acc Account) Sequence() int           { return acc.data.Sequence }
func (acc Account) Balance() int64          { return acc.data.Balance }

func (acc *Account) SubtractFromBalance(amt int64) {
	acc.data.Balance -= amt
}

func (acc *Account) AddToBalance(amt int64) {
	acc.data.Balance += amt
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

// GenerateTestAccount generates an account for testing purpose
func GenerateTestAccount(number int) (*Account, crypto.PrivateKey) {
	a, _, priv := crypto.GenerateTestKeyPair()
	acc := NewAccount(a, number)
	acc.data.Balance = util.RandInt64(10000000)
	acc.data.Sequence = util.RandInt(100)
	return acc, priv
}
