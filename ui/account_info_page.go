package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageAccountDetails = "AccountDetails"

type acctDetailsPage struct {
	theme                    *decredmaterial.Theme
	acctDetailsPageContainer layout.List
	backButton               decredmaterial.IconButton
	acctInfo                 **wallet.Account
	line                     *decredmaterial.Line
	editAccount              *widget.Clickable
}

func (win *Window) AcctDetailsPage(common pageCommon) layout.Widget {
	pg := &acctDetailsPage{
		acctDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		acctInfo:    &win.walletAccount,
		theme:       common.theme,
		backButton:  common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		line:        common.theme.Line(),
		editAccount: new(widget.Clickable),
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
			return pg.header(gtx, &common)
		},
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

	body := common.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			if *pg.acctInfo == nil {
				return layout.Dimensions{}
			}
			return pg.acctDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.Inset{}.Layout(gtx, widgets[i])
			})
		})
	})

	return body
}

func (pg *acctDetailsPage) header(gtx layout.Context, common *pageCommon) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
							return pg.backButton.Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.H6("")
					if *pg.acctInfo != nil {
						txt.Text = (*pg.acctInfo).Name
					} else {
						txt.Text = "Not found"
					}

					txt.Alignment = text.Middle
					return txt.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					edit := common.icons.editIcon
					edit.Scale = 0.25
					return layout.E.Layout(gtx, func(gtx C) D {
						return decredmaterial.Clickable(gtx, pg.editAccount, func(gtx C) D {
							return edit.Layout(gtx)
						})
					})
				}),
			)
		})
	})
}

func (pg *acctDetailsPage) accountBalanceLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
	acctBalLayout := func(balType string, mainBalance, subBalance decredmaterial.Label) layout.Dimensions {
		return layout.Inset{
			Right: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return mainBalance.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return subBalance.Layout(gtx)
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

	tMain, tSub := breakBalance((*pg.acctInfo).TotalBalance)
	spendable := dcrutil.Amount((*pg.acctInfo).SpendableBalance).String()
	sMain, sSub := breakBalance(spendable)

	return pg.pageSections(gtx, func(gtx C) D {
		accountIcon := common.icons.accountIcon
		if (*pg.acctInfo).Name == "imported" {
			accountIcon = common.icons.importedAccountIcon
		}
		accountIcon.Scale = 0.8

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						m := values.MarginPadding10
						inset := layout.Inset{
							Right: m,
							Top:   m,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								mainLabel := pg.theme.H4(tMain)
								subLabel := pg.theme.Body1(tSub)
								return acctBalLayout("Total Balance", mainLabel, subLabel)
							}),
							layout.Rigid(func(gtx C) D {
								mainLabel := pg.theme.Body1(sMain)
								subLabel := pg.theme.Caption(sSub)
								inset := layout.Inset{
									Top: values.MarginPadding5,
								}
								return inset.Layout(gtx, func(gtx C) D {
									return acctBalLayout("Spendable", mainLabel, subLabel)
								})
							}),
						)
					}),
				)
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
		*common.page = PageWallet
	}

	if pg.editAccount.Clicked() {
		fmt.Println("edit icon was clicked")
	}
}
