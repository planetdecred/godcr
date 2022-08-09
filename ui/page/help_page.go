package page

import (
	"image"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const HelpPageID = "Help"

type HelpPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	documentation   *decredmaterial.Clickable
	copyRedirectURL *decredmaterial.Clickable
	shadowBox       *decredmaterial.Shadow
	backButton      decredmaterial.IconButton
}

func NewHelpPage(l *load.Load) *HelpPage {
	pg := &HelpPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(HelpPageID),
		documentation:    l.Theme.NewClickable(true),
		copyRedirectURL:  l.Theme.NewClickable(false),
	}

	pg.shadowBox = l.Theme.Shadow()
	pg.shadowBox.SetShadowRadius(14)

	pg.documentation.Radius = decredmaterial.Radius(14)
	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
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
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *HelpPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding20}.Layout(gtx, pg.pageContentLayout)
			}),
		)
	})
}

func (pg *HelpPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrHelp),
			SubTitle:   values.String(values.StrHelpInfo),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween, WeightSum: 1}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return pg.pageSectionsMobile(gtx, pg.Theme.Icons.DocumentationIcon, pg.documentation, values.String(values.StrDocumentation))
						}),
					)
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformMobile(gtx, false, false, body)
}

func (pg *HelpPage) pageHeaderLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding16,
							Top:   values.MarginPaddingMinus2,
						}.Layout(gtx, pg.backButton.Layout)
					}),
					layout.Rigid(pg.Theme.Label(values.TextSize20, values.String(values.StrHelp)).Layout),
				)
			})
		}),
	)
}

func (pg *HelpPage) pageContentLayout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(values.MarginPadding550)
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return layout.Flex{WeightSum: 3, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					sub := pg.Load.Theme.Label(values.TextSize14, values.String(values.StrHelpInfo))
					sub.Color = pg.Load.Theme.Color.GrayText2
					return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, sub.Layout)
				}),
				layout.Flexed(1, pg.document()),
			)
		})
	})
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

func (pg *HelpPage) pageSectionsMobile(gtx C, icon *decredmaterial.Image, action *decredmaterial.Clickable, title string) D {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Orientation: layout.Horizontal,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   action,
			Direction:   layout.W,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return icon.Layout24dp(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:  values.MarginPadding2,
					Left: values.MarginPadding18,
				}.Layout(gtx, func(gtx C) D {
					return pg.Theme.Body1(title).Layout(gtx)
				})
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
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}

	if pg.backButton.Button.Clicked() {
		pg.ParentNavigator().CloseCurrentPage()
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
