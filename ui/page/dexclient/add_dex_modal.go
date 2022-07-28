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

type AddDexModal struct {
	*load.Load
	*decredmaterial.Modal
	addDexServerBtn  decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cancelBtn        decredmaterial.Button
	materialLoader   material.LoaderStyle
	cert             decredmaterial.Editor

	onDexAdded func()
}

func NewAddDexModal(l *load.Load) *AddDexModal {
	md := &AddDexModal{
		Load:             l,
		Modal:            l.Theme.ModalFloatTitle("add_dex_modal"),
		dexServerAddress: l.Theme.Editor(&widget.Editor{Submit: true}, strDexAddr),
		cert:             l.Theme.Editor(new(widget.Editor), strTLSCert),
		addDexServerBtn:  l.Theme.Button(values.String(values.StrContinue)),
		cancelBtn:        l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:   material.Loader(l.Theme.Base),
	}
	md.addDexServerBtn.SetEnabled(false)

	return md
}

func (md *AddDexModal) OnDismiss() {}

func (md *AddDexModal) OnResume() {}

func (md *AddDexModal) OnDexAdded(callback func()) *AddDexModal {
	md.onDexAdded = callback
	return md
}

func (md *AddDexModal) validateInputs() (bool, string) {
	if md.isSending {
		return false, ""
	}

	dexServer := md.dexServerAddress.Editor.Text()
	if dexServer == "" {
		md.addDexServerBtn.SetEnabled(false)
		return false, ""
	}

	md.addDexServerBtn.SetEnabled(true)
	return true, dexServer
}

func (md *AddDexModal) Handle() {
	canSubmit, dexServer := md.validateInputs()

	if isDexSubmit, _ := decredmaterial.HandleEditorEvents(md.dexServerAddress.Editor); isDexSubmit && canSubmit {
		md.doAddDexServer(dexServer)
	}

	if canSubmit && md.addDexServerBtn.Clicked() {
		md.doAddDexServer(dexServer)
	}

	if md.cancelBtn.Clicked() && !md.isSending {
		md.Dismiss()
	}
}

func (md *AddDexModal) doAddDexServer(serverAddr string) {
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

		cert := []byte(md.cert.Editor.Text())

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

func (md *AddDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		md.Load.Theme.Label(values.TextSize20, values.String(values.StrAddDexServer)).Layout,
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(md.dexServerAddress.Layout),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Max.Y = 300
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, md.cert.Layout)
				}),
			)
		},
		md.Theme.Separator().Layout,
		func(gtx C) D {
			customServerText := md.Theme.Label(values.TextSize16, strCustomServer)
			customServerText.Color = md.Theme.Color.Primary
			return customServerText.Layout(gtx)
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

func (md *AddDexModal) payFeeAndRegister(dexServer *core.Exchange, cert []byte) {
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
					md.saveDexServer(dexServer.Host, cert)
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

// saveDexServer after pay the fee success save the host and cert to db.
func (md *AddDexModal) saveDexServer(host string, cert []byte) {
	dexServer := new(components.DexServer)
	err := md.Load.WL.MultiWallet.ReadUserConfigValue(components.KnownDexServersConfigKey, &dexServer)
	if err != nil {
		return
	}
	if dexServer.SavedHosts == nil {
		dexServer.SavedHosts = make(map[string][]byte)
	}
	dexServer.SavedHosts[host] = cert
	md.Load.WL.MultiWallet.SaveUserConfigValue(components.KnownDexServersConfigKey, dexServer)
}
