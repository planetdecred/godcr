package uidex

import (
	"errors"
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"

	"github.com/planetdecred/godcr/dex"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

// Dex represents the Dex UI. There should only be one.
type Dex struct {
	ops      op.Ops
	theme    *decredmaterial.Theme
	dexc     *dex.Dex
	userInfo *dex.User

	// market current selected
	market *selectedMaket

	current, previous string

	err string

	pages       map[string]layout.Widget
	internalLog chan string

	// Toggle between wallet and dex view mode
	switchView *int
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
func NewDexUI(dexc *dex.Dex, decredIcons map[string]image.Image, collection []text.FontFace, internalLog chan string, v *int) (*Dex, error) {
	d := new(Dex)
	d.dexc = dexc
	theme := decredmaterial.NewTheme(collection, decredIcons)
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	d.ops = op.Ops{}
	d.theme = theme
	d.internalLog = internalLog

	d.userInfo = new(dex.User)
	d.market = new(selectedMaket)
	d.current = PageMarkets
	d.switchView = v
	d.addPages(decredIcons)

	return d, nil
}

func (d *Dex) Ops() *op.Ops {
	return &d.ops
}

func (d *Dex) changePage(page string) {
	d.current = page
	d.refresh()
}

func (d *Dex) refresh() {
}

func (d *Dex) setReturnPage(from string) {
	d.previous = from
	d.refresh()
}

// Run runs main event handling and page rendering loop
func (d *Dex) Run(shutdown chan int, w *app.Window) {
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

func (d *Dex) HandlerDestroy(shutdown chan int) {

}

func (d *Dex) HandlerPages(gtx layout.Context) {
	d.theme.Background(gtx, d.pages[d.current])
}
