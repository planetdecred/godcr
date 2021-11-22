package dexclient

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DexLoginPageID = "DexLogin"

type DexLoginPage struct {
	*load.Load
	submit      decredmaterial.Button
	appPassword decredmaterial.Editor
}

func NewDexLoginPage(l *load.Load) *DexLoginPage {
	pg := &DexLoginPage{
		Load:        l,
		submit:      l.Theme.Button("Login"),
		appPassword: l.Theme.EditorPassword(new(widget.Editor), "App password"),
	}

	pg.submit.TextSize = values.TextSize12
	pg.submit.Background = l.Theme.Color.Primary
	pg.appPassword.Editor.SingleLine = true

	return pg
}

func (pg *DexLoginPage) ID() string {
	return DexLoginPageID
}

func (pg *DexLoginPage) OnClose() {}

func (pg *DexLoginPage) OnResume() {
}

func (pg *DexLoginPage) Handle() {
	if pg.submit.Button.Clicked() {
		err := pg.Dexc().Login([]byte(pg.appPassword.Editor.Text()))
		if err != nil {
			pg.Toast.NotifyError(err.Error())
			return
		}

		pg.ChangeFragment(NewMarketPage(pg.Load))
	}
}

func (pg *DexLoginPage) Layout(gtx layout.Context) D {
	body := func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							return pg.Theme.H6("Login").Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
									return pg.appPassword.Layout(gtx)
								})
							}),
							layout.Rigid(pg.submit.Layout),
						)
					}),
				)
			})
		})
	}

	return components.UniformPadding(gtx, body)
}
