// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Editor struct {
	t *Theme
	material.EditorStyle

	TitleLabel      Label
	errorLabel      Label
	LineColor       color.NRGBA
	TitleLabelColor color.NRGBA

	//IsRequired if true, displays a required field text at the buttom of the editor.
	IsRequired bool
	//IsTitleLabel if true makes the title label visible.
	IsTitleLabel bool
	//Bordered if true makes the adds a border around the editor.
	IsCustomButton bool
	CustomButton   Button

	Bordered bool
	//IsPassword if true, displays the show and hide button.
	isPassword bool

	requiredErrorText string

	showHidePassword IconButton

	m2 unit.Value
	m5 unit.Value
}

func (t *Theme) EditorPassword(editor *widget.Editor, hint string) Editor {
	editor.Mask = '*'
	e := t.Editor(editor, hint)
	e.isPassword = true
	return e
}

func (t *Theme) Editor(editor *widget.Editor, hint string) Editor {
	errorLabel := t.Caption("")
	errorLabel.Color = t.Color.Danger

	m := material.Editor(t.Base, editor, hint)
	m.TextSize = t.TextSize
	m.Color = t.Color.Text
	m.Hint = hint
	m.HintColor = t.Color.Hint

	var m0 = unit.Dp(0)

	return Editor{
		t:                 t,
		EditorStyle:       m,
		TitleLabel:        t.Body2(""),
		IsTitleLabel:      true,
		Bordered:          true,
		LineColor:         t.Color.BorderColor,
		TitleLabelColor:   t.Color.IconGray,
		errorLabel:        errorLabel,
		requiredErrorText: "Field is required",

		m2: unit.Dp(2),
		m5: unit.Dp(5),

		showHidePassword: IconButton{
			material.IconButtonStyle{
				Icon:       mustIcon(widget.NewIcon(icons.ActionVisibilityOff)),
				Size:       values.MarginPadding24,
				Background: color.NRGBA{},
				Color:      t.Color.Gray,
				Inset:      layout.UniformInset(m0),
				Button:     new(widget.Clickable),
			},
		},
		CustomButton: t.Button(new(widget.Clickable), ""),
	}
}

func (e Editor) Layout(gtx layout.Context) layout.Dimensions {
	e.handleEvents()

	if e.Editor.Len() > 0 {
		e.TitleLabel.Text = e.Hint
	}

	if e.Editor.Focused() {
		e.TitleLabel.Text = e.Hint
		e.TitleLabel.Color = e.TitleLabelColor
		e.Hint = ""
	}

	if e.IsRequired && !e.Editor.Focused() && e.Editor.Len() == 0 {
		e.errorLabel.Text = e.requiredErrorText
		e.LineColor = e.t.Color.Danger
	}

	if e.errorLabel.Text != "" {
		e.LineColor, e.TitleLabel.Color = e.t.Color.Danger, e.t.Color.Danger
	}

	return layout.UniformInset(e.m2).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return e.editorLayout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								if e.errorLabel.Text != "" {
									inset := layout.Inset{
										Top:  e.m2,
										Left: e.m5,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return e.errorLabel.Layout(gtx)
									})
								}
								return layout.Dimensions{}
							}),
						)
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						if e.IsTitleLabel {
							return layout.Inset{
								Top:  values.MarginPaddingMinus10,
								Left: values.MarginPadding10,
							}.Layout(gtx, func(gtx C) D {
								return Card{Color: e.t.Color.Surface}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return e.TitleLabel.Layout(gtx)
								})
							})
						}
						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}

func (e Editor) editorLayout(gtx C) D {
	if e.Bordered {
		border := widget.Border{Color: e.LineColor, CornerRadius: e.m5, Width: unit.Dp(2)}
		return border.Layout(gtx, func(gtx C) D {
			inset := layout.Inset{
				Top:    e.m5,
				Bottom: e.m5,
				Left:   values.MarginPadding10,
				Right:  e.m5,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return e.editor(gtx)
			})
		})
	}

	return e.editor(gtx)
}

func (e Editor) editor(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Top:    e.m5,
						Bottom: e.m5,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return e.EditorStyle.Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			if e.isPassword {
				inset := layout.Inset{
					Top:  e.m2,
					Left: e.m5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					icon := mustIcon(widget.NewIcon(icons.ActionVisibilityOff))
					if e.Editor.Mask == '*' {
						icon = mustIcon(widget.NewIcon(icons.ActionVisibility))
					}
					e.showHidePassword.Icon = icon
					return e.showHidePassword.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			if e.IsCustomButton {
				inset := layout.Inset{
					Top:   e.m5,
					Left:  e.m5,
					Right: e.m5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					e.CustomButton.TextSize = unit.Dp(10)
					return e.CustomButton.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (e Editor) handleEvents() {
	if e.showHidePassword.Button.Clicked() {
		if e.Editor.Mask == '*' {
			e.Editor.Mask = 0
		} else if e.Editor.Mask == 0 {
			e.Editor.Mask = '*'
		}
	}

	if e.errorLabel.Text != "" {
		e.LineColor = e.t.Color.Danger
	} else {
		e.LineColor = e.t.Color.BorderColor
	}

	if e.requiredErrorText != "" {
		e.LineColor = e.t.Color.Danger
	} else {
		e.LineColor = e.t.Color.BorderColor
	}
}

func (e *Editor) SetRequiredErrorText(txt string) {
	e.requiredErrorText = txt
}

func (e *Editor) SetError(text string) {
	e.errorLabel.Text = text
}

func (e *Editor) ClearError() {
	e.errorLabel.Text = ""
}

func (e *Editor) IsDirty() bool {
	return e.errorLabel.Text == ""
}
