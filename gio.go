package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr-gio/giolog"
	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/pages/common"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type (
	desktop struct {
		window          *gioapp.Window
		displayName     string
		pages           []navPage
		standalonePages map[string]standalonePageHandler
		currentPage     string
		pageChanged     bool
		appDisplayName  string
		multiWallet     *helper.MultiWallet
		syncer          *common.Syncer
	}
)

const (
	windowWidth  = 520
	windowHeight = 530

	navSectionHeight = 70
)

func launchUserInterface(appDisplayName, appDataDir, netType string) {
	logger, err := dcrlibwallet.RegisterLogger("GIOL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot register logger: %v", err)
		return
	}
	giolog.UseLogger(logger)

	// initialize theme
	helper.Initialize()

	app := &desktop{
		currentPage: "overview",
		pageChanged: true,
	}

	theme := helper.GetTheme()
	err = helper.InitImages(theme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot load assets: %v", err)
		return
	}

	multiWallet, shouldCreateOrRestoreWallet, shouldPromptForPass, err := helper.LoadWallet(appDataDir, netType)
	if err != nil {
		// todo show error in UI
		giolog.Log.Errorf(err.Error())
		return
	}

	app.multiWallet = multiWallet
	if shouldCreateOrRestoreWallet {
		app.currentPage = "welcome"
	} else if shouldPromptForPass {
		//app.currentPage = "passphrase"
		// TODO prompt for passphrase
	}

	app.syncer = common.NewSyncer(app.multiWallet, app.refreshWindow)
	app.multiWallet.AddSyncProgressListener(app.syncer, app.appDisplayName)

	app.prepareHandlers()
	go func() {
		app.window = gioapp.NewWindow(
			gioapp.Size(unit.Dp(windowWidth), unit.Dp(windowHeight)),
			gioapp.Title(app.displayName),
		)

		if err := app.renderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	// run app
	gioapp.Main()
}

func (d *desktop) prepareHandlers() {
	// set standalone page
	d.standalonePages = getStandalonePages(d.multiWallet)

	// set navPages
	d.pages = getNavPages()
	if len(d.pages) > 0 && d.currentPage == "" {
		d.changePage(d.pages[0].name)
	}
}

func (d *desktop) changePage(pageName string) {
	if d.currentPage == pageName {
		return
	}
	d.currentPage = pageName
	d.pageChanged = true
}

func (d *desktop) renderLoop() error {
	ctx := layout.NewContext(d.window.Queue())

	for {
		e := <-d.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ctx.Reset(e.Config, e.Size)
			d.render(ctx)
			e.Frame(ctx.Ops)
		}
	}
}

func (d *desktop) render(ctx *layout.Context) {
	helper.PaintArea(ctx, helper.BackgroundColor, windowWidth, windowHeight)

	// first check if current page is standalone and render
	if page, ok := d.standalonePages[d.currentPage]; ok {
		d.renderStandalonePage(page, ctx)
	} else {
		var page navPage
		for i := range d.pages {
			if d.pages[i].name == d.currentPage {
				page = d.pages[i]
				break
			}
		}

		if d.pageChanged {
			d.pageChanged = false
			page.handler.BeforeRender(d.syncer, d.multiWallet)
		}
		d.renderNavPage(page, ctx)
	}
	d.refreshWindow()
}

func (d *desktop) renderNavPage(page navPage, ctx *layout.Context) {
	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
				layout.Rigid(func() {
					helper.LogoSymbol.Layout(ctx)
				}),
				layout.Flexed(1, func() {
					inset := layout.Inset{
						Top: unit.Dp(17),
					}
					inset.Layout(ctx, func() {
						widgets.NewLabel(page.label).
							SetSize(6).
							SetWeight(text.Bold).
							Draw(ctx)
					})
				}),
			)
		}),
		layout.Expanded(func() {
			inset := layout.Inset{
				Top:   unit.Dp(55),
				Left:  unit.Dp(15),
				Right: unit.Dp(15),
			}
			inset.Layout(ctx, func() {
				page.handler.Render(ctx, d.changePage)
			})
		}),
		layout.Stacked(func() {
			d.renderNavSection(ctx)
		}),
	)
}

func (d *desktop) renderStandalonePage(page standalonePageHandler, ctx *layout.Context) {
	inset := layout.Inset{}
	inset.Layout(ctx, func() {
		layout.Stack{Alignment: layout.NW}.Layout(ctx,
			layout.Stacked(func() {
				inset := layout.Inset{Top: unit.Dp(helper.StandaloneScreenPadding)}
				inset.Layout(ctx, func() {
					page.Render(ctx, d.refreshWindow, d.changePage)
				})
			}),
		)
	})
}

func (d *desktop) renderNavSection(ctx *layout.Context) {
	inset := layout.Inset{
		Top: unit.Dp(windowHeight - navSectionHeight),
	}
	inset.Layout(ctx, func() {
		helper.PaintFooter(ctx, helper.WhiteColor, windowWidth, navSectionHeight)

		navItemWidth := ctx.Constraints.Width.Max / 4
		(&layout.List{Axis: layout.Horizontal}).Layout(ctx, 4, func(i int) {
			if i > 3 {
				return
			}
			d.pages[i].button.DrawNavItem(ctx, d.pages[i].icon, navItemWidth, func() {
				d.changePage(d.pages[i].name)
			})
		})
	})
}

func (d *desktop) refreshWindow() {
	d.window.Invalidate()
}
