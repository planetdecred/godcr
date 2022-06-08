package wallets

import (
	"fmt"

	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	CreateRestorePageID = "Restore"
	numberOfSeeds       = 32
)

type seedEditors struct {
	focusIndex int
	editors    []decredmaterial.RestoreEditor
}

type seedItemMenu struct {
	text   string
	button decredmaterial.Button
}

type Restore struct {
	*load.Load
	isRestoring     bool
	restoreComplete func()

	seedList *layout.List

	backButton      decredmaterial.IconButton
	validateSeed    decredmaterial.Button
	resetSeedFields decredmaterial.Button
	optionsMenuCard decredmaterial.Card

	suggestions    []string
	allSuggestions []string
	focused        []int
	seedMenu       []seedItemMenu

	seedPhrase string

	openPopupIndex  int
	selected        int
	suggestionLimit int

	seedClicked  bool
	isLastEditor bool

	seedEditors              seedEditors
	nextcurrentCaretPosition int // caret position when seed editor is switched
	currentCaretPosition     int // current caret position
	selectedSeedEditor       int // stores the current focus index of seed editors
}

func NewRestorePage(l *load.Load, onRestoreComplete func()) *Restore {
	pg := &Restore{
		Load:            l,
		restoreComplete: onRestoreComplete,

		seedList: &layout.List{Axis: layout.Vertical},

		suggestionLimit: 3,
		openPopupIndex:  -1,
	}

	pg.optionsMenuCard = decredmaterial.Card{Color: pg.Theme.Color.Surface}
	pg.optionsMenuCard.Radius = decredmaterial.Radius(8)

	pg.validateSeed = l.Theme.Button("Validate wallet seeds")
	pg.validateSeed.Font.Weight = text.Medium

	pg.resetSeedFields = l.Theme.OutlineButton("Clear all")
	pg.resetSeedFields.Font.Weight = text.Medium

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.backButton.Icon = pg.Theme.Icons.ContentClear

	for i := 0; i <= numberOfSeeds; i++ {
		widgetEditor := new(widget.Editor)
		widgetEditor.SingleLine, widgetEditor.Submit = true, true
		pg.seedEditors.editors = append(pg.seedEditors.editors, l.Theme.RestoreEditor(widgetEditor, "", fmt.Sprintf("%d", i+1)))
	}
	pg.seedEditors.focusIndex = -1
	pg.seedEditors.editors[0].Edit.Editor.Focus()

	// init suggestion buttons
	pg.initSeedMenu()

	// set suggestions
	pg.allSuggestions = dcrlibwallet.PGPWordList()

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *Restore) ID() string {
	return CreateRestorePageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Restore) OnNavigatedTo() {}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *Restore) Layout(gtx layout.Context) layout.Dimensions {
	pg.Load.SetCurrentAppWidth(gtx.Constraints.Max.X)
	if pg.Load.GetCurrentAppWidth() <= gtx.Px(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}

	return pg.layoutDesktop(gtx)
}

func (pg *Restore) layoutDesktop(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Restore wallet",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.restore(gtx)
			},
		}
		return sp.Layout(gtx)
	}

	pg.resetSeedFields.SetEnabled(pg.updateSeedResetBtn())
	seedValid, _ := pg.validateSeeds()
	pg.validateSeed.SetEnabled(seedValid)

	return components.UniformPadding(gtx, body)
}

func (pg *Restore) layoutMobile(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Restore wallet",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.restoreMobile(gtx)
			},
		}
		return sp.Layout(gtx)
	}

	pg.resetSeedFields.SetEnabled(pg.updateSeedResetBtn())
	seedValid, _ := pg.validateSeeds()
	pg.validateSeed.SetEnabled(seedValid)

	return components.UniformPadding(gtx, body)
}

func (pg *Restore) restore(gtx layout.Context) layout.Dimensions {
	dims := layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Orientation: layout.Vertical,
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Background:  pg.Theme.Color.Surface,
				Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
				Padding:     layout.UniformInset(values.MarginPadding15)}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Bottom: values.MarginPadding10,
					}.Layout(gtx, pg.Theme.Body1("Enter your seed phrase").Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.seedEditorViewDesktop(gtx)
				}),
				layout.Rigid(pg.resetSeedFields.Layout),
			)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1, Top: values.MarginPadding20}.Layout(gtx, pg.restoreButtonSection)
			})
		}),
	)
	return dims
}

