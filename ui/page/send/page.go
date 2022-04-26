package send

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gioui.org/io/key"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
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
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

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

	moreOptionIsOpen       bool
	isFetchingExchangeRate bool

	exchangeRate        float64
	usdExchangeSet      bool
	exchangeRateMessage string
	confirmTxModal      *sendConfirmModal

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
		keyEvent:       make(chan *key.Event),
	}

	// Source account picker
	pg.sourceAccountSelector = components.NewAccountSelector(l, nil).
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
		pg.sourceAccountSelector.SelectFirstWalletValidAccount(nil, []int32{}) // refresh source account
	})

	pg.sendDestination.addressChanged = func() {
		pg.validateAndConstructTx()
	}

	pg.amount.amountChanged = func() {
		pg.validateAndConstructTxAmountOnly()
	}

	pg.initLayoutWidgets()

	return pg
}

// RestyleWidgets restyles select widgets to match the current theme. This is
// especially necessary when the dark mode setting is changed.
func (pg *Page) RestyleWidgets() {
	pg.amount.styleWidgets()
	pg.sendDestination.styleWidgets()
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *Page) ID() string {
	return PageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedTo() {
	pg.RestyleWidgets()

	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.sourceAccountSelector.ListenForTxNotifications(pg.ctx)
	pg.sendDestination.destinationAccountSelector.SelectFirstWalletValidAccount(nil, []int32{})
	pg.sourceAccountSelector.SelectFirstWalletValidAccount(nil, []int32{})
	pg.sendDestination.destinationAddressEditor.Editor.Focus()

	currencyExchangeValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyExchangeValue == values.USDExchangeValue {
		pg.usdExchangeSet = true
		go pg.fetchExchangeRate()
	} else {
		pg.usdExchangeSet = false
	}
	pg.Load.SubscribeKeyEvent(pg.keyEvent, pg.ID())
}

// OnDarkModeChanged is triggered whenever the dark mode setting is changed
// to enable restyling UI elements where necessary.
// Satisfies the load.DarkModeChangeHandler interface.
func (pg *Page) OnDarkModeChanged(isDarkModeOn bool) {
	pg.amount.styleWidgets()
}

func (pg *Page) fetchExchangeRate() {
	if pg.isFetchingExchangeRate {
		return
	}
	maxAttempts := 5
	delayBtwAttempts := 2 * time.Second
	pg.isFetchingExchangeRate = true
	desc := "for getting dcrUsdtBittrex exchange rate value"
	pg.exchangeRateMessage = "fetching exchange rate..."

	var dcrUsdtBittrex load.DCRUSDTBittrex
	attempts, err := components.RetryFunc(maxAttempts, delayBtwAttempts, desc, func() error {
		return load.GetUSDExchangeValue(&dcrUsdtBittrex)
	})
	if err != nil {
		pg.exchangeRateMessage = "Exchange rate not fetched. Kindly check internet connection."
		log.Printf("error fetching usd exchange rate value after %d attempts: %v", attempts, err)
	} else if dcrUsdtBittrex.LastTradeRate == "" {
		log.Printf("no error while fetching usd exchange rate in %d tries, but no rate was fetched", attempts)
		pg.exchangeRateMessage = "Exchange rate not fetched."
	} else {
		log.Printf("exchange rate value fetched: %s", dcrUsdtBittrex.LastTradeRate)
		pg.exchangeRateMessage = ""
		exchangeRate, err := strconv.ParseFloat(dcrUsdtBittrex.LastTradeRate, 64)
		if err != nil {
			pg.exchangeRateMessage = err.Error()
		} else {
			pg.exchangeRate = exchangeRate
			pg.amount.setExchangeRate(exchangeRate)
			pg.validateAndConstructTx() // convert estimates to usd
		}
	}
	pg.isFetchingExchangeRate = false
	pg.RefreshWindow()
}

func (pg *Page) validateAndConstructTx() {
	if pg.validate() {
		pg.constructTx(false)
	} else {
		pg.clearEstimates()
	}
}

func (pg *Page) validateAndConstructTxAmountOnly() {
	if !pg.sendDestination.validate() && pg.amount.amountIsValid() {
		pg.constructTx(true)
	} else {
		pg.validateAndConstructTx()
	}
}

func (pg *Page) validate() bool {
	amountIsValid := pg.amount.amountIsValid()
	addressIsValid := pg.sendDestination.validate()

	validForSending := amountIsValid && addressIsValid

	return validForSending
}

func (pg *Page) constructTx(useDefaultParams bool) {
	destinationAddress, err := pg.sendDestination.destinationAddress(useDefaultParams)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}
	destinationAccount := pg.sendDestination.destinationAccount(useDefaultParams)

	amountAtom, SendMax, err := pg.amount.validAmount()
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

	err = unsignedTx.AddSendDestination(destinationAddress, amountAtom, SendMax)
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
	if SendMax {
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

	if SendMax {
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

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *Page) HandleUserInteractions() {
	pg.nextButton.SetEnabled(pg.validate())
	pg.sendDestination.handle()
	pg.amount.handle()

	if pg.backButton.Button.Clicked() {
		pg.PopFragment()
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title("Send DCR").
			Body("Input or scan the destination wallet address and input the amount to send funds.").
			PositiveButton("Got it", func(isChecked bool) {})
		pg.ShowModal(info)
	}

	for pg.moreOption.Button.Clicked() {
		pg.moreOptionIsOpen = !pg.moreOptionIsOpen
	}

	for pg.retryExchange.Clicked() {
		go pg.fetchExchangeRate()
	}

	for pg.nextButton.Clicked() {
		if pg.txAuthor != nil {
			pg.confirmTxModal = newSendConfirmModal(pg.Load, pg.authoredTxData).SetParent(pg)
			pg.confirmTxModal.exchangeRateSet = pg.exchangeRate != -1 && pg.usdExchangeSet

			pg.confirmTxModal.txSent = func() {
				pg.resetFields()
				pg.clearEstimates()
			}

			pg.Load.UnsubscribeKeyEvent(pg.ID())
			pg.confirmTxModal.Show()
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

	modalShown := pg.confirmTxModal != nil && pg.confirmTxModal.IsShown()

	currencyValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyValue != values.USDExchangeValue {
		switch {
		case !pg.sendDestination.sendToAddress:
			if !pg.amount.dcrAmountEditor.Editor.Focused() && !modalShown {
				pg.amount.dcrAmountEditor.Editor.Focus()
			}
			decredmaterial.SwitchEditors(pg.keyEvent, pg.amount.dcrAmountEditor.Editor)
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
			if !modalShown {
				pg.amount.dcrAmountEditor.Editor.Focus()
			}
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

	// if destination switch is equal to Address
	if pg.sendDestination.sendToAddress {
		if pg.sendDestination.validate() {
			if currencyValue != values.USDExchangeValue {
				if len(pg.amount.dcrAmountEditor.Editor.Text()) == 0 {
					pg.amount.SendMax = false
				}
			} else {
				if len(pg.amount.dcrAmountEditor.Editor.Text()) == 0 {
					pg.amount.usdAmountEditor.Editor.SetText("")
					pg.amount.SendMax = false
				}
			}
		}
	} else {
		if currencyValue != values.USDExchangeValue {
			if len(pg.amount.dcrAmountEditor.Editor.Text()) == 0 {
				pg.amount.SendMax = false
			}
		} else {
			if len(pg.amount.dcrAmountEditor.Editor.Text()) == 0 {
				pg.amount.usdAmountEditor.Editor.SetText("")
				pg.amount.SendMax = false
			}
		}
	}

	if len(pg.amount.dcrAmountEditor.Editor.Text()) > 0 && pg.sourceAccountSelector.Changed() {
		pg.amount.validateDCRAmount()
		pg.validateAndConstructTxAmountOnly()
	}

	if pg.amount.IsMaxClicked() {
		pg.amount.setError("")
		pg.amount.SendMax = true
		pg.amount.amountChanged()
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedFrom() {
	pg.Load.UnsubscribeKeyEvent(pg.ID())
	pg.ctxCancel()
}
