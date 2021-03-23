// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Switch struct {
	material.SwitchStyle
}

type SwitchButtonText struct {
	t                        *Theme
	activeTextColor          color.NRGBA
	inactiveTextColor        color.NRGBA
	activeBtn, inactiveBtn   *widget.Clickable
	activeCard, inactiveCard Card
	inactivetxt              string
	activeTxt                string
	isActive, isInactive     bool
}

func (t *Theme) Switch(swtch *widget.Bool) Switch {
	return Switch{material.Switch(t.Base, swtch)}
}

func (t *Theme) SwitchButtonText(activeTxt, inactivetxt string, activeBtn, inactiveBtn *widget.Clickable) *SwitchButtonText {
	sw := &SwitchButtonText{
		t:           t,
		activeBtn:   activeBtn,
		inactiveBtn: inactiveBtn,

		inactivetxt:  inactivetxt,
		activeTxt:    activeTxt,
		isActive:     true,
		activeCard:   t.Card(),
		inactiveCard: t.Card(),
	}

	sw.activeCard.Color = sw.t.Color.Surface
	sw.inactiveCard.Color = color.NRGBA{}

	raduis := CornerRadius{NE: 7, NW: 7, SE: 7, SW: 7}
	sw.activeCard.Radius, sw.inactiveCard.Radius = raduis, raduis

	sw.activeTextColor = sw.t.Color.DeepBlue
	sw.inactiveTextColor = sw.t.Color.IconColor
	return sw
}

func (s *SwitchButtonText) Layout(gtx layout.Context) layout.Dimensions {
	s.handleClickEvent()
	m8 := unit.Dp(8)
	m4 := unit.Dp(4)
	card := s.t.Card()
	card.Color = s.t.Color.BorderColor
	card.Radius = CornerRadius{NE: 8, NW: 8, SE: 8, SW: 8}
	return card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return s.activeCard.Layout(gtx, func(gtx C) D {
								return layout.Inset{
									Left:   m8,
									Bottom: m4,
									Right:  m8,
									Top:    m4,
								}.Layout(gtx, func(gtx C) D {
									txt := s.t.Body2(s.activeTxt)
									txt.Color = s.activeTextColor
									if !s.isActive {
										txt.Color = s.inactiveTextColor
									}
									return txt.Layout(gtx)
								})
							})
						}),
						layout.Expanded(s.activeBtn.Layout),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return s.inactiveCard.Layout(gtx, func(gtx C) D {
								return layout.Inset{
									Left:   m8,
									Bottom: m4,
									Right:  m8,
									Top:    m4,
								}.Layout(gtx, func(gtx C) D {
									txt := s.t.Body2(s.inactivetxt)
									txt.Color = s.activeTextColor
									if !s.isInactive {
										txt.Color = s.inactiveTextColor
									}
									return txt.Layout(gtx)
								})
							})
						}),
						layout.Expanded(s.inactiveBtn.Layout),
					)
				}),
			)
		})
	})
}

func (s *SwitchButtonText) handleClickEvent() {
	for s.inactiveBtn.Clicked() {
		s.inactiveCard.Color = s.t.Color.Surface
		s.activeCard.Color = color.NRGBA{}
		s.isActive = false
		s.isInactive = true
	}

	for s.activeBtn.Clicked() {
		s.inactiveCard.Color = color.NRGBA{}
		s.activeCard.Color = s.t.Color.Surface
		s.isActive = true
		s.isInactive = false
	}
}
