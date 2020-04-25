// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
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
	theme *Theme
	//title is the title of the editor input field
	title     string
	LineColor color.RGBA

	editorMaterial Editor
	flexWidth      float32
	editor         *widget.Editor

	IsVisibleBtn     bool
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

func (t *Theme) EditorCustom(hint string, editor *widget.Editor) EditorCustom {
	return EditorCustom{
		theme:     t,
		title:     hint,
		flexWidth: 1,
		LineColor: t.Color.Text,

		editorMaterial: t.Editor(hint),
		editor:         editor,

		pasteBtnMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentContentPaste)),
			Size:       unit.Dp(30),
			Background: color.RGBA{0, 0, 0, 0},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},

		clearBtMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentClear)),
			Size:       unit.Dp(30),
			Background: color.RGBA{0, 0, 0, 0},
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
	if e.IsVisibleBtn {
		e.flexWidth = 0.93
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			if e.editor.Text() == "" {
				e.theme.Body1("").Layout(gtx)
			} else {
				e.theme.Body1(e.title).Layout(gtx)
			}
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							inset := layout.Inset{
								Top:    unit.Dp(6),
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
					)
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(10),
					}
					inset.Layout(gtx, func() {
						if e.IsVisibleBtn {
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
}

func (e EditorCustom) handleEvents(gtx *layout.Context) {
	data, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	for e.pasteBtnWidget.Clicked(gtx) {
		e.editor.SetText(data)
	}
	for e.clearBtnWidget.Clicked(gtx) {
		e.editor.SetText("")
	}

	for _, evt := range e.editor.Events(gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			e.editorMaterial.HintColor = e.theme.Color.Hint
			return
		}
	}
}
