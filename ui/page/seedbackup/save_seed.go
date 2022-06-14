package seedbackup

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const SaveSeedPageID = "save_seed"

type saveSeedRow struct {
	rowIndex int
	word1    string
	word2    string
	word3    string
}

type SaveSeedPage struct {
	*load.Load
	wallet *dcrlibwallet.Wallet

	backButton   decredmaterial.IconButton
	actionButton decredmaterial.Button
	seedList     *widget.List

	infoText   string
	seed       string
	rows       []saveSeedRow
	mobileRows []saveSeedRow
}

func NewSaveSeedPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SaveSeedPage {
	pg := &SaveSeedPage{
		Load:         l,
		wallet:       wallet,
		infoText:     "You will be asked to enter the seed word on the next screen.",
		actionButton: l.Theme.Button("I have written down all 33 words"),
		seedList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.backButton.Icon = l.Theme.Icons.ContentClear

	pg.actionButton.Font.Weight = text.Medium

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *SaveSeedPage) ID() string {
	return SaveSeedPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SaveSeedPage) OnNavigatedTo() {

	modal.NewPasswordModal(pg.Load).
		Title("Confirm to show seed").
		PositiveButton("Confirm", func(password string, m *modal.PasswordModal) bool {
			go func() {
				seed, err := pg.wallet.DecryptSeed([]byte(password))
				if err != nil {
					m.SetLoading(false)
					m.SetError(err.Error())
					return
				}

				m.Dismiss()

				pg.seed = seed

				wordList := strings.Split(seed, " ")
				row1 := wordList[:11]
				row2 := wordList[11:22]
				row3 := wordList[22:]

				//for mobile
				rowMobile1 := wordList[:17]
				rowMobile2 := wordList[17:]
				mobileRows := make([]saveSeedRow, 0)
				for i := range rowMobile1 {
					r2 := ""
					if i < len(rowMobile2) {
						r2 = rowMobile2[i]
					}
					mobileRows = append(mobileRows, saveSeedRow{
						rowIndex: i + 1,
						word1:    rowMobile1[i],
						word2:    r2,
					})
				}

				rows := make([]saveSeedRow, 0)
				for i := range row1 {
					rows = append(rows, saveSeedRow{
						rowIndex: i + 1,
						word1:    row1[i],
						word2:    row2[i],
						word3:    row3[i],
					})
				}
				pg.rows = rows
				pg.mobileRows = mobileRows
			}()

			return false
		}).
		NegativeButton("Cancel", func() {
			pg.PopToFragment(components.WalletsPageID)
		}).Show()

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SaveSeedPage) HandleUserInteractions() {
	for pg.actionButton.Clicked() {
		pg.ChangeFragment(NewVerifySeedPage(pg.Load, pg.wallet, pg.seed))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SaveSeedPage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SaveSeedPage) Layout(gtx C) D {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *SaveSeedPage) layoutDesktop(gtx C) D {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Write down seed word",
		SubTitle:   "Step 1/2",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load)
		},
		Body: func(gtx C) D {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, "Write down all 33 words in the correct order.")
					label.Color = pg.Theme.Color.GrayText1
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize14, "Your 33-word seed word")
					label.Color = pg.Theme.Color.GrayText1
					return decredmaterial.LinearLayout{
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.WrapContent,
						Orientation: layout.Vertical,
						Background:  pg.Theme.Color.Surface,
						Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
						// bottom margin accounts for action button's height + components.UniformPadding bottom margin 24dp + 16dp
						Margin:  layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding120},
						Padding: layout.Inset{Top: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding8, Left: values.MarginPadding16},
					}.Layout(gtx,
						layout.Rigid(label.Layout),
						layout.Rigid(func(gtx C) D {
							return pg.Theme.List(pg.seedList).Layout(gtx, len(pg.rows), func(gtx C, index int) D {
								return pg.desktopSeedRow(gtx, pg.rows[index])
							})
						}),
					)
				}),
			)
		},
	}

	return container(gtx, false, *pg.Theme, sp.Layout, pg.infoText, pg.actionButton)
}

