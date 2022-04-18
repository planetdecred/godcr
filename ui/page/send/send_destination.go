package send

import (
	"fmt"
	"image/color"
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

	dst.destinationAddressEditor = l.Theme.Editor(new(widget.Editor), "Destination Address")
	dst.destinationAddressEditor.Editor.SingleLine = true
	dst.destinationAddressEditor.Editor.SetText("")

	dst.accountSwitch = l.Theme.SwitchButtonText([]decredmaterial.SwitchItem{{Text: "Address"}, {Text: "My account"}})

	// Destination account picker
	dst.destinationAccountSelector = components.NewAccountSelector(dst.Load, nil).
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

func (dst *destination) destinationAddress(useDefaultParams bool) (string, error) {
	destinationAccount := dst.destinationAccountSelector.SelectedAccount()
	wal := dst.WL.MultiWallet.WalletWithID(destinationAccount.WalletID)

	if useDefaultParams {
		return wal.CurrentAddress(destinationAccount.Number)
	}

	if dst.sendToAddress {
		valid, address := dst.validateDestinationAddress()
		if valid {
			return address, nil
		}

		return "", fmt.Errorf("invalid address")
	}

	return wal.CurrentAddress(destinationAccount.Number)
}

func (dst *destination) destinationAccount(useDefaultParams bool) *dcrlibwallet.Account {
	if useDefaultParams {
		return dst.destinationAccountSelector.SelectedAccount()
	}

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

// styleWidgets sets the appropriate colors for the destination widgets.
func (dst *destination) styleWidgets() {
	dst.accountSwitch.Active, dst.accountSwitch.Inactive = dst.Theme.Color.Surface, color.NRGBA{}
	dst.accountSwitch.ActiveTextColor, dst.accountSwitch.InactiveTextColor = dst.Theme.Color.GrayText1, dst.Theme.Color.Text
	dst.destinationAddressEditor.EditorStyle.Color = dst.Theme.Color.Text
}
