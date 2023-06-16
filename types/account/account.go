// Package account provides functionality for managing account information.
package account

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

// The Account struct represents a account object.
type Account struct {
	data accountData
}

// accountData contains the data associated with a account.
type accountData struct {
	Number   int32
	Sequence int32
	Balance  int64
}

// NewAccount constructs a new account from the given number.
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

// Number returns the number of the account.
func (acc Account) Number() int32 {
	return acc.data.Number
}

// Sequence returns the sequence number of the account.
func (acc Account) Sequence() int32 {
	return acc.data.Sequence
}

// Balance returns the balance of the account.
func (acc Account) Balance() int64 {
	return acc.data.Balance
}

// SubtractFromBalance subtracts the given amount from the account's balance.
func (acc *Account) SubtractFromBalance(amt int64) {
	acc.data.Balance -= amt
}

// AddToBalance adds the given amount to the account's balance.
func (acc *Account) AddToBalance(amt int64) {
	acc.data.Balance += amt
}

// IncSequence increases the sequence anytime this account signs a transaction.
func (acc *Account) IncSequence() {
	acc.data.Sequence++
}

// Hash calculates and returns the hash of the account.
func (acc *Account) Hash() hash.Hash {
	bs, err := acc.Bytes()
	if err != nil {
		panic(err)
	}
	return hash.CalcHash(bs)
}

// SerializeSize returns the size in bytes required to serialize the account.
func (acc *Account) SerializeSize() int {
	return 16 // 4+4+8
}

// Bytes returns the serialized byte representation of the account.
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

// Clone creates a deep copy of the account.
func (acc *Account) Clone() *Account {
	cloned := new(Account)
	*cloned = *acc
	return cloned
}

// GenerateTestAccount generates an account for testing purposes.
func GenerateTestAccount(number int32) (*Account, crypto.Signer) {
	signer := bls.GenerateTestSigner()
	acc := NewAccount(number)
	acc.data.Balance = util.RandInt64(100 * 1e14)
	acc.data.Sequence = util.RandInt32(1000)
	return acc, signer
}
