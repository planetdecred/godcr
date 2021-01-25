package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageCreateRestore = "CreateRestore"

type (
	seedEditors struct {
		focusIndex int
		editors    []decredmaterial.Editor
	}
)

type createRestore struct {
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
	unlock             decredmaterial.Button
	showPasswordModal  decredmaterial.Button
	hidePasswordModal  decredmaterial.Button
	showRestoreWallet  decredmaterial.Button
	showResetModal     decredmaterial.Button
	resetSeedFields    decredmaterial.Button
	hideResetModal     decredmaterial.Button

	spendingPassword      decredmaterial.Editor
	walletName            decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	addWallet             decredmaterial.Button
	errLabel              decredmaterial.Label

	seedEditors seedEditors

	seedListLeft     *layout.List
	seedListRight    *layout.List
	autoCompleteList *layout.List

	seedSuggestions []decredmaterial.Button

	createModal     *decredmaterial.Modal
	warningModal    *decredmaterial.Modal
	modalTitleLabel decredmaterial.Label
	modalSeparator  *decredmaterial.Line
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage(common pageCommon) layout.Widget {
	pg := createRestore{
		theme:    common.theme,
		wal:      common.wallet,
		info:     common.info,
		keyEvent: common.keyEvents,
		errChan:  common.errorChannels[PageCreateRestore],

		errLabel:              common.theme.Body1(""),
		spendingPassword:      common.theme.Editor(new(widget.Editor), "Enter password"),
		walletName:            common.theme.Editor(new(widget.Editor), "Enter wallet name"),
		matchSpendingPassword: common.theme.Editor(new(widget.Editor), "Enter password again"),
		addWallet:             common.theme.Button(new(widget.Clickable), "create wallet"),
		hideResetModal:        common.theme.Button(new(widget.Clickable), "cancel"),
		suggestionLimit:       3,
		createModal:           common.theme.Modal(),
		warningModal:          common.theme.Modal(),
		modalTitleLabel:       common.theme.H6(""),
		modalSeparator:        common.theme.Line(),
	}

	pg.create = common.theme.Button(new(widget.Clickable), "create wallet")
	pg.unlock = common.theme.Button(new(widget.Clickable), "unlock wallet")
	pg.unlock.Background = common.theme.Color.Success

	pg.showPasswordModal = common.theme.Button(new(widget.Clickable), "proceed")
	pg.showRestoreWallet = common.theme.Button(new(widget.Clickable), "Restore an existing wallet")
	pg.showRestoreWallet.Background = color.NRGBA{}
	pg.showRestoreWallet.Color = common.theme.Color.Hint

	pg.closeCreateRestore = common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationArrowBack)))
	pg.closeCreateRestore.Background = color.NRGBA{}
	pg.closeCreateRestore.Color = common.theme.Color.Hint

	pg.hideRestoreWallet = common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationArrowBack)))
	pg.hideRestoreWallet.Background = color.NRGBA{}
	pg.hideRestoreWallet.Color = common.theme.Color.Hint

	pg.hidePasswordModal = common.theme.Button(new(widget.Clickable), "cancel")
	pg.hidePasswordModal.Color = common.theme.Color.Danger
	pg.hidePasswordModal.Background = color.NRGBA{R: 238, G: 238, B: 238, A: 255}

	pg.showResetModal = common.theme.Button(new(widget.Clickable), "reset")
	pg.showResetModal.Color = common.theme.Color.Hint
	pg.showResetModal.Background = color.NRGBA{}

	pg.resetSeedFields = common.theme.Button(new(widget.Clickable), "yes, reset")
	pg.resetSeedFields.Color = common.theme.Color.Danger
	pg.resetSeedFields.Background = color.NRGBA{R: 238, G: 238, B: 238, A: 255}

	pg.errLabel.Color = pg.theme.Color.Danger

	for i := 0; i <= 32; i++ {
		// pg.seedEditors = append(pg.seedEditors, common.theme.Editor(new(widget.Editor), fmt.Sprintf("%d", i+1)))
		widgetEditor := new(widget.Editor)
		widgetEditor.SingleLine, widgetEditor.Submit = true, true
		pg.seedEditors.editors = append(pg.seedEditors.editors, win.theme.Editor(widgetEditor, fmt.Sprintf("%d", i+1)))
	}
	pg.seedEditors.focusIndex = -1

	// init suggestion buttons
	for i := 0; i < pg.suggestionLimit; i++ {
		pg.seedSuggestions = append(pg.seedSuggestions, win.theme.Button(new(widget.Clickable), ""))
	}

	pg.seedListLeft, pg.seedListRight = &layout.List{Axis: layout.Vertical}, &layout.List{Axis: layout.Vertical}
	pg.spendingPassword.Editor.Mask, pg.matchSpendingPassword.Editor.Mask = '*', '*'
	pg.spendingPassword.Editor.SingleLine, pg.matchSpendingPassword.Editor.SingleLine = true, true
	pg.walletName.Editor.SingleLine = true

	pg.autoCompleteList = &layout.List{Axis: layout.Horizontal}

	pg.allSuggestions = dcrlibwallet.PGPWordList()
	pg.lastOffsetRight = pg.seedListRight.Position.Offset
	pg.lastOffsetLeft = pg.seedListLeft.Position.Offset

	return func(gtx C) D {
		pg.handle(common)
		return pg.layout(gtx, common)
	}
}

