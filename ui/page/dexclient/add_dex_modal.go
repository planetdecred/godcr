package dexclient

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const addDexModalID = "add_dex_modal"

const testDexHost = "dex-test.ssgen.io:7232"

type addDexModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cert             decredmaterial.Editor
	done             func()
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newAddDexModal(l *load.Load) *addDexModal {
	md := &addDexModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServer:     l.Theme.Button("Submit"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
	}

	md.addDexServer.TextSize = values.TextSize12
	md.addDexServer.Background = l.Theme.Color.Primary

	md.dexServerAddress.Editor.SetText(testDexHost)
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
	md.dexServerAddress.Editor.SetText("")
}

func (md *addDexModal) OnResume() {
	md.dexServerAddress.Editor.Focus()
}

func (md *addDexModal) Handle() {
	if md.addDexServer.Button.Clicked() {
		if md.dexServerAddress.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			cert := []byte(md.cert.Editor.Text())
			// TODO: Use DiscoverAccount instead of GetDEXConfig to enable account
			// recovery without re-paying the fee. This is only relevant when the
			// dex client supports restoring from seed. Requires a field for app
			// password.
			dex, err := md.Dexc.GetDEXConfig(md.dexServerAddress.Editor.Text(), cert)
			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
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

			completeRegistration := func() {
				cfModal := newConfirmRegisterModal(md.Load, dex, cert, feeAssetName)
				cfModal.confirmed = func() {
					md.done()
				}
				cfModal.Show()
			}

			// Dismiss this modal before displaying a new one for adding a wallet
			// or completing the registration.
			md.Dismiss()

			if md.Dexc.WalletState(feeAsset.ID) != nil {
				completeRegistration()
			} else {
				newWalletModal := newCreateWalletModal(md.Load, &walletInfoWidget{
					image:    coinImageBySymbol(&md.Load.Icons, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				})
				newWalletModal.walletCreated = func() {
					completeRegistration()
				}
				newWalletModal.Show()
			}
		}()
	}
}

func (md *addDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return md.dexServerAddress.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.cert.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return md.addDexServer.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w)
}

// coinImageBySymbol returns image widget for supported asset coins.
func coinImageBySymbol(icons *load.Icons, coinName string) *decredmaterial.Image {
	m := map[string]*decredmaterial.Image{
		"btc": icons.BTC,
		"dcr": icons.DCR,
		"bch": icons.BCH,
		"ltc": icons.LTC,
	}
	coin, ok := m[strings.ToLower(coinName)]

	if !ok {
		return nil
	}

	return coin
}
