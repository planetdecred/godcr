package ui

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

var testPhrase = "tissue recover scorecard Istanbul solo dinosaur framework forever freedom typewriter spheroid " +
	"Capricorn standard paperweight drainage informant steamship gossamer klaxon conformist quota provincial erase " +
	"paperweight soybean universe blowtorch sandalwood drumbeat dictator unearth bravado lockup"

const (
	PageSeedBackup = "seedbackup"
	infoView       = iota
	seedView
	verifyView
)

type backupPage struct {
	gtx   *layout.Context
	theme *decredmaterial.Theme
	wal   *wallet.Wallet
	info  *wallet.MultiWalletInfo

	backButton     decredmaterial.IconButton
	titleLabel     decredmaterial.Label
	viewSeedPhrase decredmaterial.Button
	checkBoxes     []decredmaterial.CheckBox

	backButtonWidget     *widget.Button
	viewSeedPhraseWidget *widget.Button
	checkBoxWidgets      []*widget.CheckBox

	container      *layout.List
	infoList       *layout.List
	seedPhraseList *layout.List

	seedPhrase     []string
	active         int
	selectedWallet int
}

func (win *Window) BackupPage(c pageCommon) layout.Widget {
	b := &backupPage{
		gtx:   c.gtx,
		theme: c.theme,
		wal:   c.wallet,
		info:  c.info,

		viewSeedPhrase: c.theme.Button("View seed phrase"),
		backButton:     c.theme.PlainIconButton(c.icons.navigationArrowBack),
		titleLabel:     c.theme.H5("Keep in mind"),

		backButtonWidget:     new(widget.Button),
		viewSeedPhraseWidget: new(widget.Button),
		container:            &layout.List{Axis: layout.Vertical},

		active:         infoView,
		selectedWallet: *c.selectedWallet,
	}

	b.backButton.Color = c.theme.Color.Hint
	b.backButton.Size = unit.Dp(32)

	b.viewSeedPhrase.Background = c.theme.Color.Hint

	b.checkBoxes = []decredmaterial.CheckBox{
		c.theme.CheckBox("The 33-word seed phrase is EXTREMELY IMPORTANT."),
		c.theme.CheckBox("Seed phrase is the only way to restore your wallet."),
		c.theme.CheckBox("It is recommended to store your seed phrase in a physical format (e.g. write down on a paper)."),
		c.theme.CheckBox("It is highly discouraged to store your seed phrase in any digital format (e.g. screenshot)."),
		c.theme.CheckBox("Anyone with your seed phrase can steal your funds. DO NOT show it to anyone."),
	}

	for _, cb := range b.checkBoxes {
		cb.IconColor = c.theme.Color.Success
		cb.Color = c.theme.Color.Success
	}

	for i := 0; i < len(b.checkBoxes); i++ {
		b.checkBoxWidgets = append(b.checkBoxWidgets, new(widget.CheckBox))
	}

	b.infoList = &layout.List{Axis: layout.Vertical}
	b.seedPhraseList = &layout.List{Axis: layout.Vertical}

	return func() {
		b.layout()
		b.handle(c)
	}
}

func (pg *backupPage) layout() {
	pg.theme.Surface(pg.gtx, func() {
		toMax(pg.gtx)
		layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(pg.gtx,
			layout.Rigid(func() {
				switch pg.active {
				case infoView:
					pg.infoView()()
				case seedView:
					pg.seedView()()
				case verifyView:
					pg.verifyView()()
				}
			}),
		)
	})
}

func (pg *backupPage) pageTitle() layout.Widget {
	gtx := pg.gtx
	return func() {
		layout.Inset{Bottom: unit.Dp(50), Top: unit.Dp(20)}.Layout(pg.gtx, func() {
			layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
				layout.Rigid(func() {
					pg.backButton.Layout(gtx, pg.backButtonWidget)
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func() {
						pg.titleLabel.Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *backupPage) viewTemplate(content layout.Widget) layout.Widget {
	return func() {
		pg.gtx.Constraints.Height.Min = pg.gtx.Constraints.Height.Max
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
					layout.Rigid(pg.pageTitle()),
					layout.Rigid(func() {
						content()
					}),
				)
			}),
			layout.Rigid(func() {
				if pg.verifyCheckBoxes() {
					pg.viewSeedPhrase.Background = pg.theme.Color.Primary
				}
				pg.viewSeedPhrase.Layout(pg.gtx, pg.viewSeedPhraseWidget)
			}),
		)
	}
}

func (pg *backupPage) infoView() layout.Widget {
	return func() {
		pg.viewTemplate(func() {
			pg.infoList.Layout(pg.gtx, len(pg.checkBoxWidgets), func(i int) {
				layout.Inset{Bottom: unit.Dp(20)}.Layout(pg.gtx, func() {
					pg.checkBoxes[i].Layout(pg.gtx, pg.checkBoxWidgets[i])
				})
			})
		})()
	}
}

func (pg *backupPage) seedView() layout.Widget {
	return func() {
		pg.viewTemplate(
			func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Flex{}.Layout(pg.gtx,
						layout.Rigid(func() {
							pg.gtx.Constraints.Width.Max = pg.gtx.Constraints.Width.Max / 2
							(&layout.List{Axis: layout.Vertical}).Layout(pg.gtx, len(pg.seedPhrase), func(i int) {
								if i < 17 {
									pg.theme.Body2(pg.seedPhrase[i]).Layout(pg.gtx)
								}
							})
						}),
						layout.Rigid(func() {
							(&layout.List{Axis: layout.Vertical}).Layout(pg.gtx, len(pg.seedPhrase), func(i int) {
								if i > 16 {
									pg.theme.Body2(pg.seedPhrase[i]).Layout(pg.gtx)
								}
							})
						}),
					)
				})
			},
		)()
	}
}

func (pg *backupPage) verifyView() layout.Widget {
	return func() {
		pg.viewTemplate(func() {
			pg.theme.Body2("verify view").Layout(pg.gtx)
		})()
	}
}

func (pg *backupPage) verifyCheckBoxes() bool {
	for _, cb := range pg.checkBoxWidgets {
		if !cb.Checked(pg.gtx) {
			return false
		}
	}
	return true
}

func (pg *backupPage) handle(c pageCommon) {
	if pg.backButtonWidget.Clicked(pg.gtx) {
		if pg.active == 1 {
			*c.page = PageWallet
		} else {
			pg.viewSeedPhrase.Text = "View seed phrase"
			pg.active -= 1
		}
	}

	if pg.viewSeedPhraseWidget.Clicked(pg.gtx) {
		if pg.verifyCheckBoxes() && pg.active == 1 {
			// seedPhrase := pg.wal.GetWalletSeedPhrase(pg.info.Wallets[pg.selectedWallet].ID)
			pg.seedPhrase = strings.Split(testPhrase, " ")
			pg.viewSeedPhrase.Text = "I have written down all 33 words"
			pg.active += 1
		} else if pg.active != verifyView {
			pg.active += 1
		}
	}
}
