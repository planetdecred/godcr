package ui

import (
	"errors"
	"sync"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/notification"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
// Window maintains an internal state of variables to determine what to display at
// any point in time.
type Window struct {
	*app.Window

	wallet               *wallet.Wallet
	walletUnspentOutputs *wallet.UnspentOutputs

	load *load.Load

	modalMutex sync.Mutex
	modals     []load.Modal

	currentPage   load.Page
	pageBackStack []load.Page

	selectedAccount int
	txAuthor        dcrlibwallet.TxAuthor

	walletAcctMixerStatus chan *wallet.AccountMixer
}

type (
	C = layout.Context
	D = layout.Dimensions
)
type WriteClipboard struct {
	Text string
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(wal *wallet.Wallet) (*Window, error) {
	var netType string
	if wal.Net == dcrlibwallet.Testnet3 {
		netType = "testnet"
	} else {
		netType = wal.Net
	}

	win := &Window{
		Window:                app.NewWindow(app.MinSize(values.AppWidth, values.AppHeight), app.Title(values.StringF(values.StrAppTitle, netType))),
		wallet:                wal,
		walletUnspentOutputs:  new(wallet.UnspentOutputs),
		walletAcctMixerStatus: make(chan *wallet.AccountMixer),
	}

	l, err := win.NewLoad()
	if err != nil {
		return nil, err
	}
	win.load = l

	return win, nil
}

func (win *Window) NewLoad() (*load.Load, error) {
	th := decredmaterial.NewTheme(assets.FontCollection(), assets.DecredIcons, false)
	if th == nil {
		return nil, errors.New("unexpected error while loading theme")
	}

	l := &load.Load{
		Theme: th,

		WL: &load.WalletLoad{
			MultiWallet:    win.wallet.GetMultiWallet(),
			Wallet:         win.wallet,
			UnspentOutputs: win.walletUnspentOutputs,
			TxAuthor:       win.txAuthor,
		},

		Toast: notification.NewToast(th),

		Printer: message.NewPrinter(language.English),
	}

	l.RefreshWindow = win.Invalidate
	l.ShowModal = win.showModal
	l.DismissModal = win.dismissModal
	l.PopWindowPage = win.popPage
	l.ChangeWindowPage = win.changePage

	// ReloadApp closes the current page active on the
	// app window. When the next FrameEvent is received,
	// a new StartPage will be initialized and displayed.
	l.ReloadApp = func() {
		l.UIMu.Lock()
		defer l.UIMu.Unlock()
		if win.currentPage != nil {
			win.currentPage.OnNavigatedFrom()
			win.currentPage = nil
		}
	}

	// DarkModeSettingChanged checks if any page or any
	// modal implements the AppSettingsChangeHandler
	l.DarkModeSettingChanged = func(isDarkModeOn bool) {
		if page, ok := win.currentPage.(load.AppSettingsChangeHandler); ok {
			page.OnDarkModeChanged(isDarkModeOn)
		}
		for _, modal := range win.modals {
			if modal, ok := modal.(load.AppSettingsChangeHandler); ok {
				modal.OnDarkModeChanged(isDarkModeOn)
			}
		}
	}

	l.LanguageSettingChanged = func() {
		if page, ok := win.currentPage.(load.AppSettingsChangeHandler); ok {
			page.OnLanguageChanged()
		}
	}

	l.CurrencySettingChanged = func() {
		if page, ok := win.currentPage.(load.AppSettingsChangeHandler); ok {
			page.OnCurrencyChanged()
		}
	}

	return l, nil
}

// HandleEvents runs main event handling and page rendering loop.
func (win *Window) HandleEvents() {

	for {
		e := <-win.Events()
		switch evt := e.(type) {

		case system.DestroyEvent:
			if win.currentPage != nil {
				win.currentPage.OnNavigatedFrom()
				win.currentPage = nil
			}
			return // exits the loop, caller will exit the program.

		case system.FrameEvent:
			win.displayWindow(evt)

		case key.Event:
			win.handleKeyEvent(&evt)

		default:
			log.Tracef("Unhandled window event %v\n", e)
		}
	}
}

// displayWindow is called when a FrameEvent is received by the active window.
// Since user actions such as button clicks also trigger FrameEvents, this
// method first checks for pending user actions before displaying the UI
// elements. This ensures that the proper interface is displayed to the user
// based on their last performed action where applicable.
func (win *Window) displayWindow(evt system.FrameEvent) {
	win.load.UIMu.Lock()
	defer win.load.UIMu.Unlock()
	// Set up the StartPage the first time a FrameEvent is received.
	if win.currentPage == nil {
		win.currentPage = page.NewStartPage(win.load)
		win.currentPage.OnNavigatedTo()
	}

	// A FrameEvent may be generated because of a user interaction
	// with the current page such as a button click. First handle
	// any such user interaction before rendering the page.
	win.currentPage.HandleUserInteractions()
	for _, modal := range win.modals {
		modal.Handle() // TODO: Just the top-most modal should do.
	}

	// Draw the window's UI components into an op.Ops.
	gtx := layout.NewContext(&op.Ops{}, evt)

	win.drawWindowUI(gtx)

	// Render the window's UI components on screen.
	evt.Frame(gtx.Ops)
}

// drawWindowUI draws the window UI components into the provided graphical
// context, preparing the context for rendering on screen.
func (win *Window) drawWindowUI(gtx C) {
	// Create a base view holder to hold all the following UI components
	// one on top the other. Components that do not take up the entire
	// window will be aligned to the top of the window.
	viewsHolder := layout.Stack{Alignment: layout.N}

	background := layout.Expanded(func(gtx C) D {
		return decredmaterial.Fill(gtx, win.load.Theme.Color.Gray4)
	})

	// TODO: Should suffice to just draw the top-most modal?
	modals := layout.Stacked(func(gtx C) D {
		modals := win.modals
		if len(modals) == 0 {
			return layout.Dimensions{}
		}

		modalLayouts := make([]layout.StackChild, 0)
		for _, modal := range modals {
			widget := modal.Layout(gtx)
			l := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return widget
			})
			modalLayouts = append(modalLayouts, l)
		}

		return layout.Stack{Alignment: layout.Center}.Layout(gtx, modalLayouts...)
	})

	viewsHolder.Layout(
		gtx,
		background,
		layout.Stacked(win.currentPage.Layout),
		modals,
		layout.Stacked(win.load.Toast.Layout),
	)
}