func (pg *SaveSeedPage) layoutMobile(gtx C) D {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Write down seed word",
		SubTitle:   "Step 1/2",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load)
		},
		Body: func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, "Write down all 33 words in the correct order.")
					label.Color = pg.Theme.Color.GrayText1
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize14, "Your 33-word seed word")
					label.Color = pg.Theme.Color.GrayText1

					return decredmaterial.LinearLayout{
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.WrapContent,
						Orientation: layout.Vertical,
						Background:  pg.Theme.Color.Surface,
						Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
						// bottom margin accounts for action button's height + components.UniformPadding bottom margin 24dp + 16dp
						Margin:  layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding120},
						Padding: layout.Inset{Top: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding8, Left: values.MarginPadding16},
					}.Layout(gtx,
						layout.Rigid(label.Layout),
						layout.Rigid(func(gtx C) D {
							return pg.Theme.List(pg.seedList).Layout(gtx, len(pg.mobileRows), func(gtx C, index int) D {
								return pg.mobileSeedRow(gtx, pg.mobileRows[index])
							})
						}),
					)
				}),
			)
		},
	}

	return container(gtx, true, *pg.Theme, sp.Layout, pg.infoText, pg.actionButton)
}

func (pg *SaveSeedPage) mobileSeedRow(gtx C, row saveSeedRow) D {
	itemWidth := gtx.Constraints.Max.X / 2 // Divide total width into 2 rows for mobile
	topMargin := values.MarginPadding8
	if row.rowIndex == 1 {
		topMargin = values.MarginPadding16
	}
	return decredmaterial.LinearLayout{
		Width:  decredmaterial.MatchParent,
		Height: decredmaterial.WrapContent,
		Margin: layout.Inset{Top: topMargin},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return seedItem(pg.Theme, gtx, itemWidth, row.rowIndex, row.word1)
		}),
		layout.Rigid(func(gtx C) D {
			if row.word2 == "" {
				return layout.Dimensions{}
			}
			return seedItem(pg.Theme, gtx, itemWidth, row.rowIndex+17, row.word2)
		}),
	)
}

func (pg *SaveSeedPage) desktopSeedRow(gtx C, row saveSeedRow) D {
	itemWidth := gtx.Constraints.Max.X / 3 // Divide total width into 3 rows for deskop
	topMargin := values.MarginPadding8
	if row.rowIndex == 1 {
		topMargin = values.MarginPadding16
	}
	return decredmaterial.LinearLayout{
		Width:  decredmaterial.MatchParent,
		Height: decredmaterial.WrapContent,
		Margin: layout.Inset{Top: topMargin},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return seedItem(pg.Theme, gtx, itemWidth, row.rowIndex, row.word1)
		}),
		layout.Rigid(func(gtx C) D {
			return seedItem(pg.Theme, gtx, itemWidth, row.rowIndex+11, row.word2)
		}),
		layout.Rigid(func(gtx C) D {
			return seedItem(pg.Theme, gtx, itemWidth, row.rowIndex+22, row.word3)
		}),
	)
}

func seedItem(theme *decredmaterial.Theme, gtx C, width, index int, word string) D {
	return decredmaterial.LinearLayout{
		Width:  width,
		Height: decredmaterial.WrapContent,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			indexLabel := theme.Label(values.TextSize16, fmt.Sprint(index))
			indexLabel.Color = theme.Color.GrayText1
			indexLabel.Font.Weight = text.Medium
			return decredmaterial.LinearLayout{
				Width:     gtx.Dp(values.MarginPadding30),
				Height:    gtx.Dp(values.MarginPadding22),
				Direction: layout.Center,
				Margin:    layout.Inset{Right: values.MarginPadding8},
				Border:    decredmaterial.Border{Radius: decredmaterial.Radius(9), Color: theme.Color.Gray3, Width: values.MarginPadding1},
			}.Layout2(gtx, indexLabel.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			seedWord := theme.Label(values.TextSize16, word)
			seedWord.Color = theme.Color.GrayText1
			seedWord.Font.Weight = text.Medium
			return seedWord.Layout(gtx)
		}),
	)
}
