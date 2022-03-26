package account

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/util"
)

// Account structure
type Account struct {
	data accountData
}

type accountData struct {
	Address  crypto.Address
	Number   int32
	Sequence int32
	Balance  int64
}

// NewAccount constructs a new account object
func NewAccount(addr crypto.Address, number int32) *Account {
	return &Account{
		data: accountData{
			Address: addr,
			Number:  number,
		},
	}
}

// FromBytes constructs a new account from byte array
func FromBytes(data []byte) (*Account, error) {
	acc := new(Account)
	r := bytes.NewReader(data)
	err := encoding.ReadElements(r,
		&acc.data.Address,
		&acc.data.Number,
		&acc.data.Sequence,
		&acc.data.Balance)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (acc Account) Address() crypto.Address { return acc.data.Address }
func (acc Account) Number() int32           { return acc.data.Number }
func (acc Account) Sequence() int32         { return acc.data.Sequence }
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
	bs, err := acc.Bytes()
	if err != nil {
		panic(err)
	}
	return hash.CalcHash(bs)
}
func (acc *Account) SerializeSize() int {
	return 37 // 21+4+4+8
}

func (acc *Account) Bytes() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, acc.SerializeSize()))
	err := encoding.WriteElements(w,
		&acc.data.Address,
		acc.data.Number,
		acc.data.Sequence,
		acc.data.Balance)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
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
func GenerateTestAccount(number int32) (*Account, crypto.Signer) {
	signer := bls.GenerateTestSigner()
	acc := NewAccount(signer.Address(), number)
	acc.data.Balance = util.RandInt64(100 * 1e14)
	acc.data.Sequence = util.RandInt32(1000)
	return acc, signer
}