func (pg *Restore) restoreMobile(gtx layout.Context) layout.Dimensions {
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(0.9, func(gtx C) D {
			return layout.Stack{Alignment: layout.N}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return decredmaterial.LinearLayout{
						Orientation: layout.Vertical,
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.WrapContent,
						Background:  pg.Theme.Color.Surface,
						Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
						Padding:     layout.UniformInset(values.MarginPadding15)}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Bottom: values.MarginPadding10,
							}.Layout(gtx, pg.Theme.Body1("Enter your seed phrase").Layout)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Flexed(0.92, func(gtx C) D {
									return pg.seedEditorViewMobile(gtx)
								}),
								layout.Flexed(0.08, func(gtx C) D {
									return pg.resetSeedFields.Layout(gtx)
								}),
							)

						}),
					)
				}),
			)
		}),
		layout.Flexed(0.1, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1, Top: values.MarginPadding20}.Layout(gtx, pg.restoreButtonSection)
			})
		}),
	)
	return dims
}

func (pg *Restore) restoreButtonSection(gtx layout.Context) layout.Dimensions {
	card := pg.Theme.Card()
	card.Radius = decredmaterial.Radius(0)
	return card.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return pg.validateSeed.Layout(gtx)
	})
}

func (pg *Restore) seedEditorViewDesktop(gtx layout.Context) layout.Dimensions {
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
}

func (pg *Restore) seedEditorViewMobile(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{
		Right: values.MarginPadding5,
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				return pg.inputsGroupMobile(gtx, pg.seedList, 33, 0)
			})
		}),
	)
}

