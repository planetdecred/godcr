package seedbackup

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
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
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	wallet *dcrlibwallet.Wallet

	backButton   decredmaterial.IconButton
	actionButton decredmaterial.Button
	seedList     *widget.List
	hexLabel     decredmaterial.Label
	copy         decredmaterial.Button

	infoText   string
	seed       string
	rows       []saveSeedRow
	mobileRows []saveSeedRow
}

func NewSaveSeedPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SaveSeedPage {
	pg := &SaveSeedPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(SaveSeedPageID),
		wallet:           wallet,
		hexLabel:         l.Theme.Label(values.TextSize12, ""),
		copy:             l.Theme.Button("Copy"),
		infoText:         "You will be asked to enter the seed word on the next screen.",
		actionButton:     l.Theme.Button("I have written down all 33 words"),
		seedList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.copy.TextSize = values.TextSize12
	pg.hexLabel.MaxLines = 1
	pg.copy.Background = color.NRGBA{}
	pg.copy.HighlightColor = pg.Theme.Color.SurfaceHighlight
	pg.copy.Color = pg.Theme.Color.Primary
	pg.copy.Inset = layout.UniformInset(values.MarginPadding16)

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.backButton.Icon = l.Theme.Icons.ContentClear

	pg.actionButton.Font.Weight = text.Medium

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SaveSeedPage) OnNavigatedTo() {

	passwordModal := modal.NewPasswordModal(pg.Load).
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
			pg.ParentNavigator().ClosePagesAfter(components.WalletsPageID)
		})
	pg.ParentWindow().ShowModal(passwordModal)

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SaveSeedPage) HandleUserInteractions() {
	for pg.actionButton.Clicked() {
		pg.ParentNavigator().Display(NewVerifySeedPage(pg.Load, pg.wallet, pg.seed))
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
			promptToExit(pg.Load, pg.ParentNavigator(), pg.ParentWindow())
		},
		Body: func(gtx C) D {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, "Write down all 33 words in the correct order.")
					label.Color = pg.Theme.Color.GrayText1
					return label.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					label := pg.Theme.Label(values.TextSize14, "Your 33-word seed word")
					label.Color = pg.Theme.Color.GrayText1
					return decredmaterial.LinearLayout{
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.WrapContent,
						Orientation: layout.Vertical,
						Background:  pg.Theme.Color.Surface,
						Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
						Margin:      layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding2},
						Padding:     layout.Inset{Top: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding8, Left: values.MarginPadding16},
					}.Layout(gtx,
						layout.Rigid(label.Layout),
						layout.Rigid(func(gtx C) D {
							return pg.Theme.List(pg.seedList).Layout(gtx, len(pg.rows), func(gtx C, index int) D {
								return pg.desktopSeedRow(gtx, pg.rows[index])
							})
						}),
					)
				}),
				layout.Flexed(1, pg.hexLayout),
			)
		},
	}

	layout := func(gtx C) D {
		return sp.Layout(pg.ParentWindow(), gtx)
	}

	return container(gtx, false, *pg.Theme, layout, pg.infoText, pg.actionButton)
}

func (pg *SaveSeedPage) layoutMobile(gtx C) D {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Write down seed word",
		SubTitle:   "Step 1/2",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load, pg.ParentNavigator(), pg.ParentWindow())
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
	layout := func(gtx C) D {
		return sp.Layout(pg.ParentWindow(), gtx)
	}

	return container(gtx, true, *pg.Theme, layout, pg.infoText, pg.actionButton)
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

func (pg *SaveSeedPage) hexLayout(gtx layout.Context) layout.Dimensions {
	pg.handleCopyEvent(gtx)
	card := decredmaterial.Card{
		Color: pg.Theme.Color.Gray4,
	}

	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Vertical,
		Background:  pg.Theme.Color.Surface,
		Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
		Margin:      layout.Inset{Top: values.MarginPadding0, Bottom: values.MarginPadding16},
		Padding:     layout.Inset{Top: values.MarginPadding5, Right: values.MarginPadding16, Bottom: values.MarginPadding16, Left: values.MarginPadding16},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := pg.Theme.Label(values.TextSize14, "Seed hex")
			label.Color = pg.Theme.Color.GrayText1

			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					card.Radius = decredmaterial.CornerRadius{TopRight: 0, TopLeft: 8, BottomRight: 0, BottomLeft: 8}

					return card.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							seedString := pg.seed
							if seedString != "" {
								hexString, _ := components.SeedWordsToHex(pg.seed)
								pg.hexLabel.Text = hexString
							}
							return pg.hexLabel.Layout(gtx)

						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					card.Radius = decredmaterial.CornerRadius{TopRight: 8, TopLeft: 0, BottomRight: 8, BottomLeft: 0}
					return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D {
						return card.Layout(gtx, pg.copy.Layout)
					})
				}),
			)
		}),
	)
}

func (pg *SaveSeedPage) handleCopyEvent(gtx layout.Context) {
	if pg.copy.Clicked() {
		clipboard.WriteOp{Text: pg.hexLabel.Text}.Add(gtx.Ops)

		pg.copy.Text = "Copied!"
		pg.copy.Color = pg.Theme.Color.Success
		time.AfterFunc(time.Second*3, func() {
			pg.copy.Text = "Copy"
			pg.copy.Color = pg.Theme.Color.Primary
			pg.ParentWindow().Reload()
		})
	}
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
