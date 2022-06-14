package page

import (
	"image"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const HelpPageID = "Help"

type HelpPage struct {
	*load.Load

	documentation   *decredmaterial.Clickable
	copyRedirectURL *decredmaterial.Clickable
	shadowBox       *decredmaterial.Shadow
	backButton      decredmaterial.IconButton
}

func NewHelpPage(l *load.Load) *HelpPage {
	pg := &HelpPage{
		Load:            l,
		documentation:   l.Theme.NewClickable(true),
		copyRedirectURL: l.Theme.NewClickable(false),
	}

	pg.shadowBox = l.Theme.Shadow()
	pg.shadowBox.SetShadowRadius(14)

	pg.documentation.Radius = decredmaterial.Radius(14)
	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *HelpPage) ID() string {
	return HelpPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *HelpPage) OnNavigatedTo() {

}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *HelpPage) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrHelp),
			SubTitle:   values.String(values.StrHelpInfo),
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween, WeightSum: 2}.Layout(gtx,
						layout.Flexed(1, pg.document()),
					)
				})
			},
		}
		return sp.Layout(gtx)
	}
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx, body)
	}
	return pg.layoutDesktop(gtx, body)
}

func (pg *HelpPage) layoutDesktop(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformPadding(gtx, body)
}

func (pg *HelpPage) layoutMobile(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformMobile(gtx, false, body)
}

func (pg *HelpPage) document() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Theme.Icons.DocumentationIcon, pg.documentation, values.String(values.StrDocumentation))
	}
}

func (pg *HelpPage) pageSections(gtx C, icon *decredmaterial.Image, action *decredmaterial.Clickable, title string) D {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Orientation: layout.Vertical,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   action,
			Direction:   layout.Center,
			Alignment:   layout.Middle,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return icon.Layout24dp(gtx)
			}),
			layout.Rigid(pg.Theme.Body1(title).Layout),
			layout.Rigid(func(gtx C) D {
				size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y}
				return D{Size: size}
			}),
		)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *HelpPage) HandleUserInteractions() {
	if pg.documentation.Clicked() {
		decredURL := "https://docs.decred.org"
		info := modal.NewInfoModal(pg.Load).
			Title("View documentation").
			Body(values.String(values.StrCopyLink)).
			SetCancelable(true).
			UseCustomWidget(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						border := widget.Border{Color: pg.Theme.Color.Gray4, CornerRadius: values.MarginPadding10, Width: values.MarginPadding2}
						wrapper := pg.Theme.Card()
						wrapper.Color = pg.Theme.Color.Gray4
						return border.Layout(gtx, func(gtx C) D {
							return wrapper.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Flexed(0.9, pg.Theme.Body1(decredURL).Layout),
										layout.Flexed(0.1, func(gtx C) D {
											return layout.E.Layout(gtx, func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
													if pg.copyRedirectURL.Clicked() {
														clipboard.WriteOp{Text: decredURL}.Add(gtx.Ops)
														pg.Toast.Notify(values.String(values.StrCopied))
													}
													return pg.copyRedirectURL.Layout(gtx, pg.Theme.Icons.CopyIcon.Layout24dp)
												})
											})
										}),
									)
								})
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						return layout.Inset{
							Top:  values.MarginPaddingMinus10,
							Left: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							label := pg.Theme.Body2(values.String(values.StrWebURL))
							label.Color = pg.Theme.Color.GrayText2
							return label.Layout(gtx)
						})
					}),
				)
			}).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) {})
		pg.ShowModal(info)
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *HelpPage) OnNavigatedFrom() {}
