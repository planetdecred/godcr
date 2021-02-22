package decredmaterial

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type Image struct {
	imageOp paint.ImageOp
}

const (
	scaleFactor = 25
)

func NewImage(img image.Image) *Image {
	return &Image{
		imageOp: paint.NewImageOp(img),
	}
}

func (i *Image) layout(gtx layout.Context, sf int) D {

	pt := f32.Point{
		X: float32(sf) / float32(i.imageOp.Size().X),
		Y: float32(sf) / float32(i.imageOp.Size().Y),
	}
	op.Affine(
		f32.Affine2D{}.Scale(f32.Point{}, pt),
	).Add(gtx.Ops)
	i.imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{
		Size: image.Point{
			X: sf,
			Y: sf,
		},
	}
}

func (i *Image) Layout(gtx layout.Context) D {
	return i.layout(gtx, scaleFactor)
}

func (i *Image) LayoutWithScaleFactor(gtx layout.Context, sf int) D {
	return i.layout(gtx, sf)
}
