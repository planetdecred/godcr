package governance

import (
	"context"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const TreasuryPageID = "Treasury"

type TreasuryPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	multiWallet   *dcrlibwallet.MultiWallet
	wallets       []*dcrlibwallet.Wallet
	treasuryItems []*components.TreasuryItem

	listContainer      *widget.List
	viewGovernanceKeys *decredmaterial.Clickable
	copyRedirectURL    *decredmaterial.Clickable
	redirectIcon       *decredmaterial.Image

	searchEditor decredmaterial.Editor
	infoButton   decredmaterial.IconButton

	isPolicyFetchInProgress bool
}

func NewTreasuryPage(l *load.Load) *TreasuryPage {
	pg := &TreasuryPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(TreasuryPageID),
		multiWallet:      l.WL.MultiWallet,
		wallets:          l.WL.SortedWalletList(),
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		redirectIcon:       l.Theme.Icons.RedirectIcon,
		viewGovernanceKeys: l.Theme.NewClickable(true),
		copyRedirectURL:    l.Theme.NewClickable(false),
	}

	pg.searchEditor = l.Theme.IconEditor(new(widget.Editor), values.String(values.StrSearch), l.Theme.Icons.SearchIcon, true)
	pg.searchEditor.Editor.SingleLine, pg.searchEditor.Editor.Submit, pg.searchEditor.Bordered = true, true, false

	_, pg.infoButton = components.SubpageHeaderButtons(l)
	pg.infoButton.Size = values.MarginPadding20

	return pg
}

func (pg *TreasuryPage) ID() string {
	return TreasuryPageID
}

func (pg *TreasuryPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.FetchPolicies()
}

func (pg *TreasuryPage) OnNavigatedFrom() {
	if pg.ctxCancel != nil {
		pg.ctxCancel()
	}
}

func (pg *TreasuryPage) HandleUserInteractions() {
	for i := range pg.treasuryItems {
		if pg.treasuryItems[i].SetChoiceButton.Clicked() {
			pg.updatePolicyPreference(pg.treasuryItems[i])
		}
	}

	if pg.infoButton.Button.Clicked() {
		infoModal := modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrTreasurySpending)).
			Body(values.String(values.StrTreasurySpendingInfo)).
			SetCancelable(true).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(infoModal)
	}

	for pg.viewGovernanceKeys.Clicked() {
		host := "https://github.com/decred/dcrd/blob/master/chaincfg/mainnetparams.go#L477"
		if pg.WL.MultiWallet.NetType() == dcrlibwallet.Testnet3 {
			host = "https://github.com/decred/dcrd/blob/master/chaincfg/testnetparams.go#L390"
		}

		info := modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrVerifyGovernanceKeys)).
			Body(values.String(values.StrCopyLink)).
			SetCancelable(true).
			UseCustomWidget(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						border := widget.Border{Color: pg.Theme.Color.Gray4, CornerRadius: values.MarginPadding10, Width: values.MarginPadding2}
						wrapper := pg.Theme.Card()
						wrapper.Color = pg.Theme.Color.Gray4
						return border.Layout(gtx, func(gtx C) D {
							return wrapper.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Flexed(0.9, pg.Theme.Body1(host).Layout),
										layout.Flexed(0.1, func(gtx C) D {
											return layout.E.Layout(gtx, func(gtx C) D {
												if pg.copyRedirectURL.Clicked() {
													clipboard.WriteOp{Text: host}.Add(gtx.Ops)
													pg.Toast.Notify(values.String(values.StrCopied))
												}
												return pg.copyRedirectURL.Layout(gtx, pg.Theme.Icons.CopyIcon.Layout24dp)
											})
										}),
									)
								})
							})
						})
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:  values.MarginPaddingMinus10,
							Left: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							label := pg.Theme.Body2(values.String(values.StrWebURL))
							label.Color = pg.Theme.Color.GrayText2
							return label.Layout(gtx)
						})
					}),
				)
			}).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}

	if pg.isPolicyFetchInProgress {
		time.AfterFunc(time.Second*1, func() {
			pg.ParentWindow().Reload()
		})
	}

	pg.searchEditor.EditorIconButtonEvent = func() {
		//TODO: treasury search functionality
	}
}

func (pg *TreasuryPage) FetchPolicies() {
	selectedWallet := pg.WL.SelectedWallet.Wallet

	pg.isPolicyFetchInProgress = true

	// Fetch (or re-fetch) treasury policies in background as this makes
	// a network call. Refresh the window once the call completes.
	go func() {
		pg.treasuryItems = components.LoadPolicies(pg.Load, selectedWallet, pg.WL.MultiWallet.PiKey(0))
		pg.isPolicyFetchInProgress = true
		pg.ParentWindow().Reload()
	}()

	// Refresh the window now to signify that the syncing
	// has started with pg.isSyncing set to true above.
	pg.ParentWindow().Reload()
}

func (pg *TreasuryPage) Layout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.Theme.Label(values.TextSize20, values.String(values.StrTreasurySpending)).Layout),
						layout.Rigid(pg.infoButton.Layout),
					)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, pg.layoutVerifyGovernanceKeys)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{
							Top: values.MarginPadding10,
						}.Layout(gtx, pg.layoutContent)
					}),
				)
			})
		}),
	)
}

func (pg *TreasuryPage) layoutVerifyGovernanceKeys(gtx C) D {
	return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
		return pg.viewGovernanceKeys.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Right: values.MarginPadding10,
					}.Layout(gtx, pg.redirectIcon.Layout16dp)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top: values.MarginPaddingMinus2,
					}.Layout(gtx, pg.Theme.Label(values.TextSize16, values.String(values.StrVerifyGovernanceKeys)).Layout)
				}),
			)
		})
	})
}

func (pg *TreasuryPage) layoutContent(gtx C) D {
	if len(pg.treasuryItems) == 0 {
		return components.LayoutNoPoliciesFound(gtx, pg.Load, pg.isPolicyFetchInProgress)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
				return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
					return list.Layout(gtx, len(pg.treasuryItems), func(gtx C, i int) D {
						return decredmaterial.LinearLayout{
							Orientation: layout.Vertical,
							Width:       decredmaterial.MatchParent,
							Height:      decredmaterial.WrapContent,
							Background:  pg.Theme.Color.Surface,
							Direction:   layout.W,
							Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
							Padding:     layout.UniformInset(values.MarginPadding15),
							Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.
							Layout2(gtx, func(gtx C) D {
								return components.TreasuryItemWidget(gtx, pg.Load, pg.treasuryItems[i])
							})
					})
				})
			})
		}),
	)
}

func (pg *TreasuryPage) updatePolicyPreference(treasuryItem *components.TreasuryItem) {
	passwordModal := modal.NewPasswordModal(pg.Load).
		Title(values.String(values.StrConfirmVote)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				selectedWallet := pg.WL.SelectedWallet.Wallet
				votingPreference := treasuryItem.OptionsRadioGroup.Value
				err := selectedWallet.SetTreasuryPolicy(treasuryItem.Policy.PiKey, votingPreference, "", []byte(password))
				if err != nil {
					if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
						pm.SetError(values.String(values.StrInvalidPassphrase))
					} else {
						pm.Toast.NotifyError(err.Error())
					}
					pm.SetLoading(false)
					return
				}
				go pg.FetchPolicies() // re-fetch policies when voting is done.
				pm.Toast.Notify(values.String(values.StrPolicySetSuccessful))
				pm.Dismiss()
			}()

			return false
		})
	pg.ParentWindow().ShowModal(passwordModal)
}
