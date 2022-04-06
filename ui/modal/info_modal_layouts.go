package modal

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	VerifyMessageInfoTemplate      = "VerifyMessageInfo"
	SignMessageInfoTemplate        = "SignMessageInfo"
	PrivacyInfoTemplate            = "PrivacyInfo"
	SetupMixerInfoTemplate         = "ConfirmSetupMixer"
	TransactionDetailsInfoTemplate = "TransactionDetailsInfoInfo"
	WalletBackupInfoTemplate       = "WalletBackupInfo"
	AllowUnmixedSpendingTemplate   = "AllowUnmixedSpending"
	TicketPriceErrorTemplate       = "TicketPriceError"
)

func verifyMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := `<span style="text-color: gray">
						After you or your counterparty has genrated a signature, you can use this form to verify the
				 		validity of the signature.
				 		<br /> Once you have entered the address, the message and the corresponding signature, you will see <font color="success">VALID</font> 
						if the signature appropriately matches the address and message, otherwise <font color="danger">INVALID</font>.
					</span>`

			return renderers.RenderHTML(text, th).Layout(gtx)
		},
	}
}

func signMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	text := `<span style="text-color: gray">
				Signing a message with an address' private key allows you to prove that 
				you are the owner of a given address  to a possible counterparty.
			</span>`

	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func privacyInfo(l *load.Load) []layout.Widget {
	ic := decredmaterial.NewIcon(l.Theme.Icons.ImageBrightness1)
	ic.Color = l.Theme.Color.Gray1
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := l.Theme.Body1("When the mixer is activated, funds will be gradually transfered from the unmixed account to the mixed account.")
					text.Color = l.Theme.Color.GrayText2
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			txt := l.Theme.Body1("Important: keep this app open while mixer is running.")
			txt.Font.Weight = text.SemiBold
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := l.Theme.Body1("The mixer routine will automatically stop when the unmixed balance is fully mixed.")
					text.Color = l.Theme.Color.GrayText2
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
	}
}

func setupMixerInfo(th *decredmaterial.Theme) []layout.Widget {
	text := `<span style="text-color: grayText2">
				Two dedicated accounts ("mixed" & "unmixed") will be created in order to use the mixer. <br></br>
				<b>This action cannot be undone.</b>
			</span>`

	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func transactionDetailsInfo(th *decredmaterial.Theme) []layout.Widget {
	text := `<span style="text-color: grayText2">Tap on <span style="text-color: primary">blue text </span> to copy the item</span>`

	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func backupInfo(th *decredmaterial.Theme) []layout.Widget {
	textGray := "Please backup your seed words and keep them in a safe place in order to recover your funds if your device gets lost or broken."
	textDanger := "Anyone who has your seed words can spend your funds! Do not share them."

	return []layout.Widget{
		func(gtx C) D {
			txt := th.Label(values.TextSize16, textGray)
			txt.Color = th.Color.GrayText1
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			txt := th.Label(values.TextSize16, textDanger)
			txt.Color = th.Color.Danger
			txt.Font.Weight = text.Medium
			return layout.Inset{Top: values.MarginPaddingMinus15}.Layout(gtx, txt.Layout)
		},
	}
}

func allowUnspendUnmixedAcct(l *load.Load) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := decredmaterial.NewIcon(l.Theme.Icons.ActionInfo)
					ic.Color = l.Theme.Color.GrayText1
					return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						return ic.Layout(gtx, values.MarginPadding18)
					})
				}),
				layout.Rigid(func(gtx C) D {
					text := l.Theme.Body1("Spendings from unmixed accounts could potentially be traced back to you")
					text.Color = l.Theme.Color.GrayText1
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			text := `<span style="text-color: grayText1">
					Please type "<span style="font-weight: bold">I understand the risks</span>
					" to allow spending from unmixed accounts.
			</span>`
			return renderers.RenderHTML(text, l.Theme).Layout(gtx)
		},
	}
}

func ticketPriceErrorInfo(l *load.Load) []layout.Widget {
	col := l.Theme.Color.GrayText2
	return []layout.Widget{
		func(gtx C) D {
			txt := "wallet needs"
			if l.WL.MultiWallet.LoadedWalletsCount() > 1 {
				txt = "wallets need"
			}

			text := l.Theme.Body1(fmt.Sprintf("Your %s to be synced before some of the staking functionality will be available.", txt))
			text.Color = col
			return text.Layout(gtx)
		},
		func(gtx C) D {
			bestBlock := l.WL.MultiWallet.GetBestBlock()
			activationHeight := l.WL.MultiWallet.DCP0001ActivationBlockHeight()
			txt := l.Theme.Body1(fmt.Sprintf("The current sync progress is %v blocks, and the minimum required is %v blocks.", bestBlock.Height, activationHeight))
			txt.Font.Weight = text.SemiBold
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			text := l.Theme.Body1("Check the Overview page for more details on sync progress.")
			text.Color = col
			return text.Layout(gtx)
		},
	}
}
