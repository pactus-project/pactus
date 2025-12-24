//go:build gtk

package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet/types"
)

type WalletWidgetModel interface {
	WalletName() string
	IsEncrypted() bool
	WalletInfo() (*types.WalletInfo, error)
	TotalBalance() (amount.Amount, error)
	TotalStake() (amount.Amount, error)
	AddressRows() []model.AddressRow
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

	handlers WalletWidgetHandlers

	timeoutID glib.SourceHandle

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWalletWidgetController(view *view.WalletWidgetView, model WalletWidgetModel) *WalletWidgetController {
	return &WalletWidgetController{view: view, model: model}
}

func (c *WalletWidgetController) Bind(h WalletWidgetHandlers) {
	c.handlers = h

	// Reset lifecycle context (in case Bind is called more than once).
	if c.cancel != nil {
		c.cancel()
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

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
			if c.handlers.OnNewAddress != nil {
				c.handlers.OnNewAddress()
			}
		},
		"on_set_default_fee": func() {
			if c.handlers.OnSetDefaultFee != nil {
				c.handlers.OnSetDefaultFee()
			}
		},
		"on_change_password": func() {
			if c.handlers.OnChangePassword != nil {
				c.handlers.OnChangePassword()
			}
		},
		"on_show_seed": func() {
			if c.handlers.OnShowSeed != nil {
				c.handlers.OnShowSeed()
			}
		},
	})

	// Context menu actions.
	c.view.MenuItemUpdateLabel.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && c.handlers.OnUpdateLabel != nil {
			c.handlers.OnUpdateLabel(addr)
		}
	})
	c.view.MenuItemDetails.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && c.handlers.OnShowDetails != nil {
			c.handlers.OnShowDetails(addr)
		}
	})
	c.view.MenuItemPrivateKey.Connect("activate", func(_ *gtk.MenuItem) {
		addr := c.selectedAddress()
		if addr != "" && c.handlers.OnShowPrivateKey != nil {
			c.handlers.OnShowPrivateKey(addr)
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
		if addr != "" && c.handlers.OnShowDetails != nil {
			c.handlers.OnShowDetails(addr)
		}
	})

	c.Refresh()
	c.timeoutID = glib.TimeoutAdd(15000, func() bool {
		c.Refresh()

		return true
	})
}

func (c *WalletWidgetController) selectedAddress() string {
	addr, ok, err := c.view.SelectionAddress(1)
	if err != nil || !ok {
		return ""
	}

	return addr
}

func (c *WalletWidgetController) Refresh() {
	// Compute in background and then update UI on main loop.
	ctx := c.ctx
	go func() {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		rows := c.model.AddressRows()

		// Update info lines.

		balance, _ := c.model.TotalBalance()
		stake, _ := c.model.TotalStake()
		balanceStr := balance.String()
		stakeStr := stake.String()
		feeStr := ""
		encryptedStr := ""
		if info, err := c.model.WalletInfo(); err == nil {
			feeStr = info.DefaultFee.String()
			encryptedStr = gtkutil.YesNo(info.Encrypted)
		}

		if ctx != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		glib.IdleAdd(func() bool {
			if ctx != nil {
				select {
				case <-ctx.Done():
					return false
				default:
				}
			}

			c.view.LabelEncrypted.SetText(encryptedStr)
			c.view.LabelTotalBalance.SetText(balanceStr)
			c.view.LabelTotalStake.SetText(stakeStr)
			c.view.LabelDefaultFee.SetText(feeStr)

			c.view.ClearRows()
			for _, item := range rows {
				c.view.AppendRow(
					[]int{0, 1, 2, 3, 4, 5},
					[]any{
						strconv.Itoa(item.No),
						item.Address,
						gtkutil.ImportedLabel(item.Label, item.Imported),
						item.Balance.String(),
						item.Stake.String(),
						gtkutil.FloatPtr(item.AvailabilityScore),
					},
				)
			}

			return false
		})
	}()
}

func (c *WalletWidgetController) Cleanup() {
	if c.timeoutID != 0 {
		glib.SourceRemove(c.timeoutID)
		c.timeoutID = 0
	}
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
}
