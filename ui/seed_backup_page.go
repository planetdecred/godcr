package ui

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"time"

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
	buttonGroup struct {
		skin   decredmaterial.Button
		button *widget.Button
	}

	seedGroup struct {
		selected int
		buttons  []buttonGroup
	}

	viewText struct {
		title       string
		action      string
		steps       string
		instruction string
	}
)

type backupPage struct {
	gtx   *layout.Context
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
	checkBoxes     []decredmaterial.CheckBox
	checkIcon      *decredmaterial.Icon

	backButtonWidget *widget.Button
	actionWidget     *widget.Button
	checkBoxWidgets  []*widget.CheckBox

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
		gtx:   c.gtx,
		theme: c.theme,
		wal:   c.wallet,
		info:  c.info,

		action:         c.theme.Button("View seed phrase"),
		backButton:     c.theme.PlainIconButton(c.icons.navigationArrowBack),
		title:          c.theme.H5("Keep in mind"),
		steps:          c.theme.Body1("Step 1/2"),
		instruction:    c.theme.H6("Write down all 33 words in the correct order"),
		successMessage: c.theme.H4("Your seed phrase backup is verified"),
		successInfo:    c.theme.Body2("Be sure to store your seed phrase backup in a secure location."),
		checkIcon:      c.icons.actionCheckCircle,

		backButtonWidget: new(widget.Button),
		actionWidget:     new(widget.Button),

		active:        infoView,
		selectedSeeds: make([]string, 0, 33),
	}

	b.checkIcon.Color = c.theme.Color.Success
	b.steps.Color = c.theme.Color.Hint
	b.successMessage.Alignment = text.Middle
	b.successInfo.Alignment = text.Middle
	b.successInfo.Color = b.theme.Color.Hint

	b.backButton.Color = c.theme.Color.Hint
	b.backButton.Size = unit.Dp(32)

	b.action.Background = c.theme.Color.Hint

	b.checkBoxes = []decredmaterial.CheckBox{
		c.theme.CheckBox("The 33-word seed phrase is EXTREMELY IMPORTANT."),
		c.theme.CheckBox("Seed phrase is the only way to restore your wallet."),
		c.theme.CheckBox("It is recommended to store your seed phrase in a physical format (e.g. write down on a paper)."),
		c.theme.CheckBox("It is highly discouraged to store your seed phrase in any digital format (e.g. screenshot)."),
		c.theme.CheckBox("Anyone with your seed phrase can steal your funds. DO NOT show it to anyone."),
	}

	b.instruction.Alignment = text.Middle
	b.allSuggestions = dcrlibwallet.PGPWordList()

	for _, cb := range b.checkBoxes {
		cb.IconColor = c.theme.Color.Success
		cb.Color = c.theme.Color.Success
	}

	for i := 0; i < len(b.checkBoxes); i++ {
		b.checkBoxWidgets = append(b.checkBoxWidgets, new(widget.CheckBox))
	}

	for i := 0; i < 33; i++ {
		var bg []buttonGroup

		for j := 0; j < 3; j++ {
			bg = append(bg, buttonGroup{
				skin:   c.theme.Button(""),
				button: new(widget.Button),
			})
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

	return func() {
		b.layout()
		b.handle(c)
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

func (pg *backupPage) layout() {
	pg.theme.Surface(pg.gtx, func() {
		toMax(pg.gtx)
		layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(pg.gtx,
			layout.Rigid(func() {
				pg.action.Background = pg.theme.Color.Hint
				pg.action.Color = pg.theme.Color.InvText
				switch pg.active {
				case infoView:
					if pg.verifyCheckBoxes() {
						pg.activeButton()
					}
					pg.infoView()()
				case seedView:
					pg.activeButton()
					pg.seedView()()
				case verifyView:
					if checkSlice(pg.selectedSeeds) {
						pg.activeButton()
					}
					pg.verifyView()()
				case successView:
					pg.activeButton()
					pg.successView()()
				}
			}),
		)
	})
}

func (pg *backupPage) pageTitle() layout.Widget {
	gtx := pg.gtx
	return func() {
		layout.Inset{Bottom: unit.Dp(5), Top: unit.Dp(20)}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
				layout.Rigid(func() {
					pg.backButton.Layout(gtx, pg.backButtonWidget)
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func() {
								pg.title.Layout(gtx)
							})
						}),
						layout.Rigid(func() {
							if pg.active != infoView {
								layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func() {
									pg.steps.Layout(gtx)
								})
							}
						}),
						layout.Rigid(func() {
							pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
							if pg.active != infoView {
								layout.Inset{Right: unit.Dp(30), Top: unit.Dp(20)}.Layout(gtx, func() {
									pg.instruction.Layout(gtx)
								})
							}
						}),
					)
				}),
			)
		})
	}
}

