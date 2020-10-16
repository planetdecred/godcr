package ui

import (
	"fmt"
	"strconv"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageTickets = "tickets"

type ticketPage struct {
	vspd                 *dcrlibwallet.VSPD
	ticketPageContainer  layout.List
	purchaseTicketButton decredmaterial.Button
	selectedWallet       wallet.InfoShort
	selectedAccount      wallet.Account
	shouldInitializeVSPD bool
	inputNumbTickets     decredmaterial.Editor
	submitButton         decredmaterial.Button
	passwordModal        *decredmaterial.Password
	password             []byte
	isPasswordModalOpen  bool
	tiketHash            string
	ticketPrice          string
}

func (win *Window) TicketPage(common pageCommon) layout.Widget {
	pg := &ticketPage{
		ticketPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		submitButton:         common.theme.Button(new(widget.Clickable), "Vote tiket"),
		inputNumbTickets:     common.theme.Editor(new(widget.Editor), "Enter tickets amount"),
		purchaseTicketButton: common.theme.Button(new(widget.Clickable), "Purchase"),
		passwordModal:        common.theme.Password(),
	}
	pg.submitButton.TextSize = values.TextSize12
	pg.purchaseTicketButton.TextSize = values.TextSize12

	pg.inputNumbTickets.IsRequired = true
	pg.inputNumbTickets.Editor.SingleLine = true

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *ticketPage) Handler(c pageCommon) {
	if pg.selectedAccount.CurrentAddress != c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount].CurrentAddress {
		pg.shouldInitializeVSPD = true
		pg.selectedAccount = c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount]
	}
	if pg.selectedWallet.ID != c.info.Wallets[*c.selectedWallet].ID {
		pg.shouldInitializeVSPD = true
		pg.selectedWallet = c.info.Wallets[*c.selectedWallet]
	}

	if pg.shouldInitializeVSPD {
		pg.shouldInitializeVSPD = false
		tkPrice := c.wallet.TicketPrice(pg.selectedWallet.ID)
		pg.ticketPrice = fmt.Sprintf("Ticket price: %s", tkPrice)
		pg.vspd = c.wallet.NewVSPD(pg.selectedWallet.ID, pg.selectedAccount.Number)
		_, err := pg.vspd.GetInfo()
		if err != nil {
			log.Error("[GetInfo] err:", err)
			return
		}
	}

	if pg.purchaseTicketButton.Button.Clicked() {
		pg.isPasswordModalOpen = true
	}

	if pg.submitButton.Button.Clicked() {
		log.Info(">>>>>>>>>>>>>> pg.password", string(pg.password), pg.tiketHash)
		resp, err := pg.vspd.GetVSPFeeAddress(pg.tiketHash, pg.password[:])
		if err != nil {
			log.Error("Get fee error", err)
			return
		}
		log.Info(">>>", resp.FeeAddress)
		log.Info(">>>", resp.FeeAmount)
		feeTxStr, err := pg.vspd.CreateTicketFeeTx(resp.FeeAmount, resp.FeeAddress, []byte("123"))
		if err != nil {
			log.Error("Create ticket eror", err)
			return
		}
		log.Info(">>>", feeTxStr)
		msg, err := pg.vspd.PayVSPFee(feeTxStr, pg.tiketHash, resp.FeeAddress, []byte("123"))
		if err != nil {
			log.Error("Pay fee error", err)
			return
		}
		log.Info("5: [Done]", msg)
	}
}

func (pg *ticketPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	margin := values.MarginPadding20

	widgets := []func(gtx C) D{
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return common.theme.H5(pg.ticketPrice).Layout(gtx)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return pg.inputNumbTickets.Layout(gtx)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return layout.Dimensions{}
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return layout.Dimensions{}
			})
		},
	}

	dims := common.LayoutWithAccounts(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, common.SelectedAccountLayout)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.ticketPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return layout.Inset{}.Layout(gtx, widgets[i])
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.purchaseTicketButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
						str := fmt.Sprintf("Vote tiket")
						if pg.tiketHash != "" {
							str = fmt.Sprintf("Vote ticket [%s]", pg.tiketHash)
							pg.submitButton.Color = common.theme.Color.Success
						}
						pg.submitButton.Text = str
						return pg.submitButton.Layout(gtx)
					})
				}),
			)
		})
	})

	if pg.isPasswordModalOpen {
		return common.Modal(gtx, dims, pg.drawPasswordModal(gtx, &common))
	}
	return dims
}

func (pg *ticketPage) drawPasswordModal(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return pg.passwordModal.Layout(gtx, func(password []byte) {
		numbTicketsStr := pg.inputNumbTickets.Editor.Text()
		numbTickets, err := strconv.Atoi(numbTicketsStr)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}

		h, err := c.wallet.PurchaseTicket(pg.selectedWallet.ID, pg.selectedAccount.Number, uint32(numbTickets), password)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}
		pg.password = password
		pg.tiketHash = h
		pg.isPasswordModalOpen = false
	}, func() {
		pg.password = nil
		pg.isPasswordModalOpen = false
	})
}
