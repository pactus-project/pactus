//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type WalletWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	ColViewAddresses    *gtk.ColumnView
	ColViewTransactions *gtk.ColumnView
	LabelName           *gtk.Label
	LabelDriver         *gtk.Label
	LabelCreatedAt      *gtk.Label
	LabelLocation       *gtk.Label
	LabelEncrypted      *gtk.Label
	LabelTotalBalance   *gtk.Label
	LabelDefaultFee     *gtk.Label
	LabelTotalStake     *gtk.Label

	BtnRefreshAddresses *gtk.Button
	BtnNewAddress       *gtk.Button
	BtnSetDefaultFee    *gtk.Button
	BtnChangePassword   *gtk.Button
	BtnShowSeed         *gtk.Button
	BtnTxRefresh        *gtk.Button
	BtnTxPrev           *gtk.Button
	BtnTxNext           *gtk.Button

	ContextMenu         *gtk.PopoverMenu
	MenuItemUpdateLabel *gtk.Button
	MenuItemDetails     *gtk.Button
	MenuItemPrivateKey  *gtk.Button

	listStore   *gio.ListStore
	txListStore *gio.ListStore
}

// func createTextColumn(title string, columnID int) *gtk.ListViewColumn {
// 	cellRenderer := gtk.NewCellRendererText()

// 	column := gtk.NewListView(title, cellRenderer, "text", columnID)

// 	column.SetResizable(true)

// 	return column
// }

func NewWalletWidgetView() *WalletWidgetView {
	builder := NewViewBuilder(assets.WalletWidgetUI)

	listViewAddresses := builder.GetColumnViewObj("id_columnview_addresses")
	listViewTransactions := builder.GetColumnViewObj("id_columnview_transactions")

	view := &WalletWidgetView{
		ViewBuilder: builder,
		Box:         builder.GetBoxObj("id_box_wallet"),

		ColViewAddresses:    listViewAddresses,
		ColViewTransactions: listViewTransactions,
		LabelName:           builder.GetLabelObj("id_label_wallet_name"),
		LabelDriver:         builder.GetLabelObj("id_label_wallet_driver"),
		LabelCreatedAt:      builder.GetLabelObj("id_label_wallet_created_at"),
		LabelLocation:       builder.GetLabelObj("id_label_wallet_location"),
		LabelEncrypted:      builder.GetLabelObj("id_label_wallet_encrypted"),

		LabelTotalBalance: builder.GetLabelObj("id_label_wallet_total_balance"),
		LabelTotalStake:   builder.GetLabelObj("id_label_wallet_total_stake"),
		LabelDefaultFee:   builder.GetLabelObj("id_label_wallet_default_fee"),

		BtnRefreshAddresses: builder.GetButtonObj("id_button_refresh_addresses"),
		BtnNewAddress:       builder.GetButtonObj("id_button_new_address"),
		BtnSetDefaultFee:    builder.GetButtonObj("id_button_set_default_fee"),
		BtnChangePassword:   builder.GetButtonObj("id_button_change_password"),
		BtnShowSeed:         builder.GetButtonObj("id_button_show_seed"),
		BtnTxRefresh:        builder.GetButtonObj("id_button_tx_refresh"),
		BtnTxPrev:           builder.GetButtonObj("id_button_tx_prev"),
		BtnTxNext:           builder.GetButtonObj("id_button_tx_next"),
	}

	// Toolbar icons for GTK4 - set icon on buttons
	view.BtnRefreshAddresses.SetChild(assets.IconRefresh16)
	view.BtnNewAddress.SetChild(assets.IconAdd16)
	view.BtnSetDefaultFee.SetChild(assets.IconFee16)
	view.BtnChangePassword.SetChild(assets.IconPassword16)
	view.BtnShowSeed.SetChild(assets.IconSeed16)
	view.BtnTxRefresh.SetChild(assets.IconRefresh16)
	view.BtnTxPrev.SetChild(assets.IconPrev16)
	view.BtnTxNext.SetChild(assets.IconNext16)
	view.BtnTxPrev.SetSensitive(false)
	view.BtnTxNext.SetSensitive(false)

	// // Build list store for address table.
	// listStore := gio.NewListStore(glib.TypeString)
	// listStore.Append()
	// // 	glib.TypeString, // No
	// // 	glib.TypeString, // Address
	// // 	glib.TypeString, // Label
	// // 	glib.TypeString, // Balance
	// // 	glib.TypeString, // Stake
	// // })

	// view.listStore = listStore
	// view.ListViewAddresses.SetModel(gtk.NewSingleSelection(listStore))

	// // Columns.
	// colNo := createTextColumn("No", 0)
	// colAddress := createTextColumn("Address", 1)
	// colLabel := createTextColumn("Label", 2)
	// colBalance := createTextColumn("Balance", 3)
	// colStake := createTextColumn("Stake", 4)

	// view.ListViewAddresses.AppendColumn(colNo)
	// view.ListViewAddresses.AppendColumn(colAddress)
	// view.ListViewAddresses.AppendColumn(colLabel)
	// view.ListViewAddresses.AppendColumn(colBalance)
	// view.ListViewAddresses.AppendColumn(colStake)

	// // Transactions list store and columns.
	// txStore := gtk.NewListStore([]glib.Type{
	// 	glib.TypeString, // no
	// 	glib.TypeString, // id
	// 	glib.TypeString, // sender
	// 	glib.TypeString, // receiver
	// 	glib.TypeString, // type
	// 	glib.TypeString, // amount
	// 	glib.TypeString, // direction
	// 	glib.TypeString, // status
	// 	glib.TypeString, // comment
	// })

	// view.txListStore = txStore
	// view.ListViewTransactions.SetModel(&txStore.TreeModel)

	// colTxNo := createTextColumn("#", 0)
	// colTxID := createTextColumn("ID", 1)
	// colTxSender := createTextColumn("Sender", 2)
	// colTxReceiver := createTextColumn("Receiver", 3)
	// colTxType := createTextColumn("Type", 4)
	// colTxAmount := createTextColumn("Amount", 5)
	// colTxDir := createTextColumn("Direction", 6)
	// colTxStatus := createTextColumn("Status", 7)
	// colTxComment := createTextColumn("Comment", 8)

	// view.ListViewTransactions.AppendColumn(colTxNo)
	// view.ListViewTransactions.AppendColumn(colTxID)
	// view.ListViewTransactions.AppendColumn(colTxSender)
	// view.ListViewTransactions.AppendColumn(colTxReceiver)
	// view.ListViewTransactions.AppendColumn(colTxType)
	// view.ListViewTransactions.AppendColumn(colTxAmount)
	// view.ListViewTransactions.AppendColumn(colTxDir)
	// view.ListViewTransactions.AppendColumn(colTxStatus)
	// view.ListViewTransactions.AppendColumn(colTxComment)

	// popover := gtk.NewPopoverMenu()
	// view.ContextMenu = popover

	// // Create menu buttons for context menu actions
	// view.MenuItemUpdateLabel = gtk.NewButtonWithLabel("Update Label")
	// view.MenuItemDetails = gtk.NewButtonWithLabel("Details")
	// view.MenuItemPrivateKey = gtk.NewButtonWithLabel("Private Key")

	return view
}

