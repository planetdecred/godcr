package page

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const AboutPageID = "About"

type AboutPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	card      decredmaterial.Card
	container *layout.List

	version        decredmaterial.Label
	versionValue   decredmaterial.Label
	buildDate      decredmaterial.Label
	buildDateValue decredmaterial.Label
	network        decredmaterial.Label
	networkValue   decredmaterial.Label
	license        decredmaterial.Label
	licenseRow     *decredmaterial.Clickable

	chevronRightIcon *decredmaterial.Icon

	backButton decredmaterial.IconButton
	shadowBox  *decredmaterial.Shadow
}

func NewAboutPage(l *load.Load) *AboutPage {
	pg := &AboutPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(AboutPageID),
		card:             l.Theme.Card(),
		container:        &layout.List{Axis: layout.Vertical},
		version:          l.Theme.Body1(values.String(values.StrVersion)),
		versionValue:     l.Theme.Body1(l.WL.Wallet.Version()),
		buildDate:        l.Theme.Body1(values.String(values.StrBuildDate)),
		buildDateValue:   l.Theme.Body1(l.WL.Wallet.BuildDate().Format("2006-01-02 15:04:05")),
		network:          l.Theme.Body1(values.String(values.StrNetwork)),
		license:          l.Theme.Body1(values.String(values.StrLicense)),
		licenseRow:       l.Theme.NewClickable(true),
		shadowBox:        l.Theme.Shadow(),
		chevronRightIcon: decredmaterial.NewIcon(l.Theme.Icons.ChevronRight),
	}

	pg.licenseRow.Radius = decredmaterial.BottomRadius(14)

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	col := pg.Theme.Color.GrayText2
	pg.versionValue.Color = col
	pg.buildDateValue.Color = col

	netType := pg.WL.Wallet.Net
	if pg.WL.Wallet.Net == dcrlibwallet.Testnet3 {
		netType = "Testnet"
	}
	pg.networkValue = l.Theme.Body1(netType)
	pg.networkValue.Color = col

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AboutPage) OnNavigatedTo() {

}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *AboutPage) Layout(gtx C) D {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *AboutPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding20}.Layout(gtx, pg.pageContentLayout)
			}),
		)
	})
}

func (pg *AboutPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (pg *AboutPage) pageHeaderLayout(gtx layout.Context) layout.Dimensions {
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
					layout.Rigid(pg.Theme.Label(values.TextSize20, values.String(values.StrAbout)).Layout),
				)
			})
		}),
	)
}

func (pg *AboutPage) pageContentLayout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(values.MarginPadding550)
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		return pg.card.Layout(gtx, pg.layoutRows)
	})
}

func (pg *AboutPage) layoutRows(gtx C) D {
	var in = layout.Inset{
		Top:    values.MarginPadding20,
		Bottom: values.MarginPadding20,
		Left:   values.MarginPadding16,
		Right:  values.MarginPadding16,
	}
	w := []func(gtx C) D{
		func(gtx C) D {
			return components.Container{Padding: in}.Layout(gtx, func(gtx C) D {
				return components.EndToEndRow(gtx, pg.version.Layout, pg.versionValue.Layout)
			})
		},
		func(gtx C) D {
			return components.Container{Padding: in}.Layout(gtx, func(gtx C) D {
				return components.EndToEndRow(gtx, pg.buildDate.Layout, pg.buildDateValue.Layout)
			})
		},
		func(gtx C) D {
			return components.Container{Padding: in}.Layout(gtx, func(gtx C) D {
				return components.EndToEndRow(gtx, pg.network.Layout, pg.networkValue.Layout)
			})
		},

		func(gtx C) D {
			licenseRowLayout := func(gtx C) D {
				return pg.licenseRow.Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, pg.license.Layout)
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return in.Layout(gtx, func(gtx C) D {
									pg.chevronRightIcon.Color = pg.Theme.Color.Gray1
									return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
								})
							})
						}),
					)
				})
			}
			if pg.licenseRow.IsHovered() {
				return pg.shadowBox.Layout(gtx, licenseRowLayout)
			}
			return licenseRowLayout(gtx)
		},
	}

	return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return layout.Inset{Bottom: values.MarginPadding3}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(w[i]),
				layout.Rigid(func(gtx C) D {
					if i == len(w)-1 {
						return D{}
					}
					return layout.Inset{
						Left: values.MarginPadding16,
					}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
			)
		})
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *AboutPage) HandleUserInteractions() {
	if pg.licenseRow.Clicked() {
		pg.ParentNavigator().Display(NewLicensePage(pg.Load))
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
func (pg *AboutPage) OnNavigatedFrom() {}
