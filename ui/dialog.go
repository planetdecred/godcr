package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"gioui.org/widget/material"
)

const (
	dialogSize        = 0.3
	confirmCancelSize = 0.3
)

// Dialog is a convenient struct for presenting a dialog modal.
type Dialog struct {
	ConfirmButton, CancelButton material.Button
	Confirm, Cancel             *widget.Button
	Active                      bool
}

// Layout renders the modal if Active is true.
// Blocks input behind the modal.
// If either Confirm or Cancel is nil, the corresponding Button is not rendered.
func (diag Dialog) Layout(gtx *layout.Context, dialog func()) {
	if !diag.Active {
		return
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1-dialogSize, func() {
			FillWithColor(gtx, FadedColor, true)
		}),
		layout.Flexed(dialogSize, func() {
			FillWithColor(gtx, WhiteColor, true)

			children := make([]layout.FlexChild, 0, 2)

			if diag.Confirm != nil {
				children = append(children,
					layout.Rigid(func() {
						layout.Align(layout.Center).Layout(gtx, func() {
							diag.ConfirmButton.Layout(gtx, diag.Confirm)
						})
					}))
			}

			if diag.Cancel != nil {
				children = append(children,
					layout.Rigid(func() {
						layout.Align(layout.Center).Layout(gtx, func() {
							diag.CancelButton.Layout(gtx, diag.Cancel)
						})
					}))
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1-confirmCancelSize, dialog),
				layout.Flexed(confirmCancelSize, func() {
					layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx, children...)
				}),
			)
		}),
	)
}
