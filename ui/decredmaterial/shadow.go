package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Shadow struct {
	color     color.NRGBA
	right     bool
	left      bool
	top       bool
	bottom    bool
}

const (
	shadowElevation = 3
)

func (t *Theme) Shadow(col color.NRGBA, left, right, top, bottom bool) *Shadow {
	return &Shadow{
		color:     col,
		right:     right,
		left:      left,
		top:       top,
		bottom:    bottom,
	}
}

func (s *Shadow) Layout(gtx C, startX, startY int, wdgt func(gtx C) D) D {
	dims := wdgt(gtx)
	s.layoutShadow(gtx, startX, startY, dims.Size)

	return D{}
}

func (s *Shadow) layoutShadow(gtx layout.Context, startX, startY int, size image.Point) {
	transparent := s.color
	transparent.A = 0

	col := color.NRGBA{A: 0x20}

	if s.left {
		func() {
			defer op.Save(gtx.Ops).Load()
			leftRect := image.Rectangle{
				Min: image.Point{
					Y: startY,
					X: startX,
				},
				Max: image.Point{
					X: startX + shadowElevation,
					Y: size.Y,
				},
			}
			clip.Rect(leftRect).Add(gtx.Ops)
			paint.LinearGradientOp{
				Color1: col,
				Color2: transparent,
				Stop1: layout.FPt(image.Point{
					X: leftRect.Min.X,
					Y: leftRect.Max.Y,
				}),
				Stop2: layout.FPt(image.Point{
					X: leftRect.Min.X,
					Y: leftRect.Max.Y,
				}),
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}()
	}

	if s.bottom {
		func() {
			defer op.Save(gtx.Ops).Load()
			bottomRect := image.Rectangle{
				Min: image.Point{
					Y: size.Y,
					X: startX,
				},
				Max: image.Point{
					X: size.X,
					Y: size.Y + shadowElevation,
				},
			}
			clip.Rect(bottomRect).Add(gtx.Ops)
			paint.LinearGradientOp{
				Color1: col,
				Color2: transparent,
				Stop1: layout.FPt(image.Point{
					X: bottomRect.Min.X,
					Y: bottomRect.Max.Y,
				}),
				Stop2: layout.FPt(image.Point{
					X: bottomRect.Min.X,
					Y: bottomRect.Max.Y,
				}),
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}()
	}

	if s.right {
		func(){
			defer op.Save(gtx.Ops).Load()
			rightRect := image.Rectangle{
				Min: image.Point{
					Y: startY,
					X: size.X,
				},
				Max: image.Point{
					X: size.X + shadowElevation,
					Y: size.Y,
				},
			}
			clip.Rect(rightRect).Add(gtx.Ops)
			paint.LinearGradientOp{
				Color1: col,
				Color2: transparent,
				Stop1: layout.FPt(image.Point{
					X: rightRect.Min.X,
					Y: rightRect.Max.Y,
				}),
				Stop2: layout.FPt(image.Point{
					X: rightRect.Min.X,
					Y: rightRect.Max.Y,
				}),
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}()
	}

	if s.top {
		func() {
			defer op.Save(gtx.Ops).Load()
			bottomRect := image.Rectangle{
				Min: image.Point{
					Y: startY - shadowElevation,
					X: startX,
				},
				Max: image.Point{
					X: size.X,
					Y: startY,
				},
			}
			clip.Rect(bottomRect).Add(gtx.Ops)
			paint.LinearGradientOp{
				Color1: col,
				Color2: transparent,
				Stop1: layout.FPt(image.Point{
					X: bottomRect.Min.X,
					Y: bottomRect.Max.Y,
				}),
				Stop2: layout.FPt(image.Point{
					X: bottomRect.Min.X,
					Y: bottomRect.Max.Y,
				}),
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}()
	}
}
