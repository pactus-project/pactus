//go:build gtk

package main

import (
	"fmt"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

type walletModel struct {
	wallet    *wallet.Wallet
	listStore *gtk.ListStore
	node      *node.Node
}

func newWalletModel(wlt *wallet.Wallet, nde *node.Node) *walletModel {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // Column no
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Label
		glib.TYPE_STRING, // Balance
		glib.TYPE_STRING, // Stake
		glib.TYPE_STRING) // Availability Score

	return &walletModel{
		wallet:    wlt,
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
			balanceStr := util.ChangeToString(balance)
			stakeStr := util.ChangeToString(stake)

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
			for _, d := range data {
				iter := model.listStore.Append()
				err := model.listStore.Set(iter,
					[]int{
						IDAddressesColumnNo,
						IDAddressesColumnAddress,
						IDAddressesColumnLabel,
						IDAddressesColumnBalance,
						IDAddressesColumnStake,
						IDAddressesColumnAvailabilityScore,
					},
					[]interface{}{
						d[0],
						d[1],
						d[2],
						d[3],
						d[4],
						d[5],
					})

				errorCheck(err)
			}

			return false
		})
	}()
}
