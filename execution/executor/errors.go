package executor

import "errors"

// ErrInsufficientFunds indicates the balance is low for the transaction.
var ErrInsufficientFunds = errors.New("insufficient funds")
