package send

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/io/key"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	PageID = "Send"
)

type moreItem struct {
	text     string
	id       string
	button   *decredmaterial.Clickable
	action   func()
	separate bool
}

type Page struct {
	*load.Load
	pageContainer *widget.List

	sourceAccountSelector *components.AccountSelector
	sendDestination       *destination
	amount                *sendAmount
	keyEvent              chan *key.Event

	backButton    decredmaterial.IconButton
	infoButton    decredmaterial.IconButton
	moreOption    decredmaterial.IconButton
	retryExchange decredmaterial.Button
	nextButton    decredmaterial.Button

	txFeeCollapsible *decredmaterial.Collapsible
	shadowBox        *decredmaterial.Shadow
	optionsMenuCard  decredmaterial.Card
	moreItems        []moreItem
	backdrop         *widget.Clickable

	moreOptionIsOpen bool

	exchangeRate   float64
	usdExchangeSet bool
	exchangeError  string

	*authoredTxData
}

type authoredTxData struct {
	txAuthor            *dcrlibwallet.TxAuthor
	destinationAddress  string
	destinationAccount  *dcrlibwallet.Account
	sourceAccount       *dcrlibwallet.Account
	txFee               string
	txFeeUSD            string
	estSignedSize       string
	totalCost           string
	totalCostUSD        string
	balanceAfterSend    string
	balanceAfterSendUSD string
	sendAmount          string
	sendAmountUSD       string
}

func NewSendPage(l *load.Load) *Page {
	pg := &Page{
		Load:            l,
		sendDestination: newSendDestination(l),
		amount:          newSendAmount(l),

		exchangeRate: -1,

		authoredTxData: &authoredTxData{},
		shadowBox:      l.Theme.Shadow(),
		backdrop:       new(widget.Clickable),
		keyEvent:       l.Receiver.KeyEvents,
	}

	// Source account picker
	pg.sourceAccountSelector = components.NewAccountSelector(l).
		Title("Sending account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {
			pg.validateAndConstructTx()
		}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != load.MaxInt32 && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				if pg.sendDestination.sendToAddress {
					// only mixed can send to address
					accountIsValid = account.Number == wal.MixedAccountNumber()
				} else {
					// send to account, check if selected destination account belongs to wallet
					destinationAccount := pg.sendDestination.destinationAccountSelector.SelectedAccount()
					if destinationAccount.WalletID != account.WalletID {
						accountIsValid = account.Number == wal.MixedAccountNumber()
					}
				}
			}
			return accountIsValid
		})

	pg.sendDestination.destinationAccountSelector.AccountSelected(func(selectedAccount *dcrlibwallet.Account) {
		pg.validateAndConstructTx()
		pg.sourceAccountSelector.SelectFirstWalletValidAccount() // refresh source account
	})

	pg.sendDestination.addressChanged = func() {
		pg.validateAndConstructTx()
	}

	pg.amount.amountChanged = func() {
		pg.validateAndConstructTx()
	}

	pg.initLayoutWidgets()

	return pg
}

func (pg *Page) ID() string {
	return PageID
}

func (pg *Page) OnResume() {
	pg.sendDestination.destinationAccountSelector.SelectFirstWalletValidAccount()
	pg.sourceAccountSelector.SelectFirstWalletValidAccount()
	pg.sendDestination.destinationAddressEditor.Editor.Focus()

	currencyExchangeValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyExchangeValue == values.USDExchangeValue {
		pg.usdExchangeSet = true
		pg.fetchExchangeValue()
	} else {
		pg.usdExchangeSet = false
	}
	pg.Load.EnableKeyEvent = true
}

func (pg *Page) fetchExchangeValue() {
	pg.exchangeError = ""
	go func() {
		var dcrUsdtBittrex load.DCRUSDTBittrex
		err := load.GetUSDExchangeValue(&dcrUsdtBittrex)
		if err != nil {
			pg.exchangeError = err.Error()
			return
		}

		exchangeRate, err := strconv.ParseFloat(dcrUsdtBittrex.LastTradeRate, 64)
		if err != nil {
			pg.exchangeError = err.Error()
			return
		}

		pg.exchangeError = ""
		pg.exchangeRate = exchangeRate
		pg.amount.setExchangeRate(exchangeRate)
		pg.validateAndConstructTx() // convert estimates to usd
	}()
}

