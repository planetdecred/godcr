package decredmaterial

import (
	"gioui.org/widget"
)

// DangerButton a button with the background set to theme.Danger
func (t *Theme) DangerButton(button *widget.Clickable, text string) Button {
	btn := t.Button(button, text)
	btn.Background = t.Color.Danger
	return btn
}
