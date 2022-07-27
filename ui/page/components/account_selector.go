package components

import (
	"context"
	"errors"
	"sync"

	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const AccoutSelectorID = "AccountSelector"

type AccountSelector struct {
	*load.Load
	*listeners.TxAndBlockNotificationListener

	selectedAccount *dcrlibwallet.Account
	accountIsValid  func(*dcrlibwallet.Account) bool
	callback        func(*dcrlibwallet.Account)

	openSelectorDialog *decredmaterial.Clickable
	selectorModal      *AccountSelectorModal

	dialogTitle  string
	totalBalance string
	changed      bool
}

// NewAccountSelector opens up a modal to select the desired account. If a
// nil value is passed for selectedWallet, then accounts for all wallets
// are shown, otherwise only accounts for the selectedWallet is shown.
func NewAccountSelector(l *load.Load) *AccountSelector {
	return &AccountSelector{
		Load:               l,
		accountIsValid:     func(*dcrlibwallet.Account) bool { return true },
		openSelectorDialog: l.Theme.NewClickable(true),
	}
}

func (as *AccountSelector) Title(title string) *AccountSelector {
	as.dialogTitle = title
	return as
}

func (as *AccountSelector) AccountValidator(accountIsValid func(*dcrlibwallet.Account) bool) *AccountSelector {
	as.accountIsValid = accountIsValid
	return as
}

func (as *AccountSelector) AccountSelected(callback func(*dcrlibwallet.Account)) *AccountSelector {
	as.callback = callback
	return as
}

func (as *AccountSelector) Changed() bool {
	changed := as.changed
	as.changed = false
	return changed
}

func (as *AccountSelector) Handle(window app.WindowNavigator) {
	for as.openSelectorDialog.Clicked() {
		as.selectorModal = newAccountSelectorModal(as.Load, as.selectedAccount).
			title(as.dialogTitle).
			accountValidator(as.accountIsValid).
			accountSelected(func(account *dcrlibwallet.Account) {
				if as.selectedAccount.Number != account.Number {
					as.changed = true
				}
				as.SetSelectedAccount(account)
				as.callback(account)
			}).
			onModalExit(func() {
				as.selectorModal = nil
			})
		window.ShowModal(as.selectorModal)
	}
}

// SelectFirstWalletValidAccount selects the first valid account from the
// first wallet in the SortedWalletList
// If selectedWallet is not nil, the first account for the selectWallet is selected.
func (as *AccountSelector) SelectFirstWalletValidAccount() error {
	if as.selectedAccount != nil && as.accountIsValid(as.selectedAccount) {
		as.UpdateSelectedAccountBalance()
		// no need to select account
		return nil
	}

	accountsResult, err := as.WL.SelectedWallet.Wallet.GetAccountsRaw()
	if err != nil {
		return err
	}

	accounts := accountsResult.Acc
	for _, account := range accounts {
		if as.accountIsValid(account) {
			as.SetSelectedAccount(account)
			as.callback(account)
			return nil
		}
	}

	return errors.New(values.String(values.StrNoValidAccountFound))
}

func (as *AccountSelector) SetSelectedAccount(account *dcrlibwallet.Account) {
	as.selectedAccount = account
	as.totalBalance = dcrutil.Amount(account.TotalBalance).String()
}

func (as *AccountSelector) UpdateSelectedAccountBalance() {
	// wal := as.multiWallet.WalletWithID(as.SelectedAccount().WalletID)
	bal, err := as.WL.SelectedWallet.Wallet.GetAccountBalance(as.SelectedAccount().Number)
	if err == nil {
		as.totalBalance = dcrutil.Amount(bal.Total).String()
	}
}

func (as *AccountSelector) SelectedAccount() *dcrlibwallet.Account {
	return as.selectedAccount
}

func (as *AccountSelector) Layout(window app.WindowNavigator, gtx C) D {
	as.Handle(window)

	return decredmaterial.LinearLayout{
		Width:   decredmaterial.MatchParent,
		Height:  decredmaterial.WrapContent,
		Padding: layout.UniformInset(values.MarginPadding12),
		Border: decredmaterial.Border{
			Width:  values.MarginPadding2,
			Color:  as.Theme.Color.Gray2,
			Radius: decredmaterial.Radius(8),
		},
		Clickable: as.openSelectorDialog,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			accountIcon := as.Theme.Icons.AccountIcon
			inset := layout.Inset{
				Right: values.MarginPadding8,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return accountIcon.Layout24dp(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return as.Theme.Body1(as.selectedAccount.Name).Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return as.Theme.Body1(as.totalBalance).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Left: values.MarginPadding15,
						}
						return inset.Layout(gtx, func(gtx C) D {
							ic := decredmaterial.NewIcon(as.Theme.Icons.DropDownIcon)
							ic.Color = as.Theme.Color.Gray1
							return ic.Layout(gtx, values.MarginPadding20)
						})
					}),
				)
			})
		}),
	)
}

