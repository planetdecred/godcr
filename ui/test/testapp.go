package main

import (
	"errors"
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

var pageContainer = &layout.List{Axis: layout.Vertical}

func main() {
	gofont.Register()
	win, err := CreateWindow()
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}

	go win.Loop()

	app.Main()
}

type TestStruct struct {
	window *app.Window
	theme  *decredmaterial.Theme
	gtx    *layout.Context

	customEditorInput struct {
		test1W, test2W, test3W, test4W         widget.Editor
		test1btn, test2btn, test3btn, test4btn widget.Button
		radiobtn                               widget.Enum
		checkbox                               widget.CheckBox
	}

	customEditorOutput struct {
		test1, test2, test3, test4             decredmaterial.Editor
		test1btn, test2btn, test3btn, test4btn decredmaterial.Button
		testOutput                             decredmaterial.Label
		radiobtn                               decredmaterial.RadioButton
		checkbox                               decredmaterial.CheckBox
		progressBar                            *decredmaterial.ProgressBar
	}
}

func CreateWindow() (*TestStruct, error) {
	win := new(TestStruct)
	win.window = app.NewWindow(app.Title("GoDcr - Test app"))
	theme := decredmaterial.NewTheme()
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme
	win.gtx = layout.NewContext(win.window.Queue())

	win.initWidgets()
	return win, nil
}

func (t *TestStruct) Loop() {
	for e := range t.window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			t.gtx.Reset(e.Config, e.Size)
			t.TestPage()
			e.Frame(t.gtx.Ops)
		}
	}
}

func (t *TestStruct) initWidgets() {
	theme := t.theme

	// Editor test
	t.customEditorOutput.test1 = theme.Editor("Enter Hint Text1")
	t.customEditorOutput.test1.IsVisible = true
	t.customEditorOutput.test1.IsRequired = true

	t.customEditorOutput.test2 = theme.Editor("Enter Hint Text2")
	t.customEditorOutput.test2.IsVisible = true
	t.customEditorInput.test1W.SingleLine = true

	t.customEditorOutput.test3 = theme.Editor("Enter Hint Text3")
	t.customEditorOutput.test3.IsRequired = true

	t.customEditorOutput.test4 = theme.Editor("Enter Hint Text3")

	t.customEditorOutput.testOutput = t.theme.H6("no button clicked yet.")

	t.customEditorOutput.test1btn = theme.Button("Text1")
	t.customEditorOutput.test2btn = theme.Button("Text2")
	t.customEditorOutput.test3btn = theme.Button("Text3")
	t.customEditorOutput.test4btn = theme.Button("Text4")
	t.customEditorOutput.radiobtn = theme.RadioButton("btn1", "test radio button")
	t.customEditorOutput.checkbox = theme.CheckBox("test checkbox")
	t.customEditorOutput.progressBar = theme.ProgressBar(60)

}

func (t *TestStruct) TestPage() {
	body := func() {
		t.testPageContents()
	}
	layout.Flex{Axis: layout.Horizontal}.Layout(t.gtx,
		layout.Flexed(1, body),
	)
}

func (t *TestStruct) testPageContents() {
	t.handleInput()
	pageContent := []func(){
		func() {
			t.theme.H4("Decrematerial Test Page").Layout(t.gtx)
		},
		func() {
			t.gtx.Constraints.Width.Min = t.gtx.Px(unit.Dp(450))
			t.gtx.Constraints.Width.Max = t.gtx.Constraints.Width.Min
			t.customEditorOutput.test1.Layout(t.gtx, &t.customEditorInput.test1W)
		},
		func() {
			t.customEditorOutput.test2.Layout(t.gtx, &t.customEditorInput.test2W)
		},
		func() {
			t.customEditorOutput.test3.Layout(t.gtx, &t.customEditorInput.test3W)
		},
		func() {
			t.customEditorOutput.test4.Layout(t.gtx, &t.customEditorInput.test4W)
		},
		func() {
			layout.Flex{}.Layout(t.gtx,
				layout.Rigid(func() {
					t.customEditorOutput.test1btn.Layout(t.gtx, &t.customEditorInput.test1btn)
				}),
				layout.Rigid(func() {
					t.customEditorOutput.test2btn.Layout(t.gtx, &t.customEditorInput.test2btn)
				}),
				layout.Rigid(func() {
					t.customEditorOutput.test3btn.Layout(t.gtx, &t.customEditorInput.test3btn)
				}),
				layout.Rigid(func() {
					t.customEditorOutput.test4btn.Layout(t.gtx, &t.customEditorInput.test4btn)
				}),
			)
		},
		func() {
			t.gtx.Constraints.Height.Max = 20
			t.gtx.Constraints.Width.Max = t.gtx.Px(unit.Dp(550))
			t.customEditorOutput.progressBar.Layout(t.gtx)
		},
		func() {
			layout.Flex{}.Layout(t.gtx,
				layout.Rigid(func() {
					t.customEditorOutput.radiobtn.Layout(t.gtx, &t.customEditorInput.radiobtn)
				}),
				layout.Rigid(func() {
					t.customEditorOutput.checkbox.Layout(t.gtx, &t.customEditorInput.checkbox)
				}),
				// layout.Rigid(func() {
				// 	t.customEditorOutput.test3btn.Layout(t.gtx, &t.customEditorInput.test3btn)
				// }),
				// layout.Rigid(func() {
				// 	t.customEditorOutput.test4btn.Layout(t.gtx, &t.customEditorInput.test4btn)
				// }),
			)
		},
		func() {
			t.customEditorOutput.testOutput.Layout(t.gtx)
		},
	}

	pageContainer.Layout(t.gtx, len(pageContent), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(t.gtx, pageContent[i])
	})
}

func (t *TestStruct) handleInput() {
	if t.customEditorInput.test1btn.Clicked(t.gtx) {
		txt := t.customEditorInput.test1W.Text()
		if txt == "" {
			t.customEditorOutput.test1.ErrorLabel.Text = "This field is required and cannot be empty."
			return
		}
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorInput.test2btn.Clicked(t.gtx) {
		txt := t.customEditorInput.test2W.Text()
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorInput.test3btn.Clicked(t.gtx) {
		txt := t.customEditorInput.test3W.Text()
		if txt == "" {
			t.customEditorOutput.test3.LineColor = t.theme.Color.Danger
		}
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorInput.test4btn.Clicked(t.gtx) {
		txt := t.customEditorInput.test3W.Text()
		t.customEditorOutput.testOutput.Text = txt
	}
}
