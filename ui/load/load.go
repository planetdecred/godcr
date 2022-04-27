// The load package contains data structures that are shared by components in the ui package. It is not a dumping ground
// for code you feel might be shared with other components in the future. Before adding code here, ask yourself, can
// the code be isolated in the package you're calling it from? Is it really needed by other packages in the ui package?
// or you're just planning for a use case that might never used.

package load

import (
	"golang.org/x/text/message"

	"gioui.org/io/key"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/notification"
	"github.com/planetdecred/godcr/wallet"
)

type DCRUSDTBittrex struct {
	LastTradeRate string
}

type Receiver struct {
	KeyEvents map[string]chan *key.Event
}

type Load struct {
	Theme *decredmaterial.Theme

	WL       *WalletLoad
	Receiver *Receiver
	Printer  *message.Printer
	Network  string

	Toast *notification.Toast

	SelectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	ToggleSync          func()
	RefreshWindow       func()
	ShowModal           func(Modal)
	DismissModal        func(Modal)
	ChangeWindowPage    func(page Page, keepBackStack bool)
	PopWindowPage       func() bool
	ChangeFragment      func(page Page)
	PopFragment         func()
	PopToFragment       func(pageID string)
	SubscribeKeyEvent   func(eventChan chan *key.Event, pageID string) // Widgets call this function to recieve key events.
	UnsubscribeKeyEvent func(pageID string) error
	ReloadApp           func()

	DarkModeSettingChanged func(bool)
	LanguageSettingChanged func()
	CurrencySettingChanged func()
}

func (l *Load) RefreshTheme() {
	isDarkModeOn := l.WL.MultiWallet.ReadBoolConfigValueForKey(DarkModeConfigKey, false)
	l.Theme.SwitchDarkMode(isDarkModeOn, assets.DecredIcons)
	l.DarkModeSettingChanged(isDarkModeOn)
	l.LanguageSettingChanged()
	l.CurrencySettingChanged()
	l.RefreshWindow()
}

func (l *Load) Dexc() *dcrlibwallet.DexClient {
	return l.WL.MultiWallet.DexClient()
}
