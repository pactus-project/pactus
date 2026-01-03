package localprovider

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/wallet/provider"
)

var _ provider.IBlockchainProvider = (*LocalBlockchainProvider)(nil)

type LocalBlockchainProvider struct {
	state state.Facade
}

func NewLocalBlockchainProvider(state state.Facade) *LocalBlockchainProvider {
	return &LocalBlockchainProvider{
		state: state,
	}
}

func (p *LocalBlockchainProvider) LastBlockHeight() (block.Height, error) {
	return block.Height(p.state.LastBlockHeight()), nil
}

func (p *LocalBlockchainProvider) GetAccount(addrStr string) (*account.Account, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return nil, err
	}

	return p.state.AccountByAddress(addr)
}

func (p *LocalBlockchainProvider) GetValidator(addrStr string) (*validator.Validator, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return nil, err
	}

	return p.state.ValidatorByAddress(addr)
}

func (p *LocalBlockchainProvider) GetTransaction(txID string) (*tx.Tx, block.Height, error) {
	idHash, err := hash.FromString(txID)
	if err != nil {
		return nil, 0, err
	}

	cTrx, err := p.state.CommittedTx(tx.ID(idHash))
	if err != nil {
		return nil, 0, err
	}

	trx, err := cTrx.ToTx()
	if err != nil {
		return nil, 0, err
	}

	return trx, block.Height(cTrx.Height), nil
}

func (p *LocalBlockchainProvider) SendTx(trx *tx.Tx) (string, error) {
	if err := p.state.AddPendingTxAndBroadcast(trx); err != nil {
		return "", err
	}

	return trx.ID().String(), nil
}

func (p *LocalBlockchainProvider) Close() error {
	return nil
}
