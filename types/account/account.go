package account

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

// Account represents a structure for an account information.
type Account struct {
	data accountData
}

type accountData struct {
	Number   int32
	Sequence int32
	Balance  int64
}

// NewAccount constructs a new account object.
func NewAccount(number int32) *Account {
	return &Account{
		data: accountData{
			Number: number,
		},
	}
}

// FromBytes constructs a new account from byte array.
func FromBytes(data []byte) (*Account, error) {
	acc := new(Account)
	r := bytes.NewReader(data)
	err := encoding.ReadElements(r,
		&acc.data.Number,
		&acc.data.Sequence,
		&acc.data.Balance)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (acc Account) Number() int32   { return acc.data.Number }
func (acc Account) Sequence() int32 { return acc.data.Sequence }
func (acc Account) Balance() int64  { return acc.data.Balance }

func (acc *Account) SubtractFromBalance(amt int64) {
	acc.data.Balance -= amt
}

func (acc *Account) AddToBalance(amt int64) {
	acc.data.Balance += amt
}

// IncSequence increments the account's sequence every time it signs a transaction.
func (acc *Account) IncSequence() {
	acc.data.Sequence++
}

// Hash returns the hash of this account.
func (acc *Account) Hash() hash.Hash {
	bs, err := acc.Bytes()
	if err != nil {
		panic(err)
	}
	return hash.CalcHash(bs)
}

func (acc *Account) SerializeSize() int {
	return 16 // 4+4+8
}

func (acc *Account) Bytes() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, acc.SerializeSize()))
	err := encoding.WriteElements(w,
		acc.data.Number,
		acc.data.Sequence,
		acc.data.Balance)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// GenerateTestAccount generates an account for testing purpose.
func GenerateTestAccount(number int32) (*Account, crypto.Signer) {
	signer := bls.GenerateTestSigner()
	acc := NewAccount(number)
	acc.data.Balance = util.RandInt64(100 * 1e14)
	acc.data.Sequence = util.RandInt32(1000)
	return acc, signer
}
