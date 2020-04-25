// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/f32"

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
	titleLabel Label
	LineColor color.RGBA

	editorMaterial Editor
	editor         *widget.Editor

	pasteButtonMaterial IconButton
	pasteButtonWidget   *widget.Button

	clearButtonMaterial IconButton
	clearButtonWidget   *widget.Button
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

func (t *Theme) EditorCustom(hint, title string, editor *widget.Editor) EditorCustom {
	return EditorCustom{
		titleLabel: t.Body1(title),
		LineColor: t.Color.Text,

		editorMaterial: t.Editor(hint),
		editor:         editor,

		pasteButtonMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentContentPaste)),
			Size:       unit.Dp(30),
			Background: color.RGBA{0, 0, 0, 0},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},

		clearButtonMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentClear)),
			Size:       unit.Dp(30),
			Background: color.RGBA{0, 0, 0, 0},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},
		pasteButtonWidget: new(widget.Button),
		clearButtonWidget: new(widget.Button),
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
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			e.titleLabel.Layout(gtx)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							inset := layout.Inset{
								Top: unit.Dp(6),
								Bottom: unit.Dp(4),
							}
							inset.Layout(gtx, func() {
								layout.Flex{}.Layout(gtx,
									layout.Flexed(0.9, func() {
										e.editorMaterial.Layout(gtx, e.editor)
									}),
								)
							})
						}),
						layout.Rigid(func(){
							layout.Flex{}.Layout(gtx,
								layout.Flexed(0.9, func() {
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
					)
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(10),
					}
					inset.Layout(gtx, func() {
						if e.editor.Text() == "" {
							e.pasteButtonMaterial.Layout(gtx, e.pasteButtonWidget)
						} else {
							e.clearButtonMaterial.Layout(gtx, e.clearButtonWidget)
						}
					})
				}),
			)
		}),
	)
}

func (e EditorCustom) handleEvents(gtx *layout.Context) {
	data, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	for e.pasteButtonWidget.Clicked(gtx) {
		e.editor.SetText(data)
	}
	for e.clearButtonWidget.Clicked(gtx) {
		e.editor.SetText("")
	}
}
