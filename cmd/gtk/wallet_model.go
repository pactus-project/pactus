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

func newWalletModel(wallet *wallet.Wallet) *walletModel {
	listStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	errorCheck(err)
	return &walletModel{
		wallet:    wallet,
		listStore: listStore,
	}
}

func (wm *walletModel) ToTreeModel() *gtk.TreeModel {
	return wm.listStore.ToTreeModel()
}

func (wm *walletModel) rebuildModel() {
	addrs := wm.wallet.Addresses()

	wm.listStore.Clear()
	for no, addr := range addrs {
		label := addr.Label
		if addr.Imported {
			label += "(Imported)"
		}
		balance, stake, _ := wm.wallet.GetBalance(addr.Address)
		//errorCheck(err)
		balanceStr := strconv.FormatInt(balance, 10)
		stakeStr := strconv.FormatInt(stake, 10)

		iter := wm.listStore.Append()
		err := wm.listStore.Set(iter,
			[]int{
				ADDRESSES_COLUMN_NO,
				ADDRESSES_COLUMN_ADDRESS,
				ADDRESSES_COLUMN_LABEL,
				ADDRESSES_COLUMN_BALANCE,
				ADDRESSES_COLUMN_STAKE},
			[]interface{}{
				no + 1,
				addr.Address,
				label,
				balanceStr,
				stakeStr,
			})

		errorCheck(err)
	}
}
