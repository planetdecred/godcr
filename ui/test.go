package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

const PageTest = "test"

func (win *Window) TestPage() {
	body := func() {
		win.testPageContents()
	}
	win.Page(body)
}

func (win *Window) testPageContents() {
	win.handleInput()
	ReceivePageContent := []func(){
		func() {
			win.outputs.customEditor.test1.Layout(win.gtx)
		},
		func() {
			win.outputs.customEditor.test2.Layout(win.gtx)
		},
		func() {
			win.outputs.customEditor.test3.Layout(win.gtx)
		},
		func() {
			win.outputs.customEditor.test4.Layout(win.gtx)
		},
		func() {
			layout.Flex{}.Layout(win.gtx,
				layout.Rigid(func() {
					win.outputs.customEditor.test1btn.Layout(win.gtx, &win.inputs.customEditor.test1btn)
				}),
				layout.Rigid(func() {
					win.outputs.customEditor.test2btn.Layout(win.gtx, &win.inputs.customEditor.test2btn)
				}),
				layout.Rigid(func() {
					win.outputs.customEditor.test3btn.Layout(win.gtx, &win.inputs.customEditor.test3btn)
				}),
				layout.Rigid(func() {
					win.outputs.customEditor.test4btn.Layout(win.gtx, &win.inputs.customEditor.test4btn)
				}),
			)
		},
		func() {
			win.outputs.customEditor.testOutput.Layout(win.gtx)
		},
	}

	pageContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(win.gtx, ReceivePageContent[i])
	})
}

func (win *Window) handleInput() {
	if win.inputs.customEditor.test1btn.Clicked(win.gtx) {
		t := win.outputs.customEditor.test1.Text()
		if t == "" {
			win.outputs.customEditor.test1.ErrorLabel.Text = "this field is required n cannot be empty."
		}
		win.outputs.customEditor.testOutput.Text = t
	}
	if win.inputs.customEditor.test2btn.Clicked(win.gtx) {
		t := win.outputs.customEditor.test2.Text()
		win.outputs.customEditor.testOutput.Text = t
	}
	if win.inputs.customEditor.test3btn.Clicked(win.gtx) {
		t := win.outputs.customEditor.test3.Text()
		if t == "" {
			win.outputs.customEditor.test3.ErrorLabel.Text = "this field is required n cannot be empty."
		}
		win.outputs.customEditor.testOutput.Text = t
	}
	if win.inputs.customEditor.test4btn.Clicked(win.gtx) {
		t := win.outputs.customEditor.test4.Text()
		win.outputs.customEditor.testOutput.Text = t
	}
}