func (pg *createRestore) layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	return common.Layout(gtx, func(gtx C) D {
		toMax(gtx)
		pd := values.MarginPadding15
		dims := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							if common.states.creating {
								return pg.processing(gtx)
							} else if pg.showRestore {
								return pg.restore(gtx)
							} else {
								return pg.mainContent(gtx)
							}
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.showPassword {
					pg.modalTitleLabel.Text = "Create Wallet"
					if pg.showRestore {
						pg.modalTitleLabel.Text = "Restore Wallet"
					}

					w := []func(gtx C) D{
						func(gtx C) D {
							return pg.modalTitleLabel.Layout(gtx)
						},
						func(gtx C) D {
							pg.modalSeparator.Width = gtx.Constraints.Max.X
							return pg.modalSeparator.Layout(gtx)
						},
						func(gtx C) D {
							if pg.showRestore {
								return layout.Dimensions{}
							}
							return pg.walletName.Layout(gtx)
						},
						func(gtx C) D {
							return pg.spendingPassword.Layout(gtx)
						},
						func(gtx C) D {
							return pg.matchSpendingPassword.Layout(gtx)
						},
						func(gtx C) D {
							return pg.errLabel.Layout(gtx)
						},
						func(gtx C) D {
							if pg.showRestore {
								pg.addWallet.Text = "restore wallet"
							} else {
								pg.addWallet.Text = "create new wallet"
							}
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
										return pg.addWallet.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
										pg.hidePasswordModal.Color = common.theme.Color.Primary
										return pg.hidePasswordModal.Layout(gtx)
									})
								}),
							)
						},
					}
					return pg.createModal.Layout(gtx, w, 1300)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx C) D {
				if pg.showWarning {
					// pg.warningModal.SetTitle("Reset Seed Input")
					var msg = "You are about clearing all the seed input fields. Are you sure you want to proceed with this action?"
					w := []func(gtx C) D{
						func(gtx C) D {
							txt := common.theme.H6(msg)
							txt.Color = common.theme.Color.Danger
							txt.Alignment = text.Middle
							return txt.Layout(gtx)
						},
						func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
											return pg.resetSeedFields.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
											pg.hidePasswordModal.Background = common.theme.Color.Primary
											pg.hidePasswordModal.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
											return pg.hidePasswordModal.Layout(gtx)
										})
									}),
								)
							})
						},
					}
					return pg.warningModal.Layout(gtx, w, 1300)
				}
				return layout.Dimensions{}
			}),
		)
		return dims
	})
}

