package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

//DemoID is the id of the ui test page
const DemoID = "ui-test"

type selectWidget struct {
	label         material.Label
	widget        *materialplus.Select
	selectedKey   material.Label
	selectedValue material.Label
}

//Demo represents the ui test page of the app
//It is solely to showcase custom widgets used by the app
type Demo struct {
	Common
	selectWidget *selectWidget

	loadMainUIButton         *widget.Button
	loadMainUIButtonMaterial material.Button
	progressBar              *materialplus.ProgressBar
	states                   map[string]interface{}
}

func (demo *Demo) Init(*materialplus.Theme, *wallet.Wallet, map[string]interface{}) {

}

func (demo *Demo) Draw(gtx *layout.Context) interface{} {

	return nil
}
