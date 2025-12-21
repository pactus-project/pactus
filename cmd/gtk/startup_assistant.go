//go:build gtk

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
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

// ---- Local shims for startup assistant only (boot flow) ----
// The rest of cmd/gtk should call gtkutil directly.

type Color = gtkutil.Color

const (
	ColorRed    = gtkutil.ColorRed
	ColorGreen  = gtkutil.ColorGreen
	ColorBlue   = gtkutil.ColorBlue
	ColorYellow = gtkutil.ColorYellow
	ColorOrange = gtkutil.ColorOrange
	ColorPurple = gtkutil.ColorPurple
	ColorGray   = gtkutil.ColorGray
)

func fatalErrorCheck(err error) { gtkutil.FatalErrorCheck(err) }

func showErrorDialog(parent gtk.IWindow, msg string) { gtkutil.ShowErrorDialog(parent, msg) }

func showInfoDialog(parent gtk.IWindow, msg string) { gtkutil.ShowInfoDialog(parent, msg) }

func showError(err error) { gtkutil.ShowError(err) }

func getTextViewContent(tv *gtk.TextView) string { return gtkutil.GetTextViewContent(tv) }

func setTextViewContent(tv *gtk.TextView, content string) { gtkutil.SetTextViewContent(tv, content) }

func getEntryText(entry *gtk.Entry) string { return gtkutil.GetEntryText(entry) }

func comboBoxActiveValue(combo *gtk.ComboBox) int { return gtkutil.ComboBoxActiveValue(combo) }

func setColoredText(label *gtk.Label, str string, color Color) {
	gtkutil.SetColoredText(label, str, color)
}

