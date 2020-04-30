package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	// "github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageTest = "test"

// type testPage struct {
// 	gtx       *layout.Context
// 	container layout.List
// 	editorW   decredmaterial.EditorCustom
// }

func (win *Window) TestPage() {
	body := func() {
		win.testPageContents()
	}
	win.Page(body)
}

func (win *Window) testPageContents() {
	ReceivePageContent := []func(){
		func() {
			win.test1.Layout(win.gtx)
		},
		func() {
			win.test2.Layout(win.gtx)
		},
		func() {
			win.test3.Layout(win.gtx)
		},
		func() {
			win.test4.Layout(win.gtx)
		},
	}

	pageContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(win.gtx, ReceivePageContent[i])
	})
}
