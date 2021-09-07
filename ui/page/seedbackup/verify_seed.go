package seedbackup

import (
	"math/rand"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const VerifySeedPageID = "verify_seed"

type shuffledSeedWords struct {
	selectedIndex int
	words         []string
	clickables    []*widget.Clickable
}

type VerifySeedPage struct {
	*load.Load
	wallet        *dcrlibwallet.Wallet
	seed          string
	multiSeedList []shuffledSeedWords

	backButton   decredmaterial.IconButton
	actionButton decredmaterial.Button
	container    *layout.List
	seedList     *layout.List
}

func NewVerifySeedPage(l *load.Load, wallet *dcrlibwallet.Wallet, seed string) *VerifySeedPage {
	pg := &VerifySeedPage{
		Load:   l,
		wallet: wallet,
		seed:   seed,

		actionButton: l.Theme.Button(new(widget.Clickable), "Verify"),
		container:    &layout.List{Axis: layout.Vertical},
		seedList:     &layout.List{Axis: layout.Vertical},
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.backButton.Icon = l.Icons.ContentClear

	return pg
}

func (pg *VerifySeedPage) ID() string {
	return SaveSeedPageID
}

func (pg *VerifySeedPage) OnResume() {
	allSeeds := dcrlibwallet.PGPWordList()

	multiSeedList := make([]shuffledSeedWords, 0)
	seedWords := strings.Split(pg.seed, " ")
	rand.Seed(time.Now().UnixNano())
	for _, word := range seedWords {
		index := seedPosition(word, allSeeds)
		shuffledSeed := pg.getMultiSeed(index, dcrlibwallet.PGPWordList()) // using allSeeds here modifies the slice
		multiSeedList = append(multiSeedList, shuffledSeed)
	}

	pg.multiSeedList = multiSeedList
}

func (pg *VerifySeedPage) getMultiSeed(realSeedIndex int, allSeeds []string) shuffledSeedWords {
	shuffledSeed := shuffledSeedWords{
		selectedIndex: -1,
		words:         make([]string, 0),
		clickables:    make([]*widget.Clickable, 0),
	}

	shuffledSeed.words = append(shuffledSeed.words, allSeeds[realSeedIndex])
	shuffledSeed.clickables = append(shuffledSeed.clickables, &widget.Clickable{})
	allSeeds = removeSeed(allSeeds, realSeedIndex)

	for i := 0; i < 3; i++ {
		randomSeed := rand.Intn(len(allSeeds))

		shuffledSeed.words = append(shuffledSeed.words, allSeeds[randomSeed])
		shuffledSeed.clickables = append(shuffledSeed.clickables, &widget.Clickable{})
		allSeeds = removeSeed(allSeeds, randomSeed)
	}

	rand.Shuffle(len(shuffledSeed.words), func(i, j int) {
		shuffledSeed.words[i], shuffledSeed.words[j] = shuffledSeed.words[j], shuffledSeed.words[i]
	})

	return shuffledSeed
}

func seedPosition(seed string, allSeeds []string) int {
	for i := range allSeeds {
		if allSeeds[i] == seed {
			return i
		}
	}
	return -1
}

func removeSeed(allSeeds []string, index int) []string {
	return append(allSeeds[:index], allSeeds[index+1:]...)
}

func (pg *VerifySeedPage) allSeedsSelected() bool {
	for _, multiSeed := range pg.multiSeedList {
		if multiSeed.selectedIndex == -1 {
			return false
		}
	}

	return true
}

func (pg *VerifySeedPage) selectedSeedPhrase() string {
	var wordList []string
	for _, multiSeed := range pg.multiSeedList {
		if multiSeed.selectedIndex != -1 {
			wordList = append(wordList, multiSeed.words[multiSeed.selectedIndex])
		}
	}

	return strings.Join(wordList, " ")
}

func (pg *VerifySeedPage) verifySeed() {
	modal.NewPasswordModal(pg.Load).
		Title("Confirm to verify seed").
		PositiveButton("Confirm", func(password string, m *modal.PasswordModal) bool {
			go func() {
				seed := pg.selectedSeedPhrase()
				_, err := pg.WL.MultiWallet.VerifySeedForWallet(pg.wallet.ID, seed, []byte(password))
				if err != nil {
					if err.Error() == dcrlibwallet.ErrInvalid {
						pg.Toast.NotifyError("Failed to verify. Please go through every word and try again.")
						m.Dismiss()
						return
					}

					m.SetLoading(false)
					m.SetError(err.Error())
					return
				}
				m.Dismiss()

				pg.ChangeFragment(NewBackupSuccessPage(pg.Load))
			}()

			return false
		}).
		NegativeButton("Cancel", func() {}).Show()
}

func (pg *VerifySeedPage) Handle() {
	for i, multiSeed := range pg.multiSeedList {
		for j, clickable := range multiSeed.clickables {
			for clickable.Clicked() {
				pg.multiSeedList[i].selectedIndex = j
			}
		}
	}

	for pg.actionButton.Clicked() {
		if pg.allSeedsSelected() {
			pg.verifySeed()
		}
	}
}

func (pg *VerifySeedPage) OnClose() {}

// - Layout

func (pg *VerifySeedPage) Layout(gtx C) D {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Verify seed phrase",
		SubTitle:   "Step 2/2",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load)
		},
		Body: func(gtx C) D {
			wdg := []layout.Widget{
				func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, "Select the correct words to verify.")
					label.Color = pg.Theme.Color.Gray3
					return label.Layout(gtx)
				},
				func(gtx C) D {
					return layout.Inset{
						Bottom: values.MarginPadding96,
					}.Layout(gtx, func(gtx C) D {
						return pg.seedList.Layout(gtx, len(pg.multiSeedList), func(gtx C, index int) D {
							return pg.seedListRow(gtx, index, pg.multiSeedList[index])
						})
					})
				},
			}

			return pg.container.Layout(gtx, len(wdg), func(gtx C, index int) D {
				return wdg[index](gtx)
			})
		},
	}

	if pg.allSeedsSelected() {
		pg.actionButton.Background = pg.Theme.Color.Primary
		pg.actionButton.Color = pg.Theme.Color.InvText
	} else {
		pg.actionButton.Background = pg.Theme.Color.InactiveGray
		pg.actionButton.Color = pg.Theme.Color.Text
	}

	return container(gtx, *pg.Theme, sp.Layout, "", pg.actionButton)
}

