package uidex

import (
	"github.com/planetdecred/godcr/dex"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	loading  bool // true if the window is in the middle of an operation that cannot be stopped
	creating bool // true if a wallet is being created or restored
}

// updateStates changes the dexc state based on the received update
func (d *Dex) updateStates(update interface{}) {
	switch e := update.(type) {
	case dex.User:
		d.userInfo = &e
		return
	}
}
