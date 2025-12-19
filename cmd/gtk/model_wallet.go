//go:build gtk

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/vault"
)

type walletModel struct {
	node       *node.Node
	walletName string
	listStore  *gtk.ListStore
}

func newWalletModel(nde *node.Node, walletName string) *walletModel {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // Column no
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Label
		glib.TYPE_STRING, // Balance
		glib.TYPE_STRING, // Stake
		glib.TYPE_STRING) // Availability Score

	return &walletModel{
		node:      nde,
		listStore: listStore,
	}
}

func (model *walletModel) ToTreeModel() *gtk.TreeModel {
	return model.listStore.ToTreeModel()
}

func (model *walletModel) rebuildModel() {
	go func() {
		data := [][]string{}
		for no, info := range model.wallet.AddressInfos() {
			label := info.Label
			if info.Path == "" {
				label += "(Imported)"
			}

			balance, _ := model.wallet.Balance(info.Address)
			stake, _ := model.wallet.Stake(info.Address)
			balanceStr := balance.String()
			stakeStr := stake.String()

			var score string

			valAddr, err := crypto.AddressFromString(info.Address)
			if err == nil {
				val := model.node.State().ValidatorByAddress(valAddr)
				if val != nil {
					score = strconv.FormatFloat(model.node.State().AvailabilityScore(val.Number()), 'f', -1, 64)
				}
			}

			data = append(data, []string{
				fmt.Sprintf("%v", no+1),
				info.Address,
				label,
				balanceStr,
				stakeStr,
				score,
			})
		}

		glib.IdleAdd(func() bool {
			model.listStore.Clear()
			for _, item := range data {
				iter := model.listStore.Append()
				_ = model.listStore.Set(iter,
					[]int{
						IDAddressesColumnNo,
						IDAddressesColumnAddress,
						IDAddressesColumnLabel,
						IDAddressesColumnBalance,
						IDAddressesColumnStake,
						IDAddressesColumnAvailabilityScore,
					},
					[]any{
						item[0],
						item[1],
						item[2],
						item[3],
						item[4],
						item[5],
					})
			}

			return false
		})
	}()
}

func (model *walletModel) IsEncrypted() bool {
	info, _ := model.WalletInfo()

	return info.Encrypted
}

func (model *walletModel) WalletInfo() wallet.Info {
	info, err := model.node.WalletManager().WalletInfo(model.walletName)
	if err != nil {
		log.Printf("failed to get wallet info: %v", err)

		return wallet.Info{}
	}

	return *info
}

func (model *walletModel) UpdatePassword(oldPassword, newPassword string) error {
	return model.node.WalletManager().UpdatePassword(model.walletName, oldPassword, newPassword)
}

func (model *walletModel) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return model.node.WalletManager().PrivateKey(model.walletName, password, addr)
}

func (model *walletModel) NewAddress(addressType crypto.AddressType, label string,
	opts ...wallet.NewAddressOption,
) (*vault.AddressInfo, error) {
	return model.node.WalletManager().NewAddress(model.walletName, label, opts...)
}

func (model *walletModel) SetDefaultFee(fee amount.Amount) error {
	return model.node.WalletManager().SetDefaultFee(model.walletName, fee)
}

func (model *walletModel) AllAccountAddresses() []vault.AddressInfo {
	return model.node.WalletManager().AllAccountAddresses(model.walletName)
}

func (model *walletModel) AllValidatorAddresses() []vault.AddressInfo {
	return model.node.WalletManager().AllValidatorAddresses(model.walletName)
}

func (model *walletModel) AddressInfo(addr string) *vault.AddressInfo {
	return model.node.WalletManager().AddressInfo(model.walletName, addr)
}

func (model *walletModel) Balance(addr string) amount.Amount {
	return model.node.WalletManager().Balance(model.walletName, addr)
}

func (model *walletModel) Stake(addr string) amount.Amount {
	return model.node.WalletManager().Stake(model.walletName, addr)
}

func (model *walletModel) MakeBondTx(sender, receiver, publicKey string, amount amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error) {
	return model.node.WalletManager().MakeBondTx(model.walletName, sender, receiver, publicKey, amount, opts...)
}
