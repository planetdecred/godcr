package dexclient

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/core"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ncruces/zenity"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const addDexModalID = "add_dex_modal"

type addDexModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	addDexServerBtn  decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cancelBtn        decredmaterial.Button
	materialLoader   material.LoaderStyle
	certFilePath     string
	fileSelectBtn    *decredmaterial.Clickable

	// appPass is the password value after login or initialize to continue processing add new DEX
	// the Add Dex Modal won't show password input on UI
	appPass      string
	appPassword  decredmaterial.Editor
	onDexCreated func(*core.Exchange)

	listServer        *widget.List
	listExchangeWdg   []*knownExchangeWidget
	selectedServer    string
	listServerBtn     decredmaterial.Button
	customServerBtn   decredmaterial.Button
	isUseCustomServer bool
}

type knownExchangeWidget struct {
	selectBtn *decredmaterial.Clickable
	host      string
}

func newAddDexModal(l *load.Load, appPass string) *addDexModal {
	clickable := func() *decredmaterial.Clickable {
		cl := l.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	md := &addDexModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServerBtn:  l.Theme.Button("Submit"),
		cancelBtn:        l.Theme.OutlineButton("Cancel"),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
		appPass:          appPass,
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App Password"),

		fileSelectBtn:   clickable(),
		listServerBtn:   l.Theme.OutlineButton("Pick a Server"),
		customServerBtn: l.Theme.OutlineButton("Custom Server"),
		listServer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	md.customServerBtn.Background = l.Theme.Color.Background
	md.listServerBtn.CornerRadius, md.customServerBtn.CornerRadius = values.MarginPadding0, values.MarginPadding0
	inset := layout.Inset{
		Top:    values.MarginPadding5,
		Bottom: values.MarginPadding5,
		Left:   values.MarginPadding9,
		Right:  values.MarginPadding9,
	}
	md.listServerBtn.Inset, md.customServerBtn.Inset = inset, inset
	md.listServerBtn.TextSize, md.customServerBtn.TextSize = values.TextSize14, values.TextSize14

	md.appPassword.Editor.SingleLine = true
	md.dexServerAddress.Editor.SingleLine = true

	return md
}

func (md *addDexModal) ModalID() string {
	return addDexModalID
}

func (md *addDexModal) Show() {
	md.ShowModal(md)
}

func (md *addDexModal) Dismiss() {
	md.DismissModal(md)
}

func (md *addDexModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *addDexModal) OnResume() {
	md.appPassword.Editor.Focus()

	clickable := func() *decredmaterial.Clickable {
		cl := md.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	// Initialize listExchangeWdg.
	certs := core.CertStore[md.Dexc().Core().Network()]
	md.listExchangeWdg = make([]*knownExchangeWidget, 0)
	for host := range certs {
		md.listExchangeWdg = append(md.listExchangeWdg, &knownExchangeWidget{
			host:      host,
			selectBtn: clickable(),
		})
	}

	if len(md.listExchangeWdg) > 0 {
		md.selectedServer = md.listExchangeWdg[0].host
	}
}

func (md *addDexModal) DexCreated(callback func(*core.Exchange)) *addDexModal {
	md.onDexCreated = callback
	return md
}

func (md *addDexModal) Handle() {
	if md.listServerBtn.Clicked() {
		md.isUseCustomServer = false
		md.listServerBtn.Background = md.Theme.Color.Surface
		md.customServerBtn.Background = md.Theme.Color.Background
	}

	if md.customServerBtn.Clicked() {
		md.isUseCustomServer = true
		md.dexServerAddress.Editor.Focus()
		md.customServerBtn.Background = md.Theme.Color.Surface
		md.listServerBtn.Background = md.Theme.Color.Background
	}

	if md.cancelBtn.Button.Clicked() && !md.isSending {
		md.Dismiss()
	}

	if md.fileSelectBtn.Clicked() {
		filePath, err := zenity.SelectFile(
			zenity.Title("Select Cert File"),
			zenity.FileFilter{
				Name:     "Cert file",
				Patterns: []string{"*.cert"},
			},
		)

		if err != nil {
			md.Toast.NotifyError(err.Error())
			return
		}

		md.certFilePath = filePath
	}

	for _, eWdg := range md.listExchangeWdg {
		if eWdg.selectBtn.Clicked() {
			if md.selectedServer == eWdg.host {
				md.selectedServer = ""
				break
			}
			md.selectedServer = eWdg.host
			break
		}
	}

	if md.addDexServerBtn.Button.Clicked() {
		if md.isSending {
			return
		}

		md.isSending = true
		md.modal.SetDisabled(true)
		go func() {
			var cert []byte
			serverAddr := md.selectedServer

			if md.isUseCustomServer {
				serverAddr = md.dexServerAddress.Editor.Text()
				c, err := getCertFromFile(md.certFilePath)
				if err != nil {
					md.Toast.NotifyError(err.Error())
					return
				}
				cert = c
			}

			appPass := md.appPass
			if appPass == "" {
				appPass = md.appPassword.Editor.Text()
			}

			if serverAddr == "" {
				md.Toast.NotifyError("Please choose a server address or set a custom server")
				return
			}

			if appPass == "" {
				md.Toast.NotifyError("Please input your application password")
				return
			}

			md.isSending = true
			dex, paid, err := md.Dexc().Core().DiscoverAccount(serverAddr, []byte(appPass), cert)
			md.isSending = false
			md.modal.SetDisabled(false)

			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			if paid {
				md.onDexCreated(dex)
				md.Dismiss()
				return
			}

			cfRegistration := &confirmRegistration{
				Load:      md.Load,
				isSending: &md.isSending,
				Show:      md.Show,
				completed: md.onDexCreated,
				Dismiss:   md.Dismiss,
			}

			// Ensure a wallet is connected that can be used to pay the fees.
			// TODO: This automatically selects the dcr wallet if the DEX
			// supports it for fee payment, otherwise picks a random wallet
			// to use for fee payment. Should instead update the modal UI
			// to show the options and let the user choose which wallet to
			// set up and use for fee payment.
			feeAssetName := "dcr"
			feeAsset := dex.RegFees[feeAssetName]
			if feeAsset == nil {
				for feeAssetName, feeAsset = range dex.RegFees {
					break
				}
			}

			// Dismiss this modal before displaying a new one for adding a wallet
			// or completing the registration.
			md.Dismiss()
			if md.Dexc().HasWallet(int32(feeAsset.ID)) {
				cfRegistration.confirm(dex, feeAssetName, appPass, cert)
				return
			}

			newCreateWalletModal(md.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(&md.Load.Icons, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				},
				appPass,
				func(wallModal *createWalletModal) {
					cfRegistration.isSending = &wallModal.isSending
					cfRegistration.Show = wallModal.Show
					cfRegistration.Dismiss = wallModal.Dismiss
					cfRegistration.confirm(dex, feeAssetName, appPass, cert)
				}).Show()
		}()
	}
}

func (md *addDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		md.Load.Theme.Label(values.TextSize20, "Add a dex").Layout,
		func(gtx C) D {
			return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(.5, md.listServerBtn.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left:  values.MarginPadding1,
						Right: values.MarginPadding1,
					}.Layout(gtx, func(gtx C) D { return D{} })
				}),
				layout.Flexed(.5, md.customServerBtn.Layout),
			)
		},
		func(gtx C) D {
			if md.isUseCustomServer {
				return md.customServerLayout(gtx)
			}
			return md.listServerLayout(gtx)
		},
		md.Theme.Separator().Layout,
		func(gtx C) D {
			if md.appPass != "" {
				return D{}
			}
			return md.appPassword.Layout(gtx)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return D{}
						}
						return layout.Inset{
							Right:  values.MarginPadding4,
							Bottom: values.MarginPadding15,
						}.Layout(gtx, md.cancelBtn.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, md.materialLoader.Layout)
						}
						return md.addDexServerBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.modal.Layout(gtx, w)
}

