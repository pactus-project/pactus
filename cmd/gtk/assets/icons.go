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

	//go:embed icons/dialog_warning.svg
	iconDialogWarningData    []byte
	IconDialogWarningTexture *gdk.Texture

	//go:embed icons/dialog_success.svg
	iconDialogSuccessData    []byte
	IconDialogSuccessTexture *gdk.Texture

	//go:embed icons/dialog_error.svg
	iconDialogErrorData    []byte
	IconDialogErrorTexture *gdk.Texture

	//go:embed icons/dialog_info.svg
	iconDialogInfoData    []byte
	IconDialogInfoTexture *gdk.Texture

	//go:embed icons/dialog_question.svg
	iconDialogQuestionData    []byte
	IconDialogQuestionTexture *gdk.Texture

	//go:embed icons/status_warning.svg
	iconStatusWarningData    []byte
	IconStatusWarningTexture *gdk.Texture

	//go:embed icons/status_success.svg
	iconStatusSuccessData    []byte
	IconStatusSuccessTexture *gdk.Texture

	//go:embed icons/status_error.svg
	iconStatusErrorData    []byte
	IconStatusErrorTexture *gdk.Texture

	//go:embed icons/theme_sun.svg
	iconThemeSunData    []byte
	IconThemeSunTexture *gdk.Texture

	//go:embed icons/theme_moon.svg
	iconThemeMoonData    []byte
	IconThemeMoonTexture *gdk.Texture

	//go:embed icons/copy.svg
	iconCopyData    []byte
	IconCopyTexture *gdk.Texture

	//go:embed icons/download.svg
	iconDownloadData    []byte
	IconDownloadTexture *gdk.Texture

	//go:embed icons/loading.svg
	iconLoadingData    []byte
	IconLoadingTexture *gdk.Texture

	//go:embed icons/list.svg
	iconListData    []byte
	IconListTexture *gdk.Texture

	//go:embed icons/folder_open.svg
	iconFolderOpenData    []byte
	IconFolderOpenTexture *gdk.Texture

	//go:embed icons/package.svg
	iconPackageData    []byte
	IconPackageTexture *gdk.Texture

	//go:embed icons/arrow_incoming.svg
	iconArrowIncomingData    []byte
	IconArrowIncomingTexture *gdk.Texture

	//go:embed icons/arrow_outgoing.svg
	iconArrowOutgoingData    []byte
	IconArrowOutgoingTexture *gdk.Texture
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
	IconDialogWarningTexture = TextureFromBytes(iconDialogWarningData)
	IconDialogSuccessTexture = TextureFromBytes(iconDialogSuccessData)
	IconDialogErrorTexture = TextureFromBytes(iconDialogErrorData)
	IconDialogInfoTexture = TextureFromBytes(iconDialogInfoData)
	IconDialogQuestionTexture = TextureFromBytes(iconDialogQuestionData)
	IconStatusWarningTexture = TextureFromBytes(iconStatusWarningData)
	IconStatusSuccessTexture = TextureFromBytes(iconStatusSuccessData)
	IconStatusErrorTexture = TextureFromBytes(iconStatusErrorData)
	IconThemeSunTexture = TextureFromBytes(iconThemeSunData)
	IconThemeMoonTexture = TextureFromBytes(iconThemeMoonData)
	IconCopyTexture = TextureFromBytes(iconCopyData)
	IconDownloadTexture = TextureFromBytes(iconDownloadData)
	IconLoadingTexture = TextureFromBytes(iconLoadingData)
	IconListTexture = TextureFromBytes(iconListData)
	IconFolderOpenTexture = TextureFromBytes(iconFolderOpenData)
	IconPackageTexture = TextureFromBytes(iconPackageData)
	IconArrowIncomingTexture = TextureFromBytes(iconArrowIncomingData)
	IconArrowOutgoingTexture = TextureFromBytes(iconArrowOutgoingData)
}
