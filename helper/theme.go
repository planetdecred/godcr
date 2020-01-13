package helper

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"gioui.org/font"
	"gioui.org/font/gofont"
)

type (
	Fonts struct {
		Regular       text.Font
		Bold          text.Font
		RegularItalic text.Font
		BoldItalic    text.Font
	}
	Theme struct {
		*material.Theme
		*text.Shaper
		Fonts *Fonts
	}
)

const (
	regularFontSize = 12
)

var (
	theme *Theme
)

func Initialize() {
	if theme == nil {
		theme = newTheme()
	}
}

func newTheme() *Theme {
	gofont.Register()

	materialTheme := &material.Theme{
		Shaper: font.Default(),
	}

	materialTheme.Color.Primary = DecredDarkBlueColor
	materialTheme.Color.Text = BlackColor
	materialTheme.Color.Hint = GrayColor
	materialTheme.TextSize = unit.Dp(10)

	return &Theme{
		Theme:  materialTheme,
		Shaper: materialTheme.Shaper,
		Fonts:  getFonts(),
	}
}

func getFonts() *Fonts {
	return &Fonts{
		Regular: text.Font{
			Size: unit.Dp(regularFontSize),
		},
		Bold: text.Font{
			Size:   unit.Dp(regularFontSize),
			Weight: text.Bold,
		},
		RegularItalic: text.Font{
			Size:  unit.Dp(regularFontSize),
			Style: text.Italic,
		},
		BoldItalic: text.Font{
			Size:   unit.Dp(regularFontSize),
			Weight: text.Bold,
			Style:  text.Italic,
		},
	}
}

func GetTheme() *Theme {
	return theme
}
