package helper

import (
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
)

type (
	clickItem struct {
		position f32.Point
		time     time.Time
	}

	Clicker struct {
		click      gesture.Click
		clicks     int
		prevClicks int
		history    []clickItem
	}
)

func NewClicker() Clicker {
	return Clicker{
		history: []clickItem{},
	}
}

func (c *Clicker) Clicked(ctx *layout.Context) bool {
	c.processEvents(ctx)
	if c.clicks > 0 {
		c.clicks--
		if c.prevClicks > 0 {
			c.prevClicks--
		}
		if c.clicks > 0 {
			// Ensure timely delivery of remaining clicks.
			op.InvalidateOp{}.Add(ctx.Ops)
		}
		return true
	}
	return false
}

func (c *Clicker) History() []clickItem {
	return c.history
}

func (c *Clicker) processEvents(ctx *layout.Context) {
	for _, e := range c.click.Events(ctx) {
		switch e.Type {
		case gesture.TypeClick:
			c.clicks++
		case gesture.TypePress:
			c.history = append(c.history, clickItem{
				position: e.Position,
				time:     ctx.Now(),
			})
		}
	}
}

func (c *Clicker) Register(ctx *layout.Context) {
	// Flush clicks from before the previous frame.
	c.clicks -= c.prevClicks
	c.prevClicks = 0
	c.processEvents(ctx)
	c.click.Add(ctx.Ops)
	for len(c.history) > 0 {
		click := c.history[0]
		if ctx.Now().Sub(click.time) < 1*time.Second {
			break
		}
		copy(c.history, c.history[1:])
		c.history = c.history[:len(c.history)-1]
	}
}
