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

	txAuthor dcrlibwallet.TxAuthor

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
			ops := win.handleFrameEvent(evt)
			evt.Frame(ops)

		default:
			log.Tracef("Unhandled window event %v\n", e)
		}
	}
}

// handleFrameEvent is called when a FrameEvent is received by the active
// window. It expects a new frame in the form of a list of operations that
// describes what to display and how to handle input. This operations list
// is returned to the caller for displaying on screen.
func (win *Window) handleFrameEvent(evt system.FrameEvent) *op.Ops {
	switch {
	case win.navigator.CurrentPage() == nil:
		// Prepare to display the StartPage if no page is currently displayed.
		win.navigator.Display(page.NewStartPage(win.load))

	default:
		// The app window may have received some user interaction such as key
		// presses, a button click, etc which triggered this FrameEvent. Handle
		// such interactions before re-displaying the UI components. This
		// ensures that the proper interface is displayed to the user based on
		// the action(s) they just performed.
		win.handleRelevantKeyPresses(evt)
		win.navigator.CurrentPage().HandleUserInteractions()
		if modal := win.navigator.TopModal(); modal != nil {
			modal.Handle()
		}
	}

	// Generate an operations list with instructions for drawing the window's UI
	// components onto the screen. Use the generated ops to request key events.
	ops := win.prepareToDisplayUI(evt)
	win.addKeyEventRequestsToOps(ops)

	return ops
}

// handleRelevantKeyPresses checks if any open modal or the current page is a
// load.KeyEventHandler AND if the provided system.FrameEvent contains key press
// events for the modal or page.
func (win *Window) handleRelevantKeyPresses(evt system.FrameEvent) {
	handleKeyPressFor := func(tag string, maybeHandler interface{}) {
		handler, ok := maybeHandler.(load.KeyEventHandler)
		if !ok {
			return
		}
		for _, event := range evt.Queue.Events(tag) {
			if keyEvent, isKeyEvent := event.(key.Event); isKeyEvent && keyEvent.State == key.Press {
				handler.HandleKeyPress(&keyEvent)
			}
		}
	}

	// Handle key events on the top modal first, if there's one.
	// Only handle key events on the current page if no modal is displayed.
	if modal := win.navigator.TopModal(); modal != nil {
		handleKeyPressFor(modal.ID(), modal)
	} else {
		handleKeyPressFor(win.navigator.CurrentPageID(), win.navigator.CurrentPage())
	}
}

// prepareToDisplayUI creates an operation list and writes the layout of all the
// window UI components into it. The created ops is returned and may be used to
// record further operations before finally being rendered on screen via
// system.FrameEvent.Frame(ops).
func (win *Window) prepareToDisplayUI(evt system.FrameEvent) *op.Ops {
	backgroundWidget := layout.Expanded(func(gtx C) D {
		return decredmaterial.Fill(gtx, win.load.Theme.Color.Gray4)
	})

	currentPageWidget := layout.Stacked(func(gtx C) D {
		if modal := win.navigator.TopModal(); modal != nil {
			gtx = gtx.Disabled()
		}
		return win.navigator.CurrentPage().Layout(gtx)
	})

	topModalLayout := layout.Stacked(func(gtx C) D {
		modal := win.navigator.TopModal()
		if modal == nil {
			return layout.Dimensions{}
		}
		return modal.Layout(gtx)
	})

	// Use a StackLayout to write the above UI components into an operations
	// list via a graphical context that is linked to the ops.
	ops := &op.Ops{}
	gtx := layout.NewContext(ops, evt)
	layout.Stack{Alignment: layout.N}.Layout(
		gtx,
		backgroundWidget,
		currentPageWidget,
		topModalLayout,
		layout.Stacked(win.load.Toast.Layout),
	)

	return ops
}

// addKeyEventRequestsToOps checks if the current page or any modal has
// registered to be notified of certain key events and updates the provided
// operations list with instructions to generate a FrameEvent if any of the
// desired keys is pressed on the window.
func (win *Window) addKeyEventRequestsToOps(ops *op.Ops) {
	requestKeyEvents := func(tag string, desiredKeys key.Set) {
		if desiredKeys == "" {
			return
		}

		// Execute the key.InputOP{}.Add operation after all other operations.
		// This is particularly important because some pages call op.Defer to
		// signfiy that some operations should be executed after all other
		// operations, which has an undesirable effect of discarding this key
		// operation unless it's done last, after all other defers are done.
		m := op.Record(ops)
		key.InputOp{Tag: tag, Keys: desiredKeys}.Add(ops)
		op.Defer(ops, m.Stop())
	}

	// Request key events on the top modal, if necessary.
	// Only request key events on the current page if no modal is displayed.
	if modal := win.navigator.TopModal(); modal != nil {
		if handler, ok := modal.(load.KeyEventHandler); ok {
			requestKeyEvents(modal.ID(), handler.KeysToHandle())
		}
	} else {
		if handler, ok := win.navigator.CurrentPage().(load.KeyEventHandler); ok {
			requestKeyEvents(win.navigator.CurrentPageID(), handler.KeysToHandle())
		}
	}
}
