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

func (p *OfflineBlockchainProvider) LastBlockHeight() (block.Height, error) {
	return 0, ErrOffline
}

func (p *OfflineBlockchainProvider) GetAccount(addrStr string) (*account.Account, error) {
	return nil, ErrOffline
}

func (p *OfflineBlockchainProvider) GetValidator(addrStr string) (*validator.Validator, error) {
	return nil, ErrOffline
}

func (p *OfflineBlockchainProvider) GetTransaction(txID string) (*tx.Tx, block.Height, error) {
	return nil, 0, ErrOffline
}

func (p *OfflineBlockchainProvider) SendTx(trx *tx.Tx) (string, error) {
	return "", ErrOffline
}

func (p *OfflineBlockchainProvider) Close() error {
	return nil
}
