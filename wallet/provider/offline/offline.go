package offlineprovider

import (
	"errors"

	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/wallet/provider"
)

var _ provider.IBlockchainProvider = (*OfflineBlockchainProvider)(nil)

// ErrOffline describes an error in which the wallet is offline.
var ErrOffline = errors.New("wallet is in offline mode")

type OfflineBlockchainProvider struct{}

func NewOfflineBlockchainProvider() *OfflineBlockchainProvider {
	return &OfflineBlockchainProvider{}
}

func (*OfflineBlockchainProvider) LastBlockHeight() (block.Height, error) {
	return 0, ErrOffline
}

func (*OfflineBlockchainProvider) GetAccount(string) (*account.Account, error) {
	return nil, ErrOffline
}

func (*OfflineBlockchainProvider) GetValidator(string) (*validator.Validator, error) {
	return nil, ErrOffline
}

func (*OfflineBlockchainProvider) GetTransaction(string) (*tx.Tx, block.Height, error) {
	return nil, 0, ErrOffline
}

func (*OfflineBlockchainProvider) SendTx(*tx.Tx) (string, error) {
	return "", ErrOffline
}

func (*OfflineBlockchainProvider) Close() error {
	return nil
}
