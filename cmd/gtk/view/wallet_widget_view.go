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

	TreeViewWallet       *gtk.TreeView
	TreeViewTransactions *gtk.TreeView
	LabelName            *gtk.Label
	LabelDriver          *gtk.Label
	LabelCreatedAt       *gtk.Label
	LabelLocation        *gtk.Label
	LabelEncrypted       *gtk.Label
	LabelTotalBalance    *gtk.Label
	LabelDefaultFee      *gtk.Label
	LabelTotalStake      *gtk.Label

	BtnRefreshAddresses *gtk.ToolButton
	BtnNewAddress       *gtk.ToolButton
	BtnSetDefaultFee    *gtk.ToolButton
	BtnChangePassword   *gtk.ToolButton
	BtnShowSeed         *gtk.ToolButton
	BtnTxRefresh        *gtk.ToolButton
	BtnTxPrev           *gtk.ToolButton
	BtnTxNext           *gtk.ToolButton

	ContextMenu         *gtk.Menu
	MenuItemUpdateLabel *gtk.MenuItem
	MenuItemDetails     *gtk.MenuItem
	MenuItemPrivateKey  *gtk.MenuItem

	listStore   *gtk.ListStore
	txListStore *gtk.ListStore
}

func createTextColumn(title string, columnID int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	gtkutil.FatalErrorCheck(err)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", columnID)
	gtkutil.FatalErrorCheck(err)

	column.SetResizable(true)

	return column
}

