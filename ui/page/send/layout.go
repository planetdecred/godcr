package send

import (
	"fmt"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func (pg *Page) initLayoutWidgets() {
	pg.pageContainer = &widget.List{
		List: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
	}

	pg.txFeeCollapsible = pg.Theme.Collapsible()

	pg.nextButton = pg.Theme.Button(values.String(values.StrNext))
	pg.nextButton.TextSize = values.TextSize18
	pg.nextButton.Inset = layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding15}
	pg.nextButton.SetEnabled(false)

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(pg.Load)
	pg.backButton.Icon = pg.Theme.Icons.ContentClear

	pg.moreOption = pg.Theme.IconButton(pg.Theme.Icons.NavMoreIcon)
	pg.moreOption.Inset = layout.UniformInset(values.MarginPadding0)

	pg.retryExchange = pg.Theme.Button(values.String(values.StrRetry))
	pg.retryExchange.Background = pg.Theme.Color.Gray1
	pg.retryExchange.Color = pg.Theme.Color.Surface
	pg.retryExchange.TextSize = values.TextSize12
	pg.retryExchange.Inset = layout.Inset{
		Top:    values.MarginPadding5,
		Right:  values.MarginPadding8,
		Bottom: values.MarginPadding5,
		Left:   values.MarginPadding8,
	}

	pg.moreItems = pg.getMoreItem()
}

func (pg *Page) topNav(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding20
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.backButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: m}.Layout(gtx, pg.Theme.H6(values.String(values.StrSend)+" DCR").Layout)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(pg.infoButton.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: m}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									if pg.moreOptionIsOpen {
										pg.layoutOptionsMenu(gtx)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(pg.moreOption.Layout),
							)
						})
					}),
				)
			})
		}),
	)
}

func (pg *Page) getMoreItem() []moreItem {
	return []moreItem{
		// TODO: temp removal till issue #658 is resolved and V1.0 is release
		// {
		// 	text:   "Advanced mode",
		// 	button: pg.Theme.NewClickable(true),
		// 	id:     UTXOPageID,
		// 	action: func() {
		// 		pg.ChangeFragment(NewUTXOPage(pg.Load, pg.sourceAccountSelector.SelectedAccount()))
		// 	},
		// },
		{
			text:   values.String(values.StrClearAll),
			button: pg.Theme.NewClickable(true),
			action: func() {
				pg.resetFields()
				pg.moreOptionIsOpen = false
			},
		},
	}
}

func (pg *Page) layoutOptionsMenu(gtx layout.Context) {
	inset := layout.Inset{
		Top:  values.MarginPadding30,
		Left: values.MarginPaddingMinus100,
	}

	m := op.Record(gtx.Ops)
	inset.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding130)
		return pg.shadowBox.Layout(gtx, func(gtx C) D {
			optionsMenuCard := decredmaterial.Card{Color: pg.Theme.Color.Surface}
			optionsMenuCard.Radius = decredmaterial.Radius(5)
			return optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.moreItems), func(gtx C, i int) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.moreItems[i].button.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return pg.Theme.Body1(pg.moreItems[i].text).Layout(gtx)
								})
							})
						}),
					)
				})
			})
		})
	})
	op.Defer(gtx.Ops, m.Stop())
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *Page) Layout(gtx layout.Context) layout.Dimensions {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *Page) layoutDesktop(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, values.String(values.StrFrom), false, func(gtx C) D {
				return pg.sourceAccountSelector.Layout(pg.ParentWindow(), gtx)
			})
		},
		func(gtx C) D {
			return pg.toSection(gtx)
		},
		func(gtx C) D {
			return pg.feeSection(gtx)
		},
	}
	dims := layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Stack{Alignment: layout.NE}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return components.UniformPadding(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
									return pg.topNav(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.Theme.List(pg.pageContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
									return layout.Inset{Bottom: values.MarginPadding16, Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
										return layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}.Layout(gtx, pageContent[i])
									})
								})
							}),
						)
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D {
					return pg.balanceSection(gtx)
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			if pg.moreOptionIsOpen {
				return pg.backdrop.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					semantic.Button.Add(gtx.Ops)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				})
			}
			return D{}
		}),
	)

	return dims
}

func (pg *Page) layoutMobile(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, values.String(values.StrFrom), false, func(gtx C) D {
				return pg.sourceAccountSelector.Layout(pg.ParentWindow(), gtx)
			})
		},
		func(gtx C) D {
			return pg.toSection(gtx)
		},
		func(gtx C) D {
			return pg.feeSection(gtx)
		},
	}

	dims := layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Stack{Alignment: layout.NE}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return components.UniformMobile(gtx, false, true, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding16, Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
									return pg.topNav(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.Theme.List(pg.pageContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
									return layout.Inset{Bottom: values.MarginPadding16, Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
										return layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}.Layout(gtx, pageContent[i])
									})
								})
							}),
						)
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D {
					return pg.balanceSection(gtx)
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			if pg.moreOptionIsOpen {
				return pg.backdrop.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					semantic.Button.Add(gtx.Ops)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				})
			}
			return D{}
		}),
	)

	return dims
}

