//go:build gtk

package main

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

type walletModel struct {
	wallet    *wallet.Wallet
	listStore *gtk.ListStore
}

func newWalletModel(wallet *wallet.Wallet) *walletModel {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // Column no
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Label
		glib.TYPE_STRING, // balance
		glib.TYPE_STRING) // Stake

	return &walletModel{
		wallet:    wallet,
		listStore: listStore,
	}
}

func (model *walletModel) ToTreeModel() *gtk.TreeModel {
	return model.listStore.ToTreeModel()
}

func (model *walletModel) rebuildModel() {
	go func() {
		data := [][]string{}
		for no, info := range model.wallet.AddressLabels() {
			label := info.Label
			if info.Imported {
				label += "(Imported)"
			}

			balanceStr := "-"
			stakeStr := "-"
			if no < 64 {
				balance, _ := model.wallet.Balance(info.Address)
				stake, _ := model.wallet.Stake(info.Address)
				balanceStr = util.ChangeToString(balance)
				stakeStr = util.ChangeToString(stake)
			}

			data = append(data, []string{
				fmt.Sprintf("%v", no+1),
				info.Address,
				label,
				balanceStr,
				stakeStr,
			})
		}

		glib.IdleAdd(func() bool {
			model.listStore.Clear()
			for _, d := range data {
				iter := model.listStore.Append()
				err := model.listStore.Set(iter,
					[]int{
						IDAddressesColumnNo,
						IDAddressesColumnAddress,
						IDAddressesColumnLabel,
						IDAddressesColumnBalance,
						IDAddressesColumnStake},
					[]interface{}{
						d[0],
						d[1],
						d[2],
						d[3],
						d[4],
					})

				errorCheck(err)
			}

			return false
		})
	}()
}

func (model *walletModel) createAddress() error {
	address, err := model.wallet.DeriveNewAddress("")
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
			fmt.Sprintf("%v", model.wallet.AddressCount()+1),
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