func (pg *createRestore) mainContent(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				if pg.info.LoadedWallets > 0 {
					return pg.closeCreateRestore.Layout(gtx)
				}
				return layout.Dimensions{}
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				title := pg.theme.H3("")
				title.Alignment = text.Middle
				if pg.info.LoadedWallets > 0 {
					title.Text = "Create or Restore Wallet"
				} else {
					title.Text = "Welcome to Decred Wallet, a secure & open-source desktop wallet."
				}
				return pg.centralize(gtx, func(gtx C) D {
					return title.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			btnPadding := values.MarginPadding10
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.wal.LoadedWalletsCount() > int32(0) {
						return layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.unlock.Layout(gtx)
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return pg.create.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return pg.showRestoreWallet.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *createRestore) restore(gtx layout.Context) layout.Dimensions {
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return pg.hideRestoreWallet.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			txt := pg.theme.H3("Restore from seed phrase")
			txt.Alignment = text.Middle
			return pg.centralize(gtx, func(gtx C) D {
				return txt.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			txt := pg.theme.H6("Enter your seed phrase in the correct order")
			txt.Alignment = text.Middle
			return pg.centralize(gtx, func(gtx C) D {
				return txt.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return pg.centralize(gtx, func(gtx C) D {
					return pg.errLabel.Layout(gtx)
				})
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Constraints.Max.X / 2
						return pg.inputsGroup(gtx, pg.seedListLeft, 16, 0)
					}),
					layout.Rigid(func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedListRight, 17, 16)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return pg.centralize(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding15,
							Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.showPasswordModal.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.showResetModal.Layout(gtx)
					}),
				)
			})
		}),
	)
	return dims
}

func (pg *createRestore) processing(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				message := pg.theme.H3("")
				message.Alignment = text.Middle
				if pg.restoring {
					message.Text = "restoring wallet..."
				} else {
					message.Text = "creating wallet..."
				}
				return message.Layout(gtx)
			})
		}))
}

func (pg *createRestore) inputsGroup(gtx layout.Context, l *layout.List, len int, startIndex int) layout.Dimensions {
	return l.Layout(gtx, len, func(gtx C, i int) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding20, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
							return pg.seedEditors.editors[i+startIndex].Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5, Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
					return pg.autoComplete(gtx, i, startIndex)
				})
			}),
		)
	})
}

func (pg *createRestore) autoComplete(gtx layout.Context, index, startIndex int) layout.Dimensions {
	if index+startIndex != pg.seedEditors.focusIndex {
		return layout.Dimensions{}
	}

	return pg.autoCompleteList.Layout(gtx, len(pg.suggestions), func(gtx C, i int) D {
		return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return pg.seedSuggestions[i].Layout(gtx)
		})
	})
}

func (pg *createRestore) onSuggestionSeedsClicked() {
	index := pg.seedEditors.focusIndex
	for _, b := range pg.seedSuggestions {
		for b.Button.Clicked() {
			pg.seedEditors.editors[index].Editor.SetText(b.Text)
			//pg.seedEditors.editors[index].Editor.Move(len(b.Text))
			pg.seedClicked = true
			if index != 32 {
				pg.seedEditors.editors[index+1].Editor.Focus()
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
			pg.seedEditors.focusIndex = -1
			pg.lastOffsetLeft = leftOffset
		}
	}
	if rightOffset > pg.lastOffsetRight || rightOffset < pg.lastOffsetRight {
		if pg.seedListRight.Position.BeforeEnd {
			pg.seedEditors.focusIndex = -1
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

	for i := 0; i < len(pg.seedEditors.editors); i++ {
		editor := &pg.seedEditors.editors[i]

		if editor.Editor.Focused() {
			focused = append(focused, i)
		}

		for _, e := range editor.Editor.Events() {
			switch e.(type) {
			case widget.ChangeEvent:
				pg.scrollUp()
				// hide suggestions if seed clicked
				if pg.seedClicked {
					pg.seedEditors.focusIndex = -1
					pg.seedClicked = false
				} else {
					pg.seedEditors.focusIndex = i
				}

				pg.suggestions = pg.suggestionSeeds(editor.Editor.Text())
				for k, s := range pg.suggestions {
					pg.seedSuggestions[k].Text = s
				}
			case widget.SubmitEvent:
				if i != 32 {
					pg.seedEditors.editors[i+1].Editor.Focus()
				}
			}
		}
	}

	if len(diff(pg.focused, focused)) > 0 {
		pg.seedEditors.focusIndex = -1
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
	pass := pg.spendingPassword.Editor.Text()
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

	match := pg.matchSpendingPassword.Editor.Text()
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
	pg.spendingPassword.Editor.SetText("")
	pg.matchSpendingPassword.Editor.SetText("")
}

func (pg *createRestore) validateSeeds() bool {
	pg.seedPhrase = ""
	pg.errLabel.Text = ""

	for i, editor := range pg.seedEditors.editors {
		if editor.Editor.Text() == "" {
			pg.seedEditors.editors[i].HintColor = pg.theme.Color.Danger
			pg.errLabel.Text = "all seed fields are required"
			return false
		}

		pg.seedPhrase += editor.Editor.Text() + " "
	}

	if !dcrlibwallet.VerifySeed(pg.seedPhrase) {
		pg.errLabel.Text = "invalid seed phrase"
		return false
	}

	return true
}

func (pg *createRestore) resetSeeds() {
	for i := 0; i < len(pg.seedEditors.editors); i++ {
		pg.seedEditors.editors[i].Editor.SetText("")
	}
}

func (pg *createRestore) resetPage() {
	pg.showPassword = false
	pg.showRestore = false
}

func (pg *createRestore) centralize(gtx layout.Context, content layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, content)
		}),
	)
}

