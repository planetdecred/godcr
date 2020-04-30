// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	// "errors"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/atotto/clipboard"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Editor struct {
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA

	shaper text.Shaper
}

type EditorCustom struct {
	hint       string
	TitleLabel Label
	ErrorLabel Label
	LineColor  color.RGBA

	editorMaterial Editor
	flexWidth      float32
	editor         *widget.Editor
	//IsVisible if false hides the paste and clear button.
	IsVisible bool
	//Required if true the input field must contain data and cannot be empty.
	Required         bool
	pasteBtnMaterial IconButton
	pasteBtnWidget   *widget.Button

	clearBtMaterial IconButton
	clearBtnWidget  *widget.Button
}

func (t *Theme) Editor(hint string) Editor {
	return Editor{
		TextSize:  t.TextSize,
		Color:     t.Color.Text,
		shaper:    t.Shaper,
		Hint:      hint,
		HintColor: t.Color.Hint,
	}
}

func (t *Theme) EditorCustom(hint string) EditorCustom {
	errorLabel := t.Caption("Field is required")
	errorLabel.Color = color.RGBA{255, 0, 0, 255}

	return EditorCustom{
		TitleLabel: t.Body1(""),
		// selected:   false,
		flexWidth:      1,
		hint:           hint,
		LineColor:      t.Color.Text,
		ErrorLabel:     errorLabel,
		editorMaterial: t.Editor(hint),
		editor:         new(widget.Editor),

		pasteBtnMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentContentPaste)),
			Size:       unit.Dp(30),
			Background: color.RGBA{237, 237, 237, 255},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},

		clearBtMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentClear)),
			Size:       unit.Dp(30),
			Background: color.RGBA{237, 237, 237, 255},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},
		pasteBtnWidget: new(widget.Button),
		clearBtnWidget: new(widget.Button),
	}
}

func (e Editor) Layout(gtx *layout.Context, editor *widget.Editor) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := widget.Label{Alignment: editor.Alignment}
	tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	macro.Stop()
	if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
		gtx.Constraints.Width.Min = w
	}
	if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
		gtx.Constraints.Height.Min = h
	}
	editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	if editor.Len() > 0 {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		editor.PaintText(gtx)
	} else {
		macro.Add()
	}
	paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
	editor.PaintCaret(gtx)
	stack.Pop()
}

func (e EditorCustom) Layout(gtx *layout.Context) {
	e.handleEvents(gtx)

	if e.IsVisible {
		e.flexWidth = 0.93
	}

	layout.UniformInset(unit.Dp(1)).Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				if e.editor.Focused() || e.editor.Len() != 0 {
					e.TitleLabel.Text = e.hint
					e.editorMaterial.Hint = ""
				}
				e.TitleLabel.Layout(gtx)
			}),
			layout.Rigid(func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func() {
								inset := layout.Inset{
									Top:    unit.Dp(4),
									Bottom: unit.Dp(4),
								}
								inset.Layout(gtx, func() {
									layout.Flex{}.Layout(gtx,
										layout.Flexed(e.flexWidth, func() {
											e.editorMaterial.Layout(gtx, e.editor)
										}),
									)
								})
							}),
							layout.Rigid(func() {
								// if e.Required {
								// 	if !e.editor.Focused() && e.editor.Len() == 0 {
								// 		e.LineColor = color.RGBA{255, 0, 0, 255}
								// 	}
								// }
								layout.Flex{}.Layout(gtx,
									layout.Flexed(e.flexWidth, func() {
										rect := f32.Rectangle{
											Max: f32.Point{
												X: float32(gtx.Constraints.Width.Max),
												Y: 1,
											},
										}
										op.TransformOp{}.Offset(f32.Point{
											X: 0,
											Y: 0,
										}).Add(gtx.Ops)
										paint.ColorOp{Color: e.LineColor}.Add(gtx.Ops)
										paint.PaintOp{Rect: rect}.Add(gtx.Ops)
									}),
								)
							}),
							layout.Rigid(func() {
								if e.Required {
									// if e.editor.Len() == 0 && !e.editor.Focused(){
									// e.ErrorLabel.Text = "Field is required"
									// }
									e.ErrorLabel.Layout(gtx)
								}

							}),
						)
					}),
					layout.Rigid(func() {
						inset := layout.Inset{
							Left: unit.Dp(10),
						}
						inset.Layout(gtx, func() {
							if e.IsVisible {
								if e.editor.Text() == "" {
									e.pasteBtnMaterial.Layout(gtx, e.pasteBtnWidget)
								} else {
									e.clearBtMaterial.Layout(gtx, e.clearBtnWidget)
								}
							}
						})
					}),
				)
			}),
		)
	})
}

func (e EditorCustom) Text() string {
	if e.Required && e.editor.Len() == 0 && !e.editor.Focused() {
		e.ErrorLabel.Text = "Field is required and cannot be empty."
		e.LineColor = color.RGBA{255, 0, 0, 255}
		return ""
	}
	return e.editor.Text()
}

func (e EditorCustom) handleEvents(gtx *layout.Context) {
	data, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}
	for e.pasteBtnWidget.Clicked(gtx) {
		e.editor.SetText(data)
		e.reset(gtx)
	}
	for e.clearBtnWidget.Clicked(gtx) {
		e.editor.SetText("")
		e.reset(gtx)
	}
}

func (e EditorCustom) reset(gtx *layout.Context) {
	for _, evt := range e.editor.Events(gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			e.LineColor = darkblue
			return
		}
	}
}
