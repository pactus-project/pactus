//go:build gtk

package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
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

//nolint:all  // complexity can't be reduced more. It needs to refactor.
func startupAssistant(workingDir string, chainType genesis.ChainType) bool {
	successful := false
	assistant, err := gtk.AssistantNew()
	fatalErrorCheck(err)

	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Pactus - Init Wizard")

	assistFunc := pageAssistant()

	// --- page_mode
	wgtWalletMode, radioRestoreWallet, pageModeName := pageWalletMode(assistant, assistFunc)

	// --- page_seed_generate
	wgtSeedGenerate, txtSeed, pageSeedGenerateName := pageSeedGenerate(assistant, assistFunc)

	// --- page_seed_confirm
	wgtSeedConfirm, pageSeedConfirmName := pageSeedConfirm(assistant, assistFunc, txtSeed)

	// -- page_seed_restore
	wgtSeedRestore, textRestoreSeed, pageSeedRestoreName := pageSeedRestore(assistant, assistFunc)

	// --- page_password
	wgtPassword, entryPassword, pagePasswordName := pagePassword(assistant, assistFunc)

	// --- page_num_validators
	wgtNumValidators, lsNumValidators, comboNumValidators,
		pageNumValidatorsName := pageNumValidators(assistant, assistFunc)

	// -- page_node_type
	wgtNodeType, gridImport, radioImport, pageNodeTypeName := pageNodeType(assistant, assistFunc)

	// --- page_summary
	wgtSummary, txtNodeInfo, pageSummaryName := pageSummary(assistant, assistFunc)

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

	assistant.SetPageType(wgtWalletMode, gtk.ASSISTANT_PAGE_INTRO)      // page 0
	assistant.SetPageType(wgtSeedGenerate, gtk.ASSISTANT_PAGE_CONTENT)  // page 1
	assistant.SetPageType(wgtSeedConfirm, gtk.ASSISTANT_PAGE_CONTENT)   // page 2
	assistant.SetPageType(wgtSeedRestore, gtk.ASSISTANT_PAGE_CONTENT)   // page 3
	assistant.SetPageType(wgtPassword, gtk.ASSISTANT_PAGE_CONTENT)      // page 4
	assistant.SetPageType(wgtNumValidators, gtk.ASSISTANT_PAGE_CONTENT) // page 5
	assistant.SetPageType(wgtNodeType, gtk.ASSISTANT_PAGE_CONTENT)      // page 6
	assistant.SetPageType(wgtSummary, gtk.ASSISTANT_PAGE_SUMMARY)       // page 7

	mnemonic := ""
	prevPageIndex := -1
	prevPageAdjust := 0
	assistant.Connect("prepare", func(assistant *gtk.Assistant, page *gtk.Widget) {
		isRestoreMode := radioRestoreWallet.GetActive()
		curPageName, err := page.GetName()
		curPageIndex := assistant.GetCurrentPage()
		fatalErrorCheck(err)

		isForward := true
		if curPageIndex > 0 && curPageIndex < prevPageIndex {
			isForward = false
		}

		log.Printf("%v (restore: %v, prev: %v, cur: %v)\n",
			curPageName, isRestoreMode, prevPageIndex, curPageIndex)
		switch curPageName {
		case pageModeName:
			assistantPageComplete(assistant, wgtWalletMode, true)

		case pageSeedGenerateName:
			if isRestoreMode {
				if isForward {
					// forward
					log.Printf("jumping forward from seedGenerate page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Printf("jumping backward from seedGenerate page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}
				assistantPageComplete(assistant, wgtSeedGenerate, false)
			} else {
				mnemonic, _ = wallet.GenerateMnemonic(128)
				setTextViewContent(txtSeed, mnemonic)
				assistantPageComplete(assistant, wgtSeedGenerate, true)
			}
		case pageSeedConfirmName:
			if isRestoreMode {
				if isForward {
					// forward
					log.Printf("jumping forward from seedConfirm page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Printf("jumping backward from seedConfirm page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}
				assistantPageComplete(assistant, wgtSeedConfirm, false)
			} else {
				assistantPageComplete(assistant, wgtSeedConfirm, false)
			}
		case pageSeedRestoreName:
			if !isRestoreMode {
				if isForward {
					// forward
					log.Printf("jumping forward from seedRestore page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Printf("jumping backward from seedRestore page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}
				assistantPageComplete(assistant, wgtSeedConfirm, false)
			} else {
				assistantPageComplete(assistant, wgtSeedRestore, true)
			}
		case pagePasswordName:
			if isRestoreMode {
				mnemonic = getTextViewContent(textRestoreSeed)

				if err := wallet.CheckMnemonic(mnemonic); err != nil {
					showErrorDialog(assistant, "mnemonic is invalid")
					assistant.PreviousPage()
				}
			}
			assistantPageComplete(assistant, wgtPassword, true)
		case pageNumValidatorsName:
			assistantPageComplete(assistant, wgtNumValidators, true)

		case pageNodeTypeName:
			assistantPageComplete(assistant, wgtNodeType, true)
			ssLabel, err := gtk.LabelNew("")
			fatalErrorCheck(err)
			setMargin(ssLabel, 5, 5, 1, 1)
			ssLabel.SetHAlign(gtk.ALIGN_START)

			listBox, err := gtk.ListBoxNew()
			fatalErrorCheck(err)
			setMargin(listBox, 5, 5, 1, 1)
			listBox.SetHAlign(gtk.ALIGN_CENTER)
			listBox.SetSizeRequest(700, -1)

			ssDLBtn, err := gtk.ButtonNewWithLabel("⏬ Download")
			fatalErrorCheck(err)
			setMargin(ssDLBtn, 10, 5, 1, 1)
			ssDLBtn.SetHAlign(gtk.ALIGN_CENTER)
			ssDLBtn.SetSizeRequest(700, -1)

			ssPBLabel, err := gtk.LabelNew("")
			fatalErrorCheck(err)
			setMargin(ssPBLabel, 5, 10, 1, 1)
			ssPBLabel.SetHAlign(gtk.ALIGN_START)

			gridImport.Attach(ssLabel, 0, 1, 1, 1)
			gridImport.Attach(listBox, 0, 2, 1, 1)
			gridImport.Attach(ssDLBtn, 0, 3, 1, 1)
			gridImport.Attach(ssPBLabel, 0, 5, 1, 1)
			ssLabel.SetVisible(false)
			listBox.SetVisible(false)
			ssDLBtn.SetVisible(false)
			ssPBLabel.SetVisible(false)

			snapshotIndex := 0

			radioImport.Connect("toggled", func() {
				if radioImport.GetActive() {
					assistantPageComplete(assistant, wgtNodeType, false)

					ssLabel.SetVisible(true)
					ssLabel.SetText("   ♻️ Please wait, loading snapshot list...")

					go func() {
						time.Sleep(1 * time.Second)

						glib.IdleAdd(func() {
							snapshotURL := cmd.DefaultSnapshotURL // TODO: make me optional...

							storeDir := filepath.Join(workingDir, "data")
							importer, err := cmd.NewImporter(
								chainType,
								snapshotURL,
								storeDir,
							)
							if err != nil {
								showError(err)

								return
							}

							ctx := context.Background()
							mdCh := getMetadata(ctx, importer, listBox)

							if md := <-mdCh; md == nil {
								ssLabel.SetText("   ❌ Failed to get snapshot list, please try again later.")
							} else {
								ssLabel.SetText("   🔽 Please select a snapshot to download:")
								listBox.SetVisible(true)

								listBox.Connect("row-selected", func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
									if row != nil {
										snapshotIndex = row.GetIndex()
										ssDLBtn.SetVisible(true)
									}
								})

								ssDLBtn.Connect("clicked", func() {
									radioGroup, _ := radioImport.GetParent()
									radioImport.SetSensitive(false)
									radioGroup.ToWidget().SetSensitive(false)
									ssLabel.SetSensitive(false)
									listBox.SetSensitive(false)
									ssDLBtn.SetSensitive(false)

									ssDLBtn.SetVisible(false)
									ssPBLabel.SetVisible(true)
									listBox.SetSelectionMode(gtk.SELECTION_NONE)

									go func() {
										log.Printf("start downloading...\n")

										err := importer.Download(ctx, &md[snapshotIndex],
											func(fileName string) func(stats downloader.Stats) {
												return func(stats downloader.Stats) {
													if !stats.Completed {
														percent := int(stats.Percent)
														glib.IdleAdd(func() {
															dlMessage := fmt.Sprintf("🌐 Downloading %s | %d%% (%s / %s)",
																fileName,
																percent,
																util.FormatBytesToHumanReadable(uint64(stats.Downloaded)),
																util.FormatBytesToHumanReadable(uint64(stats.TotalSize)),
															)
															ssPBLabel.SetText("   " + dlMessage)
														})
													}
												}
											},
										)

										glib.IdleAdd(func() {
											fatalErrorCheck(err)

											log.Printf("extracting data...\n")
											ssPBLabel.SetText("   " + "📂 Extracting downloaded files...")
											err := importer.ExtractAndStoreFiles()
											fatalErrorCheck(err)

											log.Printf("moving data...\n")
											ssPBLabel.SetText("   " + "📑 Moving data...")
											err = importer.MoveStore()
											fatalErrorCheck(err)

											log.Printf("cleanup...\n")
											err = importer.Cleanup()
											fatalErrorCheck(err)

											ssPBLabel.SetText("   " + "✅ Import completed.")
											assistantPageComplete(assistant, wgtNodeType, true)
										})
									}()
								})
							}
						})
					}()
				} else {
					assistantPageComplete(assistant, wgtNodeType, true)
					ssLabel.SetVisible(false)
					listBox.SetVisible(false)
					ssDLBtn.SetVisible(false)
					ssPBLabel.SetVisible(false)
				}
			})
		case pageSummaryName:
			iter, err := comboNumValidators.GetActiveIter()
			fatalErrorCheck(err)

			val, err := lsNumValidators.GetValue(iter, 0)
			fatalErrorCheck(err)

			valueInterface, err := val.GoValue()
			fatalErrorCheck(err)

			numValidators := valueInterface.(int)
			walletPassword, err := entryPassword.GetText()
			fatalErrorCheck(err)

			validatorAddrs, rewardAddrs, err := cmd.CreateNode(context.Background(), numValidators,
				chainType, workingDir, mnemonic, walletPassword, nil)
			fatalErrorCheck(err)

			// Done! showing the node information
			successful = true
			nodeInfo := fmt.Sprintf("Working directory: %s\n", workingDir)
			nodeInfo += fmt.Sprintf("Network: %s\n", chainType.String())
			nodeInfo += "\nValidator addresses:\n"
			for i, addr := range validatorAddrs {
				nodeInfo += fmt.Sprintf("%v- %s\n", i+1, addr)
			}

			nodeInfo += "\nReward address:\n"
			nodeInfo += fmt.Sprintf("%s", rewardAddrs)

			setTextViewContent(txtNodeInfo, nodeInfo)
		}
		prevPageIndex = curPageIndex + prevPageAdjust
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
		assistant.AppendPage(page)
		assistant.SetPageTitle(page, title)

		return page.ToWidget()
	}
}

func pageWalletMode(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.RadioButton, string) {
	var mode *gtk.Widget
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

	pageModeName := "page_wallet_mode"
	pageModeTitle := "Wallet Mode"
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
	var pageWidget *gtk.Widget
	textViewSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewSeed, 6, 6, 6, 6)
	textViewSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewSeed.SetEditable(false)
	textViewSeed.SetMonospace(true)
	textViewSeed.SetSizeRequest(0, 80)

	pageSeedName := "page_seed_generate"
	pageSeedTitle := "Wallet Seed"
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

func pageSeedRestore(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	var pageWidget *gtk.Widget
	textViewRestoreSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewRestoreSeed, 6, 6, 6, 6)
	textViewRestoreSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewRestoreSeed.SetEditable(true)
	textViewRestoreSeed.SetMonospace(true)
	textViewRestoreSeed.SetSizeRequest(0, 80)

	pageSeedName := "page_seed_restore"
	pageSeedTitle := "Wallet Seed Restore"
	pageSeedSubject := "Enter your wallet seed:"
	pageSeedDesc := "Please enter your 12 words mnemonics backup to restore your wallet."

	pageWidget = assistFunc(
		assistant,
		textViewRestoreSeed,
		pageSeedName,
		pageSeedTitle,
		pageSeedSubject,
		pageSeedDesc)

	return pageWidget, textViewRestoreSeed, pageSeedName
}

func pageSeedConfirm(assistant *gtk.Assistant, assistFunc assistantFunc,
	textViewSeed *gtk.TextView,
) (*gtk.Widget, string) {
	pageWidget := new(gtk.Widget)
	textViewConfirmSeed, err := gtk.TextViewNew()
	fatalErrorCheck(err)

	setMargin(textViewConfirmSeed, 6, 6, 6, 6)
	textViewConfirmSeed.SetWrapMode(gtk.WRAP_WORD)
	textViewConfirmSeed.SetEditable(true)
	textViewConfirmSeed.SetMonospace(true)
	textViewConfirmSeed.SetSizeRequest(0, 80)

	textViewConfirmSeed.Connect("paste_clipboard", func(_ *gtk.TextView) {
		showInfoDialog(assistant, "Opps, no copy paste!")
		textViewConfirmSeed.StopEmission("paste_clipboard")
	})

	seedConfirmTextBuffer, err := textViewConfirmSeed.GetBuffer()
	fatalErrorCheck(err)

	seedConfirmTextBuffer.Connect("changed", func(_ *gtk.TextBuffer) {
		mnemonic1 := getTextViewContent(textViewSeed)
		mnemonic2 := getTextViewContent(textViewConfirmSeed)
		space := regexp.MustCompile(`\s+`)
		mnemonic2 = space.ReplaceAllString(mnemonic2, " ")
		mnemonic2 = strings.TrimSpace(mnemonic2)
		if mnemonic1 == mnemonic2 {
			assistantPageComplete(assistant, pageWidget, true)
		} else {
			assistantPageComplete(assistant, pageWidget, false)
		}
	})

	pageSeedConfirmName := "page_seed_confirm"
	pageSeedConfirmTitle := "Confirm Seed"
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

func pageNodeType(assistant *gtk.Assistant, assistFunc assistantFunc) (
	*gtk.Widget,
	*gtk.Grid,
	*gtk.RadioButton,
	string,
) {
	var pageWidget *gtk.Widget

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	fatalErrorCheck(err)

	grid, err := gtk.GridNew()
	fatalErrorCheck(err)

	btnFullNode, err := gtk.RadioButtonNewWithLabel(nil, "Full node")
	fatalErrorCheck(err)
	btnFullNode.SetActive(true)

	btnPruneNode, err := gtk.RadioButtonNewWithLabelFromWidget(btnFullNode, "Pruned node")
	fatalErrorCheck(err)

	radioBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	fatalErrorCheck(err)

	radioBox.Add(btnFullNode)
	setMargin(btnFullNode, 6, 6, 6, 6)
	radioBox.Add(btnPruneNode)
	setMargin(btnPruneNode, 6, 10, 6, 6)

	grid.Attach(radioBox, 0, 0, 1, 1)

	vbox.PackStart(grid, true, true, 0)

	pageName := "page_node_type"
	pageTitle := "Node Type"
	pageSubject := "How do you want to start your node?"
	pageDesc := `A pruned node doesn’t keep all the historical data.
Instead, it only retains the most recent part of the blockchain, deleting older data to save disk space.
Offline data is available at: <a href="https://snapshot.pactus.org/">https://snapshot.pactus.org/</a>.`

	// Create and return the page widget using assistFunc
	pageWidget = assistFunc(
		assistant,
		vbox,
		pageName,
		pageTitle,
		pageSubject,
		pageDesc,
	)

	return pageWidget, grid, btnPruneNode, pageName
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
			assistantPageComplete(assistant, pageWidget, true)
		} else {
			assistantPageComplete(assistant, pageWidget, false)
		}
	}
	entryPassword.Connect("changed", func(_ *gtk.Entry) {
		validatePassword()
	})

	entryConfirmPassword.Connect("changed", func(_ *gtk.Entry) {
		validatePassword()
	})

	pagePasswordName := "page_password"
	pagePasswordTitle := "Wallet Password"
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

func pageNumValidators(assistant *gtk.Assistant,
	assistFunc assistantFunc,
) (*gtk.Widget, *gtk.ListStore, *gtk.ComboBox, string) {
	var pageWidget *gtk.Widget
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
	pageNumValidatorsTitle := "Number of Validators"
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

func pageSummary(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	var pageWidget *gtk.Widget
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

	pageFinalName := "page_summary"
	pageFinalTitle := "Summary"
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

func assistantPageComplete(assistant *gtk.Assistant, page gtk.IWidget, completed bool) {
	assistant.SetPageComplete(page, completed)
	assistant.UpdateButtonsState()
}

func getMetadata(
	ctx context.Context,
	importer *cmd.Importer,
	listBox *gtk.ListBox,
) <-chan []cmd.Metadata {
	mdCh := make(chan []cmd.Metadata, 1)

	go func() {
		defer close(mdCh)

		children := listBox.GetChildren()
		for children.Length() > 0 {
			child := children.Data().(*gtk.Widget)
			listBox.Remove(child)
			children = children.Next()
		}

		metadata, err := importer.GetMetadata(ctx)
		if err != nil {
			mdCh <- nil

			return
		}

		for _, md := range metadata {
			label, err := gtk.LabelNew(fmt.Sprintf("snapshot %s (%s)",
				md.CreatedAtTime().Format("2006-01-02"),
				util.FormatBytesToHumanReadable(md.Data.Size),
			))
			fatalErrorCheck(err)

			listBoxRow, err := gtk.ListBoxRowNew()
			fatalErrorCheck(err)

			listBoxRow.Add(label)
			listBox.Add(listBoxRow)
		}
		listBox.ShowAll()
		mdCh <- metadata
	}()

	return mdCh
}
