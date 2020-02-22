package page

import (
	"errors"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

const (
	// WalletsID is the id of the wallets page
	WalletsID = "wallets"

	//StateDeletedWallet is the map key for a deleted walet event
	StateDeletedWallet = "deleted"

	StateDeleteErrorWallet = "deleteerror"
)

// Wallets contains the Wallet, Theme and various widgets
type Wallets struct {
	wal   *wallet.Wallet
	theme *materialplus.Theme

	states map[string]interface{}

	add, sync widget.Button
	selected  int

	expandBtns []widget.Button

	addIcon, expandIcon                 *material.Icon
	delete, dialogCancel, dialogConfirm widget.Button

	editor widget.Editor

	dialog, deleting bool
}

// Init stores the theme and Wallet
func (pg *Wallets) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.wal = wal
	pg.theme = theme

	expand, err := material.NewIcon(icons.NavigationExpandMore)
	pg.expandIcon = expand

	pg.addIcon, err = material.NewIcon(icons.ContentAdd)

	pg.states = states
	pg.expandBtns = make([]widget.Button, 0)

	pg.editor = widget.Editor{SingleLine: true, Submit: true}

	log.Tracef("Handle %v", err)
}

// Draw layouts out the widgets on the given context
func (pg *Wallets) Draw(gtx *layout.Context) (event interface{}) {
	info := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	stateErr, haveError := pg.states[StateError].(error)
	numWallets := len(info.Wallets)
	if numWallets == 0 {
		pg.theme.Label(units.Label, "Retrieving Wallet Info").Layout(gtx)
		return nil
	}

	if len(pg.expandBtns) != numWallets {
		pg.expandBtns = make([]widget.Button, numWallets)
	}

	// if info.Synced {
	// 	pg.syncWdg.Text = "Synced"
	// } else if info.Syncing {
	// 	pg.syncWdg.Text = "Cancel Sync"
	// } else {
	// 	pg.syncWdg.Text = "Start Sync"
	// }

	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(.35, func() {
			layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func() {
					pg.theme.H2("Wallets").Layout(gtx)
				}),
				layout.Rigid(func() {
					pg.theme.IconButton(pg.addIcon).Layout(gtx, &pg.add)
					if pg.add.Clicked(gtx) {
						//log.Debugf("{%s} AddWallet Btn Clicked", WalletsID)
						event = EventNav{
							Current: WalletsID,
							Next:    LandingID,
						}
					}
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(info.Wallets), func(i int) {
						//pg.theme.LayoutWithBackGround(gtx, false, func() {
						//gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
						layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Flexed(.3, func() {
								ui.Center(gtx, func() { pg.theme.H4(info.Wallets[i].Name).Layout(gtx) })
							}),
							layout.Rigid(func() {
								ui.Center(gtx, func() {
									pg.theme.H5(info.Wallets[i].Balance.String()).Layout(gtx)
								})
							}),
							layout.Rigid(func() {
								ui.Center(gtx, func() {
									if pg.selected != i {
										pg.theme.IconButton(pg.expandIcon).Layout(gtx, &pg.expandBtns[i])
										if pg.expandBtns[i].Clicked(gtx) {
											pg.selected = i
										}
									}
								})
							}),
						)
						//})
					})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(.2, func() {
							layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx,
								layout.Rigid(func() {
									ui.Center(gtx, func() { pg.theme.H3("Balance").Layout(gtx) })
								}),
								layout.Rigid(func() {
									ui.Center(gtx, func() { pg.theme.H4(info.Wallets[pg.selected].Balance.String()).Layout(gtx) })
								}),
							)
						}),
						layout.Rigid(func() {
							del := pg.theme.Button("Delete wallet")
							del.Background = pg.theme.Danger
							del.Layout(gtx, &pg.delete)
							if pg.delete.Clicked(gtx) {
								pg.dialog = true
							}
						}),
					)
				}),
			)
		}),
	)

	confirm := pg.theme.Button("Confirm")
	if pg.deleting {
		confirm.Background = ui.GrayColor
	}
	lbl := pg.theme.H3("Are you sure?")
	if haveError {
		lbl.Color = pg.theme.Danger
		if errors.Is(stateErr, wallet.ErrBadPass) {
			lbl.Text = "Invalid password"
		} else {
			lbl.Text = "Something went wrong. See log for details"
			//log.Error(stateErr.Error())
		}
		pg.deleting = false
	}

	if _, ok := pg.states[StateDeletedWallet]; ok {
		lbl.Text = "Wallet deleted"
	}

	ui.Dialog{
		ConfirmButton: confirm,
		Confirm:       &pg.dialogConfirm,
		CancelButton:  pg.theme.ButtonWithColor("Cancel", pg.theme.Danger),
		Cancel:        &pg.dialogCancel,
	}.Layout(gtx, pg.dialog, func() {
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func() {
				ui.Center(gtx, func() { lbl.Layout(gtx) })
			}),
			layout.Rigid(func() {
				ui.Center(gtx, func() { pg.theme.Editor("Enter spending password").Layout(gtx, &pg.editor) })
			}),
		)
	})

	for pg.dialogCancel.Clicked(gtx) {
		pg.dialog = false
		pg.deleting = false
		pg.editor.SetText("")
		delete(pg.states, StateError)
		delete(pg.states, StateDeletedWallet)
	}

	for pg.dialogConfirm.Clicked(gtx) {
		pass := pg.editor.Text()
		//log.Debug("Confirmed with pass " + pass)
		if pass != "" && !pg.deleting {
			pg.deleting = true
			pg.wal.DeleteWallet(info.Wallets[pg.selected].ID, pass)
			if stateErr != nil {
				delete(pg.states, StateError)
			}
		} else {

		}
	}

	if pg.sync.Clicked(gtx) {
		if !info.Synced {
			log.Info("Starting sync")
			if err := pg.wal.StartSync(); err != nil {
				log.Error(err)
				//pg.syncWdg.Text = "Error starting sync"
			} else {
				//pg.syncWdg.Text = "Cancel Sync"
			}
		}
		if info.Syncing {
			log.Info("Cancelling sync")
			pg.wal.CancelSync()
		}
	}
	return
}
