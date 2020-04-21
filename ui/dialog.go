package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

type dialogWidgets struct {
	modal                     decredmaterial.Modal
	active                    bool
	dialog                    layout.Widget
	cancel, confirm           widget.Button
	cancelW                   decredmaterial.IconButton
	confirmW                  decredmaterial.Button
	password, matchPassword   widget.Editor
	passwordW, matchPasswordW decredmaterial.Editor
}

func newDialogWidgets(common pageCommon) *dialogWidgets {
	return &dialogWidgets{
		cancelW:   common.theme.PlainIconButton(common.icons.contentClear),
		passwordW: common.theme.Editor("Enter password"),
		confirmW:  common.theme.Button("Confirm"),
	}
}

func (wdgs *dialogWidgets) LayoutIfActive(gtx *layout.Context, body layout.Widget) {
	if wdgs.active {
		wdgs.modal.Layout(gtx, body, wdgs.dialog)
	} else {
		body()
	}
	if wdgs.cancel.Clicked(gtx) {
		wdgs.active = false
	}
}

func (wdgs *dialogWidgets) SetDialog(diag layout.Widget) {
	wdgs.active = true
	wdgs.dialog = diag
}
