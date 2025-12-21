//go:build gtk

package model

import (
	"errors"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/types"
)

type WalletModel struct {
	Node      *node.Node
	WalletKey string
}

// AddressRow is a UI-friendly but UI-agnostic representation of an address entry.
// Formatting (strings/markup) should be done by presenters/controllers, not here.
type AddressRow struct {
	No                int
	Address           string
	Label             string
	Path              string
	Imported          bool
	Balance           amount.Amount
	Stake             amount.Amount
	AvailabilityScore *float64
}

func NewWalletModel(n *node.Node, walletName string) (*WalletModel, error) {
	if err := n.WalletManager().LoadWallet(walletName, n.GRPC().Address()); err != nil &&
		!errors.Is(err, wltmgr.ErrWalletAlreadyLoaded) {
		return nil, err
	}

	return &WalletModel{Node: n, WalletKey: walletName}, nil
}

// WalletName returns the display name used in the UI.
func (model *WalletModel) WalletName() string {
	return model.WalletKey
}

func (model *WalletModel) WalletPath() string {
	// Prefer the wallet manager directory (it knows configured wallets dir).
	return model.Node.WalletManager().WalletPath(model.WalletKey)
}

func (model *WalletModel) IsEncrypted() bool {
	info, err := model.Node.WalletManager().WalletInfo(model.WalletKey)
	if err != nil || info == nil {
		return false
	}

	return info.Encrypted
}

func (model *WalletModel) WalletInfo() (types.WalletInfo, error) {
	info, err := model.Node.WalletManager().WalletInfo(model.WalletKey)
	if err != nil {
		return types.WalletInfo{}, err
	}
	if info == nil {
		return types.WalletInfo{}, wltmgr.ErrWalletNotLoaded
	}

	return *info, nil
}

func (model *WalletModel) TotalBalance() (amount.Amount, error) {
	return model.Node.WalletManager().TotalBalance(model.WalletKey)
}

func (model *WalletModel) TotalStake() (amount.Amount, error) {
	return model.Node.WalletManager().TotalStake(model.WalletKey)
}

func (model *WalletModel) AddressInfo(addr string) *types.AddressInfo {
	info, err := model.Node.WalletManager().AddressInfo(model.WalletKey, addr)
	if err != nil {
		return nil
	}

	return info
}

func (model *WalletModel) ListAddresses(opts ...types.ListAddressOption) []types.AddressInfo {
	infos, err := model.Node.WalletManager().ListAddresses(model.WalletKey, opts...)
	if err != nil {
		return nil
	}

	return infos
}

func (model *WalletModel) Balance(addr string) (amount.Amount, error) {
	return model.Node.WalletManager().Balance(model.WalletKey, addr)
}

func (model *WalletModel) Stake(addr string) (amount.Amount, error) {
	return model.Node.WalletManager().Stake(model.WalletKey, addr)
}

func (model *WalletModel) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return model.Node.WalletManager().PrivateKey(model.WalletKey, password, addr)
}

func (model *WalletModel) Mnemonic(password string) (string, error) {
	return model.Node.WalletManager().Mnemonic(model.WalletKey, password)
}

func (model *WalletModel) UpdatePassword(oldPassword, newPassword string) error {
	return model.Node.WalletManager().UpdatePassword(model.WalletKey, oldPassword, newPassword)
}

func (model *WalletModel) SetDefaultFee(fee amount.Amount) error {
	return model.Node.WalletManager().SetDefaultFee(model.WalletKey, fee)
}

func (model *WalletModel) NewAddress(
	addressType crypto.AddressType,
	label string,
	opts ...types.NewAddressOption,
) (*types.AddressInfo, error) {
	return model.Node.WalletManager().NewAddress(model.WalletKey, addressType, label, opts...)
}

func (model *WalletModel) AddressLabel(addr string) string {
	label, err := model.Node.WalletManager().AddressLabel(model.WalletKey, addr)
	if err != nil {
		return ""
	}

	return label
}

func (model *WalletModel) SetAddressLabel(addr, label string) error {
	return model.Node.WalletManager().SetAddressLabel(model.WalletKey, addr, label)
}

// AddressRows returns typed address rows with domain data only.
func (model *WalletModel) AddressRows() []AddressRow {
	rows := make([]AddressRow, 0)
	infos, err := model.Node.WalletManager().ListAddresses(model.WalletKey)
	if err != nil {
		return rows
	}
	for no, info := range infos {
		balance, _ := model.Node.WalletManager().Balance(model.WalletKey, info.Address)
		stake, _ := model.Node.WalletManager().Stake(model.WalletKey, info.Address)

		var scorePtr *float64
		valAddr, err := crypto.AddressFromString(info.Address)
		if err == nil {
			val := model.Node.State().ValidatorByAddress(valAddr)
			if val != nil {
				score := model.Node.State().AvailabilityScore(val.Number())
				scorePtr = &score
			}
		}

		rows = append(rows, AddressRow{
			No:                no + 1,
			Address:           info.Address,
			Label:             info.Label,
			Path:              info.Path,
			Imported:          info.Path == "",
			Balance:           balance,
			Stake:             stake,
			AvailabilityScore: scorePtr,
		})
	}

	return rows
}

func (model *WalletModel) MakeTransferTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.Node.WalletManager().MakeTransferTx(model.WalletKey, sender, receiver, amt, opts...)
}

func (model *WalletModel) MakeBondTx(
	sender, receiver, publicKey string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.Node.WalletManager().MakeBondTx(model.WalletKey, sender, receiver, publicKey, amt, opts...)
}

func (model *WalletModel) MakeUnbondTx(validator string, opts ...wallet.TxOption) (*tx.Tx, error) {
	return model.Node.WalletManager().MakeUnbondTx(model.WalletKey, validator, opts...)
}

func (model *WalletModel) MakeWithdrawTx(
	sender, receiver string,
	amt amount.Amount,
	opts ...wallet.TxOption,
) (*tx.Tx, error) {
	return model.Node.WalletManager().MakeWithdrawTx(model.WalletKey, sender, receiver, amt, opts...)
}

func (model *WalletModel) SignTransaction(password string, trx *tx.Tx) error {
	return model.Node.WalletManager().SignTransaction(model.WalletKey, password, trx)
}

func (model *WalletModel) BroadcastTransaction(trx *tx.Tx) (string, error) {
	return model.Node.WalletManager().BroadcastTransaction(model.WalletKey, trx)
}
