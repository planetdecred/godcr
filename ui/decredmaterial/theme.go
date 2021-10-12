// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
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
		Primary          color.NRGBA
		Primary50        color.NRGBA
		PrimaryHighlight color.NRGBA
		Secondary        color.NRGBA
		Text             color.NRGBA
		Hint             color.NRGBA
		Overlay          color.NRGBA
		InvText          color.NRGBA
		Success          color.NRGBA
		Success2         color.NRGBA
		Danger           color.NRGBA
		Background       color.NRGBA
		Surface          color.NRGBA
		SurfaceHighlight color.NRGBA
		Gray             color.NRGBA
		Black            color.NRGBA
		DeepBlue         color.NRGBA
		LightBlue        color.NRGBA
		LightBlue2       color.NRGBA
		LightBlue3       color.NRGBA
		LightBlue4       color.NRGBA
		LightBlue5       color.NRGBA
		LightBlue6       color.NRGBA
		BlueProgressTint color.NRGBA
		LightGray        color.NRGBA
		InactiveGray     color.NRGBA
		ActiveGray       color.NRGBA
		Gray1            color.NRGBA
		Gray2            color.NRGBA
		Gray3            color.NRGBA
		Orange           color.NRGBA
		Orange2          color.NRGBA
		Orange3          color.NRGBA
		OrangeRipple     color.NRGBA
		Gray4            color.NRGBA
		Gray5            color.NRGBA
		Gray6            color.NRGBA
		Green50          color.NRGBA
		Green500         color.NRGBA
		Turquoise100     color.NRGBA
		Turquoise300     color.NRGBA
		Turquoise700     color.NRGBA
		Turquoise800     color.NRGBA
		Yellow           color.NRGBA
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
	expandIcon            *Image
	collapseIcon          *Image

	dropDownMenus []*DropDown

	DarkMode bool
}

func (t *Theme) setColorMode(darkMode bool) {
	if darkMode {
		t.DarkMode = true
		t.Color.Primary = rgb(0x57B6FF)
		t.Color.Primary50 = rgb(0xE3F2FF)
		t.Color.PrimaryHighlight = rgb(0x1B41B3)
		t.Color.Text = argb(0x99FFFFFF)
		t.Color.Hint = rgb(0x8997A5)
		t.Color.InvText = rgb(0xffffff)
		t.Color.Overlay = rgb(0x000000)
		t.Color.Surface = rgb(0x252525)
		t.Color.SurfaceHighlight = rgb(0x3D3D3D)
		t.Color.Success = green
		t.Color.Success2 = rgb(0xE1F8EF)
		t.Color.Danger = rgb(0xed6d47)
		t.Color.Gray = argb(0x99FFFFFF)
		t.Color.Gray1 = rgb(0x1E1E1E)
		t.Color.Gray2 = rgb(0x8997a5)
		t.Color.Gray3 = argb(0xDEFFFFFF)
		t.Color.Gray4 = argb(0x99FFFFFF)
		t.Color.Gray5 = argb(0x61FFFFFF)
		t.Color.Gray6 = argb(0xCCFFFFFF)
		t.Color.Green50 = rgb(0xE8F7EA)
		t.Color.Green500 = rgb(0x41BE53)
		t.Color.LightGray = rgb(0x121212)
		t.Color.ActiveGray = rgb(0x363636)
		t.Color.DeepBlue = argb(0xDEFFFFFF)
		t.Color.InactiveGray = rgb(0xc4cbd2)
		t.Color.Black = rgb(0x000000)
		t.Color.Background = argb(0x22444444)
		t.Color.LightBlue = rgb(0xe4f6ff)
		t.Color.LightBlue2 = rgb(0x75D8FF)
		t.Color.LightBlue3 = rgb(0xBCE8FF)
		t.Color.LightBlue4 = rgb(0xBBDEFF)
		t.Color.LightBlue5 = rgb(0x70CBFF)
		t.Color.LightBlue6 = rgb(0x4B91D8)
		t.Color.BlueProgressTint = rgb(0x73d7ff)
		t.Color.Orange = rgb(0xD34A21)
		t.Color.Orange2 = rgb(0xF8E8E7)
		t.Color.Orange3 = rgb(0xF8CABC)
		t.Color.OrangeRipple = rgb(0xD32F2F)
		t.Color.Turquoise100 = rgb(0xB6EED7)
		t.Color.Turquoise300 = rgb(0x2DD8A3)
		t.Color.Turquoise700 = rgb(0x00A05F)
		t.Color.Turquoise800 = rgb(0x008F52)
		t.Color.Yellow = rgb(0xffc84e)
		t.TextSize = unit.Sp(16)
	} else {
		t.DarkMode = false
		t.Color.Primary = keyblue
		t.Color.Primary50 = rgb(0xE3F2FF)
		t.Color.PrimaryHighlight = rgb(0x1B41B3)
		t.Color.Text = darkblue
		t.Color.Hint = rgb(0x8997A5)
		t.Color.InvText = rgb(0xffffff)
		t.Color.Overlay = rgb(0x000000)
		t.Color.Surface = rgb(0xffffff)
		t.Color.SurfaceHighlight = rgb(0xE6EAED)
		t.Color.Success = green
		t.Color.Success2 = rgb(0xE1F8EF)
		t.Color.Danger = rgb(0xed6d47)
		t.Color.Gray = rgb(0x596D81)
		t.Color.Gray1 = rgb(0xe6eaed)
		t.Color.Gray2 = rgb(0x8997a5)
		t.Color.Gray3 = rgb(0x3d5873)
		t.Color.Gray4 = rgb(0x3d5873)
		t.Color.Gray5 = rgb(0x3d5873)
		t.Color.Gray6 = rgb(0x091440)
		t.Color.Green50 = rgb(0xE8F7EA)
		t.Color.Green500 = rgb(0x41BE53)
		t.Color.LightGray = rgb(0xf3f5f6)
		t.Color.ActiveGray = rgb(0xf3f5f6)
		t.Color.DeepBlue = rgb(0x091440)
		t.Color.InactiveGray = rgb(0xc4cbd2)
		t.Color.Black = rgb(0x000000)
		t.Color.Background = argb(0x22444444)
		t.Color.LightBlue = rgb(0xe4f6ff)
		t.Color.LightBlue2 = rgb(0x75D8FF)
		t.Color.LightBlue3 = rgb(0xBCE8FF)
		t.Color.LightBlue4 = rgb(0xBBDEFF)
		t.Color.LightBlue5 = rgb(0x70CBFF)
		t.Color.LightBlue6 = rgb(0x4B91D8)
		t.Color.BlueProgressTint = rgb(0x73d7ff)
		t.Color.Orange = rgb(0xD34A21)
		t.Color.Orange2 = rgb(0xF8E8E7)
		t.Color.Orange3 = rgb(0xF8CABC)
		t.Color.OrangeRipple = rgb(0xD32F2F)
		t.Color.Turquoise100 = rgb(0xB6EED7)
		t.Color.Turquoise300 = rgb(0x2DD8A3)
		t.Color.Turquoise700 = rgb(0x00A05F)
		t.Color.Turquoise800 = rgb(0x008F52)
		t.Color.Yellow = rgb(0xffc84e)
		t.TextSize = unit.Sp(16)
	}
}

