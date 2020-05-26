package ui

import (
	"fmt"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type SendPage struct {
	pageContainer   layout.List
	theme           *decredmaterial.Theme
	wallet          *wallet.Wallet
	wallets         []wallet.InfoShort
	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	selectedWallet  wallet.InfoShort
	selectedAccount wallet.Account

	destinationAddressEditor           *widget.Editor
	sendAmountEditor                   *widget.Editor
	nextButtonWidget                   *widget.Button
	closeConfirmationModalButtonWidget *widget.Button
	confirmButtonWidget                *widget.Button
	copyIconWidget                     *widget.Button

	transactionFeeValueLabel   decredmaterial.Label
	totalCostValueLabel        decredmaterial.Label
	balanceAfterSendValueLabel decredmaterial.Label

	destinationAddressEditorMaterial     decredmaterial.Editor
	sendAmountEditorMaterial             decredmaterial.Editor
	nextButtonMaterial                   decredmaterial.Button
	closeConfirmationModalButtonMaterial decredmaterial.Button
	confirmButtonMaterial                decredmaterial.Button
	accountsTab                          *decredmaterial.Tabs
	walletsTab                           *decredmaterial.Tabs

	copyIconMaterial decredmaterial.IconButton

	remainingBalance   int64
	sendErrorText      string
	txHashText         string
	txHash             string
	calculateErrorText string

	passwordModal *decredmaterial.Password

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool
	hasInitializedTxAuthor    bool
	hasCopiedTxHash           bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error
}

const (
	PageSend = "send"
)

func (win *Window) SendPage(common pageCommon) layout.Widget {
	page := &SendPage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},

		theme:           common.theme,
		wallet:          common.wallet,
		wallets:         common.info.Wallets,
		txAuthor:        &win.txAuthor,
		broadcastResult: &win.broadcastResult,

		destinationAddressEditor:           new(widget.Editor),
		sendAmountEditor:                   new(widget.Editor),
		nextButtonWidget:                   new(widget.Button),
		closeConfirmationModalButtonWidget: new(widget.Button),
		confirmButtonWidget:                new(widget.Button),
		copyIconWidget:                     new(widget.Button),

		sendErrorText: "",
		txHashText:    "",

		closeConfirmationModalButtonMaterial: common.theme.Button("Close"),
		nextButtonMaterial:                   common.theme.Button("Next"),
		confirmButtonMaterial:                common.theme.Button("Confirm"),
		transactionFeeValueLabel:             common.theme.Body2("0 DCR"),
		totalCostValueLabel:                  common.theme.Body2("0 DCR"),
		balanceAfterSendValueLabel:           common.theme.Body2("0 DCR"),
		copyIconMaterial:                     common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy))),

		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		hasInitializedTxAuthor:    false,
		hasCopiedTxHash:           false,
		isBroadcastingTransaction: false,

		passwordModal: common.theme.Password(),
		accountsTab:   decredmaterial.NewTabs(),
		walletsTab:    decredmaterial.NewTabs(),

		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	page.walletsTab.Position = decredmaterial.Top
	page.accountsTab.Position = decredmaterial.Top

	page.destinationAddressEditorMaterial = common.theme.Editor("Destination Address")
	page.destinationAddressEditorMaterial.SetRequiredErrorText("")
	page.destinationAddressEditorMaterial.IsRequired = true
	page.destinationAddressEditorMaterial.IsVisible = true

	page.sendAmountEditorMaterial = common.theme.Editor("Amount to be sent")
	page.sendAmountEditorMaterial.SetRequiredErrorText("")
	page.sendAmountEditorMaterial.IsRequired = true

	page.closeConfirmationModalButtonMaterial.Background = common.theme.Color.Gray
	page.destinationAddressEditor.SetText("")

	page.copyIconMaterial.Background = common.theme.Color.Background
	page.copyIconMaterial.Color = common.theme.Color.Text
	page.copyIconMaterial.Size = unit.Dp(35)
	page.copyIconMaterial.Padding = unit.Dp(5)

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (pg *SendPage) Handle(common pageCommon) {
	pg.validate(true)
	pg.watchForBroadcastResult()

	if pg.walletsTab.Changed() {
		pg.selectedWallet = pg.wallets[pg.walletsTab.Selected]
		pg.selectedAccount = pg.selectedWallet.Accounts[0]
		pg.accountsTab.Selected = 0

		pg.setAccountTabs()
		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
	}

	if pg.accountsTab.Changed() {
		pg.selectedAccount = pg.selectedWallet.Accounts[pg.accountsTab.Selected]
		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
	}

	if pg.hasCopiedTxHash {
		time.AfterFunc(3*time.Second, func() {
			pg.hasCopiedTxHash = false
		})
	}

	if pg.isBroadcastingTransaction {
		col := pg.theme.Color.Gray
		col.A = 150
		pg.confirmButtonMaterial.Text = "Sending..."
		pg.confirmButtonMaterial.Background = col
	} else {
		pg.confirmButtonMaterial.Text = "Send"
		pg.confirmButtonMaterial.Background = pg.theme.Color.Primary
	}

	for pg.nextButtonWidget.Clicked(common.gtx) {
		if pg.validate(false) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButtonWidget.Clicked(common.gtx) {
		pg.sendErrorText = ""
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButtonWidget.Clicked(common.gtx) {
		pg.sendErrorText = ""
		pg.isConfirmationModalOpen = false
	}

	for range pg.destinationAddressEditor.Events(common.gtx) {
		go pg.calculateValues()
	}

	for range pg.sendAmountEditor.Events(common.gtx) {
		go pg.calculateValues()
	}

	for pg.copyIconWidget.Clicked(common.gtx) {
		clipboard.WriteAll(pg.txHash)
		pg.hasCopiedTxHash = true
	}

	select {
	case err := <-pg.txAuthorErrChan:
		pg.calculateErrorText = err.Error()
	case err := <-pg.broadcastErrChan:
		pg.sendErrorText = err.Error()
		pg.isBroadcastingTransaction = false
	default:
	}
}

func (pg *SendPage) drawWalletsTab(common pageCommon, body func()) {
	wallets := make([]decredmaterial.TabItem, len(pg.wallets))
	for i := range pg.wallets {
		wallets[i] = decredmaterial.TabItem{
			Label: pg.theme.Body1(pg.wallets[i].Name),
		}
	}
	pg.walletsTab.SetTabs(wallets)

	pg.setAccountTabs()
	pg.walletsTab.Layout(common.gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(common.gtx,
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(10), Right: unit.Dp(10)}.Layout(common.gtx, func() {
					pg.theme.H6("Accounts: ").Layout(common.gtx)
				})
			}),
			layout.Rigid(func() {
				pg.accountsTab.Layout(common.gtx, body)
			}),
		)
	})
}

