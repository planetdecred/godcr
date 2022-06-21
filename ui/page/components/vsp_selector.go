package components

import (
	"context"
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type VSPSelector struct {
	*load.Load

	dialogTitle string

	changed      bool
	showVSPModal *decredmaterial.Clickable
	selectedVSP  *dcrlibwallet.VSP
}

func NewVSPSelector(l *load.Load) *VSPSelector {
	v := &VSPSelector{
		Load:         l,
		showVSPModal: l.Theme.NewClickable(true),
	}
	return v
}

func (v *VSPSelector) Title(title string) *VSPSelector {
	v.dialogTitle = title
	return v
}

func (v *VSPSelector) Changed() bool {
	changed := v.changed
	v.changed = false
	return changed
}

func (v *VSPSelector) SelectVSP(vspHost string) {
	for _, vsp := range v.WL.MultiWallet.KnownVSPs() {
		if vsp.Host == vspHost {
			v.changed = true
			v.selectedVSP = vsp
			break
		}
	}
}

func (v *VSPSelector) SelectedVSP() *dcrlibwallet.VSP {
	return v.selectedVSP
}

func (v *VSPSelector) handle(window app.WindowNavigator) {
	if v.showVSPModal.Clicked() {
		modal := newVSPSelectorModal(v.Load).
			title(values.String(values.StrVotingServiceProvider)).
			vspSelected(func(info *dcrlibwallet.VSP) {
				v.SelectVSP(info.Host)
			})
		window.ShowModal(modal)
	}
}

func (v *VSPSelector) Layout(window app.WindowNavigator, gtx layout.Context) layout.Dimensions {
	v.handle(window)

	border := widget.Border{
		Color:        v.Theme.Color.Gray2,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return v.showVSPModal.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if v.selectedVSP == nil {
							txt := v.Theme.Label(values.TextSize16, values.String(values.StrSelectVSP))
							txt.Color = v.Theme.Color.GrayText3
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
									txt := v.Theme.Label(values.TextSize16, fmt.Sprintf("%v%%", v.selectedVSP.FeePercentage))
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										ic := decredmaterial.NewIcon(v.Theme.Icons.DropDownIcon)
										ic.Color = v.Theme.Color.Gray1
										return ic.Layout(gtx, values.MarginPadding20)
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

type vspSelectorModal struct {
	*load.Load
	*decredmaterial.Modal

	dialogTitle string

	inputVSP decredmaterial.Editor
	addVSP   decredmaterial.Button

	selectedVSP *dcrlibwallet.VSP
	vspList     *decredmaterial.ClickableList

	vspSelectedCallback func(*dcrlibwallet.VSP)
}

func newVSPSelectorModal(l *load.Load) *vspSelectorModal {
	v := &vspSelectorModal{
		Load:  l,
		Modal: l.Theme.ModalFloatTitle("VSPSelectorModal"),

		inputVSP: l.Theme.Editor(new(widget.Editor), values.String(values.StrAddVSP)),
		addVSP:   l.Theme.Button(values.String(values.StrSave)),
		vspList:  l.Theme.NewClickableList(layout.Vertical),
	}
	v.inputVSP.Editor.SingleLine = true

	v.addVSP.SetEnabled(false)

	return v
}

func (v *vspSelectorModal) OnResume() {
	if len(v.WL.MultiWallet.KnownVSPs()) == 0 {
		go func() {
			v.WL.MultiWallet.ReloadVSPList(context.TODO())
			v.ParentWindow().Reload()
		}()
	}
}

func (v *vspSelectorModal) Handle() {
	v.addVSP.SetEnabled(v.editorsNotEmpty(v.inputVSP.Editor))
	if v.addVSP.Clicked() {
		go func() {
			err := v.WL.MultiWallet.SaveVSP(v.inputVSP.Editor.Text())
			if err != nil {
				v.Toast.NotifyError(err.Error())
			} else {
				v.inputVSP.Editor.SetText("")
			}
		}()
	}

	if v.Modal.BackdropClicked(true) {
		v.Dismiss()
	}

	if clicked, selectedItem := v.vspList.ItemClicked(); clicked {
		v.selectedVSP = v.WL.MultiWallet.KnownVSPs()[selectedItem]
		v.vspSelectedCallback(v.selectedVSP)
		v.Dismiss()
	}
}

func (v *vspSelectorModal) title(title string) *vspSelectorModal {
	v.dialogTitle = title
	return v
}

func (v *vspSelectorModal) vspSelected(callback func(*dcrlibwallet.VSP)) *vspSelectorModal {
	v.vspSelectedCallback = callback
	v.Dismiss()
	return v
}

func (v *vspSelectorModal) Layout(gtx layout.Context) layout.Dimensions {
	return v.Modal.Layout(gtx, []layout.Widget{
		func(gtx C) D {
			title := v.Theme.Label(values.TextSize20, v.dialogTitle)
			title.Font.Weight = text.SemiBold
			return title.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := v.Theme.Label(values.TextSize14, values.String(values.StrAddress))
					txt.Color = v.Theme.Color.GrayText2
					txtFee := v.Theme.Label(values.TextSize14, values.String(values.StrFee))
					txtFee.Color = v.Theme.Color.GrayText2
					return EndToEndRow(gtx, txt.Layout, txtFee.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					// if no vsp loaded, display a no vsp text
					vsps := v.WL.MultiWallet.KnownVSPs()
					if len(vsps) == 0 {
						noVsp := v.Theme.Label(values.TextSize14, values.String(values.StrNoVSPLoaded))
						noVsp.Color = v.Theme.Color.GrayText2
						return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, noVsp.Layout)
					}

					return v.vspList.Layout(gtx, len(vsps), func(gtx C, i int) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(0.8, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
									txt := v.Theme.Label(values.TextSize14, fmt.Sprintf("%v%%", vsps[i].FeePercentage))
									txt.Color = v.Theme.Color.GrayText1
									return EndToEndRow(gtx, v.Theme.Label(values.TextSize16, vsps[i].Host).Layout, txt.Layout)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if v.selectedVSP == nil || v.selectedVSP.Host != vsps[i].Host {
									return layout.Dimensions{}
								}
								ic := decredmaterial.NewIcon(v.Theme.Icons.NavigationCheck)
								return ic.Layout(gtx, values.MarginPadding20)
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
	})
}

func (v *vspSelectorModal) editorsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if strings.TrimSpace(e.Text()) == "" {
			return false
		}
	}

	return true
}

func (v *vspSelectorModal) OnDismiss() {}
