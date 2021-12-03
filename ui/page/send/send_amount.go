package send

import (
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const invalidAmountErr = "Invalid amount" //TODO: use localized strings

type sendAmount struct {
	*load.Load

	dcrAmountEditor decredmaterial.Editor
	usdAmountEditor decredmaterial.Editor

	sendMax               bool
	dcrSendMaxChangeEvent bool
	usdSendMaxChangeEvent bool
	amountChanged         func()

	amountErrorText string

	exchangeRate float64
}

func newSendAmount(l *load.Load) *sendAmount {

	sa := &sendAmount{
		Load:         l,
		exchangeRate: -1,
	}

	sa.dcrAmountEditor = l.Theme.Editor(new(widget.Editor), "Amount (DCR)")
	sa.dcrAmountEditor.Editor.SetText("")
	sa.dcrAmountEditor.HasCustomButton = true
	sa.dcrAmountEditor.Editor.SingleLine = true
	sa.dcrAmountEditor.CustomButton.Background = l.Theme.Color.Gray1
	sa.dcrAmountEditor.CustomButton.Color = l.Theme.Color.Surface
	sa.dcrAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	sa.dcrAmountEditor.CustomButton.Text = "Max"
	sa.dcrAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	sa.usdAmountEditor = l.Theme.Editor(new(widget.Editor), "Amount (USD)")
	sa.usdAmountEditor.Editor.SetText("")
	sa.usdAmountEditor.HasCustomButton = true
	sa.usdAmountEditor.Editor.SingleLine = true
	sa.usdAmountEditor.CustomButton.Background = l.Theme.Color.Gray1
	sa.usdAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	sa.usdAmountEditor.CustomButton.Text = "Max"
	sa.usdAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	return sa
}

func (sa *sendAmount) setExchangeRate(exchangeRate float64) {
	sa.exchangeRate = exchangeRate
	sa.validateDCRAmount() // convert dcr input to usd
}

func (sa *sendAmount) setAmount(amount int64) {
	// TODO: this workaround ignores the change events from the
	// amount input to avoid construct tx cycle.
	sa.dcrSendMaxChangeEvent = sa.sendMax
	sa.dcrAmountEditor.Editor.SetText(fmt.Sprintf("%.8f", dcrutil.Amount(amount).ToCoin()))

	if sa.exchangeRate != -1 {
		usdAmount := load.DCRToUSD(sa.exchangeRate, dcrutil.Amount(amount).ToCoin())

		sa.usdSendMaxChangeEvent = true
		sa.usdAmountEditor.Editor.SetText(fmt.Sprintf("%.2f", usdAmount))

	}
}

func (sa *sendAmount) amountIsValid() bool {
	_, err := strconv.ParseFloat(sa.dcrAmountEditor.Editor.Text(), 64)
	return err == nil || sa.sendMax
}

func (sa *sendAmount) validAmount() (int64, bool, error) {
	if sa.sendMax {
		return 0, sa.sendMax, nil
	}

	amount, err := strconv.ParseFloat(sa.dcrAmountEditor.Editor.Text(), 64)
	if err != nil {
		return -1, sa.sendMax, err
	}

	return dcrlibwallet.AmountAtom(amount), sa.sendMax, nil
}

func (sa *sendAmount) validateDCRAmount() {
	sa.amountErrorText = ""
	if sa.inputsNotEmpty(sa.dcrAmountEditor.Editor) {
		dcrAmount, err := strconv.ParseFloat(sa.dcrAmountEditor.Editor.Text(), 64)
		if err != nil {
			// empty usd input
			sa.usdAmountEditor.Editor.SetText("")
			sa.amountErrorText = invalidAmountErr
			// todo: invalid decimal places error
			return
		}

		if sa.exchangeRate != -1 {
			usdAmount := load.DCRToUSD(sa.exchangeRate, dcrAmount)
			sa.usdAmountEditor.Editor.SetText(fmt.Sprintf("%.2f", usdAmount)) // 2 decimal places
		}

		return
	}

	// empty usd input since this is empty
	sa.usdAmountEditor.Editor.SetText("")
}

// validateUSDAmount is called when usd text changes
func (sa *sendAmount) validateUSDAmount() bool {
	sa.amountErrorText = ""
	if sa.inputsNotEmpty(sa.usdAmountEditor.Editor) {
		usdAmount, err := strconv.ParseFloat(sa.usdAmountEditor.Editor.Text(), 64)
		if err != nil {
			// empty dcr input
			sa.dcrAmountEditor.Editor.SetText("")
			sa.amountErrorText = invalidAmountErr
			return false
		}

		if sa.exchangeRate != -1 {
			dcrAmount := load.USDToDCR(sa.exchangeRate, usdAmount)
			sa.dcrAmountEditor.Editor.SetText(fmt.Sprintf("%.8f", dcrAmount)) // 8 decimal places
		}

		return true
	}

	// empty dcr input since this is empty
	sa.dcrAmountEditor.Editor.SetText("")
	return false
}

func (sa *sendAmount) inputsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}

func (sa *sendAmount) setError(err string) {
	sa.amountErrorText = err
}

func (sa *sendAmount) resetFields() {
	sa.sendMax = false

	sa.clearAmount()
}

func (sa *sendAmount) clearAmount() {
	sa.amountErrorText = ""
	sa.dcrAmountEditor.Editor.SetText("")
	sa.usdAmountEditor.Editor.SetText("")
}

func (sa *sendAmount) handle() {
	sa.dcrAmountEditor.SetError(sa.amountErrorText)

	if sa.amountErrorText != "" {
		sa.dcrAmountEditor.LineColor = sa.Theme.Color.Danger
		sa.usdAmountEditor.LineColor = sa.Theme.Color.Danger
	} else {
		sa.dcrAmountEditor.LineColor = sa.Theme.Color.Gray2
		sa.usdAmountEditor.LineColor = sa.Theme.Color.Gray2
	}

	if sa.sendMax {
		sa.dcrAmountEditor.CustomButton.Background = sa.Theme.Color.Primary
		sa.usdAmountEditor.CustomButton.Background = sa.Theme.Color.Primary
	} else {
		sa.dcrAmountEditor.CustomButton.Background = sa.Theme.Color.Gray1
		sa.usdAmountEditor.CustomButton.Background = sa.Theme.Color.Gray1
	}

	for _, evt := range sa.dcrAmountEditor.Editor.Events() {
		if sa.dcrAmountEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				if sa.dcrSendMaxChangeEvent {
					sa.dcrSendMaxChangeEvent = false
					continue
				}
				sa.sendMax = false
				sa.validateDCRAmount()
				sa.amountChanged()

			}
		}
	}

	for _, evt := range sa.usdAmountEditor.Editor.Events() {
		if sa.usdAmountEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				if sa.usdSendMaxChangeEvent {
					sa.usdSendMaxChangeEvent = false
					continue
				}
				sa.sendMax = false
				sa.validateUSDAmount()
				sa.amountChanged()
			}
		}
	}

	for sa.dcrAmountEditor.CustomButton.Clicked() ||
		sa.usdAmountEditor.CustomButton.Clicked() {
		sa.setError("")
		sa.sendMax = true
		sa.amountChanged()
	}
}
