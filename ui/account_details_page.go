package ui

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageAccountDetails = "AccountDetails"

type acctDetailsPage struct {
	common                   *pageCommon
	wallet                   *wallet.Wallet
	current                  wallet.InfoShort
	theme                    *decredmaterial.Theme
	acctDetailsPageContainer layout.List
	backButton               decredmaterial.IconButton
	acctInfo                 **wallet.Account
	editAccount              *widget.Clickable
	errorReceiver            chan error
}

func AcctDetailsPage(common *pageCommon) Page {
	pg := &acctDetailsPage{
		acctDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		common:        common,
		wallet:        common.wallet,
		acctInfo:      common.walletAccount,
		theme:         common.theme,
		editAccount:   new(widget.Clickable),
		errorReceiver: make(chan error),
	}

	pg.backButton, _ = common.SubPageHeaderButtons()

	return pg
}

func (pg *acctDetailsPage) OnResume() {

}

func (pg *acctDetailsPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common

	widgets := []func(gtx C) D{
		func(gtx C) D {
			return pg.accountBalanceLayout(gtx, common)
		},
		func(gtx C) D {
			m := values.MarginPadding10
			return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
				return pg.theme.Separator().Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.accountInfoLayout(gtx)
		},
	}

	title := "Not found"
	if *pg.acctInfo != nil {
		title = (*pg.acctInfo).Name
	}
	acctName := strings.Title(strings.ToLower(title))
	body := func(gtx C) D {
		page := SubPage{
			title:      acctName,
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			backButton: pg.backButton,
			back: func() {
				common.changePage(PageWallet)
			},
			body: func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
					return pg.theme.Card().Layout(gtx, func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
							if *pg.acctInfo == nil {
								return layout.Dimensions{}
							}
							return pg.acctDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
								return layout.Inset{}.Layout(gtx, widgets[i])
							})
						})
					})
				})
			},
			extraItem: pg.editAccount,
			extra: func(gtx C) D {
				return layout.Inset{}.Layout(gtx, func(gtx C) D {
					edit := common.icons.editIcon
					edit.Scale = 1
					return layout.E.Layout(gtx, edit.Layout)
				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}
	return pg.common.UniformPadding(gtx, body)
}

func (pg *acctDetailsPage) accountBalanceLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
	totalBalanceMain, totalBalanceSub := breakBalance(common.printer, (*pg.acctInfo).TotalBalance)
	spendable := dcrutil.Amount((*pg.acctInfo).SpendableBalance).String()
	spendableMain, spendableSub := breakBalance(common.printer, spendable)
	immatureRewards := dcrutil.Amount((*pg.acctInfo).Balance.ImmatureReward).String()
	rewardBalanceMain, rewardBalanceSub := breakBalance(common.printer, immatureRewards)
	lockedByTickets := dcrutil.Amount((*pg.acctInfo).Balance.LockedByTickets).String()
	lockBalanceMain, lockBalanceSub := breakBalance(common.printer, lockedByTickets)
	votingAuthority := dcrutil.Amount((*pg.acctInfo).Balance.VotingAuthority).String()
	voteBalanceMain, voteBalanceSub := breakBalance(common.printer, votingAuthority)
	immatureStakeGen := dcrutil.Amount((*pg.acctInfo).Balance.ImmatureStakeGeneration).String()
	stakeBalanceMain, stakeBalanceSub := breakBalance(common.printer, immatureStakeGen)

	return pg.pageSections(gtx, func(gtx C) D {
		accountIcon := common.icons.accountIcon
		if (*pg.acctInfo).Name == "imported" { //TODO
			accountIcon = common.icons.importedAccountIcon
		}
		accountIcon.Scale = 1

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						m := values.MarginPadding10
						return layout.Inset{
							Right: m,
							Top:   m,
						}.Layout(gtx, accountIcon.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return pg.acctBalLayout(gtx, "Total Balance", totalBalanceMain, totalBalanceSub, true)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.acctBalLayout(gtx, "Spendable", spendableMain, spendableSub, false)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.acctBalLayout(gtx, "Immature Rewards", rewardBalanceMain, rewardBalanceSub, false)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.acctBalLayout(gtx, "Locked By Tickets", lockBalanceMain, lockBalanceSub, false)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.acctBalLayout(gtx, "Voting Authority", voteBalanceMain, voteBalanceSub, false)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.acctBalLayout(gtx, "Immature Stake Gen", stakeBalanceMain, stakeBalanceSub, false)
			}),
		)
	})
}

func (pg *acctDetailsPage) acctBalLayout(gtx layout.Context, balType string, mainBalance, subBalance string, isFirst bool) layout.Dimensions {
	mainLabel := pg.theme.Body1(mainBalance)
	subLabel := pg.theme.Caption(subBalance)
	subLabel.Color = pg.theme.Color.DeepBlue
	marginTop := values.MarginPadding16
	marginLeft := values.MarginPadding35
	if isFirst {
		mainLabel = pg.theme.H4(mainBalance)
		subLabel = pg.theme.Body1(subBalance)
		marginTop = values.MarginPadding0
		marginLeft = values.MarginPadding0
	}
	return layout.Inset{
		Right: values.MarginPadding10,
		Top:   marginTop,
		Left:  marginLeft,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
					layout.Rigid(mainLabel.Layout),
					layout.Rigid(subLabel.Layout),
				)
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body2(balType)
				txt.Color = pg.theme.Color.Gray
				return txt.Layout(gtx)
			}),
		)
	})
}

func (pg *acctDetailsPage) accountInfoLayout(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		m := values.MarginPadding10
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.acctInfoLayout(gtx, "Account Number", fmt.Sprint((*pg.acctInfo).Number))
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Top:    m,
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.acctInfoLayout(gtx, "HD Path", (*pg.acctInfo).HDPath)
				})
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					ext := (*pg.acctInfo).Keys.External
					int := (*pg.acctInfo).Keys.Internal
					imp := (*pg.acctInfo).Keys.Imported
					return pg.acctInfoLayout(gtx, "Key", ext+" external, "+int+" internal, "+imp+" imported")
				})
			}),
		)
	})
}

func (pg *acctDetailsPage) acctInfoLayout(gtx layout.Context, leftText, rightText string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					leftTextLabel := pg.theme.Label(values.TextSize14, leftText)
					leftTextLabel.Color = pg.theme.Color.Gray
					return leftTextLabel.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, pg.theme.Body1(rightText).Layout)
		}),
	)
}

func (pg *acctDetailsPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	m := values.MarginPadding20
	mtb := values.MarginPadding5
	return layout.Inset{Left: m, Right: m, Top: mtb, Bottom: mtb}.Layout(gtx, body)
}

func (pg *acctDetailsPage) handle() {
	common := pg.common

	if pg.editAccount.Clicked() {
		pg.current = common.info.Wallets[*common.selectedWallet]
		textModal := newTextInputModal(common).
			hint("Account name").
			positiveButton(values.String(values.StrRename), func(newName string, tim *textInputModal) bool {
				pg.wallet.RenameAccount(pg.current.ID, (*pg.acctInfo).Number, newName, pg.errorReceiver)
				(*pg.acctInfo).Name = newName
				return true
			})

		textModal.title("Rename account").
			negativeButton(values.String(values.StrCancel), func() {})
		textModal.Show()
	}
}

func (pg *acctDetailsPage) onClose() {}