func (pg *backupPage) viewTemplate(content layout.Widget) layout.Widget {
	return func() {
		layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(pg.gtx, func() {
			layout.Stack{}.Layout(pg.gtx,
				layout.Stacked(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
						layout.Rigid(func() {
							if pg.active != successView {
								pg.pageTitle()()
							}
						}),
						layout.Rigid(func() {
							layout.Inset{Bottom: unit.Dp(50)}.Layout(pg.gtx, func() {
								content()
							})
						}),
					)
				}),
				layout.Stacked(func() {
					pg.gtx.Constraints.Height.Min = pg.gtx.Constraints.Height.Max
					layout.S.Layout(pg.gtx, func() {
						pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
						layout.Inset{Bottom: unit.Dp(10)}.Layout(pg.gtx, func() {
							pg.action.Layout(pg.gtx, pg.actionWidget)
						})
					})
				}),
				layout.Stacked(func() {
					if len(pg.error) > 0 {
						layout.Inset{Top: unit.Dp(75)}.Layout(pg.gtx, func() {
							pg.theme.ErrorAlert(pg.gtx, pg.error)
						})
					}
				}),
			)
		})
	}
}

func (pg *backupPage) infoView() layout.Widget {
	return func() {
		pg.viewTemplate(func() {
			pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
			pg.gtx.Constraints.Height.Min = pg.gtx.Constraints.Height.Max
			layout.Center.Layout(pg.gtx, func() {
				layout.Inset{Bottom: unit.Dp(60)}.Layout(pg.gtx, func() {
					pg.infoList.Layout(pg.gtx, len(pg.checkBoxWidgets), func(i int) {
						layout.Inset{Bottom: unit.Dp(20)}.Layout(pg.gtx, func() {
							pg.checkBoxes[i].Layout(pg.gtx, pg.checkBoxWidgets[i])
						})
					})
				})
			})
		})()
	}
}

func (pg *backupPage) seedView() layout.Widget {
	return func() {
		pg.viewTemplate(
			func() {
				pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
				layout.Center.Layout(pg.gtx, func() {
					pg.viewList.Layout(pg.gtx, 1, func(i int) {
						layout.Flex{Axis: layout.Horizontal}.Layout(pg.gtx,
							layout.Rigid(func() {
								pg.gtx.Constraints.Width.Max = pg.gtx.Constraints.Width.Max / 2
								pg.seedPhraseListLeft.Layout(pg.gtx, len(pg.seedPhrase), func(i int) {
									if i < 17 {
										pg.seedText(i)
									}
								})
							}),
							layout.Rigid(func() {
								pg.seedPhraseListRight.Layout(pg.gtx, len(pg.seedPhrase), func(i int) {
									if i > 16 {
										pg.seedText(i)
									}
								})
							}),
						)
					})
				})
			},
		)()
	}
}

