package ui

import (
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const CreateWalletTemplate = "CreateWallet"

type ModalTemplate struct {
}

func createNewWallet(th *decredmaterial.Theme) []func(gtx C) D {
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
			password := th.Editor(new(widget.Editor), "Enter password")
			password.Editor.Mask, password.Editor.SingleLine = '*', true
			return password.Layout(gtx)
		},
		func(gtx C) D {
			matchingPassword := th.Editor(new(widget.Editor), "Enter password")
			matchingPassword.Editor.Mask, matchingPassword.Editor.SingleLine = '*', true
			return matchingPassword.Layout(gtx)
		},
	}
}

func modalLayout(th *decredmaterial.Theme, template string) []func(gtx C) D {
	switch template {
	case CreateWalletTemplate:
		return createNewWallet(th)
	}
	return []func(gtx C) D{}
}
