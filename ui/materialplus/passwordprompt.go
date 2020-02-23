package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/layouts"
	"github.com/raedahgroup/godcr-gio/ui/styles"
)

type PasswordDialog struct {
	layouts.Dialog
	editor    *widget.Editor
	editorWdg material.Editor
	prompt    material.Label
	confirmed bool
}

func (t *Theme) PasswordDialog(prompt string) PasswordDialog {
	return PasswordDialog{
		editorWdg: t.Editor("Enter password"),
		prompt:    t.H3(prompt),
		Dialog: layouts.Dialog{
			Background:    t.Background,
			Confirm:       new(widget.Button),
			Cancel:        new(widget.Button),
			ConfirmButton: t.Button("Confirm"),
			CancelButton:  t.DangerButton("Cancel"),
		},
	}
}

func (diag PasswordDialog) Layout(gtx *layout.Context, editor *widget.Editor) {
	editor.SingleLine = true
	diag.Dialog.Layout(gtx, func() {
		layout.Flex{}.Layout(gtx,
			layouts.FlexedWithStyle(gtx, styles.Centered, 0.5, func() {
				diag.prompt.Layout(gtx)
			}),
			layouts.FlexedWithStyle(gtx, styles.Centered, 0.5, func() {
				diag.editorWdg.Layout(gtx, editor)
			}),
		)
	})
}

func (diag PasswordDialog) Clicked(gtx *layout.Context) (bool, bool) {
	return diag.Confirm.Clicked(gtx), diag.Cancel.Clicked(gtx)
}

func (diag PasswordDialog) Password() *string {
	if diag.confirmed {
		text := diag.editor.Text()
		return &text
	}
	return nil
}
