package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/io/key"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageCreateRestore = "createrestore"

var (
	inputGroupContainerLeft  = &layout.List{Axis: layout.Vertical}
	inputGroupContainerRight = &layout.List{Axis: layout.Vertical}
)

type (
	seedEditors struct {
		focusIndex int
		editors    []widget.Editor
	}
	seedSuggestions struct {
		text   string
		button widget.Button
	}
)

type createRestore struct {
	gtx             *layout.Context
	theme           *decredmaterial.Theme
	wal             *wallet.Wallet
	keyEvent        chan *key.Event
	err             func()
	walletExists    bool
	showRestore     bool
	showPassword    bool

	closeCreateRestore    decredmaterial.IconButton
	backToMain            decredmaterial.IconButton
	create                decredmaterial.Button
	showPasswordModal     decredmaterial.Button
	hidePasswordModal     decredmaterial.Button
	showRestoreWallet     decredmaterial.Button
	seedEditors           []decredmaterial.Editor
	seedSuggestionButtons []decredmaterial.Button
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	addWallet             decredmaterial.Button

	seedEditorWidgets           seedEditors
	seedSuggestionButtonWidgets []seedSuggestions
	toCreateWalletWidget        *widget.Button
	togglePasswordModalWidget   *widget.Button
	backCreateRestoreWidget     *widget.Button
	toggleDisplayRestoreWidget  *widget.Button
	spendingPasswordWidget      *widget.Editor
	matchSpendingPasswordWidget *widget.Editor
	addWalletWidget             *widget.Button
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage(common pageCommon) layout.Widget {
	pg := createRestore{
		gtx:                         common.gtx,
		theme:                       common.theme,
		wal:                         common.wallet,
		keyEvent:                    common.keyEvents,
		err:                         win.Err,
		walletExists:                win.walletInfo.LoadedWallets > 0,
		toCreateWalletWidget:        new(widget.Button),
		backCreateRestoreWidget:     new(widget.Button),
		toggleDisplayRestoreWidget:  new(widget.Button),
		togglePasswordModalWidget:   new(widget.Button),
		spendingPasswordWidget:      new(widget.Editor),
		matchSpendingPasswordWidget: new(widget.Editor),
		addWalletWidget:             new(widget.Button),

		spendingPassword:      common.theme.Editor("Enter password"),
		matchSpendingPassword: common.theme.Editor("Enter password again"),
		addWallet:             common.theme.Button("create wallet"),
	}

	pg.matchSpendingPasswordWidget.SingleLine = true
	pg.create = common.theme.Button("create wallet")
	pg.showPasswordModal = common.theme.Button("proceed")
	pg.showRestoreWallet = common.theme.Button("Restore an existing wallet")
	pg.showRestoreWallet.Background = color.RGBA{}
	pg.showRestoreWallet.Color = common.theme.Color.Hint

	pg.closeCreateRestore = common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationArrowBack)))
	pg.closeCreateRestore.Background = color.RGBA{}
	pg.closeCreateRestore.Color = common.theme.Color.Primary

	pg.backToMain = common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationArrowBack)))
	pg.backToMain.Background = color.RGBA{}
	pg.backToMain.Color = common.theme.Color.Hint
	pg.backToMain.Animated = false

	pg.hidePasswordModal = common.theme.Button("cancel")
	pg.hidePasswordModal.Color = common.theme.Color.Danger
	pg.hidePasswordModal.Background = color.RGBA{R: 238, G: 238, B: 238, A: 255}

	for i := 0; i <= 32; i++ {
		pg.seedEditors = append(pg.seedEditors, common.theme.Editor(fmt.Sprintf("Input word %d...", i+1)))
		pg.seedEditorWidgets.focusIndex = -1
		pg.seedEditorWidgets.editors = append(pg.seedEditorWidgets.editors, widget.Editor{SingleLine: true, Submit: true})
	}

	return func() {
		pg.layout()
		pg.handle(common, win)
	}
}

