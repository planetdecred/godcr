package send

import (
	"fmt"
	"strings"

	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
)

type destination struct {
	*load.Load

	addressChanged             func()
	destinationAddressEditor   decredmaterial.Editor
	destinationAccountSelector *components.AccountSelector

	sendToAddress bool
	accountSwitch *decredmaterial.SwitchButtonText
}

func newSendDestination(l *load.Load) *destination {
	dst := &destination{
		Load: l,
	}

	dst.destinationAddressEditor = l.Theme.Editor(new(widget.Editor), "Address")
	dst.destinationAddressEditor.Editor.SingleLine = true
	dst.destinationAddressEditor.Editor.SetText("")

	dst.accountSwitch = l.Theme.SwitchButtonText([]decredmaterial.SwitchItem{{Text: "Address"}, {Text: "My account"}})

	// Destination account picker
	dst.destinationAccountSelector = components.NewAccountSelector(dst.Load).
		Title("Receiving account").
		AccountValidator(func(account *dcrlibwallet.Account) bool {

			// Filter out imported account and mixed.
			wal := dst.Load.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == components.MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}

			return true
		})

	return dst
}

func (dst *destination) destinationAddress() (string, error) {
	if dst.sendToAddress {
		valid, address := dst.validateDestinationAddress()
		if valid {
			return address, nil
		}

		return "", fmt.Errorf("invalid address")
	}

	destinationAccount := dst.destinationAccountSelector.SelectedAccount()
	wal := dst.WL.MultiWallet.WalletWithID(destinationAccount.WalletID)

	return wal.CurrentAddress(destinationAccount.Number)
}

func (dst *destination) destinationAccount() *dcrlibwallet.Account {
	if dst.sendToAddress {
		return nil
	}

	return dst.destinationAccountSelector.SelectedAccount()
}

func (dst *destination) validateDestinationAddress() (bool, string) {

	address := dst.destinationAddressEditor.Editor.Text()
	address = strings.TrimSpace(address)

	if len(address) == 0 {
		dst.destinationAddressEditor.SetError("")
		return false, address
	}

	if dst.WL.MultiWallet.IsAddressValid(address) {
		dst.destinationAddressEditor.SetError("")
		return true, address
	}

	dst.destinationAddressEditor.SetError("Invalid address")
	return false, address
}

func (dst *destination) validate() bool {
	if dst.sendToAddress {
		validAddress, _ := dst.validateDestinationAddress()
		return validAddress
	}

	return true
}

func (dst *destination) clearAddressInput() {
	dst.destinationAddressEditor.SetError("")
	dst.destinationAddressEditor.Editor.SetText("")
}

func (dst *destination) handle() {
	sendToAddress := dst.accountSwitch.SelectedIndex() == 1
	if sendToAddress != dst.sendToAddress { // switch changed
		dst.sendToAddress = sendToAddress
		dst.addressChanged()
	}

	for _, evt := range dst.destinationAddressEditor.Editor.Events() {
		if dst.destinationAddressEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				dst.addressChanged()
			}
		}
	}
}
