package values

import "image/color"

type Color struct {
	Primary          color.NRGBA
	Primary50        color.NRGBA
	PrimaryHighlight color.NRGBA

	// text colors
	Text      color.NRGBA // default color #091440
	InvText   color.NRGBA // inverted default color #ffffff
	GrayText1 color.NRGBA // darker shade #3D5873
	GrayText2 color.NRGBA // lighter shade of GrayText1 #596D81
	GrayText3 color.NRGBA // lighter shade of GrayText2 #8997A5 (hint)
	GrayText4 color.NRGBA // lighter shade of GrayText3 ##C4CBD2
	GreenText color.NRGBA // green text #41BE53

	// background colors
	Background       color.NRGBA
	Black            color.NRGBA
	BlueProgressTint color.NRGBA
	Danger           color.NRGBA
	DeepBlue         color.NRGBA
	LightBlue        color.NRGBA
	LightBlue2       color.NRGBA
	LightBlue3       color.NRGBA
	LightBlue4       color.NRGBA
	LightBlue5       color.NRGBA
	LightBlue6       color.NRGBA
	Gray1            color.NRGBA
	Gray2            color.NRGBA
	Gray3            color.NRGBA
	Gray4            color.NRGBA
	Green50          color.NRGBA
	Green500         color.NRGBA
	Orange           color.NRGBA
	Orange2          color.NRGBA
	Orange3          color.NRGBA
	OrangeRipple     color.NRGBA
	Success          color.NRGBA
	Success2         color.NRGBA
	Surface          color.NRGBA
	SurfaceHighlight color.NRGBA
	Turquoise100     color.NRGBA
	Turquoise300     color.NRGBA
	Turquoise700     color.NRGBA
	Turquoise800     color.NRGBA
	Yellow           color.NRGBA
}

func NewThemeColor(darkMode bool) *Color {
	c := &Color{}
	c.defualtThemeColors()
	if darkMode {
		c.darkThemeColors()
	}
	return c
}

func (c *Color) darkThemeColors() {
	c.Primary = rgb(0x57B6FF)

	// text colors
	c.Text = argb(0x99FFFFFF)
	c.GrayText1 = argb(0xDEFFFFFF)
	c.GrayText2 = argb(0x99FFFFFF)
	c.GrayText3 = argb(0x61FFFFFF)
	c.GrayText4 = argb(0x61FFFFFF)

	// background colors
	c.DeepBlue = argb(0x99FFFFFF)
	c.Gray1 = argb(0x99FFFFFF)
	c.Gray2 = rgb(0x3D3D3D)
	c.Gray3 = rgb(0x8997a5)
	c.Gray4 = argb(0x121212)
	c.Surface = rgb(0x252525)
}

func (c *Color) defualtThemeColors() *Color {
	c.Primary = rgb(0x2970ff)
	c.Primary50 = rgb(0xE3F2FF)
	c.PrimaryHighlight = rgb(0x1B41B3)

	// text colors
	c.Text = rgb(0x091440)
	c.InvText = rgb(0xffffff)
	c.GrayText1 = rgb(0x3d5873)
	c.GrayText2 = rgb(0x596D81)
	c.GrayText3 = rgb(0x8997a5) //hint
	c.GrayText4 = rgb(0xc4cbd2)
	c.GreenText = rgb(0x41BE53)

	// background colors
	c.Background = argb(0x22444444)
	c.Black = rgb(0x000000)
	c.BlueProgressTint = rgb(0x73d7ff)
	c.Danger = rgb(0xed6d47)
	c.DeepBlue = rgb(0x091440)
	c.LightBlue = rgb(0xe4f6ff)
	c.LightBlue2 = rgb(0x75D8FF)
	c.LightBlue3 = rgb(0xBCE8FF)
	c.LightBlue4 = rgb(0xBBDEFF)
	c.LightBlue5 = rgb(0x70CBFF)
	c.LightBlue6 = rgb(0x4B91D8)
	c.Gray1 = rgb(0x3d5873) // darkest gray #3D5873 (icon color)
	c.Gray2 = rgb(0xe6eaed) // light 0xe6eaed
	c.Gray3 = rgb(0xc4cbd2) // InactiveGray #C4CBD2
	c.Gray4 = rgb(0xf3f5f6) //active n light gray combined f3f5f6
	c.Green50 = rgb(0xE8F7EA)
	c.Green500 = rgb(0x41BE53)
	c.Orange = rgb(0xD34A21)
	c.Orange2 = rgb(0xF8E8E7)
	c.Orange3 = rgb(0xF8CABC)
	c.OrangeRipple = rgb(0xD32F2F)
	c.Success = rgb(0x41bf53)
	c.Success2 = rgb(0xE1F8EF)
	c.Surface = rgb(0xffffff)
	c.Turquoise100 = rgb(0xB6EED7)
	c.Turquoise300 = rgb(0x2DD8A3)
	c.Turquoise700 = rgb(0x00A05F)
	c.Turquoise800 = rgb(0x008F52)
	c.Yellow = rgb(0xffc84e)

	return c
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
