package uidex

import (
	"errors"
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"

	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

// Dex represents the Dex UI. There should only be one.
type DexUI struct {
	ops      op.Ops
	theme    *decredmaterial.Theme
	dexc     *dexc.Dexc
	userInfo *dexc.User

	// market current selected
	market           *selectedMaket
	maxOrderEstimate *dexc.MaxOrderEstimate

	current, previous string

	err string

	pages map[string]layout.Widget

	// Toggle between wallet and dex view mode
	switchView *int

	refreshWindow func()
}
type selectedMaket struct {
	host          string
	name          string
	marketBase    string
	marketQuote   string
	marketBaseID  uint32
	marketQuoteID uint32
}

// NewDexUI creates and initializes a new walletUI with start
func NewDexUI(dc *dexc.Dexc, decredIcons map[string]image.Image, collection []text.FontFace, v *int, invalidate func()) (*DexUI, error) {
	d := new(DexUI)
	d.dexc = dc
	theme := decredmaterial.NewTheme(collection, decredIcons)
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	d.ops = op.Ops{}
	d.theme = theme

	d.userInfo = new(dexc.User)
	d.market = new(selectedMaket)
	d.current = PageMarkets
	d.switchView = v
	d.addPages(decredIcons)

	d.refreshWindow = invalidate

	return d, nil
}

func (d *DexUI) Ops() *op.Ops {
	return &d.ops
}

func (d *DexUI) changePage(page string) {
	d.current = page
	d.refresh()
}

func (d *DexUI) refresh() {
	d.refreshWindow()
}

func (d *DexUI) setReturnPage(from string) {
	d.previous = from
	d.refresh()
}

// Run runs main event handling and page rendering loop
func (d *DexUI) Run(shutdown chan int, w *app.Window) {
	for {
		select {
		case e := <-d.dexc.Send:
			if e.Err != nil {
				log.Error(e.Err)
			}
			d.updateStates(e.Resp)
		}
	}
}

func (d *DexUI) HandlerDestroy(shutdown chan int) {

}

func (d *DexUI) HandlerPages(gtx layout.Context) {
	d.theme.Background(gtx, d.pages[d.current])
}
