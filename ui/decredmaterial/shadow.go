package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Shadow struct {
	color  color.NRGBA
	right  bool
	left   bool
	top    bool
	bottom bool
}

const (
	shadowElevation = 4
)

func (t *Theme) Shadow(left, right, top, bottom bool) *Shadow {
	return &Shadow{
		color:  color.NRGBA{A: 0x10},
		right:  right,
		left:   left,
		top:    top,
		bottom: bottom,
	}
}

func (s *Shadow) Layout(gtx C, startX, startY int, wdgt func(gtx C) D) D {
	dims := wdgt(gtx)
	size := dims.Size

	if s.top {
		rect := image.Rectangle{
			Min: image.Point{
				Y: startY - shadowElevation,
				X: startX,
			},
			Max: image.Point{
				X: size.X,
				Y: startY,
			},
		}
		s.drawShadowOverRectangle(gtx, rect)
	}

	if s.bottom {
		bottomStartX := startX
		bottomEndX := size.X
		if s.left {
			bottomStartX -= shadowElevation
		}

		if s.right {
			bottomEndX += shadowElevation
		}

		rect := image.Rectangle{
			Min: image.Point{
				Y: size.Y,
				X: bottomStartX,
			},
			Max: image.Point{
				X: bottomEndX,
				Y: size.Y + shadowElevation,
			},
		}
		s.drawShadowOverRectangle(gtx, rect)
	}

	if s.left {
		rect := image.Rectangle{
			Min: image.Point{
				Y: startY,
				X: startX - shadowElevation,
			},
			Max: image.Point{
				X: startX,
				Y: size.Y,
			},
		}
		s.drawShadowOverRectangle(gtx, rect)
	}

	if s.right {
		rect := image.Rectangle{
			Min: image.Point{
				Y: startY,
				X: size.X,
			},
			Max: image.Point{
				X: size.X + shadowElevation,
				Y: size.Y,
			},
		}
		s.drawShadowOverRectangle(gtx, rect)
	}
	return D{}
}

func (s *Shadow) drawShadowOverRectangle(gtx C, rect image.Rectangle) {
	transparent := s.color
	transparent.A = 0

	defer op.Save(gtx.Ops).Load()
	clip.Rect(rect).Add(gtx.Ops)
	paint.LinearGradientOp{
		Color1: s.color,
		Color2: transparent,
		Stop1: layout.FPt(image.Point{
			X: rect.Min.X,
			Y: rect.Max.Y,
		}),
		Stop2: layout.FPt(image.Point{
			X: rect.Min.X,
			Y: rect.Max.Y,
		}),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
