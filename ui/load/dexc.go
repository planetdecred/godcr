package load

import (
	"decred.org/dcrdex/client/core"
	"github.com/planetdecred/godcr/dexc"
)

type DexcLoad struct {
	*core.Core
	Dexc       *dexc.Dexc
	IsLoggedIn bool // Keep user logged in state
}
