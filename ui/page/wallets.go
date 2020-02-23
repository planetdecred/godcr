package page

import (
	"errors"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/layouts"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/icons"
	"github.com/raedahgroup/godcr-gio/ui/styles"
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

	add, sync widget.Button
	selected  int

	expandBtns []widget.Button

	delete, dialogCancel, dialogConfirm widget.Button

	editor widget.Editor

	dialog, processing, deleted bool
	prompt                      string
	onDialogConfirm             func(string)
	err                         error
	list                        *layout.List
}

// Init stores the theme and Wallet
func (pg *Wallets) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.wal = wal
	pg.theme = theme

	pg.states = states
	pg.expandBtns = make([]widget.Button, 0)

	pg.editor = widget.Editor{SingleLine: true, Submit: true}
	pg.list = (&layout.List{Axis: layout.Vertical})
}

func (pg *Wallets) drawHeader(gtx *layout.Context, info *wallet.MultiWalletInfo) (event interface{}) {
	layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(.3, func() {
			pg.theme.H2("Wallets").Layout(gtx)
		}),
		layout.Rigid(func() {
			pg.theme.IconButton(icons.ContentAdd).Layout(gtx, &pg.add)
			if pg.add.Clicked(gtx) {
				log.Tracef("{%s} AddWallet Btn Clicked", WalletsID)
				event = EventNav{
					Current: WalletsID,
					Next:    LandingID,
				}
			}
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
		}),
	)
	return
}

func (pg *Wallets) drawWalletInfo(gtx *layout.Context, info *wallet.MultiWalletInfo) (event interface{}) {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(.2, func() {
			layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
				layouts.RigidWithStyle(gtx, styles.Centered, func() {
					pg.theme.H3("Balance").Layout(gtx)
				}),
				layouts.RigidWithStyle(gtx, styles.Centered, func() {
					pg.theme.H4(info.Wallets[pg.selected].Balance.String()).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(.6, func() {

		}),
		layout.Flexed(.2, styles.WithStyle(gtx, styles.Centered, func() {
			pg.theme.DangerButton("Delete wallet").Layout(gtx, &pg.delete)
			if pg.delete.Clicked(gtx) {
				pg.showDialog("Delete Wallet?", func(pass string) {
					pg.wal.DeleteWallet(info.Wallets[pg.selected].ID, pass)
				})
			}
		})),
	)
	return
}

func (pg *Wallets) drawDialog(gtx *layout.Context, info *wallet.MultiWalletInfo) {
	dconfirm := &pg.dialogConfirm
	dcancel := &pg.dialogCancel

	label := pg.theme.H3(pg.prompt)
	confirm := pg.theme.Button("Confirm")
	cancel := pg.theme.Button("Cancel")
	cancel.Background = pg.theme.Danger

	labelLayout := func() {
		label.Layout(gtx)
	}

	diag := func() {
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layouts.RigidWithStyle(gtx, styles.Centered, labelLayout),
			layouts.RigidWithStyle(gtx, styles.Centered, func() {
				pg.theme.Editor("Enter spending password").Layout(gtx, &pg.editor)
			}),
		)
	}

	if pg.err != nil {
		label.Color = pg.theme.Danger
		if errors.Is(pg.err, wallet.ErrBadPass) {
			label.Text = "Invalid password"
		} else {
			label.Text = "Something went wrong. See log for details"
		}
		pg.processing = false
	} else if pg.processing {
		confirm.Background = pg.theme.Disabled
		cancel.Background = pg.theme.Disabled
	} else if pg.deleted {
		label.Text = "Wallet deleted"
		label.Color = pg.theme.Success
		diag = labelLayout
		dconfirm = nil
		cancel.Text = "Close dialog"
		cancel.Background = pg.theme.Theme.Color.Primary
	}

	layouts.Dialog{
		ConfirmButton: confirm,
		Confirm:       dconfirm,
		CancelButton:  cancel,
		Cancel:        dcancel,
		Active:        pg.dialog,
	}.Layout(gtx, diag)

	if !pg.processing {
		for pg.dialogCancel.Clicked(gtx) {
			pg.closeDialog()
			delete(pg.states, StateError)
			delete(pg.states, StateDeletedWallet)
		}
		for pg.dialogConfirm.Clicked(gtx) {
			pg.onDialogConfirm(pg.editor.Text())
			pg.processing = true
		}
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

	if len(pg.expandBtns) != numWallets {
		pg.expandBtns = make([]widget.Button, numWallets)
	}

	//selector := pg.drawSelector(gtx, info)
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(.15, func() {
			evt := pg.drawHeader(gtx, info)
			if event == nil {
				event = evt
			}
		}),
		layout.Flexed(.85, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(.3, func() {
					(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(info.Wallets), func(i int) {
						layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(.7, styles.WithStyles(gtx, func() {
								pg.theme.Label(unit.Dp(20), info.Wallets[i].Name+"\n"+info.Wallets[i].Balance.String()).Layout(gtx)
							}, styles.Centered)),
							layout.Flexed(.3, styles.WithStyles(gtx, func() {
								btn := pg.theme.IconButton(icons.NavigationArrowForward)
								if pg.selected == i {
									styles.FillWithColor(gtx, pg.theme.Color.Primary, false)
								}
								btn.Layout(gtx, &pg.expandBtns[i])
								if pg.expandBtns[i].Clicked(gtx) {
									pg.selected = i
								}
							}, styles.Centered)),
						)
					})
				}),
				layout.Flexed(.7, styles.WithStyle(gtx, styles.Background{Color: pg.theme.Color.Primary}, func() {
					evt := pg.drawWalletInfo(gtx, info)
					if event == nil {
						event = evt
					}
				})),
			)
		}),
	)
	pg.drawDialog(gtx, info)
	return
}

func (pg *Wallets) showDialog(prompt string, onConfirm func(string)) {
	pg.dialog = true
	pg.prompt = prompt
	pg.onDialogConfirm = onConfirm
}

func (pg *Wallets) closeDialog() {
	pg.dialog = false
	pg.deleted = false
	pg.prompt = ""
	pg.editor.SetText("")
	pg.onDialogConfirm = func(_ string) {}
}
