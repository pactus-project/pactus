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

type assistantFunc func(assistant *gtk.Assistant, content gtk.IWidget, name,
	title, subject, desc string) *gtk.Widget

func setMargin(widget gtk.IWidget, top, bottom, start, end int) {
	widget.ToWidget().SetMarginTop(top)
	widget.ToWidget().SetMarginBottom(bottom)
	widget.ToWidget().SetMarginStart(start)
	widget.ToWidget().SetMarginEnd(end)
}

func startupAssistant(workingDir string, chain genesis.ChainType) bool {
	successful := false
	assistant, err := gtk.AssistantNew()
	fatalErrorCheck(err)

	assistant.Hide()
	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Pactus - Init Wizard")

	assistFunc := pageAssistant()

	// --- PageMode
	mode, restoreRadio, pageModeName := pageMode(assistant, assistFunc)

	// --- seedGenerate
	seedGenerate, textViewSeed, pageSeedGenerateName := pageSeedGenerate(assistant, assistFunc)

	// -- seedRestore
	seedRestore, pageSeedRestoreName := pageSeedRestore(assistant, assistFunc)

	// --- seedConfirm
	seedConfirm, pageSeedConfirmName := pageSeedConfirm(assistant, assistFunc, textViewSeed)

	// --- PagePassword
	password, entryPassword, pagePasswordName := pagePassword(assistant, assistFunc)

	// --- numValidators
	numValidators, lsNumValidators, comboNumValidators, pageNumValidatorsName := pageNumValidators(assistant, assistFunc)

	// --- final
	final, textViewNodeInfo, pageFinalName := pageFinal(assistant, assistFunc)

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

	assistant.AppendPage(mode)          // page 0
	assistant.AppendPage(seedGenerate)  // page 1
	assistant.AppendPage(seedConfirm)   // page 2
	assistant.AppendPage(seedRestore)   // page 3
	assistant.AppendPage(password)      // page 4
	assistant.AppendPage(numValidators) // page 5
	assistant.AppendPage(final)         // page 6

	assistant.SetPageType(mode, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(seedGenerate, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(seedConfirm, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(seedRestore, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(password, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(numValidators, gtk.ASSISTANT_PAGE_CONTENT)
	assistant.SetPageType(final, gtk.ASSISTANT_PAGE_SUMMARY)

	mnemonic := ""
	prevPageName := ""
	assistant.Connect("prepare", func(assistant *gtk.Assistant, page *gtk.Widget) {
		isRestoreMode := restoreRadio.GetActive()
		curPageName, err := page.GetName()
		fatalErrorCheck(err)

		log.Printf("%v - %v (restore: %v, prev: %v)\n",
			assistant.GetCurrentPage(), curPageName, isRestoreMode, prevPageName)
		switch curPageName {
		case pageModeName:
			{
				assistant.SetPageComplete(mode, true)
				assistant.UpdateButtonsState()
			}

		case pageSeedGenerateName:
			{
				if isRestoreMode {
					if prevPageName == pageSeedGenerateName {
						// backward
						log.Printf("jumping backward from seedGenerate page")
						assistant.PreviousPage()
					} else if prevPageName == pageModeName {
						// forward
						log.Printf("jumping forward from seedGenerate page")
						assistant.NextPage()
					} else {
						log.Fatalf("invalid1 page order, pageName: %v, prevPageName: %v",
							curPageName, prevPageName)
					}
				} else {
					mnemonic = wallet.GenerateMnemonic(128)
					setTextViewContent(textViewSeed, mnemonic)
				}
				assistant.SetPageComplete(seedGenerate, true)
			}
		case pageSeedConfirmName:
			{
				if isRestoreMode {
					if prevPageName == pageSeedGenerateName {
						// backward
						log.Printf("jumping backward from seedConfirm page")
						assistant.PreviousPage()
					} else if prevPageName == pageModeName {
						// forward
						log.Printf("jumping forward from seedConfirm page")
						assistant.NextPage()
					} else {
						log.Fatalf("invalid2 page order, pageName: %v, prevPageName: %v",
							curPageName, prevPageName)
					}
				} else {
					assistant.SetPageComplete(seedConfirm, false)
				}
			}
		case pageSeedRestoreName:
			{
				if !isRestoreMode {
					if prevPageName == pageSeedRestoreName {
						// backward
						log.Printf("jumping backward from seedRestore page")
						assistant.PreviousPage()
					} else if prevPageName == pageSeedConfirmName {
						// forward
						log.Printf("jumping forward from seedRestore page")
						assistant.NextPage()
					} else {
						log.Fatalf("invalid page order, pageName: %v, prevPageName: %v",
							curPageName, prevPageName)
					}
				} else {
					assistant.SetPageComplete(seedGenerate, true)
				}

			}
		case pagePasswordName:
			{
				assistant.SetPageComplete(password, true)
			}
		case pageNumValidatorsName:
			{
				assistant.SetPageComplete(numValidators, true)
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
		prevPageName = curPageName
	})

	assistant.SetModal(true)
	assistant.ShowAll()

	gtk.Main()
	return successful
}

func pageAssistant() assistantFunc {
	return func(assistant *gtk.Assistant, content gtk.IWidget, name, title, subject, desc string) *gtk.Widget {
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

		return page.ToWidget()
	}
}

func pageMode(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.RadioButton, string) {
	mode := new(gtk.Widget)
	newWalletRadio, err := gtk.RadioButtonNewWithLabel(nil, "Create a new wallet from the scratch")
	fatalErrorCheck(err)

	restoreWalletRadio, err := gtk.RadioButtonNewWithLabelFromWidget(newWalletRadio,
		"Restore a wallet from the seed phrase")
	fatalErrorCheck(err)

	radioBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	fatalErrorCheck(err)

	radioBox.Add(newWalletRadio)
	setMargin(newWalletRadio, 6, 6, 6, 6)
	radioBox.Add(restoreWalletRadio)
	setMargin(restoreWalletRadio, 6, 6, 6, 6)

	pageModeName := "page_mode"
	pageModeTitle := "Initialize mode"
	pageModeSubject := "How to create your wallet?"
	pageModeDesc := "If you are running the node for the first time, choose the first option."
	mode = assistFunc(
		assistant,
		radioBox,
		pageModeName,
		pageModeTitle,
		pageModeSubject,
		pageModeDesc)
	return mode, restoreWalletRadio, pageModeName
}

func pageSeedGenerate(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	pageWidget := new(gtk.Widget)
	textViewSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewSeed, 6, 6, 6, 6)
	textViewSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewSeed.SetEditable(false)
	textViewSeed.SetMonospace(true)
	textViewSeed.SetSizeRequest(0, 80)

	pageSeedName := "page_seed_generate"
	pageSeedTitle := "Wallet seed"
	pageSeedSubject := "Your wallet generation seed is:"
	pageSeedDesc := `Please write these 12 words on paper.
This seed will allow you to recover your wallet in case of computer failure.
<b>WARNING:</b>
  - Never disclose your seed.
  - Never type it on a website.
  - Do not store it electronically.`

	pageWidget = assistFunc(
		assistant,
		textViewSeed,
		pageSeedName,
		pageSeedTitle,
		pageSeedSubject,
		pageSeedDesc)
	return pageWidget, textViewSeed, pageSeedName
}

func pageSeedRestore(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, string) {
	pageWidget := new(gtk.Widget)
	textViewRestoreSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewRestoreSeed, 6, 6, 6, 6)
	textViewRestoreSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewRestoreSeed.SetEditable(true)
	textViewRestoreSeed.SetMonospace(true)
	textViewRestoreSeed.SetSizeRequest(0, 80)

	textViewRestoreSeed.Connect("paste_clipboard", func(textView *gtk.TextView) {
		showInfoDialog(assistant, "Opps, no copy paste!")
		textViewRestoreSeed.StopEmission("paste_clipboard")
	})

	seedConfirmTextBuffer, err := textViewRestoreSeed.GetBuffer()
	fatalErrorCheck(err)

	seedConfirmTextBuffer.Connect("changed", func(buf *gtk.TextBuffer) {

	})

	pageSeedName := "page_seed_restore"
	pageSeedTitle := "Wallet seed restore"
	pageSeedSubject := "Please enter your seed:"
	pageSeedDesc := "Please enter your 12 words mnemonics backup to restore your wallet."

	pageWidget = assistFunc(
		assistant,
		textViewRestoreSeed,
		pageSeedName,
		pageSeedTitle,
		pageSeedSubject,
		pageSeedDesc)
	return pageWidget, pageSeedName
}

func pageSeedConfirm(assistant *gtk.Assistant, assistFunc assistantFunc, textViewSeed *gtk.TextView) (*gtk.Widget, string) {
	pageWidget := new(gtk.Widget)
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
			assistant.SetPageComplete(pageWidget, true)
		} else {
			assistant.SetPageComplete(pageWidget, false)
		}
	})

	pageSeedConfirmName := "page_seed_confirm"
	pageSeedConfirmTitle := "Confirm seed"
	pageSeedConfirmSubject := "What was your seed?"
	pageSeedConfirmDesc := `Your seed is important!
To make sure that you have properly saved your seed, please retype it here.`

	pageWidget = assistFunc(
		assistant,
		textViewConfirmSeed,
		pageSeedConfirmName,
		pageSeedConfirmTitle,
		pageSeedConfirmSubject,
		pageSeedConfirmDesc)
	return pageWidget, pageSeedConfirmName
}

func pagePassword(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.Entry, string) {
	pageWidget := new(gtk.Widget)
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
			assistant.SetPageComplete(pageWidget, true)
		} else {
			assistant.SetPageComplete(pageWidget, false)
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

	pageWidget = assistFunc(
		assistant,
		grid,
		pagePasswordName,
		pagePasswordTitle,
		pagePasswordSubject,
		pagePsswrdDesc)
	return pageWidget, entryPassword, pagePasswordName
}

func pageNumValidators(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.ListStore, *gtk.ComboBox, string) {
	pageWidget := new(gtk.Widget)
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

	grid, err := gtk.GridNew()
	fatalErrorCheck(err)

	grid.Add(labelNumValidators)
	grid.Attach(comboNumValidators, 1, 0, 1, 1)

	pageNumValidatorsName := "page_num_validators"
	pageNumValidatorsTitle := "Number of validators"
	pageNumValidatorsSubject := "How many validators do you want to create?"
	pageNumValidatorsDesc := `Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.
You can define validators based on the amount of coins you want to stake.
For more information, look <a href="https://pactus.org/user-guides/run-pactus-gui/">here</a>`

	pageWidget = assistFunc(
		assistant,
		grid,
		pageNumValidatorsName,
		pageNumValidatorsTitle,
		pageNumValidatorsSubject,
		pageNumValidatorsDesc)
	return pageWidget, lsNumValidators, comboNumValidators, pageNumValidatorsName
}

func pageFinal(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	pageWidget := new(gtk.Widget)
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

	pageWidget = assistFunc(
		assistant,
		scrolledWindow,
		pageFinalName,
		pageFinalTitle,
		pageFinalSubject,
		pageFinalDesc)
	return pageWidget, textViewNodeInfo, pageFinalName
}
