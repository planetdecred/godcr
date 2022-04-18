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

	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/draw"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Theme struct {
	Shaper text.Shaper
	Base   *material.Theme
	Color  *values.Color
	Styles *values.WidgetStyles

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
}

func NewTheme(fontCollection []text.FontFace, decredIcons map[string]image.Image, isDarkModeOn bool) *Theme {
	t := &Theme{
		Shaper:   text.NewCache(fontCollection),
		Base:     material.NewTheme(fontCollection),
		Color:    &values.Color{},
		Styles:   values.DefaultWidgetStyles(),
		TextSize: values.TextSize16,
	}
	t.SwitchDarkMode(isDarkModeOn, decredIcons)
	t.checkBoxCheckedIcon = MustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = MustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = MustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = MustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))
	t.chevronUpIcon = MustIcon(widget.NewIcon(icons.NavigationExpandLess))
	t.chevronDownIcon = MustIcon(widget.NewIcon(icons.NavigationExpandMore))
	t.navMoreIcon = MustIcon(widget.NewIcon(icons.NavigationMoreHoriz))
	t.navigationCheckIcon = MustIcon(widget.NewIcon(icons.NavigationCheck))
	t.dropDownIcon = MustIcon(widget.NewIcon(icons.NavigationArrowDropDown))

	return t
}

func (t *Theme) SwitchDarkMode(isDarkModeOn bool, decredIcons map[string]image.Image) {
	t.Color.DefualtThemeColors()
	expandIcon := "expand_icon"
	collapseIcon := "collapse_icon"
	if isDarkModeOn {
		t.Color.DarkThemeColors() // override defaults with dark themed colors
		expandIcon = "expand_dm"
		collapseIcon = "collapse_dm"
	}

	t.expandIcon = NewImage(decredIcons[expandIcon])
	t.collapseIcon = NewImage(decredIcons[collapseIcon])

	t.updateStyles(isDarkModeOn)
}

// UpdateStyles update the style definition for different widgets. This should
// be done whenever the base theme changes to ensure that the style definitions
// use the values for the latest theme.
func (t *Theme) updateStyles(isDarkModeOn bool) {
	// update switch style colors
	t.Styles.SwitchStyle.ActiveColor = t.Color.Primary
	t.Styles.SwitchStyle.InactiveColor = t.Color.Gray3
	t.Styles.SwitchStyle.ThumbColor = t.Color.White

	// update icon button style colors
	t.Styles.IconButtonColorStyle.Background = color.NRGBA{}
	t.Styles.IconButtonColorStyle.Foreground = t.Color.Gray1

	// update Collapsible widget style colors
	t.Styles.CollapsibleStyle.Background = t.Color.Surface
	t.Styles.CollapsibleStyle.Foreground = color.NRGBA{}

	// update clickable colors
	t.Styles.ClickableStyle.Color = t.Color.SurfaceHighlight
	t.Styles.ClickableStyle.HoverColor = t.Color.Gray5

	// dropdown clickable colors
	t.Styles.DropdownClickableStyle.Color = t.Color.SurfaceHighlight
	col := t.Color.Gray3
	if isDarkModeOn {
		col = t.Color.Gray5
	}
	t.Styles.DropdownClickableStyle.HoverColor = Hovered(col)
}

func (t *Theme) Background(gtx layout.Context, w layout.Widget) {
	layout.Stack{
		Alignment: layout.N,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, t.Color.Gray4)
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

// isOpenDropdownGroup iterate over Dropdowns registered as a member
// of {group}, returns true if any of the drop down state is open.
func (t *Theme) isOpenDropdownGroup(group uint) bool {
	for _, dropDown := range t.dropDownMenus {
		if dropDown.group == group {
			if dropDown.isOpen {
				return true
			}
		}
	}
	return false
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
