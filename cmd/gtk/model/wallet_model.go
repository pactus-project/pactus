//go:build gtk

package model

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/types"
)

type WalletModel struct {
	manager    wltmgr.IManager
	walletName string
}

// AddressRow is a UI-friendly but UI-agnostic representation of an address entry.
// Formatting (strings/markup) should be done by presenters/controllers, not here.
type AddressRow struct {
	No       int
	Address  string
	Label    string
	Path     string
	Imported bool
	Balance  amount.Amount
	Stake    amount.Amount
}

func NewWalletModel(manager wltmgr.IManager, walletName string) (*WalletModel, error) {
	return &WalletModel{manager: manager, walletName: walletName}, nil
}

// WalletName returns the display name used in the UI.
func (model *WalletModel) WalletName() string {
	return model.walletName
}

func (model *WalletModel) IsEncrypted() bool {
	info, err := model.manager.WalletInfo(model.walletName)
	if err != nil {
		return false
	}

	return info.Encrypted
}

func (model *WalletModel) WalletInfo() (*types.WalletInfo, error) {
	info, err := model.manager.WalletInfo(model.walletName)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (model *WalletModel) TotalBalance() (amount.Amount, error) {
	return model.manager.TotalBalance(model.walletName)
}

func (model *WalletModel) TotalStake() (amount.Amount, error) {
	return model.manager.TotalStake(model.walletName)
}

func (model *WalletModel) AddressInfo(addr string) *types.AddressInfo {
	info, err := model.manager.AddressInfo(model.walletName, addr)
	if err != nil {
		return nil
	}

	return info
}

func (model *WalletModel) ListAddresses(opts ...wallet.ListAddressOption) []types.AddressInfo {
	infos, err := model.manager.ListAddresses(model.walletName, opts...)
	if err != nil {
		return nil
	}

	return infos
}

func (model *WalletModel) Balance(addr string) (amount.Amount, error) {
	return model.manager.Balance(model.walletName, addr)
}

func (model *WalletModel) Stake(addr string) (amount.Amount, error) {
	return model.manager.Stake(model.walletName, addr)
}

func (model *WalletModel) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return model.manager.PrivateKey(model.walletName, password, addr)
}

func (model *WalletModel) Mnemonic(password string) (string, error) {
	return model.manager.Mnemonic(model.walletName, password)
}

func (model *WalletModel) UpdatePassword(oldPassword, newPassword string) error {
	return model.manager.UpdatePassword(model.walletName, oldPassword, newPassword)
}

func (model *WalletModel) SetDefaultFee(fee amount.Amount) error {
	return model.manager.SetDefaultFee(model.walletName, fee)
}

func (model *WalletModel) NewAddress(
	addressType crypto.AddressType,
	label string,
	opts ...wallet.NewAddressOption,
) (*types.AddressInfo, error) {
	return model.manager.NewAddress(model.walletName, addressType, label, opts...)
}

func (model *WalletModel) AddressLabel(addr string) string {
	label, err := model.manager.AddressLabel(model.walletName, addr)
	if err != nil {
		return ""
	}

	return label
}

func (model *WalletModel) SetAddressLabel(addr, label string) error {
	return model.manager.SetAddressLabel(model.walletName, addr, label)
}

// AddressRows returns typed address rows with domain data only.
func (model *WalletModel) AddressRows() []AddressRow {
	rows := make([]AddressRow, 0)
	infos, err := model.manager.ListAddresses(model.walletName)
	if err != nil {
		return rows
	}
	for no, info := range infos {
		balance, _ := model.manager.Balance(model.walletName, info.Address)
		stake, _ := model.manager.Stake(model.walletName, info.Address)

		rows = append(rows, AddressRow{
			No:       no + 1,
			Address:  info.Address,
			Label:    info.Label,
			Path:     info.Path,
			Imported: info.Path == "",
			Balance:  balance,
			Stake:    stake,
		})
	}

	return rows
}

func (model *WalletModel) MakeTransferTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.manager.MakeTransferTx(model.walletName, sender, receiver, amt, opts...)
}

func (model *WalletModel) MakeBondTx(
	sender, receiver, publicKey string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.manager.MakeBondTx(model.walletName, sender, receiver, publicKey, amt, opts...)
}

func (model *WalletModel) MakeUnbondTx(validator string, opts ...wallet.TxOption) (*tx.Tx, error) {
	return model.manager.MakeUnbondTx(model.walletName, validator, opts...)
}

func (model *WalletModel) MakeWithdrawTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.manager.MakeWithdrawTx(model.walletName, sender, receiver, amt, opts...)
}

func (model *WalletModel) SignTransaction(password string, trx *tx.Tx) error {
	return model.manager.SignTransaction(model.walletName, password, trx)
}

func (model *WalletModel) BroadcastTransaction(trx *tx.Tx) (string, error) {
	return model.manager.BroadcastTransaction(model.walletName, trx)
}

func (model *WalletModel) Transactions(count, skip int) []*types.TransactionInfo {
	txs, err := model.manager.ListTransactions(
		model.walletName,
		wallet.WithCount(count),
		wallet.WithSkip(skip),
	)
	if err != nil {
		return nil
	}

	return txs
}