func (win *Window) handleKeyEvent(evt *key.Event) {
	// Handle key events on any displayed modals first.
	for _, modal := range win.modals {
		if handler, ok := modal.(load.KeyEventHandler); ok {
			handler.HandleKeyEvent(evt)
		}
	}
	// Only handle key events on the current page if no modal is displayed.
	if len(win.modals) == 0 {
		if handler, ok := win.currentPage.(load.KeyEventHandler); ok {
			handler.HandleKeyEvent(evt)
		}
	}
}

// changePage displays the provided page on the window and optionally adds
// the current page to the backstack. This automatically refreshes the display,
// callers should not re-refresh the display.
func (win *Window) changePage(page load.Page, keepBackStack bool) {
	if win.currentPage != nil && keepBackStack {
		win.currentPage.OnNavigatedFrom()
		win.pageBackStack = append(win.pageBackStack, win.currentPage)
	}

	win.currentPage = page
	win.currentPage.OnNavigatedTo()
	win.Invalidate()
}

// popPage goes back to the previous page. This automatically refreshes the
// display, callers should not re-refresh the display.
// Returns true if page was popped.
func (win *Window) popPage() bool {
	if len(win.pageBackStack) == 0 {
		return false
	}

	// get and remove last page
	previousPageIndex := len(win.pageBackStack) - 1
	previousPage := win.pageBackStack[previousPageIndex]
	win.pageBackStack = win.pageBackStack[:previousPageIndex]

	// close the current page and display the previous page
	win.currentPage.OnNavigatedFrom()
	previousPage.OnNavigatedTo()
	win.currentPage = previousPage
	win.Invalidate()

	return true
}

// TODO: showModal should refresh display, callers shouldn't.
func (win *Window) showModal(modal load.Modal) {
	modal.OnResume() // setup display data
	win.modalMutex.Lock()
	win.modals = append(win.modals, modal)
	win.modalMutex.Unlock()
}

func (win *Window) dismissModal(modal load.Modal) {
	win.modalMutex.Lock()
	defer win.modalMutex.Unlock()
	for i, m := range win.modals {
		if m.ModalID() == modal.ModalID() {
			modal.OnDismiss() // do garbage collection in modal
			win.modals = append(win.modals[:i], win.modals[i+1:]...)
			win.Invalidate()
			return
		}
	}
}
