// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/widget/material"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/atotto/clipboard"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/math/fixed"
)

type Editor struct {
	material.EditorStyle

	shaper           text.Shaper
	StartHint        string
	TitleLabel       Label
	ErrorLabel       Label
	LineColor        color.RGBA
	LineColorFocused color.RGBA

	flexWidth float32
	//IsVisible if true, displays the paste and clear button.
	IsVisible bool
	//IsRequired if true, displays a required field text at the buttom of the editor.
	IsRequired bool

	requiredErrorText string

	pasteBtnMaterial IconButton
	clearBtMaterial  IconButton

	hasCalculated                                   bool
	editorLineHeight                                int
	titleLabelDims, editorDims, lineDims, errorDims layout.Dimensions
}

func (t *Theme) Editor(editor *widget.Editor, hint string) *Editor {
	errorLabel := t.Caption("")
	errorLabel.Color = t.Color.Danger

	m := material.Editor(t.Base, editor, hint)
	m.TextSize = t.TextSize
	m.Color = t.Color.Text
	m.HintColor = t.Color.Hint

	e := Editor{
		EditorStyle:       m,
		shaper:            t.Shaper,
		TitleLabel:        t.Body2(""),
		StartHint:         hint,
		flexWidth:         1,
		LineColor:         t.Color.Text,
		LineColorFocused:  t.Color.Primary,
		ErrorLabel:        errorLabel,
		requiredErrorText: "Field is required",
		hasCalculated:     false,

		pasteBtnMaterial: IconButton{
			material.IconButtonStyle{
				Icon:       mustIcon(widget.NewIcon(icons.ContentContentPaste)),
				Size:       unit.Dp(30),
				Background: color.RGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(unit.Dp(5)),
				Button:     new(widget.Clickable),
			},
		},

		clearBtMaterial: IconButton{
			material.IconButtonStyle{
				Icon:       mustIcon(widget.NewIcon(icons.ContentClear)),
				Size:       unit.Dp(30),
				Background: color.RGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(unit.Dp(5)),
				Button:     new(widget.Clickable),
			},
		},
	}
	e.TitleLabel.Text = m.Hint

	return &e
}

func (e *Editor) calculateDims(gtx layout.Context) {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			e.titleLabelDims = e.TitleLabel.Layout(gtx)
			return e.titleLabelDims
		}),
		layout.Rigid(func(gtx C) D {
			e.editorDims = e.EditorStyle.Layout(gtx)
			e.editorLineHeight = e.editorDims.Size.Y
			return e.editorDims
		}),
		layout.Rigid(func(gtx C) D {
			e.lineDims = layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
			}.Layout(gtx, func(gtx C) D {
				return e.editorLine(gtx)
			})
			return e.lineDims
		}),
	)
}

func (e *Editor) calculateEditorHeight(gtx layout.Context) {
	r := strings.NewReader(e.Editor.Text())
	lines, _ := e.shaper.Layout(e.Font, fixed.I(gtx.Px(e.TextSize)), gtx.Constraints.Max.X, r)
	e.editorDims.Size.Y = e.editorLineHeight * len(lines)
}

// Layout renders the editor to screen. The editor line is able to retain
// it's relative position whether or not the hint or title labels are displayed
// or not because their dimensions are pre-calculated before hand
func (e *Editor) Layout(gtx layout.Context) layout.Dimensions {
	if !e.hasCalculated {
		e.calculateDims(gtx)
		e.hasCalculated = true
	}

	if !e.Editor.SingleLine {
		e.calculateEditorHeight(gtx)
	}

	e.handleEvents()
	if e.IsVisible {
		e.flexWidth = 0.93
	}

	if e.Editor.Focused() || e.Editor.Len() != 0 {
		e.Hint = ""
	}

	if !e.Editor.Focused() {
		e.Hint = e.StartHint
	}

	if e.IsRequired && !e.Editor.Focused() && e.Editor.Len() == 0 {
		e.ErrorLabel.Text = e.requiredErrorText
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if e.Editor.Focused() || e.Editor.Len() != 0 {
				e.TitleLabel.Color = color.RGBA{41, 112, 255, 255}
				e.TitleLabel.Layout(gtx)
			}
			return e.titleLabelDims
		}),
		layout.Rigid(func(gtx C) D {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(0.9, func(gtx C) D {
					return e.EditorStyle.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top: unit.Dp(-15),
					}.Layout(gtx, func(gtx C) D {
						if e.IsVisible {
							if e.Editor.Text() == "" {
								return e.pasteBtnMaterial.Layout(gtx)
							}
							return e.clearBtMaterial.Layout(gtx)
						}
						return layout.Dimensions{}
					})
				}),
			)
			return e.editorDims
		}),
		layout.Rigid(func(gtx C) D {
			layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
			}.Layout(gtx, func(gtx C) D {
				return e.editorLine(gtx)
			})
			return e.lineDims
		}),
		layout.Rigid(func(gtx C) D {
			if e.ErrorLabel.Text != "" {
				inset := layout.Inset{
					Top: unit.Dp(3),
				}
				return inset.Layout(gtx, func(gtx C) D {
					return e.ErrorLabel.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (e Editor) editorLine(gtx C) D {
	col := e.LineColor
	if e.Editor.Focused() {
		col = e.LineColorFocused
	}

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(e.flexWidth, func(gtx C) D {
			dims := image.Point{
				X: gtx.Constraints.Max.X,
				Y: 2,
			}
			rect := f32.Rectangle{
				Max: layout.FPt(dims),
			}
			op.Offset(f32.Point{
				X: 0,
				Y: 0,
			}).Add(gtx.Ops)
			paint.ColorOp{Color: col}.Add(gtx.Ops)
			paint.PaintOp{Rect: rect}.Add(gtx.Ops)
			return layout.Dimensions{Size: dims}
		}),
	)
}

func (e *Editor) handleEvents() {
	for e.pasteBtnMaterial.Button.Clicked() {
		data, err := clipboard.ReadAll()
		if err != nil {
			panic(err)
		}
		e.Editor.SetText(data)
	}
	for e.clearBtMaterial.Button.Clicked() {
		e.Editor.SetText("")
	}
}

func (e *Editor) SetRequiredErrorText(txt string) {
	e.requiredErrorText = txt
}

func (e *Editor) SetError(errorText string) {
	e.ErrorLabel.Text = errorText
}

func (e *Editor) ClearError() {
	e.ErrorLabel.Text = ""
}

func (e *Editor) IsDirty() bool {
	return e.ErrorLabel.Text == ""
}
