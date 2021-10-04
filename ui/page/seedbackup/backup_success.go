package seedbackup

import (
	"gioui.org/layout"
	"gioui.org/text"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const BackupSuccessPageID = "backup_success"

type BackupSuccessPage struct {
	*load.Load
	actionButton decredmaterial.Button
}

func NewBackupSuccessPage(l *load.Load) *BackupSuccessPage {
	pg := &BackupSuccessPage{
		Load: l,

		actionButton: l.Theme.OutlineButton("Back to Wallets"),
	}
	pg.actionButton.Font.Weight = text.Medium

	return pg
}

func (pg *BackupSuccessPage) ID() string {
	return BackupSuccessPageID
}

func (pg *BackupSuccessPage) OnResume() {}

func (pg *BackupSuccessPage) Handle() {
	for pg.actionButton.Clicked() {
		pg.PopToFragment(components.WalletsPageID)
	}
}

func (pg *BackupSuccessPage) OnClose() {}

// - Layout

func (pg *BackupSuccessPage) Layout(gtx C) D {

	return components.UniformPadding(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.MatchParent,
			Orientation: layout.Vertical,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.MatchParent,
					Orientation: layout.Vertical,
					Alignment:   layout.Middle,
					Direction:   layout.Center,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						successIcon := pg.Icons.ActionCheckCircle
						return successIcon.Layout(gtx, values.MarginPadding64)
					}),
					layout.Rigid(func(gtx C) D {
						label := pg.Theme.Label(values.TextSize24, "Your seed phrase backup is verified")
						label.Color = pg.Theme.Color.DeepBlue

						return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, label.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						label := pg.Theme.Label(values.TextSize16, "Be sure to store your seed phrase backup in a secure location.")
						label.Color = pg.Theme.Color.Gray3

						return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, label.Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X

				return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, pg.actionButton.Layout)
			}),
		)
	})
}
