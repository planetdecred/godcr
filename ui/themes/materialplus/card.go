package materialplus

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"

	"github.com/raedahgroup/godcr-gio/ui/values"
)

// Card is a widget that uses a rectangle as flat colored background
type Card struct {
	Height int
	Width  int
	Color  color.RGBA
}

// Layout defines the rectangle with dimensions, fills it with a color
// and draws it.
func (c *Card) Layout(gtx *layout.Context, borderRadius float32) {
	br := borderRadius
	if c.Width == 0 {
		c.Width = gtx.Constraints.Width.Max
	}
	if c.Height == 0 {
		c.Height = gtx.Constraints.Height.Max
	}

	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(c.Width),
			Y: float32(c.Height),
		},
	}
	if br > 0 {
		clip.Rect{
			Rect: rect,
			NE:   br, NW: br, SE: br, SW: br,
		}.Op(gtx.Ops).Add(gtx.Ops)
	}

	Fill(gtx, c.Color, c.Width, c.Height)
}

// Card returns an instance of Card
func (t *Theme) Card() Card {
	return Card{
		Color: values.DefaultCardGray,
	}
}