func (as *AccountSelector) ListenForTxNotifications(ctx context.Context, window app.WindowNavigator) {
	if as.TxAndBlockNotificationListener != nil {
		return
	}
	as.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := as.WL.MultiWallet.AddTxAndBlockNotificationListener(as.TxAndBlockNotificationListener, true, AccoutSelectorID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-as.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.BlockAttached:
					// refresh wallet account and balance on every new block
					// only if sync is completed.
					if as.WL.MultiWallet.IsSynced() {
						as.UpdateSelectedAccountBalance()
						if as.selectorModal != nil {
							as.selectorModal.setupWalletAccounts()
						}
						window.Reload()
					}
				case listeners.NewTransaction:
					// refresh accounts list when new transaction is received
					as.UpdateSelectedAccountBalance()
					if as.selectorModal != nil {
						as.selectorModal.setupWalletAccounts()
					}
					window.Reload()
				}
			case <-ctx.Done():
				as.WL.MultiWallet.RemoveTxAndBlockNotificationListener(AccoutSelectorID)
				close(as.TxAndBlockNotifChan)
				as.TxAndBlockNotificationListener = nil
				return
			}
		}
	}()
}

type AccountSelectorModal struct {
	*load.Load
	*decredmaterial.Modal

	accountIsValid func(*dcrlibwallet.Account) bool
	callback       func(*dcrlibwallet.Account)
	onExit         func()

	walletInfoButton decredmaterial.IconButton
	accountsList     layout.List

	currentSelectedAccount *dcrlibwallet.Account
	accounts               []*selectorAccount // key = wallet id
	eventQueue             event.Queue
	walletMu               sync.Mutex

	dialogTitle string

	isCancelable bool
}

type selectorAccount struct {
	*dcrlibwallet.Account
	clickable *decredmaterial.Clickable
}

func newAccountSelectorModal(l *load.Load, currentSelectedAccount *dcrlibwallet.Account) *AccountSelectorModal {
	asm := &AccountSelectorModal{
		Load:         l,
		Modal:        l.Theme.ModalFloatTitle("AccountSelectorModal"),
		accountsList: layout.List{Axis: layout.Vertical},

		currentSelectedAccount: currentSelectedAccount,
		isCancelable:           true,
	}

	asm.walletInfoButton = l.Theme.IconButton(asm.Theme.Icons.ActionInfo)
	asm.walletInfoButton.Size = values.MarginPadding15
	asm.walletInfoButton.Inset = layout.UniformInset(values.MarginPadding0)

	asm.Modal.ShowScrollbar(true)
	return asm
}

func (asm *AccountSelectorModal) OnResume() {
	asm.setupWalletAccounts()
}

func (asm *AccountSelectorModal) setupWalletAccounts() {
	walletAccounts := make([]*selectorAccount, 0)

	if !asm.WL.SelectedWallet.Wallet.IsWatchingOnlyWallet() {
		accountsResult, err := asm.WL.SelectedWallet.Wallet.GetAccountsRaw()
		if err != nil {
			log.Errorf("Error getting accounts:", err)
			return
		}

		accounts := accountsResult.Acc
		walletAccounts = make([]*selectorAccount, 0)
		for _, account := range accounts {
			if asm.accountIsValid(account) {
				walletAccounts = append(walletAccounts, &selectorAccount{
					Account:   account,
					clickable: asm.Theme.NewClickable(true),
				})
			}
		}
	}
	asm.accounts = walletAccounts
}

func (asm *AccountSelectorModal) SetCancelable(min bool) *AccountSelectorModal {
	asm.isCancelable = min
	return asm
}

func (asm *AccountSelectorModal) Handle() {
	if asm.eventQueue != nil {
		for _, account := range asm.accounts {
			for account.clickable.Clicked() {
				asm.callback(account.Account)
				asm.onExit()
				asm.Dismiss()
			}
		}
	}

	if asm.Modal.BackdropClicked(asm.isCancelable) {
		asm.onExit()
		asm.Dismiss()
	}
}

