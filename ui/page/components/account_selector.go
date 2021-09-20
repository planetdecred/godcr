package components

import (
	"errors"
	"fmt"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type AccountSelector struct {
	*load.Load
	multiWallet *dcrlibwallet.MultiWallet
	dialogTitle string

	accountIsValid func(*dcrlibwallet.Account) bool
	callback       func(*dcrlibwallet.Account)

	openSelectorDialog *widget.Clickable

	wallets            []*dcrlibwallet.Wallet
	selectedAccount    *dcrlibwallet.Account
	selectedWalletName string
	totalBalance       string
}

func NewAccountSelector(l *load.Load) *AccountSelector {

	return &AccountSelector{
		Load:               l,
		multiWallet:        l.WL.MultiWallet,
		accountIsValid:     func(*dcrlibwallet.Account) bool { return true },
		openSelectorDialog: new(widget.Clickable),

		wallets: l.WL.SortedWalletList(),
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

func (as *AccountSelector) Handle() {
	for as.openSelectorDialog.Clicked() {
		newAccountSelectorModal(as.Load, as.selectedAccount, as.wallets).
			title(as.dialogTitle).
			accountValidator(as.accountIsValid).
			accountSelected(func(account *dcrlibwallet.Account) {
				as.setupSelectedAccount(account)
				as.callback(account)
			}).Show()
	}
}

func (as *AccountSelector) SelectFirstWalletValidAccount() error {
	if as.selectedAccount != nil && as.accountIsValid(as.selectedAccount) {
		// no need to select account
		return nil
	}

	for _, wal := range as.wallets {
		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			return err
		}

		accounts := accountsResult.Acc
		for _, account := range accounts {
			if as.accountIsValid(account) {
				as.setupSelectedAccount(account)
				as.callback(account)
				return nil
			}
		}
	}

	return errors.New("no valid account found")
}

func (as *AccountSelector) setupSelectedAccount(account *dcrlibwallet.Account) {
	wal := as.multiWallet.WalletWithID(account.WalletID)

	as.selectedAccount = account
	as.selectedWalletName = wal.Name
	as.totalBalance = dcrutil.Amount(account.TotalBalance).String()
}

func (as *AccountSelector) SelectedAccount() *dcrlibwallet.Account {
	return as.selectedAccount
}

func (as *AccountSelector) Layout(gtx layout.Context) layout.Dimensions {
	as.Handle()

	border := widget.Border{
		Color:        as.Theme.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, as.openSelectorDialog, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						accountIcon := as.Icons.AccountIcon
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
								Color: as.Theme.Color.LightGray,
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
									text.Color = as.Theme.Color.Gray
									return text.Layout(gtx)
								})
							})
						})
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := as.Theme.Body1(as.totalBalance)
									txt.Color = as.Theme.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return as.Icons.DropDownIcon.Layout(gtx, values.MarginPadding20)
									})
								}),
							)
						})
					}),
				)
			})
		})
	})
}

const ModalAccountSelector = "AccountSelectorModal"

type AccountSelectorModal struct {
	*load.Load
	dialogTitle string

	isCancelable bool

	accountIsValid func(*dcrlibwallet.Account) bool
	callback       func(*dcrlibwallet.Account)

	modal            decredmaterial.Modal
	walletInfoButton decredmaterial.IconButton
	walletsList      layout.List
	accountsList     layout.List

	currentSelectedAccount *dcrlibwallet.Account
	wallets                []*dcrlibwallet.Wallet // TODO sort array instead
	filteredWallets        []*dcrlibwallet.Wallet
	accounts               map[int][]*selectorAccount // key = wallet id
	eventQueue             event.Queue
}

type selectorAccount struct {
	*dcrlibwallet.Account
	button *widget.Clickable
}

func newAccountSelectorModal(l *load.Load, currentSelectedAccount *dcrlibwallet.Account, wallets []*dcrlibwallet.Wallet) *AccountSelectorModal {
	asm := &AccountSelectorModal{
		Load:         l,
		modal:        *l.Theme.ModalFloatTitle(),
		walletsList:  layout.List{Axis: layout.Vertical},
		accountsList: layout.List{Axis: layout.Vertical},

		currentSelectedAccount: currentSelectedAccount,
		wallets:                wallets,
		isCancelable:           true,
	}

	asm.walletInfoButton = l.Theme.PlainIconButton(new(widget.Clickable), asm.Icons.ActionInfo)
	asm.walletInfoButton.Color = asm.Theme.Color.Gray3
	asm.walletInfoButton.Size = values.MarginPadding15
	asm.walletInfoButton.Inset = layout.UniformInset(values.MarginPadding0)

	return asm
}

