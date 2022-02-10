package dexclient

import (
	"fmt"
	"image"
	"image/color"

	"decred.org/dcrdex/client/core"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	scaleHeight float32 = 1
	scaleWidth  float32 = 1e-4
	strokeWidth float32 = 2
)

// depthChart implements depth chart display sell, buy orders.
type depthChart struct {
	strokeBuyColor  color.NRGBA
	strokeSellColor color.NRGBA
	fillBuyColor    color.NRGBA
	fillSellColor   color.NRGBA

	zoom *zoomWdg
}

type zoomWdg struct {
	increase decredmaterial.IconButton
	decrease decredmaterial.IconButton
	text     decredmaterial.Label
	value    int
}

func newDepthChart(l *load.Load) *depthChart {
	increase := l.Theme.IconButton(l.Icons.ContentAdd)
	decrease := l.Theme.IconButton(l.Icons.ContentRemove)
	colorStyle := &values.ColorStyle{Background: l.Theme.Color.Danger, Foreground: l.Theme.Color.Text}
	increase.ChangeColorStyle(colorStyle)
	decrease.ChangeColorStyle(colorStyle)

	return &depthChart{
		strokeBuyColor:  l.Theme.Color.Success,
		fillBuyColor:    l.Theme.Color.Success2,
		strokeSellColor: l.Theme.Color.Danger,
		fillSellColor:   l.Theme.Color.Orange2,
		zoom: &zoomWdg{
			increase: increase,
			decrease: decrease,
			text:     l.Theme.Body1("5"),
			value:    5,
		},
	}
}

// Layout draws the Board and accepts input for adding alive cells.
func (chart depthChart) layout(gtx C, buys, sells []*core.MiniOrder) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y / 3
	chart.handler()
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return chart.depthChartLayout(gtx, buys, sells)
		}),
		layout.Stacked(func(gtx C) D {
			return chart.zoomButtonLayout(gtx)
		}),
	)
}

func (chart depthChart) zoomButtonLayout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return chart.zoom.decrease.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Left:  values.MarginPadding8,
					Right: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return chart.zoom.text.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return chart.zoom.increase.Layout(gtx)
			}),
		)
	})
}

func (chart depthChart) depthChartLayout(gtx C, buys, sells []*core.MiniOrder) D {
	c := func(miniOrders []*core.MiniOrder, isSell bool) layout.Widget {
		return func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return chart.filledLayout(gtx, miniOrders, isSell)
				}),
				layout.Stacked(func(gtx C) D {
					return chart.strokeLayout(gtx, miniOrders, isSell)
				}),
			)
		}
	}

	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	return layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Flexed(.5, c(buys, false)),
		layout.Flexed(.5, c(sells, true)),
	)
}

func (chart depthChart) strokeLayout(gtx C, orderBooks []*core.MiniOrder, isSell bool) D {
	var sizeX int
	var color color.NRGBA
	if isSell {
		sizeX = 0
		color = chart.strokeSellColor
	} else {
		sizeX = gtx.Constraints.Max.X
		color = chart.strokeBuyColor
	}

	size := image.Point{X: sizeX, Y: gtx.Constraints.Max.Y / 3}

	var x float32 = float32(size.X)
	var y float32 = float32(size.Y)
	var stroke clip.Path
	stroke.Begin(gtx.Ops)

	for i := 0; i < len(orderBooks); i++ {
		var lineLength float32 = 0
		var depth float32 = float32(orderBooks[i].Qty) * scaleHeight * float32(chart.zoom.value)
		if i < len(orderBooks)-1 {
			if isSell {
				lineLength = (float32(orderBooks[i+1].MsgRate) - float32(orderBooks[i].MsgRate)) * scaleWidth * float32(chart.zoom.value)
			} else {
				lineLength = (float32(orderBooks[i].MsgRate) - float32(orderBooks[i+1].MsgRate)) * scaleWidth * float32(chart.zoom.value)
			}
		}

		stroke.MoveTo(f32.Pt(x, y))
		y -= depth
		stroke.LineTo(f32.Pt(x, y))
		stroke.Close()
		stroke.MoveTo(f32.Pt(x, y))
		if isSell {
			x += lineLength
		} else {
			x -= lineLength
		}
		stroke.LineTo(f32.Pt(x, y))
		stroke.Close()
	}

	if isSell {
		stroke.MoveTo(f32.Pt(x, y))
		stroke.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), y))
	} else {
		stroke.MoveTo(f32.Pt(x, y))
		stroke.LineTo(f32.Pt(0, y))
	}
	stroke.Close()

	defer clip.Stroke{
		Path:  stroke.End(),
		Width: strokeWidth,
	}.Op().Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, color)

	return D{Size: size}
}

func (chart depthChart) filledLayout(gtx C, orderBooks []*core.MiniOrder, isSell bool) D {
	var sizeX int
	var color color.NRGBA
	var filled clip.Path

	if isSell {
		sizeX = 0
		color = chart.fillSellColor
	} else {
		sizeX = gtx.Constraints.Max.X
		color = chart.fillBuyColor
	}
	size := image.Point{X: sizeX, Y: gtx.Constraints.Max.Y / 3}

	var x float32 = float32(size.X)
	var y float32 = float32(size.Y)
	filled.Begin(gtx.Ops)

	for i := 0; i < len(orderBooks); i++ {
		var nextX float32
		var filledLength float32 = 0
		var depth float32 = float32(orderBooks[i].Qty) * scaleHeight * float32(chart.zoom.value)

		if i < len(orderBooks)-1 {
			if isSell {
				filledLength = (float32(orderBooks[i+1].MsgRate) - float32(orderBooks[i].MsgRate)) * scaleWidth * float32(chart.zoom.value)
			} else {
				filledLength = (float32(orderBooks[i].MsgRate) - float32(orderBooks[i+1].MsgRate)) * scaleWidth * float32(chart.zoom.value)
			}
		}

		if isSell {
			nextX = x + filledLength
		} else {
			nextX = x - filledLength
		}
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(nextX, float32(size.Y)))
		y -= depth
		filled.LineTo(f32.Pt(nextX, y))
		filled.LineTo(f32.Pt(x, y))
		if isSell {
			x += filledLength
		} else {
			x -= filledLength
		}
	}

	// Fill the rest of the chart
	if isSell {
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(size.Y)))
		filled.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), y))
		filled.LineTo(f32.Pt(x, y))
	} else {
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(0, float32(size.Y)))
		filled.LineTo(f32.Pt(0, y))
		filled.LineTo(f32.Pt(x, y))
	}
	filled.Close()

	defer clip.Outline{Path: filled.End()}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: size}
}

func (chart *depthChart) handler() {
	if chart.zoom.increase.Button.Clicked() {
		chart.zoom.value++
		chart.zoom.text.Text = fmt.Sprintf("%d", chart.zoom.value)
	}

	if chart.zoom.decrease.Button.Clicked() {
		if chart.zoom.value > 1 {
			chart.zoom.value--
			chart.zoom.text.Text = fmt.Sprintf("%d", chart.zoom.value)
		}
	}
}
