package app

import (
	"gioui.org/layout"
)

// Page defines methods that control the appearance and functionality of
// UI components displayed on a window.
type Page interface {
	// ID is a unique string that identifies the page and may be used
	// to differentiate this page from other pages.
	ID() string
	// OnAttachedToNavigator is called when navigation occurs; i.e. when a page
	// or modal is pushed into the window's display. The navigator parameter is
	// the PageNavigator or WindowNavigator object that is used to display the
	// content. This is called just before OnNavigatedTo() is called.
	OnAttachedToNavigator(navigator PageNavigator)
	// OnNavigatedTo is called when the page is about to be displayed and may be
	// used to initialize page features that are only relevant when the page is
	// displayed. This is called just before HandleUserInteractions() and
	// Layout() are called (in that order).
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
	// ID is a unique string that identifies the modal and may be used
	// to differentiate this modal from other modals.
	ID() string
	// OnAttachedToNavigator is called when navigation occurs; i.e. when a page
	// or modal is pushed into the window's display. The navigator parameter is
	// the PageNavigator or WindowNavigator object that is used to display the
	// content. This is called just before OnResume() is called.
	OnAttachedToNavigator(navigator PageNavigator)
	// OnResume is called to initialize data and get UI elements ready to be
	// displayed. This is called just before Handle() and Layout() are called (in
	// that order).
	OnResume()
	// Handle is called just before Layout() to determine if any user
	// interaction recently occurred on the modal and may be used to update the
	// page's UI components shortly before they are displayed.
	Handle()
	// Layout draws the modal's UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(gtx layout.Context) layout.Dimensions
	// OnDismiss is called after the modal is dismissed.
	// NOTE: The modal may be re-displayed on the app's window, in which case
	// OnResume() will be called again. This method should not destroy UI
	// components unless they'll be recreated in the OnResume() method.
	OnDismiss()
}

// Closable should be implemented by pages and modals that want to know when
// they are closed in order to perform some cleanup actions.
type Closable interface {
	// OnClosed is called to indicate that a specific instance of a page or
	// modal has been dismissed and will no longer be displayed.
	OnClosed()
}

// PageNavigator defines methods for navigating between pages in a window or a
// MasterPage.
type PageNavigator interface {
	// CurrentPage returns the page that is at the top of the stack. Returns nil
	// if the stack is empty.
	CurrentPage() Page
	// CurrentPageID returns the ID of the current page or an empty string if no
	// page is displayed.
	CurrentPageID() string
	// Display causes the specified page to be displayed on the parent window or
	// page. All other instances of this same page will be closed and removed
	// from the backstack.
	Display(page Page)
	// CloseCurrentPage dismisses the page at the top of the stack and gets the
	// next page ready for display.
	CloseCurrentPage()
	// ClosePagesAfter dismisses all pages from the top of the stack until (and
	// excluding) the page with the specified ID. If no page is found with the
	// provided ID, no page will be popped. The page with the specified ID will
	// be displayed after the other pages are popped.
	ClosePagesAfter(keepPageID string)
	// ClearStackAndDisplay dismisses all pages in the stack and displays the
	// specified page.
	ClearStackAndDisplay(page Page)
	// CloseAllPages dismisses all pages in the stack.
	CloseAllPages()
}

// WindowNavigator defines methods for page navigation, displaying modals and
// reloading the entire window display.
type WindowNavigator interface {
	PageNavigator
	// ShowModal displays a modal over the current page. Any previously
	// displayed modal will be hidden by this new modal.
	ShowModal(Modal)
	// DismissModal dismisses the modal with the specified ID, if it was
	// previously displayed by this WindowNavigator. If there are more than 1
	// modal with the specified ID, only the top-most instance is dismissed.
	DismissModal(modalID string)
	// TopModal returns the top-most modal in display or nil if there is no
	// modal in display.
	TopModal() Modal
	// Reload causes the entire window display to be reloaded. If a page is
	// currently displayed, this should call the page's HandleUserInteractions()
	// method. If a modal is displayed, the modal's Handle() method should also
	// be called. Finally, the current page and modal's Layout methods should be
	// called to render the entire window's display.
	Reload()
}
