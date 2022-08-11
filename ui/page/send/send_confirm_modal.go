package send

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type sendConfirmModal struct {
	*load.Load
	*decredmaterial.Modal

	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	passwordEditor               decredmaterial.Editor

	txSent    func()
	isSending bool

	*authoredTxData
	exchangeRateSet bool
}

func newSendConfirmModal(l *load.Load, data *authoredTxData) *sendConfirmModal {
	scm := &sendConfirmModal{
		Load:  l,
		Modal: l.Theme.ModalFloatTitle("send_confirm_modal"),

		authoredTxData: data,
	}

	scm.closeConfirmationModalButton = l.Theme.OutlineButton(values.String(values.StrCancel))
	scm.closeConfirmationModalButton.Font.Weight = text.Medium

	scm.confirmButton = l.Theme.Button("")
	scm.confirmButton.Font.Weight = text.Medium
	scm.confirmButton.SetEnabled(false)

	scm.passwordEditor = l.Theme.EditorPassword(new(widget.Editor), values.String(values.StrSpendingPassword))
	scm.passwordEditor.Editor.SetText("")
	scm.passwordEditor.Editor.SingleLine = true
	scm.passwordEditor.Editor.Submit = true

	return scm
}

func (scm *sendConfirmModal) OnResume() {
	scm.passwordEditor.Editor.Focus()
}

func (scm *sendConfirmModal) OnDismiss() {}

func (scm *sendConfirmModal) broadcastTransaction() {
	password := scm.passwordEditor.Editor.Text()
	if password == "" || scm.isSending {
		return
	}

	scm.isSending = true
	scm.Modal.SetDisabled(true)
	go func() {
		_, err := scm.authoredTxData.txAuthor.Broadcast([]byte(password))
		scm.isSending = false
		scm.Modal.SetDisabled(false)
		if err != nil {
			errModal := modal.NewErrorModal(scm.Load, err.Error(), func(isChecked bool) bool {
				return true
			})
			scm.ParentWindow().ShowModal(errModal)
			return
		}
		successModal := modal.NewSuccessModal(scm.Load, values.String(values.StrTxSent), func(isChecked bool) bool {
			return true
		})
		scm.ParentWindow().ShowModal(successModal)

		scm.txSent()
		scm.Dismiss()
	}()
}

func (scm *sendConfirmModal) Handle() {
	for _, evt := range scm.passwordEditor.Editor.Events() {
		if scm.passwordEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				scm.confirmButton.SetEnabled(scm.passwordEditor.Editor.Text() != "")
			case widget.SubmitEvent:
				scm.broadcastTransaction()
			}
		}
	}

	for scm.confirmButton.Clicked() {
		scm.broadcastTransaction()
	}

	for scm.closeConfirmationModalButton.Clicked() {
		if !scm.isSending {
			scm.Dismiss()
		}
	}
}

func (scm *sendConfirmModal) Layout(gtx layout.Context) D {

	w := []layout.Widget{
		func(gtx C) D {
			return scm.Theme.H6(values.String(values.StrConfirmSend)).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := scm.Theme.Icons.SendIcon
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding2, Right: values.MarginPadding16}.Layout(gtx, icon.Layout24dp)
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
											txt.Color = scm.Theme.Color.GrayText2
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
							icon := decredmaterial.NewIcon(scm.Theme.Icons.NavigationArrowForward)
							icon.Color = scm.Theme.Color.Gray1
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
											card.Radius = decredmaterial.Radius(0)
											card.Color = scm.Theme.Color.Gray4
											inset := layout.Inset{
												Left: values.MarginPadding5,
											}
											return inset.Layout(gtx, func(gtx C) D {
												return card.Layout(gtx, func(gtx C) D {
													return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
														destinationWallet := scm.WL.MultiWallet.WalletWithID(scm.destinationAccount.WalletID)
														txt := scm.Theme.Caption(destinationWallet.Name)
														txt.Color = scm.Theme.Color.GrayText1
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
					return scm.contentRow(gtx, values.String(values.StrSendingFrom), scm.sourceAccount.Name, sendWallet.Name)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						txFeeText := scm.txFee
						if scm.exchangeRateSet {
							txFeeText = fmt.Sprintf("%s (%s)", scm.txFee, scm.txFeeUSD)
						}

						return scm.contentRow(gtx, values.String(values.StrFee), txFeeText, "")
					})
				}),
				layout.Rigid(func(gtx C) D {
					totalCostText := scm.totalCost
					if scm.exchangeRateSet {
						totalCostText = fmt.Sprintf("%s (%s)", scm.totalCost, scm.totalCostUSD)
					}

					return scm.contentRow(gtx, values.String(values.StrTotalCost), totalCostText, "")
				}),
			)
		},
		func(gtx C) D {
			return scm.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := decredmaterial.NewIcon(scm.Theme.Icons.ActionInfo)
					icon.Color = scm.Theme.Color.Gray1
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return icon.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := scm.Theme.Body2(values.String(values.StrSendWarning))
					txt.Color = scm.Theme.Color.GrayText1
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
							if scm.isSending {
								return D{}
							}
							return scm.closeConfirmationModalButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if scm.isSending {
							return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx C) D {
								return material.Loader(scm.Theme.Base).Layout(gtx)
							})
						}
						scm.confirmButton.Text = fmt.Sprintf("%s %s", values.String(values.StrSend), scm.totalCost)
						return scm.confirmButton.Layout(gtx)
					}),
				)
			})
		},
	}

	return scm.Modal.Layout(gtx, w)
}

func (scm *sendConfirmModal) contentRow(gtx layout.Context, leftValue, rightValue, walletName string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := scm.Theme.Body2(leftValue)
			txt.Color = scm.Theme.Color.GrayText2
			return txt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(scm.Theme.Body1(rightValue).Layout),
					layout.Rigid(func(gtx C) D {
						if walletName != "" {
							card := scm.Theme.Card()
							card.Radius = decredmaterial.Radius(0)
							card.Color = scm.Theme.Color.Gray4
							inset := layout.Inset{
								Left: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := scm.Theme.Caption(walletName)
										txt.Color = scm.Theme.Color.GrayText2
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