func (pg *Page) validateAndConstructTx() {
	if pg.validate() {
		pg.constructTx()
	} else {
		pg.clearEstimates()
	}
}

func (pg *Page) validate() bool {

	amountIsValid := pg.amount.amountIsValid()
	addressIsValid := pg.sendDestination.validate()

	validForSending := amountIsValid && addressIsValid
	pg.nextButton.SetEnabled(validForSending)

	return validForSending
}

func (pg *Page) constructTx() {
	destinationAddress, err := pg.sendDestination.destinationAddress()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}
	destinationAccount := pg.sendDestination.destinationAccount()

	amountAtom, sendMax, err := pg.amount.validAmount()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	sourceAccount := pg.sourceAccountSelector.SelectedAccount()
	unsignedTx, err := pg.WL.MultiWallet.NewUnsignedTx(sourceAccount.WalletID, sourceAccount.Number)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	err = unsignedTx.AddSendDestination(destinationAddress, amountAtom, sendMax)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	feeAndSize, err := unsignedTx.EstimateFeeAndSize()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	feeAtom := feeAndSize.Fee.AtomValue
	if sendMax {
		amountAtom = sourceAccount.Balance.Spendable - feeAtom
	}

	totalSendingAmount := dcrutil.Amount(amountAtom + feeAtom)
	balanceAfterSend := dcrutil.Amount(sourceAccount.Balance.Spendable - int64(totalSendingAmount))

	// populate display data
	pg.txFee = dcrutil.Amount(feeAtom).String()
	pg.estSignedSize = fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize)
	pg.totalCost = totalSendingAmount.String()
	pg.balanceAfterSend = balanceAfterSend.String()
	pg.sendAmount = dcrutil.Amount(amountAtom).String()
	pg.destinationAddress = destinationAddress
	pg.destinationAccount = destinationAccount
	pg.sourceAccount = sourceAccount

	if sendMax {
		// TODO: this workaround ignores the change events from the
		// amount input to avoid construct tx cycle.
		pg.amount.setAmount(amountAtom)
	}

	if pg.exchangeRate != -1 && pg.usdExchangeSet {
		pg.txFeeUSD = fmt.Sprintf("$%.4f", load.DCRToUSD(pg.exchangeRate, feeAndSize.Fee.DcrValue))
		pg.totalCostUSD = load.FormatUSDBalance(pg.Printer, load.DCRToUSD(pg.exchangeRate, totalSendingAmount.ToCoin()))
		pg.balanceAfterSendUSD = load.FormatUSDBalance(pg.Printer, load.DCRToUSD(pg.exchangeRate, balanceAfterSend.ToCoin()))

		usdAmount := load.DCRToUSD(pg.exchangeRate, dcrutil.Amount(amountAtom).ToCoin())
		pg.sendAmountUSD = load.FormatUSDBalance(pg.Printer, usdAmount)
	}

	pg.txAuthor = unsignedTx
}

func (pg *Page) feeEstimationError(err string) {
	if err == dcrlibwallet.ErrInsufficientBalance {
		pg.amount.setError("Not enough funds")
	} else if strings.Contains(err, invalidAmountErr) {
		pg.amount.setError(invalidAmountErr)
	} else {
		pg.amount.setError(err)
		pg.Toast.NotifyError("Error estimating transaction: " + err)
	}

	pg.clearEstimates()
}

func (pg *Page) clearEstimates() {
	pg.txAuthor = nil
	pg.txFee = " - "
	pg.txFeeUSD = " - "
	pg.estSignedSize = " - "
	pg.totalCost = " - "
	pg.totalCostUSD = " - "
	pg.balanceAfterSend = " - "
	pg.balanceAfterSendUSD = " - "
	pg.sendAmount = " - "
	pg.sendAmountUSD = " - "
}

