package components

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

func mixerInfoStatusTextLayout(gtx C, l *load.Load, mixerActive bool) D {
	txt := l.Theme.H6("Mixer")
	subtxt := l.Theme.Body2("Ready to mix")
	subtxt.Color = l.Theme.Color.GrayText2
	iconVisibility := false

	if mixerActive {
		txt.Text = "Mixer is running..."
		subtxt.Text = "Keep this app opened"
		iconVisibility = true
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(txt.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if !iconVisibility {
						return layout.Dimensions{}
					}

					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, l.Theme.Icons.AlertGray.Layout16dp)
				}),
				layout.Rigid(func(gtx C) D {
					return subtxt.Layout(gtx)
				}),
			)
		}),
	)
}

func MixerInfoLayout(gtx C, l *load.Load, mixerActive bool, button layout.Widget, mixerInfo layout.Widget) D {
	return l.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := l.Theme.Icons.Mixer
								return ic.Layout24dp(gtx)
							}),
							layout.Flexed(1, func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
									return mixerInfoStatusTextLayout(gtx, l, mixerActive)
								})
							}),
							layout.Rigid(button),
						)
					})
				}),
				layout.Rigid(mixerInfo),
				layout.Rigid(func(gtx C) D {
					if mixerActive {
						txt := l.Theme.Body2("The mixer will automatically stop when unmixed balance are fully mixed.")
						txt.Color = l.Theme.Color.GrayText2
						return txt.Layout(gtx)
					}
					return D{}
				}),
			)
		})
	})
}

func MixerInfoContentWrapper(gtx C, l *load.Load, content layout.Widget) D {
	card := l.Theme.Card()
	card.Color = l.Theme.Color.Gray4
	return card.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, content)
	})
}
