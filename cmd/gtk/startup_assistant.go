//go:build gtk

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/wallet"
)

func getTextViewcontent(tv *gtk.TextView) string {
	buf, _ := tv.GetBuffer()
	startIter, endIter := buf.GetBounds()
	content, _ := buf.GetText(startIter, endIter, true)
	return content
}

func setTextViewcontent(tv *gtk.TextView, content string) {
	buf, err := tv.GetBuffer()
	errorCheck(err)
	buf.SetText(content)
}

func setMargin(widget gtk.IWidget, top, bottom, start, end int) {
	widget.ToWidget().SetMarginTop(top)
	widget.ToWidget().SetMarginBottom(bottom)
	widget.ToWidget().SetMarginStart(start)
	widget.ToWidget().SetMarginEnd(end)
}

func startupAssistant(workingDir string) bool {
	gtk.Init(nil)

	successful := false
	createPage := func(assistant *gtk.Assistant, content gtk.IWidget, name, title, subject, desc string) *gtk.Widget {
		page, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
		errorCheck(err)
		page.SetHExpand(true)
		titleLabel, err := gtk.LabelNew(title)
		errorCheck(err)
		setMargin(titleLabel, 0, 20, 0, 0)
		frame, err := gtk.FrameNew(subject)
		errorCheck(err)
		frame.SetHExpand(true)

		descLabel, err := gtk.LabelNew("")
		errorCheck(err)
		descLabel.SetUseMarkup(true)
		descLabel.SetMarkup(desc)
		descLabel.SetVExpand(true)
		descLabel.SetVAlign(gtk.ALIGN_END)
		descLabel.SetHAlign(gtk.ALIGN_START)
		setMargin(descLabel, 0, 0, 0, 0)
		frame.Add(content)

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
		errorCheck(err)
		box.Add(frame)
		box.Add(descLabel)

		page.Add(titleLabel)
		page.Add(box)

		page.SetName(name)
		assistant.AppendPage(page)
		assistant.SetPageTitle(page, title)

		return page.ToWidget()
	}

	assistant, err := gtk.AssistantNew()
	errorCheck(err)

	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Zarb - Init Wizard")

	var pageMode *gtk.Widget
	var pagePassword *gtk.Widget
	var pageSeed *gtk.Widget
	var pageSeedConfirm *gtk.Widget
	var pageFinal *gtk.Widget

	// --- PageMode
	newWalletRadio, err := gtk.RadioButtonNewWithLabel(nil, "Create a new wallet from the scratch")
	errorCheck(err)
	recoverWalletRadio, err := gtk.RadioButtonNewWithLabelFromWidget(newWalletRadio, "Restore a wallet from the seed phrase")
	errorCheck(err)
	recoverWalletRadio.SetSensitive(false)

	radioBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	errorCheck(err)
	radioBox.Add(newWalletRadio)
	setMargin(newWalletRadio, 6, 6, 6, 6)
	radioBox.Add(recoverWalletRadio)
	setMargin(recoverWalletRadio, 6, 6, 6, 6)

	pageModeName := "page_mode"
	pageModeTitle := "Initialize mode"
	pageModeSubject := "How to create your wallet?"
	pageModeDesc := "If you are running the node for the first time, choose the first option."
	pageMode = createPage(
		assistant,
		radioBox,
		pageModeName,
		pageModeTitle,
		pageModeSubject,
		pageModeDesc)

	// --- pageSeed
	seedTextView, err := gtk.TextViewNew()
	errorCheck(err)
	setMargin(seedTextView, 6, 6, 6, 6)
	seedTextView.SetWrapMode(gtk.WRAP_WORD)
	seedTextView.SetEditable(false)
	seedTextView.SetMonospace(true)
	seedTextView.SetSizeRequest(0, 80)

	pageSeedName := "page_seed"
	pageSeedTitle := "Wallet seed"
	pageSeedSubject := "Your wallet generation seed is:"
	pageSeedDesc := `<span allow_breaks="true">Please write these 12 words on paper.
This seed will allow you to recover your wallet in case of computer failure.
<b>WARNING:</b>
  - Never disclose your seed.
  - Never type it on a website.
  - Do not store it electronically.</span>`

	pageSeed = createPage(
		assistant,
		seedTextView,
		pageSeedName,
		pageSeedTitle,
		pageSeedSubject,
		pageSeedDesc)

	// --- pageSeedConfirm
	seedConfirmTextView, err := gtk.TextViewNew()
	errorCheck(err)
	setMargin(seedConfirmTextView, 6, 6, 6, 6)
	seedConfirmTextView.SetWrapMode(gtk.WRAP_WORD)
	seedConfirmTextView.SetEditable(true)
	seedConfirmTextView.SetMonospace(true)
	seedConfirmTextView.SetSizeRequest(0, 80)

	seedConfirmTextView.Connect("paste_clipboard", func(textView *gtk.TextView) {
		showInfoDialog("Opps, no copy paste!")
		seedConfirmTextView.StopEmission("paste_clipboard")
	})

	seedConfirmTextBuffer, err := seedConfirmTextView.GetBuffer()
	errorCheck(err)
	seedConfirmTextBuffer.Connect("changed", func(buf *gtk.TextBuffer) {
		mnemonic1 := getTextViewcontent(seedTextView)
		mnemonic2 := getTextViewcontent(seedConfirmTextView)
		space := regexp.MustCompile(`\s+`)
		mnemonic2 = space.ReplaceAllString(mnemonic2, " ")
		mnemonic2 = strings.TrimSpace(mnemonic2)
		if mnemonic1 == mnemonic2 {
			assistant.SetPageComplete(pageSeedConfirm, true)
		} else {
			assistant.SetPageComplete(pageSeedConfirm, false)
		}
	})

	pageSeedConfirmName := "page_seed_confirm"
	pageSeedConfirmTitle := "Confirm seed"
	pageSeedConfirmSubject := "What was your seed?"
	pageSeedConfirmDesc := `Your seed is important!
To make sure that you have properly saved your seed, please retype it here.`

	pageSeedConfirm = createPage(
		assistant,
		seedConfirmTextView,
		pageSeedConfirmName,
		pageSeedConfirmTitle,
		pageSeedConfirmSubject,
		pageSeedConfirmDesc)

	// --- PagePassword
	passwordEntry, err := gtk.EntryNew()
	errorCheck(err)
	setMargin(passwordEntry, 6, 6, 6, 6)
	passwordEntry.SetVisibility(false)
	passwordLabel, err := gtk.LabelNew("Password: ")
	errorCheck(err)
	passwordLabel.SetHAlign(gtk.ALIGN_START)
	setMargin(passwordLabel, 6, 6, 6, 6)

	passwordConfirmEntry, err := gtk.EntryNew()
	errorCheck(err)
	setMargin(passwordConfirmEntry, 6, 6, 6, 6)
	passwordConfirmEntry.SetVisibility(false)
	confirmationLineLabel, err := gtk.LabelNew("Confirmation: ")
	errorCheck(err)
	confirmationLineLabel.SetHAlign(gtk.ALIGN_START)
	setMargin(confirmationLineLabel, 6, 6, 6, 6)

	grid, err := gtk.GridNew()
	errorCheck(err)
	grid.Add(passwordLabel)
	grid.Attach(passwordEntry, 1, 0, 1, 1)
	grid.AttachNextTo(confirmationLineLabel, passwordLabel, gtk.POS_BOTTOM, 1, 1)
	grid.AttachNextTo(passwordConfirmEntry, passwordEntry, gtk.POS_BOTTOM, 1, 1)

	validatePassword := func() {
		pass1, err := passwordEntry.GetText()
		errorCheck(err)
		pass2, err := passwordConfirmEntry.GetText()
		errorCheck(err)
		if pass1 == pass2 {
			assistant.SetPageComplete(pagePassword, true)
		} else {
			assistant.SetPageComplete(pagePassword, false)
		}
	}
	passwordEntry.Connect("changed", func(entry *gtk.Entry) {
		validatePassword()
	})

	passwordConfirmEntry.Connect("changed", func(entry *gtk.Entry) {
		validatePassword()
	})

	pagePasswordName := "page_password"
	pagePasswordTitle := "Wallet password"
	pagePasswordSubject := "Enter password for your wallet:"
	pagePsswrdDesc := "Please choose a strong password for your wallet."

	pagePassword = createPage(
		assistant,
		grid,
		pagePasswordName,
		pagePasswordTitle,
		pagePasswordSubject,
		pagePsswrdDesc)

	// --- pageFinal
	NodeInfoTextView, err := gtk.TextViewNew()
	errorCheck(err)
	setMargin(NodeInfoTextView, 6, 6, 6, 6)
	NodeInfoTextView.SetWrapMode(gtk.WRAP_WORD)
	NodeInfoTextView.SetEditable(false)
	NodeInfoTextView.SetMonospace(true)
	NodeInfoTextView.SetSizeRequest(0, 160)

	pageFinalName := "page_final"
	pageFinalTitle := "Run the node"
	pageFinalSubject := "Your node information:"
	pageFinalDesc := `Congratulation. Your node is initialized successfully.
Now you are ready to start the node!`

	pageFinal = createPage(
		assistant,
		NodeInfoTextView,
		pageFinalName,
		pageFinalTitle,
		pageFinalSubject,
		pageFinalDesc)

	assistant.Connect("cancel", func() {
		assistant.Close()
		assistant.Destroy()
		gtk.MainQuit()
	})
	assistant.Connect("close", func() {
		assistant.Close()
		assistant.Destroy()
		gtk.MainQuit()
	})

	mnemonic := wallet.GenerateMnemonic()

	assistant.Connect("prepare", func(assistant *gtk.Assistant, page *gtk.Widget) {
		name, err := page.GetName()
		errorCheck(err)
		fmt.Printf("%v - %v\n", assistant.GetCurrentPage(), name)
		switch name {
		case pageModeName:
			{
				assistant.SetPageComplete(pageMode, true)
			}
		case pageSeedName:
			{
				if getTextViewcontent(seedTextView) == "" {
					setTextViewcontent(seedTextView, mnemonic)
				}
				assistant.SetPageComplete(pageSeed, true)
			}
		case pageSeedConfirmName:
			{
			}
		case pagePasswordName:
			{
				assistant.SetPageComplete(pagePassword, true)
			}

		case pageFinalName:
			{
				defaultWallet, err := wallet.FromMnemonic(
					cmd.ZarbDefaultWalletPath(workingDir),
					mnemonic,
					"",
					0)
				errorCheck(err)
				valAddr, err := defaultWallet.MakeNewAddress("", "Validator address")
				errorCheck(err)
				rewardAddr, err := defaultWallet.MakeNewAddress("", "Reward address")
				errorCheck(err)

				// To make process faster we set password after generating addresses
				walletPassword, err := passwordEntry.GetText()
				errorCheck(err)
				err = defaultWallet.UpdatePassword("", walletPassword)
				errorCheck(err)
				err = defaultWallet.Save()
				errorCheck(err)
				err = genesis.Testnet().SaveToFile(cmd.ZarbGenesisPath(workingDir))
				errorCheck(err)

				conf := config.DefaultConfig()
				conf.Network.Name = "perdana-testnet"
				conf.Network.Bootstrap.Addresses = []string{"/ip4/172.104.169.94/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq"}
				conf.Network.Bootstrap.MinThreshold = 4
				conf.Network.Bootstrap.MaxThreshold = 8
				conf.State.RewardAddress = rewardAddr
				err = conf.SaveToFile(cmd.ZarbConfigPath(workingDir))
				errorCheck(err)

				successful = true
				nodeInfo := fmt.Sprintf("Working directory:\n  %s\n\n", workingDir)
				nodeInfo += fmt.Sprintf("Validator address:\n  %s\n\n", valAddr)
				nodeInfo += fmt.Sprintf("Reward address:\n  %s\n", rewardAddr)

				setTextViewcontent(NodeInfoTextView, nodeInfo)
			}
		}
	})

	assistant.SetPageType(pageMode, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageSeed, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageSeedConfirm, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pagePassword, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageFinal, gtk.ASSISTANT_PAGE_SUMMARY)

	assistant.SetModal(true)
	assistant.ShowAll()

	gtk.Main()
	return successful
}
