package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
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
	errorReceiver   chan error
	showRestore     bool
	restoring       bool
	showPassword    bool
	seedPhrase      string
	suggestionLimit int
	suggestions     []string
	allSuggestions  []string
	seedClicked     bool
	focused         []int

	closeCreateRestore decredmaterial.IconButton
	hideRestoreWallet  decredmaterial.IconButton
	create             decredmaterial.Button
	unlock             decredmaterial.Button
	restoreWalletBtn   decredmaterial.Button
	hidePasswordModal  decredmaterial.Button
	showRestoreWallet  decredmaterial.Button
	resetSeedFields    decredmaterial.Button
	hideResetModal     decredmaterial.Button

	spendingPassword      decredmaterial.Editor
	walletName            decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	addWallet             decredmaterial.Button
	errLabel              decredmaterial.Label

	passwordStrength decredmaterial.ProgressBarStyle

	seedEditors seedEditors

	seedList         *layout.List
	autoCompleteList *layout.List

	seedSuggestions []decredmaterial.Button

	createModal     *decredmaterial.Modal
	warningModal    *decredmaterial.Modal
	modalTitleLabel decredmaterial.Label

	alertIcon *widget.Image
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage(common pageCommon) layout.Widget {
	pg := createRestore{
		theme:         common.theme,
		wal:           common.wallet,
		info:          common.info,
		keyEvent:      common.keyEvents,
		errorReceiver: make(chan error),

		errLabel:              common.theme.Body1(""),
		spendingPassword:      common.theme.EditorPassword(new(widget.Editor), "Enter password"),
		walletName:            common.theme.Editor(new(widget.Editor), "Wallet name (optional)"),
		matchSpendingPassword: common.theme.EditorPassword(new(widget.Editor), "Enter password again"),
		addWallet:             common.theme.Button(new(widget.Clickable), "create wallet"),
		hideResetModal:        common.theme.Button(new(widget.Clickable), "cancel"),
		suggestionLimit:       3,
		createModal:           common.theme.Modal(),
		warningModal:          common.theme.Modal(),
		modalTitleLabel:       common.theme.H6(""),
		passwordStrength:      win.theme.ProgressBar(0),
	}

	pg.create = common.theme.Button(new(widget.Clickable), "create wallet")
	pg.unlock = common.theme.Button(new(widget.Clickable), "unlock wallet")
	pg.unlock.Background = common.theme.Color.Success

	pg.restoreWalletBtn = common.theme.Button(new(widget.Clickable), "Restore")
	pg.showRestoreWallet = common.theme.Button(new(widget.Clickable), "Restore an existing wallet")
	pg.showRestoreWallet.Background = color.NRGBA{}
	pg.showRestoreWallet.Color = common.theme.Color.Hint

	pg.closeCreateRestore = common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationArrowBack)))
	pg.closeCreateRestore.Background = color.NRGBA{}
	pg.closeCreateRestore.Color = common.theme.Color.Hint

	pg.hideRestoreWallet = common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationClose)))
	pg.hideRestoreWallet.Background = color.NRGBA{}
	pg.hideRestoreWallet.Color = common.theme.Color.Hint

	pg.hidePasswordModal = common.theme.Button(new(widget.Clickable), "cancel")
	pg.hidePasswordModal.Color = common.theme.Color.Danger
	pg.hidePasswordModal.Background = color.NRGBA{R: 238, G: 238, B: 238, A: 255}

	pg.resetSeedFields = common.theme.Button(new(widget.Clickable), "Clear all")
	pg.resetSeedFields.Color = common.theme.Color.Hint
	pg.resetSeedFields.Background = color.NRGBA{}

	pg.alertIcon = common.icons.alertGray
	pg.alertIcon.Scale = 1.0

	pg.restoreWalletBtn.Inset = layout.Inset{
		Top:    values.MarginPadding12,
		Bottom: values.MarginPadding12,
		Right:  values.MarginPadding50,
		Left:   values.MarginPadding50,
	}
	pg.restoreWalletBtn.Background = common.theme.Color.Gray1
	pg.restoreWalletBtn.TextSize = values.TextSize16
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

	pg.seedList = &layout.List{Axis: layout.Vertical}
	pg.spendingPassword.Editor.SingleLine, pg.matchSpendingPassword.Editor.SingleLine = true, true
	pg.walletName.Editor.SingleLine = true

	pg.autoCompleteList = &layout.List{Axis: layout.Vertical}

	pg.allSuggestions = dcrlibwallet.PGPWordList()

	return func(gtx C) D {
		pg.handle(common)
		return pg.layout(gtx, common)
	}
}

