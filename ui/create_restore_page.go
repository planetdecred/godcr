package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/io/key"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial/editor"
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
		button *widget.Button
		skin   decredmaterial.Button
	}
)

type createRestore struct {
	gtx             *layout.Context
	theme           *decredmaterial.Theme
	info            *wallet.MultiWalletInfo
	wal             *wallet.Wallet
	keyEvent        chan *key.Event
	errChan         chan error
	showRestore     bool
	restoring       bool
	showPassword    bool
	showWarning     bool
	seedPhrase      string
	suggestionLimit int
	suggestions     []string
	allSuggestions  []string
	seedClicked     bool
	lastOffsetRight int
	lastOffsetLeft  int
	focused         []int

	closeCreateRestore decredmaterial.IconButton
	hideRestoreWallet  decredmaterial.IconButton
	create             decredmaterial.Button
	showPasswordModal  decredmaterial.Button
	hidePasswordModal  decredmaterial.Button
	showRestoreWallet  decredmaterial.Button
	showReset          decredmaterial.Button

	seedEditors           []decredmaterial.Editor
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	addWallet             decredmaterial.Button
	errLabel              decredmaterial.Label

	seedEditorWidgets           seedEditors
	toCreateWalletWidget        *widget.Button
	showPasswordModalWidget     *widget.Button
	hidePasswordModalWidget     *widget.Button
	backCreateRestoreWidget     *widget.Button
	showRestoreWidget           *widget.Button
	hideRestoreWidget           *widget.Button
	spendingPasswordWidget      *editor.Editor
	matchSpendingPasswordWidget *editor.Editor
	addWalletWidget             *widget.Button
	showResetWidget             *widget.Button

	seedListLeft     *layout.List
	seedListRight    *layout.List
	autoCompleteList *layout.List

	seedSuggestions []seedSuggestions
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
		spendingPasswordWidget:      new(editor.Editor),
		matchSpendingPasswordWidget: new(editor.Editor),
		addWalletWidget:             new(widget.Button),
		showResetWidget:             new(widget.Button),

		errLabel:              common.theme.Body1(""),
		spendingPassword:      common.theme.Editor("Enter password"),
		matchSpendingPassword: common.theme.Editor("Enter password again"),
		addWallet:             common.theme.Button("create wallet"),
		suggestionLimit:       3,
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

	pg.showReset = common.theme.Button("reset")
	pg.showReset.Color = common.theme.Color.Hint
	pg.showReset.Background = color.RGBA{}

	pg.errLabel.Color = pg.theme.Color.Danger

	for i := 0; i <= 32; i++ {
		pg.seedEditors = append(pg.seedEditors, common.theme.Editor(fmt.Sprintf("%d", i+1)))
		pg.seedEditorWidgets.editors = append(pg.seedEditorWidgets.editors, widget.Editor{SingleLine: true, Submit: true})
	}
	pg.seedEditorWidgets.focusIndex = -1

	// init suggestion buttons
	for i := 0; i < pg.suggestionLimit; i++ {
		pg.seedSuggestions = append(pg.seedSuggestions, seedSuggestions{
			button: new(widget.Button),
			skin:   win.theme.Button(""),
		})
	}

	pg.seedListLeft, pg.seedListRight = &layout.List{Axis: layout.Vertical}, &layout.List{Axis: layout.Vertical}
	pg.spendingPasswordWidget.Mask, pg.matchSpendingPasswordWidget.Mask = '*', '*'
	pg.spendingPasswordWidget.SingleLine, pg.matchSpendingPasswordWidget.SingleLine = true, true

	pg.autoCompleteList = &layout.List{Axis: layout.Horizontal}

	pg.allSuggestions = dcrlibwallet.PGPWordList()
	pg.lastOffsetRight = pg.seedListRight.Position.Offset
	pg.lastOffsetLeft = pg.seedListLeft.Position.Offset

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
					pg.spendingPassword.LayoutPasswordEditor(pg.gtx, pg.spendingPasswordWidget)
				},
				func() {
					pg.matchSpendingPassword.LayoutPasswordEditor(pg.gtx, pg.matchSpendingPasswordWidget)
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
								pg.hidePasswordModal.Color = common.theme.Color.Primary
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
				var msg = "You are about clearing all the seed input fields. Are you sure you want to proceed with this action?."

				if pg.showWarning {
					pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
					pg.theme.ErrorAlert(pg.gtx, msg)
				}
			}),
			layout.Rigid(func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Flex{Alignment: layout.Middle}.Layout(pg.gtx,
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Right: unit.Dp(10)}.Layout(pg.gtx, func() {
								pg.showPasswordModal.Layout(pg.gtx, pg.showPasswordModalWidget)
							})
						}),
						layout.Rigid(func() {
							pg.showReset.Layout(pg.gtx, pg.showResetWidget)
						}),
					)
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
	if index+startIndex != pg.seedEditorWidgets.focusIndex {
		return
	}

	pg.autoCompleteList.Layout(pg.gtx, len(pg.suggestions), func(i int) {
		layout.Inset{Right: unit.Dp(4)}.Layout(pg.gtx, func() {
			pg.seedSuggestions[i].skin.Layout(pg.gtx, pg.seedSuggestions[i].button)
		})
	})
}

func currentFocus(focusedList []int) int {
	f := 0
	for i, e := range focusedList {
		if i == 0 || e > f {
			f = e
		}
	}
	return f
}

