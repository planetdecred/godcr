package tickets

import (
	"fmt"
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/wallet"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type vspSelector struct {
	*load.Load

	dialogTitle string

	changed      bool
	showVSPModal *decredmaterial.Clickable
	selectedVSP  *wallet.VSPInfo
}

func newVSPSelector(l *load.Load) *vspSelector {
	v := &vspSelector{
		Load:         l,
		showVSPModal: l.Theme.NewClickable(true),
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
	for _, vsp := range v.WL.VspInfo.List {
		if vsp.Host == vspHost {
			v.changed = true
			v.selectedVSP = vsp
			break
		}
	}
}

func (v *vspSelector) SelectedVSP() *wallet.VSPInfo {
	return v.selectedVSP
}

func (v *vspSelector) handle() {
	if v.showVSPModal.Clicked() {
		newVSPSelectorModal(v.Load).
			title("Voting service provider").
			vspSelected(func(info *wallet.VSPInfo) {
				v.selectVSP(info.Host)
			}).
			Show()
	}
}

func (v *vspSelector) Layout(gtx layout.Context) layout.Dimensions {
	v.handle()

	border := widget.Border{
		Color:        v.Theme.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return v.showVSPModal.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if v.selectedVSP == nil {
							txt := v.Theme.Label(values.TextSize16, "Select VSP...")
							txt.Color = v.Theme.Color.Gray2
							return txt.Layout(gtx)
						}
						return v.Theme.Label(values.TextSize16, v.selectedVSP.Host).Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									if v.selectedVSP == nil {
										return layout.Dimensions{}
									}
									txt := v.Theme.Label(values.TextSize16, fmt.Sprintf("%v%%", v.selectedVSP.Info.FeePercentage))
									txt.Color = v.Theme.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										gtx.Constraints.Min.X = gtx.Px(values.MarginPadding20)
										return v.Icons.DropDownIcon.Layout(gtx, v.Theme.Color.Gray3)
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
	*load.Load

	dialogTitle string

	modal    decredmaterial.Modal
	inputVSP decredmaterial.Editor
	addVSP   decredmaterial.Button

	vspHosts    *layout.List
	selectVSP   []*gesture.Click
	selectedVSP *wallet.VSPInfo

	vspSelectedCallback func(*wallet.VSPInfo)
}

func newVSPSelectorModal(l *load.Load) *vspSelectorModal {
	v := &vspSelectorModal{
		Load: l,

		inputVSP: l.Theme.Editor(new(widget.Editor), "Add a new VSP..."),
		addVSP:   l.Theme.Button("Save"),
		vspHosts: &layout.List{Axis: layout.Vertical},
		modal:    *l.Theme.ModalFloatTitle(),
	}

	v.addVSP.SetEnabled(false)

	return v
}

func (v *vspSelectorModal) OnResume() {

}

func (v *vspSelectorModal) ModalID() string {
	return VSPSelectorModalID
}

func (v *vspSelectorModal) Show() {
	v.ShowModal(v)
}

func (v *vspSelectorModal) Dismiss() {
	v.DismissModal(v)
}

func (v *vspSelectorModal) Handle() {
	v.addVSP.SetEnabled(v.editorsNotEmpty(v.inputVSP.Editor))
	if v.addVSP.Clicked() {
		go func() {
			err := v.WL.AddVSP(v.inputVSP.Editor.Text())
			if err != nil {
				v.Toast.NotifyError(err.Error())
			} else {
				v.inputVSP.Editor.SetText("")
			}
		}()
	}

	vspList := v.WL.VspInfo.List
	if len(vspList) != len(v.selectVSP) {
		v.selectVSP = createClickGestures(len(vspList))
	}
}

func (v *vspSelectorModal) title(title string) *vspSelectorModal {
	v.dialogTitle = title
	return v
}

func (v *vspSelectorModal) vspSelected(callback func(*wallet.VSPInfo)) *vspSelectorModal {
	v.vspSelectedCallback = callback
	v.Dismiss()
	return v
}

func (v *vspSelectorModal) OnDismiss() {}

func (v *vspSelectorModal) Layout(gtx layout.Context) layout.Dimensions {
	return v.modal.Layout(gtx, []layout.Widget{
		func(gtx C) D {
			return v.Theme.Label(values.TextSize20, v.dialogTitle).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := v.Theme.Label(values.TextSize14, "Address")
					txt.Color = v.Theme.Color.Gray2
					txtFee := v.Theme.Label(values.TextSize14, "Fee")
					txtFee.Color = v.Theme.Color.Gray2
					return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
						return components.EndToEndRow(gtx, txt.Layout, txtFee.Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					listVSP := v.WL.VspInfo.List
					return v.vspHosts.Layout(gtx, len(listVSP), func(gtx C, i int) D {
						click := v.selectVSP[i]
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
						click.Add(gtx.Ops)
						v.handlerSelectVSP(click.Events(gtx), listVSP[i])

						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(0.8, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
									txt := v.Theme.Label(values.TextSize14, fmt.Sprintf("%v", listVSP[i].Info.FeePercentage)+"%")
									txt.Color = v.Theme.Color.Gray2
									return components.EndToEndRow(gtx, v.Theme.Label(values.TextSize16, listVSP[i].Host).Layout, txt.Layout)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if v.selectedVSP != nil || v.selectedVSP != listVSP[i] {
									return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
										return layout.Dimensions{}
									})
								}
								return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
									gtx.Constraints.Min.X = gtx.Px(values.MarginPadding20)
									return v.Icons.NavigationCheck.Layout(gtx, v.Theme.Color.Gray3)
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

func (v *vspSelectorModal) handlerSelectVSP(events []gesture.ClickEvent, info *wallet.VSPInfo) {
	for _, e := range events {
		if e.Type == gesture.TypeClick {
			v.selectedVSP = info
			v.vspSelectedCallback(info)
			v.Dismiss()
		}
	}
}

func (v *vspSelectorModal) editorsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}

	return true
}
