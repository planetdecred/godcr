package send

import (
	"fmt"
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const ModalSendConfirm = "send_confirm_modal"

type sendConfirmModal struct {
	*load.Load
	modal *decredmaterial.Modal

	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	passwordEditor               decredmaterial.Editor

	txSent    func()
	isSending bool

	*authoredTxData
	exchangeRateSet    bool
	destinationAddress string
	destinationAccount *dcrlibwallet.Account
	sourceAccount      *dcrlibwallet.Account
}

func newSendConfirmModal(l *load.Load, data *authoredTxData) *sendConfirmModal {
	scm := &sendConfirmModal{
		Load:  l,
		modal: l.Theme.ModalFloatTitle(),

		authoredTxData: data,
	}

	scm.closeConfirmationModalButton = l.Theme.Button(new(widget.Clickable), "Cancel")
	scm.closeConfirmationModalButton.Background = color.NRGBA{}
	scm.closeConfirmationModalButton.Color = l.Theme.Color.Primary

	scm.confirmButton = l.Theme.Button(new(widget.Clickable), "")
	scm.confirmButton.Background = scm.Theme.Color.InactiveGray

	scm.passwordEditor = l.Theme.EditorPassword(new(widget.Editor), "Spending password")
	scm.passwordEditor.Editor.SetText("")
	scm.passwordEditor.Editor.SingleLine = true
	scm.passwordEditor.Editor.Submit = true

	return scm
}

func (scm *sendConfirmModal) ModalID() string {
	return ModalSendConfirm
}

func (scm *sendConfirmModal) Show() {
	scm.ShowModal(scm)
}

func (scm *sendConfirmModal) Dismiss() {
	scm.DismissModal(scm)
}

func (scm *sendConfirmModal) OnResume() {
	scm.passwordEditor.Editor.Focus()
}

func (scm *sendConfirmModal) OnDismiss() {

}

func (scm *sendConfirmModal) broadcastTransaction() {
	password := scm.passwordEditor.Editor.Text()
	if password == "" || scm.isSending {
		return
	}

	scm.isSending = true
	go func() {
		_, err := scm.authoredTxData.txAuthor.Broadcast([]byte(password))
		scm.isSending = false
		if err != nil {
			scm.CreateToast(err.Error(), false)
			return
		}
		scm.CreateToast("Transaction sent!", true)

		scm.txSent()
		scm.Dismiss()
	}()
}

func (scm *sendConfirmModal) Handle() {
	for _, evt := range scm.passwordEditor.Editor.Events() {
		if scm.passwordEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				if scm.passwordEditor.Editor.Text() == "" {
					scm.confirmButton.Background = scm.Theme.Color.InactiveGray
				} else {
					scm.confirmButton.Background = scm.Theme.Color.Primary
				}
			case widget.SubmitEvent:
				scm.broadcastTransaction()
			}
		}
	}

	for scm.confirmButton.Button.Clicked() {
		scm.broadcastTransaction()
	}

	for scm.closeConfirmationModalButton.Button.Clicked() {
		if !scm.isSending {
			scm.Dismiss()
		}
	}
}

func (scm *sendConfirmModal) Layout(gtx layout.Context) D {

	w := []layout.Widget{
		func(gtx C) D {
			return scm.Theme.H6("Confim to send").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := scm.Icons.SendIcon
					icon.Scale = 0.7
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding2, Right: values.MarginPadding16}.Layout(gtx, icon.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return components.LayoutBalance(gtx, scm.Load, scm.sendAmount)
								}),
								layout.Flexed(1, func(gtx C) D {
									if scm.exchangeRateSet {
										return layout.E.Layout(gtx, func(gtx C) D {
											txt := scm.Theme.Body1(scm.sendAmountUSD)
											txt.Color = scm.Theme.Color.Gray
											return txt.Layout(gtx)
										})
									}
									return layout.Dimensions{}
								}),
							)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							icon := scm.Icons.NavigationArrowForward
							icon.Color = scm.Theme.Color.Gray3
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								return icon.Layout(gtx, values.MarginPadding15)
							})
						}),
						layout.Rigid(func(gtx C) D {
							if scm.destinationAccount != nil {
								return layout.E.Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return scm.Theme.Body2(scm.destinationAccount.Name).Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											card := scm.Theme.Card()
											card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
											card.Color = scm.Theme.Color.LightGray
											inset := layout.Inset{
												Left: values.MarginPadding5,
											}
											return inset.Layout(gtx, func(gtx C) D {
												return card.Layout(gtx, func(gtx C) D {
													return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
														destinationWallet := scm.WL.MultiWallet.WalletWithID(scm.destinationAccount.WalletID)
														txt := scm.Theme.Caption(destinationWallet.Name)
														txt.Color = scm.Theme.Color.Gray
														return txt.Layout(gtx)
													})
												})
											})
										}),
									)
								})
							}
							return scm.Theme.Body2(scm.destinationAddress).Layout(gtx)
						}),
					)
				}),
			)
		},
		func(gtx C) D {
			return scm.Theme.Separator().Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					sendWallet := scm.WL.MultiWallet.WalletWithID(scm.sourceAccount.WalletID)
					return scm.contentRow(gtx, "Sending from", scm.sourceAccount.Name, sendWallet.Name)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						txFeeText := scm.txFee
						if scm.exchangeRateSet {
							txFeeText = fmt.Sprintf("%s (%s)", scm.txFee, scm.txFeeUSD)
						}

						return scm.contentRow(gtx, "Fee", txFeeText, "")
					})
				}),
				layout.Rigid(func(gtx C) D {
					totalCostText := scm.totalCost
					if scm.exchangeRateSet {
						totalCostText = fmt.Sprintf("%s (%s)", scm.totalCost, scm.totalCostUSD)
					}

					return scm.contentRow(gtx, "Total cost", totalCostText, "")
				}),
			)
		},
		func(gtx C) D {
			return scm.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := scm.Icons.ActionInfo
					icon.Color = scm.Theme.Color.Gray
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return icon.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := scm.Theme.Body2("Your DCR will be sent after this step.")
					txt.Color = scm.Theme.Color.Gray3
					return txt.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding8,
						}.Layout(gtx, func(gtx C) D {
							return scm.closeConfirmationModalButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if scm.isSending {
							th := material.NewTheme(gofont.Collection())
							return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx C) D {
								return material.Loader(th).Layout(gtx)
							})
						}
						scm.confirmButton.Text = fmt.Sprintf("Send %s", scm.totalCost)
						return scm.confirmButton.Layout(gtx)
					}),
				)
			})
		},
	}

	return scm.modal.Layout(gtx, w, 900)
}

func (scm *sendConfirmModal) contentRow(gtx layout.Context, leftValue, rightValue, walletName string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := scm.Theme.Body2(leftValue)
			txt.Color = scm.Theme.Color.Gray
			return txt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(scm.Theme.Body1(rightValue).Layout),
					layout.Rigid(func(gtx C) D {
						if walletName != "" {
							card := scm.Theme.Card()
							card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
							card.Color = scm.Theme.Color.LightGray
							inset := layout.Inset{
								Left: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := scm.Theme.Caption(walletName)
										txt.Color = scm.Theme.Color.Gray
										return txt.Layout(gtx)
									})
								})
							})
						}
						return layout.Dimensions{}
					}),
				)
			})
		}),
	)
}
