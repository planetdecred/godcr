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
	gtx          *layout.Context
	theme        *decredmaterial.Theme
	info         *wallet.MultiWalletInfo
	wal          *wallet.Wallet
	keyEvent     chan *key.Event
	errChan      chan error
	showRestore  bool
	restoring    bool
	showPassword bool
	seedPhrase   string

	closeCreateRestore    decredmaterial.IconButton
	hideRestoreWallet     decredmaterial.IconButton
	create                decredmaterial.Button
	showPasswordModal     decredmaterial.Button
	hidePasswordModal     decredmaterial.Button
	showRestoreWallet     decredmaterial.Button
	seedEditors           []decredmaterial.Editor
	seedSuggestionButtons []decredmaterial.Button
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	addWallet             decredmaterial.Button
	errLabel              decredmaterial.Label

	seedEditorWidgets           seedEditors
	seedSuggestionButtonWidgets []seedSuggestions
	toCreateWalletWidget        *widget.Button
	showPasswordModalWidget     *widget.Button
	hidePasswordModalWidget     *widget.Button
	backCreateRestoreWidget     *widget.Button
	showRestoreWidget           *widget.Button
	hideRestoreWidget           *widget.Button
	spendingPasswordWidget      *widget.Editor
	matchSpendingPasswordWidget *widget.Editor
	addWalletWidget             *widget.Button

	seedListLeft  *layout.List
	seedListRight *layout.List
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage(common pageCommon) layout.Widget {
	pg := createRestore{
		gtx:                         common.gtx,
		theme:                       common.theme,
		wal:                         common.wallet,
		info:                        common.info,
		keyEvent:                    common.keyEvents,
		errChan:                     common.errorChannels[PageCreateRestore],
		toCreateWalletWidget:        new(widget.Button),
		backCreateRestoreWidget:     new(widget.Button),
		showRestoreWidget:           new(widget.Button),
		hideRestoreWidget:           new(widget.Button),
		showPasswordModalWidget:     new(widget.Button),
		hidePasswordModalWidget:     new(widget.Button),
		spendingPasswordWidget:      new(widget.Editor),
		matchSpendingPasswordWidget: new(widget.Editor),
		addWalletWidget:             new(widget.Button),

		errLabel:              common.theme.Body1(""),
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
	pg.closeCreateRestore.Color = common.theme.Color.Hint

	pg.hideRestoreWallet = common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationArrowBack)))
	pg.hideRestoreWallet.Background = color.RGBA{}
	pg.hideRestoreWallet.Color = common.theme.Color.Hint

	pg.hidePasswordModal = common.theme.Button("cancel")
	pg.hidePasswordModal.Color = common.theme.Color.Danger
	pg.hidePasswordModal.Background = color.RGBA{R: 238, G: 238, B: 238, A: 255}

	pg.errLabel.Color = pg.theme.Color.Danger
	for i := 0; i <= 32; i++ {
		pg.seedEditors = append(pg.seedEditors, common.theme.Editor(fmt.Sprintf("%d", i+1)))
		pg.seedEditorWidgets.focusIndex = -1
		pg.seedEditorWidgets.editors = append(pg.seedEditorWidgets.editors, widget.Editor{SingleLine: true, Submit: true})
	}

	pg.seedListLeft, pg.seedListRight = &layout.List{Axis: layout.Vertical}, &layout.List{Axis: layout.Vertical}

	return func() {
		pg.layout(common)
		pg.handle(common)
	}
}

func (pg *createRestore) layout(common pageCommon) {
	pg.theme.Surface(pg.gtx, func() {
		toMax(pg.gtx)
		pd := unit.Dp(15)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(pg.gtx,
			layout.Flexed(1, func() {
				layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(pg.gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
						layout.Flexed(1, func() {
							if common.states.creating {
								pg.processing()()
							} else if pg.showRestore {
								pg.restore()()
							} else {
								pg.mainContent()()
							}
						}),
					)
				})
			}),
		)
		if pg.showPassword {
			modalTitle := "Create Wallet"
			if pg.showRestore {
				modalTitle = "Restore Wallet"
			}

			w := []func(){
				func() {
					pg.spendingPassword.Layout(pg.gtx, pg.spendingPasswordWidget)
				},
				func() {
					pg.matchSpendingPassword.Layout(pg.gtx, pg.matchSpendingPasswordWidget)
				},
				func() {
					pg.errLabel.Layout(pg.gtx)
				},
				func() {
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
								pg.hidePasswordModal.Layout(pg.gtx, pg.hidePasswordModalWidget)
							})
						}),
					)
				},
			}
			pg.theme.Modal(pg.gtx, modalTitle, w)
		}
	})
}

func (pg *createRestore) mainContent() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.W.Layout(pg.gtx, func() {
					if pg.info.LoadedWallets > 0 {
						pg.closeCreateRestore.Layout(pg.gtx, pg.backCreateRestoreWidget)
					}
				})
			}),
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					title := pg.theme.H3("")
					title.Alignment = text.Middle
					if pg.info.LoadedWallets > 0 {
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
							pg.showRestoreWallet.Layout(pg.gtx, pg.showRestoreWidget)
						})
					}),
				)
			}),
		)
	}
}

