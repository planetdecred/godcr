package ui

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/text"

	"github.com/raedahgroup/dcrlibwallet"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const (
	PageSeedBackup = "seedbackup"
	infoView       = iota
	seedView
	verifyView
	successView
)

type (
	seedGroup struct {
		selected int
		buttons  []decredmaterial.Button
	}

	viewText struct {
		title       string
		action      string
		steps       string
		instruction string
	}
)

type backupPage struct {
	theme *decredmaterial.Theme
	wal   *wallet.Wallet
	info  *wallet.MultiWalletInfo

	backButton     decredmaterial.IconButton
	title          decredmaterial.Label
	steps          decredmaterial.Label
	instruction    decredmaterial.Label
	successMessage decredmaterial.Label
	successInfo    decredmaterial.Label
	action         decredmaterial.Button
	checkBoxes     []decredmaterial.CheckBoxStyle
	checkIcon      *widget.Icon

	infoList            *layout.List
	viewList            *layout.List
	seedPhraseListLeft  *layout.List
	seedPhraseListRight *layout.List
	verifyList          *layout.List
	suggestionList      *layout.List

	suggestions []seedGroup

	seedPhrase     []string
	selectedSeeds  []string
	allSuggestions []string
	active         int
	error          string
}

func (win *Window) BackupPage(c pageCommon) layout.Widget {
	b := &backupPage{
		theme: c.theme,
		wal:   c.wallet,
		info:  c.info,

		action:         c.theme.Button(new(widget.Clickable), "View seed phrase"),
		backButton:     c.theme.PlainIconButton(new(widget.Clickable), c.icons.navigationArrowBack),
		title:          c.theme.H5("Keep in mind"),
		steps:          c.theme.Body1("Step 1/2"),
		instruction:    c.theme.H6("Write down all 33 words in the correct order"),
		successMessage: c.theme.H4("Your seed phrase backup is verified"),
		successInfo:    c.theme.Body2("Be sure to store your seed phrase backup in a secure location."),
		checkIcon:      c.icons.actionCheckCircle,

		active:        infoView,
		selectedSeeds: make([]string, 0, 33),
	}

	b.checkIcon.Color = c.theme.Color.Success
	b.steps.Color = c.theme.Color.Hint
	b.successMessage.Alignment = text.Middle
	b.successInfo.Alignment = text.Middle
	b.successInfo.Color = b.theme.Color.Hint

	b.backButton.Color = c.theme.Color.Hint
	b.backButton.Size = values.MarginPadding30

	b.action.Background = c.theme.Color.Hint

	b.checkBoxes = []decredmaterial.CheckBoxStyle{
		c.theme.CheckBox(new(widget.Bool), "The 33-word seed phrase is EXTREMELY IMPORTANT."),
		c.theme.CheckBox(new(widget.Bool), "Seed phrase is the only way to restore your wallet."),
		c.theme.CheckBox(new(widget.Bool), "It is recommended to store your seed phrase in a physical format (e.g. write down on a paper)."),
		c.theme.CheckBox(new(widget.Bool), "It is highly discouraged to store your seed phrase in any digital format (e.g. screenshot)."),
		c.theme.CheckBox(new(widget.Bool), "Anyone with your seed phrase can steal your funds. DO NOT show it to anyone."),
	}

	b.instruction.Alignment = text.Middle
	b.allSuggestions = dcrlibwallet.PGPWordList()

	for _, cb := range b.checkBoxes {
		cb.IconColor = c.theme.Color.Success
		cb.Color = c.theme.Color.Success
	}

	for i := 0; i < 33; i++ {
		var bg []decredmaterial.Button

		for j := 0; j < 3; j++ {
			bg = append(bg, c.theme.Button(new(widget.Clickable), ""))
		}
		b.suggestions = append(b.suggestions, seedGroup{selected: -1, buttons: bg})
		b.selectedSeeds = append(b.selectedSeeds, "-")
	}

	b.infoList = &layout.List{Axis: layout.Vertical}
	b.viewList = &layout.List{Axis: layout.Vertical}
	b.seedPhraseListLeft = &layout.List{Axis: layout.Vertical}
	b.seedPhraseListRight = &layout.List{Axis: layout.Vertical}
	b.verifyList = &layout.List{Axis: layout.Vertical}
	b.suggestionList = &layout.List{Axis: layout.Horizontal}

	return func(gtx C) layout.Dimensions {
		b.handle(c)
		return b.layout(gtx)
	}
}

