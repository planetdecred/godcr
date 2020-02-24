// +build !dev

package ui

import (
	"time"

	"gioui.org/io/system"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.wallet.Send:
			log.Debugf("Recieved event %+v", e)
			if e.Err != nil {
				win.window.Invalidate()
				break
			}
			switch evt := e.Resp.(type) {
			case *wallet.LoadedWallets:
				win.wallet.GetMultiWalletInfo()
				if evt.Count == 0 {
					win.current = Loading
				}
				//else {
				//	win.current =
				//}
			case *wallet.MultiWalletInfo:
				*win.walletInfo = *evt
			default:
				win.updateState(e.Resp)
			}
			// set error if it exists
			if e.Err != nil {
				//win.states[page.StateError] = e.Err
			}
			win.window.Invalidate()
		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				close(shutdown)
				return
			case system.FrameEvent:
				win.gtx.Reset(evt.Config, evt.Size)
				start := time.Now()
				win.current(win.gtx, win.theme, win.walletInfo)
				log.Tracef("Page {%s} rendered in %v", win.current, time.Since(start))
				evt.Frame(win.gtx.Ops)
				win.HandleInputs()
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}
}
