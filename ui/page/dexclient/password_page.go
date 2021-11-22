package dexclient

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DexPasswordPageID = "DexPassword"

type DexPasswordPage struct {
	*load.Load
	createPassword                decredmaterial.Button
	appPassword, appPasswordAgain decredmaterial.Editor
	isSending                     bool
}

func NewDexPasswordPage(l *load.Load) *DexPasswordPage {
	pg := &DexPasswordPage{
		Load:             l,
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "Password"),
		appPasswordAgain: l.Theme.EditorPassword(new(widget.Editor), "Password again"),
		createPassword:   l.Theme.Button("Create password"),
	}

	pg.createPassword.TextSize = values.TextSize12
	pg.createPassword.Background = l.Theme.Color.Primary
	pg.appPassword.Editor.SingleLine = true
	pg.appPasswordAgain.Editor.SingleLine = true

	return pg
}

func (pg *DexPasswordPage) ID() string {
	return DexPasswordPageID
}

func (pg *DexPasswordPage) OnClose() {}

func (pg *DexPasswordPage) OnResume() {
}

func (pg *DexPasswordPage) Handle() {
	if pg.createPassword.Button.Clicked() {
		if pg.appPasswordAgain.Editor.Text() != pg.appPassword.Editor.Text() || pg.isSending {
			return
		}

		pg.isSending = true
		go func() {
			err := pg.Dexc().InitializeWithPassword([]byte(pg.appPassword.Editor.Text()))
			pg.isSending = false
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}
			pg.ChangeFragment(NewAddDexPage(pg.Load))
		}()
	}
}

func (pg *DexPasswordPage) Layout(gtx layout.Context) D {
	body := func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Label(values.TextSize20, "Set App Password").Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Label(values.TextSize14, "Set your app password. This password will protect your DEX account keys and connected wallets.").Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.appPassword.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.appPasswordAgain.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.createPassword.Layout(gtx)
						})
					}),
				)
			})
		})
	}

	return components.UniformPadding(gtx, body)
}
