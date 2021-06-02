package modals

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/text"

	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PrivacyInfo struct {
	title string
	*common
}

const PrivacyInfoModal = "PrivacyInfo"

func (m *Modals) registerPrivacyInfoModal() {
	m.modals[PrivacyInfoModal] = &PrivacyInfo{
		title:  "Privacy Info",
		common: m.common,
	}
}

func (m *PrivacyInfo) getTitle() string {
	return m.title
}

func (m *PrivacyInfo) onCancel()  {}
func (m *PrivacyInfo) onConfirm() {}

func (m *PrivacyInfo) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.theme.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := m.theme.Body1("When you turn on the mixer, your unmixed DCRs in this wallet (unmixed balance) will be gradually mixed.")
					text.Color = m.theme.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			txt := m.theme.Body1("Important: keep this app opened while mixer is running.")
			txt.Font.Weight = text.Bold
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.theme.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := m.theme.Body1("Mixer will automatically stop when unmixed balance are fully mixed.")
					text.Color = m.theme.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
	}
}
