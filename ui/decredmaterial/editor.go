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

	TitleLabel Label
	ErrorLabel Label
	LineColor  color.RGBA

	flexWidth float32
	//IsVisible if true, displays the paste and clear button.
	IsVisible bool
	//IsRequired if true, displays a required field text at the buttom of the editor.
	IsRequired bool
	//IsTitleLabel if true makes the title lable visible.
	IsTitleLabel bool

	pasteBtnMaterial IconButton
	pasteBtnWidget   *widget.Button

	clearBtMaterial IconButton
	clearBtnWidget  *widget.Button

	editor *widget.Editor
}

func (t *Theme) Editor(hint string) Editor {
	errorLabel := t.Caption("")
	errorLabel.Color = t.Color.Danger

	return Editor{
		TextSize:     t.TextSize,
		Color:        t.Color.Text,
		shaper:       t.Shaper,
		Hint:         hint,
		HintColor:    t.Color.Hint,
		TitleLabel:   t.Body2(""),
		flexWidth:    1,
		IsTitleLabel: true,
		LineColor:    t.Color.Text,
		ErrorLabel:   errorLabel,

		pasteBtnMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentContentPaste)),
			Size:       unit.Dp(30),
			Background: color.RGBA{},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},

		clearBtMaterial: IconButton{
			Icon:       mustIcon(NewIcon(icons.ContentClear)),
			Size:       unit.Dp(30),
			Background: color.RGBA{},
			Color:      t.Color.Text,
			Padding:    unit.Dp(5),
		},
		pasteBtnWidget: new(widget.Button),
		clearBtnWidget: new(widget.Button),
	}
}

func (e Editor) Layout(gtx *layout.Context, editor *widget.Editor) {
	e.editor = editor
	e.handleEvents(gtx)
	if e.IsVisible {
		e.flexWidth = 0.93
	}
	if e.editor.Focused() || e.editor.Len() != 0 {
		e.TitleLabel.Text = e.Hint
		e.LineColor = color.RGBA{0, 137, 123, 255}
		e.Hint = ""
	}

	layout.UniformInset(unit.Dp(2)).Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				if e.IsTitleLabel {
					if e.editor.Focused() {
						e.TitleLabel.Color = color.RGBA{0, 137, 123, 255}
					}
					e.TitleLabel.Layout(gtx)
				}
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
											e.EditorWidget(gtx, editor)
										}),
									)
								})
							}),
							layout.Rigid(func() {
								e.editorLine(gtx)
							}),
							layout.Rigid(func() {
								if e.IsRequired {
									if !e.editor.Focused() && e.editor.Len() == 0 {
										e.ErrorLabel.Text = "Field is required"
									}
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

func (e Editor) EditorWidget(gtx *layout.Context, editor *widget.Editor) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.editor.Alignment}
	tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	macro.Stop()
	if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
		gtx.Constraints.Width.Min = w
	}
	if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
		gtx.Constraints.Height.Min = h
	}
	e.editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	if e.editor.Len() > 0 {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		e.editor.PaintText(gtx)
	} else {
		macro.Add()
	}
	paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
	e.editor.PaintCaret(gtx)
	stack.Pop()
}

func (e Editor) editorLine(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(e.flexWidth, func() {
			rect := f32.Rectangle{
				Max: f32.Point{
					X: float32(gtx.Constraints.Width.Max),
					Y: 2,
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
}

func (e Editor) Text() string {
	if e.IsRequired && e.editor.Len() == 0 && !e.editor.Focused() {
		e.ErrorLabel.Text = "Field is required and cannot be empty."
		e.LineColor = color.RGBA{255, 0, 0, 255}
		return ""
	}
	return e.editor.Text()
}

func (e Editor) SetText(text string) {
	if text == "" {
		return
	}
	e.editor.SetText(text)
}

func (e Editor) Len() int {
	return e.editor.Len()
}

func (e Editor) Clear() {
	e.editor.SetText("")
}

func (e Editor) Focused() bool {
	return e.editor.Focused()
}

func (e Editor) handleEvents(gtx *layout.Context) {
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
}