func (asm *AccountSelectorModal) title(title string) *AccountSelectorModal {
	asm.dialogTitle = title
	return asm
}

func (asm *AccountSelectorModal) accountValidator(accountIsValid func(*dcrlibwallet.Account) bool) *AccountSelectorModal {
	asm.accountIsValid = accountIsValid
	return asm
}

func (asm *AccountSelectorModal) accountSelected(callback func(*dcrlibwallet.Account)) *AccountSelectorModal {
	asm.callback = callback
	return asm
}

func (asm *AccountSelectorModal) Layout(gtx C) D {
	asm.eventQueue = gtx

	w := []layout.Widget{
		func(gtx C) D {
			title := asm.Theme.H6(asm.dialogTitle)
			title.Color = asm.Theme.Color.Text
			title.Font.Weight = text.SemiBold
			return title.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Stack{Alignment: layout.NW}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					accounts := asm.accounts
					return asm.accountsList.Layout(gtx, len(accounts), func(gtx C, aindex int) D {
						return asm.walletAccountLayout(gtx, accounts[aindex])
					})
				}),
				layout.Stacked(func(gtx C) D {
					if false { //TODO
						inset := layout.Inset{
							Top:  values.MarginPadding20,
							Left: values.MarginPaddingMinus75,
						}
						return inset.Layout(gtx, func(gtx C) D {
							// return page.walletInfoPopup(gtx)
							return D{}
						})
					}
					return D{}
				}),
			)
		},
	}

	return asm.Modal.Layout(gtx, w)
}

func (asm *AccountSelectorModal) walletAccountLayout(gtx C, account *selectorAccount) D {
	accountIcon := asm.Theme.Icons.AccountIcon

	return decredmaterial.LinearLayout{
		Width:     decredmaterial.MatchParent,
		Height:    decredmaterial.WrapContent,
		Margin:    layout.Inset{Bottom: values.MarginPadding4},
		Padding:   layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8},
		Clickable: account.clickable,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(0.1, func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding18,
			}.Layout(gtx, func(gtx C) D {
				return accountIcon.Layout24dp(gtx)
			})
		}),
		layout.Flexed(0.8, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					acct := asm.Theme.Label(values.TextSize18, account.Name)
					acct.Color = asm.Theme.Color.Text
					return EndToEndRow(gtx, acct.Layout, func(gtx C) D {
						return LayoutBalance(gtx, asm.Load, dcrutil.Amount(account.TotalBalance).String())
					})
				}),
				layout.Rigid(func(gtx C) D {
					spendable := asm.Theme.Label(values.TextSize14, values.String(values.StrLabelSpendable))
					spendable.Color = asm.Theme.Color.GrayText2
					spendableBal := asm.Theme.Label(values.TextSize14, dcrutil.Amount(account.Balance.Spendable).String())
					spendableBal.Color = asm.Theme.Color.GrayText2
					return EndToEndRow(gtx, spendable.Layout, spendableBal.Layout)
				}),
			)
		}),

		layout.Flexed(0.1, func(gtx C) D {
			inset := layout.Inset{
				Right: values.MarginPadding2,
				Top:   values.MarginPadding10,
				Left:  values.MarginPadding10,
			}
			sections := func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						ic := decredmaterial.NewIcon(asm.Theme.Icons.NavigationCheck)
						ic.Color = asm.Theme.Color.Gray1
						return ic.Layout(gtx, values.MarginPadding20)
					})
				})
			}

			if account.Number == asm.currentSelectedAccount.Number {
				return sections(gtx)
			}
			return D{}
		}),
	)
}

func (asm *AccountSelectorModal) walletInfoPopup(gtx C) D {
	// TODO: currently not used.. skipping str localization
	title := "Some accounts are hidden."
	desc := "Some accounts are disabled by StakeShuffle settings to protect your privacy."
	card := asm.Theme.Card()
	card.Radius = decredmaterial.Radius(7)
	gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding280)
	return card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := asm.Theme.Body2(title)
							txt.Font.Weight = text.SemiBold
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := asm.Theme.Body2("Tx direction")
							txt.Color = asm.Theme.Color.GrayText2
							return txt.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					txt := asm.Theme.Body2(desc)
					txt.Color = asm.Theme.Color.GrayText2
					return txt.Layout(gtx)
				}),
			)
		})
	})
}

func (asm *AccountSelectorModal) onModalExit(exit func()) *AccountSelectorModal {
	asm.onExit = exit
	return asm
}

func (asm *AccountSelectorModal) OnDismiss() {}
