package dexclient

import (
	"fmt"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

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
	cert             decredmaterial.Editor

	// dexClientPassword will be used if set, without requesting password input on UI.
	dexClientPassword     string
	appPassword           decredmaterial.Editor
	onDexAdded            func(*core.Exchange)
	knownServers          map[string]*decredmaterial.Clickable
	selectedServer        string
	pickServerFromListBtn decredmaterial.Button
	useCustomServerBtn    decredmaterial.Button
	isUseCustomServer     bool
}

func newAddDexModal(l *load.Load) *addDexModal {
	tabButton := func(text string, active bool) decredmaterial.Button {
		btn := l.Theme.OutlineButton(text)
		btn.CornerRadius = values.MarginPadding0
		btn.Inset = layout.Inset{
			Top:    values.MarginPadding5,
			Bottom: values.MarginPadding5,
			Left:   values.MarginPadding9,
			Right:  values.MarginPadding9,
		}
		btn.TextSize = values.TextSize14
		if !active {
			btn.Background = l.Theme.Color.Background
		}
		return btn
	}

	md := &addDexModal{
		Load:                  l,
		modal:                 l.Theme.ModalFloatTitle(),
		dexServerAddress:      l.Theme.Editor(&widget.Editor{Submit: true}, strDexAddr),
		cert:                  l.Theme.Editor(new(widget.Editor), strTLSCert),
		addDexServerBtn:       l.Theme.Button(strSubmit),
		cancelBtn:             l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:        material.Loader(material.NewTheme(gofont.Collection())),
		appPassword:           l.Theme.EditorPassword(&widget.Editor{Submit: true}, strAppPassword),
		pickServerFromListBtn: tabButton(strPickAServer, true),
		useCustomServerBtn:    tabButton(strCustomServer, false),
	}
	md.addDexServerBtn.SetEnabled(false)

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

func (md *addDexModal) WithAppPassword(appPass string) *addDexModal {
	md.dexClientPassword = appPass
	return md
}

func (md *addDexModal) OnResume() {
	md.appPassword.Editor.Focus()

	clickable := func() *decredmaterial.Clickable {
		cl := md.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	dexServers := sortExchanges(core.CertStore[md.Dexc().Core().Network()])
	md.knownServers = make(map[string]*decredmaterial.Clickable, len(dexServers))
	for _, server := range dexServers {
		md.knownServers[server] = clickable()
	}
	if len(dexServers) > 0 {
		md.selectedServer = dexServers[0]
	}
}

func (md *addDexModal) OnDexAdded(callback func(*core.Exchange)) *addDexModal {
	md.onDexAdded = callback
	return md
}

func (md *addDexModal) validateInputs() (bool, string, string) {
	if md.isSending {
		return false, "", ""
	}

	appPass := md.dexClientPassword
	if appPass == "" {
		appPass = md.appPassword.Editor.Text()
	}

	dexServer := md.selectedServer
	if md.isUseCustomServer {
		dexServer = md.dexServerAddress.Editor.Text()
	}

	if appPass == "" || dexServer == "" {
		md.addDexServerBtn.SetEnabled(false)
		return false, "", ""
	}

	md.addDexServerBtn.SetEnabled(true)
	return true, appPass, dexServer
}

func (md *addDexModal) Handle() {
	canSubmit, appPass, dexServer := md.validateInputs()

	if isDexServerSubmit, _ := decredmaterial.HandleEditorEvents(md.dexServerAddress.Editor); isDexServerSubmit {
		if canSubmit {
			md.doAddDexServer(dexServer, appPass)
		} else if md.dexClientPassword == "" {
			md.appPassword.Editor.Focus()
		}
	}

	isSubmit, _ := decredmaterial.HandleEditorEvents(md.appPassword.Editor)
	if canSubmit && (md.addDexServerBtn.Button.Clicked() || isSubmit) {
		md.doAddDexServer(dexServer, appPass)
	}

	if md.pickServerFromListBtn.Clicked() {
		md.isUseCustomServer = false
		md.pickServerFromListBtn.Background = md.Theme.Color.Surface
		md.useCustomServerBtn.Background = md.Theme.Color.Background
	}

	if md.useCustomServerBtn.Clicked() {
		md.isUseCustomServer = true
		md.dexServerAddress.Editor.Focus()
		md.useCustomServerBtn.Background = md.Theme.Color.Surface
		md.pickServerFromListBtn.Background = md.Theme.Color.Background
	}

	if md.cancelBtn.Button.Clicked() && !md.isSending {
		md.Dismiss()
	}

	for host, cl := range md.knownServers {
		if cl.Clicked() {
			if md.selectedServer == host {
				md.selectedServer = ""
				break
			}
			md.selectedServer = host
			break
		}
	}
}

func (md *addDexModal) doAddDexServer(serverAddr, appPass string) {
	if md.isSending {
		return
	}

	md.isSending = true
	md.modal.SetDisabled(true)
	go func() {
		defer func() {
			md.isSending = false
			md.modal.SetDisabled(false)
		}()

		var cert []byte
		if md.isUseCustomServer {
			cert = []byte(md.cert.Editor.Text())
		}

		dexServer, paid, err := md.Dexc().Core().DiscoverAccount(serverAddr, []byte(appPass), cert)
		if err != nil {
			md.Toast.NotifyError(err.Error())
			return
		}

		md.Dismiss()
		if paid {
			md.onDexAdded(dexServer)
			return
		}

		newFeeAssetSelectorModal(md.Load, dexServer).
			OnAssetSelected(func(asset *core.SupportedAsset) {
				cfReg := &confirmRegistration{
					Load:      md.Load,
					dexServer: dexServer,
					isSending: &md.isSending,
					Show:      md.Show,
					completed: md.onDexAdded,
					Dismiss:   md.Dismiss,
				}
				feeAssetName := asset.Symbol
				if asset.Wallet != nil {
					cfReg.confirm(feeAssetName, appPass, cert)
					return
				}
				newCreateWalletModal(md.Load,
					&walletInfoWidget{
						image:    components.CoinImageBySymbol(&md.Icons, feeAssetName),
						coinName: feeAssetName,
						coinID:   asset.ID,
					},
					appPass,
					func(wallModal *createWalletModal) {
						cfReg.isSending = &wallModal.isSending
						cfReg.Show = wallModal.Show
						cfReg.Dismiss = wallModal.Dismiss
						cfReg.confirm(feeAssetName, appPass, cert)
					}).
					SetRegisterAction(true).
					Show()
			}).
			Show()
	}()
}

func (md *addDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		md.Load.Theme.Label(values.TextSize20, strAddADex).Layout,
		func(gtx C) D {
			return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(.5, md.pickServerFromListBtn.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left:  values.MarginPadding1,
						Right: values.MarginPadding1,
					}.Layout(gtx, func(gtx C) D { return D{} })
				}),
				layout.Flexed(.5, md.useCustomServerBtn.Layout),
			)
		},
		func(gtx C) D {
			if md.isUseCustomServer {
				return md.customServerLayout(gtx)
			}
			return md.serversLayout(gtx)
		},
		md.Theme.Separator().Layout,
		func(gtx C) D {
			if md.dexClientPassword != "" {
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

func (md *addDexModal) serversLayout(gtx C) D {
	dexServers := sortExchanges(core.CertStore[md.Dexc().Core().Network()])
	var childrens = make([]layout.FlexChild, 0, len(dexServers))

	for i := 0; i < len(dexServers); i++ {
		host := dexServers[i]

		childrens = append(childrens, layout.Rigid(func(gtx C) D {
			return md.knownServers[host].Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Inset{
					Top:    values.MarginPadding8,
					Bottom: values.MarginPadding8,
					Left:   values.MarginPadding12,
					Right:  values.MarginPadding12,
				}.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.Y = 45
					return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(md.Theme.Label(values.MarginPadding14, host).Layout),
						layout.Rigid(func(gtx C) D {
							if md.selectedServer != host {
								return D{}
							}
							ic := decredmaterial.NewIcon(md.Icons.NavigationCheck)
							ic.Color = md.Theme.Color.Success
							return ic.Layout(gtx, values.MarginPadding20)
						}),
					)
				})
			})
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
}

func (md *addDexModal) customServerLayout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(md.dexServerAddress.Layout),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.Y = 300
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, md.cert.Layout)
		}),
	)
}

