package page

import (
	"gioui.org/layout"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
)

// CreateWalletID is the id of the createwallet page
const CreateWalletID = "createWallet"

// CreateWallet represents the wallet creation page
type CreateWallet struct {
	passwordAndPin *materialplus.PasswordAndPin
}

// Init initializes the create wallet page widgets
func (pg *CreateWallet) Init(theme *materialplus.Theme) {
	pg.passwordAndPin = theme.PasswordAndPin()
}

// Draw renders the page's widgets to screen
func (pg *CreateWallet) Draw(gtx *layout.Context, _ event.Event) (evt event.Event) {
	pg.passwordAndPin.Draw(gtx, pg.createFunc, pg.cancelFunc)
	return nil
}

func (pg *CreateWallet) cancelFunc() {

}

func (pg *CreateWallet) createFunc(password string) {

}
