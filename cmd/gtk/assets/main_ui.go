//go:build gtk

package assets

import _ "embed"

// Main window / widgets UI and CSS.

//go:embed ui/main_window.ui
var MainWindowUI []byte

//go:embed ui/widget_node.ui
var NodeWidgetUI []byte

//go:embed ui/widget_wallet.ui
var WalletWidgetUI []byte

