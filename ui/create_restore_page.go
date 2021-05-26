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
	"gioui.org/widget/material"

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
		editors    []decredmaterial.RestoreEditor
	}
)

type seedItemMenu struct {
	text   string
	button *widget.Clickable
}

type createRestore struct {
	common          pageCommon
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
	seedMenu        []seedItemMenu
	openPopupIndex  int

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
	optionsMenuCard       decredmaterial.Card

	passwordStrength decredmaterial.ProgressBarStyle

	seedEditors seedEditors

	seedList         *layout.List
	restoreContainer layout.List

	createModal     *decredmaterial.Modal
	warningModal    *decredmaterial.Modal
	modalTitleLabel decredmaterial.Label

	alertIcon *widget.Image
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage(common pageCommon) Page {
	pg := &createRestore{
		common:        common,
		theme:         common.theme,
		wal:           common.wallet,
		info:          common.info,
		keyEvent:      common.keyEvents,
		errorReceiver: make(chan error),

		errLabel:              common.theme.Body1(""),
		spendingPassword:      common.theme.EditorPassword(new(widget.Editor), "Spending password"),
		walletName:            common.theme.Editor(new(widget.Editor), "Wallet name (optional)"),
		matchSpendingPassword: common.theme.EditorPassword(new(widget.Editor), "Confirm spending password"),
		addWallet:             common.theme.Button(new(widget.Clickable), "create wallet"),
		hideResetModal:        common.theme.Button(new(widget.Clickable), "cancel"),
		suggestionLimit:       3,
		createModal:           common.theme.Modal(),
		warningModal:          common.theme.Modal(),
		modalTitleLabel:       common.theme.H6(""),
		passwordStrength:      win.theme.ProgressBar(0),
		openPopupIndex:        -1,
		restoreContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
	}

	pg.optionsMenuCard = decredmaterial.Card{Color: pg.theme.Color.Surface}
	pg.optionsMenuCard.Radius = decredmaterial.CornerRadius{NE: 5, NW: 5, SE: 5, SW: 5}

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
	pg.restoreWalletBtn.Background = common.theme.Color.InactiveGray
	pg.restoreWalletBtn.TextSize = values.TextSize16
	pg.errLabel.Color = pg.theme.Color.Danger

	pg.passwordStrength.Color = pg.theme.Color.LightGray

	for i := 0; i <= 32; i++ {
		widgetEditor := new(widget.Editor)
		widgetEditor.SingleLine, widgetEditor.Submit = true, true
		pg.seedEditors.editors = append(pg.seedEditors.editors, win.theme.RestoreEditor(widgetEditor, "", fmt.Sprintf("%d", i+1)))
	}
	pg.seedEditors.focusIndex = -1

	// init suggestion buttons
	pg.initSeedMenu()

	pg.seedList = &layout.List{Axis: layout.Vertical}
	pg.spendingPassword.Editor.SingleLine, pg.matchSpendingPassword.Editor.SingleLine = true, true
	pg.walletName.Editor.SingleLine = true

	pg.allSuggestions = dcrlibwallet.PGPWordList()

	return pg
}

func (pg *createRestore) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common

	if pg.info.LoadedWallets > 0 {
		pg.restoring = true
		pg.showRestore = true
	}

	return common.Layout(gtx, func(gtx C) D {
		pd := values.MarginPadding15
		dims := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
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
			layout.Rigid(func(gtx C) D {
				if pg.showPassword {
					pg.modalTitleLabel.Text = "Create Wallet"
					if pg.showRestore {
						pg.modalTitleLabel.Text = "Restore Wallet"
					}

					w := []layout.Widget{
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
									return layout.Inset{Top: values.MarginPadding6}.Layout(gtx, func(gtx C) D {
										return pg.hideRestoreWallet.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return pg.theme.H6("Restore wallet").Layout(gtx)
									})
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding24
					v := values.MarginPadding6
					return Container{padding: layout.Inset{Right: m, Left: m, Top: v, Bottom: m}}.Layout(gtx, func(gtx C) D {
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
						return layout.Inset{Bottom: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.restoreContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
								return layout.Inset{Bottom: values.MarginPadding4}.Layout(gtx, pageContent[i])
							})
						})
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D {
					return pg.restoreButtonSection(gtx)
				})
			})
		}),
	)
	return dims
}

