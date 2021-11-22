package dexclient

import "decred.org/dcrdex/client/core"

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			if n.Type() == core.NoteTypeFeePayment {
				pg.RefreshWindow()
			}
		case <-pg.ctx.Done():
			return
		}
	}
}
