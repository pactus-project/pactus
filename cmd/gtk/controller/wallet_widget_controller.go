//go:build gtk

package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/gopkg/scheduler"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// addressRow wraps the model addressRow for display purposes.
type addressRow struct {
	no   int
	addr *pactus.AddressInfo
}

type transactionRow struct {
	trx *pactus.WalletTransactionInfo
}

type WalletWidgetController struct {
	view           *view.WalletWidgetView
	model          *model.WalletModel
	lsAddresses    *gioutil.ListModel[addressRow]
	lsTransactions *gioutil.ListModel[transactionRow]
	txSkip         int
	txCount        int
}

func NewWalletWidgetController(view *view.WalletWidgetView, model *model.WalletModel) *WalletWidgetController {
	lsAddresses := gioutil.NewListModel[addressRow]()
	lsTransactions := gioutil.NewListModel[transactionRow]()

	view.ColViewAddresses.SetModel(gtk.NewSingleSelection(lsAddresses))
	view.ColViewTransactions.SetModel(gtk.NewSingleSelection(lsTransactions))

	return &WalletWidgetController{
		view:           view,
		model:          model,
		lsAddresses:    lsAddresses,
		lsTransactions: lsTransactions,
		txCount:        20,
	}
}

func (c *WalletWidgetController) BuildView(ctx context.Context, nav *Navigator) error {
	gtkutil.IdleAddSync(func() {
		gtkutil.ColumnViewSetup(c.view.ColViewAddresses, c.lsAddresses)
		gtkutil.ColumnViewSetup(c.view.ColViewTransactions, c.lsTransactions)

		// Setup address columns
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "No", func(row addressRow) string {
			return strconv.Itoa(row.no)
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "Address", func(row addressRow) string {
			return row.addr.Address
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "Type", func(row addressRow) string {
			typ := crypto.AddressType(row.addr.AddressType)
			if typ == crypto.AddressTypeTreasury {
				return ""
			}

			return typ.String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "Label", func(row addressRow) string {
			return row.addr.Label
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "Balance", func(row addressRow) string {
			return amount.Amount(row.addr.Balance).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewAddresses, "Stake", func(row addressRow) string {
			return amount.Amount(row.addr.Stake).String()
		})
		// Setup transaction columns
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "No", func(row transactionRow) string {
			return strconv.Itoa(int(row.trx.No))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "ID", func(row transactionRow) string {
			return row.trx.TxId
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Sender", func(row transactionRow) string {
			return row.trx.Sender
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Receiver", func(row transactionRow) string {
			return row.trx.Receiver
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Type", func(row transactionRow) string {
			return payload.Type(row.trx.PayloadType).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Amount", func(row transactionRow) string {
			return amount.Amount(row.trx.Amount).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Direction", func(row transactionRow) string {
			return getDirectionTextWithIcon(types.TxDirection(row.trx.Direction))
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Status", func(row transactionRow) string {
			return types.TransactionStatus(row.trx.Status).String()
		})
		gtkutil.ColumnViewAppendTextColumn(c.view.ColViewTransactions, "Comment", func(row transactionRow) string {
			return row.trx.Comment
		})
	})

	info, err := c.model.WalletInfo()

	gtkutil.IdleAddSync(func() {
		if err == nil {
			c.view.LabelName.SetText(c.model.WalletName())
			c.view.LabelDriver.SetText(info.Driver)
			c.view.LabelCreatedAt.SetText(time.Unix(info.CreatedAt, 0).Format(time.RFC1123))
			c.view.LabelLocation.SetText(info.Path)
		}

		gtkutil.ConnectButtonSignal(c.view.BtnNewAddress, nav.ShowWalletNewAddress)
		gtkutil.ConnectButtonSignal(c.view.BtnShowSeed, nav.ShowWalletShowSeed)
		gtkutil.ConnectButtonSignal(c.view.BtnChangePassword, nav.ShowWalletChangePassword)
		gtkutil.ConnectButtonSignal(c.view.BtnSetDefaultFee, nav.ShowWalletSetDefaultFee)
		gtkutil.ConnectButtonSignal(c.view.BtnRefreshAddresses, c.RefreshAddresses)

		gtkutil.ConnectButtonSignal(c.view.BtnTxRefresh, c.RefreshTransactions)
		gtkutil.ConnectButtonSignal(c.view.BtnTxNext, c.nextTransactionsPage)
		gtkutil.ConnectButtonSignal(c.view.BtnTxPrev, c.prevTransactionsPage)

		gtkutil.CaptureDoubleClick(&c.view.ColViewAddresses.Widget, c.ShowAddressDetails)
	})

	totalBalance1, _ := c.model.TotalBalance()
	scheduler.Every(refreshWalletInterval).Do(ctx, func(context.Context) {
		totalBalance2, _ := c.model.TotalBalance()

		if totalBalance1 != totalBalance2 {
			c.RefreshInfo()

			if gtkutil.IsWidgetShowing(&c.view.ColViewAddresses.Widget) {
				gtkutil.Logf("refreshing addresses")
				c.RefreshAddresses()
			}

			if gtkutil.IsWidgetShowing(&c.view.ColViewTransactions.Widget) {
				gtkutil.Logf("refreshing transactions")
				c.RefreshTransactions()
			}

			totalBalance1 = totalBalance2
		}
	})

	c.RefreshInfo()
	c.RefreshAddresses()
	c.RefreshTransactions()

	return nil
}

func (c *WalletWidgetController) SetupMenu(appWindow *gtk.ApplicationWindow) {
	gtkutil.CreateContextMenu(appWindow, &c.view.ColViewAddresses.Widget, []gtkutil.ContextMenuItem{
		{
			Label:  "Update _Label",
			Action: c.ShowUpdateLabel,
		},
		{
			Label:  "_Details",
			Action: c.ShowAddressDetails,
		},
		{
			Label:  "_Private Key",
			Action: c.ShowPrivateKey,
		},
	})
}

// getDirectionTextWithIcon returns formatted text with icon for the transaction direction.
func getDirectionTextWithIcon(dir types.TxDirection) string {
	switch dir {
	case types.TxDirectionIncoming:
		return "incoming ↘"
	case types.TxDirectionOutgoing:
		return "outgoing ↗"
	case types.TxDirectionAny:
		return "any"
	default:
		return "unknown"
	}
}

func (c *WalletWidgetController) RefreshInfo() {
	balance, _ := c.model.TotalBalance()
	stake, _ := c.model.TotalStake()
	balanceStr := balance.String()
	stakeStr := stake.String()

	info, err := c.model.WalletInfo()
	if err != nil {
		return
	}

	gtkutil.IdleAddSync(func() {
		c.view.LabelEncrypted.SetText(gtkutil.YesNo(info.Encrypted))
		c.view.LabelEncrypted.SetText(gtkutil.YesNo(info.Encrypted))
		c.view.LabelTotalBalance.SetText(balanceStr)
		c.view.LabelTotalStake.SetText(stakeStr)
		c.view.LabelDefaultFee.SetText(amount.Amount(info.DefaultFee).String())
	})
}

// RefreshAddresses updates the address list from the model.
func (c *WalletWidgetController) RefreshAddresses() {
	infos, err := c.model.ListAddresses()
	if err != nil {
		return
	}

	gtkutil.IdleAddSync(func() {
		gtkutil.ClearListModel(c.lsAddresses)

		for _, info := range infos {
			row := addressRow{
				addr: info,
			}
			c.lsAddresses.Append(row)
		}
	})
}

// RefreshTransactions updates the transaction list.
func (c *WalletWidgetController) RefreshTransactions() {
	infos, err := c.model.Transactions(c.txCount, c.txSkip)
	if err != nil {
		return
	}

	gtkutil.IdleAddAsync(func() {
		gtkutil.ClearListModel(c.lsTransactions)
		for _, info := range infos {
			row := transactionRow{
				trx: info,
			}
			c.lsTransactions.Append(row)
		}

		hasNext := len(infos) == c.txCount
		c.view.SetTxPager(c.txSkip > 0, hasNext)
	})
}

func (c *WalletWidgetController) ShowUpdateLabel() {
	dlgView := view.NewAddressLabelDialogView()
	dlgCtrl := NewAddressLabelDialogController(dlgView, c.model)
	dlgCtrl.Show(c.getSelectedAddress().Address, func() {
		go c.RefreshAddresses()
	})
}

func (c *WalletWidgetController) ShowAddressDetails() {
	dlgView := view.NewAddressDetailsDialogView()
	dlgCtrl := NewAddressDetailsDialogController(dlgView, c.model)
	dlgCtrl.Show(c.getSelectedAddress())
}

func (c *WalletWidgetController) ShowPrivateKey() {
	dlgView := view.NewAddressPrivateKeyDialogView()
	dlgCtrl := NewAddressPrivateKeyDialogController(dlgView, c.model)
	dlgCtrl.Show(c.getSelectedAddress().Address)
}

func (c *WalletWidgetController) prevTransactionsPage() {
	if c.txSkip >= c.txCount {
		c.txSkip -= c.txCount
	} else {
		c.txSkip = 0
	}
	c.RefreshTransactions()
}

func (c *WalletWidgetController) nextTransactionsPage() {
	c.txSkip += c.txCount
	c.RefreshTransactions()
}

func (c *WalletWidgetController) getSelectedAddress() *pactus.AddressInfo {
	row := gtkutil.ColumnViewGetSelectedItem(c.view.ColViewAddresses, c.lsAddresses)

	return row.addr
}
