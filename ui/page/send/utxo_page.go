package send

import (
	"fmt"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

const UTXOPageID = "unspentTransactionOutput"

type UTXOPage struct {
	*load.Load
	utxoListContainer      layout.List
	txAuthor               *dcrlibwallet.TxAuthor
	backButton             decredmaterial.IconButton
	useUTXOButton          decredmaterial.Button
	unspentOutputs         **wallet.UnspentOutputs
	unspentOutputsSelected *map[int]map[int32]map[string]*wallet.UnspentOutput
	checkboxes             []decredmaterial.CheckBoxStyle
	copyButtons            []decredmaterial.IconButton
	selectAllChexBox       decredmaterial.CheckBoxStyle
	separator              decredmaterial.Line

	txnFee            string
	txnAmount         string
	txnAmountAfterFee string

	selectedWalletID  int
	selectedAccountID int32
}

func NewUTXOPage(l *load.Load, account *dcrlibwallet.Account) *UTXOPage {
	pg := &UTXOPage{
		Load:           l,
		unspentOutputs: &l.WL.UnspentOutputs,
		utxoListContainer: layout.List{
			Axis: layout.Vertical,
		},
		// txAuthor:               &l.WL.TxAuthor,
		unspentOutputsSelected: &l.SelectedUTXO,
		selectAllChexBox:       l.Theme.CheckBox(new(widget.Bool), ""),
		separator:              l.Theme.Separator(),
		selectedWalletID:       account.WalletID,
		selectedAccountID:      account.Number,
	}

	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)
	pg.useUTXOButton = l.Theme.Button("OK")

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *UTXOPage) ID() string {
	return UTXOPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *UTXOPage) OnNavigatedTo() {

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *UTXOPage) HandleUserInteractions() {
	if len(pg.checkboxes) != len((*pg.unspentOutputs).List) {
		pg.checkboxes = make([]decredmaterial.CheckBoxStyle, len((*pg.unspentOutputs).List))
		pg.copyButtons = make([]decredmaterial.IconButton, len((*pg.unspentOutputs).List))

		for i := 0; i < len((*pg.unspentOutputs).List); i++ {
			utxo := (*pg.unspentOutputs).List[i]
			pg.checkboxes[i] = pg.Theme.CheckBox(new(widget.Bool), "")
			if (*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID][utxo.UTXO.OutputKey] != nil {
				pg.checkboxes[i].CheckBox.Value = true
			}
			icoBtn := pg.Theme.IconButton(decredmaterial.MustIcon(widget.NewIcon(icons.ContentContentCopy)))
			icoBtn.Inset, icoBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding20
			icoBtn.ChangeColorStyle(&values.ColorStyle{Background: pg.Theme.Color.Gray4})
			pg.copyButtons[i] = icoBtn
		}
		pg.calculateAmountAndFeeUTXO()
	}

	if pg.backButton.Button.Clicked() {
		pg.clearPageData()
		pg.PopFragment()
	}

	/*if pg.useUTXOButton.Button.Clicked() {
		//TODO
		//pg.ChangePage(send.PageID)
	}
	*/

	if pg.selectAllChexBox.CheckBox.Changed() {
		for i, utxo := range (*pg.unspentOutputs).List {
			if pg.selectAllChexBox.CheckBox.Value {
				pg.checkboxes[i].CheckBox.Value = true
				(*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID][utxo.UTXO.OutputKey] = utxo
			} else {
				delete((*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID], utxo.UTXO.OutputKey)
				pg.checkboxes[i].CheckBox.Value = false
			}
		}
		pg.calculateAmountAndFeeUTXO()
	}
}

func (pg *UTXOPage) handlerCheckboxes(cb *decredmaterial.CheckBoxStyle, utxo *wallet.UnspentOutput) {
	if cb.CheckBox.Changed() {
		if cb.CheckBox.Value {
			(*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID][utxo.UTXO.OutputKey] = utxo
		} else {
			delete((*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID], utxo.UTXO.OutputKey)
		}
		pg.calculateAmountAndFeeUTXO()
	}
}

