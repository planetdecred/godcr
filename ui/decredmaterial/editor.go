// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/atotto/clipboard"
	"github.com/raedahgroup/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type LineStyle uint8

const (
	RoundedRectangle LineStyle = iota
	SingleUnderLine
	NoLine
)

type Editor struct {
	t *Theme
	material.EditorStyle

	LineStyle LineStyle

	TitleLabel Label
	ErrorLabel Label
	LineColor  color.RGBA

	flexWidth float32
	//IsVisible if true, displays the paste and clear button.
	IsVisible bool
	//IsRequired if true, displays a required field text at the buttom of the editor.
	IsRequired bool
	//IsTitleLabel if true makes the title label visible.
	IsTitleLabel bool

	requiredErrorText string

	pasteBtnMaterial IconButton
	clearBtMaterial  IconButton
}

func (t *Theme) Editor(editor *widget.Editor, hint string) Editor {
	errorLabel := t.Caption("")
	errorLabel.Color = t.Color.Danger

	m := material.Editor(t.Base, editor, hint)
	m.TextSize = t.TextSize
	m.Color = t.Color.Text
	m.Hint = hint
	m.HintColor = t.Color.Hint

	return Editor{
		t:                 t,
		EditorStyle:       m,
		TitleLabel:        t.Body2(""),
		flexWidth:         0,
		IsTitleLabel:      true,
		LineColor:         t.Color.Hint,
		ErrorLabel:        errorLabel,
		requiredErrorText: "Field is required",

		pasteBtnMaterial: IconButton{
			material.IconButtonStyle{
				Icon:       mustIcon(widget.NewIcon(icons.ContentContentPaste)),
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},

		clearBtMaterial: IconButton{
			material.IconButtonStyle{
				Icon:       mustIcon(widget.NewIcon(icons.ContentClear)),
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
	}
}

func (e Editor) Layout(gtx layout.Context) layout.Dimensions {
	e.handleEvents()
	if e.IsVisible {
		e.flexWidth = 20
	}
	if e.Editor.Focused() || e.Editor.Len() != 0 {
		e.TitleLabel.Text = e.Hint
		e.LineColor = color.RGBA{41, 112, 255, 255}
		e.Hint = ""
	}

	if e.IsRequired && !e.Editor.Focused() && e.Editor.Len() == 0 {
		e.ErrorLabel.Text = e.requiredErrorText
		e.LineColor = e.t.Color.Danger
	}

	if e.ErrorLabel.Text != "" && e.Editor.Focused() && e.Editor.Len() != 0 {
		e.LineColor = e.t.Color.Danger
	}

	return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if e.IsTitleLabel {
					if e.Editor.Focused() {
						e.TitleLabel.Color = color.RGBA{41, 112, 255, 255}
					}
					return e.TitleLabel.Layout(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return e.editorLayout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								if e.ErrorLabel.Text != "" {
									inset := layout.Inset{
										Top: values.MarginPadding2,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return e.ErrorLabel.Layout(gtx)
									})
								}
								return layout.Dimensions{}
							}),
						)
					}),
				)
			}),
		)
	})
}

func (e Editor) editorLayout(gtx C) D {
	var dims layout.Dimensions
	switch e.LineStyle {
	case RoundedRectangle:
		dims = e.editorRectangle(gtx, func(gtx C) D {
			return e.editorSection(gtx, false)
		})
	case SingleUnderLine:
		dims = e.editorSection(gtx, true)
	case NoLine:
		dims = e.editorSection(gtx, false)
	}
	return dims
}

func (e Editor) editorSection(gtx layout.Context, underline bool) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding5
					inset := layout.Inset{
						Top:    m,
						Bottom: m,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return e.EditorStyle.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if underline {
						return e.editorLine(gtx)
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			if e.IsVisible {
				inset := layout.Inset{
					Top:  values.MarginPadding2,
					Left: values.MarginPadding5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					if e.Editor.Text() == "" {
						return e.pasteBtnMaterial.Layout(gtx)
					}
					return e.clearBtMaterial.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (e Editor) editorRectangle(gtx layout.Context, body layout.Widget) layout.Dimensions {
	border := widget.Border{Color: e.LineColor, CornerRadius: values.MarginPadding7, Width: values.MarginPadding1}
	return border.Layout(gtx, func(gtx C) D {
		mtb := values.MarginPadding2
		mlr := values.MarginPadding5
		return layout.Inset{Top: mtb, Bottom: mtb, Left: mlr, Right: mlr}.Layout(gtx, body)
	})
}

func (e Editor) editorLine(gtx C) D {
	line := e.t.Line()
	line.Color = e.LineColor
	line.Height = 2
	line.Width = gtx.Constraints.Max.X
	return line.Layout(gtx)
}

func (e Editor) handleEvents() {
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

	if e.ErrorLabel.Text != "" {
		e.LineColor = e.t.Color.Danger
	} else {
		e.LineColor = e.t.Color.Hint
	}

	if e.requiredErrorText != "" {
		e.LineColor = e.t.Color.Danger
	} else {
		e.LineColor = e.t.Color.Hint
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
