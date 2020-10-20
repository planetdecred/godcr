package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Combo struct {
	items         []ComboItem
	isOpen        bool
	selectedIndex int
	color         color.RGBA
}

type ComboItem struct {
	Text   string
	Icon   image.Image
	button Button
	label  Label
}

func (t *Theme) Combo(items []ComboItem) *Combo {
	c := &Combo{
		isOpen: false,
		items:  make([]ComboItem, len(items)+1),
		color:  mulAlpha(t.Color.Gray, 50),
	}

	for i := range items {
		items[i].button = t.Button(new(widget.Clickable), items[i].Text)
		items[i].label = t.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}
	c.items[0] = ComboItem{
		Text:   items[0].Text,
		label:  t.Body1(items[0].Text),
		button: t.Button(new(widget.Clickable), items[0].Text),
	}

	return c
}

func (c *Combo) Selected() string {
	return c.items[c.selectedIndex].Text
}

func (c *Combo) SelectedIndex() int {
	return c.selectedIndex
}

func (c *Combo) handleEvents() {
	for c.items[0].button.Button.Clicked() {
		c.isOpen = !c.isOpen
	}

	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Button.Clicked() {
				c.selectedIndex = index
				c.items[0].label.Text = c.items[index].Text
				c.isOpen = false
			}
		}
	}
}

func (c *Combo) layoutIcon(gtx layout.Context, itemIndex int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		if c.items[itemIndex].Icon == nil {
			return layout.Dimensions{}
		}

		img := widget.Image{Src: paint.NewImageOp(c.items[itemIndex].Icon)}
		img.Scale = 0.085

		return img.Layout(gtx)
	})
}

func (c *Combo) layoutText(gtx layout.Context, index int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return c.items[index].label.Layout(gtx)
	})
}

func (c *Combo) layoutActiveIcon(gtx layout.Context) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Dimensions{}
	})
}

func (c *Combo) layoutOption(gtx layout.Context, itemIndex int) layout.Dimensions {
	btn := c.items[itemIndex].button

	min := gtx.Constraints.Min
	min.X = 100
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
			}.Add(gtx.Ops)
			return fill(gtx, c.color)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = min
				iconLayout := c.layoutIcon(gtx, itemIndex)
				textLayout := c.layoutText(gtx, itemIndex)
				activeIconLayout := c.layoutActiveIcon(gtx)

				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, iconLayout, textLayout, activeIconLayout)
			})
		}),
		layout.Expanded(btn.Button.Layout),
	)
}

func (c *Combo) Layout(gtx layout.Context) layout.Dimensions {
	c.handleEvents()

	children := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return c.layoutOption(gtx, 0)
		}),
	}

	if c.isOpen {
		items := c.items[1:]
		for i := range items {
			index := i
			children = append(children, layout.Rigid(func(gtx C) D {
				return c.layoutOption(gtx, index+1)
			}))
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