func (pg *createRestore) layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	return common.Layout(gtx, func(gtx C) D {
		pd := values.MarginPadding15
		dims := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						if common.states.creating {
							return layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(gtx, func(gtx C) D {
								return pg.processing(gtx)
							})
						} else if pg.showRestore {
							return pg.restore(gtx)
						} else {
							return layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(gtx, func(gtx C) D {
								return pg.mainContent(gtx)
							})
						}
					}),
				)
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
							return pg.theme.Separator().Layout(gtx)
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
		)
		return common.UniformPadding(gtx, func(gtx C) D {
			return dims
		})
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
	op.TransformOp{}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, pg.theme.Color.LightGray)
	dims := layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return pg.hideRestoreWallet.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
										return pg.theme.H6("Restore wallet").Layout(gtx)
									})
								}),
							)
						})
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
						return layout.N.Layout(gtx, func(gtx C) D {
							pageContent := []func(gtx C) D{
								func(gtx C) D {
									return pg.restorePageSections(gtx, "Enter your seed phase", "1/3", func(gtx C) D {
										return pg.enterSeedPhase(gtx)
									})
								},
								func(gtx C) D {
									return pg.restorePageSections(gtx, "Create spending password", "2/3", func(gtx C) D {
										return pg.createPasswordPhase(gtx)
									})
								},
								func(gtx C) D {
									return pg.restorePageSections(gtx, "Chose a wallet name", "3/3", func(gtx C) D {
										return pg.renameWalletPhase(gtx)
									})
								},
							}
							return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pageContent), func(gtx C, i int) D {
								return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, pageContent[i])
							})
						})
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			card := pg.theme.Card()
			card.Radius = decredmaterial.CornerRadius{
				NE: 0,
				NW: 0,
				SE: 0,
				SW: 0,
			}
			return card.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding16,
								Bottom: values.MarginPadding16,
								Right:  values.MarginPadding16,
							}.Layout(gtx, func(gtx C) D {
								return pg.restoreWalletBtn.Layout(gtx)
							})
						})
					}),
				)
			})
		}),
	)
	return dims
}

func (pg *createRestore) enterSeedPhase(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Right: values.MarginPadding5,
			}
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 7, 0, 5)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 7, 1, 5)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 7, 2, 5)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 6, 3, 5)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return pg.inputsGroup(gtx, pg.seedList, 6, 4, 5)
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.errLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.resetSeedFields.Layout(gtx)
		}),
	)

}

func (pg *createRestore) createPasswordPhase(gtx layout.Context) layout.Dimensions {
	phaseContents := []func(gtx C) D{
		func(gtx C) D {
			card := pg.theme.Card()
			card.Color = pg.theme.Color.LightGray
			msg := "This spending password is required to sign transactions. Make sure to use a strong password and keep it safe."
			return card.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
								Top:   values.MarginPadding3,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.alertIcon.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.theme.Body1(msg).Layout(gtx)
						}),
					)
				})
			})
		},
		func(gtx C) D {
			return pg.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return pg.passwordStrength.Layout(gtx)
		},
		func(gtx C) D {
			return pg.matchSpendingPassword.Layout(gtx)
		},
	}

	return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(phaseContents), func(gtx C, i int) D {
		return layout.UniformInset(values.MarginPadding5).Layout(gtx, phaseContents[i])
	})
}

func (pg *createRestore) renameWalletPhase(gtx layout.Context) layout.Dimensions {
	return pg.walletName.Layout(gtx)
}

func (pg *createRestore) restorePageSections(gtx layout.Context, title string, phaseProgress string, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Bottom: values.MarginPadding16,
							Left:   values.MarginPadding5,
							Right:  values.MarginPadding5,
						}.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							txt := pg.theme.Body1(title)
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceBetween,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									border := widget.Border{
										Color:        pg.theme.Color.Gray1,
										CornerRadius: values.MarginPadding14,
										Width:        values.MarginPadding2,
									}
									phase := pg.theme.Body1(phaseProgress)
									return border.Layout(gtx, func(gtx C) D {
										return layout.Inset{
											Top:    values.MarginPadding5,
											Bottom: values.MarginPadding5,
											Left:   values.MarginPadding8,
											Right:  values.MarginPadding8,
										}.Layout(gtx, func(gtx C) D {
											return phase.Layout(gtx)
										})
									})
								}),
							)
						})
					}),
					layout.Rigid(body),
				)
			})
		})
	})
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

