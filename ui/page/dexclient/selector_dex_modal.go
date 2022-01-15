package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const selectorDexModalID = "selector_dex_modal"

type selectorDexModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	selectDexHostBtn map[string]*decredmaterial.Clickable
	selectedHost     string
	cancelBtn        decredmaterial.Button
	callback         func(*core.Exchange)
}

func newSelectorDexModal(l *load.Load, selectedHost string) *selectorDexModal {
	md := &selectorDexModal{
		Load:         l,
		selectedHost: selectedHost,
		modal:        l.Theme.ModalFloatTitle(),
		cancelBtn:    l.Theme.OutlineButton("Cancel"),
	}

	return md
}

func (md *selectorDexModal) ModalID() string {
	return selectorDexModalID
}

func (md *selectorDexModal) Show() {
	md.ShowModal(md)
}

func (md *selectorDexModal) Dismiss() {
	md.DismissModal(md)
}

func (md *selectorDexModal) OnDismiss() {
}

func (md *selectorDexModal) OnResume() {
	md.initDEXServersWidget()
}

func (md *selectorDexModal) OnDexSelected(callback func(dex *core.Exchange)) *selectorDexModal {
	md.callback = callback
	return md
}

func (md *selectorDexModal) Handle() {
	if md.cancelBtn.Button.Clicked() {
		md.Dismiss()
	}

	for host, selectBtn := range md.selectDexHostBtn {
		if selectBtn.Clicked() {
			md.callback(md.Dexc().DEXServers()[host])
			md.Dismiss()
			return
		}
	}
}

func (md *selectorDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Select Dex").Layout(gtx)
		},
		func(gtx C) D {
			var childrens = make([]layout.FlexChild, 0, len(md.selectDexHostBtn))
			exchanges := sliceExchanges(md.Dexc().DEXServers())
			for i := 0; i < len(exchanges); i++ {
				host := exchanges[i].Host
				childrens = append(childrens, layout.Rigid(func(gtx C) D {
					return md.selectDexHostBtn[host].Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{
							Top: values.MarginPadding4, Bottom: values.MarginPadding4,
							Left: values.MarginPadding8, Right: values.MarginPadding8,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(md.Theme.Label(values.TextSize14, host).Layout),
								layout.Rigid(func(gtx C) D {
									if md.selectedHost != host {
										return D{}
									}
									gtx.Constraints.Min.X = 30
									ic := md.Icons.NavigationCheck
									return ic.Layout(gtx, md.Theme.Color.Success)
								}),
							)
						})
					})
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, md.cancelBtn.Layout)
		},
	}

	return md.modal.Layout(gtx, w)
}

func (md *selectorDexModal) initDEXServersWidget() {
	exchanges := sliceExchanges(md.Dexc().DEXServers())
	md.selectDexHostBtn = make(map[string]*decredmaterial.Clickable, len(exchanges))
	for i := 0; i < len(exchanges); i++ {
		cl := md.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		md.selectDexHostBtn[exchanges[i].Host] = cl
	}
}
