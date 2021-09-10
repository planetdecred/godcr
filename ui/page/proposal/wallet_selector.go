package proposal

import (
	"errors"
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type WalletSelector struct {
	*load.Load
	multiWallet *dcrlibwallet.MultiWallet
	dialogTitle string

	walletIsValid func(*dcrlibwallet.Wallet) bool
	callback      func(*dcrlibwallet.Wallet)

	openSelectorDialog *widget.Clickable

	wallets        []*dcrlibwallet.Wallet
	selectedWallet *dcrlibwallet.Wallet
	totalBalance   string
}

func NewWalletSelector(l *load.Load) *WalletSelector {

	return &WalletSelector{
		Load:               l,
		multiWallet:        l.WL.MultiWallet,
		walletIsValid:      func(*dcrlibwallet.Wallet) bool { return true },
		openSelectorDialog: new(widget.Clickable),

		wallets: l.WL.SortedWalletList(),
	}
}

func (as *WalletSelector) Title(title string) *WalletSelector {
	as.dialogTitle = title
	return as
}

func (as *WalletSelector) WalletValidator(walletIsValid func(*dcrlibwallet.Wallet) bool) *WalletSelector {
	as.walletIsValid = walletIsValid
	return as
}

func (as *WalletSelector) WalletSelected(callback func(*dcrlibwallet.Wallet)) *WalletSelector {
	as.callback = callback
	return as
}

func (as *WalletSelector) Handle() {
	for as.openSelectorDialog.Clicked() {
		newWalletSelectorModal(as.Load, as.selectedWallet, as.wallets).
			title(as.dialogTitle).
			accountValidator(as.walletIsValid).
			accountSelected(func(wallet *dcrlibwallet.Wallet) {
				as.selectedWallet = wallet
				as.setupSelectedWallet(wallet)
				as.callback(wallet)
			}).Show()

	}
}

func (as *WalletSelector) SelectFirstValidWallet() error {
	if as.selectedWallet != nil && as.walletIsValid(as.selectedWallet) {
		// no need to select account
		return nil
	}

	for _, wal := range as.wallets {
		if as.walletIsValid(wal) {
			as.selectedWallet = wal
			as.setupSelectedWallet(wal)
			as.callback(wal)
			return nil
		}
	}

	return errors.New("no valid wallet found")
}

func (as *WalletSelector) setupSelectedWallet(wallet *dcrlibwallet.Wallet) {

	totalBalance, err := as.WL.TotalWalletBalance(wallet.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	as.totalBalance = totalBalance.String()
}

func (as *WalletSelector) SelectedWallet() *dcrlibwallet.Wallet {
	return as.selectedWallet
}

func (as *WalletSelector) Layout(gtx layout.Context) layout.Dimensions {
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
						return as.Theme.Body1(as.selectedWallet.Name).Layout(gtx)
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

const ModalWalletSelector = "WalletSelectorModal"

type WalletSelectorModal struct {
	*load.Load
	dialogTitle string

	walletIsValid func(*dcrlibwallet.Wallet) bool
	callback      func(*dcrlibwallet.Wallet)

	modal       decredmaterial.Modal
	walletsList *decredmaterial.ClickableList

	currentSelectedWallet *dcrlibwallet.Wallet
	wallets               []*dcrlibwallet.Wallet // TODO sort array instead
	filteredWallets       []*dcrlibwallet.Wallet
}

func newWalletSelectorModal(l *load.Load, currentSelectedAccount *dcrlibwallet.Wallet, wallets []*dcrlibwallet.Wallet) *WalletSelectorModal {
	asm := &WalletSelectorModal{
		Load:        l,
		modal:       *l.Theme.ModalFloatTitle(),
		walletsList: l.Theme.NewClickableList(layout.Vertical),

		currentSelectedWallet: currentSelectedAccount,
		wallets:               wallets,
	}

	return asm
}

func (asm *WalletSelectorModal) OnResume() {
	wallets := make([]*dcrlibwallet.Wallet, 0)

	for _, wal := range asm.wallets {
		if asm.walletIsValid(wal) {
			wallets = append(wallets, wal)
		}
	}

	asm.filteredWallets = wallets
}

func (asm *WalletSelectorModal) ModalID() string {
	return ModalWalletSelector
}

func (asm *WalletSelectorModal) Show() {
	asm.ShowModal(asm)
}

func (asm *WalletSelectorModal) Dismiss() {
	asm.DismissModal(asm)
}

func (asm *WalletSelectorModal) Handle() {
	if clicked, index := asm.walletsList.ItemClicked(); clicked {
		asm.callback(asm.filteredWallets[index])
		asm.Dismiss()
	}
}

func (asm *WalletSelectorModal) title(title string) *WalletSelectorModal {
	asm.dialogTitle = title
	return asm
}

func (asm *WalletSelectorModal) accountValidator(walletIsValid func(*dcrlibwallet.Wallet) bool) *WalletSelectorModal {
	asm.walletIsValid = walletIsValid
	return asm
}

func (asm *WalletSelectorModal) accountSelected(callback func(*dcrlibwallet.Wallet)) *WalletSelectorModal {
	asm.callback = callback
	return asm
}

func (asm *WalletSelectorModal) OnDismiss() {

}

func (asm *WalletSelectorModal) Layout(gtx layout.Context) layout.Dimensions {

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
						return asm.walletAccountLayout(gtx, wal)
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

func (asm *WalletSelectorModal) walletAccountLayout(gtx layout.Context, wallet *dcrlibwallet.Wallet) layout.Dimensions {

	walletTotalBalance, _ := asm.WL.TotalWalletBalance(wallet.ID)
	walletSpendableBalance, _ := asm.WL.SpendableWalletBalance(wallet.ID)

	return layout.Inset{
		Bottom: values.MarginPadding20,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(0.1, func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding18,
						}.Layout(gtx, func(gtx C) D {
							accountIcon := asm.Icons.AccountIcon
							return accountIcon.Layout24dp(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								acct := asm.Theme.Label(values.TextSize18, wallet.Name)
								acct.Color = asm.Theme.Color.Text
								return components.EndToEndRow(gtx, acct.Layout, func(gtx C) D {
									return components.LayoutBalance(gtx, asm.Load, walletTotalBalance.String())
								})
							}),
							layout.Rigid(func(gtx C) D {
								spendable := asm.Theme.Label(values.TextSize14, values.String(values.StrLabelSpendable))
								spendable.Color = asm.Theme.Color.Gray
								//TODO
								spendableBal := asm.Theme.Label(values.TextSize14, walletSpendableBalance.String())
								spendableBal.Color = asm.Theme.Color.Gray
								return components.EndToEndRow(gtx, spendable.Layout, spendableBal.Layout)
							}),
						)
					}),

					layout.Flexed(0.1, func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding10,
						}
						sections := func(gtx layout.Context) layout.Dimensions {
							return layout.E.Layout(gtx, func(gtx C) D {
								return inset.Layout(gtx, func(gtx C) D {
									return asm.Icons.NavigationCheck.Layout(gtx, values.MarginPadding20)
								})
							})
						}

						if wallet.ID == asm.currentSelectedWallet.ID {
							return sections(gtx)
						}
						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}
