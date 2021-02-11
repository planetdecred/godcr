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
	t     *Theme
	Color color.NRGBA

	activeBtn, inactiveBtn *widget.Clickable

	activeCard, inactiveCard Card

	inactivetxt          string
	activeTxt            string
	isActive, isInactive bool
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

	return sw
}

func (s *SwitchButtonText) Layout(gtx layout.Context) layout.Dimensions {
	s.handleClickEvent(gtx)
	card := s.t.Card()
	card.Color = s.t.Color.LightGray
	m10 := unit.Dp(10)
	m5 := unit.Dp(5)
	return card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return s.activeCard.Layout(gtx, func(gtx C) D {
								return layout.Inset{
									Left:   m10,
									Bottom: m5,
									Right:  m10,
									Top:    m5,
								}.Layout(gtx, func(gtx C) D {
									txt := s.t.Body1(s.activeTxt)
									txt.Color = s.t.Color.Text
									if !s.isActive {
										txt.Color = s.t.Color.Gray
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
									Left:   m10,
									Bottom: m5,
									Right:  m10,
									Top:    m5,
								}.Layout(gtx, func(gtx C) D {
									txt := s.t.Body1(s.inactivetxt)
									txt.Color = s.t.Color.Text
									if !s.isInactive {
										txt.Color = s.t.Color.Gray
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

func (s *SwitchButtonText) handleClickEvent(gtx layout.Context) {
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
