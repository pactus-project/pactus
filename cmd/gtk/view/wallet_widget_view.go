//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletWidgetView struct {
	ViewBuilder

	Box             *gtk.Box
	BoxInfo         *gtk.Box
	BoxAddresses    *gtk.Box
	BoxTransactions *gtk.Box

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
}

func NewWalletWidgetView() *WalletWidgetView {
	builder := NewViewBuilder(assets.WalletWidgetUI)

	view := &WalletWidgetView{
		ViewBuilder:     builder,
		Box:             builder.GetBoxObj("id_box_wallet"),
		BoxInfo:         builder.GetBoxObj("id_box_wallet_info"),
		BoxAddresses:    builder.GetBoxObj("id_box_wallet_addresses"),
		BoxTransactions: builder.GetBoxObj("id_box_wallet_transactions"),

		ColViewAddresses:    builder.GetColumnViewObj("id_columnview_addresses"),
		ColViewTransactions: builder.GetColumnViewObj("id_columnview_transactions"),
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

	gtkutil.ColumnViewSetDefaultProperties(view.ColViewAddresses)
	gtkutil.ColumnViewSetDefaultProperties(view.ColViewTransactions)

	gtkutil.ExtendImageButton(view.BtnRefreshAddresses, "_Refresh",
		"Refresh wallet addresses", assets.IconRefreshTexture)
	gtkutil.ExtendImageButton(view.BtnNewAddress, "_New Address",
		"Create a new address", assets.IconAddTexture)
	gtkutil.ExtendImageButton(view.BtnSetDefaultFee, "Default _Fee",
		"Set the default transaction fee", assets.IconFeeTexture)
	gtkutil.ExtendImageButton(view.BtnChangePassword, "Change _Password",
		"Change the wallet password", assets.IconPasswordTexture)
	gtkutil.ExtendImageButton(view.BtnShowSeed, "Show _Seed",
		"Display the wallet seed phrase", assets.IconSeedTexture)
	gtkutil.ExtendImageButton(view.BtnTxRefresh, "_Refresh",
		"Refresh transaction list", assets.IconRefreshTexture)
	gtkutil.ExtendImageButton(view.BtnTxPrev, "_Previous",
		"Show previous transactions", assets.IconPrevTexture)
	gtkutil.ExtendImageButton(view.BtnTxNext, "_Next",
		"Show next transactions", assets.IconNextTexture)
	view.BtnTxPrev.SetSensitive(false)
	view.BtnTxNext.SetSensitive(false)

	return view
}

func (view *WalletWidgetView) SetTxPager(prevEnabled, nextEnabled bool) {
	view.BtnTxPrev.SetSensitive(prevEnabled)
	view.BtnTxNext.SetSensitive(nextEnabled)
}
