package components

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const AccoutSelectorID = "AccountSelector"

type AccountSelector struct {
	*load.Load
	*listeners.TxAndBlockNotificationListener

	multiWallet     *dcrlibwallet.MultiWallet
	selectedAccount *dcrlibwallet.Account
	selectedWallet  *dcrlibwallet.Wallet

	accountIsValid func(*dcrlibwallet.Account) bool
	callback       func(*dcrlibwallet.Account)

	openSelectorDialog *decredmaterial.Clickable
	selectorModal      *AccountSelectorModal

	dialogTitle        string
	selectedWalletName string
	totalBalance       string
	changed            bool
}

// NewAccountSelector opens up a modal to select the desired account. If a
// nil value is passed for selectedWallet, then accounts for all wallets
// are shown, otherwise only accounts for the selectedWallet is shown.
func NewAccountSelector(l *load.Load, selectedWallet *dcrlibwallet.Wallet) *AccountSelector {
	return &AccountSelector{
		Load:               l,
		multiWallet:        l.WL.MultiWallet,
		selectedWallet:     selectedWallet,
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

func (as *AccountSelector) Handle() {
	for as.openSelectorDialog.Clicked() {
		as.selectorModal = newAccountSelectorModal(as.Load, as.selectedAccount, as.selectedWallet).
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
		as.ShowModal(as.selectorModal)
	}
}

// SelectFirstWalletValidAccount selects the first valid account from the
// first wallet in the SortedWalletList. If selectedWallet is not nil,
// the first account for the selectWallet is selected.
func (as *AccountSelector) SelectFirstWalletValidAccount(selectedWallet *dcrlibwallet.Wallet) error {
	if as.selectedAccount != nil && as.accountIsValid(as.selectedAccount) {
		as.UpdateSelectedAccountBalance()
		// no need to select account
		return nil
	}

	if selectedWallet != nil {
		accountsResult, err := selectedWallet.GetAccountsRaw()
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
	}

	for _, wal := range as.WL.SortedWalletList() {
		accountsResult, err := wal.GetAccountsRaw()
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
	}

	return errors.New("no valid account found")
}

// SelectValidAccountExcept selects a valid account from the selectedWallet
// except the account with accountID.
func (as *AccountSelector) SelectValidAccountExcept(selectedWallet *dcrlibwallet.Wallet, accountID int32) error {
	if selectedWallet != nil {
		accountsResult, err := selectedWallet.GetAccountsRaw()
		if err != nil {
			return err
		}

		accounts := accountsResult.Acc
		for _, account := range accounts {
			if as.accountIsValid(account) && account.Number != accountID {
				as.SetSelectedAccount(account)
				as.callback(account)
				return nil
			}
		}
	}

	return errors.New("no valid account found")
}

func (as *AccountSelector) SetSelectedAccount(account *dcrlibwallet.Account) {
	wal := as.multiWallet.WalletWithID(account.WalletID)

	as.selectedAccount = account
	as.selectedWalletName = wal.Name
	as.totalBalance = dcrutil.Amount(account.TotalBalance).String()
}

func (as *AccountSelector) UpdateSelectedAccountBalance() {
	wal := as.multiWallet.WalletWithID(as.SelectedAccount().WalletID)
	bal, err := wal.GetAccountBalance(as.SelectedAccount().Number)
	if err == nil {
		as.totalBalance = dcrutil.Amount(bal.Total).String()
	}
}

func (as *AccountSelector) SelectedAccount() *dcrlibwallet.Account {
	return as.selectedAccount
}

func (as *AccountSelector) Layout(gtx C) D {
	as.Handle()

	return decredmaterial.LinearLayout{
		Width:     decredmaterial.MatchParent,
		Height:    decredmaterial.WrapContent,
		Padding:   layout.UniformInset(values.MarginPadding12),
		Border:    decredmaterial.Border{Width: values.MarginPadding2, Color: as.Theme.Color.Gray2, Radius: decredmaterial.Radius(8)},
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
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Left: values.MarginPadding4,
				Top:  values.MarginPadding2,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return decredmaterial.Card{
					Color: as.Theme.Color.Gray4,
				}.Layout(gtx, func(gtx C) D {
					m2 := values.MarginPadding2
					m4 := values.MarginPadding4
					inset := layout.Inset{
						Left:   m4,
						Top:    m2,
						Bottom: m2,
						Right:  m4,
					}
					return inset.Layout(gtx, func(gtx C) D {
						text := as.Theme.Caption(as.selectedWalletName)
						text.Color = as.Theme.Color.GrayText2
						return text.Layout(gtx)
					})
				})
			})
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

func (as *AccountSelector) ListenForTxNotifications(ctx context.Context) {
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
						as.RefreshWindow()
					}
				case listeners.NewTransaction:
					// refresh accounts list when new transaction is received
					as.UpdateSelectedAccountBalance()
					if as.selectorModal != nil {
						as.selectorModal.setupWalletAccounts()
					}
					as.RefreshWindow()
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

const ModalAccountSelector = "AccountSelectorModal"

type AccountSelectorModal struct {
	*load.Load

	accountIsValid func(*dcrlibwallet.Account) bool
	callback       func(*dcrlibwallet.Account)
	onExit         func()

	modal            decredmaterial.Modal
	walletInfoButton decredmaterial.IconButton
	walletsList      layout.List
	accountsList     layout.List

	currentSelectedAccount *dcrlibwallet.Account
	selectedWallet         *dcrlibwallet.Wallet
	filteredWallet         []*dcrlibwallet.Wallet
	accounts               map[int][]*selectorAccount // key = wallet id
	eventQueue             event.Queue
	walletMu               sync.Mutex

	dialogTitle string

	isCancelable bool
}

type selectorAccount struct {
	*dcrlibwallet.Account
	clickable *decredmaterial.Clickable
}

func newAccountSelectorModal(l *load.Load, currentSelectedAccount *dcrlibwallet.Account, selectedWallet *dcrlibwallet.Wallet) *AccountSelectorModal {
	asm := &AccountSelectorModal{
		Load:         l,
		modal:        *l.Theme.ModalFloatTitle(),
		walletsList:  layout.List{Axis: layout.Vertical},
		accountsList: layout.List{Axis: layout.Vertical},

		currentSelectedAccount: currentSelectedAccount,
		selectedWallet:         selectedWallet,
		isCancelable:           true,
	}

	asm.walletInfoButton = l.Theme.IconButton(asm.Theme.Icons.ActionInfo)
	asm.walletInfoButton.Size = values.MarginPadding15
	asm.walletInfoButton.Inset = layout.UniformInset(values.MarginPadding0)

	asm.modal.ShowScrollbar(true)
	return asm
}

func (asm *AccountSelectorModal) OnResume() {
	asm.setupWalletAccounts()
}

func (asm *AccountSelectorModal) setupWalletAccounts() {
	wals := make([]*dcrlibwallet.Wallet, 0)
	walletAccounts := make(map[int][]*selectorAccount)

	for _, wal := range asm.WL.SortedWalletList() {
		if wal.IsWatchingOnlyWallet() {
			continue
		}

		if asm.selectedWallet == nil {
			wals = append(wals, wal)

			accountsResult, err := wal.GetAccountsRaw()
			if err != nil {
				fmt.Println("Error getting accounts:", err)
				continue
			}

			accounts := accountsResult.Acc
			walletAccounts[wal.ID] = make([]*selectorAccount, 0)
			for _, account := range accounts {
				if asm.accountIsValid(account) {
					walletAccounts[wal.ID] = append(walletAccounts[wal.ID], &selectorAccount{
						Account:   account,
						clickable: asm.Theme.NewClickable(true),
					})
				}
			}
		} else if wal.ID == asm.selectedWallet.ID {
			accountsResult, err := wal.GetAccountsRaw()
			if err != nil {
				fmt.Println("Error getting accounts:", err)
				continue
			}

			accounts := accountsResult.Acc
			walletAccounts[wal.ID] = make([]*selectorAccount, 0)
			for _, account := range accounts {
				if asm.accountIsValid(account) {
					walletAccounts[wal.ID] = append(walletAccounts[wal.ID], &selectorAccount{
						Account:   account,
						clickable: asm.Theme.NewClickable(true),
					})
				}
			}
		}
	}
	asm.filteredWallet = wals
	asm.accounts = walletAccounts
}

func (asm *AccountSelectorModal) ModalID() string {
	return ModalAccountSelector
}

func (asm *AccountSelectorModal) Show() {
	asm.ShowModal(asm)
}

func (asm *AccountSelectorModal) Dismiss() {
	asm.DismissModal(asm)
}

func (asm *AccountSelectorModal) SetCancelable(min bool) *AccountSelectorModal {
	asm.isCancelable = min
	return asm
}

func (asm *AccountSelectorModal) Handle() {
	if asm.eventQueue != nil {
		for _, accounts := range asm.accounts {
			for _, account := range accounts {
				for account.clickable.Clicked() {
					asm.callback(account.Account)
					asm.onExit()
					asm.Dismiss()
				}
			}
		}
	}

	if asm.modal.BackdropClicked(asm.isCancelable) {
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

	asm.walletMu.Lock()
	wallets := asm.filteredWallet
	asm.walletMu.Unlock()

	wallAcctGroup := func(gtx C, title string, body layout.Widget) D {
		return layout.Inset{
			Bottom: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := asm.Theme.Body2(title)
							txt.Color = asm.Theme.Color.Text
							inset := layout.Inset{
								Bottom: values.MarginPadding15,
							}
							return inset.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							showInfoBtn := false
							if showInfoBtn {
								inset := layout.Inset{
									Top: values.MarginPadding2,
								}
								return inset.Layout(gtx, func(gtx C) D {
									return asm.walletInfoButton.Layout(gtx)
								})
							}
							return D{}
						}),
					)
				}),
				layout.Rigid(body),
			)
		})
	}

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
					if asm.selectedWallet != nil {
						return asm.walletsList.Layout(gtx, 1, func(gtx C, windex int) D {
							wal := asm.selectedWallet
							return wallAcctGroup(gtx, wal.Name, func(gtx C) D {
								accounts := asm.accounts[wal.ID]
								return asm.accountsList.Layout(gtx, len(accounts), func(gtx C, aindex int) D {
									return asm.walletAccountLayout(gtx, accounts[aindex])
								})
							})
						})
					}

					return asm.walletsList.Layout(gtx, len(wallets), func(gtx C, windex int) D {
						wal := wallets[windex]
						return wallAcctGroup(gtx, wal.Name, func(gtx C) D {
							accounts := asm.accounts[wal.ID]
							return asm.accountsList.Layout(gtx, len(accounts), func(gtx C, aindex int) D {
								return asm.walletAccountLayout(gtx, accounts[aindex])
							})
						})
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

	return asm.modal.Layout(gtx, w)
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

			if account.Number == asm.currentSelectedAccount.Number &&
				account.WalletID == asm.currentSelectedAccount.WalletID {
				return sections(gtx)
			}
			return D{}
		}),
	)
}

func (asm *AccountSelectorModal) walletInfoPopup(gtx C) D {
	title := "Some accounts are hidden."
	desc := "Some accounts are disabled by StakeShuffle settings to protect your privacy."
	card := asm.Theme.Card()
	card.Radius = decredmaterial.Radius(7)
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding280)
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
