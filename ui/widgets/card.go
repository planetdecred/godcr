package widgets

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"image/color"

	"github.com/raedahgroup/godcr-gio/ui/values"
)

// Card is a widget that uses a rectangle as flat colored background
type Card struct {
	height int
	width  int
	color  color.RGBA
}

// SetColor specifies the color of the card widget
func (c *Card) SetColor(col color.RGBA) {
	c.color = col
}

// SetHeight sets the height of the card widget
func (c *Card) SetHeight(height int) {
	c.height = height
}

// SetWidth sets the width of the card widget
func (c *Card) SetWidth(width int) {
	c.width = width
}

// Layout defines the rectangle with dimensions, fills it with a color
// and draws it.
func (c *Card) Layout(gtx *layout.Context) {
	if c.width == 0 {
		c.width = gtx.Constraints.Width.Max
	}
	if c.height == 0 {
		c.height = gtx.Constraints.Height.Max
	}

	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(c.width),
			Y: float32(c.height),
		},
	}
	Fill(gtx, c.color, rect)
}

// NewCard creates a new card object
func NewCard() *Card {
	return &Card{
		color: values.DefaultCardGray,
	}
}
