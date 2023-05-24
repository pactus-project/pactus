//go:build gtk

package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet"
)

func setMargin(widget gtk.IWidget, top, bottom, start, end int) {
	widget.ToWidget().SetMarginTop(top)
	widget.ToWidget().SetMarginBottom(bottom)
	widget.ToWidget().SetMarginStart(start)
	widget.ToWidget().SetMarginEnd(end)
}

func startupAssistant(workingDir string, chain genesis.ChainType) bool {
	successful := false
	createPage := func(assistant *gtk.Assistant, content gtk.IWidget, name,
		title, subject, desc string) *gtk.Widget {
		page, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
		fatalErrorCheck(err)

		page.SetHExpand(true)

		frame, err := gtk.FrameNew(subject)
		fatalErrorCheck(err)

		frame.SetHExpand(true)

		labelDesc, err := gtk.LabelNew("")
		fatalErrorCheck(err)

		labelDesc.SetUseMarkup(true)
		labelDesc.SetMarkup("<span allow_breaks='true'>" + desc + "</span>")
		labelDesc.SetVExpand(true)
		labelDesc.SetVAlign(gtk.ALIGN_END)
		labelDesc.SetHAlign(gtk.ALIGN_START)
		setMargin(labelDesc, 0, 0, 0, 0)
		frame.Add(content)

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
		fatalErrorCheck(err)

		box.Add(frame)
		box.Add(labelDesc)
		page.Add(box)

		page.SetName(name)
		assistant.AppendPage(page)
		assistant.SetPageTitle(page, title)

		return page.ToWidget()
	}

	assistant, err := gtk.AssistantNew()
	fatalErrorCheck(err)

	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Pactus - Init Wizard")

	var pageMode *gtk.Widget
	var pagePassword *gtk.Widget
	var pageSeed *gtk.Widget
	var pageSeedConfirm *gtk.Widget
	var pageNumValidators *gtk.Widget
	var pageFinal *gtk.Widget

	// --- PageMode
	newWalletRadio, err := gtk.RadioButtonNewWithLabel(nil, "Create a new wallet from the scratch")
	fatalErrorCheck(err)

	recoverWalletRadio, err := gtk.RadioButtonNewWithLabelFromWidget(newWalletRadio,
		"Restore a wallet from the seed phrase")
	fatalErrorCheck(err)

	recoverWalletRadio.SetSensitive(false)

	radioBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	fatalErrorCheck(err)

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
	textViewSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewSeed, 6, 6, 6, 6)
	textViewSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewSeed.SetEditable(false)
	textViewSeed.SetMonospace(true)
	textViewSeed.SetSizeRequest(0, 80)

	pageSeedName := "page_seed"
	pageSeedTitle := "Wallet seed"
	pageSeedSubject := "Your wallet generation seed is:"
	pageSeedDesc := `Please write these 12 words on paper.
This seed will allow you to recover your wallet in case of computer failure.
<b>WARNING:</b>
  - Never disclose your seed.
  - Never type it on a website.
  - Do not store it electronically.`

	pageSeed = createPage(
		assistant,
		textViewSeed,
		pageSeedName,
		pageSeedTitle,
		pageSeedSubject,
		pageSeedDesc)

	// --- pageSeedConfirm
	textViewConfirmSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewConfirmSeed, 6, 6, 6, 6)
	textViewConfirmSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewConfirmSeed.SetEditable(true)
	textViewConfirmSeed.SetMonospace(true)
	textViewConfirmSeed.SetSizeRequest(0, 80)

	textViewConfirmSeed.Connect("paste_clipboard", func(textView *gtk.TextView) {
		showInfoDialog(assistant, "Opps, no copy paste!")
		textViewConfirmSeed.StopEmission("paste_clipboard")
	})

	seedConfirmTextBuffer, err := textViewConfirmSeed.GetBuffer()
	fatalErrorCheck(err)

	seedConfirmTextBuffer.Connect("changed", func(buf *gtk.TextBuffer) {
		mnemonic1 := getTextViewContent(textViewSeed)
		mnemonic2 := getTextViewContent(textViewConfirmSeed)
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
		textViewConfirmSeed,
		pageSeedConfirmName,
		pageSeedConfirmTitle,
		pageSeedConfirmSubject,
		pageSeedConfirmDesc)

	// --- PagePassword
	entryPassword, err := gtk.EntryNew()
	fatalErrorCheck(err)

	setMargin(entryPassword, 6, 6, 6, 6)
	entryPassword.SetVisibility(false)
	labelConfirmPassword, err := gtk.LabelNew("Password: ")
	fatalErrorCheck(err)

	labelConfirmPassword.SetHAlign(gtk.ALIGN_START)
	setMargin(labelConfirmPassword, 6, 6, 6, 6)

	entryConfirmPassword, err := gtk.EntryNew()
	fatalErrorCheck(err)

	setMargin(entryConfirmPassword, 6, 6, 6, 6)
	entryConfirmPassword.SetVisibility(false)
	labelConfirmation, err := gtk.LabelNew("Confirmation: ")
	fatalErrorCheck(err)

	labelConfirmation.SetHAlign(gtk.ALIGN_START)
	setMargin(labelConfirmation, 6, 6, 6, 6)

	grid, err := gtk.GridNew()
	fatalErrorCheck(err)

	grid.Add(labelConfirmPassword)
	grid.Attach(entryPassword, 1, 0, 1, 1)
	grid.AttachNextTo(labelConfirmation, labelConfirmPassword, gtk.POS_BOTTOM, 1, 1)
	grid.AttachNextTo(entryConfirmPassword, entryPassword, gtk.POS_BOTTOM, 1, 1)

	validatePassword := func() {
		pass1, err := entryPassword.GetText()
		fatalErrorCheck(err)

		pass2, err := entryConfirmPassword.GetText()
		fatalErrorCheck(err)

		if pass1 == pass2 {
			assistant.SetPageComplete(pagePassword, true)
		} else {
			assistant.SetPageComplete(pagePassword, false)
		}
	}
	entryPassword.Connect("changed", func(entry *gtk.Entry) {
		validatePassword()
	})

	entryConfirmPassword.Connect("changed", func(entry *gtk.Entry) {
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

	// --- pageNumValidators
	lsNumValidators, err := gtk.ListStoreNew(glib.TYPE_INT)
	fatalErrorCheck(err)

	for i := 0; i < 32; i++ {
		iter := lsNumValidators.Append()
		err = lsNumValidators.SetValue(iter, 0, i+1)
		fatalErrorCheck(err)
	}

	comboNumValidators, err := gtk.ComboBoxNewWithModel(lsNumValidators)
	fatalErrorCheck(err)

	cellRenderer, err := gtk.CellRendererTextNew()
	fatalErrorCheck(err)

	// Set the default selected value to 7 (index 6)
	comboNumValidators.SetActive(6)

	comboNumValidators.PackStart(cellRenderer, true)
	comboNumValidators.AddAttribute(cellRenderer, "text", 0)

	labelNumValidators, err := gtk.LabelNew("Number of validators: ")
	fatalErrorCheck(err)

	setMargin(labelNumValidators, 6, 6, 6, 6)
	setMargin(comboNumValidators, 6, 6, 6, 6)

	grid, err = gtk.GridNew()
	fatalErrorCheck(err)

	grid.Add(labelNumValidators)
	grid.Attach(comboNumValidators, 1, 0, 1, 1)

	pageNumValidatorsName := "page_num_validators"
	pageNumValidatorsTitle := "Number of validators"
	pageNumValidatorsSubject := "How many validators do you want to create?"
	pageNumValidatorsDesc := `Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.
You can define validators based on the amount of coins you want to stake.
For more information, look <a href="https://pactus.org/user-guides/run-pactus-gui/">here</a>`

	pageNumValidators = createPage(
		assistant,
		grid,
		pageNumValidatorsName,
		pageNumValidatorsTitle,
		pageNumValidatorsSubject,
		pageNumValidatorsDesc)

	// --- pageFinal
	textViewNodeInfo, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewNodeInfo, 6, 6, 6, 6)
	textViewNodeInfo.SetWrapMode(gtk.WRAP_WORD)
	textViewNodeInfo.SetEditable(false)
	textViewNodeInfo.SetMonospace(true)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	fatalErrorCheck(err)

	scrolledWindow.SetSizeRequest(0, 300)
	scrolledWindow.Add(textViewNodeInfo)

	pageFinalName := "page_final"
	pageFinalTitle := "Node info"
	pageFinalSubject := "Your node information:"
	pageFinalDesc := `Congratulation. Your node is initialized successfully.
Now you are ready to start the node!`

	pageFinal = createPage(
		assistant,
		scrolledWindow,
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

	mnemonic := wallet.GenerateMnemonic(128)

	assistant.Connect("prepare", func(assistant *gtk.Assistant, page *gtk.Widget) {
		name, err := page.GetName()
		fatalErrorCheck(err)

		log.Printf("%v - %v\n", assistant.GetCurrentPage(), name)
		switch name {
		case pageModeName:
			{
				assistant.SetPageComplete(pageMode, true)
			}
		case pageSeedName:
			{
				text := getTextViewContent(textViewSeed)
				if text == "" {
					setTextViewContent(textViewSeed, mnemonic)
				}
				assistant.SetPageComplete(pageSeed, true)
			}
		case pageSeedConfirmName:
			{
				assistant.SetPageComplete(pageSeedConfirm, false)
			}
		case pagePasswordName:
			{
				assistant.SetPageComplete(pagePassword, true)
			}

		case pageNumValidatorsName:
			{
				assistant.SetPageComplete(pageNumValidators, true)
			}

		case pageFinalName:
			{
				iter, err := comboNumValidators.GetActiveIter()
				fatalErrorCheck(err)

				val, err := lsNumValidators.GetValue(iter, 0)
				fatalErrorCheck(err)

				valueInterface, err := val.GoValue()
				fatalErrorCheck(err)

				numValidators := valueInterface.(int)

				fmt.Println("number of validators:", numValidators)

				walletPassword, err := entryPassword.GetText()
				fatalErrorCheck(err)

				validatorAddrs, rewardAddrs, err := cmd.CreateNode(numValidators, chain, workingDir, mnemonic, walletPassword)
				fatalErrorCheck(err)

				// Done! showing the node information
				successful = true
				nodeInfo := fmt.Sprintf("Working directory: %s\n", workingDir)
				nodeInfo += fmt.Sprintf("Network: %s\n", chain.String())
				nodeInfo += "\nValidator addresses:\n"
				for i, addr := range validatorAddrs {
					nodeInfo += fmt.Sprintf("%v- %s\n", i+1, addr)
				}

				nodeInfo += "\nReward addresses:\n"
				for i, addr := range rewardAddrs {
					nodeInfo += fmt.Sprintf("%v- %s\n", i+1, addr)
				}

				setTextViewContent(textViewNodeInfo, nodeInfo)
			}
		}
	})

	assistant.SetPageType(pageMode, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageSeed, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageSeedConfirm, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pagePassword, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageNumValidators, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(pageFinal, gtk.ASSISTANT_PAGE_SUMMARY)

	assistant.SetModal(true)
	assistant.ShowAll()

	gtk.Main()
	return successful
}
