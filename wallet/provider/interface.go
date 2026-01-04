package provider

import (
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
)

type IBlockchainProvider interface {
	LastBlockHeight() (block.Height, error)
	GetAccount(addrStr string) (*account.Account, error)
	GetValidator(addrStr string) (*validator.Validator, error)
	GetTransaction(txID string) (*tx.Tx, block.Height, error)

	SendTx(trx *tx.Tx) (string, error)

	Close() error
}
