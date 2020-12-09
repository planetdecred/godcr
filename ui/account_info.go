package ui

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"

	// "github.com/decred/dcrd/dcrutil"
	// "github.com/planetdecred/dcrlibwallet"
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
	infoBtn                  decredmaterial.IconButton
	accountIcon              *widget.Image
	line                     *decredmaterial.Line
}

func (win *Window) AcctDetailsPage(common pageCommon) layout.Widget {
	pg := &acctDetailsPage{
		acctDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		acctInfo:   &win.walletAccount,
		theme:      common.theme,
		backButton: common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		line:       common.theme.Line(),
	}

	pg.line.Color = common.theme.Color.Background
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)
	pg.infoBtn = common.theme.IconButton(new(widget.Clickable), common.icons.actionInfo)
	pg.infoBtn.Color = common.theme.Color.Gray
	pg.infoBtn.Background = common.theme.Color.Surface
	pg.infoBtn.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.Handler(gtx, common)
		return pg.Layout(gtx, common)
	}
}

func (pg *acctDetailsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	widgets := []func(gtx C) D{
		func(gtx C) D {
			return pg.header(gtx)
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
			return pg.accountInfoLayout(gtx, &common)
		},
	}

	body := common.Layout(gtx, func(gtx C) D {
		return decredmaterial.Card{Color: common.theme.Color.Surface, CornerStyle: decredmaterial.RoundedEdge}.Layout(gtx, func(gtx C) D {
			// if *pg.acctInfo == nil {
			// 	return layout.Dimensions{}
			// }
			return pg.acctDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.Inset{}.Layout(gtx, widgets[i])
			})
		})
	})

	return body
}

func (pg *acctDetailsPage) header(gtx layout.Context) layout.Dimensions {
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
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.infoBtn.Layout(gtx)
					})
				}),
			)
		})
	})
}

func (pg *acctDetailsPage) accountBalanceLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
	// txnWidgets := transactionWdg{}
	// initTxnWidgets(common, *pg.txnInfo, &txnWidgets)
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

	main, sub := breakBalance("7.89087676")

	return pg.pageSections(gtx, func(gtx C) D {
		pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.accountIcon)}
		// if name == "imported" {
		// 	pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.importedAccountIcon)}
		// }
		pg.accountIcon.Scale = 0.8

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding5,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return pg.accountIcon.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								mainLabel := pg.theme.H4(main)
								subLabel := pg.theme.Body1(sub)
								return acctBalLayout("Total Balance", mainLabel, subLabel)
							}),
							layout.Rigid(func(gtx C) D {
								mainLabel := pg.theme.Body1(main)
								subLabel := pg.theme.Caption(sub)
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

func (pg *acctDetailsPage) accountInfoLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
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
		pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.accountIcon)}
		// if name == "imported" {
		// 	pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.importedAccountIcon)}
		// }
		pg.accountIcon.Scale = 0.8
		m := values.MarginPadding10

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return acctInfoLayout(gtx, "Account Number", "0")
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Top:    m,
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return acctInfoLayout(gtx, "HD Path", "m/ 44' / 1' / 0' ")
				})
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return acctInfoLayout(gtx, "Key", "21 external, 20 internal, 0 imported")
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
}