func (asm *AccountSelectorModal) OnResume() {
	wallets := make([]*dcrlibwallet.Wallet, 0)
	walletAccounts := make(map[int][]*selectorAccount)

	// TODO use a sorted wallet list
	for _, wal := range asm.wallets {
		// filter all accounts
		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			fmt.Println("Error getting accounts:", err)
			return
		}

		accounts := accountsResult.Acc
		walletAccounts[wal.ID] = make([]*selectorAccount, 0)
		for _, account := range accounts {
			if asm.accountIsValid(account) {
				walletAccounts[wal.ID] = append(walletAccounts[wal.ID], &selectorAccount{
					Account: account,
					button:  new(widget.Clickable),
				})
			}
		}

		// add wallet only if it has valid accounts
		if len(walletAccounts[wal.ID]) > 0 {
			wallets = append(wallets, wal)
		}
	}

	asm.filteredWallets = wallets
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
				if account.button.Clicked() {
					asm.callback(account.Account)
					asm.Dismiss()
				}
			}
		}
	}

	if asm.modal.BackdropClicked(asm.isCancelable) {
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

func (asm *AccountSelectorModal) OnDismiss() {

}

func (asm *AccountSelectorModal) Layout(gtx layout.Context) layout.Dimensions {
	asm.eventQueue = gtx

	wallAcctGroup := func(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
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
							var showInfoBtn bool = false
							if showInfoBtn {
								inset := layout.Inset{
									Top: values.MarginPadding2,
								}
								return inset.Layout(gtx, func(gtx C) D {
									return asm.walletInfoButton.Layout(gtx)
								})
							}
							return layout.Dimensions{}
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
			title.Font.Weight = text.Bold
			return title.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Stack{Alignment: layout.NW}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return asm.walletsList.Layout(gtx, len(asm.filteredWallets), func(gtx C, windex int) D {
						wal := asm.filteredWallets[windex]
						return wallAcctGroup(gtx, wal.Name, func(gtx C) D {
							accounts := asm.accounts[wal.ID]

							card := asm.Theme.Card()
							card.Color = color.NRGBA{}
							card.Radius = decredmaterial.Radius(0)
							return asm.accountsList.Layout(gtx, len(accounts), func(gtx C, aindex int) D {
								return card.HoverableLayout(gtx, accounts[aindex].button, func(gtx C) D {
									return asm.walletAccountLayout(gtx, accounts[aindex])
								})
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
							return layout.Dimensions{}
						})
					}
					return layout.Dimensions{}
				}),
			)
		},
	}

	return asm.modal.Layout(gtx, w, 850)
}

func (asm *AccountSelectorModal) walletAccountLayout(gtx layout.Context, account *selectorAccount) layout.Dimensions {
	accountIcon := asm.Icons.AccountIcon
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return Container{Padding: layout.UniformInset(values.MarginPadding8)}.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(0.1, func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding18,
								}.Layout(gtx, func(gtx C) D {
									return accountIcon.Layout24dp(gtx)
								})
							}),
							layout.Flexed(0.7, func(gtx C) D {
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
										spendable.Color = asm.Theme.Color.Gray
										spendableBal := asm.Theme.Label(values.TextSize14, dcrutil.Amount(account.Balance.Spendable).String())
										spendableBal.Color = asm.Theme.Color.Gray
										return EndToEndRow(gtx, spendable.Layout, spendableBal.Layout)
									}),
								)
							}),
							layout.Flexed(0.1, func(gtx C) D {
								inset := layout.Inset{
									Top: values.MarginPadding10,
								}
								sections := func(gtx layout.Context) layout.Dimensions {
									return layout.E.Layout(gtx, func(gtx C) D {
										return inset.Layout(gtx, func(gtx C) D {
											return asm.Icons.NavigationCheck.Layout(gtx, values.MarginPadding20)
										})
									})
								}

								if account.Number == asm.currentSelectedAccount.Number &&
									account.WalletID == asm.currentSelectedAccount.WalletID {
									return sections(gtx)
								}
								return layout.Dimensions{}
							}),
						)
					}),
				)
			}),
			layout.Expanded(account.button.Layout),
		)
	})
}

func (asm *AccountSelectorModal) walletInfoPopup(gtx layout.Context) layout.Dimensions {
	title := "Some accounts are hidden."
	desc := "Some accounts are disabled by StakeShuffle settings to protect your privacy."
	card := asm.Theme.Card()
	card.Radius = decredmaterial.Radius(7)
	border := widget.Border{Color: asm.Theme.Color.Background, CornerRadius: values.MarginPadding7, Width: values.MarginPadding1}
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding280)
	return border.Layout(gtx, func(gtx C) D {
		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := asm.Theme.Body2(title)
								txt.Color = asm.Theme.Color.DeepBlue
								txt.Font.Weight = text.Bold
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								txt := asm.Theme.Body2("Tx direction")
								txt.Color = asm.Theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						txt := asm.Theme.Body2(desc)
						txt.Color = asm.Theme.Color.Gray
						return txt.Layout(gtx)
					}),
				)
			})
		})
	})
}