func (pg *SendPage) setAccountTabs() {
	accounts := make([]decredmaterial.TabItem, len(pg.selectedWallet.Accounts))
	for i := range pg.selectedWallet.Accounts {
		if pg.selectedWallet.Accounts[i].Name == "imported" {
			continue
		}
		accounts[i] = decredmaterial.TabItem{
			Label: pg.theme.Body1(pg.selectedWallet.Accounts[i].Name),
		}
	}
	pg.accountsTab.SetTabs(accounts)
}

func (pg *SendPage) Layout(common pageCommon) {
	if len(common.info.Wallets) == 0 {
		// show no wallets message
		return
	}

	pg.wallets = common.info.Wallets
	if !pg.hasInitializedTxAuthor {
		pg.selectedWallet = pg.wallets[*common.selectedWallet]
		pg.selectedAccount = pg.selectedWallet.Accounts[0]

		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
		pg.hasInitializedTxAuthor = true
	}

	common.Layout(common.gtx, func() {
		pg.drawPageContents(common)
	})

	if pg.isConfirmationModalOpen {
		pg.drawConfirmationModal(common.gtx)

		if pg.isPasswordModalOpen {
			pg.drawPasswordModal(common.gtx)
		}
	}
}

func (pg *SendPage) drawPageContents(common pageCommon) {
	pageContent := []func(){
		func() {
			pg.drawSuccessSection(common.gtx)
		},
		func() {
			pg.drawCopiedLabelSection(common.gtx)
		},
		func() {
			pg.drawSelectedAccountSection(common.gtx)
		},
		func() {
			pg.destinationAddressEditorMaterial.Layout(common.gtx, pg.destinationAddressEditor)
		},
		func() {
			pg.sendAmountEditorMaterial.Layout(common.gtx, pg.sendAmountEditor)
		},
		func() {
			pg.drawTransactionDetailWidgets(common.gtx)
		},
		func() {
			if pg.calculateErrorText != "" {
				common.gtx.Constraints.Width.Min = common.gtx.Constraints.Width.Max
				pg.theme.ErrorAlert(common.gtx, pg.calculateErrorText)
			}
		},
		func() {
			common.gtx.Constraints.Width.Min = common.gtx.Constraints.Width.Max
			pg.nextButtonMaterial.Layout(common.gtx, pg.nextButtonWidget)
		},
	}

	w := func() {
		inset := layout.Inset{
			Left: unit.Dp(-110),
		}
		inset.Layout(common.gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(7)).Layout(common.gtx, func() {
						pg.pageContainer.Layout(common.gtx, len(pageContent), func(i int) {
							layout.Inset{Top: unit.Dp(5)}.Layout(common.gtx, pageContent[i])
						})
					})
				}),
			)
		})
	}
}