func (pg *createRestore) layout() {
	pg.theme.Surface(pg.gtx, func() {
		toMax(pg.gtx)
		pd := unit.Dp(15)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(pg.gtx,
			layout.Flexed(1, func() {
				layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(pg.gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
						layout.Rigid(func() {
							layout.W.Layout(pg.gtx, func() {
								if pg.walletExists {
									pg.closeCreateRestore.Layout(pg.gtx, pg.backCreateRestoreWidget)
								}
							})
						}),
						layout.Flexed(1, func() {
							if pg.showRestore {
								pg.Restore()()
							} else {
								pg.mainContent()()
							}
						}),
					)
				})
			}),
		)
		if pg.showPassword {
			pg.theme.Modal(pg.gtx, func() {
				layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(pg.gtx,
					layout.Flexed(1, func() {
						layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(pg.gtx, func() {
							layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(pg.gtx,
								layout.Rigid(func() {
									d := pg.theme.H3("Create Wallet")
									d.Layout(pg.gtx)
								}),
								layout.Rigid(func() {
									pg.spendingPassword.Layout(pg.gtx, pg.spendingPasswordWidget)
								}),
								layout.Rigid(func() {
									pg.matchSpendingPassword.Layout(pg.gtx, pg.matchSpendingPasswordWidget)
								}),
								layout.Rigid(func() {
									// win.Err()
								}),
							)
						})
					}),
					layout.Rigid(func() {
						if pg.showRestore {
							pg.addWallet.Text = "restore wallet"
						} else {
							pg.addWallet.Text = "create new wallet"
						}
						layout.Flex{Axis: layout.Horizontal}.Layout(pg.gtx,
							layout.Rigid(func() {
								layout.UniformInset(unit.Dp(5)).Layout(pg.gtx, func() {
									pg.addWallet.Layout(pg.gtx, pg.addWalletWidget)
								})
							}),
							layout.Rigid(func() {
								layout.UniformInset(unit.Dp(5)).Layout(pg.gtx, func() {
									pg.hidePasswordModal.Layout(pg.gtx, pg.togglePasswordModalWidget)
								})
							}),
						)
					}),
				)
			})
		}
	})
}

func (pg *createRestore) mainContent() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					title := pg.theme.H3("")
					title.Alignment = text.Middle
					if pg.walletExists {
						title.Text = "Create or Restore Wallet"
					} else {
						title.Text = "Welcome to Decred Wallet, a secure & open-source desktop wallet."
					}
					title.Layout(pg.gtx)
				})
			}),
			layout.Rigid(func() {
				btnPadding := unit.Dp(10)
				layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
					layout.Rigid(func() {
						layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(pg.gtx, func() {
							pg.create.Layout(pg.gtx, pg.toCreateWalletWidget)
						})
					}),
					layout.Rigid(func() {
						layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(pg.gtx, func() {
							pg.showRestoreWallet.Layout(pg.gtx, pg.toggleDisplayRestoreWidget)
						})
					}),
				)
			}),
		)
	}
}

func (pg *createRestore) Restore() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.W.Layout(pg.gtx, func() {
					pg.backToMain.Layout(pg.gtx, pg.toggleDisplayRestoreWidget)
				})
			}),
			layout.Rigid(func() {
				txt := pg.theme.H3("Restore from seed phrase")
				txt.Alignment = text.Middle
				txt.Layout(pg.gtx)
			}),
			layout.Rigid(func() {
				txt := pg.theme.H6("Enter your seed phrase in the correct order")
				txt.Alignment = text.Middle
				txt.Layout(pg.gtx)
			}),
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(pg.gtx, func() {
					layout.Center.Layout(pg.gtx, func() {
						pg.err()
					})
				})
			}),
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Flex{}.Layout(pg.gtx,
						layout.Rigid(func() {
							pg.gtx.Constraints.Width.Max = pg.gtx.Constraints.Width.Max / 2
							pg.inputsGroup(inputGroupContainerLeft, 16, 0)
						}),
						layout.Rigid(func() {
							//fmt.Printf("max %v min %v \n", pg.gtx.Constraints.Width.Max, pg.gtx.Constraints.Width.Min)
							pg.inputsGroup(inputGroupContainerRight, 17, 16)
						}),
					)
				})
			}),
			layout.Rigid(func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15)}.Layout(pg.gtx, func() {
						pg.showPasswordModal.Layout(pg.gtx, pg.togglePasswordModalWidget)
					})
				})
			}),
		)
	}
}

func (pg *createRestore) inputsGroup(l *layout.List, len int, startIndex int) {
	l.Layout(pg.gtx, len, func(i int) {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(pg.gtx,
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(pg.gtx, func() {
							pg.seedEditors[i+startIndex].Layout(pg.gtx, &pg.seedEditorWidgets.editors[i+startIndex])
						})
					}),
				)
			}),
			layout.Rigid(func() {
				pg.autoComplete(i, startIndex)
			}),
		)
	})
}

func (pg *createRestore) autoComplete(index, startIndex int) {
	if !pg.seedEditorWidgets.editors[index+startIndex].Focused() {
		return
	}

	(&layout.List{Axis: layout.Horizontal}).Layout(pg.gtx, len(pg.seedSuggestionButtonWidgets), func(i int) {
		layout.Inset{Right: unit.Dp(4)}.Layout(pg.gtx, func() {
			pg.seedSuggestionButtons[i].Layout(pg.gtx, &pg.seedSuggestionButtonWidgets[i].button)
		})
	})
}

func (pg *createRestore) onSuggestionSeedsClicked() {
	for i := 0; i < len(pg.seedSuggestionButtonWidgets); i++ {
		btn := pg.seedSuggestionButtonWidgets[i]
		if btn.button.Clicked(pg.gtx) {
			for i := 0; i < len(pg.seedEditorWidgets.editors); i++ {
				editor := &pg.seedEditorWidgets.editors[i]
				if editor.Focused() {
					editor.SetText(btn.text)
					editor.Move(len(btn.text))

					if i < len(pg.seedEditorWidgets.editors)-1 {
						pg.seedEditorWidgets.editors[i+1].Focus()
					} else {
						pg.seedEditorWidgets.focusIndex = -1
					}
				}
			}
		}
	}
}

