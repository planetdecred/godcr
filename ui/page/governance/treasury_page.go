package governance

import (
	"context"
	"time"

	// "gioui.org/font/gofont"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	// "gioui.org/unit"
	"gioui.org/widget"

	// "gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const TreasuryPageID = "Treasury"

type TreasuryPage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	multiWallet   *dcrlibwallet.MultiWallet
	wallets       []*dcrlibwallet.Wallet
	LiveTickets   []*dcrlibwallet.Transaction
	treasuryItems []*components.TreasuryItem

	listContainer       *widget.List
	syncButton          *widget.Clickable
	viewVotingDashboard *decredmaterial.Clickable
	copyRedirectURL     *decredmaterial.Clickable
	redirectIcon        *decredmaterial.Image

	walletDropDown *decredmaterial.DropDown
	treasuryList   *decredmaterial.ClickableList

	searchEditor      decredmaterial.Editor
	infoButton        decredmaterial.IconButton
	optionsRadioGroup *widget.Enum

	setChoiceButton decredmaterial.Button

	syncCompleted bool
	isSyncing     bool
}

func NewTreasuryPage(l *load.Load) *TreasuryPage {
	pg := &TreasuryPage{
		Load:         l,
		multiWallet:  l.WL.MultiWallet,
		wallets:      l.WL.SortedWalletList(),
		treasuryList: l.Theme.NewClickableList(layout.Vertical),
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		syncButton:          new(widget.Clickable),
		optionsRadioGroup:   new(widget.Enum),
		redirectIcon:        l.Theme.Icons.RedirectIcon,
		viewVotingDashboard: l.Theme.NewClickable(true),
		copyRedirectURL:     l.Theme.NewClickable(false),
	}

	pg.setChoiceButton = l.Theme.Button("Set Choice")

	pg.searchEditor = l.Theme.IconEditor(new(widget.Editor), values.String(values.StrSearch), l.Theme.Icons.SearchIcon, true)
	pg.searchEditor.Editor.SingleLine, pg.searchEditor.Editor.Submit, pg.searchEditor.Bordered = true, true, false

	_, pg.infoButton = components.SubpageHeaderButtons(l)
	pg.infoButton.Size = values.MarginPadding20

	pg.walletDropDown = components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets, values.TxDropdownGroup, 0)

	// update agenda options prefrence to that of the selected wallet
	// selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
	// treasuryItems := components.LoadPolicies(pg.Load, selectedWallet, dcrlibwallet.PiKey)
	// for _, treasuryItem := range treasuryItems {
	// 	// voteChoices := make([]string, len(consensusItem.Agenda.Choices))
	// 	voteChoices := [...]string{"Yes", "No", "Abstain"}
	// 	initialValue := treasuryItem.Policy.Policy
	// 	treasuryItem.VoteChoices = voteChoices

	// 	treasuryItem.VoteChoices = voteChoices
	// 	// avm.initialValue = consensusItem.Agenda.VotingPreference
	// 	// avm.optionsRadioGroup.Value = avm.initialValue

	// }

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
	for pg.walletDropDown.Changed() {
		pg.FetchPolicies()
	}

	// for i := range pg.treasuryItems {
	// 	if pg.treasuryItems[i].VoteButton.Clicked() {
	// 		newAgendaVoteModal(pg.Load, &pg.treasuryItems[i].Agenda, func() {
	// 			go pg.FetchAgendas() // re-fetch agendas when modal is dismissed
	// 		}).Show()
	// 	}
	// }

	for pg.syncButton.Clicked() {
		go pg.FetchPolicies()
	}

	if pg.infoButton.Button.Clicked() {
		modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrTreasurySpending)).
			Body(values.String(values.StrTreasurySpendingInfo)).
			SetCancelable(true).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) {}).Show()
	}

	for pg.viewVotingDashboard.Clicked() {
		host := "https://github.com/decred/dcrd/blob/master/chaincfg/mainnetparams.go#L485"
		// if pg.WL.MultiWallet.NetType() == dcrlibwallet.Testnet3 {
		// 	host = "https://voting.decred.org/testnet"
		// }

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
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) {})
		pg.ShowModal(info)
	}

	if pg.syncCompleted {
		time.AfterFunc(time.Second*1, func() {
			pg.syncCompleted = false
			pg.RefreshWindow()
		})
	}

	pg.searchEditor.EditorIconButtonEvent = func() {
		//TODO: treasury search functionality
	}
}

func (pg *TreasuryPage) FetchPolicies() {
	// newestFirst := pg.orderDropDown.SelectedIndex() == 0
	selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]

	pg.isSyncing = true

	// Fetch (or re-fetch) treasury policies in background as this makes
	// a network call. Refresh the window once the call completes.
	go func() {
		pg.treasuryItems = components.LoadPolicies(pg.Load, selectedWallet, dcrlibwallet.PiKey)
		pg.isSyncing = false
		pg.syncCompleted = true
		pg.RefreshWindow()
	}()

	// Refresh the window now to signify that the syncing
	// has started with pg.isSyncing set to true above.
	pg.RefreshWindow()
}

func (pg *TreasuryPage) Layout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.Theme.Label(values.TextSize20, values.String(values.StrTreasurySpending)).Layout), // Do we really need to display the title? nav is proposals already
						layout.Rigid(pg.infoButton.Layout),
					)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, pg.layoutRedirectVoting)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{
							Top: values.MarginPadding60,
						}.Layout(gtx, pg.layoutContent)
					}),
					layout.Expanded(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							card := pg.Theme.Card()
							card.Radius = decredmaterial.Radius(8)
							return card.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding8).Layout(gtx, func(gtx C) D {
									return D{}
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return pg.walletDropDown.Layout(gtx, 45, true)
					}),
				)
			})
		}),
	)
}

func (pg *TreasuryPage) lineSeparator(inset layout.Inset) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, pg.Theme.Separator().Layout)
	}
}

func (pg *TreasuryPage) layoutRedirectVoting(gtx C) D {
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.viewVotingDashboard.Layout(gtx, func(gtx C) D {
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
		}),
	)
}

func (pg *TreasuryPage) layoutContent(gtx C) D {
	if len(pg.treasuryItems) == 0 {
		return components.LayoutNoPoliciesFound(gtx, pg.Load, pg.isSyncing)
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
