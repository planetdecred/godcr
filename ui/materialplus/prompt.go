package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/icons"
)

type PasswordDialog struct {
	Password material.Editor
	Prompt   material.Label
	ConfirmCancel
}

func (t *Theme) PasswordDialog(prompt string) PasswordDialog {
	return PasswordDialog{
		Password: t.Editor("password"),
		Prompt:   t.Label(unit.Dp(20), prompt),
		ConfirmCancel: ConfirmCancel{
			Confirm: t.Button("Confirm"),
			Cancel:  t.IconButton(icons.ContentAdd),
		},
	}
}

func (diag PasswordDialog) Layout(gtx *layout.Context, editor *widget.Editor, confirm, cancel *widget.Button) {
	body := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() { diag.Prompt.Layout(gtx) }),
			layout.Rigid(func() { diag.Password.Layout(gtx, editor) }),
		)
	}
	diag.ConfirmCancel.Body = body
	diag.ConfirmCancel.Layout(gtx, confirm, cancel)
}

func (diag PasswordDialog) LayoutWithMatch(gtx *layout.Context, editor, match *widget.Editor, confirm, cancel *widget.Button) {
	body := func() {

	}
	diag.ConfirmCancel.Body = body
	diag.ConfirmCancel.Layout(gtx, confirm, cancel)
}