func (pg *Page) pageSections(gtx layout.Context, title string, showAccountSwitch bool, body layout.Widget) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Bottom: values.MarginPadding16,
							}
							return inset.Layout(gtx, pg.Theme.Body1(title).Layout)
						}),
						layout.Flexed(1, func(gtx C) D {
							if showAccountSwitch {
								return layout.E.Layout(gtx, func(gtx C) D {
									inset := layout.Inset{
										Top: values.MarginPaddingMinus5,
									}
									return inset.Layout(gtx, pg.sendDestination.accountSwitch.Layout)
								})
							}
							return layout.Dimensions{}
						}),
					)
				}),
				layout.Rigid(body),
			)
		})
	})
}

func (pg *Page) toSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, values.String(values.StrTo), true, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					if !pg.sendDestination.sendToAddress {
						return pg.sendDestination.destinationAccountSelector.Layout(pg.ParentWindow(), gtx)
					}
					return pg.sendDestination.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.exchangeRate != -1 && pg.usdExchangeSet {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							return pg.amount.dcrAmountEditor.Layout(gtx)
						}),
						layout.Flexed(0.1, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								icon := pg.Theme.Icons.CurrencySwapIcon
								return icon.Layout12dp(gtx)
							})
						}),
						layout.Flexed(0.45, func(gtx C) D {
							return pg.amount.usdAmountEditor.Layout(gtx)
						}),
					)
				}
				return pg.amount.dcrAmountEditor.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				if pg.exchangeRateMessage == "" {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							gtx.Constraints.Min.Y = gtx.Dp(values.MarginPadding1)
							return decredmaterial.Fill(gtx, pg.Theme.Color.Gray1)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								label := pg.Theme.Body2(pg.exchangeRateMessage)
								label.Color = pg.Theme.Color.Danger
								if pg.isFetchingExchangeRate {
									label.Color = pg.Theme.Color.Primary
								}
								return label.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								if pg.isFetchingExchangeRate {
									return layout.Dimensions{}
								}
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return layout.E.Layout(gtx, pg.retryExchange.Layout)
							}),
						)
					}),
				)
			}),
		)
	})
}

func (pg *Page) feeSection(gtx layout.Context) layout.Dimensions {
	collapsibleHeader := func(gtx C) D {
		feeText := pg.txFee
		if pg.exchangeRate != -1 && pg.usdExchangeSet {
			feeText = fmt.Sprintf("%s (%s)", pg.txFee, pg.txFeeUSD)
		}
		return pg.Theme.Body1(feeText).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.Gray4
		inset := layout.Inset{
			Top: values.MarginPadding10,
		}
		return inset.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							//TODO
							return pg.contentRow(gtx, values.String(values.StrEstimatedTime), "10 minutes (2 blocks)")
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Top:    values.MarginPadding5,
								Bottom: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.contentRow(gtx, values.String(values.StrEstimatedSize), pg.estSignedSize)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.contentRow(gtx, values.String(values.StrFee)+" "+values.String(values.StrRate), "10 atoms/Byte")
						}),
					)
				})
			})
		})
	}
	inset := layout.Inset{
		Bottom: values.MarginPadding75,
	}
	return inset.Layout(gtx, func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrFee), false, func(gtx C) D {
			return pg.txFeeCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
		})
	})
}

func (pg *Page) balanceSection(gtx layout.Context) layout.Dimensions {
	c := pg.Theme.Card()
	c.Radius = decredmaterial.Radius(0)
	return c.Layout(gtx, func(gtx C) D {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(0.6, func(gtx C) D {
					inset := layout.Inset{
						Right: values.MarginPadding15,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								inset := layout.Inset{
									Bottom: values.MarginPadding10,
								}
								return inset.Layout(gtx, func(gtx C) D {
									totalCostText := pg.totalCost
									if pg.exchangeRate != -1 && pg.usdExchangeSet {
										totalCostText = fmt.Sprintf("%s (%s)", pg.totalCost, pg.totalCostUSD)
									}
									return pg.contentRow(gtx, values.String(values.StrTotalCost), totalCostText)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.contentRow(gtx, values.String(values.StrBalanceAfter), pg.balanceAfterSend)
							}),
						)
					})
				}),
				layout.Flexed(0.3, func(gtx C) D {
					return pg.nextButton.Layout(gtx)
				}),
			)
		})
	})
}

func (pg *Page) contentRow(gtx layout.Context, leftValue, rightValue string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(leftValue)
			txt.Color = pg.Theme.Color.GrayText2
			return txt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(pg.Theme.Body1(rightValue).Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Dimensions{}
					}),
				)
			})
		}),
	)
}
