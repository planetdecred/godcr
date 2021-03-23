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
	wallet                   *wallet.Wallet
	current                  wallet.InfoShort
	theme                    *decredmaterial.Theme
	acctDetailsPageContainer layout.List
	backButton               decredmaterial.IconButton
	acctInfo                 **wallet.Account
	line                     *decredmaterial.Line
	editAccount              *widget.Clickable
	errorReceiver            chan error
}

func (win *Window) AcctDetailsPage(common pageCommon) layout.Widget {
	pg := &acctDetailsPage{
		acctDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		wallet:        common.wallet,
		acctInfo:      &win.walletAccount,
		theme:         common.theme,
		backButton:    common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		line:          common.theme.Line(),
		editAccount:   new(widget.Clickable),
		errorReceiver: make(chan error),
	}

	pg.line.Color = common.theme.Color.Background
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.Handler(gtx, common)
		return pg.Layout(gtx, common)
	}
}

func (pg *acctDetailsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	widgets := []func(gtx C) D{
		func(gtx C) D {
			return pg.accountBalanceLayout(gtx, &common)
		},
		func(gtx C) D {
			pg.line.Width = gtx.Constraints.Max.X
			pg.line.Height = 2
			m := values.MarginPadding5
			return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
				return pg.line.Layout(gtx)
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
			back: func() {
				common.changePage(PageWallet)
			},
			body: func(gtx C) D {
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
			},
			extraItem: pg.editAccount,
			extra: func(gtx C) D {
				return layout.Inset{}.Layout(gtx, func(gtx C) D {
					edit := common.icons.editIcon
					edit.Scale = 1
					return layout.E.Layout(gtx, func(gtx C) D {
						return edit.Layout(gtx)
					})
				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *acctDetailsPage) accountBalanceLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
	totalBalanceMain, totalBalanceSub := breakBalance((*pg.acctInfo).TotalBalance)
	spendable := dcrutil.Amount((*pg.acctInfo).SpendableBalance).String()
	spendableMain, spendableSub := breakBalance(spendable)
	immatureRewards := dcrutil.Amount((*pg.acctInfo).Balance.ImmatureReward).String()
	rewardBalanceMain, rewardBalanceSub := breakBalance(immatureRewards)
	lockedByTickets := dcrutil.Amount((*pg.acctInfo).Balance.LockedByTickets).String()
	lockBalanceMain, lockBalanceSub := breakBalance(lockedByTickets)
	votingAuthority := dcrutil.Amount((*pg.acctInfo).Balance.VotingAuthority).String()
	voteBalanceMain, voteBalanceSub := breakBalance(votingAuthority)
	immatureStakeGen := dcrutil.Amount((*pg.acctInfo).Balance.ImmatureStakeGeneration).String()
	stakeBalanceMain, stakeBalanceSub := breakBalance(immatureStakeGen)

	return pg.pageSections(gtx, func(gtx C) D {
		accountIcon := common.icons.accountIcon
		if (*pg.acctInfo).Name == "imported" {
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
						}.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
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
	marginTop := values.MarginPadding15
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
					layout.Rigid(func(gtx C) D {
						return mainLabel.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return subLabel.Layout(gtx)
					}),
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
	acctInfoLayout := func(gtx layout.Context, leftText, rightText string) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						leftTextLabel := pg.theme.Body1(leftText)
						leftTextLabel.Color = pg.theme.Color.Gray
						return leftTextLabel.Layout(gtx)
					}),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					inset := layout.Inset{
						Right: values.MarginPadding10,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return pg.theme.Body1(rightText).Layout(gtx)
					})
				})
			}),
		)
	}

	return pg.pageSections(gtx, func(gtx C) D {
		m := values.MarginPadding10
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return acctInfoLayout(gtx, "Account Number", fmt.Sprint((*pg.acctInfo).Number))
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Top:    m,
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return acctInfoLayout(gtx, "HD Path", (*pg.acctInfo).HDPath)
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
					return acctInfoLayout(gtx, "Key", ext+" external, "+int+" internal, "+imp+" imported")
				})
			}),
		)
	})
}

func (pg *acctDetailsPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	m := values.MarginPadding20
	mtb := values.MarginPadding5
	return layout.Inset{Left: m, Right: m, Top: mtb, Bottom: mtb}.Layout(gtx, body)
}

func (pg *acctDetailsPage) Handler(gtx layout.Context, common pageCommon) {
	if pg.backButton.Button.Clicked() {
		common.changePage(PageWallet)
	}

	if pg.editAccount.Clicked() {
		pg.current = common.info.Wallets[*common.selectedWallet]
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RenameAccountTemplate,
				title:    "Rename account",
				confirm: func(name string) {
					pg.wallet.RenameAccount(pg.current.ID, (*pg.acctInfo).Number, name, pg.errorReceiver)
					(*pg.acctInfo).Name = name
				},
				confirmText: "Rename",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}
}