//nolint:all  // complexity can't be reduced more. It needs to refactor.
func startupAssistant(workingDir string, chainType genesis.ChainType) bool {
	successful := false
	assistant, err := gtk.AssistantNew()
	fatalErrorCheck(err)

	assistant.SetDefaultSize(600, 400)
	assistant.SetTitle("Pactus - Node Setup Wizard")

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
	wgtNumValidators, comboNumValidators,
		pageNumValidatorsName := pageNumValidators(assistant, assistFunc)

	// -- page_node_type
	wgtNodeType, gridImport, radioImport, pageNodeTypeName := pageNodeType(assistant, assistFunc)

	// --- page_address_recovery
	wgtAddressRecovery, txtRecoveryLog, btnCancelRecovery, lblRecoveryStatus,
		pageAddressRecoveryName := pageAddressRecovery(assistant, assistFunc)

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

	assistant.SetPageType(wgtWalletMode, gtk.ASSISTANT_PAGE_INTRO)        // page 0
	assistant.SetPageType(wgtSeedGenerate, gtk.ASSISTANT_PAGE_CONTENT)    // page 1
	assistant.SetPageType(wgtSeedConfirm, gtk.ASSISTANT_PAGE_CONTENT)     // page 2
	assistant.SetPageType(wgtSeedRestore, gtk.ASSISTANT_PAGE_CONTENT)     // page 3
	assistant.SetPageType(wgtPassword, gtk.ASSISTANT_PAGE_CONTENT)        // page 4
	assistant.SetPageType(wgtNumValidators, gtk.ASSISTANT_PAGE_CONTENT)   // page 5
	assistant.SetPageType(wgtNodeType, gtk.ASSISTANT_PAGE_CONTENT)        // page 6
	assistant.SetPageType(wgtAddressRecovery, gtk.ASSISTANT_PAGE_CONTENT) // page 7
	assistant.SetPageType(wgtSummary, gtk.ASSISTANT_PAGE_SUMMARY)         // page 8

	mnemonic := ""
	prevPageIndex := -1
	prevPageAdjust := 0
	rewardAddr := ""
	recoveredAddrs := []string{}
	nodeCreated := false
	addressedRecovered := false
	var nodeWallet *wallet.Wallet

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
					log.Print("jumping forward from seedGenerate page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Print("jumping backward from seedGenerate page")
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
					log.Print("jumping forward from seedConfirm page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Print("jumping backward from seedConfirm page")
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
					log.Print("jumping forward from seedRestore page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					log.Print("jumping backward from seedRestore page")
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
					showErrorDialog(assistant, "Invalid seed phrase. Please check your seed phrase and try again.")
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
			setMargin(ssLabel, 6, 6, 6, 6)
			ssLabel.SetHAlign(gtk.ALIGN_START)

			listBox, err := gtk.ListBoxNew()
			fatalErrorCheck(err)
			setMargin(listBox, 6, 6, 6, 6)
			listBox.SetHAlign(gtk.ALIGN_CENTER)
			listBox.SetSizeRequest(700, -1)

			ssDLBtn, err := gtk.ButtonNewWithLabel("⏬ Download")
			fatalErrorCheck(err)
			setMargin(ssDLBtn, 6, 6, 6, 6)
			ssDLBtn.SetHAlign(gtk.ALIGN_CENTER)
			ssDLBtn.SetSizeRequest(700, -1)

			ssPBLabel, err := gtk.LabelNew("")
			fatalErrorCheck(err)
			setMargin(ssPBLabel, 6, 6, 6, 6)
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

			if !nodeCreated {
				numValidators := comboBoxActiveValue(comboNumValidators)
				walletPassword := getEntryText(entryPassword)

				nodeWallet, rewardAddr, err = cmd.CreateNode(numValidators, chainType, workingDir, mnemonic, walletPassword)
				if err != nil {
					showError(err)

					return
				}

				// Prevent re-entry
				nodeCreated = true
			}

			radioImport.Connect("toggled", func() {
				if radioImport.GetActive() {
					assistantPageComplete(assistant, wgtNodeType, false)

					ssLabel.SetVisible(true)
					ssLabel.SetText("♻️ Please wait, loading snapshot list...")

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

							if metadata := <-mdCh; metadata == nil {
								setColoredText(ssLabel, "❌ Failed to get snapshot list. Please try again later.", ColorRed)
							} else {
								ssLabel.SetText("🔽 Please select a snapshot to download:")
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
										log.Print("start downloading...\n")

										time.Sleep(5 * time.Second)
										err := importer.Download(ctx, &metadata[snapshotIndex],
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
															ssPBLabel.SetText(dlMessage)
														})
													}
												}
											},
										)

										glib.IdleAdd(func() {
											if err != nil {
												setColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), ColorRed)

												return
											}

											log.Print("extracting data...\n")
											ssPBLabel.SetText("📂 Extracting downloaded files...")
											err := importer.ExtractAndStoreFiles()
											if err != nil {
												setColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), ColorRed)

												return
											}

											log.Print("moving data...\n")
											ssPBLabel.SetText("📑 Moving data...")
											err = importer.MoveStore()
											if err != nil {
												setColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), ColorRed)

												return
											}

											log.Print("cleanup...\n")
											err = importer.Cleanup()
											if err != nil {
												setColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), ColorRed)

												return
											}

											setColoredText(ssPBLabel, "✅ Import completed.", ColorGreen)
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
		case pageAddressRecoveryName:
			// Only handle recovery for restore mode
			if !isRestoreMode {
				// Skip this page for new wallets
				if isForward {
					log.Print("jumping forward from addressRecovery page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					log.Print("jumping backward from addressRecovery page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}

				return
			}

			if !addressedRecovered {
				// Prevent re-entry
				addressedRecovered = true

				// Disable next button initially
				assistantPageComplete(assistant, wgtAddressRecovery, false)

				lblRecoveryStatus.SetText("Processing...")

				// Reset recovery context
				recoveryCtx, cancelRecovery := context.WithCancel(context.Background())

				// Setup cancel recovery button handler
				btnCancelRecovery.Connect("clicked", func() {
					lblRecoveryStatus.SetText("Cancelling recovery...")
					cancelRecovery()
					btnCancelRecovery.SetSensitive(false)
				})

				go func() {
					walletPassword := getEntryText(entryPassword)

					recoveryIndex := 0
					err = nodeWallet.RecoveryAddresses(recoveryCtx, walletPassword, func(addr string) {
						glib.IdleAdd(func() {
							currentText := getTextViewContent(txtRecoveryLog)
							newText := fmt.Sprintf("%s%d. %s\n", currentText, recoveryIndex+1, addr)
							setTextViewContent(txtRecoveryLog, newText)
							recoveredAddrs = append(recoveredAddrs, addr)
							recoveryIndex++
						})
					})

					glib.IdleAdd(func() {
						if err != nil {
							if errors.Is(err, context.Canceled) {
								setColoredText(lblRecoveryStatus, "Address recovery aborted", ColorYellow)
								btnCancelRecovery.SetVisible(false)
								assistantPageComplete(assistant, wgtAddressRecovery, true)
							} else {
								setColoredText(lblRecoveryStatus, fmt.Sprintf("Address recovery failed: %v", err), ColorRed)
								btnCancelRecovery.SetVisible(false)
								assistantPageComplete(assistant, wgtAddressRecovery, true)
							}
						} else {
							setColoredText(lblRecoveryStatus, "✅ Wallet addresses successfully recovered", ColorGreen)
							btnCancelRecovery.SetVisible(false)
							assistantPageComplete(assistant, wgtAddressRecovery, true)
						}
					})
				}()
			}

		case pageSummaryName:
			// Done! showing the node information
			successful = true
			nodeInfo := ""

			nodeInfo += "🔄 Recovered Addresses:\n"
			for i, addr := range recoveredAddrs {
				nodeInfo += fmt.Sprintf("%v- %s\n", i+1, addr)
			}

			nodeInfo += "\n🏛️ Validator Addresses:\n"
			for i, info := range nodeWallet.ListValidatorAddresses() {
				nodeInfo += fmt.Sprintf("%v- %s\n", i+1, info.Address)
			}

			nodeInfo += "\n💰 Reward Address:\n"
			nodeInfo += fmt.Sprintf("%s\n", rewardAddr)

			nodeInfo += fmt.Sprintf("\n📁 Working Directory: %s", workingDir)
			nodeInfo += fmt.Sprintf("\n🌐 Network: %s\n", chainType.String())

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
	newWalletRadio, err := gtk.RadioButtonNewWithLabel(nil, "Create a new wallet from scratch")
	fatalErrorCheck(err)

	restoreWalletRadio, err := gtk.RadioButtonNewWithLabelFromWidget(newWalletRadio,
		"Restore a wallet from seed phrase")
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
	pageModeDesc := "If you are setting up the node for the first time, choose the first option."
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
	pageSeedSubject := "Your wallet seed phrase:"
	pageSeedDesc := `<b>⚠️ CRITICAL: Write down this seed phrase and store it safely!</b>
     This is the ONLY way to recover your wallet if needed.
     Never share it with anyone or store it electronically.`

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
	pageSeedSubject := "Enter your wallet seed phrase:"
	pageSeedDesc := "Please enter your wallet seed phrase to restore your wallet."

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
		showInfoDialog(assistant, "Copy and paste is not allowed")
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
	pageSeedConfirmDesc := `Your seed phrase is critical for wallet recovery!
To ensure you have properly saved your seed phrase, please retype it here.`

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
	pageDesc := `A pruned node doesn't keep all the historical blockchain data.
Instead, it only retains the most recent part of the blockchain, deleting older data to save disk space.
Historical data is available at: <a href="https://snapshot.pactus.org/">https://snapshot.pactus.org/</a>`

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
	labelPassword, err := gtk.LabelNew("Password: ")
	fatalErrorCheck(err)

	labelPassword.SetHAlign(gtk.ALIGN_START)
	setMargin(labelPassword, 6, 6, 6, 6)

	entryConfirmPassword, err := gtk.EntryNew()
	fatalErrorCheck(err)

	setMargin(entryConfirmPassword, 6, 6, 6, 6)
	entryConfirmPassword.SetVisibility(false)
	labelConfirmPassword, err := gtk.LabelNew("Confirmation: ")
	fatalErrorCheck(err)

	labelConfirmPassword.SetHAlign(gtk.ALIGN_START)
	setMargin(labelConfirmPassword, 6, 6, 6, 6)

	grid, err := gtk.GridNew()
	fatalErrorCheck(err)

	labelMessage, err := gtk.LabelNew("")
	fatalErrorCheck(err)

	grid.Attach(labelPassword, 0, 0, 1, 1)
	grid.Attach(entryPassword, 1, 0, 1, 1)
	grid.Attach(labelConfirmPassword, 0, 1, 1, 1)
	grid.Attach(entryConfirmPassword, 1, 1, 1, 1)
	grid.Attach(labelMessage, 1, 2, 1, 1)

	validatePassword := func() {
		pass1 := getEntryText(entryPassword)
		pass2 := getEntryText(entryConfirmPassword)

		if pass1 == pass2 {
			labelMessage.SetText("")
			assistantPageComplete(assistant, pageWidget, true)
		} else {
			if pass2 != "" {
				setColoredText(labelMessage, "Passwords do not match", ColorYellow)
			}
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
	pagePsswrdDesc := "Please choose a strong password to protect your wallet."

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
) (*gtk.Widget, *gtk.ComboBox, string) {
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

	return pageWidget, comboNumValidators, pageNumValidatorsName
}

func pageAddressRecovery(assistant *gtk.Assistant, assistFunc assistantFunc) (
	*gtk.Widget, *gtk.TextView, *gtk.Button, *gtk.Label, string,
) {
	var pageWidget *gtk.Widget

	// Create a vertical box to hold all elements
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	fatalErrorCheck(err)

	// Create TextView for recovery log (read-only, scrollable)
	textViewRecoveryLog, err := gtk.TextViewNew()
	fatalErrorCheck(err)
	setMargin(textViewRecoveryLog, 6, 6, 6, 6)
	textViewRecoveryLog.SetWrapMode(gtk.WRAP_WORD)
	textViewRecoveryLog.SetEditable(false)
	textViewRecoveryLog.SetMonospace(true)

	// Create scrolled window for the text view
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	fatalErrorCheck(err)
	scrolledWindow.SetSizeRequest(0, 200)
	scrolledWindow.Add(textViewRecoveryLog)

	// Create status label
	lblRecoveryStatus, err := gtk.LabelNew("")
	fatalErrorCheck(err)
	setMargin(lblRecoveryStatus, 6, 6, 6, 6)
	lblRecoveryStatus.SetHAlign(gtk.ALIGN_START)

	// Create cancel button
	btnCancelRecovery, err := gtk.ButtonNewWithLabel("Cancel")
	fatalErrorCheck(err)
	setMargin(btnCancelRecovery, 6, 6, 6, 6)
	btnCancelRecovery.SetHAlign(gtk.ALIGN_CENTER)
	btnCancelRecovery.SetSizeRequest(150, -1)

	// Add widgets to vbox
	vbox.PackStart(scrolledWindow, true, true, 0)
	vbox.PackStart(lblRecoveryStatus, false, false, 0)
	vbox.PackStart(btnCancelRecovery, false, false, 0)

	pageAddressRecoveryName := "page_address_recovery"
	pageAddressRecoveryTitle := "Address Recovery"
	pageAddressRecoverySubject := "Recovered Addresses"
	pageAddressRecoveryDesc := `Please wait while wallet addresses are recovered...`

	pageWidget = assistFunc(
		assistant,
		vbox,
		pageAddressRecoveryName,
		pageAddressRecoveryTitle,
		pageAddressRecoverySubject,
		pageAddressRecoveryDesc)

	return pageWidget, textViewRecoveryLog, btnCancelRecovery, lblRecoveryStatus, pageAddressRecoveryName
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