func (pg *createRestore) restoreButtonSection(gtx layout.Context) layout.Dimensions {
	card := pg.theme.Card()
	card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
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
						return pg.inputsGroup(gtx, pg.seedList, 7, 0)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 7, 1)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 7, 2)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return pg.inputsGroup(gtx, pg.seedList, 6, 3)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return pg.inputsGroup(gtx, pg.seedList, 6, 4)
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
						m := values.MarginPadding10
						v := values.MarginPadding5
						return Container{padding: layout.Inset{Right: v, Left: v, Bottom: m}}.Layout(gtx, func(gtx C) D {
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
										Width:        values.MarginPadding1,
									}
									phase := pg.theme.Body2(phaseProgress)
									return border.Layout(gtx, func(gtx C) D {
										m := values.MarginPadding8
										v := values.MarginPadding5
										return Container{padding: layout.Inset{Right: m, Left: m, Top: v, Bottom: v}}.Layout(gtx, func(gtx C) D {
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

func (pg *createRestore) inputsGroup(gtx layout.Context, l *layout.List, len, startIndex int) layout.Dimensions {
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return l.Layout(gtx, len, func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
							pg.layoutSeedMenu(gtx, i*5+startIndex)
							return pg.seedEditors.editors[i*5+startIndex].Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (pg *createRestore) onSuggestionSeedsClicked() {
	index := pg.seedEditors.focusIndex
	for _, b := range pg.seedMenu {
		for b.button.Clicked() {
			pg.seedEditors.editors[index].Edit.Editor.SetText(b.text)
			pg.seedEditors.editors[index].Edit.Editor.MoveCaret(len(b.text), 0)
			pg.seedClicked = true
			if index != 32 {
				pg.seedEditors.editors[index+1].Edit.Editor.Focus()
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

		if editor.Edit.Editor.Focused() {
			focused = append(focused, i)
		}

		for _, e := range editor.Edit.Editor.Events() {
			switch e.(type) {
			case widget.ChangeEvent:
				// hide suggestions if seed clicked
				if pg.seedClicked {
					pg.seedEditors.focusIndex = -1
					pg.seedClicked = false
				} else {
					pg.seedEditors.focusIndex = i
				}
				text := editor.Edit.Editor.Text()
				if text == "" {
					pg.openPopupIndex = -1
				} else {
					pg.openPopupIndex = i
				}

				pg.resetSeedFields.Color = pg.theme.Color.Primary
				pg.suggestions = pg.suggestionSeeds(text)
				for k, s := range pg.suggestions {
					pg.seedMenu[k] = seedItemMenu{
						text:   s,
						button: new(widget.Clickable),
					}
				}
			case widget.SubmitEvent:
				if i != 32 {
					pg.seedEditors.editors[i+1].Edit.Editor.Focus()
				}
			}
		}
	}

	if len(diff(pg.focused, focused)) > 0 {
		pg.seedEditors.focusIndex = -1
	}
	pg.focused = focused
}

func (pg *createRestore) initSeedMenu() {
	for i := 0; i < pg.suggestionLimit; i++ {
		pg.seedMenu = append(pg.seedMenu, seedItemMenu{
			text:   "",
			button: new(widget.Clickable),
		})
	}
}

func (pg *createRestore) layoutSeedMenu(gtx layout.Context, optionsSeedMenuIndex int) {
	if pg.openPopupIndex != optionsSeedMenuIndex || pg.openPopupIndex != pg.seedEditors.focusIndex {
		return
	}

	inset := layout.Inset{
		Top:  values.MarginPadding35,
		Left: values.MarginPadding0,
	}

	m := op.Record(gtx.Ops)
	inset.Layout(gtx, func(gtx C) D {
		border := widget.Border{Color: pg.theme.Color.LightGray, CornerRadius: values.MarginPadding5, Width: values.MarginPadding2}
		return border.Layout(gtx, func(gtx C) D {
			return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.seedMenu), func(gtx C, i int) D {
					return material.Clickable(gtx, pg.seedMenu[i].button, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
							return pg.theme.Body2(pg.seedMenu[i].text).Layout(gtx)
						})
					})
				})
			})
		})
	})
	op.Defer(gtx.Ops, m.Stop())
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
		if editor.Edit.Editor.Text() == "" {
			pg.seedEditors.editors[i].Edit.HintColor = pg.theme.Color.Danger
			pg.errLabel.Text = "All seed fields are required"
			return false
		}

		pg.seedPhrase += editor.Edit.Editor.Text() + " "
	}

	if !dcrlibwallet.VerifySeed(pg.seedPhrase) {
		pg.errLabel.Text = "invalid seed phrase"
		return false
	}

	return true
}

func (pg *createRestore) resetSeeds() {
	for i := 0; i < len(pg.seedEditors.editors); i++ {
		pg.seedEditors.editors[i].Edit.Editor.SetText("")
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

func (pg *createRestore) handle() {
	common := pg.common

	for pg.hideRestoreWallet.Button.Clicked() {
		if pg.info.LoadedWallets <= 0 {
			pg.showRestore = false
			pg.restoring = false
			pg.errLabel.Text = ""
		} else {
			pg.resetSeeds()
			common.changePage(PageWallet)
		}
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

	if pg.restoreWalletBtn.Button.Clicked() {
		pass := pg.validatePasswords()
		if !pg.validateSeeds() || pass == "" {
			return
		}
		pg.wal.RestoreWallet(pg.seedPhrase, pass, pg.errorReceiver)
		pg.resetSeeds()
		common.states.creating = true
		pg.resetPasswords()
		pg.resetPage()
	}

	for pg.hidePasswordModal.Button.Clicked() {
		pg.showPassword = false
		pg.errLabel.Text = ""
		pg.resetPasswords()
	}

	if pg.matchSpendingPassword.Editor.Len() > 0 && pg.spendingPassword.Editor.Len() > 0 {
		pg.restoreWalletBtn.Background = pg.theme.Color.Primary
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
				pg.seedEditors.editors[focus].Edit.Editor.SetText(pg.suggestions[0])
				pg.seedClicked = true
				pg.seedEditors.editors[focus].Edit.Editor.MoveCaret(len(pg.suggestions[0]), -1)
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

	computePasswordStrength(&pg.passwordStrength, common.theme, pg.spendingPassword.Editor)
	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
}

func (pg *createRestore) onClose() {}