func (md *addDexModal) listServerLayout(gtx C) D {
	return md.Theme.List(md.listServer).Layout(gtx, len(md.listExchangeWdg), func(gtx C, i int) D {
		return md.listExchangeWdg[i].selectBtn.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{
				Top:    values.MarginPadding8,
				Bottom: values.MarginPadding8,
				Left:   values.MarginPadding12,
				Right:  values.MarginPadding12,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(md.Theme.Label(values.MarginPadding14, md.listExchangeWdg[i].host).Layout),
					layout.Rigid(func(gtx C) D {
						if md.selectedServer != md.listExchangeWdg[i].host {
							return D{}
						}
						gtx.Constraints.Min.X = 30
						ic := md.Load.Icons.NavigationCheck
						return ic.Layout(gtx, md.Theme.Color.Success)
					}),
				)
			})
		})
	})
}

func (md *addDexModal) customServerLayout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(md.dexServerAddress.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, md.Theme.Label(values.MarginPadding16, "TLS Certificate").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					fileName := "None file selected"
					if md.certFilePath != "" {
						fileName = filepath.Base(md.certFilePath)
					}
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, md.Theme.Label(values.MarginPadding16, fileName).Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return widget.Border{
						Color:        md.Theme.Color.Gray2,
						CornerRadius: values.MarginPadding4,
						Width:        values.MarginPadding1,
					}.Layout(gtx, func(gtx C) D {
						labelBtn := "add a file"
						if md.certFilePath != "" {
							labelBtn = "Choose another file"
						}
						return md.fileSelectBtn.Layout(gtx, func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding4,
								Bottom: values.MarginPadding4,
								Left:   values.MarginPadding10,
								Right:  values.MarginPadding10,
							}.Layout(gtx, md.Theme.Label(values.MarginPadding14, labelBtn).Layout)
						})
					})
				}),
			)
		}),
	)
}