func (pg *UTXOPage) calculateAmountAndFeeUTXO() {
	unsignedTx, err := pg.WL.MultiWallet.NewUnsignedTx(pg.selectedWalletID, pg.selectedAccountID)
	if err != nil {
		return
	}

	var utxoKeys []string
	var totalAmount int64
	for utxoKey, utxo := range (*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID] {
		utxoKeys = append(utxoKeys, utxoKey)
		totalAmount += utxo.UTXO.Amount
	}

	err = unsignedTx.UseInputs(utxoKeys)
	if err != nil {
		return
	}
	feeAndSize, err := unsignedTx.EstimateFeeAndSize()
	if err != nil {
		return
	}
	pg.txnAmount = dcrutil.Amount(totalAmount).String()
	pg.txnFee = dcrutil.Amount(feeAndSize.Fee.AtomValue).String()
	pg.txnAmountAfterFee = dcrutil.Amount(totalAmount - feeAndSize.Fee.AtomValue).String()

	pg.txAuthor = unsignedTx
}

func (pg *UTXOPage) clearPageData() {
	pg.checkboxes = nil
	pg.txnFee = ""
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *UTXOPage) Layout(gtx C) D {
	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, pg.backButton.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding10,
							// Top:  values.MarginPaddingMinus10,
						}.Layout(gtx, pg.Theme.H6("Coin Control").Layout)
					}),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Flexed(0.25, func(gtx C) D {
										utxos := (*pg.unspentOutputsSelected)[pg.selectedWalletID][pg.selectedAccountID]
										return pg.textData(gtx, "Selected:  ", fmt.Sprintf("%d", len(utxos)))
									}),
									layout.Flexed(0.25, func(gtx C) D {
										return pg.textData(gtx, "Amount:  ", pg.txnAmount)
									}),
									layout.Flexed(0.25, func(gtx C) D {
										return pg.textData(gtx, "Fee:  ", pg.txnFee)
									}),
									layout.Flexed(0.25, func(gtx C) D {
										return pg.textData(gtx, "After Fee:  ", pg.txnAmountAfterFee)
									}),
								)
							})
						}),
						layout.Rigid(pg.separator.Layout),
						layout.Rigid(func(gtx C) D {
							return pg.utxoRowHeader(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							if len(pg.checkboxes) == 0 {
								return D{}
							}
							return pg.utxoListContainer.Layout(gtx, len((*pg.unspentOutputs).List), func(gtx C, index int) D {
								utxo := (*pg.unspentOutputs).List[index]
								pg.handlerCheckboxes(&pg.checkboxes[index], utxo)
								return pg.utxoRow(gtx, utxo, index)
							})
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.useUTXOButton.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (pg *UTXOPage) textData(gtx C, txt, value string) D {
	txt1 := pg.Theme.Label(values.MarginPadding15, txt)
	txt2 := pg.Theme.Label(values.MarginPadding15, value)
	txt1.MaxLines, txt2.MaxLines = 1, 1
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(txt1.Layout),
		layout.Rigid(txt2.Layout),
	)
}

func (pg *UTXOPage) utxoRowHeader(gtx C) D {
	txt := pg.Theme.Label(values.MarginPadding15, "")
	txt.MaxLines = 1
	return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(pg.selectAllChexBox.Layout),
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
	})
}

func (pg *UTXOPage) utxoRow(gtx C, data *wallet.UnspentOutput, index int) D {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(pg.checkboxes[index].Layout),
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(data.Amount)
			txt.MaxLines = 1
			txt.Alignment = text.Start
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(data.UTXO.Addresses)
			txt.MaxLines = 1
			gtx.Constraints.Max.X = gtx.Px(values.MarginPadding200)
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding200)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(data.DateTime)
			txt.MaxLines = 1
			txt.Alignment = text.End
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(fmt.Sprintf("%d", data.UTXO.Confirmations))
			txt.MaxLines = 1
			txt.Alignment = text.End
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if pg.copyButtons[index].Button.Clicked() {
				clipboard.WriteOp{Text: data.UTXO.Addresses}.Add(gtx.Ops)
			}
			return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, pg.copyButtons[index].Layout)
		}),
	)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *UTXOPage) OnNavigatedFrom() {}