func NewWalletWidgetView() *WalletWidgetView {
	builder := NewViewBuilder(assets.WalletWidgetUI)

	treeViewWallet := builder.GetTreeViewObj("id_treeview_addresses")
	treeViewTransactions := builder.GetTreeViewObj("id_treeview_transactions")

	view := &WalletWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_wallet"),

		TreeViewWallet:       treeViewWallet,
		TreeViewTransactions: treeViewTransactions,
		LabelName:            builder.GetLabelObj("id_label_wallet_name"),
		LabelDriver:          builder.GetLabelObj("id_label_wallet_driver"),
		LabelCreatedAt:       builder.GetLabelObj("id_label_wallet_created_at"),
		LabelLocation:        builder.GetLabelObj("id_label_wallet_location"),
		LabelEncrypted:       builder.GetLabelObj("id_label_wallet_encrypted"),

		LabelTotalBalance: builder.GetLabelObj("id_label_wallet_total_balance"),
		LabelTotalStake:   builder.GetLabelObj("id_label_wallet_total_stake"),
		LabelDefaultFee:   builder.GetLabelObj("id_label_wallet_default_fee"),

		BtnRefreshAddresses: builder.GetToolButtonObj("id_button_refresh_addresses"),
		BtnNewAddress:       builder.GetToolButtonObj("id_button_new_address"),
		BtnSetDefaultFee:    builder.GetToolButtonObj("id_button_set_default_fee"),
		BtnChangePassword:   builder.GetToolButtonObj("id_button_change_password"),
		BtnShowSeed:         builder.GetToolButtonObj("id_button_show_seed"),
		BtnTxRefresh:        builder.GetToolButtonObj("id_button_tx_refresh"),
		BtnTxPrev:           builder.GetToolButtonObj("id_button_tx_prev"),
		BtnTxNext:           builder.GetToolButtonObj("id_button_tx_next"),
	}

	// Toolbar icons.
	view.BtnRefreshAddresses.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconRefreshPixbuf16))
	view.BtnNewAddress.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconAddPixbuf16))
	view.BtnSetDefaultFee.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconFeePixbuf16))
	view.BtnChangePassword.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconPasswordPixbuf16))
	view.BtnShowSeed.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconSeedPixbuf16))
	view.BtnTxRefresh.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconRefreshPixbuf16))
	view.BtnTxPrev.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconPrevPixbuf16))
	view.BtnTxNext.SetIconWidget(gtkutil.ImageFromPixbuf(assets.IconNextPixbuf16))
	view.BtnTxPrev.SetSensitive(false)
	view.BtnTxNext.SetSensitive(false)

	// Build list store for address table.
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // No
		glib.TYPE_STRING, // Address
		glib.TYPE_STRING, // Label
		glib.TYPE_STRING, // Balance
		glib.TYPE_STRING, // Stake
	)
	gtkutil.FatalErrorCheck(err)

	view.listStore = listStore
	view.TreeViewWallet.SetModel(listStore.ToTreeModel())

	// Columns.
	colNo := createTextColumn("No", 0)
	colAddress := createTextColumn("Address", 1)
	colLabel := createTextColumn("Label", 2)
	colBalance := createTextColumn("Balance", 3)
	colStake := createTextColumn("Stake", 4)

	view.TreeViewWallet.AppendColumn(colNo)
	view.TreeViewWallet.AppendColumn(colAddress)
	view.TreeViewWallet.AppendColumn(colLabel)
	view.TreeViewWallet.AppendColumn(colBalance)
	view.TreeViewWallet.AppendColumn(colStake)

	// Transactions list store and columns.
	txStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // no
		glib.TYPE_STRING, // id
		glib.TYPE_STRING, // sender
		glib.TYPE_STRING, // receiver
		glib.TYPE_STRING, // type
		glib.TYPE_STRING, // amount
		glib.TYPE_STRING, // direction
		glib.TYPE_STRING, // status
		glib.TYPE_STRING, // comment
	)
	gtkutil.FatalErrorCheck(err)

	view.txListStore = txStore
	view.TreeViewTransactions.SetModel(txStore.ToTreeModel())

	colTxNo := createTextColumn("#", 0)
	colTxID := createTextColumn("ID", 1)
	colTxSender := createTextColumn("Sender", 2)
	colTxReceiver := createTextColumn("Receiver", 3)
	colTxType := createTextColumn("Type", 4)
	colTxAmount := createTextColumn("Amount", 5)
	colTxDir := createTextColumn("Direction", 6)
	colTxStatus := createTextColumn("Status", 7)
	colTxComment := createTextColumn("Comment", 8)

	view.TreeViewTransactions.AppendColumn(colTxNo)
	view.TreeViewTransactions.AppendColumn(colTxID)
	view.TreeViewTransactions.AppendColumn(colTxSender)
	view.TreeViewTransactions.AppendColumn(colTxReceiver)
	view.TreeViewTransactions.AppendColumn(colTxType)
	view.TreeViewTransactions.AppendColumn(colTxAmount)
	view.TreeViewTransactions.AppendColumn(colTxDir)
	view.TreeViewTransactions.AppendColumn(colTxStatus)
	view.TreeViewTransactions.AppendColumn(colTxComment)

	// Context menu (actions are wired by controller).
	menu, err := gtk.MenuNew()
	gtkutil.FatalErrorCheck(err)

	view.ContextMenu = menu

	item, err := gtk.MenuItemNewWithLabel("Update _Label")
	gtkutil.FatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemUpdateLabel = item

	item, err = gtk.MenuItemNewWithLabel("_Details")
	gtkutil.FatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemDetails = item

	item, err = gtk.MenuItemNewWithLabel("_Private key")
	gtkutil.FatalErrorCheck(err)

	item.SetUseUnderline(true)
	item.Show()
	menu.Append(item)
	view.MenuItemPrivateKey = item

	return view
}

func (view *WalletWidgetView) ClearRows() {
	view.listStore.Clear()
}

func (view *WalletWidgetView) AppendRow(cols []int, values []any) {
	iter := view.listStore.Append()
	_ = view.listStore.Set(iter, cols, values)
}

func (view *WalletWidgetView) ClearTxRows() {
	view.txListStore.Clear()
}

func (view *WalletWidgetView) AppendTxRow(cols []int, values []any) {
	iter := view.txListStore.Append()
	_ = view.txListStore.Set(iter, cols, values)
}

func (view *WalletWidgetView) SetTxPager(prevEnabled, nextEnabled bool) {
	view.BtnTxPrev.SetSensitive(prevEnabled)
	view.BtnTxNext.SetSensitive(nextEnabled)
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
