// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/draw"
)

var (
	// decred primary colors

	keyblue  = rgb(0x2970ff)
	darkblue = rgb(0x091440)

	// decred complemetary colors
	green = rgb(0x41bf53)
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Theme struct {
	Shaper text.Shaper
	Base   *material.Theme
	Color  struct {
		Primary      color.NRGBA
		Secondary    color.NRGBA
		Text         color.NRGBA
		Hint         color.NRGBA
		Overlay      color.NRGBA
		InvText      color.NRGBA
		Success      color.NRGBA
		Success2     color.NRGBA
		Danger       color.NRGBA
		Background   color.NRGBA
		Surface      color.NRGBA
		Gray         color.NRGBA
		Black        color.NRGBA
		DeepBlue     color.NRGBA
		LightBlue    color.NRGBA
		LightGray    color.NRGBA
		InactiveGray color.NRGBA
		Gray1        color.NRGBA
		Gray2        color.NRGBA
		Gray3        color.NRGBA
	}
	Icon struct {
		ContentCreate *widget.Icon
		ContentAdd    *widget.Icon
	}
	TextSize              unit.Value
	checkBoxCheckedIcon   *widget.Icon
	checkBoxUncheckedIcon *widget.Icon
	radioCheckedIcon      *widget.Icon
	radioUncheckedIcon    *widget.Icon
	chevronUpIcon         *widget.Icon
	dropDownIcon          *widget.Icon
	chevronDownIcon       *widget.Icon
	navigationCheckIcon   *widget.Icon
	navMoreIcon           *widget.Icon
	expandIcon            *widget.Image
	collapseIcon          *widget.Image

	dropDownMenus []*DropDown
}

func NewTheme(fontCollection []text.FontFace, decredIcons map[string]image.Image) *Theme {
	t := &Theme{
		Shaper: text.NewCache(fontCollection),
		Base:   material.NewTheme(fontCollection),
	}
	t.Color.Primary = keyblue
	t.Color.Text = darkblue
	t.Color.Hint = rgb(0x8997A5)
	t.Color.InvText = rgb(0xffffff)
	t.Color.Overlay = rgb(0x000000)
	t.Color.Surface = rgb(0xffffff)
	t.Color.Success = green
	t.Color.Success2 = rgb(0xE1F8EF)
	t.Color.Danger = rgb(0xed6d47)
	t.Color.Gray = rgb(0x596D81)
	t.Color.Gray1 = rgb(0xe6eaed)
	t.Color.Gray2 = rgb(0x8997a5)
	t.Color.Gray3 = rgb(0x3d5873)
	t.Color.LightGray = rgb(0xf3f5f6)
	t.Color.DeepBlue = rgb(0x091440)
	t.Color.InactiveGray = rgb(0xc4cbd2)
	t.Color.Black = rgb(0x000000)
	t.Color.Background = argb(0x22444444)
	t.Color.LightBlue = rgb(0xe4f6ff)
	t.TextSize = unit.Sp(16)

	t.checkBoxCheckedIcon = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))
	t.chevronUpIcon = mustIcon(widget.NewIcon(icons.NavigationExpandLess))
	t.chevronDownIcon = mustIcon(widget.NewIcon(icons.NavigationExpandMore))
	t.navMoreIcon = mustIcon(widget.NewIcon(icons.NavigationMoreHoriz))
	t.navigationCheckIcon = mustIcon(widget.NewIcon(icons.NavigationCheck))
	t.dropDownIcon = mustIcon(widget.NewIcon(icons.NavigationArrowDropDown))

	t.expandIcon = &widget.Image{Src: paint.NewImageOp(decredIcons["expand_icon"])}
	t.collapseIcon = &widget.Image{Src: paint.NewImageOp(decredIcons["collapse_icon"])}
	return t
}

func (t *Theme) Background(gtx layout.Context, w layout.Widget) {
	layout.Stack{
		Alignment: layout.N,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, t.Color.LightGray)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) Surface(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Stack{
		Alignment: layout.Center,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, t.Color.Surface)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) ImageIcon(gtx layout.Context, icon image.Image, size int) layout.Dimensions {
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	draw.ApproxBiLinear.Scale(img, img.Bounds(), icon, icon.Bounds(), draw.Src, nil)

	i := widget.Image{Src: paint.NewImageOp(img)}
	i.Scale = float32(size) / float32(gtx.Px(unit.Dp(float32(size))))
	return i.Layout(gtx)
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func fillMax(gtx layout.Context, col color.NRGBA, radius CornerRadius) D {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
	st := op.Save(gtx.Ops)
	track := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}

	clip.RRect{
		Rect: track,
		NE:   radius.NE, NW: radius.NW, SE: radius.SE, SW: radius.SW,
	}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()

	return layout.Dimensions{Size: d}
}

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: d,
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()

	return layout.Dimensions{Size: d}
}

func Fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	return fill(gtx, col)
}

// mulAlpha scales all color components by alpha/255.
func mulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	a := uint16(alpha)
	return color.NRGBA{
		A: uint8(uint16(c.A) * a / 255),
		R: uint8(uint16(c.R) * a / 255),
		G: uint8(uint16(c.G) * a / 255),
		B: uint8(uint16(c.B) * a / 255),
	}
}

func (t *Theme) closeAllDropdownMenus(group uint) {
	for _, dropDown := range t.dropDownMenus {
		if dropDown.group == group {
			dropDown.isOpen = false
		}
	}
}
