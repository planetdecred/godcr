package load

import (
	"gioui.org/io/key"
	"gioui.org/layout"
)

// Page defines methods that control the appearance and functionality of
// UI components displayed on a window.
type Page interface {
	// ID is a unique string that identifies the page and may be used
	// to differentiate this page from other pages.
	ID() string
	// OnNavigatedTo is called when the page is about to be displayed and
	// may be used to initialize page features that are only relevant when
	// the page is displayed.
	OnNavigatedTo()
	// HandleUserInteractions is called just before Layout() to determine
	// if any user interaction recently occurred on the page and may be
	// used to update the page's UI components shortly before they are
	// displayed.
	HandleUserInteractions()
	// Layout draws the page UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(layout.Context) layout.Dimensions
	// OnNavigatedFrom is called when the page is about to be removed from
	// the displayed window. This method should ideally be used to disable
	// features that are irrelevant when the page is NOT displayed.
	// NOTE: The page may be re-displayed on the app's window, in which case
	// OnNavigatedTo() will be called again. This method should not destroy UI
	// components unless they'll be recreated in the OnNavigatedTo() method.
	OnNavigatedFrom()
}

type Modal interface {
	ModalID() string
	OnResume()
	Layout(gtx layout.Context) layout.Dimensions
	OnDismiss()
	Dismiss()
	Show()
	Handle()
}

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
	// HandleKeyEvent is called when a key is pressed on the current window.
	HandleKeyEvent(*key.Event)
}
