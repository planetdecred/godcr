package governance

import (
	"errors"
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
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

	openSelectorDialog *decredmaterial.Clickable

	wallets        []*dcrlibwallet.Wallet
	selectedWallet *dcrlibwallet.Wallet
	totalBalance   string
}

// TODO: merge this into the account selector modal.
func NewWalletSelector(l *load.Load) *WalletSelector {

	return &WalletSelector{
		Load:               l,
		multiWallet:        l.WL.MultiWallet,
		walletIsValid:      func(*dcrlibwallet.Wallet) bool { return true },
		openSelectorDialog: l.Theme.NewClickable(true),

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

func (as *WalletSelector) Handle(window app.WindowNavigator) {
	for as.openSelectorDialog.Clicked() {
		walletSelectorModal := newWalletSelectorModal(as.Load, as.selectedWallet).
			title(as.dialogTitle).
			accountValidator(as.walletIsValid).
			accountSelected(func(wallet *dcrlibwallet.Wallet) {
				as.selectedWallet = wallet
				as.setupSelectedWallet(wallet)
				as.callback(wallet)
			})
		window.ShowModal(walletSelectorModal)
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

	return errors.New(values.String(values.StrnoValidWalletFound))
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

func (as *WalletSelector) Layout(gtx layout.Context, window app.WindowNavigator) layout.Dimensions {
	as.Handle(window)

	border := widget.Border{
		Color:        as.Theme.Color.Gray2,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return as.openSelectorDialog.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						accountIcon := as.Theme.Icons.AccountIcon
						return layout.Inset{
							Right: values.MarginPadding8,
						}.Layout(gtx, accountIcon.Layout24dp)
					}),
					layout.Rigid(as.Theme.Body1(as.selectedWallet.Name).Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(as.Theme.Body1(as.totalBalance).Layout),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										ic := decredmaterial.NewIcon(as.Theme.Icons.DropDownIcon)
										return ic.Layout(gtx, values.MarginPadding20)
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

type WalletSelectorModal struct {
	*load.Load
	*decredmaterial.Modal

	dialogTitle string

	isCancelable bool

	walletIsValid func(*dcrlibwallet.Wallet) bool
	callback      func(*dcrlibwallet.Wallet)

	walletsList *decredmaterial.ClickableList

	currentSelectedWallet *dcrlibwallet.Wallet
	filteredWallets       []*dcrlibwallet.Wallet
}

func newWalletSelectorModal(l *load.Load, currentSelectedAccount *dcrlibwallet.Wallet) *WalletSelectorModal {
	asm := &WalletSelectorModal{
		Load:        l,
		Modal:       l.Theme.ModalFloatTitle("WalletSelectorModal"),
		walletsList: l.Theme.NewClickableList(layout.Vertical),

		currentSelectedWallet: currentSelectedAccount,
		isCancelable:          true,
	}

	return asm
}

func (asm *WalletSelectorModal) OnResume() {
	wallets := make([]*dcrlibwallet.Wallet, 0)

	for _, wal := range asm.WL.SortedWalletList() {
		if asm.walletIsValid(wal) {
			wallets = append(wallets, wal)
		}
	}

	asm.filteredWallets = wallets
}

func (asm *WalletSelectorModal) Handle() {
	if clicked, index := asm.walletsList.ItemClicked(); clicked {
		asm.callback(asm.filteredWallets[index])
		asm.Dismiss()
	}

	if asm.Modal.BackdropClicked(asm.isCancelable) {
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
			title.Font.Weight = text.SemiBold
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

	return asm.Modal.Layout(gtx, w)
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
							accountIcon := asm.Theme.Icons.AccountIcon
							return accountIcon.Layout24dp(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								acct := asm.Theme.Label(values.TextSize18, wallet.Name)
								acct.Color = asm.Theme.Color.Text
								return components.EndToEndRow(gtx, acct.Layout, func(gtx C) D {
									return components.LayoutBalanceWithUnit(gtx, asm.Load, walletTotalBalance.String())
								})
							}),
							layout.Rigid(func(gtx C) D {
								spendable := asm.Theme.Label(values.TextSize14, values.String(values.StrLabelSpendable))
								spendable.Color = asm.Theme.Color.GrayText2
								//TODO
								spendableBal := asm.Theme.Label(values.TextSize14, walletSpendableBalance.String())
								spendableBal.Color = asm.Theme.Color.GrayText2
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
									ic := decredmaterial.NewIcon(asm.Theme.Icons.NavigationCheck)
									return ic.Layout(gtx, values.MarginPadding20)
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
