//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	TreeViewWallet    *gtk.TreeView
	LabelName         *gtk.Label
	LabelLocation     *gtk.Label
	LabelEncrypted    *gtk.Label
	LabelTotalBalance *gtk.Label
	LabelDefaultFee   *gtk.Label
	LabelTotalStake   *gtk.Label

	BtnNewAddress     *gtk.ToolButton
	BtnSetDefaultFee  *gtk.ToolButton
	BtnChangePassword *gtk.ToolButton
	BtnShowSeed       *gtk.ToolButton

	ContextMenu         *gtk.Menu
	MenuItemUpdateLabel *gtk.MenuItem
	MenuItemDetails     *gtk.MenuItem
	MenuItemPrivateKey  *gtk.MenuItem

	listStore *gtk.ListStore
}

func createTextColumn(title string, columnID int) (*gtk.TreeViewColumn, error) {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", columnID)
	if err != nil {
		return nil, err
	}

	column.SetResizable(true)

	return column, nil
}

func NewWalletWidgetView(columnTypes ...glib.Type) (*WalletWidgetView, error) {
	builder := NewViewBuilder(assets.WalletWidgetUI)

	box := builder.GetBoxObj("id_box_wallet")
	treeViewWallet := builder.GetTreeViewObj("id_treeview_addresses")

	view := &WalletWidgetView{
		Box: box,

		TreeViewWallet: treeViewWallet,
		LabelName:      builder.GetLabelObj("id_label_wallet_name"),
		LabelLocation:  builder.GetLabelObj("id_label_wallet_location"),
		LabelEncrypted: builder.GetLabelObj("id_label_wallet_encrypted"),

		LabelTotalBalance: builder.GetLabelObj("id_label_wallet_total_balance"),
		LabelTotalStake:   builder.GetLabelObj("id_label_wallet_total_stake"),
		LabelDefaultFee:   builder.GetLabelObj("id_label_wallet_default_fee"),

		BtnNewAddress:     builder.GetToolButtonObj("id_button_new_address"),
		BtnSetDefaultFee:  builder.GetToolButtonObj("id_button_set_default_fee"),
		BtnChangePassword: builder.GetToolButtonObj("id_button_change_password"),
		BtnShowSeed:       builder.GetToolButtonObj("id_button_show_seed"),

		ViewBuilder: builder,
	}

	// Toolbar icons.
	view.BtnNewAddress.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconAddPixbuf16))
	view.BtnSetDefaultFee.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconFeePixbuf16))
	view.BtnChangePassword.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconPasswordPixbuf16))
	view.BtnShowSeed.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconSeedPixbuf16))

	// Build list store for address table.
	if len(columnTypes) == 0 {
		columnTypes = []glib.Type{
			glib.TYPE_STRING,
			glib.TYPE_STRING,
			glib.TYPE_STRING,
			glib.TYPE_STRING,
			glib.TYPE_STRING,
			glib.TYPE_STRING,
		}
	}
	ls, err := gtk.ListStoreNew(columnTypes...)
	if err != nil {
		return nil, err
	}
	view.listStore = ls
	view.TreeViewWallet.SetModel(ls.ToTreeModel())

	// Columns.
	colNo, err := createTextColumn("No", 0)
	if err != nil {
		return nil, err
	}
	colAddress, err := createTextColumn("Address", 1)
	if err != nil {
		return nil, err
	}
	colLabel, err := createTextColumn("Label", 2)
	if err != nil {
		return nil, err
	}
	colBalance, err := createTextColumn("Balance", 3)
	if err != nil {
		return nil, err
	}
	colStake, err := createTextColumn("Stake", 4)
	if err != nil {
		return nil, err
	}
	colScore, err := createTextColumn("Availability Score", 5)
	if err != nil {
		return nil, err
	}

	view.TreeViewWallet.AppendColumn(colNo)
	view.TreeViewWallet.AppendColumn(colAddress)
	view.TreeViewWallet.AppendColumn(colLabel)
	view.TreeViewWallet.AppendColumn(colBalance)
	view.TreeViewWallet.AppendColumn(colStake)
	view.TreeViewWallet.AppendColumn(colScore)

	// Context menu (actions are wired by controller).
	menu, err := gtk.MenuNew()
	if err != nil {
		return nil, err
	}
	view.ContextMenu = menu

	item, err := gtk.MenuItemNewWithLabel("Update _Label")
	if err != nil {
		return nil, err
	}
	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemUpdateLabel = item

	item, err = gtk.MenuItemNewWithLabel("_Details")
	if err != nil {
		return nil, err
	}
	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemDetails = item

	item, err = gtk.MenuItemNewWithLabel("_Private key")
	if err != nil {
		return nil, err
	}
	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemPrivateKey = item

	return view, nil
}

func (view *WalletWidgetView) ClearRows() {
	view.listStore.Clear()
}

func (view *WalletWidgetView) AppendRow(cols []int, values []any) {
	iter := view.listStore.Append()
	_ = view.listStore.Set(iter, cols, values)
}

func (view *WalletWidgetView) SelectionAddress(addressColumn int) (string, bool, error) {
	selection, err := view.TreeViewWallet.GetSelection()
	if err != nil {
		return "", false, err
	}
	if selection == nil {
		return "", false, nil
	}
	model, iter, ok := selection.GetSelected()
	if !ok {
		return "", false, nil
	}
	val, err := model.(*gtk.TreeModel).GetValue(iter, addressColumn)
	if err != nil {
		return "", false, err
	}
	addr, err := val.GetString()
	if err != nil {
		return "", false, err
	}

	return addr, true, nil
}
