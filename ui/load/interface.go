package load

import (
	"gioui.org/io/key"
)

// AppSettingsChangeHandler defines a method that can be implemented by pages and
// modals to watch for real-time changes to the dark mode setting and modify
// widget appearance accordingly.
type AppSettingsChangeHandler interface {
	// OnDarkModeChanged is triggered whenever the dark mode setting is changed
	// to enable restyling UI elements where necessary.
	OnDarkModeChanged(bool)
	// OnCurrencyChanged is triggered whenever the currency setting is changed
	// to enable app refresh where necessary especially on the main page.
	OnCurrencyChanged()
	// OnLanguageChanged is triggered whenever the language setting is changed
	// to enable UI language update where necessary especially on page Nav
	OnLanguageChanged()
}

// KeyEventHandler is implemented by pages and modals that require key event
// notifications.
type KeyEventHandler interface {
	// KeysToHandle returns an expression that describes a set of key
	// combinations that the implementer of this interface wishes to capture.
	// The HandleKeyPress() method will only be called when any of these key
	// combinations is pressed.
	KeysToHandle() key.Set
	// HandleKeyPress is called when one or more keys are pressed on the current
	// window that match any of the key combinations returned by KeysToHandle().
	HandleKeyPress(*key.Event)
}
