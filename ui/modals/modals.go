package modals

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type (
	C = layout.Context
	D = layout.Dimensions

	common struct {
		theme                 *decredmaterial.Theme
		walletName            decredmaterial.Editor
		oldSpendingPassword   decredmaterial.Editor
		spendingPassword      decredmaterial.Editor
		matchSpendingPassword decredmaterial.Editor
		extendedPublicKey     decredmaterial.Editor
		passwordStrength      decredmaterial.ProgressBarStyle
		alert                 *widget.Icon
		confirm               decredmaterial.Button
		cancel                decredmaterial.Button
		cancelWidget          *widget.Clickable
		submitWidget          *widget.Clickable
	}

	modal interface {
		getTitle() string
		onCancel()
		onConfirm()
		Layout(gtx C) []layout.Widget
	}

	Modals struct {
		theme     *decredmaterial.Theme
		common    *common
		container *decredmaterial.Modal
		current   modal
		modals    map[string]modal

		isLoading     bool
		cancelWidget  *widget.Clickable
		confirmWidget *widget.Clickable
		cancel        decredmaterial.Button
		confirm       decredmaterial.Button
	}
)

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func New(theme *decredmaterial.Theme) *Modals {
	m := &Modals{
		theme:         theme,
		container:     theme.Modal(),
		cancelWidget:  new(widget.Clickable),
		confirmWidget: new(widget.Clickable),
		isLoading:     false,
	}

	m.cancel = theme.Button(m.cancelWidget, "Cancel")
	m.cancel.Font.Weight = text.Bold
	m.cancel.Background = theme.Color.Surface
	m.cancel.Color = theme.Color.Primary

	m.confirm = theme.Button(m.confirmWidget, "Confirm")

	m.initCommon()
	m.registerModals()

	return m
}

func (m *Modals) initCommon() {
	cancel := m.theme.Button(new(widget.Clickable), "Cancel")
	confirm := m.theme.Button(new(widget.Clickable), "Confirm")
	cancel.TextSize, confirm.TextSize = values.TextSize16, values.TextSize16

	spendingPassword := m.theme.EditorPassword(new(widget.Editor), "Spending password")
	spendingPassword.Editor.SingleLine, spendingPassword.Editor.Submit = true, true

	matchSpendingPassword := m.theme.EditorPassword(new(widget.Editor), "Confirm spending password")
	matchSpendingPassword.Editor.Submit, matchSpendingPassword.Editor.SingleLine = true, true

	oldSpendingPassword := m.theme.EditorPassword(new(widget.Editor), "Old spending password")
	oldSpendingPassword.Editor.Submit, oldSpendingPassword.Editor.SingleLine = true, true

	walletName := m.theme.Editor(new(widget.Editor), "")
	walletName.Editor.SingleLine, walletName.Editor.Submit = true, true

	extendedPublicKey := m.theme.Editor(new(widget.Editor), "Extended public key")
	extendedPublicKey.Editor.Submit = true

	m.common = &common{
		theme:                 m.theme,
		confirm:               confirm,
		cancel:                cancel,
		spendingPassword:      spendingPassword,
		matchSpendingPassword: matchSpendingPassword,
		oldSpendingPassword:   oldSpendingPassword,
		walletName:            walletName,
		alert:                 mustIcon(widget.NewIcon(icons.AlertError)),
		extendedPublicKey:     extendedPublicKey,
		passwordStrength:      m.theme.ProgressBar(0),
	}
}

func (m *Modals) registerModals() {
	m.modals = make(map[string]modal)
	m.registerImportWatchOnlyWalletModal()
	m.registerConfirmRemoveWalletModal()
	m.registerCreateWalletModal()
	m.registerRenameWalletModal()
	m.registerCreateAccountModal()
	m.registerRemoveWalletModal()
	m.registerChangePasswordModal()
	m.registerSetStartupPasswordModal()
	m.registerVerifyMessageInfoModal()
	m.registerSignMessageInfoModal()
	m.registerRescanWalletModal()
	m.registerSecurityToolsInfoModal()
	m.registerSecurityToolsInfoModal() 
	m.registerReceiveInfoModal()
	m.registerPrivacyInfoModal()
	m.registerSetupMixerInfoModal()
	m.registerWarnExistsMixerAccountModal()
	m.registerUnlockWalletRestoreModal()
	m.registerTransactionDetailsInfoModal()
}

func (m *Modals) Open(name string, onSubmit, onCancel func()) {
	if current, ok := m.modals[name]; ok {
		m.current = current
	}
}

func (m *Modals) layoutActionButtons(gtx C) layout.Widget {
	return func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return m.cancel.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return m.confirm.Layout(gtx)
				}),
			)
		})
	}
}

func (m *Modals) Layout(gtx C) D {
	m.handle(gtx)
	if m.current != nil {
		w := []layout.Widget{
			func(gtx C) D {
				return m.theme.H6(m.current.getTitle()).Layout(gtx)
			},
		}
		w = append(w, m.current.Layout(gtx)...)
		w = append(w, m.layoutActionButtons(gtx))
		return m.container.Layout(gtx, w, 900)
	}

	return D{}
}

func (m *Modals) onCancel() {
	m.current.onCancel()
	m.current = nil
}

func (m *Modals) onConfirm() {
	m.current.onConfirm()
}

func (m *Modals) handle(gtx C) {
	for m.cancelWidget.Clicked() {
		m.onCancel()

	}

	for m.confirmWidget.Clicked() {
		m.onConfirm()
	}
}
