//go:build gtk

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/node/config"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/wallet"
)

func getTextViewcontent(tv *gtk.TextView) (string, error) {
	buf, _ := tv.GetBuffer()
	startIter, endIter := buf.GetBounds()
	content, err := buf.GetText(startIter, endIter, true)
	if err != nil {
		return "", err
	}
	return content, nil
}

func setTextViewcontent(tv *gtk.TextView, content string) error {
	buf, err := tv.GetBuffer()
	if err != nil {
		return err
	}
	buf.SetText(content)
	return err
}

func setMargin(widget gtk.IWidget, top, bottom, start, end int) {
	widget.ToWidget().SetMarginTop(top)
	widget.ToWidget().SetMarginBottom(bottom)
	widget.ToWidget().SetMarginStart(start)
	widget.ToWidget().SetMarginEnd(end)
}

func startupAssistant(workingDir string, testnet bool) bool {
	gtk.Init(nil)

	successful := false
	createPage := func(assistant *gtk.Assistant, content gtk.IWidget, name, title, subject, desc string) *gtk.Widget {
		page, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
		errorCheck(assistant, err)
		page.SetHExpand(true)
		titleLabel, err := gtk.LabelNew(title)
		errorCheck(assistant, err)
		setMargin(titleLabel, 0, 20, 0, 0)
		frame, err := gtk.FrameNew(subject)
		errorCheck(assistant, err)
		frame.SetHExpand(true)

		descLabel, err := gtk.LabelNew("")
		errorCheck(assistant, err)
		descLabel.SetUseMarkup(true)
		descLabel.SetMarkup(desc)
		descLabel.SetVExpand(true)
		descLabel.SetVAlign(gtk.ALIGN_END)
		descLabel.SetHAlign(gtk.ALIGN_START)
		setMargin(descLabel, 0, 0, 0, 0)
		frame.Add(content)

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
		errorCheck(assistant, err)
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
	errorCheck(assistant, err)

	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Zarb - Init Wizard")

	var pageMode *gtk.Widget
	var pagePassword *gtk.Widget
	var pageSeed *gtk.Widget
	var pageSeedConfirm *gtk.Widget
	var pageFinal *gtk.Widget

	// --- PageMode
	newWalletRadio, err := gtk.RadioButtonNewWithLabel(nil, "Create a new wallet from the scratch")
	errorCheck(assistant, err)
	recoverWalletRadio, err := gtk.RadioButtonNewWithLabelFromWidget(newWalletRadio, "Restore a wallet from the seed phrase")
	errorCheck(assistant, err)
	recoverWalletRadio.SetSensitive(false)

	radioBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	errorCheck(assistant, err)
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
	errorCheck(assistant, err)
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
	errorCheck(assistant, err)
	setMargin(seedConfirmTextView, 6, 6, 6, 6)
	seedConfirmTextView.SetWrapMode(gtk.WRAP_WORD)
	seedConfirmTextView.SetEditable(true)
	seedConfirmTextView.SetMonospace(true)
	seedConfirmTextView.SetSizeRequest(0, 80)

	seedConfirmTextView.Connect("paste_clipboard", func(textView *gtk.TextView) {
		showInfoDialog(assistant, "Opps, no copy paste!")
		seedConfirmTextView.StopEmission("paste_clipboard")
	})

	seedConfirmTextBuffer, err := seedConfirmTextView.GetBuffer()
	errorCheck(assistant, err)
	seedConfirmTextBuffer.Connect("changed", func(buf *gtk.TextBuffer) {
		mnemonic1, err := getTextViewcontent(seedTextView)
		errorCheck(assistant, err)
		mnemonic2, err := getTextViewcontent(seedConfirmTextView)
		errorCheck(assistant, err)
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
	errorCheck(assistant, err)
	setMargin(passwordEntry, 6, 6, 6, 6)
	passwordEntry.SetVisibility(false)
	passwordLabel, err := gtk.LabelNew("Password: ")
	errorCheck(assistant, err)
	passwordLabel.SetHAlign(gtk.ALIGN_START)
	setMargin(passwordLabel, 6, 6, 6, 6)

	passwordConfirmEntry, err := gtk.EntryNew()
	errorCheck(assistant, err)
	setMargin(passwordConfirmEntry, 6, 6, 6, 6)
	passwordConfirmEntry.SetVisibility(false)
	confirmationLineLabel, err := gtk.LabelNew("Confirmation: ")
	errorCheck(assistant, err)
	confirmationLineLabel.SetHAlign(gtk.ALIGN_START)
	setMargin(confirmationLineLabel, 6, 6, 6, 6)

	grid, err := gtk.GridNew()
	errorCheck(assistant, err)
	grid.Add(passwordLabel)
	grid.Attach(passwordEntry, 1, 0, 1, 1)
	grid.AttachNextTo(confirmationLineLabel, passwordLabel, gtk.POS_BOTTOM, 1, 1)
	grid.AttachNextTo(passwordConfirmEntry, passwordEntry, gtk.POS_BOTTOM, 1, 1)

	validatePassword := func() {
		pass1, err := passwordEntry.GetText()
		errorCheck(assistant, err)
		pass2, err := passwordConfirmEntry.GetText()
		errorCheck(assistant, err)
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
	errorCheck(assistant, err)
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
		errorCheck(assistant, err)
		fmt.Printf("%v - %v\n", assistant.GetCurrentPage(), name)
		switch name {
		case pageModeName:
			{
				assistant.SetPageComplete(pageMode, true)
			}
		case pageSeedName:
			{
				text, _ := getTextViewcontent(seedTextView)
				if text == "" {
					err := setTextViewcontent(seedTextView, mnemonic)
					if err != nil {
						errorCheck(assistant, err)
					}
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
				network := wallet.NetworkMainNet
				if testnet {
					network = wallet.NetworkTestNet
				}
				defaultWallet, err := wallet.FromMnemonic(
					cmd.ZarbDefaultWalletPath(workingDir),
					mnemonic,
					"",
					network)
				errorCheck(assistant, err)
				valAddr, err := defaultWallet.MakeNewAddress("", "Validator address")
				errorCheck(assistant, err)
				rewardAddr, err := defaultWallet.MakeNewAddress("", "Reward address")
				errorCheck(assistant, err)

				var gen *genesis.Genesis
				confFile := cmd.ZarbConfigPath(workingDir)

				if testnet {
					gen = genesis.Testnet()

					// Save config for testnet
					if err := config.SaveTestnetConfig(confFile, rewardAddr); err != nil {
						cmd.PrintErrorMsg("Failed to write config file: %v", err)
						return
					}
				} else {
					panic("not yet!")
					// gen = genesis.Mainnet()

					// // Save config for mainnet
					// if err := config.SaveMainnetConfig(confFile, rewardAddr); err != nil {
					// 	cmd.PrintErrorMsg("Failed to write config file: %v", err)
					// 	return
					// }
				}

				// Save genesis file
				genFile := cmd.ZarbGenesisPath(workingDir)
				err = gen.SaveToFile(genFile)
				errorCheck(assistant, err)

				// To make process faster we set password after generating addresses
				walletPassword, err := passwordEntry.GetText()
				errorCheck(assistant, err)

				err = defaultWallet.UpdatePassword("", walletPassword)
				errorCheck(assistant, err)

				// Save wallet
				err = defaultWallet.Save()
				errorCheck(assistant, err)

				// Done! showing the node information
				successful = true
				nodeInfo := fmt.Sprintf("Working directory:\n  %s\n\n", workingDir)
				nodeInfo += fmt.Sprintf("Validator address:\n  %s\n\n", valAddr)
				nodeInfo += fmt.Sprintf("Reward address:\n  %s\n", rewardAddr)

				err = setTextViewcontent(NodeInfoTextView, nodeInfo)
				errorCheck(assistant, err)
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
