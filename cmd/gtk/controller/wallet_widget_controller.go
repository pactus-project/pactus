//go:build gtk

package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/ezex-io/gopkg/scheduler"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type WalletWidgetModel interface {
	WalletName() string
	IsEncrypted() bool
	WalletInfo() (*types.WalletInfo, error)
	TotalBalance() (amount.Amount, error)
	TotalStake() (amount.Amount, error)
	AddressRows() []model.AddressRow
	Transactions(count, skip int) []*pactus.WalletTransactionInfo
}

type WalletWidgetHandlers struct {
	OnNewAddress     func()
	OnSetDefaultFee  func()
	OnChangePassword func()
	OnShowSeed       func()

	OnUpdateLabel    func(address string)
	OnShowDetails    func(address string)
	OnShowPrivateKey func(address string)
}

type WalletWidgetController struct {
	view *view.WalletWidgetView

	model WalletWidgetModel

	txSkip  int
	txCount int
}

func NewWalletWidgetController(view *view.WalletWidgetView, model WalletWidgetModel) *WalletWidgetController {
	return &WalletWidgetController{
		view:    view,
		model:   model,
		txCount: 20,
	}
}

func (c *WalletWidgetController) Bind(ctx context.Context, handlers WalletWidgetHandlers) {
	info, err := c.model.WalletInfo()
	if err == nil {
		c.view.LabelName.SetText(c.model.WalletName())
		c.view.LabelDriver.SetText(info.Driver)
		c.view.LabelCreatedAt.SetText(info.CreatedAt.Format(time.RFC1123))
		c.view.LabelLocation.SetText(info.Path)
	}

	// Toolbar actions via glade signals.
	c.view.ConnectSignals(map[string]any{
		"on_new_address": func() {
			handlers.OnNewAddress()
		},
		"on_address_refresh": func() {
			c.RefreshAddresses()
		},
		"on_set_default_fee": func() {
			handlers.OnSetDefaultFee()
		},
		"on_change_password": func() {
			handlers.OnChangePassword()
		},
		"on_show_seed": func() {
			handlers.OnShowSeed()
		},
		"on_tx_refresh": func() {
			c.RefreshTransactions()
		},
		"on_tx_prev": func() {
			if c.txSkip >= c.txCount {
				c.txSkip -= c.txCount
			} else {
				c.txSkip = 0
			}
			c.RefreshTransactions()
		},
		"on_tx_next": func() {
			c.txSkip += c.txCount
			c.RefreshTransactions()
		},
	})

	// Context menu actions.
	c.view.MenuItemUpdateLabel.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && handlers.OnUpdateLabel != nil {
			handlers.OnUpdateLabel(addr)
		}
	})
	c.view.MenuItemDetails.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && handlers.OnShowDetails != nil {
			handlers.OnShowDetails(addr)
		}
	})
	c.view.MenuItemPrivateKey.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && handlers.OnShowPrivateKey != nil {
			handlers.OnShowPrivateKey(addr)
		}
	})

	// Right-click popup.
	c.view.TreeViewWallet.Connect("button-press-event", func(_ *gtk.TreeView, event *gdk.Event) bool {
		eventButton := gdk.EventButtonNewFromEvent(event)
		if eventButton.Type() == gdk.EVENT_BUTTON_PRESS && eventButton.Button() == gdk.BUTTON_SECONDARY {
			c.view.ContextMenu.PopupAtPointer(event)
		}

		return false
	})

	// Double-click opens details.
	c.view.TreeViewWallet.Connect("row-activated", func(_ *gtk.TreeView, _ *gtk.TreePath, _ *gtk.TreeViewColumn) {
		addr := c.selectedAddress()
		if addr != "" && handlers.OnShowDetails != nil {
			handlers.OnShowDetails(addr)
		}
	})

	totalBalance1, _ := c.model.TotalBalance()
	scheduler.Every(ctx, 15*time.Second).Do(func() {
		totalBalance2, _ := c.model.TotalBalance()

		if totalBalance1 != totalBalance2 {
			c.Refresh()

			totalBalance1 = totalBalance2
		}
	})

	c.Refresh()
}

func (c *WalletWidgetController) selectedAddress() string {
	addr, ok, err := c.view.SelectionAddress(1)
	if err != nil || !ok {
		return ""
	}

	return addr
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

func (c *WalletWidgetController) Refresh() {
	c.RefreshInfo()
	c.RefreshAddresses()
	c.RefreshTransactions()
}

func (c *WalletWidgetController) RefreshInfo() {
	// Update info lines.
	balance, _ := c.model.TotalBalance()
	stake, _ := c.model.TotalStake()
	balanceStr := balance.String()
	stakeStr := stake.String()

	info, err := c.model.WalletInfo()
	if err != nil {
		return
	}

	glib.IdleAdd(func() bool {
		c.view.LabelEncrypted.SetText(gtkutil.YesNo(info.Encrypted))
		c.view.LabelTotalBalance.SetText(balanceStr)
		c.view.LabelTotalStake.SetText(stakeStr)
		c.view.LabelDefaultFee.SetText(info.DefaultFee.String())

		return false
	})
}

func (c *WalletWidgetController) RefreshAddresses() {
	rows := c.model.AddressRows()

	glib.IdleAdd(func() bool {
		c.view.ClearRows()
		for _, item := range rows {
			c.view.AppendRow(
				[]int{0, 1, 2, 3, 4},
				[]any{
					strconv.Itoa(item.No),
					item.Address,
					gtkutil.ImportedLabel(item.Label, item.Imported),
					item.Balance.String(),
					item.Stake.String(),
				},
			)
		}

		return false
	})
}

func (c *WalletWidgetController) RefreshTransactions() {
	trxs := c.model.Transactions(c.txCount, c.txSkip)
	hasNext := len(trxs) == c.txCount

	glib.IdleAdd(func() bool {
		c.view.ClearTxRows()

		for _, trx := range trxs {
			c.view.AppendTxRow(
				[]int{0, 1, 2, 3, 4, 5, 6, 7, 8},
				[]any{
					trx.No,
					cmd.ShortHash(trx.TxId),
					cmd.ShortAddress(trx.Sender),
					cmd.ShortAddress(trx.Receiver),
					trx.PayloadType.String(),
					amount.Amount(trx.Amount).String(),
					getDirectionTextWithIcon(types.TxDirection(trx.Direction)),
					trx.Status.String(),
					trx.Comment,
				},
			)
		}

		c.view.SetTxPager(c.txSkip > 0, hasNext)

		return false
	})
}
