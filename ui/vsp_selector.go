package ui

import (
	"fmt"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/wallet"
	"image"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type vspSelector struct {
	*pageCommon

	dialogTitle string

	changed   bool
	showVSPModal        *widget.Clickable
	selectedVSP         wallet.VSPInfo
}

func newVSPSelector(common *pageCommon) *vspSelector {
	v := &vspSelector{
		pageCommon: common,
		showVSPModal: new(widget.Clickable),
	}
	return v
}

func (v *vspSelector) title(title string) *vspSelector {
	v.dialogTitle = title
	return v
}

func (v *vspSelector) Changed() bool {
	changed := v.changed
	v.changed = false
	return changed
}

func (v *vspSelector) selectVSP(vspHost string) {
	for _, vsp := range (*v.vspInfo).List {
		if vsp.Host == vspHost {
			v.changed = true
			v.selectedVSP = vsp
			break
		}
	}
}

func (v *vspSelector) SelectedVSP() wallet.VSPInfo {
	return v.selectedVSP
}

func (v *vspSelector) handle() {
	if v.showVSPModal.Clicked() {
		newVSPSelectorModal(v.pageCommon).
			title("Voting service provider").
			vspSelected(func(info wallet.VSPInfo) {
				v.selectVSP(info.Host)
			}).
			Show()
	}
}

func (v *vspSelector) Layout(gtx layout.Context) layout.Dimensions {
	v.handle()

	border := widget.Border{
		Color:        v.theme.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, v.showVSPModal, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if v.selectedVSP.Host == "" {
							txt := v.theme.Label(values.TextSize16, "Select VSP...")
							txt.Color = v.theme.Color.Gray2
							return txt.Layout(gtx)
						}
						return v.theme.Label(values.TextSize16, v.selectedVSP.Host).Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									if v.selectedVSP.Info == nil {
										return layout.Dimensions{}
									}
									txt := v.theme.Label(values.TextSize16, fmt.Sprintf("%v%%", v.selectedVSP.Info.FeePercentage))
									txt.Color = v.theme.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return v.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
									})
								}),
							)
						})
					}),
				)
			})
		})
	})
}

const VSPSelectorModalID = "VSPSelectorModal"

type vspSelectorModal struct {
	*pageCommon
	dialogTitle string

	modal    decredmaterial.Modal
	inputVSP decredmaterial.Editor
	addVSP   decredmaterial.Button

	vspHosts    *layout.List
	selectVSP   []*gesture.Click
	selectedVSP wallet.VSPInfo

	vspSelectedCallback func(wallet.VSPInfo)
}

func newVSPSelectorModal(common *pageCommon) *vspSelectorModal {
	v := &vspSelectorModal{
		pageCommon: common,

		inputVSP: common.theme.Editor(new(widget.Editor), "Add a new VSP..."),
		addVSP:   common.theme.Button(new(widget.Clickable), "Save"),
		vspHosts: &layout.List{Axis: layout.Vertical},
		modal:    *common.theme.ModalFloatTitle(),
	}

	return v
}

func (v *vspSelectorModal) OnResume() {

}

func (v *vspSelectorModal) modalID() string {
	return VSPSelectorModalID
}

func (v *vspSelectorModal) Show() {
	v.showModal(v)
}

func (v *vspSelectorModal) Dismiss() {
	v.dismissModal(v)
}

func (v *vspSelectorModal) handle() {
	if v.editorsNotEmpty(&v.addVSP, v.inputVSP.Editor) && v.addVSP.Button.Clicked() {
		go func() {
			err := v.AddVSP(v.inputVSP.Editor.Text())
			if err != nil {
				v.notify(err.Error(), false)
			} else {
				v.inputVSP.Editor.SetText("")
			}
		}()
	}

	vspList := (*v.vspInfo).List
	if len(vspList) != len(v.selectVSP) {
		v.selectVSP = createClickGestures(len(vspList))
	}
}

func (v *vspSelectorModal) title(title string) *vspSelectorModal {
	v.dialogTitle = title
	return v
}

func (v *vspSelectorModal) vspSelected(callback func(wallet.VSPInfo)) *vspSelectorModal {
	v.vspSelectedCallback = callback
	v.Dismiss()
	return v
}

func (v *vspSelectorModal) OnDismiss() {}

func (v *vspSelectorModal) Layout(gtx layout.Context) layout.Dimensions {
	return v.modal.Layout(gtx, []layout.Widget{
		func(gtx C) D {
			return v.theme.Label(values.TextSize20, v.dialogTitle).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := v.theme.Label(values.TextSize14, "Address")
					txt.Color = v.theme.Color.Gray2
					txtFee := v.theme.Label(values.TextSize14, "Fee")
					txtFee.Color = v.theme.Color.Gray2
					return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
						return endToEndRow(gtx, txt.Layout, txtFee.Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					listVSP := (*v.vspInfo).List
					return v.vspHosts.Layout(gtx, len(v.selectVSP), func(gtx C, i int) D {
						click := v.selectVSP[i]
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
						click.Add(gtx.Ops)
						v.handlerSelectVSP(click.Events(gtx), listVSP[i])

						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(0.8, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
									txt := v.theme.Label(values.TextSize14, fmt.Sprintf("%v", listVSP[i].Info.FeePercentage)+"%")
									txt.Color = v.theme.Color.Gray2
									return endToEndRow(gtx, v.theme.Label(values.TextSize16, listVSP[i].Host).Layout, txt.Layout)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if v.selectedVSP.Host != listVSP[i].Host {
									return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
										return layout.Dimensions{}
									})
								}
								return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
									return v.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
								})
							}),
						)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, v.inputVSP.Layout),
				layout.Rigid(v.addVSP.Layout),
			)
		},
	}, 900)
}

func (v *vspSelectorModal) handlerSelectVSP(events []gesture.ClickEvent, info wallet.VSPInfo) {
	for _, e := range events {
		if e.Type == gesture.TypeClick {
			v.selectedVSP = info
			v.vspSelectedCallback(info)
			v.Dismiss()
		}
	}
}

func (v *vspSelectorModal) editorsNotEmpty(btn *decredmaterial.Button, editors ...*widget.Editor) bool {
	btn.Color = v.theme.Color.Surface
	for _, e := range editors {
		if e.Text() == "" {
			btn.Background = v.theme.Color.Hint
			return false
		}
	}

	btn.Background = v.theme.Color.Primary
	return true
}
