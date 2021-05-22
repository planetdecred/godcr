package uidex

import (
	"errors"
	"image"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"

	"github.com/planetdecred/godcr/dex"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
)

// Dex represents the Dex UI. There should only be one.
type Dex struct {
	ops      op.Ops
	theme    *decredmaterial.Theme
	dexc     *dex.Dex
	userInfo *dex.User

	current, previous string

	signatureResult *wallet.Signature

	err string

	pages                 map[string]layout.Widget
	sysDestroyWithSync    bool
	walletAcctMixerStatus chan *wallet.AccountMixer
	internalLog           chan string

	// Toggle between wallet and dex view mode
	switchView *int
}

type WriteClipboard struct {
	Text string
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

func (d *Dex) unloaded(w *app.Window) {
	lbl := d.theme.H3("Multiwallet not loaded\nIs another instance open?")
	var ops op.Ops

	for {
		e := <-w.Events()
		switch evt := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, evt)
			lbl.Layout(gtx)
			evt.Frame(gtx.Ops)
		}
	}
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
