//go:build gtk

package main

import (
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/wallet"
)

type walletModel struct {
	wallet    *wallet.Wallet
	listStore *gtk.ListStore
}

func newWalletModel(wallet *wallet.Wallet) (*walletModel, error) {
	listStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return nil, err
	}
	return &walletModel{
		wallet:    wallet,
		listStore: listStore,
	}, nil
}

func (model *walletModel) ToTreeModel() *gtk.TreeModel {
	return model.listStore.ToTreeModel()
}

func (model *walletModel) rebuildModel() error {
	model.listStore.Clear()
	for no, info := range model.wallet.AddressInfos() {
		label := info.Label
		if info.Imported {
			label += "(Imported)"
		}
		balance, _ := model.wallet.Balance(info.Address)
		stake, _ := model.wallet.Balance(info.Address)
		//errorCheck(parent, err)
		balanceStr := strconv.FormatInt(balance, 10)
		stakeStr := strconv.FormatInt(stake, 10)

		iter := model.listStore.Append()
		err := model.listStore.Set(iter,
			[]int{
				IDAddressesColumnNo,
				IDAddressesColumnAddress,
				IDAddressesColumnLabel,
				IDAddressesColumnBalance,
				IDAddressesColumnStake},
			[]interface{}{
				no + 1,
				info.Address,
				label,
				balanceStr,
				stakeStr,
			})

		if err != nil {
			return err
		}
	}

	return nil
}

func (model *walletModel) createAddress(password string) error {
	address, err := model.wallet.MakeNewAddress(password, "")
	if err != nil {
		return err
	}

	iter := model.listStore.Append()
	err = model.listStore.Set(iter,
		[]int{
			IDAddressesColumnNo,
			IDAddressesColumnAddress,
			IDAddressesColumnLabel,
			IDAddressesColumnBalance,
			IDAddressesColumnStake},
		[]interface{}{
			model.wallet.AddressCount() + 1,
			address,
			"",
			"0",
			"0",
		})
	if err != nil {
		return err
	}

	err = model.wallet.Save()
	if err != nil {
		return err
	}

	return nil
}
