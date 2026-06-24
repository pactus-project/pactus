//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/downloader"
	"github.com/pactus-project/pactus/wallet"
)

type assistantFunc func(assistant *gtk.Assistant, content *gtk.Widget, name,
	title, subject, desc string) *gtk.Widget

func setDefaultMargin(widget *gtk.Widget) {
	widget.SetMarginTop(4)
	widget.SetMarginBottom(4)
	widget.SetMarginStart(4)
	widget.SetMarginEnd(4)
}

// startupAssistant runs the node setup wizard.
// Returns true if setup completed successfully.
//
//nolint:all // Needs refactor.
func startupAssistant(ctx context.Context, workingDir string, chain genesis.ChainType, snapshotURL string) bool {
	successful := false
	assistant := gtk.NewAssistant()

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
	wgtNumValidators, comboNumValidators, pageNumValidatorsName := pageNumValidators(assistant, assistFunc)

	// -- page_node_type
	wgtNodeType, gridImport, radioImport, pageNodeTypeName := pageNodeType(assistant, assistFunc)

	// --- page_address_recovery
	wgtAddressRecovery, txtRecoveryLog, btnCancelRecovery, lblRecoveryStatus,
		pageAddressRecoveryName := pageAddressRecovery(assistant, assistFunc)

	// --- page_summary
	wgtSummary, txtNodeInfo, pageSummaryName := pageSummary(assistant, assistFunc)

	// Set page types
	assistant.SetPageType(wgtWalletMode, gtk.AssistantPageIntro)        // page 0
	assistant.SetPageType(wgtSeedGenerate, gtk.AssistantPageContent)    // page 1
	assistant.SetPageType(wgtSeedConfirm, gtk.AssistantPageContent)     // page 2
	assistant.SetPageType(wgtSeedRestore, gtk.AssistantPageContent)     // page 3
	assistant.SetPageType(wgtPassword, gtk.AssistantPageContent)        // page 4
	assistant.SetPageType(wgtNumValidators, gtk.AssistantPageContent)   // page 5
	assistant.SetPageType(wgtNodeType, gtk.AssistantPageContent)        // page 6
	assistant.SetPageType(wgtAddressRecovery, gtk.AssistantPageContent) // page 7
	assistant.SetPageType(wgtSummary, gtk.AssistantPageSummary)         // page 8

	mnemonic := ""
	prevPageIndex := -1
	prevPageAdjust := 0
	rewardAddr := ""
	recoveredAddrs := []string{}
	nodeCreated := false
	addressedRecovered := false
	var nodeWallet *wallet.Wallet

	// Prepare signal – handles dynamic page flow
	assistant.Connect("prepare", func(assistant *gtk.Assistant, page *gtk.Widget) {
		isRestoreMode := radioRestoreWallet.Active()
		curPageName := page.Name()
		curPageIndex := assistant.CurrentPage()

		isForward := true
		if curPageIndex > 0 && curPageIndex < prevPageIndex {
			isForward = false
		}

		gtkutil.Logf("%v (restore: %v, prev: %v, cur: %v)\n",
			curPageName, isRestoreMode, prevPageIndex, curPageIndex)

		switch curPageName {
		case pageModeName:
			assistantPageComplete(assistant, wgtWalletMode, true)

		case pageSeedGenerateName:
			if isRestoreMode {
				if isForward {
					// forward
					gtkutil.Logf("jumping forward from seedGenerate page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					gtkutil.Logf("jumping backward from seedGenerate page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}
				assistantPageComplete(assistant, wgtSeedGenerate, false)
			} else {
				mnemonic, _ = wallet.GenerateMnemonic(128)
				gtkutil.SetTextViewContent(txtSeed, mnemonic)
				assistantPageComplete(assistant, wgtSeedGenerate, true)
			}

		case pageSeedConfirmName:
			if isRestoreMode {
				if isForward {
					// forward
					gtkutil.Logf("jumping forward from seedConfirm page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					gtkutil.Logf("jumping backward from seedConfirm page")
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
					gtkutil.Logf("jumping forward from seedRestore page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					// backward
					gtkutil.Logf("jumping backward from seedRestore page")
					assistant.PreviousPage()
					prevPageAdjust = -1
				}
				assistantPageComplete(assistant, wgtSeedRestore, false)
			} else {
				assistantPageComplete(assistant, wgtSeedRestore, true)
			}

		case pagePasswordName:
			if isRestoreMode {
				mnemonic = gtkutil.GetTextViewContent(textRestoreSeed)
				if err := wallet.CheckMnemonic(mnemonic); err != nil {
					gtkutil.ShowErrorDialog(&assistant.Window,
						"Invalid seed phrase. Please check your seed phrase and try again.", nil)
					assistant.PreviousPage()
				}
			}
			assistantPageComplete(assistant, wgtPassword, true)

		case pageNumValidatorsName:
			assistantPageComplete(assistant, wgtNumValidators, true)

		case pageNodeTypeName:
			assistantPageComplete(assistant, wgtNodeType, true)

			// --- dynamic snapshot UI (only shown when pruned node is selected)
			ssLabel := gtk.NewLabel("")
			ssLabel.SetHAlign(gtk.AlignStart)
			setDefaultMargin(&ssLabel.Widget)

			listBox := gtk.NewListBox()
			listBox.SetHAlign(gtk.AlignCenter)
			listBox.SetSizeRequest(700, -1)
			setDefaultMargin(&listBox.Widget)

			ssDLBtn := gtk.NewButtonWithLabel("⏬ Download")
			ssDLBtn.SetHAlign(gtk.AlignCenter)
			ssDLBtn.SetSizeRequest(700, -1)
			setDefaultMargin(&ssDLBtn.Widget)

			ssPBLabel := gtk.NewLabel("")
			ssPBLabel.SetHAlign(gtk.AlignStart)
			ssPBLabel.SetWrap(true)
			setDefaultMargin(&ssPBLabel.Widget)

			gridImport.Attach(ssLabel, 0, 1, 1, 1)
			gridImport.Attach(listBox, 0, 2, 1, 1)
			gridImport.Attach(ssDLBtn, 0, 3, 1, 1)
			gridImport.Attach(ssPBLabel, 0, 5, 1, 1)

			ssLabel.SetVisible(false)
			listBox.SetVisible(false)
			ssDLBtn.SetVisible(false)
			ssPBLabel.SetVisible(false)

			snapshotIndex := 0

			// Create wallet and node once (if not already created)
			if !nodeCreated {
				numValidatorsStr := comboNumValidators.SelectedItem().Cast().(*gtk.StringObject).String()
				numValidators, _ := strconv.Atoi(numValidatorsStr)
				walletPassword := gtkutil.EntryGetText(entryPassword)

				var err error
				nodeWallet, rewardAddr, err = cmd.CreateNode(ctx, numValidators, chain, workingDir, mnemonic, walletPassword)
				if err != nil {
					gtkutil.ShowErrorDialog(&assistant.Window, err.Error(), nil)

					return
				}
				nodeCreated = true
			}

			// Toggle handler for "Pruned node" radio button
			radioImport.Connect("toggled", func() {
				ssLabel.SetVisible(true)
				radioImport.SetSensitive(true)
				listBox.SetSensitive(true)
				ssDLBtn.SetSensitive(true)
				listBox.SetSelectionMode(gtk.SelectionSingle)

				if radioImport.Active() {
					assistantPageComplete(assistant, wgtNodeType, false)

					ssLabel.SetText("♻️ Please wait, loading snapshot list...")

					go func() {
						time.Sleep(1 * time.Second)

						glib.IdleAdd(func() {
							storeDir := filepath.Join(workingDir, "data")
							importer, err := cmd.NewImporter(chain, snapshotURL, storeDir)
							if err != nil {
								gtkutil.ShowErrorDialog(&assistant.Window, err.Error(), nil)

								return
							}

							mdCh := getMetadata(ctx, importer, listBox)

							if metadata := <-mdCh; metadata == nil {
								gtkutil.SetColoredText(ssLabel, "❌ Failed to get snapshot list. Please try again later.", gtkutil.ColorRed)
							} else {
								ssLabel.SetText("🔽 Please select a snapshot to download:")
								listBox.SetVisible(true)

								listBox.ConnectRowSelected(func(row *gtk.ListBoxRow) {
									if row != nil && row.IsSelected() {
										snapshotIndex = row.Index()
										ssDLBtn.SetVisible(true)
									}
								})

								ssDLBtn.ConnectClicked(func() {
									radioImport.SetSensitive(false)
									ssLabel.SetSensitive(false)
									listBox.SetSensitive(false)
									ssDLBtn.SetSensitive(false)

									ssDLBtn.SetVisible(false)
									ssPBLabel.SetVisible(true)
									listBox.SetSelectionMode(gtk.SelectionNone)

									gtkutil.ClearLable(ssPBLabel)

									go func() {
										gtkutil.Logf("start downloading...\n")
										time.Sleep(1 * time.Second)

										snapshot := metadata[snapshotIndex]
										err := importer.Download(
											ctx, &snapshot,
											func(stats downloader.Stats) {
												if !stats.Completed {
													percent := int(stats.Percent)
													glib.IdleAdd(func() {
														dlMessage := fmt.Sprintf(
															"🌐 Downloading %s | %d%% (%s / %s)",
															snapshot.Data.Name, percent,
															util.FormatBytesToHumanReadable(uint64(stats.Downloaded)),
															util.FormatBytesToHumanReadable(uint64(stats.TotalSize)),
														)
														ssPBLabel.SetText(dlMessage)
													})
												}
											},
										)

										gtkutil.Logf("downloaded, error: %v", err)

										glib.IdleAdd(func() {
											if err != nil {
												gtkutil.SetColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), gtkutil.ColorRed)

												return
											}

											gtkutil.Logf("extracting data...\n")
											ssPBLabel.SetText("📂 Extracting downloaded files...")
											err = importer.ExtractAndStoreFiles()
											if err != nil {
												gtkutil.SetColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), gtkutil.ColorRed)

												return
											}

											gtkutil.Logf("moving data...\n")
											ssPBLabel.SetText("📑 Moving data...")
											err = importer.MoveStore()
											if err != nil {
												gtkutil.SetColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), gtkutil.ColorRed)

												return
											}

											gtkutil.Logf("cleanup...\n")
											err = importer.Cleanup()
											if err != nil {
												gtkutil.SetColoredText(ssPBLabel, fmt.Sprintf("❌ Import failed: %v", err), gtkutil.ColorRed)

												return
											}

											gtkutil.SetColoredText(ssPBLabel, "✅ Import completed.", gtkutil.ColorGreen)
											assistantPageComplete(assistant, wgtNodeType, true)
										})
									}()
								})
							}
						})
					}()
				} else {
					// Full node selected – no snapshot needed
					assistantPageComplete(assistant, wgtNodeType, true)
					ssLabel.SetVisible(false)
					listBox.SetVisible(false)
					ssDLBtn.SetVisible(false)
					ssPBLabel.SetVisible(false)
				}
			})

		case pageAddressRecoveryName:
			if !isRestoreMode {
				// Skip address recovery for new wallets
				if isForward {
					gtkutil.Logf("jumping forward from addressRecovery page")
					assistant.NextPage()
					prevPageAdjust = 1
				} else {
					gtkutil.Logf("jumping backward from addressRecovery page")
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
				recoveryCtx, cancelRecovery := context.WithCancel(ctx)

				// Setup cancel recovery button handler
				btnCancelRecovery.Connect("clicked", func() {
					lblRecoveryStatus.SetText("Cancelling recovery...")
					cancelRecovery()
					btnCancelRecovery.SetSensitive(false)
				})

				go func() {
					walletPassword := gtkutil.EntryGetText(entryPassword)

					recoveryIndex := 0
					err := nodeWallet.RecoveryAddresses(recoveryCtx, walletPassword, func(addr string) {
						glib.IdleAdd(func() {
							currentText := gtkutil.GetTextViewContent(txtRecoveryLog)
							newText := fmt.Sprintf("%s%d. %s\n", currentText, recoveryIndex+1, addr)
							gtkutil.SetTextViewContent(txtRecoveryLog, newText)
							recoveredAddrs = append(recoveredAddrs, addr)
							recoveryIndex++
						})
					})

					glib.IdleAdd(func() {
						if err != nil {
							if errors.Is(err, context.Canceled) {
								gtkutil.SetColoredText(lblRecoveryStatus, "Address recovery aborted", gtkutil.ColorYellow)
								btnCancelRecovery.SetVisible(false)
								assistantPageComplete(assistant, wgtAddressRecovery, true)
							} else {
								gtkutil.SetColoredText(lblRecoveryStatus, fmt.Sprintf("Address recovery failed: %v", err), gtkutil.ColorRed)
								btnCancelRecovery.SetVisible(false)
								assistantPageComplete(assistant, wgtAddressRecovery, true)
							}
						} else {
							gtkutil.SetColoredText(lblRecoveryStatus, "✅ Wallet addresses successfully recovered", gtkutil.ColorGreen)
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

			if len(recoveredAddrs) > 0 {
				nodeInfo += "🔄 Recovered Addresses:\n"
				for i, addr := range recoveredAddrs {
					nodeInfo += fmt.Sprintf("%v- %s\n", i+1, addr)
				}
			}

			nodeInfo += "\n🏛️ Validator Addresses:\n"
			for i, info := range nodeWallet.ListAddresses(wallet.OnlyValidatorAddresses()) {
				nodeInfo += fmt.Sprintf("%v- %s\n", i+1, info.Address)
			}

			nodeInfo += "\n💰 Reward Address:\n" + rewardAddr + "\n"
			nodeInfo += fmt.Sprintf("\n📁 Working Directory: %s", workingDir)
			nodeInfo += fmt.Sprintf("\n🌐 Network: %s\n", chain.String())

			gtkutil.SetTextViewContent(txtNodeInfo, nodeInfo)
		}

		prevPageIndex = curPageIndex + prevPageAdjust
	})

	// Use a channel to wait for the assistant to close
	mainLoop := glib.NewMainLoop(nil, false)
	assistant.Connect("cancel", func() {
		assistant.Destroy()
		mainLoop.Quit()
	})
	assistant.Connect("close", func() {
		assistant.Destroy()
		mainLoop.Quit()
	})
	assistant.SetModal(true)
	assistant.Present()
	mainLoop.Run()

	if nodeWallet != nil {
		nodeWallet.Close()
	}

	return successful
}

// pageAssistant is a helper that creates a standard assistant page layout.
func pageAssistant() assistantFunc {
	return func(assistant *gtk.Assistant, content *gtk.Widget, name, title, subject, desc string) *gtk.Widget {
		page := gtk.NewBox(gtk.OrientationVertical, 20)
		page.SetHExpand(true)

		frameLabel := gtk.NewLabel("")
		frameLabel.AddCSSClass("assistant-frame-label")
		frameLabel.SetMarkup(fmt.Sprintf("<b>%s</b>", subject))

		frame := gtk.NewFrame("")
		frame.SetLabelWidget(frameLabel)
		frame.AddCSSClass("assistant-frame")
		frame.SetHExpand(true)
		frame.SetChild(content)

		labelDesc := gtk.NewLabel("")
		labelDesc.AddCSSClass("assistant-frame-desc")
		labelDesc.SetUseMarkup(true)
		labelDesc.SetMarkup(desc)
		labelDesc.SetVExpand(true)
		labelDesc.SetVAlign(gtk.AlignEnd)
		labelDesc.SetHAlign(gtk.AlignStart)

		box := gtk.NewBox(gtk.OrientationVertical, 10)
		box.Append(frame)
		box.Append(labelDesc)
		page.Append(box)

		page.SetName(name)
		assistant.AppendPage(page)
		assistant.SetPageTitle(page, title)

		return &page.Widget
	}
}

func pageWalletMode(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.CheckButton, string) {
	newWalletRadio := gtk.NewCheckButtonWithLabel("Create a new wallet from scratch")
	restoreWalletRadio := gtk.NewCheckButtonWithLabel("Restore a wallet from seed phrase")
	restoreWalletRadio.SetGroup(newWalletRadio)
	newWalletRadio.SetActive(true)

	radioBox := gtk.NewBox(gtk.OrientationVertical, 0)
	radioBox.Append(newWalletRadio)
	radioBox.Append(restoreWalletRadio)

	pageName := "page_wallet_mode"
	pageTitle := "Wallet Mode"
	pageSubject := "How to create your wallet?"
	pageDesc := "If you are setting up the node for the first time, choose the first option."

	mode := assistFunc(assistant, &radioBox.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return mode, restoreWalletRadio, pageName
}

func pageSeedGenerate(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	textViewSeed := gtk.NewTextView()
	textViewSeed.SetWrapMode(gtk.WrapWord)
	textViewSeed.SetEditable(false)
	textViewSeed.SetMonospace(true)
	textViewSeed.SetSizeRequest(0, 80)

	pageName := "page_seed_generate"
	pageTitle := "Wallet Seed"
	pageSubject := "Your wallet seed phrase:"
	pageDesc := `<b>⚠️ CRITICAL: Write down this seed phrase and store it safely!</b>
This is the ONLY way to recover your wallet if needed.
Never share it with anyone or store it electronically.`

	page := assistFunc(assistant, &textViewSeed.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, textViewSeed, pageName
}

func pageSeedRestore(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	textViewRestoreSeed := gtk.NewTextView()
	textViewRestoreSeed.SetWrapMode(gtk.WrapWord)
	textViewRestoreSeed.SetEditable(true)
	textViewRestoreSeed.SetMonospace(true)
	textViewRestoreSeed.SetSizeRequest(0, 80)

	pageName := "page_seed_restore"
	pageTitle := "Wallet Seed Restore"
	pageSubject := "Enter your wallet seed phrase:"
	pageDesc := "Please enter your wallet seed phrase to restore your wallet."

	page := assistFunc(assistant, &textViewRestoreSeed.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, textViewRestoreSeed, pageName
}

func pageSeedConfirm(assistant *gtk.Assistant, assistFunc assistantFunc,
	textViewSeed *gtk.TextView,
) (*gtk.Widget, string) {
	var page *gtk.Widget
	textViewConfirmSeed := gtk.NewTextView()
	textViewConfirmSeed.SetWrapMode(gtk.WrapWord)
	textViewConfirmSeed.SetEditable(true)
	textViewConfirmSeed.SetMonospace(true)
	textViewConfirmSeed.SetSizeRequest(0, 80)

	// Disable paste
	textViewConfirmSeed.Connect("paste-clipboard", func(_ *gtk.TextView) {
		gtkutil.ShowInfoDialog(&assistant.Window, "Copy and paste is not allowed", nil)
		textViewConfirmSeed.StopEmission("paste-clipboard")
	})

	buffer := textViewConfirmSeed.Buffer()
	buffer.Connect("changed", func(_ *gtk.TextBuffer) {
		mnemonic1 := gtkutil.GetTextViewContent(textViewSeed)
		mnemonic2 := gtkutil.GetTextViewContent(textViewConfirmSeed)
		space := regexp.MustCompile(`\s+`)
		mnemonic2 = space.ReplaceAllString(mnemonic2, " ")
		mnemonic2 = strings.TrimSpace(mnemonic2)

		complete := mnemonic1 == mnemonic2
		assistantPageComplete(assistant, page, complete)
	})

	pageName := "page_seed_confirm"
	pageTitle := "Confirm Seed"
	pageSubject := "What was your seed?"
	pageDesc := `Your seed phrase is critical for wallet recovery!
To ensure you have properly saved your seed phrase, please retype it here.`

	page = assistFunc(assistant, &textViewConfirmSeed.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, pageName
}

func pageNodeType(assistant *gtk.Assistant, assistFunc assistantFunc) (
	*gtk.Widget, *gtk.Grid, *gtk.CheckButton, string,
) {
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)
	grid := gtk.NewGrid()

	btnFullNode := gtk.NewCheckButtonWithLabel("Full node")
	btnFullNode.SetActive(true)
	btnPruneNode := gtk.NewCheckButtonWithLabel("Pruned node")
	btnPruneNode.SetGroup(btnFullNode)

	radioBox := gtk.NewBox(gtk.OrientationVertical, 0)
	radioBox.Append(btnFullNode)
	radioBox.Append(btnPruneNode)

	grid.Attach(radioBox, 0, 0, 1, 1)
	vbox.Append(grid)

	pageName := "page_node_type"
	pageTitle := "Node Type"
	pageSubject := "How do you want to start your node?"
	pageDesc := `A pruned node doesn't keep all the historical blockchain data.
Instead, it only retains the most recent part of the blockchain, deleting older data to save disk space.
Snapshots are download from: <a href="https://snapshot.pactus.org/">https://snapshot.pactus.org/</a>`

	page := assistFunc(assistant, &vbox.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, grid, btnPruneNode, pageName
}

func pagePassword(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.Entry, string) {
	var page *gtk.Widget
	entryPassword := gtk.NewEntry()
	entryPassword.SetVisibility(false)
	setDefaultMargin(&entryPassword.Widget)

	labelPassword := gtk.NewLabel("Password: ")
	labelPassword.SetHAlign(gtk.AlignStart)

	entryConfirmPassword := gtk.NewEntry()
	entryConfirmPassword.SetVisibility(false)
	setDefaultMargin(&entryConfirmPassword.Widget)

	labelConfirmPassword := gtk.NewLabel("Confirmation: ")
	labelConfirmPassword.SetHAlign(gtk.AlignStart)

	grid := gtk.NewGrid()
	labelMessage := gtk.NewLabel("")
	setDefaultMargin(&labelMessage.Widget)

	grid.Attach(labelPassword, 0, 0, 1, 1)
	grid.Attach(entryPassword, 1, 0, 1, 1)
	grid.Attach(labelConfirmPassword, 0, 1, 1, 1)
	grid.Attach(entryConfirmPassword, 1, 1, 1, 1)
	grid.Attach(labelMessage, 1, 2, 1, 1)

	validatePassword := func() {
		pass1 := gtkutil.EntryGetText(entryPassword)
		pass2 := gtkutil.EntryGetText(entryConfirmPassword)

		if pass1 == pass2 {
			labelMessage.SetText("")
			assistantPageComplete(assistant, page, true)
		} else {
			if pass2 != "" {
				gtkutil.SetColoredText(labelMessage, "Passwords do not match", gtkutil.ColorYellow)
			}
			assistantPageComplete(assistant, page, false)
		}
	}

	entryPassword.Connect("changed", func(_ *gtk.Entry) { validatePassword() })
	entryConfirmPassword.Connect("changed", func(_ *gtk.Entry) { validatePassword() })

	pageName := "page_password"
	pageTitle := "Wallet Password"
	pageSubject := "Enter password for your wallet:"
	pageDesc := "Please choose a strong password to protect your wallet."

	page = assistFunc(assistant, &grid.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, entryPassword, pageName
}

func pageNumValidators(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.DropDown, string) {
	// Create a list of strings "1" to "32"
	items := make([]string, 32)
	for i := 1; i <= 32; i++ {
		items[i-1] = strconv.Itoa(i)
	}
	stringList := gtk.NewStringList(items)

	dropDown := gtk.NewDropDown(stringList, nil)
	dropDown.SetSelected(31) // default to 32 (index 31)

	label := gtk.NewLabel("Number of validators: ")
	label.SetHAlign(gtk.AlignStart)

	grid := gtk.NewGrid()
	grid.Attach(label, 0, 0, 1, 1)
	grid.Attach(dropDown, 1, 0, 1, 1)

	pageName := "page_num_validators"
	pageTitle := "Number of Validators"
	pageSubject := "How many validators do you want to create?"
	pageDesc := `Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.
You can define validators based on the amount of coins you want to stake.`

	page := assistFunc(assistant, &grid.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, dropDown, pageName
}

func pageAddressRecovery(assistant *gtk.Assistant, assistFunc assistantFunc) (
	*gtk.Widget, *gtk.TextView, *gtk.Button, *gtk.Label, string,
) {
	vbox := gtk.NewBox(gtk.OrientationVertical, 10)

	textViewRecoveryLog := gtk.NewTextView()
	textViewRecoveryLog.SetWrapMode(gtk.WrapWord)
	textViewRecoveryLog.SetEditable(false)
	textViewRecoveryLog.SetMonospace(true)

	scrolledWindow := gtk.NewScrolledWindow()
	scrolledWindow.SetSizeRequest(0, 200)
	scrolledWindow.SetChild(textViewRecoveryLog)

	lblRecoveryStatus := gtk.NewLabel("")
	lblRecoveryStatus.SetHAlign(gtk.AlignStart)

	btnCancelRecovery := gtk.NewButtonWithLabel("Cancel")
	btnCancelRecovery.SetHAlign(gtk.AlignCenter)
	btnCancelRecovery.SetSizeRequest(150, -1)

	vbox.Append(scrolledWindow)
	vbox.Append(lblRecoveryStatus)
	vbox.Append(btnCancelRecovery)

	pageName := "page_address_recovery"
	pageTitle := "Address Recovery"
	pageSubject := "Recovered Addresses"
	pageDesc := "Please wait while wallet addresses are recovered..."

	page := assistFunc(assistant, &vbox.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, textViewRecoveryLog, btnCancelRecovery, lblRecoveryStatus, pageName
}

func pageSummary(assistant *gtk.Assistant, assistFunc assistantFunc) (*gtk.Widget, *gtk.TextView, string) {
	textViewNodeInfo := gtk.NewTextView()
	textViewNodeInfo.SetWrapMode(gtk.WrapWord)
	textViewNodeInfo.SetEditable(false)
	textViewNodeInfo.SetMonospace(true)

	scrolledWindow := gtk.NewScrolledWindow()
	scrolledWindow.SetSizeRequest(0, 300)
	scrolledWindow.SetChild(textViewNodeInfo)

	pageName := "page_summary"
	pageTitle := "Summary"
	pageSubject := "Your node information:"
	pageDesc := `Congratulation. Your node is initialized successfully.
Now, you are ready to start the node!`

	page := assistFunc(assistant, &scrolledWindow.Widget, pageName, pageTitle, pageSubject, pageDesc)

	return page, textViewNodeInfo, pageName
}

func assistantPageComplete(assistant *gtk.Assistant, page *gtk.Widget, completed bool) {
	assistant.SetPageComplete(page, completed)
	assistant.UpdateButtonsState()
}

// getMetadata fetches snapshot metadata and populates the ListBox.
func getMetadata(ctx context.Context, importer *cmd.Importer, listBox *gtk.ListBox) <-chan []cmd.Metadata {
	mdCh := make(chan []cmd.Metadata, 1)

	go func() {
		defer close(mdCh)

		// Clear existing rows
		listBox.RemoveAll()

		metadata, err := importer.GetMetadata(ctx)
		if err != nil {
			mdCh <- nil

			return
		}

		for _, md := range metadata {
			label := gtk.NewLabel(fmt.Sprintf(
				"Snapshot %s (%s)",
				md.CreatedAtTime().Format("2006-01-02"),
				util.FormatBytesToHumanReadable(md.Data.Size),
			))
			row := gtk.NewListBoxRow()
			row.SetChild(label)
			listBox.Append(row)
		}
		listBox.SetVisible(true)
		mdCh <- metadata
	}()

	return mdCh
}