func (pg *backupPage) verifyView() layout.Widget {
	return func() {
		pg.viewTemplate(func() {
			toMax(pg.gtx)
			pg.verifyList.Layout(pg.gtx, len(pg.suggestions), func(i int) {
				s := pg.suggestions[i]
				layout.Center.Layout(pg.gtx, func() {
					layout.Inset{Bottom: unit.Dp(30)}.Layout(pg.gtx, func() {
						layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(15), Bottom: unit.Dp(15)}.Layout(pg.gtx, func() {
									pg.theme.H6(fmt.Sprintf("%d. %s", i+1, pg.selectedSeeds[i])).Layout(pg.gtx)
								})
							}),
							layout.Rigid(func() {
								layout.Flex{Axis: layout.Horizontal}.Layout(pg.gtx,
									layout.Flexed(0.3, func() {
										pg.suggestionButtonGroup(s, 0)
									}),
									layout.Flexed(0.3, func() {
										pg.suggestionButtonGroup(s, 1)
									}),
									layout.Flexed(0.3, func() {
										pg.suggestionButtonGroup(s, 2)
									}),
								)
							}),
						)
					})
				})
			})
		})()
	}
}

func (pg *backupPage) successView() layout.Widget {
	return func() {
		pg.viewTemplate(func() {
			layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
				layout.Rigid(func() {
					pg.gtx.Constraints.Height.Min = pg.gtx.Constraints.Height.Max
					layout.Center.Layout(pg.gtx, func() {
						layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
							layout.Rigid(func() {
								pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
								layout.Inset{Bottom: unit.Dp(50), Right: unit.Dp(50)}.Layout(pg.gtx, func() {
									layout.Center.Layout(pg.gtx, func() {
										layout.UniformInset(unit.Dp(20)).Layout(pg.gtx, func() {
											pg.checkIcon.Color = pg.theme.Color.Success
											pg.checkIcon.Layout(pg.gtx, unit.Px(float32(150)))
										})
									})
								})
							}),
							layout.Rigid(func() {
								pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
								pg.successMessage.Layout(pg.gtx)
							}),
							layout.Rigid(func() {
								pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max
								pg.successInfo.Layout(pg.gtx)
							}),
						)
					})
				}),
			)
		})()
	}
}

func (pg *backupPage) seedText(index int) {
	layout.Inset{Bottom: unit.Dp(10), Left: unit.Dp(20)}.Layout(pg.gtx,
		func() {
			seedLabel := pg.theme.H6(fmt.Sprintf("%d.  %s", index+1, pg.seedPhrase[index]))
			seedLabel.Alignment = text.Middle
			seedLabel.Layout(pg.gtx)
		},
	)
}

func (pg *backupPage) suggestionButtonGroup(sg seedGroup, buttonIndex int) {
	button := sg.buttons[buttonIndex]
	button.skin.Background = pg.theme.Color.Hint
	button.skin.TextSize = unit.Dp(18)
	if sg.selected == buttonIndex {
		button.skin.Background = pg.theme.Color.Primary
	}
	layout.Inset{Right: unit.Dp(15), Left: unit.Dp(15)}.Layout(pg.gtx, func() {
		button.skin.Layout(pg.gtx, sg.buttons[buttonIndex].button)
	})
}

func (pg *backupPage) verifyCheckBoxes() bool {
	for _, cb := range pg.checkBoxWidgets {
		if !cb.Checked(pg.gtx) {
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
			s.buttons[i].skin.Text = seeds[i]
		}
		s.buttons[rand.Intn(len(seeds))].skin.Text = pg.seedPhrase[k]
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
	for _, cb := range pg.checkBoxWidgets {
		cb.SetChecked(false)
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
	if pg.backButtonWidget.Clicked(pg.gtx) {
		pg.resetPage(c)
	}

	if pg.actionWidget.Clicked(pg.gtx) && pg.verifyCheckBoxes() {
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
			if suggestion.buttons[s].button.Clicked(pg.gtx) {
				pg.suggestions[i].selected = s
				pg.selectedSeeds[i] = suggestion.buttons[s].skin.Text
			}
		}
	}
}
