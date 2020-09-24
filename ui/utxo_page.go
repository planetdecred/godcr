package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageUTXO = "unspentTransactionOutput"

type utxoPage struct {
	utxoPageContainer  layout.List
	utxoListContainer  layout.List
	backButton         decredmaterial.IconButton
	useUTXOButton      decredmaterial.Button
	outputsCollapsible *decredmaterial.Collapsible
	inputsCollapsible  *decredmaterial.Collapsible
	unspentOutputs     **wallet.UnspentOutputs
	checkboxes         []decredmaterial.CheckBoxStyle
}

func (win *Window) UTXOPage(common pageCommon) layout.Widget {
	pg := &utxoPage{
		unspentOutputs: &win.walletUnspentOutputs,
		utxoPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		utxoListContainer: layout.List{
			Axis: layout.Vertical,
		},
		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),
	}

	pg.backButton = common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack)
	pg.backButton.Color = common.theme.Color.Hint
	pg.backButton.Size = values.MarginPadding30
	pg.useUTXOButton = common.theme.Button(new(widget.Clickable), "OK")

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *utxoPage) Handler(common pageCommon) {
	if len(pg.checkboxes) != len((*pg.unspentOutputs).List) {
		pg.checkboxes = make([]decredmaterial.CheckBoxStyle, len((*pg.unspentOutputs).List))
		for i := 0; i < len((*pg.unspentOutputs).List); i++ {
			pg.checkboxes[i] = common.theme.CheckBox(new(widget.Bool), "")
		}
	}
	if pg.backButton.Button.Clicked() {
		*common.page = PageSend
	}
}

func (pg *utxoPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return c.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							return pg.backButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding10, Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return c.theme.H5("Coin Control").Layout(gtx)
						})
					}),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.utxoRowHeader(gtx, &c)
						}),
						layout.Flexed(1, func(gtx C) D {
							return pg.utxoListContainer.Layout(gtx, len((*pg.unspentOutputs).List), func(gtx C, index int) D {
								utxo := (*pg.unspentOutputs).List[index]
								return pg.utxoRow(gtx, utxo, &c, index)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.useUTXOButton.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (pg *utxoPage) utxoRowHeader(gtx layout.Context, c *pageCommon) layout.Dimensions {
	txt := c.theme.Label(values.MarginPadding15, "")
	txt.MaxLines = 1
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding35)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Text = "Amount"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding200)
			txt.Text = "Address"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			txt.Text = "Date (UTC)"
			txt.Alignment = text.End
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			txt.Text = "Confirmations"
			return txt.Layout(gtx)
		}),
	)
}

func (pg *utxoPage) utxoRow(gtx layout.Context, data *wallet.UnspentOutput, c *pageCommon, index int) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.checkboxes[index].Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := c.theme.Body1(data.Amount)
				txt.MaxLines = 1
				txt.Alignment = text.Start
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := c.theme.Body1(data.UTXO.Address)
				txt.MaxLines = 1
				gtx.Constraints.Max.X = gtx.Px(values.MarginPadding200)
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding200)
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := c.theme.Body1(data.DateTime)
				txt.MaxLines = 1
				txt.Alignment = text.End
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := c.theme.Body1(fmt.Sprintf("%d", data.UTXO.Confirmations))
				txt.MaxLines = 1
				txt.Alignment = text.End
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
				return txt.Layout(gtx)
			}),
		)
	})
}
