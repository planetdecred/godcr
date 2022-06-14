package dexclient

import (
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const testDexHost = "dex-test.ssgen.io:7232"

type addDexModal struct {
	*load.Load
	*decredmaterial.Modal
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cancelBtn        decredmaterial.Button
	materialLoader   material.LoaderStyle
	cert             decredmaterial.Editor

	// dexClientPassword will be used if set, without requesting password input on UI.
	dexClientPassword     string
	appPassword           decredmaterial.Editor
	onDexAdded            func()
	knownServers          map[string]*decredmaterial.Clickable
	selectedServer        string
	pickServerFromListBtn decredmaterial.Button
	useCustomServerBtn    decredmaterial.Button
	isUseCustomServer     bool
}

func newAddDexModal(l *load.Load) *addDexModal {
	md := &addDexModal{
		Load:             l,
		Modal:            l.Theme.ModalFloatTitle("add_dex_modal"),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
		addDexServer:     l.Theme.Button("Submit"),
		cancel:           l.Theme.OutlineButton("Cancel"),
		materialLoader:   material.Loader(l.Theme.Base),
	}

	md := &addDexModal{
		Load:                  l,
		modal:                 l.Theme.ModalFloatTitle(),
		dexServerAddress:      l.Theme.Editor(&widget.Editor{Submit: true}, strDexAddr),
		cert:                  l.Theme.Editor(new(widget.Editor), strTLSCert),
		addDexServerBtn:       l.Theme.Button(strSubmit),
		cancelBtn:             l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:        material.Loader(l.Theme.Base),
		appPassword:           l.Theme.EditorPassword(&widget.Editor{Submit: true}, strAppPassword),
		pickServerFromListBtn: tabButton(strPickAServer, true),
		useCustomServerBtn:    tabButton(strCustomServer, false),
	}
	md.addDexServerBtn.SetEnabled(false)

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

	return md
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
}

func (md *addDexModal) OnDexAdded(callback func()) *addDexModal {
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
	if canSubmit && (md.addDexServerBtn.Clicked() || isSubmit) {
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

	if md.cancelBtn.Clicked() && !md.isSending {
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

		md.isSending = true
		md.Modal.SetDisabled(true)
		go func() {
			cert := []byte(md.cert.Editor.Text())
			dex, err := md.Dexc().DEXServerInfo(md.dexServerAddress.Editor.Text(), cert)
			md.isSending = false
			md.Modal.SetDisabled(false)

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
			}

			// Dismiss this modal before displaying a new one for adding a wallet
			// or completing the registration.
			md.Dismiss()
			if md.Dexc().HasWallet(int32(feeAsset.ID)) {
				md.completeRegistration(dex, feeAssetName, cert)
				return
			}

			createWalletModal := newCreateWalletModal(md.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(md.Load, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				},
				func() {
					md.completeRegistration(dex, feeAssetName, cert)
				})
			md.ParentWindow().ShowModal(createWalletModal)
		}()
	}
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

	return md.Modal.Layout(gtx, w)
}

func (md *addDexModal) serversLayout(gtx C) D {
	dexServers := sortExchanges(core.CertStore[md.Dexc().Core().Network()])
	serverWidgets := make([]layout.FlexChild, len(dexServers))

	for i := 0; i < len(dexServers); i++ {
		host := dexServers[i]
		serverWidgets[i] = layout.Rigid(func(gtx C) D {
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
							ic := decredmaterial.NewIcon(md.Theme.Icons.NavigationCheck)
							ic.Color = md.Theme.Color.Success
							return ic.Layout(gtx, values.MarginPadding20)
						}),
					)
				})
			})
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, serverWidgets...)
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

func payFeeAndRegister(l *load.Load, dexServer *core.Exchange, cert []byte, appPass string, onDexAdded func()) {
	// Create the assetSelectorModal now, it'll remain open/visible
	// until the fee is paid and registration is completed or the
	// user manually closes it.
	assetSelectorModal := newFeeAssetSelectorModal(l, dexServer)

	confirmAndRegister := func(feeAsset *core.SupportedAsset) {
		modal.NewInfoModal(l).
			Title(strConfirmReg).
			Body(confirmRegisterModalDesc(dexServer, feeAsset.Symbol)).
			SetCancelable(false).
			NegativeButton(values.String(values.StrCancel), assetSelectorModal.Show).
			PositiveButton(strRegister, func(_ bool) {
				assetSelectorModal.Show()
				go func() {
					assetSelectorModal.modal.SetDisabled(true) // prevent re-selecting a fee asset
					regFeeAsset := dexServer.RegFees[feeAsset.Symbol]
					_, err := l.Dexc().RegisterWithDEXServer(dexServer.Host,
						cert,
						int64(regFeeAsset.Amt),
						int32(regFeeAsset.ID),
						[]byte(appPass))
					if err != nil {
						assetSelectorModal.modal.SetDisabled(false) // re-enable fee asset selection
						assetSelectorModal.Toast.NotifyError(err.Error())
						return
					}
					assetSelectorModal.Dismiss()
					onDexAdded()
				}()
			}).Show()
	}

	assetSelectorModal.
		OnAssetSelected(func(asset *core.SupportedAsset) {
			if asset.Wallet != nil {
				confirmAndRegister(asset)
				return
			}

			feeAssetName := asset.Symbol
			newCreateWalletModal(l,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(l, feeAssetName),
					coinName: feeAssetName,
					coinID:   asset.ID,
				}, appPass).
				WalletCreated(func() {
					confirmAndRegister(asset)
				}).
				CancelClicked(assetSelectorModal.Show).
				SetRegisterAction(true).
				Show()
		}).
		Show()
}

func confirmRegisterModalDesc(dexServer *core.Exchange, selectedFeeAsset string) string {
	feeAsset := dexServer.RegFees[selectedFeeAsset]
	feeAmt := formatAmountUnit(feeAsset.ID, selectedFeeAsset, feeAsset.Amt)
	return fmt.Sprintf("Confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
}
