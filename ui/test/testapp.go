package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gioui.org/op"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

var pageContainer *decredmaterial.ScrollContainer

func main() {
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

	customEditorOutput struct {
		test1, test2, test3, test4             decredmaterial.Editor
		test1btn, test2btn, test3btn, test4btn decredmaterial.Button
		testOutput                             decredmaterial.Label
		radiobtn                               decredmaterial.RadioButton
		checkbox                               decredmaterial.CheckBoxStyle
		progressBar                            decredmaterial.ProgressBarStyle
		outline                                decredmaterial.Outline
	}

	collapsible *decredmaterial.Collapsible
	dropDown    *decredmaterial.DropDown
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func getAbsoultePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %s", err.Error())
	}

	exSym, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("error getting filepath after evaluating sym links")
	}

	return path.Dir(exSym), nil
}

func CreateWindow() (*TestStruct, error) {
	win := new(TestStruct)
	win.window = app.NewWindow(app.Title("GoDcr - Test app"))

	absoluteWdPath, err := getAbsoultePath()
	if err != nil {
		panic(err)
	}

	decredIcons := make(map[string]image.Image)
	err = filepath.Walk(filepath.Join(absoluteWdPath, "../assets/decredicons"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := os.Open(path)
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		split := strings.Split(info.Name(), ".")
		decredIcons[split[0]] = img
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	theme := decredmaterial.NewTheme(gofont.Collection(), decredIcons)
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme

	win.initWidgets()
	return win, nil
}

func (t *TestStruct) Loop() {
	var ops op.Ops

	for e := range t.window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			t.TestPage(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (t *TestStruct) initWidgets() {
	theme := t.theme

	pageContainer = theme.ScrollContainer()

	// Editor test
	t.customEditorOutput.test1 = theme.Editor(new(widget.Editor), "Enter Hint Text1")
	t.customEditorOutput.test1.IsVisible = true
	t.customEditorOutput.test1.IsRequired = true

	t.customEditorOutput.test2 = theme.Editor(new(widget.Editor), "Enter Hint Text2")
	t.customEditorOutput.test2.IsVisible = true
	t.customEditorOutput.test1.Editor.SingleLine = true

	t.customEditorOutput.test3 = theme.Editor(new(widget.Editor), "Enter Hint Text3")
	t.customEditorOutput.test3.IsRequired = true

	t.customEditorOutput.test4 = theme.Editor(new(widget.Editor), "Enter Hint Text3")

	t.customEditorOutput.testOutput = t.theme.H6("no button clicked yet.")

	t.customEditorOutput.test1btn = theme.Button(new(widget.Clickable), "Text1")
	t.customEditorOutput.test2btn = theme.Button(new(widget.Clickable), "Text2")
	t.customEditorOutput.test3btn = theme.Button(new(widget.Clickable), "Text3")
	t.customEditorOutput.test4btn = theme.Button(new(widget.Clickable), "Text4")
	t.customEditorOutput.radiobtn = theme.RadioButton(new(widget.Enum), "btn1", "test radio button")
	t.customEditorOutput.checkbox = theme.CheckBox(new(widget.Bool), "test checkbox")
	t.customEditorOutput.progressBar = theme.ProgressBar(60)
	t.customEditorOutput.outline = theme.Outline()
	// t.customEditorOutput.outline.Color = theme.Color.Primary

	t.collapsible = theme.Collapsible()

	dropDownItems := []decredmaterial.DropDownItem{
		{
			Text: "All",
		},
		{
			Text: "Not All",
		},
		{
			Text: "Semi All",
		},
	}
	t.dropDown = theme.DropDown(dropDownItems, 1)
}

func (t *TestStruct) TestPage(gtx layout.Context) {
	body := func(gtx C) D {
		return t.testPageContents(gtx)
	}
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, body),
	)
}

func (t *TestStruct) testPageContents(gtx layout.Context) layout.Dimensions {
	t.handleInput()
	header := func(gtx layout.Context) layout.Dimensions {
		return t.theme.Body1("Collapsible Widget").Layout(gtx)
	}
	content := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return t.theme.Body2("Hidden item 1").Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return t.theme.Body2("Hidden item 2").Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return t.theme.Body2("Hidden item 3").Layout(gtx)
			}),
		)
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return t.theme.H4("Decrematerial Test Page").Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(unit.Dp(450))
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
			return t.customEditorOutput.test1.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test2.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test3.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test4.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test1btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test2btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test3btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test4btn.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = 20
			gtx.Constraints.Max.X = gtx.Px(unit.Dp(550))
			return t.customEditorOutput.progressBar.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.radiobtn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.checkbox.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return t.customEditorOutput.testOutput.Layout(gtx)
		},

		func(gtx C) D {
			header := func(gtx layout.Context) layout.Dimensions {
				return t.theme.Body1("Collapsible Widget").Layout(gtx)
			}
			content := func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return t.theme.Body2("Hidden item 1").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return t.theme.Body2("Hidden item 2").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return t.theme.Body2("Hidden item 3").Layout(gtx)
					}),
				)
			}
			return t.collapsible.Layout(gtx, header, content, nil)
		},

		func(gtx C) D {
			return t.customEditorOutput.outline.Layout(gtx, func(gtx C) D {
				return t.customEditorOutput.testOutput.Layout(gtx)
			})
		},
		func(gtx C) D {
			return t.theme.H4("Decrematerial Test Page").Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(unit.Dp(450))
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
			return t.customEditorOutput.test1.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test2.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test3.Layout(gtx)
		},
		func(gtx C) D {
			return t.customEditorOutput.test4.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test1btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test2btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test3btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.test4btn.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = 20
			gtx.Constraints.Max.X = gtx.Px(unit.Dp(550))
			return t.customEditorOutput.progressBar.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.radiobtn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return t.customEditorOutput.checkbox.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return t.customEditorOutput.testOutput.Layout(gtx)
		},

		func(gtx C) D {
			return t.collapsible.Layout(gtx, header, content, nil)
		},
		func(gtx C) D {
			return t.customEditorOutput.outline.Layout(gtx, func(gtx C) D {
				return t.customEditorOutput.testOutput.Layout(gtx)
			})
		},
	}

	pageContent = append(pageContent, pageContent...)
	return pageContainer.Layout(gtx, pageContent)
}

func (t *TestStruct) handleInput() {
	if t.customEditorOutput.test1btn.Button.Clicked() {
		txt := t.customEditorOutput.test1.Editor.Text()
		if txt == "" {
			t.customEditorOutput.test1.SetError("This field is required and cannot be empty.")
			return
		}
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorOutput.test2btn.Button.Clicked() {
		txt := t.customEditorOutput.test2.Editor.Text()
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorOutput.test3btn.Button.Clicked() {
		txt := t.customEditorOutput.test3.Editor.Text()
		if txt == "" {
			t.customEditorOutput.test3.LineColor = t.theme.Color.Danger
		}
		t.customEditorOutput.testOutput.Text = txt
	}
	if t.customEditorOutput.test4btn.Button.Clicked() {
		txt := t.customEditorOutput.test3.Editor.Text()
		t.customEditorOutput.testOutput.Text = txt
	}
}
