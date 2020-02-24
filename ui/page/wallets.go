package page

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/icons"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/styles"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

const (
	// WalletsID is the id of the wallets page
	WalletsID = "wallets"

	//StateDeletedWallet is the map key for a deleted walet event
	StateDeletedWallet = "deleted"
)

// Wallets contains the Wallet, Theme and various widgets
type Wallets struct {
	wal   *wallet.Wallet
	theme *materialplus.Theme

	states map[string]interface{}

	add, sync, confirm, cancel, delete widget.Button
	selected                           int

	tabBtns []*widget.Button

	editor *widget.Editor

	processing, deleted, deleting bool

	err error
}

// Init stores the theme and Wallet
func (pg *Wallets) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.wal = wal
	pg.theme = theme

	pg.states = states
	pg.tabBtns = make([]*widget.Button, 0)

	pg.editor = &widget.Editor{SingleLine: true, Submit: true}
}

func (pg *Wallets) drawHeader(gtx *layout.Context, info *wallet.MultiWalletInfo) (event interface{}) {
	layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(.3, func() {
			pg.theme.H2("Wallets").Layout(gtx)
		}),
		layout.Rigid(func() {
			pg.theme.IconButton(icons.ContentAdd).Layout(gtx, &pg.add)
		}),
		layout.Rigid(func() {
			status := pg.theme.Label(unit.Dp(20), "Not Synced")
			if info.Synced {
				status.Text = "Synced"
			} else if info.Syncing {
				status.Text = "Syncing"
			}
			status.Layout(gtx)
		}),
		layout.Rigid(func() {
			btn := pg.theme.IconButton(icons.NavigationRefresh)
			if info.Synced {
				btn = pg.theme.IconButton(icons.NavigationCheck)
				btn.Background = pg.theme.Success
			} else if info.Syncing {
				btn = pg.theme.IconButton(icons.NavigationClose)
				btn.Background = pg.theme.Danger
			}
			btn.Layout(gtx, &pg.sync)
		}),
	)
	return
}

func (pg *Wallets) drawWalletInfo(gtx *layout.Context, info wallet.InfoShort) {
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(.2, func() {
			layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func() {
					pg.theme.H3("Balance").Layout(gtx)
				}),
				layout.Rigid(func() {
					pg.theme.H4(info.Balance.String()).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(.6, func() {

		}),
		layout.Flexed(.2, styles.WithStyles(gtx, func() {
			pg.theme.DangerButton("Delete wallet").Layout(gtx, &pg.delete)
		}, styles.Maxed)),
	)

}

func (pg *Wallets) drawDialog(gtx *layout.Context, info *wallet.MultiWalletInfo) {
	if pg.deleting {
		diag := pg.theme.PasswordDialog("Enter spending password")
		diag.Layout(gtx, pg.editor, &pg.confirm, &pg.cancel)
	}

}

// Draw layouts out the widgets on the given context
func (pg *Wallets) Draw(gtx *layout.Context) (event interface{}) {
	info := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)

	if _, ok := pg.states[StateDeletedWallet]; ok {
		pg.selected = 0
		pg.processing = false
		pg.deleted = true
		pg.err = nil
		delete(pg.states, StateDeletedWallet)
		delete(pg.states, StateError)
	}

	if err, ok := pg.states[StateError].(error); ok {
		pg.err = err
	}

	numWallets := len(info.Wallets)
	if numWallets == 0 {
		pg.theme.Label(units.Label, "Retrieving Wallet Info").Layout(gtx)
		return nil
	}

	if len(pg.tabBtns) != numWallets {
		pg.tabBtns = make([]*widget.Button, numWallets)
		for i := range pg.tabBtns {
			pg.tabBtns[i] = new(widget.Button)
		}
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(.15, func() {
			evt := pg.drawHeader(gtx, info)
			if event == nil {
				event = evt
			}
		}),
		layout.Flexed(.85,
			pg.theme.Tabbed(gtx,
				&pg.selected,
				pg.tabBtns,
				styles.WithStyle(gtx, styles.Background(pg.theme.Background), func() {

				}),
				func(i int) {
					pg.theme.Label(unit.Dp(50), info.Wallets[i].Name).Layout(gtx)
				},
				styles.WithStyle(gtx, styles.Background(pg.theme.Background), func() {
					pg.drawWalletInfo(gtx, info.Wallets[pg.selected])
				}),
			)),
	)

	pg.drawDialog(gtx, info)
	for i, btn := range pg.tabBtns {
		if btn.Clicked(gtx) {
			log.Debugf("Tab %d selected", i)
			pg.selected = i
		}
	}

	for pg.sync.Clicked(gtx) {
		if !info.Synced {
			if err := pg.wal.StartSync(); err != nil {
				log.Error(err)
				pg.err = err
			}
		}
		if info.Syncing {
			log.Info("Cancelling sync")
			pg.wal.CancelSync()
		}
	}

	if pg.add.Clicked(gtx) {
		log.Tracef("{%s} AddWallet Btn Clicked", WalletsID)
		event = EventNav{
			Current: WalletsID,
			Next:    LandingID,
		}
	}
	if pg.delete.Clicked(gtx) {
		pg.deleting = true
	}
	return
}