func (pg *createRestore) onSuggestionSeedsClicked() {
	index := pg.seedEditorWidgets.focusIndex
	for _, b := range pg.seedSuggestions {
		for b.button.Clicked(pg.gtx) {
			pg.seedEditorWidgets.editors[index].SetText(b.skin.Text)
			pg.seedEditorWidgets.editors[index].Move(len(b.skin.Text))
			pg.seedClicked = true
			if index != 32 {
				pg.seedEditorWidgets.editors[index+1].Focus()
			}
		}
	}
}

// scrollUp scrolls up the editor list to display seed suggestions if focused editor is the last
func (pg *createRestore) scrollUp() {
	if !pg.seedListLeft.Position.BeforeEnd {
		pg.seedListLeft.Position.Offset += 100
		pg.lastOffsetLeft += 100
	}

	if !pg.seedListRight.Position.BeforeEnd {
		pg.seedListRight.Position.Offset += 100
		pg.lastOffsetRight += 100
	}
}

func (pg *createRestore) hideSuggestionsOnScroll() {
	leftOffset := pg.seedListLeft.Position.Offset
	rightOffset := pg.seedListRight.Position.Offset
	if leftOffset > pg.lastOffsetLeft || leftOffset < pg.lastOffsetLeft {
		if pg.seedListLeft.Position.BeforeEnd {
			pg.seedEditorWidgets.focusIndex = -1
			pg.lastOffsetLeft = leftOffset
		}
	}
	if rightOffset > pg.lastOffsetRight || rightOffset < pg.lastOffsetRight {
		if pg.seedListRight.Position.BeforeEnd {
			pg.seedEditorWidgets.focusIndex = -1
			pg.lastOffsetRight = rightOffset
		}
	}
}

func diff(a, b []int) []int {
	temp := map[int]int{}
	for _, s := range a {
		temp[s]++
	}
	for _, s := range b {
		temp[s]--
	}

	var f []int
	for s, v := range temp {
		if v != 0 {
			f = append(f, s)
		}
	}
	return f
}

func (pg *createRestore) editorSeedsEventsHandler() {
	var focused []int

	for i := 0; i < len(pg.seedEditorWidgets.editors); i++ {
		editor := &pg.seedEditorWidgets.editors[i]

		if editor.Focused() {
			focused = append(focused, i)
		}

		for _, e := range editor.Events(pg.gtx) {
			switch e.(type) {
			case widget.ChangeEvent:
				pg.scrollUp()
				pg.showWarning = false
				pg.resetBtn()
				pg.errLabel.Text = ""

				// hide suggestions if seed clicked
				if pg.seedClicked {
					pg.seedEditorWidgets.focusIndex = -1
					pg.seedClicked = false
				} else {
					pg.seedEditorWidgets.focusIndex = i
				}

				pg.suggestions = pg.suggestionSeeds(editor.Text())
				for k, s := range pg.suggestions {
					pg.seedSuggestions[k].skin.Text = s
				}
			case widget.SubmitEvent:
				if i != 32 {
					pg.seedEditorWidgets.editors[i+1].Focus()
				}
			}
		}
	}

	if len(diff(pg.focused, focused)) > 0 {
		pg.seedEditorWidgets.focusIndex = -1
	}
	pg.focused = focused
	pg.hideSuggestionsOnScroll()
}

func (pg createRestore) suggestionSeeds(text string) []string {
	var seeds []string
	if text == "" {
		return seeds
	}

	for _, word := range pg.allSuggestions {
		if strings.HasPrefix(strings.ToLower(word), strings.ToLower(text)) {
			if len(seeds) < pg.suggestionLimit {
				seeds = append(seeds, word)
			}
		}
	}
	return seeds
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

func (pg *createRestore) resetBtn() {
	pg.showReset.Text = "reset"
	pg.showReset.Color = pg.theme.Color.Hint
	pg.showReset.Background = color.RGBA{}
	pg.seedEditorWidgets.focusIndex = -1
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
		pg.showWarning = false
		pg.errLabel.Text = ""
		pg.resetPasswords()
	}

	for !pg.showWarning && pg.showResetWidget.Clicked(gtx) {
		pg.showWarning = true
		pg.showReset.Text = "Yes reset"
		pg.showReset.Color = color.RGBA{255, 255, 255, 255}
		pg.showReset.Background = pg.theme.Color.Danger
	}

	if pg.showWarning && pg.showResetWidget.Clicked(gtx) {
		pg.resetSeeds()
		pg.showWarning = false
		pg.resetBtn()
	}

	if pg.addWalletWidget.Clicked(gtx) {
		pass := pg.validatePasswords()
		if pass == "" {
			return
		}

		if pg.showRestore {
			pg.wal.RestoreWallet(pg.seedPhrase, pass, pg.errChan)
			pg.resetSeeds()
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
			if len(pg.suggestions) == 1 {
				focus := pg.seedEditorWidgets.focusIndex
				pg.seedEditorWidgets.editors[focus].SetText(pg.suggestions[0])
				pg.seedClicked = true
				pg.seedEditorWidgets.editors[focus].Move(len(pg.suggestions[0]))
			}
		}
	case err := <-pg.errChan:
		pg.errLabel.Text = err.Error()
	default:
	}

	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}
