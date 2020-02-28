package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Select struct {
	item     material.Button
	size     float32
	changed  bool
	open     bool
	selected int
	openbtn  widget.Button
	Options  []string
	btns     []*widget.Button
	list     layout.List
}

func (t *Theme) Select() *Select {

	return &Select{
		item: t.Button(""),
		size: 0.2,
		list: layout.List{Axis: layout.Vertical},
	}
}

func (s *Select) Selected() int {
	return s.selected
}

func (s *Select) Open() bool {
	return s.open
}

func (s *Select) Changed() bool {
	return s.changed
}

func (s *Select) Layout(gtx *layout.Context, w layout.Widget) {
	if s.openbtn.Clicked(gtx) {
		s.open = !s.open
	}
	if len(s.Options) != len(s.btns) {
		s.btns = make([]*widget.Button, len(s.Options))
		for i := range s.btns {
			s.btns[i] = new(widget.Button)
		}
	}
	layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Stacked(func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(s.size, func() {

				}),
				layout.Rigid(w),
			)
		}),
		layout.Stacked(func() {
			bd := func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func() {
						lbl := s.item
						lbl.Text = s.Options[s.selected]
						lbl.Layout(gtx, &s.openbtn)
					}),
					layout.Rigid(func() {
						if !s.open {
							return
						}
						(&s.list).Layout(gtx, len(s.Options), func(i int) {
							lbl := s.item
							lbl.Text = s.Options[i]
							lbl.Layout(gtx, s.btns[i])

						})
					}),
				)
			}
			layout.NE.Layout(gtx, bd)
		}),
	)

	for i := range s.btns {
		if s.btns[i].Clicked(gtx) {
			s.changed = true
			s.open = false
			s.selected = i
			return
		}
	}
	s.changed = false
}