func (pg *VerifySeedPage) seedListRow(gtx C, index int, multiSeed shuffledSeedWords) D {
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Vertical,
		Background:  pg.Theme.Color.Surface,
		Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
		Margin:      layout.Inset{Top: values.MarginPadding4, Bottom: values.MarginPadding4},
		Padding:     layout.Inset{Top: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding8, Left: values.MarginPadding16},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			text := "-"
			if multiSeed.selectedIndex != -1 {
				text = multiSeed.words[multiSeed.selectedIndex]
			}
			return seedItem(pg.Theme, gtx, gtx.Constraints.Max.X, index+1, text)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D { return pg.seedButton(gtx, 0, multiSeed) }),
					layout.Rigid(func(gtx C) D { return pg.seedButton(gtx, 1, multiSeed) }),
					layout.Rigid(func(gtx C) D { return pg.seedButton(gtx, 2, multiSeed) }),
					layout.Rigid(func(gtx C) D { return pg.seedButton(gtx, 3, multiSeed) }),
				)
			})
		}),
	)
}

func (pg *VerifySeedPage) seedButton(gtx C, index int, multiSeed shuffledSeedWords) D {
	borderColor := pg.Theme.Color.Gray1
	textColor := pg.Theme.Color.Gray2
	if index == multiSeed.selectedIndex {
		borderColor = pg.Theme.Color.Primary
		textColor = pg.Theme.Color.Primary
	}

	return decredmaterial.Clickable(gtx, multiSeed.clickables[index], func(gtx C) D {

		return decredmaterial.LinearLayout{
			Width:      gtx.Px(values.MarginPadding100),
			Height:     gtx.Px(values.MarginPadding40),
			Background: pg.Theme.Color.Surface,
			Direction:  layout.Center,
			Border:     decredmaterial.Border{Radius: decredmaterial.Radius(8), Color: borderColor, Width: values.MarginPadding2},
		}.Layout2(gtx, func(gtx C) D {
			label := pg.Theme.Label(values.TextSize16, multiSeed.words[index])
			label.Color = textColor
			return label.Layout(gtx)
		})
	})
}
