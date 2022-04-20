package page

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const AboutPageID = "About"

type AboutPage struct {
	*load.Load
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
		card:             l.Theme.Card(),
		container:        &layout.List{Axis: layout.Vertical},
		version:          l.Theme.Body1("Version"),
		versionValue:     l.Theme.Body1(l.WL.Wallet.Version()),
		buildDate:        l.Theme.Body1("Build date"),
		buildDateValue:   l.Theme.Body1(l.WL.Wallet.BuildDate().Format("2006-01-02 15:04:05")),
		network:          l.Theme.Body1("Network"),
		license:          l.Theme.Body1("License"),
		licenseRow:       l.Theme.NewClickable(true),
		shadowBox:        l.Theme.Shadow(),
		chevronRightIcon: decredmaterial.NewIcon(l.Theme.Icons.ChevronRight),
	}

	pg.licenseRow.Radius = decredmaterial.CornerRadius{
		BottomRight: values.MarginPadding14.V,
		BottomLeft:  values.MarginPadding14.V,
	}

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

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *AboutPage) ID() string {
	return AboutPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AboutPage) OnNavigatedTo() {

}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *AboutPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "About",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.card.Layout(gtx, func(gtx C) D {
					return pg.layoutRows(gtx)
				})
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *AboutPage) layoutRows(gtx layout.Context) layout.Dimensions {
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
						return layout.Dimensions{}
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
		pg.ChangeFragment(NewLicensePage(pg.Load))
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
