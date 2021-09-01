package components

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/values"
)

type SubPage struct {
	*load.Load
	Title        string
	SubTitle     string
	WalletName   string
	Back         func()
	Body         layout.Widget
	InfoTemplate string
	ExtraItem    *widget.Clickable
	Extra        layout.Widget
	ExtraText    string
	HandleExtra  func()

	BackButton decredmaterial.IconButton
	InfoButton decredmaterial.IconButton
}

func SubpageHeaderButtons(l *load.Load) (decredmaterial.IconButton, decredmaterial.IconButton) {
	backButton := l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.NavigationArrowBack)
	infoButton := l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.ActionInfo)

	zeroInset := layout.UniformInset(values.MarginPadding4)
	backButton.Color, infoButton.Color = l.Theme.Color.Gray3, l.Theme.Color.Gray3

	m24 := values.MarginPadding24
	backButton.Size, infoButton.Size = m24, m24
	backButton.Inset, infoButton.Inset = zeroInset, zeroInset

	return backButton, infoButton
}

func (sp *SubPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return sp.Header(gtx)
			})
		}),
		layout.Rigid(sp.Body),
	)
}

func (sp *SubPage) Header(gtx layout.Context) layout.Dimensions {
	sp.EventHandler()

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, sp.BackButton.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := sp.Load.Theme.Label(values.TextSize20, sp.Title)
			title.Color = sp.Load.Theme.Color.DeepBlue

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(title.Layout),
				layout.Rigid(func(gtx C) D {
					if sp.SubTitle == "" {
						return D{}
					}

					sub := sp.Load.Theme.Label(values.TextSize14, sp.SubTitle)
					sub.Color = sp.Load.Theme.Color.Gray
					return sub.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.WalletName != "" {
				return layout.Inset{Left: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return decredmaterial.Card{
						Color: sp.Theme.Color.Surface,
					}.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
							walletText := sp.Theme.Caption(sp.WalletName)
							walletText.Color = sp.Theme.Color.Gray
							return walletText.Layout(gtx)
						})
					})
				})
			}
			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.E.Layout(gtx, func(gtx C) D {
				if sp.InfoTemplate != "" {
					return sp.InfoButton.Layout(gtx)
				} else if sp.ExtraItem != nil {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if sp.ExtraText != "" {
								return layout.Inset{Right: values.MarginPadding10, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									text := sp.Theme.Caption(sp.ExtraText)
									text.Color = sp.Theme.Color.DeepBlue
									return text.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return decredmaterial.Clickable(gtx, sp.ExtraItem, sp.Extra)
						}),
					)
				}
				return layout.Dimensions{}
			})
		}),
	)
}

func (sp *SubPage) SplitLayout(gtx layout.Context) layout.Dimensions {
	card := sp.Theme.Card()
	card.Color = color.NRGBA{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D { return sp.Header(gtx) }),
			layout.Rigid(sp.Body),
		)
	})
}

func (sp *SubPage) EventHandler() {
	if sp.InfoTemplate != "" {
		if sp.InfoButton.Button.Clicked() {
			modal.NewInfoModal(sp.Load).
				Title(sp.Title).
				SetupWithTemplate(sp.InfoTemplate).
				NegativeButton("Got it", func() {}).Show()
		}
	}

	if sp.BackButton.Button.Clicked() {
		sp.Back()
	}

	if sp.ExtraItem != nil && sp.ExtraItem.Clicked() {
		sp.HandleExtra()
	}
}
