package components

/*
func MixerInfoStatusTextLayout(gtx C, l *load.Load, mixerActive bool) D {
	txt := l.Theme.H6(values.String(values.StrMixer))
	subtxt := l.Theme.Body2(values.String(values.StrReadyToMix))
	subtxt.Color = l.Theme.Color.GrayText2
	iconVisibility := false

	if mixerActive {
		txt.Text = values.String(values.StrMixerRunning)
		subtxt.Text = values.String(values.StrKeepAppOpen)
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

func MixerInfoContentWrapper(gtx C, l *load.Load, content layout.Widget) D {
	card := l.Theme.Card()
	card.Color = l.Theme.Color.Gray4
	return card.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, content)
	})
}
*/
