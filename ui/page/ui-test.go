package page

import (
	"fmt"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/values"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// UITestID is the id of the ui test page
const UITestID = "ui-test"

type selectWidget struct {
	label         material.Label
	widget        *materialplus.Select
	selectedKey   material.Label
	selectedValue material.Label
}

// UITest represents the ui test page of the app
// It is solely to showcase custom widgets used by the app
type UITest struct {
	selectWidget *selectWidget

	loadMainUIButton         *widget.Button
	loadMainUIButtonMaterial material.Button
	progressBar              *materialplus.ProgressBar
	states                   map[string]interface{}
	card  					 materialplus.Card
}

// Init initializes all available custom widgets
func (pg *UITest) Init(theme *materialplus.Theme, _ *wallet.Wallet, states map[string]interface{}) {
	selectOptions := []materialplus.SelectItem{
		{
			Key:  10,
			Text: "Select item 1",
		},
		{
			Key:  20,
			Text: "Select item 2",
		},
		{
			Key:  30,
			Text: "Select item 3",
		},
	}

	pg.selectWidget = &selectWidget{
		widget:        theme.Select(selectOptions),
		label:         theme.Body1("Select Widget"),
		selectedKey:   theme.Body2(""),
		selectedValue: theme.Body2(""),
	}

	pg.loadMainUIButton = new(widget.Button)
	pg.loadMainUIButtonMaterial = theme.Button("Load Main UI")
	pg.progressBar = theme.ProgressBar()
	pg.card = theme.Card()

	pg.states = states
}

// Draw renders the widgets to screen
func (pg *UITest) Draw(gtx *layout.Context) (res interface{}) {
	widgets := []func(){
		func() {
			pg.selectWidget.label.Layout(gtx)
			inset := layout.Inset{
				Top: unit.Dp(25),
			}
			inset.Layout(gtx, func() {
				pg.selectWidget.widget.Layout(gtx)
			})
		},
		func() {
			selected := pg.selectWidget.widget.GetSelected()

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.selectWidget.selectedKey.Text = fmt.Sprintf("Selected Key: %d", selected.Key)
					pg.selectWidget.selectedKey.Layout(gtx)
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(15),
					}
					inset.Layout(gtx, func() {
						pg.selectWidget.selectedValue.Text = fmt.Sprintf("Selected Value: %s", selected.Text)
						pg.selectWidget.selectedValue.Layout(gtx)
					})
				}),
			)
		},
	}

	pageWidgets := []func(){
		func() {
			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(widgets), func(i int) {
				layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
			})
		},
		func() {
			pg.progressBar.Layout(gtx, 25)
		},
		func() {
			for pg.loadMainUIButton.Clicked(gtx) {
				walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
				ev := EventNav{
					Current: UITestID,
				}

				if walletInfo.LoadedWallets > 0 {
					ev.Next = WalletsID
				} else {
					ev.Next = LandingID
				}

				res = ev
			}
			pg.loadMainUIButtonMaterial.Layout(gtx, pg.loadMainUIButton)
		},
		func() {
			pg.card.Width = gtx.Px(units.WalletSyncBoxWidthMin)
			pg.card.Height = gtx.Px(units.WalletSyncBoxHeightMin)
			pg.card.Color = values.ProgressBarGreen
			pg.card.Layout(gtx, 0)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(pageWidgets), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, pageWidgets[i])
	})

	return
}