func (pg *Restore) inputsGroup(gtx layout.Context, l *layout.List, len, startIndex int) layout.Dimensions {
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

func (pg *Restore) inputsGroupMobile(gtx layout.Context, l *layout.List, len, startIndex int) layout.Dimensions {
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return l.Layout(gtx, len, func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
							pg.layoutSeedMenu(gtx, i*1+startIndex)
							return pg.seedEditors.editors[i*1+startIndex].Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (pg *Restore) onSuggestionSeedsClicked() {
	index := pg.seedEditors.focusIndex
	if index != -1 {
		for i, b := range pg.seedMenu {
			for pg.seedMenu[i].button.Clicked() {
				pg.seedEditors.editors[index].Edit.Editor.SetText(b.text)
				pg.seedEditors.editors[index].Edit.Editor.MoveCaret(len(b.text), 0)
				pg.seedClicked = true
				if index != numberOfSeeds {
					pg.seedEditors.editors[index+1].Edit.Editor.Focus()
				}

				if index == numberOfSeeds {
					pg.isLastEditor = true
				}
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

func (pg *Restore) editorSeedsEventsHandler() {
	var focused []int

	seedEvent := func(i int, text string) {
		if pg.seedClicked {
			pg.seedEditors.focusIndex = -1
			pg.seedClicked = false
		} else {
			pg.seedEditors.focusIndex = i
		}

		if text == "" {
			pg.isLastEditor = false
			pg.openPopupIndex = -1
		} else {
			pg.openPopupIndex = i
		}

		if i != numberOfSeeds {
			pg.isLastEditor = false
		}
	}

	for i := 0; i < len(pg.seedEditors.editors); i++ {
		editor := &pg.seedEditors.editors[i]
		text := editor.Edit.Editor.Text()

		if editor.Edit.Editor.Focused() {
			seedEvent(i, text)
			focused = append(focused, i)
		}

		for _, e := range editor.Edit.Editor.Events() {
			switch e.(type) {
			case widget.ChangeEvent:
				seedEvent(i, text)

			case widget.SubmitEvent:
				if i != numberOfSeeds {
					pg.seedEditors.editors[i+1].Edit.Editor.Focus()
					pg.selected = 0
				}

				if i == numberOfSeeds {
					pg.selected = 0
					pg.isLastEditor = true
				}
			}
		}
	}

	if len(diff(pg.focused, focused)) > 0 {
		pg.seedEditors.focusIndex = -1
	}
	pg.focused = focused
}

func (pg *Restore) initSeedMenu() {
	for i := 0; i < pg.suggestionLimit; i++ {
		btn := pg.Theme.Button("")
		btn.Background, btn.Color = color.NRGBA{}, pg.Theme.Color.Text
		pg.seedMenu = append(pg.seedMenu, seedItemMenu{
			text:   "",
			button: btn,
		})
	}
}

func (pg *Restore) suggestionSeedEffect() {
	for k := range pg.suggestions {
		if pg.selected == k || pg.seedMenu[k].button.Hovered() {
			pg.seedMenu[k].button.Background = pg.Theme.Color.Gray4
		} else {
			pg.seedMenu[k].button.Background = color.NRGBA{}
		}
	}
}

func (pg *Restore) layoutSeedMenu(gtx layout.Context, optionsSeedMenuIndex int) {
	if pg.openPopupIndex != optionsSeedMenuIndex || pg.openPopupIndex != pg.seedEditors.focusIndex ||
		pg.isLastEditor {
		return
	}

	inset := layout.Inset{
		Top:  values.MarginPadding35,
		Left: values.MarginPadding0,
	}

	m := op.Record(gtx.Ops)
	_, caretPos := pg.seedEditors.editors[pg.seedEditors.focusIndex].Edit.Editor.CaretPos()
	inset.Layout(gtx, func(gtx C) D {
		border := widget.Border{Color: pg.Theme.Color.Gray4, CornerRadius: values.MarginPadding5, Width: values.MarginPadding2}
		if !pg.seedEditorChanged() && pg.nextcurrentCaretPosition != caretPos {
			pg.nextcurrentCaretPosition = -1
			return border.Layout(gtx, func(gtx C) D {
				return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.seedMenu), func(gtx C, i int) D {
						return layout.UniformInset(values.MarginPadding0).Layout(gtx, pg.seedMenu[i].button.Layout)
					})
				})
			})
		}
		return D{}
	})
	op.Defer(gtx.Ops, m.Stop())
}

func (pg Restore) suggestionSeeds(text string) []string {
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

func (pg *Restore) updateSeedResetBtn() bool {
	for _, editor := range pg.seedEditors.editors {
		return editor.Edit.Editor.Text() != ""
	}
	return false
}

func (pg *Restore) validateSeeds() (bool, string) {
	focus := pg.seedEditors.focusIndex
	seedMatchCounter := 0
	seedPhrase := ""
	for j := 0; j < len(pg.allSuggestions); j++ {
		if focus != -1 {
			if pg.seedEditors.editors[pg.seedEditors.focusIndex].Edit.Editor.Text() == pg.allSuggestions[j] {
				seedMatchCounter = 1
			}
		}
	}

	for i, editor := range pg.seedEditors.editors {
		if editor.Edit.Editor.Text() == "" || seedMatchCounter == 0 {
			pg.seedEditors.editors[i].Edit.HintColor = pg.Theme.Color.Danger
			return false, ""
		}

		seedPhrase += editor.Edit.Editor.Text() + " "
	}

	return true, seedPhrase
}

func (pg *Restore) verifySeeds() bool {
	isValid, seedphrase := pg.validateSeeds()
	pg.seedPhrase = ""

	if isValid {
		pg.seedPhrase = seedphrase
		if !dcrlibwallet.VerifySeed(pg.seedPhrase) {
			pg.Toast.NotifyError("invalid seed phrase")
			return false
		}
	}

	// Compare seed with existing wallets seed. On positive match abort import
	// to prevent duplicate wallet. walletWithSameSeed >= 0 if there is a match.
	walletWithSameSeed, err := pg.WL.MultiWallet.WalletWithSeed(pg.seedPhrase)
	if err != nil {
		log.Error(err)
		return false
	}

	if walletWithSameSeed != -1 {
		pg.Toast.NotifyError("A wallet with an identical seed already exists.")
		return false
	}

	return true
}

func (pg *Restore) resetSeeds() {
	for i := 0; i < len(pg.seedEditors.editors); i++ {
		pg.seedEditors.editors[i].Edit.Editor.SetText("")
	}
}

func switchSeedEditors(editors []decredmaterial.RestoreEditor) {
	for i := 0; i < len(editors); i++ {
		if editors[i].Edit.Editor.Focused() {
			if i == len(editors)-1 {
				editors[0].Edit.Editor.Focus()
			} else {
				editors[i+1].Edit.Editor.Focus()
			}
		}
	}
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *Restore) HandleUserInteractions() {
	focus := pg.seedEditors.focusIndex
	if focus != -1 {
		pg.suggestions = pg.suggestionSeeds(pg.seedEditors.editors[focus].Edit.Editor.Text())
		pg.seedMenu = pg.seedMenu[:len(pg.suggestions)]
		if !pg.seedEditorChanged() {
			for k, s := range pg.suggestions {
				pg.seedMenu[k].text, pg.seedMenu[k].button.Text = s, s
			}
		}
	}

	for pg.backButton.Button.Clicked() {
		pg.PopWindowPage()
	}

	for pg.validateSeed.Clicked() {
		if !pg.verifySeeds() {
			return
		}

		pg.isRestoring = true
		modal.NewCreatePasswordModal(pg.Load).
			Title("Enter wallet details").
			EnableName(true).
			ShowWalletInfoTip(true).
			SetParent(pg).
			PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					_, err := pg.WL.MultiWallet.RestoreWallet(walletName, pg.seedPhrase, password, dcrlibwallet.PassphraseTypePass)
					if err != nil {
						m.SetError(components.TranslateErr(err))
						m.SetLoading(false)
						pg.isRestoring = false
						return
					}

					pg.Toast.Notify("Wallet restored")
					pg.resetSeeds()
					m.Dismiss()
					// Close this page and return to the previous page (most likely wallets page)
					// if there's no restoreComplete callback function.
					if pg.restoreComplete == nil {
						pg.PopWindowPage()
					} else {
						pg.restoreComplete()
					}
				}()
				return false
			}).Show()
	}

	for pg.resetSeedFields.Clicked() {
		pg.resetSeeds()
		pg.seedEditors.focusIndex = -1
	}

	pg.editorSeedsEventsHandler()
	pg.onSuggestionSeedsClicked()
	pg.suggestionSeedEffect()

	if pg.seedEditorChanged() {
		pg.suggestions = nil
		_, caretPos := pg.seedEditors.editors[pg.seedEditors.focusIndex].Edit.Editor.CaretPos()
		pg.currentCaretPosition = caretPos
		pg.nextcurrentCaretPosition = caretPos
	}

	if pg.currentCaretPositionChanged() {
		pg.selected = 0
	}
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (pg *Restore) HandleKeyEvent(evt *key.Event) {
	// if pg.isRestoring {
	// 	return
	// }
	// if evt.Name == key.NameTab && evt.Modifiers != key.ModShift && evt.State == key.Press {
	// 	if len(pg.suggestions) > 0 {
	// 		focus := pg.seedEditors.focusIndex
	// 		pg.seedEditors.editors[focus].Edit.Editor.SetText(pg.suggestions[pg.selected])
	// 		pg.seedClicked = true
	// 		pg.seedEditors.editors[focus].Edit.Editor.MoveCaret(len(pg.suggestions[pg.selected]), -1)
	// 	}
	// 	switchSeedEditors(pg.seedEditors.editors)
	// }
	// if evt.Name == key.NameTab && evt.Modifiers == key.ModShift && evt.State == key.Press {
	// 	for i := 0; i < len(pg.seedEditors.editors); i++ {
	// 		if pg.seedEditors.editors[i].Edit.Editor.Focused() {
	// 			if i == 0 {
	// 				pg.seedEditors.editors[len(pg.seedEditors.editors)-1].Edit.Editor.Focus()
	// 			} else {
	// 				pg.seedEditors.editors[i-1].Edit.Editor.Focus()
	// 			}
	// 		}
	// 	}
	// }
	// if evt.Name == key.NameUpArrow && pg.openPopupIndex != -1 && evt.State == key.Press {
	// 	pg.selected--
	// 	if pg.selected < 0 {
	// 		pg.selected = 0
	// 	}
	// }
	// if evt.Name == key.NameDownArrow && pg.openPopupIndex != -1 && evt.State == key.Press {
	// 	pg.selected++
	// 	if pg.selected >= len(pg.suggestions) {
	// 		pg.selected = len(pg.suggestions) - 1
	// 	}
	// }
	// if (evt.Name == key.NameReturn || evt.Name == key.NameEnter) && pg.openPopupIndex != -1 && evt.State == key.Press && len(pg.suggestions) != 0 {
	// 	if pg.seedEditors.focusIndex == -1 && len(pg.suggestions) == 1 {
	// 		return
	// 	}

	// 	pg.seedMenu[pg.selected].button.Click()
	// }
}

func (pg *Restore) HandleEvent(gtx C) {
	// declare the needed keys
	key.InputOp{Tag: pg, Keys: "(ModShift)-[⌤, ⏎]| [↑,↓] | (ModShift)-Tab"}.Add(gtx.Ops)

	// range through the context events and call handleKeyEvent
	for _, evt := range gtx.Events(pg) {
		if e, ok := evt.(key.Event); ok && evt.State == key.Press {
			switch {
			case e.Name == key.NameTab && e.Modifiers == 0:
				fmt.Println(e.Name, "1")
			case e.Name == key.NameTab && e.Modifiers == key.ModShift:
				fmt.Println(e.Name, "2")
			case e.Name == key.NameUpArrow && e.Modifiers == 0:
				fmt.Println(e.Name, "3")
			case e.Name == key.NameDownArrow && e.Modifiers == 0:
				fmt.Println(e.Name, "4")
			case e.Name == key.NameLeftArrow && e.Modifiers == 0:
				fmt.Println(e.Name, "1")
			case e.Name == key.NameRightArrow && e.Modifiers == 0:
				fmt.Println(e.Name, "1")
			default:
			}
		}

		// switch evt := evt.(type) {
		// case key.Event:
		// 	if pg.isRestoring {
		// 		return
		// 	}
		// 	if evt.Name == key.NameTab && evt.Modifiers != key.ModShift && evt.State == key.Press {
		// 		if len(pg.suggestions) > 0 {
		// 			focus := pg.seedEditors.focusIndex
		// 			pg.seedEditors.editors[focus].Edit.Editor.SetText(pg.suggestions[pg.selected])
		// 			pg.seedClicked = true
		// 			pg.seedEditors.editors[focus].Edit.Editor.MoveCaret(len(pg.suggestions[pg.selected]), -1)
		// 		}
		// 		switchSeedEditors(pg.seedEditors.editors)
		// 	}
		// 	if evt.Name == key.NameTab && evt.Modifiers == key.ModShift && evt.State == key.Press {
		// 		for i := 0; i < len(pg.seedEditors.editors); i++ {
		// 			if pg.seedEditors.editors[i].Edit.Editor.Focused() {
		// 				if i == 0 {
		// 					pg.seedEditors.editors[len(pg.seedEditors.editors)-1].Edit.Editor.Focus()
		// 				} else {
		// 					pg.seedEditors.editors[i-1].Edit.Editor.Focus()
		// 				}
		// 			}
		// 		}
		// 	}
		// 	if evt.Name == key.NameUpArrow && pg.openPopupIndex != -1 && evt.State == key.Press {
		// 		pg.selected--
		// 		if pg.selected < 0 {
		// 			pg.selected = 0
		// 		}
		// 	}
		// 	if evt.Name == key.NameDownArrow && pg.openPopupIndex != -1 && evt.State == key.Press {
		// 		pg.selected++
		// 		if pg.selected >= len(pg.suggestions) {
		// 			pg.selected = len(pg.suggestions) - 1
		// 		}
		// 	}
		// 	if (evt.Name == key.NameReturn || evt.Name == key.NameEnter) && pg.openPopupIndex != -1 && evt.State == key.Press && len(pg.suggestions) != 0 {
		// 		if pg.seedEditors.focusIndex == -1 && len(pg.suggestions) == 1 {
		// 			return
		// 		}

		// 		pg.seedMenu[pg.selected].button.Click()
		// 	}
		// }
	}
}

func (pg *Restore) currentCaretPositionChanged() bool {
	focus := pg.seedEditors.focusIndex
	if !pg.seedEditorChanged() {
		if focus == -1 {
			return false
		}
		_, caretPos := pg.seedEditors.editors[pg.seedEditors.focusIndex].Edit.Editor.CaretPos()
		if pg.currentCaretPosition != caretPos {
			pg.currentCaretPosition = caretPos
			return true
		}
	}

	return false
}

func (pg *Restore) seedEditorChanged() bool {
	focus := pg.seedEditors.focusIndex
	if pg.selectedSeedEditor != focus {
		if focus == -1 {
			return false
		}
		pg.selectedSeedEditor = focus
		return true
	}

	return false
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *Restore) OnNavigatedFrom() {}
