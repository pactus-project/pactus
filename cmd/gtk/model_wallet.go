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
	info, err := model.node.WalletManager().WalletInfo(model.walletName)
	if err != nil {
		log.Println("failed to get wallet info: %s", err.Error())

		return false
	}

	return info.Encrypted
}

func (model *walletModel) UpdatePassword(oldPassword, newPassword string) error {
	model.node.WalletManager().UpdatePassword(model.walletName, oldPassword, newPassword)
	if err != nil {
		log.Println("failed to load wallet: %s", err.Error())

		return err
	}

	return model.node.WalletManager().UpdatePassword(model.walletName, oldPassword, newPassword)
}
