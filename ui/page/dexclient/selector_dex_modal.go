package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const selectorDexModalID = "selector_dex_modal"

type selectorDexModal struct {
	*load.Load
	modal                *decredmaterial.Modal
	list                 *widget.List
	cancelBtn            decredmaterial.Button
	selectorExchangeWdgs []*selectorExchangeWidget
	selectedHost         string
	callback             func(*core.Exchange)
}

type selectorExchangeWidget struct {
	selectBtn *decredmaterial.Clickable
	*core.Exchange
}

func newSelectorDexModal(l *load.Load, selectedHost string) *selectorDexModal {
	md := &selectorDexModal{
		Load:         l,
		selectedHost: selectedHost,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		modal:     l.Theme.ModalFloatTitle(),
		cancelBtn: l.Theme.OutlineButton("Cancel"),
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

	for _, selectorExchangeWdg := range md.selectorExchangeWdgs {
		if selectorExchangeWdg.selectBtn.Clicked() {
			md.callback(selectorExchangeWdg.Exchange)
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
			return md.Theme.List(md.list).Layout(gtx, len(md.selectorExchangeWdgs), func(gtx C, i int) D {
				selectorExchangeWdg := md.selectorExchangeWdgs[i]
				return selectorExchangeWdg.selectBtn.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Inset{
						Top: values.MarginPadding4, Bottom: values.MarginPadding4,
						Left: values.MarginPadding8, Right: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(md.Theme.Label(values.TextSize14, selectorExchangeWdg.Host).Layout),
							layout.Rigid(func(gtx C) D {
								if md.selectedHost != selectorExchangeWdg.Host {
									return D{}
								}
								gtx.Constraints.Min.X = 30
								ic := md.Load.Icons.NavigationCheck
								return ic.Layout(gtx, md.Theme.Color.Success)
							}),
						)
					})
				})
			})
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, md.cancelBtn.Layout)
		},
	}

	return md.modal.Layout(gtx, w)
}

func (md *selectorDexModal) initDEXServersWidget() {
	exchanges := sliceExchanges(md.Dexc().DEXServers())
	var selectorExchangeWdgs []*selectorExchangeWidget
	for i := 0; i < len(exchanges); i++ {
		exchange := exchanges[i]
		cl := md.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		selectorExchangeWdgs = append(selectorExchangeWdgs, &selectorExchangeWidget{
			selectBtn: cl,
			Exchange:  exchange,
		})
	}
	md.selectorExchangeWdgs = selectorExchangeWdgs
}
