package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type modalTemplate struct {
	walletName    decredmaterial.Editor
	password      decredmaterial.Editor
	matchPassword decredmaterial.Editor
	confirm       decredmaterial.Button
	cancel        decredmaterial.Button
}

type modalLoad struct {
	title    string
	confirm  interface{}
	cancel   interface{}
}

func (win *Window) LoadTemplates(th *decredmaterial.Theme) *modalTemplate {
	return &modalTemplate{
		confirm:       th.Button(new(widget.Clickable), "some text"),
		cancel:        th.Button(new(widget.Clickable), "some text cancel"),
		walletName:    th.Editor(new(widget.Editor), "Wallet name"),
		password:      th.Editor(new(widget.Editor), "Password"),
		matchPassword: th.Editor(new(widget.Editor), "Matching password"),
	}
}

func (m *modalTemplate) createNewWallet(th *decredmaterial.Theme) []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return th.H6("Create new wallet").Layout(gtx)
		},
		func(gtx C) D {
			separator := th.Line()
			separator.Width = gtx.Constraints.Max.X
			return separator.Layout(gtx)
		},
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			m.password.Editor.Mask, m.password.Editor.SingleLine = '*', true
			return m.password.Layout(gtx)
		},
		func(gtx C) D {
			m.matchPassword.Editor.Mask, m.matchPassword.Editor.SingleLine = '*', true
			return m.matchPassword.Layout(gtx)
		},
	}
}

func (m *modalTemplate) Layout(th *decredmaterial.Theme, load *modalLoad) []func(gtx C) D {
	var w []func(gtx C) D

	switch t := load.confirm.(type) {
	case func(string, string):
		w = m.createNewWallet(th)
		if m.confirm.Button.Clicked() {
			t(m.walletName.Editor.Text(), m.password.Editor.Text())
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
	}

	action := []func(gtx C) D{
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							return m.confirm.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							return m.cancel.Layout(gtx)
						})
					}),
				)
			})
		},
	}

	w = append(w, action...)
	return w
}