type confirmRegistration struct {
	*load.Load
	isSending *bool
	Show      func()
	Dismiss   func()
	completed func(*core.Exchange)
}

func (cf *confirmRegistration) confirm(dex *core.Exchange, feeAssetName string, password string, cert []byte) {
	modal.NewInfoModal(cf.Load).
		Title("Confirm Registration").
		Body(confirmRegisterModalDesc(dex, feeAssetName)).
		SetCancelable(false).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Register", func() {
			go func() {
				// Show previous modal and display loading status or error messages
				cf.Show()
				*cf.isSending = true
				_, err := cf.Dexc().RegisterWithDEXServer(dex.Host,
					cert,
					int64(dex.Fee.Amt),
					int32(dex.Fee.ID),
					[]byte(password))
				if err != nil {
					*cf.isSending = false
					cf.Toast.NotifyError(err.Error())
					return
				}
				*cf.isSending = false
				cf.completed(dex)
				cf.Dismiss()
			}()
		}).Show()
}

func confirmRegisterModalDesc(dex *core.Exchange, selectedFeeAsset string) string {
	feeAsset := dex.RegFees[selectedFeeAsset]
	feeAmt := formatAmount(feeAsset.ID, selectedFeeAsset, feeAsset.Amt)
	txt := fmt.Sprintf("Confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
	markets := make([]string, 0, len(dex.Markets))
	for _, mkt := range dex.Markets {
		lotSize := formatAmount(mkt.BaseID, mkt.BaseSymbol, mkt.LotSize)
		markets = append(markets, fmt.Sprintf("Base: %s\tQuote: %s\tLot Size: %s", strings.ToUpper(mkt.BaseSymbol), strings.ToUpper(mkt.QuoteSymbol), lotSize))
	}
	return fmt.Sprintf("%s\n\nThis DEX supports the following markets. All trades are in multiples of each market's lot size.\n\n%s", txt, strings.Join(markets, "\n"))
}

func formatAmount(assetID uint32, assetName string, amount uint64) string {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Sprintf("%d [%s units]", amount, assetName)
	}
	unitInfo := assetInfo.UnitInfo
	convertedLotSize := float64(amount) / float64(unitInfo.Conventional.ConversionFactor)
	return fmt.Sprintf("%s %s", strconv.FormatFloat(convertedLotSize, 'f', -1, 64), unitInfo.Conventional.Unit)
}
