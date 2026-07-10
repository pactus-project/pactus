//go:build gtk

package assets

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

var (
	//go:embed icons/add.svg
	iconAddData    []byte
	IconAddTexture *gdk.Texture

	//go:embed icons/ok.svg
	iconOKData    []byte
	IconOkTexture *gdk.Texture

	//go:embed icons/cancel.svg
	iconCancelData    []byte
	IconCancelTexture *gdk.Texture

	//go:embed icons/password.svg
	iconPasswordData    []byte
	IconPasswordTexture *gdk.Texture

	//go:embed icons/seed.svg
	iconSeedData    []byte
	IconSeedTexture *gdk.Texture

	//go:embed icons/close.svg
	iconCloseData    []byte
	IconCloseTexture *gdk.Texture

	//go:embed icons/send.svg
	iconSendData    []byte
	IconSendTexture *gdk.Texture

	//go:embed icons/fee.svg
	iconFeeData    []byte
	IconFeeTexture *gdk.Texture

	//go:embed icons/refresh.svg
	iconRefreshData    []byte
	IconRefreshTexture *gdk.Texture

	//go:embed icons/prev.svg
	iconPrevData    []byte
	IconPrevTexture *gdk.Texture

	//go:embed icons/next.svg
	iconNextData    []byte
	IconNextTexture *gdk.Texture

	//go:embed icons/save.svg
	iconSaveData    []byte
	IconSaveTexture *gdk.Texture

	//go:embed icons/nav_overview.svg
	iconNavOverviewData    []byte
	IconNavOverviewTexture *gdk.Texture

	//go:embed icons/nav_committee.svg
	iconNavCommitteeData    []byte
	IconNavCommitteeTexture *gdk.Texture

	//go:embed icons/nav_network.svg
	iconNavNetworkData    []byte
	IconNavNetworkTexture *gdk.Texture

	//go:embed icons/nav_validators.svg
	iconNavValidatorsData    []byte
	IconNavValidatorsTexture *gdk.Texture

	//go:embed icons/nav_wallet.svg
	iconNavWalletData    []byte
	IconNavWalletTexture *gdk.Texture
)

func initIcons() {
	IconAddTexture = TextureFromBytes(iconAddData)
	IconOkTexture = TextureFromBytes(iconOKData)
	IconCancelTexture = TextureFromBytes(iconCancelData)
	IconPasswordTexture = TextureFromBytes(iconPasswordData)
	IconSeedTexture = TextureFromBytes(iconSeedData)
	IconCloseTexture = TextureFromBytes(iconCloseData)
	IconSendTexture = TextureFromBytes(iconSendData)
	IconFeeTexture = TextureFromBytes(iconFeeData)
	IconRefreshTexture = TextureFromBytes(iconRefreshData)
	IconPrevTexture = TextureFromBytes(iconPrevData)
	IconNextTexture = TextureFromBytes(iconNextData)
	IconSaveTexture = TextureFromBytes(iconSaveData)
	IconNavOverviewTexture = TextureFromBytes(iconNavOverviewData)
	IconNavCommitteeTexture = TextureFromBytes(iconNavCommitteeData)
	IconNavNetworkTexture = TextureFromBytes(iconNavNetworkData)
	IconNavValidatorsTexture = TextureFromBytes(iconNavValidatorsData)
	IconNavWalletTexture = TextureFromBytes(iconNavWalletData)
}
