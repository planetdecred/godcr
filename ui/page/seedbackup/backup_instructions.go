package seedbackup

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
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
	*load.Load
	wallet *dcrlibwallet.Wallet

	backButton  decredmaterial.IconButton
	viewSeedBtn decredmaterial.Button
	checkBoxes  []decredmaterial.CheckBoxStyle
	infoList    *layout.List
}

func NewBackupInstructionsPage(l *load.Load, wallet *dcrlibwallet.Wallet) *BackupInstructionsPage {
	bi := &BackupInstructionsPage{
		Load:   l,
		wallet: wallet,

		viewSeedBtn: l.Theme.Button(new(widget.Clickable), "View seed phrase"),
	}

	bi.backButton, _ = components.SubpageHeaderButtons(l)
	bi.backButton.Icon = l.Icons.ContentClear

	bi.checkBoxes = []decredmaterial.CheckBoxStyle{
		l.Theme.CheckBox(new(widget.Bool), "The 33-word seed phrase is EXTREMELY IMPORTANT."),
		l.Theme.CheckBox(new(widget.Bool), "Seed phrase is the only way to restore your wallet."),
		l.Theme.CheckBox(new(widget.Bool), "It is recommended to store your seed phrase in a physical format (e.g. write down on a paper)."),
		l.Theme.CheckBox(new(widget.Bool), "It is highly discouraged to store your seed phrase in any digital format (e.g. screenshot)."),
		l.Theme.CheckBox(new(widget.Bool), "Anyone with your seed phrase can steal your funds. DO NOT show it to anyone."),
	}

	for i := range bi.checkBoxes {
		bi.checkBoxes[i].TextSize = values.TextSize16
	}

	bi.infoList = &layout.List{Axis: layout.Vertical}

	return bi
}

func (pg *BackupInstructionsPage) ID() string {
	return BackupInstructionsPageID
}

func (pg *BackupInstructionsPage) OnResume() {

}

func (pg *BackupInstructionsPage) Handle() {
	for pg.viewSeedBtn.Clicked() {
		if pg.verifyCheckBoxes() {
			pg.ChangeFragment(NewSaveSeedPage(pg.Load, pg.wallet))
		}
	}

}
func promptToExit(load *load.Load) {
	modal.NewInfoModal(load).
		Title("Exit?").
		Body("Are you sure you want to exit the seed backup process?").
		NegativeButton("No", func() {}).
		PositiveButton("Yes", func() {
			load.PopToFragment(components.WalletsPageID)
		}).
		Show()
}
func (pg *BackupInstructionsPage) OnClose() {}

// - Layout

func (pg *BackupInstructionsPage) Layout(gtx layout.Context) layout.Dimensions {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Keep in mind",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load)
		},
		Body: func(gtx C) D {
			return pg.infoList.Layout(gtx, len(pg.checkBoxes), func(gtx C, i int) D {
				return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, pg.checkBoxes[i].Layout)
			})
		},
	}

	if pg.verifyCheckBoxes() {
		pg.viewSeedBtn.Background = pg.Theme.Color.Primary
		pg.viewSeedBtn.Color = pg.Theme.Color.InvText
	} else {
		pg.viewSeedBtn.Background = pg.Theme.Color.InactiveGray
		pg.viewSeedBtn.Color = pg.Theme.Color.Text
	}

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
					Padding:     layout.UniformInset(values.MarginPadding16),
					Margin:      layout.Inset{Left: values.Size0_5},
					Border:      decredmaterial.Border{Radius: decredmaterial.Radius(4)},
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if !components.StringNotEmpty(infoText) {
							return D{}
						}
						label := theme.Label(values.TextSize14, infoText)
						label.Color = theme.Color.Gray3
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
