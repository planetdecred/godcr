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
	NavyBlue         color.NRGBA
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
	Gray5            color.NRGBA
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
	White            color.NRGBA
}

func (c *Color) DarkThemeColors() {
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
	c.Gray4 = rgb(0x121212)
	c.Gray5 = rgb(0x363636)
	c.Surface = rgb(0x252525)
}

func (c *Color) DefaultThemeColors() *Color {
	cl := Color{
		Primary:          rgb(0x2970ff),
		Primary50:        rgb(0xE3F2FF),
		PrimaryHighlight: rgb(0x1B41B3),

		// text colors
		Text:      rgb(0x091440),
		InvText:   rgb(0xffffff),
		GrayText1: rgb(0x3d5873),
		GrayText2: rgb(0x596D81),
		GrayText3: rgb(0x8997a5), //hint
		GrayText4: rgb(0xc4cbd2),
		GreenText: rgb(0x41BE53),

		// background colors
		Background:       argb(0x22444444),
		Black:            rgb(0x000000),
		BlueProgressTint: rgb(0x73d7ff),
		Danger:           rgb(0xed6d47),
		DeepBlue:         rgb(0x091440),
		NavyBlue:         rgb(0x1F45B0),
		LightBlue:        rgb(0xe4f6ff),
		LightBlue2:       rgb(0x75D8FF),
		LightBlue3:       rgb(0xBCE8FF),
		LightBlue4:       rgb(0xBBDEFF),
		LightBlue5:       rgb(0x70CBFF),
		LightBlue6:       rgb(0x4B91D8),
		Gray1:            rgb(0x3d5873), // darkest gray #3D5873 (icon color),
		Gray2:            rgb(0xe6eaed), // light 0xe6eaed
		Gray3:            rgb(0xc4cbd2), // InactiveGray #C4CBD2
		Gray4:            rgb(0xf3f5f6), //active n light gray combined f3f5f6
		Gray5:            rgb(0xf3f5f6),
		Green50:          rgb(0xE8F7EA),
		Green500:         rgb(0x41BE53),
		Orange:           rgb(0xD34A21),
		Orange2:          rgb(0xF8E8E7),
		Orange3:          rgb(0xF8CABC),
		OrangeRipple:     rgb(0xD32F2F),
		Success:          rgb(0x41bf53),
		Success2:         rgb(0xE1F8EF),
		Surface:          rgb(0xffffff),
		Turquoise100:     rgb(0xB6EED7),
		Turquoise300:     rgb(0x2DD8A3),
		Turquoise700:     rgb(0x00A05F),
		Turquoise800:     rgb(0x008F52),
		Yellow:           rgb(0xffc84e),
		White:            rgb(0xffffff),
	}

	return &cl
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
