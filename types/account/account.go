// Package account provides functionality for managing account information.
package account

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

// The Account struct represents an account object.
type Account struct {
	data accountData
}

// accountData contains the data associated with an account.
type accountData struct {
	Number  int32
	Balance amount.Amount
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

// Balance returns the balance of the account.
func (acc Account) Balance() amount.Amount {
	return acc.data.Balance
}

// SubtractFromBalance subtracts the given amount from the account's balance.
func (acc *Account) SubtractFromBalance(amt amount.Amount) {
	acc.data.Balance -= amt
}

// AddToBalance adds the given amount to the account's balance.
func (acc *Account) AddToBalance(amt amount.Amount) {
	acc.data.Balance += amt
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
	return 12 // 4+8
}

// Bytes returns the serialized byte representation of the account.
func (acc *Account) Bytes() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, acc.SerializeSize()))
	err := encoding.WriteElements(w,
		acc.data.Number,
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