func (pg *Page) resetFields() {
	pg.sendDestination.clearAddressInput()

	pg.amount.resetFields()
}

func (pg *Page) Handle() {
	pg.sendDestination.handle()
	pg.amount.handle()

	if pg.backButton.Button.Clicked() {
		pg.PopFragment()
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title("Send DCR").
			Body("Input or scan the destination wallet address and input the amount to send funds.").
			PositiveButton("Got it", func() {})
		pg.ShowModal(info)
	}

	for pg.moreOption.Button.Clicked() {
		pg.moreOptionIsOpen = !pg.moreOptionIsOpen
	}

	for pg.retryExchange.Clicked() {
		pg.fetchExchangeValue()
	}

	for pg.nextButton.Clicked() {
		if pg.txAuthor != nil {
			confirmTxModal := newSendConfirmModal(pg.Load, pg.authoredTxData)
			confirmTxModal.exchangeRateSet = pg.exchangeRate != -1 && pg.usdExchangeSet

			confirmTxModal.txSent = func() {
				pg.resetFields()
				pg.clearEstimates()
			}

			confirmTxModal.Show()
		}
	}

	for pg.backdrop.Clicked() {
		pg.moreOptionIsOpen = false
	}

	for _, menu := range pg.moreItems {
		if menu.button.Clicked() {
			menu.action()
		}
	}

	currencyValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyValue != values.USDExchangeValue {
		switch {
		case !pg.sendDestination.sendToAddress:
			if !pg.amount.dcrAmountEditor.Editor.Focused() {
				pg.amount.dcrAmountEditor.Editor.Focus()
			}
		default:
			if pg.sendDestination.accountSwitch.Changed() {
				if !pg.sendDestination.validate() {
					pg.sendDestination.destinationAddressEditor.Editor.Focus()
				} else {
					pg.amount.dcrAmountEditor.Editor.Focus()
				}

			}

			decredmaterial.SwitchEditors(pg.keyEvent, pg.sendDestination.destinationAddressEditor.Editor, pg.amount.dcrAmountEditor.Editor)
		}
	} else {
		switch {
		case !pg.sendDestination.sendToAddress && !(pg.amount.dcrAmountEditor.Editor.Focused() || pg.amount.usdAmountEditor.Editor.Focused()):
			pg.amount.dcrAmountEditor.Editor.Focus()
		case !pg.sendDestination.sendToAddress && (pg.amount.dcrAmountEditor.Editor.Focused() || pg.amount.usdAmountEditor.Editor.Focused()):
			decredmaterial.SwitchEditors(pg.keyEvent, pg.amount.usdAmountEditor.Editor, pg.amount.dcrAmountEditor.Editor)
		default:
			if pg.sendDestination.accountSwitch.Changed() {
				if !pg.sendDestination.validate() {
					pg.sendDestination.destinationAddressEditor.Editor.Focus()
				} else {
					pg.amount.dcrAmountEditor.Editor.Focus()
				}
			}
			decredmaterial.SwitchEditors(pg.keyEvent, pg.sendDestination.destinationAddressEditor.Editor, pg.amount.dcrAmountEditor.Editor, pg.amount.usdAmountEditor.Editor)
		}
	}

	if pg.sendDestination.sendToAddress {
		addEditor := pg.sendDestination.destinationAddressEditor
		if pg.amount.IsMaxClicked() {
			validAddress := pg.sendDestination.validate()
			if validAddress {
				pg.amount.setError("")
				pg.amount.Sendmax = true
				pg.amount.amountChanged()
			} else if len(addEditor.Editor.Text()) < 1 {
				pg.Toast.NotifyError("Set destination amount")
			}

		}
	} else {
		if pg.amount.IsMaxClicked() {
			pg.amount.setError("")
			pg.amount.Sendmax = true
			pg.amount.amountChanged()
		}
	}
}

func (pg *Page) OnClose() {
	pg.Load.EnableKeyEvent = false
}
