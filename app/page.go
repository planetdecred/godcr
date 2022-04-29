package app

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type PageMaker func(mw *dcrlibwallet.MultiWallet, theme *decredmaterial.Theme, navigator Navigator) (Page, error)

type IPage interface {
	// Construct / New / Instance / Instantiate
	Construct(mw *dcrlibwallet.MultiWallet, theme *decredmaterial.Theme, navigator Navigator) error

	Page
}

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

type KeyEventHandler interface {
	// HandleKeyEvent is called when a key is pressed on the current window.
	HandleKeyEvent(*key.Event)
}
