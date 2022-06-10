package ui

import (
	"errors"

	giouiApp "gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
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
	*giouiApp.Window
	navigator app.WindowNavigator

	wallet               *wallet.Wallet
	walletUnspentOutputs *wallet.UnspentOutputs

	load *load.Load

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

	giouiWindow := giouiApp.NewWindow(giouiApp.MinSize(values.AppWidth, values.AppHeight), giouiApp.Title(values.StringF(values.StrAppTitle, netType)))
	win := &Window{
		Window:                giouiWindow,
		navigator:             app.NewSimpleWindowNavigator(giouiWindow.Invalidate),
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

	mw := win.wallet.GetMultiWallet()

	// Set the user-configured theme colors on app load.
	isDarkModeOn := mw.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false)
	th.SwitchDarkMode(isDarkModeOn, assets.DecredIcons)

	l := &load.Load{
		Theme: th,

		WL: &load.WalletLoad{
			MultiWallet:    mw,
			Wallet:         win.wallet,
			UnspentOutputs: win.walletUnspentOutputs,
			TxAuthor:       win.txAuthor,
		},

		Toast: notification.NewToast(th),

		Printer: message.NewPrinter(language.English),
	}

	// DarkModeSettingChanged checks if any page or any
	// modal implements the AppSettingsChangeHandler
	l.DarkModeSettingChanged = func(isDarkModeOn bool) {
		if page, ok := win.navigator.CurrentPage().(load.AppSettingsChangeHandler); ok {
			page.OnDarkModeChanged(isDarkModeOn)
		}
		if modal := win.navigator.TopModal(); modal != nil {
			if modal, ok := modal.(load.AppSettingsChangeHandler); ok {
				modal.OnDarkModeChanged(isDarkModeOn)
			}
		}
	}

	l.LanguageSettingChanged = func() {
		if page, ok := win.navigator.CurrentPage().(load.AppSettingsChangeHandler); ok {
			page.OnLanguageChanged()
		}
	}

	l.CurrencySettingChanged = func() {
		if page, ok := win.navigator.CurrentPage().(load.AppSettingsChangeHandler); ok {
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
			win.navigator.CloseAllPages()
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
	// Set up the StartPage the first time a FrameEvent is received.
	if win.navigator.CurrentPage() == nil {
		win.navigator.Display(page.NewStartPage(win.load))
		return
	}

	// A FrameEvent may be generated because of a user interaction
	// with the current page such as a button click. First handle
	// any such user interaction before rendering the page.
	win.navigator.CurrentPage().HandleUserInteractions()
	if modal := win.navigator.TopModal(); modal != nil {
		modal.Handle()
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

	topModalLayout := layout.Stacked(func(gtx C) D {
		modal := win.navigator.TopModal()
		if modal == nil {
			return layout.Dimensions{}
		}
		return modal.Layout(gtx)
	})

	viewsHolder.Layout(
		gtx,
		background,
		layout.Stacked(win.navigator.CurrentPage().Layout),
		topModalLayout,
		layout.Stacked(win.load.Toast.Layout),
	)
}

func (win *Window) handleKeyEvent(evt *key.Event) {
	// Handle key events on the top modal if a modal is displayed.
	// Only handle key events on the current page if no modal is displayed.
	if modal := win.navigator.TopModal(); modal != nil {
		if handler, ok := modal.(load.KeyEventHandler); ok {
			handler.HandleKeyEvent(evt)
		}
	} else {
		if handler, ok := win.navigator.CurrentPage().(load.KeyEventHandler); ok {
			handler.HandleKeyEvent(evt)
		}
	}
}