func (pg *createRestore) handle(common pageCommon) {
	for pg.hideRestoreWallet.Button.Clicked() {
		pg.showRestore = false
		pg.restoring = false
		pg.errLabel.Text = ""
	}

	for pg.showRestoreWallet.Button.Clicked() {
		pg.restoring = true
		pg.showRestore = true
	}

	for pg.closeCreateRestore.Button.Clicked() {
		pg.resetSeeds()
		*common.page = PageWallet
	}

	for pg.unlock.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: UnlockWalletTemplate,
				title:    "Enter startup wallet password",
				confirm: func(pass string) {
					pg.wal.OpenWallets(pass, pg.errChan)
				},
				confirmText: "Confirm",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	for pg.create.Button.Clicked() {
		// pg.showPassword = true
		go func() {
			common.modalReceiver <- &modalLoad{
				template: CreateWalletTemplate,
				title:    "Create new wallet",
				confirm: func(wallet, pass string) {
					pg.wal.CreateWallet(wallet, pass, pg.errChan)
				},
				confirmText: "Create",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	for pg.showPasswordModal.Button.Clicked() {
		if pg.showRestore {
			if !pg.validateSeeds() {
				return
			}
		}
		pg.showPassword = true
	}

	for pg.hidePasswordModal.Button.Clicked() {
		pg.showPassword = false
		pg.showWarning = false
		pg.errLabel.Text = ""
		pg.resetPasswords()
	}

	for pg.showResetModal.Button.Clicked() {
		pg.showWarning = true
	}

	for pg.resetSeedFields.Button.Clicked() {
		pg.resetSeeds()
		pg.seedEditors.focusIndex = -1
		pg.showWarning = false
	}

	if pg.addWallet.Button.Clicked() {
		pass := pg.validatePasswords()
		if pass == "" {
			return
		}

		if pg.showRestore {
			pg.wal.RestoreWallet(pg.seedPhrase, pass, pg.errChan)
			pg.resetSeeds()
		} else {
			pg.wal.CreateWallet(pg.walletName.Editor.Text(), pass, pg.errChan)
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
				focus := pg.seedEditors.focusIndex
				pg.seedEditors.editors[focus].Editor.SetText(pg.suggestions[0])
				pg.seedClicked = true
				//pg.seedEditors.editors[focus].Editor.Move(len(pg.suggestions[0]))
			}
		}
	case err := <-pg.errChan:
		common.states.creating = false
		errText := err.Error()
		if err.Error() == "exists" {
			errText = "Wallet name already exists"
		}
		common.Notify(errText, false)
	default:
	}

	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}
