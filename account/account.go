package account

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
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

// NewAccount constructs a new account object
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

func (acc *Account) Hash() hash.Hash {
	bs, err := acc.Encode()
	if err != nil {
		panic(err)
	}
	return hash.CalcHash(bs)
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

func (acc Account) Fingerprint() string {
	return fmt.Sprintf("{ %s %d %v}",
		acc.Address().Fingerprint(),
		acc.Sequence(),
		acc.Balance())
}

// GenerateTestAccount generates an account for testing purpose
func GenerateTestAccount(number int) (*Account, crypto.Signer) {
	signer := bls.GenerateTestSigner()
	acc := NewAccount(signer.Address(), number)
	acc.data.Balance = util.RandInt64(100 * 1e14)
	acc.data.Sequence = util.RandInt(100)
	return acc, signer
}
