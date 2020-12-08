package store

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type StoreReader interface {
	BlockByHeight(height int) (*block.Block, error)
	BlockByHash(hash crypto.Hash) (*block.Block, int, error)
	BlockHeight(hash crypto.Hash) (int, error)
	Tx(hash crypto.Hash) (*tx.Tx, *tx.Receipt, error)
	HasAccount(crypto.Address) bool
	Account(addr crypto.Address) (*account.Account, error)
	HasValidator(crypto.Address) bool
	Validator(addr crypto.Address) (*validator.Validator, error)
}
