package app

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	giouiApp "gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/notification"
	"github.com/planetdecred/godcr/ui/values"
)

// App is app. LOL
type App struct {
	net, version string
	mw           *dcrlibwallet.MultiWallet

	window        *giouiApp.Window
	currentPage   Page
	pageBackStack []Page // investigate use, also if mutex required

	modalMutex sync.Mutex
	modals     []Modal

	// TODO: Unexport?
	Theme *decredmaterial.Theme
	Toast *notification.Toast
}

// Init initializes an app state with a MultiWallet instance.
func Init(rootDir, net, version string) (*App, error) {
	if rootDir == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}

	appTheme := decredmaterial.NewTheme(assets.FontCollection(), assets.DecredIcons, false)
	if appTheme == nil {
		return nil, errors.New("unexpected error while loading theme")
	}

	appTitle := giouiApp.Title(values.StringF(values.StrAppTitle, net)) // before renaming testnet to testnet3 below

	politeiaHost := dcrlibwallet.PoliteiaMainnetHost
	if net == "testnet" {
		net = dcrlibwallet.Testnet3
		politeiaHost = dcrlibwallet.PoliteiaTestnetHost
	}

	mw, err := dcrlibwallet.NewMultiWallet(rootDir, "bdb", net, politeiaHost)
	if err != nil {
		return nil, err
	}

	// Restore/set theme colors based on saved user pref.
	isDarkModeOn := mw.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false)
	appTheme.SwitchDarkMode(isDarkModeOn, assets.DecredIcons)

	return &App{
		net:     net,
		version: version,
		mw:      mw,
		window:  giouiApp.NewWindow(giouiApp.MinSize(values.AppWidth, values.AppHeight), appTitle),
		Theme:   appTheme,
		Toast:   notification.NewToast(appTheme),
	}, nil
}

// MultiWallet ensures read-only access to the MultiWallet instance.
func (app *App) MultiWallet() *dcrlibwallet.MultiWallet {
	return app.mw
}

// Run blocks until the app is exited.
func (app *App) Run(startPage Page) {
	app.currentPage = startPage
	app.currentPage.OnNavigatedTo()

	go func() {
		app.handleEvents()
		app.mw.Shutdown()
		os.Exit(0)
	}()
	giouiApp.Main()
}

// handleEvents runs main event handling and page rendering loop.
func (app *App) handleEvents() {
	for {
		e := <-app.window.Events()
		switch evt := e.(type) {

		case system.DestroyEvent:
			if app.currentPage != nil {
				app.currentPage.OnNavigatedFrom()
				app.currentPage = nil
			}
			return // exits the loop, caller will exit the program.

		case system.FrameEvent:
			app.displayWindow(evt)

		case key.Event:
			app.handleKeyEvent(&evt) // TODO: Use pointer?

		default:
			// log.Tracef("Unhandled window event %+v\n", e)
		}
	}
}

// displayWindow is called when a FrameEvent is received by the active window.
// Since user actions such as button clicks also trigger FrameEvents, this
// method first checks for pending user actions before displaying the UI
// elements. This ensures that the proper interface is displayed to the user
// based on their last performed action where applicable.
func (app *App) displayWindow(evt system.FrameEvent) {
	// A FrameEvent may be generated because of a user interaction
	// with the current page such as a button click. First handle
	// any such user interaction before rendering the page.
	app.currentPage.HandleUserInteractions()
	// app.modalMutex.Lock()
	for _, modal := range app.modals {
		modal.Handle() // TODO: Just the top-most modal should do.
	}
	// app.modalMutex.Unlock()
	// TODO!

	// Draw the window's UI components into an op.Ops.
	gtx := layout.NewContext(&op.Ops{}, evt)
	app.drawWindowUI(gtx)

	// Render the window's UI components on screen.
	evt.Frame(gtx.Ops)
}

// drawWindowUI draws the window UI components into the provided graphical
// context, preparing the context for rendering on screen.
func (app *App) drawWindowUI(gtx layout.Context) {
	// Create a base view holder to hold all the following UI components
	// one on top the other. Components that do not take up the entire
	// window will be aligned to the top of the window.
	viewsHolder := layout.Stack{Alignment: layout.N}

	background := layout.Expanded(func(gtx layout.Context) layout.Dimensions {
		return decredmaterial.Fill(gtx, app.Theme.Color.Gray4)
	})

	// TODO: Should suffice to just draw the top-most modal?
	modals := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		// TODO!
		// app.modalMutex.Lock()
		// defer app.modalMutex.Unlock()

		if len(app.modals) == 0 {
			return layout.Dimensions{}
		}

		modalLayouts := make([]layout.StackChild, 0)
		for _, modal := range app.modals {
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
		layout.Stacked(app.currentPage.Layout),
		modals,
		layout.Stacked(app.Toast.Layout),
	)
}

func (app *App) handleKeyEvent(evt *key.Event) {
	if handler, ok := app.currentPage.(KeyEventHandler); ok {
		handler.HandleKeyEvent(evt)
	}
	for _, modal := range app.modals { // TODO: Lock
		if handler, ok := modal.(KeyEventHandler); ok {
			handler.HandleKeyEvent(evt)
		}
	}
}

// TODO: Rename to ReloadApp or something and doc.
func (app *App) RefreshWindow() {
	app.window.Invalidate()
}

// ReloadApp closes the current page active on the
// app window. When the next FrameEvent is received,
// a new StartPage will be initialized and displayed.
// TODO: Callers shouldn't need to refresh window.
func (app *App) ReloadApp() { // TODO: Reload()?
	if app.currentPage != nil {
		app.currentPage.OnNavigatedFrom()
		app.currentPage = nil
	}
	app.window.Invalidate()
}

// TotalBalance is the sum of total and spendable balances of all controlled
// wallets.
func (app *App) TotalBalance() (dcrutil.Amount, dcrutil.Amount, error) {
	totalBalance := int64(0)
	spandableBalance := int64(0)

	for _, wallet := range app.mw.AllWallets() {
		accountsResult, err := wallet.GetAccountsRaw()
		if err != nil {
			return 0, 0, err
		}
		for _, account := range accountsResult.Acc {
			totalBalance += account.TotalBalance
			spandableBalance += account.Balance.Spendable
		}
	}

	return dcrutil.Amount(totalBalance), dcrutil.Amount(spandableBalance), nil
}

// Wallets returns the slice of all controlled wallets sorted by ID.
func (app *App) Wallets() []*dcrlibwallet.Wallet {
	wallets := app.mw.AllWallets()
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].ID < wallets[j].ID
	})
	return wallets
}
