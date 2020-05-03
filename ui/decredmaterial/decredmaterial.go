// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	// "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	// "golang.org/x/exp/shiny/materialdesign/icons"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(650)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := NewTheme()
	gtx := layout.NewContext(w.Queue())

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)
				kitchen(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

var (
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	button            = new(widget.Button)
	greenButton       = new(widget.Button)
	iconTextButton    = new(widget.Button)
	iconButton        = new(widget.Button)
	radioButtonsGroup = new(widget.Enum)
	list              = &layout.List{
		Axis: layout.Vertical,
	}
	progress            = 0
	gree               = true
	topLabel            = "Godcr, Decredmaterial"
	icon                *material.Icon
	checkbox            = new(widget.CheckBox)
)

func kitchen(gtx *layout.Context, th *material.Theme) {
	widgets := []func(){
		func() {
			th.EditorCustom(topLabel).Layout(gtx)
		},
		func() {
			th.H3(topLabel).Layout(gtx)
		},
		func() {
			th.H3(topLabel).Layout(gtx)
		},
		func() {
			th.H3(topLabel).Layout(gtx)
		},
		// func() {
		// 	gtx.Constraints.Height.Max = gtx.Px(unit.Dp(200))
		// 	th.Editor("Hint").Layout(gtx, editor)
		// },
		// func() {
		// 	e := th.Editor("Hint")
		// 	e.Font.Style = text.Italic
		// 	e.Layout(gtx, lineEditor)
		// 	for _, e := range lineEditor.Events(gtx) {
		// 		if e, ok := e.(widget.SubmitEvent); ok {
		// 			topLabel = e.Text
		// 			lineEditor.SetText("")
		// 		}
		// 	}
		// },
		func() {
			in := layout.UniformInset(unit.Dp(8))
			layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func() {
					in.Layout(gtx, func() {
						th.IconButton(icon).Layout(gtx, iconButton)
					})
				}),
				layout.Rigid(func() {
					in.Layout(gtx, func() {
						// iconAndTextButton{th}.Layout(gtx, iconTextButton, icon, "Horizontal button")
					})
				}),
				layout.Rigid(func() {
					in.Layout(gtx, func() {
						for button.Clicked(gtx) {
							gree = !gree
						}
						th.Button("Click me!").Layout(gtx, button)
					})
				}),
				layout.Rigid(func() {
					in.Layout(gtx, func() {
						var btn material.Button
						btn = th.Button("Green button")
						if gree {
							btn.Background = color.RGBA{A: 0xff, R: 0x9e, G: 0x9d, B: 0x24}
						}
						btn.Layout(gtx, greenButton)
					})
				}),
			)
		},
		func() {
			th.CheckBox("Checkbox").Layout(gtx, checkbox)
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					th.RadioButton("r1", "RadioButton1").Layout(gtx, radioButtonsGroup)
				}),
				layout.Rigid(func() {
					th.RadioButton("r2", "RadioButton2").Layout(gtx, radioButtonsGroup)
				}),
				layout.Rigid(func() {
					th.RadioButton("r3", "RadioButton3").Layout(gtx, radioButtonsGroup)
				}),
			)
		},
	}

	list.Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

