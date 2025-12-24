//go:build gtk

package assets

import _ "embed"

// Dialogs/UI used by the GTK app.

// About dialogs.

//go:embed ui/dialog_about.ui
var DialogAboutUI []byte

//go:embed ui/dialog_about_gtk.ui
var DialogAboutGTKUI []byte

// Wallet dialogs.

//go:embed ui/dialog_wallet_password.ui
var WalletPasswordDialogUI []byte

//go:embed ui/dialog_wallet_create_address.ui
var WalletCreateAddressDialogUI []byte

//go:embed ui/dialog_wallet_change_password.ui
var WalletChangePasswordDialogUI []byte

//go:embed ui/dialog_wallet_set_default_fee.ui
var WalletSetDefaultFeeDialogUI []byte

//go:embed ui/dialog_wallet_show_seed.ui
var WalletShowSeedDialogUI []byte

// Address dialogs.

//go:embed ui/dialog_address_label.ui
var AddressLabelDialogUI []byte

//go:embed ui/dialog_address_details.ui
var AddressDetailsDialogUI []byte

//go:embed ui/dialog_address_private_key.ui
var AddressPrivateKeyDialogUI []byte

// Transaction dialogs.

//go:embed ui/dialog_transaction_transfer.ui
var TxTransferDialogUI []byte

//go:embed ui/dialog_transaction_bond.ui
var TxBondDialogUI []byte

//go:embed ui/dialog_transaction_unbond.ui
var TxUnbondDialogUI []byte

//go:embed ui/dialog_transaction_withdraw.ui
var TxWithdrawDialogUI []byte

