package decredmaterial

import (
	"image"

	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"
)

type Image struct {
	*widget.Image
}

func NewImage(src image.Image) *Image {
	return &Image{
		Image: &widget.Image{
			Src: paint.NewImageOp(src),
		},
	}
}

func (img *Image) Layout16dp(gtx C) D {
	return img.LayoutSize(gtx, values.MarginPadding16)
}

func (img *Image) Layout24dp(gtx C) D {
	return img.LayoutSize(gtx, values.MarginPadding24)
}

func (img *Image) Layout36dp(gtx C) D {
	return img.LayoutSize(gtx, values.MarginPadding36)
}

func (img *Image) Layout48dp(gtx C) D {
	return img.LayoutSize(gtx, values.MarginPadding48)
}

func (img *Image) LayoutSize(gtx C, size unit.Value) D {
	width := float32(img.Src.Size().X)
	scale := size.V / width
	img.Scale = scale
	return img.Layout(gtx)
}
