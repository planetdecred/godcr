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
	seed   string
	rows   []saveSeedRow

	backButton   decredmaterial.IconButton
	actionButton decredmaterial.Button
	container    *layout.List
	seedList     *layout.List
}

func NewSaveSeedPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SaveSeedPage {
	pg := &SaveSeedPage{
		Load:   l,
		wallet: wallet,

		actionButton: l.Theme.Button(new(widget.Clickable), "I have written down all 33 words"),
		container:    &layout.List{Axis: layout.Vertical},
		seedList:     &layout.List{Axis: layout.Vertical},
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *SaveSeedPage) ID() string {
	return SaveSeedPageID
}

func (pg *SaveSeedPage) OnResume() {

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
			}()

			return false
		}).
		NegativeButton("Cancel", func() {}).Show()

}

func (pg *SaveSeedPage) Handle() {
	for pg.actionButton.Clicked() {
		pg.ChangeFragment(NewVerifySeedPage(pg.Load, pg.wallet, pg.seed))
	}
}

func (pg *SaveSeedPage) OnClose() {}

// - Layout

func (pg *SaveSeedPage) Layout(gtx C) D {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Write down seed phrase",
		SubTitle:   "Step 1/2",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			//TODO
			pg.PopToFragment("Wallets")
		},
		Body: func(gtx C) D {

			wdg := []layout.Widget{
				func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, "Write down all 33 words in the correct order.")
					label.Color = pg.Theme.Color.Gray3
					return label.Layout(gtx)
				},
				func(gtx C) D {
					label := pg.Theme.Label(values.TextSize14, "Your 33-word seed phrase")
					label.Color = pg.Theme.Color.Gray3

					return decredmaterial.LinearLayout{
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.WrapContent,
						Orientation: layout.Vertical,
						Background:  pg.Theme.Color.Surface,
						Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
						Margin:      layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding96}, // bottom margin accounts for action button's height + 16dp
						Padding:     layout.Inset{Top: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding8, Left: values.MarginPadding16},
					}.Layout(gtx,
						layout.Rigid(label.Layout),
						layout.Rigid(func(gtx C) D {
							return pg.seedList.Layout(gtx, len(pg.rows), func(gtx C, index int) D {
								return pg.seedRow(gtx, pg.rows[index])
							})
						}),
					)
				},
			}

			return pg.container.Layout(gtx, len(wdg), func(gtx C, index int) D {
				return wdg[index](gtx)
			})
		},
	}

	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return components.UniformPadding(gtx, sp.Layout)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Orientation: layout.Vertical,
				Direction:   layout.S,
				Alignment:   layout.Baseline,
				Background:  sp.Theme.Color.Surface,
				Padding:     layout.UniformInset(values.MarginPadding16),
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize14, "You will be asked to enter the seed phrase on the next screen.")
					label.Color = pg.Theme.Color.Gray3
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					pg.actionButton.Background = pg.Theme.Color.Primary
					pg.actionButton.Color = pg.Theme.Color.InvText

					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, pg.actionButton.Layout)
				}))
		}),
	)
}

func (pg *SaveSeedPage) seedRow(gtx C, row saveSeedRow) D {
	itemWidth := gtx.Constraints.Max.X / 3
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
			indexLabel.Color = theme.Color.Gray3
			indexLabel.Font.Weight = text.Medium
			return decredmaterial.LinearLayout{
				Width:     gtx.Px(values.MarginPadding30),
				Height:    gtx.Px(values.MarginPadding22),
				Direction: layout.Center,
				Margin:    layout.Inset{Right: values.MarginPadding8},
				Border:    decredmaterial.Border{Radius: decredmaterial.Radius(9), Color: theme.Color.Gray3, Width: values.MarginPadding1},
			}.Layout2(gtx, indexLabel.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			seedWord := theme.Label(values.TextSize16, word)
			seedWord.Color = theme.Color.Gray3
			seedWord.Font.Weight = text.Medium
			return seedWord.Layout(gtx)
		}),
	)
}
