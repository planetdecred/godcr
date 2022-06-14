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
	addDexServerBtn  decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cancelBtn        decredmaterial.Button
	materialLoader   material.LoaderStyle
	cert             decredmaterial.Editor

	onDexAdded            func()
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
		Modal:                 l.Theme.ModalFloatTitle("add_dex_modal"),
		dexServerAddress:      l.Theme.Editor(&widget.Editor{Submit: true}, strDexAddr),
		cert:                  l.Theme.Editor(new(widget.Editor), strTLSCert),
		addDexServerBtn:       l.Theme.Button(strSubmit),
		cancelBtn:             l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:        material.Loader(l.Theme.Base),
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

func (md *addDexModal) OnDismiss() {}

func (md *addDexModal) OnResume() {}

func (md *addDexModal) OnDexAdded(callback func()) *addDexModal {
	md.onDexAdded = callback
	return md
}

func (md *addDexModal) validateInputs() (bool, string) {
	if md.isSending {
		return false, ""
	}

	dexServer := md.selectedServer
	if md.isUseCustomServer {
		dexServer = md.dexServerAddress.Editor.Text()
	}

	if dexServer == "" {
		md.addDexServerBtn.SetEnabled(false)
		return false, ""
	}

	md.addDexServerBtn.SetEnabled(true)
	return true, dexServer
}

func (md *addDexModal) Handle() {
	canSubmit, dexServer := md.validateInputs()

	if isDexSubmit, _ := decredmaterial.HandleEditorEvents(md.dexServerAddress.Editor); isDexSubmit && canSubmit {
		md.doAddDexServer(dexServer)
	}

	if canSubmit && md.addDexServerBtn.Clicked() {
		md.doAddDexServer(dexServer)
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

func (md *addDexModal) doAddDexServer(serverAddr string) {
	if md.isSending {
		return
	}

	md.isSending = true
	md.Modal.SetDisabled(true)
	go func() {
		defer func() {
			md.isSending = false
			md.Modal.SetDisabled(false)
		}()

		var cert []byte
		if md.isUseCustomServer {
			cert = []byte(md.cert.Editor.Text())
		}

		dexServer, paid, err := md.Dexc().Core().DiscoverAccount(serverAddr, []byte(DEXClientPass), cert)
		if err != nil {
			md.Toast.NotifyError(err.Error())
			return
		}

		md.Dismiss()
		if paid {
			md.onDexAdded()
			return
		}

		md.payFeeAndRegister(dexServer, cert)
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
						layout.Rigid(md.Theme.Label(values.TextSize14, host).Layout),
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

func (md *addDexModal) payFeeAndRegister(dexServer *core.Exchange, cert []byte) {
	// Create the assetSelectorModal now, it'll remain open/visible
	// until the fee is paid and registration is completed or the
	// user manually closes it.
	assetSelectorModal := newFeeAssetSelectorModal(md.Load, dexServer)

	confirmAndRegister := func(feeAsset *core.SupportedAsset) {
		infoModal := modal.NewInfoModal(md.Load).
			Title(strConfirmReg).
			Body(confirmRegisterModalDesc(dexServer, feeAsset.Symbol)).
			SetCancelable(false).
			NegativeButton(values.String(values.StrCancel), func() {
				md.ParentWindow().ShowModal(assetSelectorModal)
			}).
			PositiveButton(strRegister, func(_ bool) bool {
				md.ParentWindow().ShowModal(assetSelectorModal)
				go func() {
					assetSelectorModal.SetLoading(true)
					assetSelectorModal.Modal.SetDisabled(true) // prevent re-selecting a fee asset
					regFeeAsset := dexServer.RegFees[feeAsset.Symbol]
					_, err := md.Load.Dexc().RegisterWithDEXServer(dexServer.Host,
						cert,
						int64(regFeeAsset.Amt),
						int32(regFeeAsset.ID),
						[]byte(DEXClientPass))
					if err != nil {
						assetSelectorModal.SetLoading(false)
						assetSelectorModal.Modal.SetDisabled(false) // re-enable fee asset selection
						assetSelectorModal.Toast.NotifyError(err.Error())
						return
					}
					assetSelectorModal.Dismiss()
					md.onDexAdded()
				}()
				return true
			})

		md.ParentWindow().ShowModal(infoModal)
	}

	assetSelectorModal.
		OnAssetSelected(func(asset *core.SupportedAsset) {
			if asset.Wallet != nil {
				confirmAndRegister(asset)
				return
			}

			feeAssetName := asset.Symbol
			createWalletModal := newCreateWalletModal(md.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(md.Load, feeAssetName),
					coinName: feeAssetName,
					coinID:   asset.ID,
				}).
				WalletCreated(func() {
					confirmAndRegister(asset)
				}).
				CancelClicked(func() {
					md.ParentWindow().ShowModal(assetSelectorModal)
				}).
				SetRegisterAction(true)

			md.ParentWindow().ShowModal(createWalletModal)
		})

	md.ParentWindow().ShowModal(assetSelectorModal)
}

func confirmRegisterModalDesc(dexServer *core.Exchange, selectedFeeAsset string) string {
	feeAsset := dexServer.RegFees[selectedFeeAsset]
	feeAmt := formatAmountUnit(feeAsset.ID, selectedFeeAsset, feeAsset.Amt)
	return fmt.Sprintf("Confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
}
