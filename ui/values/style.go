package values

import (
	"image/color"
)

// SwitchStyle defines display properties that may be used to style a
// Switch widget.
type SwitchStyle struct {
	ActiveColor       color.NRGBA
	InactiveColor     color.NRGBA
	ThumbColor        color.NRGBA
	ActiveTextColor   color.NRGBA
	InactiveTextColor color.NRGBA
}

// ColorStyle defines backgorund and foreground colors that may be used to
// style a widget that requires either or both colors.
type ColorStyle struct {
	Background color.NRGBA
	Foreground color.NRGBA
}

// ClickableStyle defines display properties that may be used to style a
// Clickable widget.
type ClickableStyle struct {
	Color      color.NRGBA
	HoverColor color.NRGBA
}

// WidgetStyles is a collection of various widget styles.
type WidgetStyles struct {
	SwitchStyle            *SwitchStyle
	IconButtonColorStyle   *ColorStyle
	CollapsibleStyle       *ColorStyle
	ClickableStyle         *ClickableStyle
	DropdownClickableStyle *ClickableStyle
}

// DefaultWidgetStyles returns a new collection of widget styles with default
// values.
func DefaultWidgetStyles() *WidgetStyles {
	return &WidgetStyles{
		SwitchStyle:            &SwitchStyle{},
		IconButtonColorStyle:   &ColorStyle{},
		CollapsibleStyle:       &ColorStyle{},
		ClickableStyle:         &ClickableStyle{},
		DropdownClickableStyle: &ClickableStyle{},
	}
}