// func (view *WalletWidgetView) ClearRows() {
// 	view.listStore.Clear()
// }

// func (view *WalletWidgetView) AppendRow(cols []int, values []any) {
// 	gtkutil.AppendRowToListStore(view.listStore, cols, values)
// }

// func (view *WalletWidgetView) ClearTxRows() {
// 	view.txListStore.Clear()
// }

// func (view *WalletWidgetView) AppendTxRow(cols []int, values []any) {
// 	gtkutil.AppendRowToListStore(view.txListStore, cols, values)
// }

func (view *WalletWidgetView) SetTxPager(prevEnabled, nextEnabled bool) {
	view.BtnTxPrev.SetSensitive(prevEnabled)
	view.BtnTxNext.SetSensitive(nextEnabled)
}

// func (view *WalletWidgetView) SelectionAddress(addressColumn int) (string, bool, error) {
// 	selection, err := view.ListViewAddresses.GetSelection()
// 	if err != nil {
// 		return "", false, err
// 	}
// 	if selection == nil {
// 		return "", false, nil
// 	}
// 	model, iter, ok := selection.GetSelected()
// 	if !ok {
// 		return "", false, nil
// 	}
// 	val, err := model.(*gtk.TreeModel).GetValue(iter, addressColumn)
// 	if err != nil {
// 		return "", false, err
// 	}
// 	addr, err := val.GetString()
// 	if err != nil {
// 		return "", false, err
// 	}

// 	return addr, true, nil
// }