func (pg *createRestore) editorSeedsEventsHandler() {
	for i := 0; i < len(pg.seedEditorWidgets.editors); i++ {
		editor := &pg.seedEditorWidgets.editors[i]

		if editor.Focused() && pg.seedEditorWidgets.focusIndex != i {
			pg.seedSuggestionButtonWidgets = nil
			pg.seedSuggestionButtons = nil
			pg.seedEditorWidgets.focusIndex = i

			return
		}

		for _, e := range editor.Events(pg.gtx) {
			switch e.(type) {
			case widget.ChangeEvent:
				pg.seedSuggestionButtonWidgets = nil
				pg.seedSuggestionButtons = nil

				if strings.Trim(editor.Text(), " ") == "" {
					return
				}

				for _, word := range dcrlibwallet.PGPWordList() {
					if strings.HasPrefix(strings.ToLower(word), strings.ToLower(editor.Text())) {
						if len(pg.seedSuggestionButtonWidgets) < 2 {
							var btn struct {
								text   string
								button widget.Button
							}

							btn.text = word
							pg.seedSuggestionButtonWidgets = append(pg.seedSuggestionButtonWidgets, btn)
							pg.seedSuggestionButtons = append(pg.seedSuggestionButtons, pg.theme.Button(word))
						}
					}
				}

			case widget.SubmitEvent:
				if i < len(pg.seedEditorWidgets.editors)-1 {
					pg.seedEditorWidgets.editors[i+1].Focus()
				}
			}
		}
	}
}

func (pg *createRestore) validatePassword() string {
	pass := pg.spendingPasswordWidget.Text()
	if pass == "" {
		pg.spendingPassword.HintColor = pg.theme.Color.Danger
		// win.err = "Wallet password required and cannot be empty."
		return ""
	}

	return pass
}

func (pg *createRestore) validatePasswords() string {
	pass := pg.validatePassword()
	if pass == "" {
		return ""
	}

	match := pg.matchSpendingPasswordWidget.Text()
	if match == "" {
		pg.matchSpendingPassword.HintColor = pg.theme.Color.Danger
		// win.err = "Enter new wallet password again and it cannot be empty."
		fmt.Printf("Enter new wallet password again and it cannot be empty.\n")
		return ""
	}

	if match != pass {
		// win.err = "New wallet passwords does no match. Try again."
		fmt.Printf("New wallet passwords does no match. Try again. \n")
		return ""
	}

	return pass
}

func (pg *createRestore) resetPasswords() {
	pg.spendingPasswordWidget.SetText("")
	pg.matchSpendingPasswordWidget.SetText("")
}

func (pg *createRestore) validateSeeds() string {
	text := ""
	// win.err = ""

	for i, editor := range pg.seedEditorWidgets.editors {
		if editor.Text() == "" {
			pg.seedEditors[i].HintColor = pg.theme.Color.Danger
			return ""
		}

		text += editor.Text() + " "
	}
	fmt.Printf("validateSeeds %v \n", text)

	if !dcrlibwallet.VerifySeed(text) {
		fmt.Printf("Invalid Seed Error \n")
		// win.err = "Invalid seed phrase"
		return ""
	}

	return text
}

func (pg *createRestore) handle(common pageCommon, win *Window) {
	gtx := common.gtx

	for pg.toggleDisplayRestoreWidget.Clicked(gtx) {
		pg.showRestore = !pg.showRestore
	}

	for pg.backCreateRestoreWidget.Clicked(gtx) {
		fmt.Printf("clicked back button")
	}

	for pg.toCreateWalletWidget.Clicked(gtx) {
		pg.showPassword = true
	}

	for pg.togglePasswordModalWidget.Clicked(gtx) {
		pg.showPassword = !pg.showPassword
		if !pg.showPassword {
			pg.resetPasswords()
		}
	}

	if pg.addWalletWidget.Clicked(gtx) {
		pass := pg.validatePasswords()
		if pass == "" {
			return
		}

		if pg.showRestore {
			pg.wal.RestoreWallet(pg.validateSeeds(), pass)
			win.states.loading = true
			log.Debug("Restore Wallet clicked")
			return
		}
		pg.wal.CreateWallet(pass)
		pg.resetPasswords()
		log.Debug("Create Wallet clicked")
		win.states.loading = true
		return
	}

	// handle key events
	select {
	case evt := <-pg.keyEvent:
		if evt.Name == key.NameTab {
			for i := 0; i < len(pg.seedEditorWidgets.editors); i++ {
				editor := &pg.seedEditorWidgets.editors[i]
				if editor.Focused() && pg.seedSuggestionButtonWidgets != nil {
					editor.SetText(pg.seedSuggestionButtonWidgets[0].text)
					editor.Move(len(pg.seedSuggestionButtonWidgets[0].text))
				}
			}
		}
	default:
	}

	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}
