package decredmaterial

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type Hoverable struct {
	hovered  bool
	position *f32.Point
}

func (t *Theme) Hoverable() *Hoverable {
	return &Hoverable{}
}

func (h *Hoverable) Hovered() bool {
	return h.hovered
}

func (h *Hoverable) Position() *f32.Point {
	return h.position
}

func (h *Hoverable) update(gtx C) {
	for _, e := range gtx.Events(h) {
		ev, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch ev.Type {
		case pointer.Enter:
			h.hovered = true
			h.position = &ev.Position
		case pointer.Leave:
			h.hovered = false
			h.position = &f32.Point{}
		}
	}
}

func (h *Hoverable) Layout(gtx C) D {
	h.update(gtx)

	defer op.Save(gtx.Ops).Load()

	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   h,
		Types: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: gtx.Constraints.Min,
	}
}