type confirmRegistration struct {
	dexServer *core.Exchange
	*load.Load
	isSending *bool
	Show      func()
	Dismiss   func()
	completed func(*core.Exchange)
}

func (cfReg *confirmRegistration) confirm(feeAssetName string, password string, cert []byte) {
	modal.NewInfoModal(cfReg.Load).
		Title(strConfirmReg).
		Body(confirmRegisterModalDesc(cfReg.dexServer, feeAssetName)).
		SetCancelable(false).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(strRegister, func() {
			go func() {
				// Show previous modal and display loading status or error messages
				cfReg.Show()
				*cfReg.isSending = true
				_, err := cfReg.Dexc().RegisterWithDEXServer(cfReg.dexServer.Host,
					cert,
					int64(cfReg.dexServer.Fee.Amt),
					int32(cfReg.dexServer.Fee.ID),
					[]byte(password))
				if err != nil {
					*cfReg.isSending = false
					cfReg.Toast.NotifyError(err.Error())
					return
				}
				cfReg.completed(cfReg.dexServer)
				cfReg.Dismiss()
			}()
		}).Show()
}

func confirmRegisterModalDesc(dexServer *core.Exchange, selectedFeeAsset string) string {
	feeAsset := dexServer.RegFees[selectedFeeAsset]
	feeAmt := formatAmountUnit(feeAsset.ID, selectedFeeAsset, feeAsset.Amt)
	txt := fmt.Sprintf("Confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
	markets := make([]string, 0, len(dexServer.Markets))
	for _, mkt := range dexServer.Markets {
		lotSize := formatAmountUnit(mkt.BaseID, mkt.BaseSymbol, mkt.LotSize)
		markets = append(markets, fmt.Sprintf("Base: %s\tQuote: %s\tLot Size: %s", strings.ToUpper(mkt.BaseSymbol), strings.ToUpper(mkt.QuoteSymbol), lotSize))
	}
	return fmt.Sprintf("%s\n\nThis DEX supports the following markets. All trades are in multiples of each market's lot size.\n\n%s", txt, strings.Join(markets, "\n"))
}
