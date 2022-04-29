package seedbackup

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const BackupInstructionsPageID = "backup_instructions"

type (
	C = layout.Context
	D = layout.Dimensions
)

type BackupInstructionsPage struct {
	*app.App
	wallet         *dcrlibwallet.Wallet
	changeFragment func(app.Page)

	backButton  decredmaterial.IconButton
	viewSeedBtn decredmaterial.Button
	checkBoxes  []decredmaterial.CheckBoxStyle
	infoList    *layout.List
}

func NewBackupInstructionsPage(app *app.App, wallet *dcrlibwallet.Wallet, changeFragment func(app.Page)) *BackupInstructionsPage {
	bi := &BackupInstructionsPage{
		App:            app,
		wallet:         wallet,
		changeFragment: changeFragment,

		viewSeedBtn: app.Theme.Button("View seed phrase"),
	}

	bi.viewSeedBtn.Font.Weight = text.Medium

	bi.backButton, _ = components.SubpageHeaderButtons(app.Theme)
	bi.backButton.Icon = app.Theme.Icons.ContentClear

	bi.checkBoxes = []decredmaterial.CheckBoxStyle{
		app.Theme.CheckBox(new(widget.Bool), "The 33-word seed word is EXTREMELY IMPORTANT."),
		app.Theme.CheckBox(new(widget.Bool), "seed word is the only way to restore your wallet."),
		app.Theme.CheckBox(new(widget.Bool), "It is recommended to store your seed word in a physical format (e.g. write down on a paper)."),
		app.Theme.CheckBox(new(widget.Bool), "It is highly discouraged to store your seed word in any digital format (e.g. screenshot)."),
		app.Theme.CheckBox(new(widget.Bool), "Anyone with your seed word can steal your funds. DO NOT show it to anyone."),
	}

	for i := range bi.checkBoxes {
		bi.checkBoxes[i].TextSize = values.TextSize16
	}

	bi.infoList = &layout.List{Axis: layout.Vertical}

	return bi
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) ID() string {
	return BackupInstructionsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) OnNavigatedTo() {

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) HandleUserInteractions() {
	for pg.viewSeedBtn.Clicked() {
		if pg.verifyCheckBoxes() {
			// TODO: Will repeat the paint cycle, just queue the next fragment to be displayed
			pg.changeFragment(NewSaveSeedPage(pg.App, pg.wallet, pg.changeFragment))
		}
	}
}

func promptToExit(app *app.App) {
	var PopToFragment func(string)
	modal.NewInfoModal(app).
		Title("Exit?").
		Body("Are you sure you want to exit the seed backup process?").
		NegativeButton("No", func() {}).
		PositiveButton("Yes", func(isChecked bool) {
			PopToFragment(components.WalletsPageID) // TODO: Will crash
		}).
		Show()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) Layout(gtx layout.Context) layout.Dimensions {
	sp := components.SubPage{
		App:        pg.App,
		Title:      "Keep in mind",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.App)
		},
		Body: func(gtx C) D {
			return pg.infoList.Layout(gtx, len(pg.checkBoxes), func(gtx C, i int) D {
				return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, pg.checkBoxes[i].Layout)
			})
		},
	}

	pg.viewSeedBtn.SetEnabled(pg.verifyCheckBoxes())

	return container(gtx, *pg.Theme, sp.Layout, "", pg.viewSeedBtn)
}

func (pg *BackupInstructionsPage) verifyCheckBoxes() bool {
	for _, cb := range pg.checkBoxes {
		if !cb.CheckBox.Value {
			return false
		}
	}
	return true
}

func container(gtx C, theme decredmaterial.Theme, body layout.Widget, infoText string, actionBtn decredmaterial.Button) D {
	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return body(gtx)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = gtx.Constraints.Max
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.WrapContent,
					Orientation: layout.Vertical,
					Direction:   layout.S,
					Alignment:   layout.Baseline,
					Background:  theme.Color.Surface,
					Shadow:      theme.Shadow(),
					Padding:     layout.UniformInset(values.MarginPadding16),
					Margin:      layout.Inset{Left: values.Size0_5},
					Border:      decredmaterial.Border{Radius: decredmaterial.Radius(4)},
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if !components.StringNotEmpty(infoText) {
							return D{}
						}
						label := theme.Label(values.TextSize14, infoText)
						label.Color = theme.Color.GrayText1
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, label.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return actionBtn.Layout(gtx)
					}))
			}),
		)
	})
}