func NewTheme(fontCollection []text.FontFace, decredIcons map[string]image.Image, isDarkModeOn bool) *Theme {
	t := &Theme{
		Shaper:   text.NewCache(fontCollection),
		Base:     material.NewTheme(fontCollection),
		DarkMode: false,
	}

	t.setColorMode(isDarkModeOn)

	t.checkBoxCheckedIcon = MustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = MustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = MustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = MustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))
	t.chevronUpIcon = MustIcon(widget.NewIcon(icons.NavigationExpandLess))
	t.chevronDownIcon = MustIcon(widget.NewIcon(icons.NavigationExpandMore))
	t.navMoreIcon = MustIcon(widget.NewIcon(icons.NavigationMoreHoriz))
	t.navigationCheckIcon = MustIcon(widget.NewIcon(icons.NavigationCheck))
	t.dropDownIcon = MustIcon(widget.NewIcon(icons.NavigationArrowDropDown))

	t.expandIcon = NewImage(decredIcons["expand_icon"])
	t.collapseIcon = NewImage(decredIcons["collapse_icon"])
	return t
}

func (t *Theme) SwitchDarkMode(isDarkModeOn bool) {
	t.setColorMode(isDarkModeOn)
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

func MustIcon(ic *widget.Icon, err error) *widget.Icon {
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
	track := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}

	defer clip.RRect{
		Rect: track,
		NE:   radius.TopRight, NW: radius.TopLeft, SE: radius.BottomRight, SW: radius.BottomLeft,
	}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, col)

	return layout.Dimensions{Size: d}
}

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	track := image.Rectangle{
		Max: d,
	}
	defer clip.Rect(track).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, col)

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

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func Disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	return color.NRGBA{
		R: byte((int(c.R)*r + int(lum)*(256-r)) / 256),
		G: byte((int(c.G)*r + int(lum)*(256-r)) / 256),
		B: byte((int(c.B)*r + int(lum)*(256-r)) / 256),
		A: byte(int(c.A) * (128 + 32) / 256),
	}
}

// Hovered blends color towards a brighter color.
func Hovered(c color.NRGBA) (d color.NRGBA) {
	const r = 0x20 // lighten ratio
	return color.NRGBA{
		R: byte(255 - int(255-c.R)*(255-r)/256),
		G: byte(255 - int(255-c.G)*(255-r)/256),
		B: byte(255 - int(255-c.B)*(255-r)/256),
		A: c.A,
	}
}

// approxLuminance is a fast approximate version of RGBA.Luminance.
func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}

func HandleEditorEvents(editors ...*widget.Editor) (bool, bool) {
	var submit, changed bool
	for _, editor := range editors {
		for _, evt := range editor.Events() {
			switch evt.(type) {
			case widget.ChangeEvent:
				changed = true
			case widget.SubmitEvent:
				submit = true
			}
		}
	}
	return submit, changed
}

//Tab key event handler for pages withe ditors
func HandleTabEvent(event chan *key.Event) bool {
	var isTabPressed bool
	select {
	case event := <-event:
		if event.Name == key.NameTab && event.State == key.Press {
			isTabPressed = true
		}
	default:
	}
	return isTabPressed
}

//Switch editors when tab key is pressed
func SwitchEditors(keyEvent chan *key.Event, editors ...*widget.Editor) {
	for i := 0; i < len(editors); i++ {
		if editors[i].Focused() {
			if HandleTabEvent(keyEvent) {
				if i == len(editors)-1 {
					editors[0].Focus()
				} else {
					editors[i+1].Focus()
				}
			}
		}
	}
}