func (pg *createRestore) restore() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.W.Layout(pg.gtx, func() {
					pg.hideRestoreWallet.Layout(pg.gtx, pg.hideRestoreWidget)
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
						pg.errLabel.Layout(pg.gtx)
					})
				})
			}),
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Flex{}.Layout(pg.gtx,
						layout.Rigid(func() {
							pg.gtx.Constraints.Width.Max = pg.gtx.Constraints.Width.Max / 2
							pg.inputsGroup(pg.seedListLeft, 16, 0)
						}),
						layout.Rigid(func() {
							pg.inputsGroup(pg.seedListRight, 17, 16)
						}),
					)
				})
			}),
			layout.Rigid(func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15)}.Layout(pg.gtx, func() {
						pg.showPasswordModal.Layout(pg.gtx, pg.showPasswordModalWidget)
					})
				})
			}),
		)
	}
}

func (pg *createRestore) processing() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					message := pg.theme.H3("")
					message.Alignment = text.Middle
					if pg.restoring {
						message.Text = "restoring wallet..."
					} else {
						message.Text = "creating wallet..."
					}
					message.Layout(pg.gtx)
				})
			}))
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
				layout.Inset{Top: unit.Dp(2), Left: unit.Dp(20)}.Layout(pg.gtx, func() {
					pg.autoComplete(i, startIndex)
				})
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
			pg.initSeedSuggestionButtons(editor)
			return
		}

		for _, e := range editor.Events(pg.gtx) {
			switch e.(type) {
			case widget.ChangeEvent:
				pg.seedSuggestionButtonWidgets = nil
				pg.seedSuggestionButtons = nil
				pg.initSeedSuggestionButtons(editor)

			case widget.SubmitEvent:
				if i < len(pg.seedEditorWidgets.editors)-1 {
					pg.seedEditorWidgets.editors[i+1].Focus()
				}
			}
		}
	}
}

func (pg *createRestore) initSeedSuggestionButtons(editor *widget.Editor) {
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
}

func (pg *createRestore) validatePassword() string {
	pass := pg.spendingPasswordWidget.Text()
	if pass == "" {
		pg.spendingPassword.HintColor = pg.theme.Color.Danger
		pg.errLabel.Text = fmt.Sprintf("wallet password required and cannot be empty")
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
		pg.errLabel.Text = fmt.Sprintf("enter new wallet password again and it cannot be empty")
		return ""
	}

	if match != pass {
		pg.errLabel.Text = fmt.Sprintf("new wallet passwords does not match")
		return ""
	}

	return pass
}

func (pg *createRestore) resetPasswords() {
	pg.spendingPasswordWidget.SetText("")
	pg.matchSpendingPasswordWidget.SetText("")
}

func (pg *createRestore) validateSeeds() bool {
	pg.seedPhrase = ""
	pg.errLabel.Text = ""

	for i, editor := range pg.seedEditorWidgets.editors {
		if editor.Text() == "" {
			pg.seedEditors[i].HintColor = pg.theme.Color.Danger
			pg.errLabel.Text = "all seed fields are required"
			return false
		}

		pg.seedPhrase += editor.Text() + " "
	}

	if !dcrlibwallet.VerifySeed(pg.seedPhrase) {
		pg.errLabel.Text = "invalid seed phrase"
		return false
	}

	return true
}

func (pg *createRestore) resetSeeds() {
	for i := 0; i < len(pg.seedEditorWidgets.editors); i++ {
		pg.seedEditorWidgets.editors[i].SetText("")
	}
}

func (pg *createRestore) resetPage() {
	pg.showPassword = false
	pg.showRestore = false
}

func (pg *createRestore) handle(common pageCommon) {
	gtx := common.gtx

	for pg.hideRestoreWidget.Clicked(gtx) {
		pg.showRestore = false
		pg.restoring = false
		pg.errLabel.Text = ""
	}

	for pg.showRestoreWidget.Clicked(gtx) {
		pg.restoring = true
		pg.showRestore = true
	}

	for pg.backCreateRestoreWidget.Clicked(gtx) {
		pg.resetSeeds()
		*common.page = PageWallet
	}

	for pg.toCreateWalletWidget.Clicked(gtx) {
		pg.showPassword = true
	}

	for pg.showPasswordModalWidget.Clicked(gtx) {
		if pg.showRestore {
			if !pg.validateSeeds() {
				return
			}
		}
		pg.showPassword = true
	}

	for pg.hidePasswordModalWidget.Clicked(gtx) {
		pg.showPassword = false
		pg.errLabel.Text = ""
		pg.resetPasswords()
	}

	if pg.addWalletWidget.Clicked(gtx) {
		pass := pg.validatePasswords()
		if pass == "" {
			return
		}

		if pg.showRestore {
			pg.wal.RestoreWallet(pg.seedPhrase, pass, pg.errChan)
		} else {
			pg.wal.CreateWallet(pass, pg.errChan)
		}
		common.states.creating = true
		pg.resetPasswords()
		pg.resetPage()
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
	case err := <-pg.errChan:
		pg.errLabel.Text = err.Error()
	default:
	}

	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}
