package modal

import (
	"gioui.org/layout"

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
	SecurityToolsInfoTemplate      = "SecurityToolsInfo"
	RemoveWalletInfoTemplate       = "RemoveWalletInfo"
)

func verifyMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrVerifyMessageInfo, `<span style="text-color: gray">`, `<br />`, `<font color="success">`, `</font>`, `<font color="danger">`, `</font>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func signMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrSignMessageInfo, `<span style="text-color: gray">`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func securityToolsInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrSecurityToolsInfo, `<span style="text-color: gray">`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func privacyInfo(l *load.Load) []layout.Widget {
	text := values.StringF(values.StrPrivacyInfo, `<span style="text-color: gray">`, `<br/><span style="font-weight: bold">`, `</span></br>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, l.Theme).Layout,
	}
}

func setupMixerInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrSetupMixerInfo, `<span style="text-color: gray">`, `<span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span>`, `<br> <span style="font-weight: bold">`, `</span></span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func transactionDetailsInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrTxdetailsInfo, `<span style="text-color: gray">`, `<span style="text-color: primary">`, `</span>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func backupInfo(th *decredmaterial.Theme) []layout.Widget {
	text := values.StringF(values.StrBackupInfo, `<span style="text-color: danger"> <span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span></span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func allowUnspendUnmixedAcct(l *load.Load) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := decredmaterial.NewIcon(l.Theme.Icons.ActionInfo)
					ic.Color = l.Theme.Color.GrayText1
					return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						return ic.Layout(gtx, values.MarginPadding18)
					})
				}),
				layout.Rigid(func(gtx C) D {
					text := values.StringF(values.StrAllowUnspendUnmixedAcct, `<span style="text-color: gray">`, `<br>`, `<span style="font-weight: bold">`, `</span>`, `</span>`)
					return renderers.RenderHTML(text, l.Theme).Layout(gtx)
				}),
			)
		},
	}
}

func removeWalletInfo(l *load.Load) []layout.Widget {
	text := values.StringF(values.StrRemoveWalletInfo, `<span style="text-color: gray">`, `<span style="font-weight: bold">`, l.WL.SelectedWallet.Wallet.Name, `</span>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, l.Theme).Layout,
	}
}
