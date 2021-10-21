package dexclient

import (
	"fmt"

	"decred.org/dcrdex/client/core"
)

func (pg *Page) connectDex(h string, password []byte) {
	go pg.readNotifications()
	// TODO: connect to dex server to listen messages
	fmt.Println(h, password[0])
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc.NotificationFeed()
	for {
		select {
		case n := <-ch:
			fmt.Println("Recv notification", n)
			fmt.Println("<INFO>", n.ID())
			fmt.Println("<INFO>", n.Severity())
			fmt.Println("<INFO>", n.Type())
			fmt.Println("<INFO>", n.DBNote())
			fmt.Println("<INFO>", n.String())
			fmt.Println("<INFO>", n.Subject())

			if n.Type() == core.NoteTypeFeePayment {
				pg.RefreshWindow()
			}
		case <-pg.ctx.Done():
			return
		}
	}
}
