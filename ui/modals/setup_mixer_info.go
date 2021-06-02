package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/values"
)

type SetupMixerInfo struct {
	title string
	*common
}

const SetupMixerInfoModal = " SetupMixerInfo"

func (m *Modals) registerSetupMixerInfoModal() {
	m.modals[PrivacyInfoModal] = &PrivacyInfo{
		title:  "Mixer Info",
		common: m.common,
	}
}

func (m *SetupMixerInfo) getTitle() string {
	return m.title
}

func (m *SetupMixerInfo) onCancel()  {}
func (m *SetupMixerInfo) onConfirm() {}

func (m *SetupMixerInfo) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			txt := m.theme.Body1("Two dedicated accounts (“mixed” & “unmixed”) will be created in order to use the mixer.")
			txt.Color = m.theme.Color.Gray
			return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, txt.Layout)
		},
		func(gtx C) D {
			txt := m.theme.Label(values.TextSize18, "This action cannot be undone.")
			return txt.Layout(gtx)
		},
	}
}