func (pg *createRestore) inputsGroup(gtx layout.Context, l *layout.List, len, startIndex, interval int) layout.Dimensions {
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return l.Layout(gtx, len, func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return pg.seedEditors.editors[i*interval+startIndex].Layout(gtx)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding5, Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
							return pg.autoComplete(gtx, i, startIndex, interval)
						})
					}),
				)
			})
		}),
	)

}

func (pg *createRestore) autoComplete(gtx layout.Context, index, startIndex, interval int) layout.Dimensions {
	if index*interval+startIndex != pg.seedEditors.focusIndex {
		return layout.Dimensions{}
	}

	return pg.autoCompleteList.Layout(gtx, len(pg.suggestions), func(gtx C, i int) D {
		return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return pg.seedSuggestions[i].Layout(gtx)
		})
	})
}

func (pg *createRestore) onSuggestionSeedsClicked() {
	index := pg.seedEditors.focusIndex
	for _, b := range pg.seedSuggestions {
		for b.Button.Clicked() {
			pg.seedEditors.editors[index].Editor.SetText(b.Text)
			pg.seedEditors.editors[index].Editor.MoveCaret(len(b.Text), 0)
			pg.seedClicked = true
			if index != 32 {
				pg.seedEditors.editors[index+1].Editor.Focus()
			}
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
				// hide suggestions if seed clicked
				if pg.seedClicked {
					pg.seedEditors.focusIndex = -1
					pg.seedClicked = false
				} else {
					pg.seedEditors.focusIndex = i
				}

				pg.resetSeedFields.Color = pg.theme.Color.Primary
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
		pg.errLabel.Text = "wallet password required and cannot be empty"
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
		pg.errLabel.Text = "Enter new wallet password again and it cannot be empty"
		return ""
	}

	if match != pass {
		pg.errLabel.Text = "Passwords does not match"
		return ""
	}

	if !pg.validateSeeds() {
		return ""
	}

	pg.restoreWalletBtn.Background = pg.theme.Color.Primary

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
			pg.errLabel.Text = "All seed fields are required"
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

func (pg *createRestore) computePasswordStrength(common pageCommon, editors ...*widget.Editor) {
	password := editors[0]
	strength := dcrlibwallet.ShannonEntropy(password.Text()) / 4.0
	pg.passwordStrength.Progress = float32(strength * 100)
	pg.passwordStrength.Color = common.theme.Color.Success
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
		common.changePage(PageWallet)
	}

	for pg.unlock.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: UnlockWalletTemplate,
				title:    "Enter startup wallet password",
				confirm: func(pass string) {
					pg.wal.OpenWallets(pass, pg.errorReceiver)
				},
				confirmText: "Confirm",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	for pg.create.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: CreateWalletTemplate,
				title:    "Create new wallet",
				confirm: func(wallet, pass string) {
					pg.wal.CreateWallet(wallet, pass, pg.errorReceiver)
				},
				confirmText: "Create",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	for pg.restoreWalletBtn.Button.Clicked() {
		if pg.showRestore {
			pass := pg.validatePasswords()
			if !pg.validateSeeds() || pass == "" {
				return
			}
		}
		pg.showPassword = true
	}

	for pg.hidePasswordModal.Button.Clicked() {
		pg.showPassword = false
		pg.errLabel.Text = ""
		pg.resetPasswords()
	}

	for pg.resetSeedFields.Button.Clicked() {
		pg.resetSeeds()
		pg.seedEditors.focusIndex = -1
	}

	if pg.addWallet.Button.Clicked() {
		pass := pg.validatePasswords()
		if pass == "" {
			return
		}

		if pg.showRestore {
			pg.wal.RestoreWallet(pg.seedPhrase, pass, pg.errorReceiver)
			pg.resetSeeds()
		} else {
			pg.wal.CreateWallet(pg.walletName.Editor.Text(), pass, pg.errorReceiver)
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
				pg.seedEditors.editors[focus].Editor.MoveCaret(len(pg.suggestions[0]), -1)
			}
		}
	case err := <-pg.errorReceiver:
		common.states.creating = false
		errText := err.Error()
		if err.Error() == "exists" {
			errText = "Wallet name already exists"
		}
		common.notify(errText, false)
	default:
	}

	pg.computePasswordStrength(common, pg.spendingPassword.Editor)
	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}
