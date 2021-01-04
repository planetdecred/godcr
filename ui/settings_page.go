package ui

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageSettings = "Settings"

type settingsPage struct {
	dummyText decredmaterial.Label
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	pg := &settingsPage{
		dummyText: common.theme.H5("Not yet implemented"),
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title:      "Settings",
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			back: func() {
				*common.page = PageWallet
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Left: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.dummyText.Layout(gtx)
							})
						}),
						// layout.Rigid(pg.editors(pg.addressEditor)),
						// layout.Rigid(pg.editors(pg.messageEditor)),
						// layout.Rigid(pg.drawButtonsRow()),
						// layout.Rigid(pg.drawResult()),
					)
				})
			},
			infoTemplate: "",
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, body)
}

func (pg *settingsPage) handle(common pageCommon) {

}
