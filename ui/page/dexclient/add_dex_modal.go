package dexclient

import (
	"fmt"
	"path/filepath"
	"strings"

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

	// defaultAppPass is the password value after login or initialize to continue processing add new DEX
	// the Add Dex Modal won't show password input on UI
	defaultAppPass string
	appPassword    decredmaterial.Editor
	onDexAdded     func(*core.Exchange)

	listServerClickable map[string]*decredmaterial.Clickable
	selectedServer      string
	listServerBtn       decredmaterial.Button
	customServerBtn     decredmaterial.Button
	isUseCustomServer   bool
}

func newAddDexModal(l *load.Load) *addDexModal {
	clickable := func() *decredmaterial.Clickable {
		cl := l.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

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
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(&widget.Editor{Submit: true}, strDexAddr),
		addDexServerBtn:  l.Theme.Button(strSubmit),
		cancelBtn:        l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
		appPassword:      l.Theme.EditorPassword(&widget.Editor{Submit: true}, strAppPassword),
		fileSelectBtn:    clickable(),
		listServerBtn:    tabButton(strPickAServer, true),
		customServerBtn:  tabButton(strCustomServer, false),
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
	md.defaultAppPass = appPass
	return md
}

func (md *addDexModal) OnResume() {
	md.appPassword.Editor.Focus()

	clickable := func() *decredmaterial.Clickable {
		cl := md.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	// Initialize listExchangeWdg.
	listServer := sliceSever(core.CertStore[md.Dexc().Core().Network()])
	md.listServerClickable = make(map[string]*decredmaterial.Clickable, len(listServer))
	for i := 0; i < len(listServer); i++ {
		md.listServerClickable[listServer[i]] = clickable()
	}
	if len(listServer) > 0 {
		md.selectedServer = listServer[0]
	}
}

func (md *addDexModal) OnDexAdded(callback func(*core.Exchange)) *addDexModal {
	md.onDexAdded = callback
	return md
}

func (md *addDexModal) validateInputs() (bool, string, string) {
	appPass := md.defaultAppPass
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
		if md.defaultAppPass != "" && canSubmit {
			md.doAddDexServer(dexServer, appPass)
		} else {
			md.appPassword.Editor.Focus()
		}
	}

	isSubmit, _ := decredmaterial.HandleEditorEvents(md.appPassword.Editor)
	if canSubmit && (md.addDexServerBtn.Button.Clicked() || isSubmit) {
		md.doAddDexServer(dexServer, appPass)
	}

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

	for host, cl := range md.listServerClickable {
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
			c, err := getCertFromFile(md.certFilePath)
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}
			cert = c
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

		newAssetSelectorModal(md.Load, dexServer).
			OnAssetSelected(func(asset *core.SupportedAsset) {
				cfReg := &confirmRegistration{
					Load:      md.Load,
					Exchange:  dexServer,
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
			if md.defaultAppPass != "" {
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
	listServer := sliceSever(core.CertStore[md.Dexc().Core().Network()])
	var childrens = make([]layout.FlexChild, 0, len(listServer))

	for i := 0; i < len(listServer); i++ {
		host := listServer[i]

		childrens = append(childrens, layout.Rigid(func(gtx C) D {
			return md.listServerClickable[host].Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Inset{
					Top:    values.MarginPadding8,
					Bottom: values.MarginPadding8,
					Left:   values.MarginPadding12,
					Right:  values.MarginPadding12,
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(md.Theme.Label(values.MarginPadding14, host).Layout),
						layout.Rigid(func(gtx C) D {
							if md.selectedServer != host {
								return D{}
							}
							gtx.Constraints.Min.X = 30
							ic := md.Load.Icons.NavigationCheck
							return ic.Layout(gtx, md.Theme.Color.Success)
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
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, md.Theme.Label(values.MarginPadding16, strTLSCert).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					fileName := strNoneFileSelect
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
						labelBtn := strAddAFile
						if md.certFilePath != "" {
							labelBtn = strChooseOtherFile
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
	*core.Exchange
	*load.Load
	isSending *bool
	Show      func()
	Dismiss   func()
	completed func(*core.Exchange)
}

func (cfReg *confirmRegistration) confirm(feeAssetName string, password string, cert []byte) {
	modal.NewInfoModal(cfReg.Load).
		Title(strConfirmReg).
		Body(confirmRegisterModalDesc(cfReg.Exchange, feeAssetName)).
		SetCancelable(false).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(strRegister, func() {
			go func() {
				// Show previous modal and display loading status or error messages
				cfReg.Show()
				*cfReg.isSending = true
				_, err := cfReg.Dexc().RegisterWithDEXServer(cfReg.Host,
					cert,
					int64(cfReg.Fee.Amt),
					int32(cfReg.Fee.ID),
					[]byte(password))
				if err != nil {
					*cfReg.isSending = false
					cfReg.Toast.NotifyError(err.Error())
					return
				}
				cfReg.completed(cfReg.Exchange)
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