func (pg *backupPage) activeButton() {
	pg.action.Background = pg.theme.Color.Primary
	pg.action.Color = pg.theme.Color.InvText
}

func (pg *backupPage) clearButton() {
	pg.action.Background = color.RGBA{}
	pg.action.Color = pg.theme.Color.Primary
}

func (pg *backupPage) layout(gtx layout.Context) layout.Dimensions {
	dims := pg.theme.Surface(gtx, func(ctx C) D {
		toMax(gtx)
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
			layout.Rigid(func(ctx C) D {
				pg.action.Background = pg.theme.Color.Hint
				pg.action.Color = pg.theme.Color.InvText
				switch pg.active {
				case infoView:
					if pg.verifyCheckBoxes() {
						pg.activeButton()
					}
					return pg.infoView(gtx)
				case seedView:
					pg.activeButton()
					return pg.seedView(gtx)
				case verifyView:
					if checkSlice(pg.selectedSeeds) {
						pg.activeButton()
					}
					return pg.verifyView(gtx)
				case successView:
					pg.activeButton()
					return pg.successView(gtx)
				default:
					if pg.verifyCheckBoxes() {
						pg.activeButton()
					}
					return pg.infoView(gtx)
				}
			}),
		)
	})
	return dims
}

func (pg *backupPage) pageTitle(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding5, Top: values.MarginPadding20}.Layout(gtx, func(ctx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
			layout.Rigid(func(ctx C) D {
				return pg.backButton.Layout(gtx)
			}),
			layout.Rigid(func(ctx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(ctx C) D {
						return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(ctx C) D {
							return pg.title.Layout(gtx)
						})
					}),
					layout.Rigid(func(ctx C) D {
						if pg.active != infoView {
							return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(ctx C) D {
								return pg.steps.Layout(gtx)
							})
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(func(ctx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						if pg.active != infoView {
							return layout.Inset{Right: values.MarginPadding30, Top: values.MarginPadding20}.Layout(gtx, func(ctx C) D {
								return pg.instruction.Layout(gtx)
							})
						}
						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}

func (pg *backupPage) viewTemplate(gtx layout.Context, content layout.Widget) layout.Dimensions {
	return layout.Inset{Left: values.MarginPadding10, Right: values.MarginPadding10}.Layout(gtx, func(ctx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(ctx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(ctx C) D {
						if pg.active != successView {
							return pg.pageTitle(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(func(ctx C) D {
						return layout.Inset{Bottom: values.MarginPadding50}.Layout(gtx, func(ctx C) D {
							return content(gtx)
						})
					}),
				)
			}),
			layout.Stacked(func(ctx C) D {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return layout.S.Layout(gtx, func(ctx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.Y
					return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(ctx C) D {
						return pg.action.Layout(gtx)
					})
				})
			}),
			layout.Stacked(func(ctx C) D {
				if len(pg.error) > 0 {
					return layout.Inset{Top: values.MarginPadding75}.Layout(gtx, func(ctx C) D {
						return pg.theme.ErrorAlert(gtx, pg.error)
					})
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func (pg *backupPage) infoView(gtx layout.Context) layout.Dimensions {
	return pg.viewTemplate(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		return layout.Center.Layout(gtx, func(ctx C) D {
			return layout.Inset{Bottom: values.MarginPadding60}.Layout(gtx, func(ctx C) D {
				return pg.infoList.Layout(gtx, len(pg.checkBoxes), func(gtx C, i int) D {
					return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(ctx C) D {
						return pg.checkBoxes[i].Layout(gtx)
					})
				})
			})
		})
	})
}

func (pg *backupPage) seedView(gtx layout.Context) layout.Dimensions {
	return pg.viewTemplate(gtx, func(ctx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Center.Layout(gtx, func(ctx C) D {
			return pg.viewList.Layout(gtx, 1, func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(ctx C) D {
						gtx.Constraints.Max.X = gtx.Constraints.Max.X / 2
						return pg.seedPhraseListLeft.Layout(gtx, len(pg.seedPhrase), func(gtx C, i int) D {
							if i < 17 {
								return pg.seedText(gtx, i)
							}
							return layout.Dimensions{}
						})
					}),
					layout.Rigid(func(ctx C) D {
						return pg.seedPhraseListRight.Layout(gtx, len(pg.seedPhrase), func(gtx C, i int) D {
							if i > 16 {
								return pg.seedText(gtx, i)
							}
							return layout.Dimensions{}
						})
					}),
				)
			})
		})
	},
	)
}

func (pg *backupPage) verifyView(gtx layout.Context) layout.Dimensions {
	return pg.viewTemplate(gtx, func(ctx C) D {
		toMax(gtx)
		return pg.verifyList.Layout(gtx, len(pg.suggestions), func(gtx C, i int) D {
			s := pg.suggestions[i]
			return layout.Center.Layout(gtx, func(ctx C) D {
				return layout.Inset{Bottom: values.MarginPadding30}.Layout(gtx, func(ctx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(ctx C) D {
							return layout.Inset{Left: values.MarginPadding15, Bottom: values.MarginPadding15}.Layout(gtx, func(ctx C) D {
								return pg.theme.H6(fmt.Sprintf("%d. %s", i+1, pg.selectedSeeds[i])).Layout(gtx)
							})
						}),
						layout.Rigid(func(ctx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(0.3, func(ctx C) D {
									return pg.suggestionButtonGroup(gtx, s, 0)
								}),
								layout.Flexed(0.3, func(ctx C) D {
									return pg.suggestionButtonGroup(gtx, s, 1)
								}),
								layout.Flexed(0.3, func(ctx C) D {
									return pg.suggestionButtonGroup(gtx, s, 2)
								}),
							)
						}),
					)
				})
			})
		})
	})
}

func (pg *backupPage) successView(gtx layout.Context) layout.Dimensions {
	return pg.viewTemplate(gtx, func(ctx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(ctx C) D {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return layout.Center.Layout(gtx, func(ctx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(ctx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Inset{
								Bottom: values.MarginPadding50,
								Right:  values.MarginPadding50,
							}.Layout(gtx, func(ctx C) D {
								return layout.Center.Layout(gtx, func(ctx C) D {
									return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(ctx C) D {
										pg.checkIcon.Color = pg.theme.Color.Success
										return pg.checkIcon.Layout(gtx, unit.Px(float32(150)))
									})
								})
							})
						}),
						layout.Rigid(func(ctx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.successMessage.Layout(gtx)
						}),
						layout.Rigid(func(ctx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.successInfo.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (pg *backupPage) seedText(gtx layout.Context, index int) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10, Left: values.MarginPadding20}.Layout(gtx,
		func(ctx C) D {
			seedLabel := pg.theme.H6(fmt.Sprintf("%d.  %s", index+1, pg.seedPhrase[index]))
			seedLabel.Alignment = text.Middle
			return seedLabel.Layout(gtx)
		},
	)
}

func (pg *backupPage) suggestionButtonGroup(gtx layout.Context, sg seedGroup, buttonIndex int) layout.Dimensions {
	button := sg.buttons[buttonIndex]
	button.Background = pg.theme.Color.Hint
	button.TextSize = values.TextSize18
	if sg.selected == buttonIndex {
		button.Background = pg.theme.Color.Primary
	}
	return layout.Inset{Right: values.MarginPadding15, Left: values.MarginPadding15}.Layout(gtx, func(ctx C) D {
		return button.Layout(gtx)
	})
}

func (pg *backupPage) verifyCheckBoxes() bool {
	for _, cb := range pg.checkBoxes {
		if !cb.CheckBox.Value {
			return false
		}
	}
	return true
}

func (pg *backupPage) randomSeeds() []string {
	var randomSeeds []string

	for i := 0; i < 3; i++ {
		random := rand.Intn(len(pg.allSuggestions))
		randomSeeds = append(randomSeeds, pg.allSuggestions[random])
	}
	return randomSeeds
}

func (pg *backupPage) populateSuggestionSeeds() {
	rand.Seed(time.Now().Unix())

	for k := range pg.seedPhrase {
		seeds := pg.randomSeeds()
		s := pg.suggestions[k]
		for i := range s.buttons {
			s.buttons[i].Text = seeds[i]
		}
		s.buttons[rand.Intn(len(seeds))].Text = pg.seedPhrase[k]
	}
}

func viewTexts(active int) viewText {
	switch active {
	case infoView:
		return viewText{
			title:  "Keep in mind",
			action: "View seed phrase",
		}
	case seedView:
		return viewText{
			title:       "Write down seed phrase",
			action:      "I have written down all 33 words",
			steps:       fmt.Sprintf("Steps %d/2", seedView-1),
			instruction: "Write down all 33 words in the correct order",
		}
	case verifyView:
		return viewText{
			title:       "Verify seed phrase",
			action:      "Verify",
			steps:       fmt.Sprintf("Steps %d/2", verifyView-1),
			instruction: "Select the correct words to verify",
		}
	case successView:
		return viewText{
			action: "Back to Wallets",
		}
	}
	return viewText{}
}

func (pg *backupPage) updateViewTexts() {
	t := viewTexts(pg.active)
	pg.title.Text = t.title
	pg.action.Text = t.action
	pg.steps.Text = t.steps
	pg.instruction.Text = t.instruction
}

func (pg *backupPage) clearError() {
	time.AfterFunc(time.Second*3, func() {
		pg.error = ""
	})
}

func checkSlice(s []string) bool {
	for _, v := range s {
		if v == "-" {
			return false
		}
	}
	return true
}

func (pg *backupPage) resetPage(c pageCommon) {
	*c.page = PageWallet
	pg.active = infoView
	pg.seedPhrase = []string{}
	pg.selectedSeeds = make([]string, 33)
	for _, cb := range pg.checkBoxes {
		cb.CheckBox.Value = false
	}
	for i := range pg.suggestions {
		pg.suggestions[i].selected = -1
	}
	for i := range pg.selectedSeeds {
		pg.selectedSeeds[i] = "-"
	}
	pg.updateViewTexts()
}

func (pg *backupPage) handle(c pageCommon) {
	if pg.backButton.Button.Clicked() {
		pg.resetPage(c)
	}

	if pg.action.Button.Clicked() && pg.verifyCheckBoxes() {
		switch pg.active {
		case infoView:
			s := pg.wal.GetWalletSeedPhrase(pg.info.Wallets[*c.selectedWallet].ID)
			pg.seedPhrase = strings.Split(s, " ")
			pg.populateSuggestionSeeds()
			pg.active++
		case verifyView:
			if !checkSlice(pg.selectedSeeds) {
				return
			}
			errMessage := "Failed to verify. Please go through every word and try again."
			s := strings.Join(pg.selectedSeeds, " ")
			if !dcrlibwallet.VerifySeed(s) {
				pg.error = errMessage
				pg.clearError()
				return
			}

			err := pg.wal.VerifyWalletSeedPhrase(pg.info.Wallets[*c.selectedWallet].ID, s)
			if err != nil {
				pg.error = errMessage
				pg.clearError()
				return
			}
			pg.info.Wallets[*c.selectedWallet].Seed = ""
			pg.active++
		case successView:
			pg.resetPage(c)
		default:
			pg.active++
		}
		pg.updateViewTexts()
	}

	for i := range pg.suggestions {
		suggestion := pg.suggestions[i]
		for s := range suggestion.buttons {
			if suggestion.buttons[s].Button.Clicked() {
				pg.suggestions[i].selected = s
				pg.selectedSeeds[i] = suggestion.buttons[s].Text
			}
		}
	}
}
