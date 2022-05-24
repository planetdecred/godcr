package wallets

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const HexRestorePageID = "hex_restore"

type HexRestore struct {
	*load.Load
	hexEditor       decredmaterial.Editor
	validateHex     decredmaterial.Button
	keyEvent        chan *key.Event
	restoreComplete func()
}

func NewHexRestorePage(l *load.Load, onRestoreComplete func()) *HexRestore {
	pg := &HexRestore{
		Load:            l,
		restoreComplete: onRestoreComplete,
		keyEvent:        make(chan *key.Event),
	}

	pg.hexEditor = l.Theme.Editor(new(widget.Editor), "Enter hex")
	pg.hexEditor.Editor.SingleLine = true
	pg.hexEditor.Editor.SetText("")

	pg.validateHex = l.Theme.Button("Validate hex")
	pg.validateHex.Font.Weight = text.Medium

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *HexRestore) ID() string {
	return HexRestorePageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *HexRestore) OnNavigatedTo() {
	pg.Load.SubscribeKeyEvent(pg.keyEvent, pg.ID())
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *HexRestore) HandleUserInteractions() {
	if pg.hexEditor.Editor.Focused() {
		hexText := pg.hexEditor.Editor.Text()
		// 16 and 64 is the minimum/maximum number of bytes allowed for a seed.
		if len(hexText) >= 16 && len(hexText) <= 64 {
			pg.validateHex.SetEnabled(true)
		}
	}

	for pg.validateHex.Clicked() {
		if !pg.verifyHex() {
			return
		}

		pg.Load.UnsubscribeKeyEvent(pg.ID())

		modal.NewCreatePasswordModal(pg.Load).
			Title("Enter wallet details").
			EnableName(true).
			ShowWalletInfoTip(true).
			SetParent(pg).
			PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					_, err := pg.WL.MultiWallet.RestoreWallet(walletName, pg.hexEditor.Editor.Text(), password, dcrlibwallet.PassphraseTypePass)
					if err != nil {
						m.SetError(components.TranslateErr(err))
						m.SetLoading(false)
						return
					}

					pg.Toast.Notify("Wallet restored")
					pg.hexEditor.Editor.SetText("")
					m.Dismiss()
					// Close this page and return to the previous page (most likely wallets page)
					// if there's no restoreComplete callback function.
					if pg.restoreComplete == nil {
						pg.PopWindowPage()
					} else {
						pg.restoreComplete()
					}
				}()
				return false
			}).Show()
	}

}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *HexRestore) OnNavigatedFrom() {
	pg.Load.UnsubscribeKeyEvent(pg.ID())
}

func (pg *HexRestore) Layout(gtx layout.Context) layout.Dimensions {
	dims := layout.Inset{
		Top: values.MarginPadding100,
	}.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return decredmaterial.LinearLayout{
					Orientation: layout.Vertical,
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.WrapContent,
					Background:  pg.Theme.Color.Surface,
					Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
					Padding:     layout.UniformInset(values.MarginPadding15)}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Bottom: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return pg.hexEditor.Layout(gtx)
						})
					}),
				)
			}),
			layout.Stacked(func(gtx C) D {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return layout.S.Layout(gtx, func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, pg.validateButtonSection)
				})
			}),
		)
	})

	pg.validateHex.SetEnabled(false)
	return dims
}

func (pg *HexRestore) validateButtonSection(gtx layout.Context) layout.Dimensions {
	card := pg.Theme.Card()
	card.Radius = decredmaterial.Radius(0)
	return card.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return pg.validateHex.Layout(gtx)
	})
}

func (pg *HexRestore) verifyHex() bool {

	hex := pg.hexEditor.Editor.Text()
	if !dcrlibwallet.VerifySeed(hex) {
		pg.Toast.NotifyError("invalid hex")
		return false
	}

	// Compare with existing wallets seed. On positive match abort import
	// to prevent duplicate wallet. walletWithSameSeed >= 0 if there is a match.
	walletWithSameSeed, err := pg.WL.MultiWallet.WalletWithSeed(hex)
	if err != nil {
		log.Error(err)
		return false
	}

	if walletWithSameSeed != -1 {
		pg.Toast.NotifyError("A wallet with an identical seed already exists.")
		return false
	}

	return true
}
