package ui

import (
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
	utxoContainer      layout.List
	backButton         decredmaterial.IconButton
	useUTXOButton      decredmaterial.Button
	outputsCollapsible *decredmaterial.Collapsible
	inputsCollapsible  *decredmaterial.Collapsible
	unspentOutputs     **wallet.UnspentOutputs
	checkbox           decredmaterial.CheckBoxStyle
}

func (win *Window) UTXOPage(common pageCommon) layout.Widget {
	pg := &utxoPage{
		unspentOutputs: &win.walletUnspentOutputs,
		utxoPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		utxoContainer: layout.List{
			Axis: layout.Vertical,
		},
		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),
		checkbox:           common.theme.CheckBox(new(widget.Bool), "test checkbox"),
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

func (pg *utxoPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return c.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.W.Layout(gtx, func(gtx C) D {
					return pg.backButton.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.txnRowHeader(gtx, &c)
			}),
			layout.Flexed(1, func(gtx C) D {
				return pg.utxoContainer.Layout(gtx, len((*pg.unspentOutputs).List), func(gtx C, index int) D {
					utxo := (*pg.unspentOutputs).List[index]
					return c.theme.Body1(utxo.Address).Layout(gtx)
				})
				// return pg.checkbox.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.useUTXOButton.Layout(gtx)
			}),
		)
	})
}

func (pg *utxoPage) txnRowHeader(gtx layout.Context, common *pageCommon) layout.Dimensions {
	txt := common.theme.Label(values.MarginPadding15, "#")
	txt.Color = common.theme.Color.Hint

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Alignment = text.Middle
			txt.Text = "Date (UTC)"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Text = "Status"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Text = "Amount"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Text = "Fee"
			return txt.Layout(gtx)
		}),
	)
}

func (pg *utxoPage) Handler(common pageCommon) {
	if pg.backButton.Button.Clicked() {
		*common.page = PageSend
	}
}
