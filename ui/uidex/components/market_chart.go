package components

import (
	"image"
	"image/color"

	"decred.org/dcrdex/client/core"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type depthChartStyle struct {
	strokeBuyColor  color.NRGBA
	strokeSellColor color.NRGBA
	fillBuyColor    color.NRGBA
	FillSellColor   color.NRGBA
}

// DepthChart implements depth chart sell or buy logic.
type DepthChart struct {
	buys  []*core.MiniOrder
	sells []*core.MiniOrder
	depthChartStyle
}

func NewDepthChart(buys, sells []*core.MiniOrder, theme *decredmaterial.Theme) DepthChart {
	return DepthChart{
		buys,
		sells,
		depthChartStyle{
			strokeBuyColor:  theme.Color.ChartBuyLine,
			fillBuyColor:    theme.Color.ChartBuyFill,
			strokeSellColor: theme.Color.ChartSellLine,
			FillSellColor:   theme.Color.ChartSellFill,
		},
	}
}

var scaleDepth float32 = 35
var scaleWidth float32 = 8
var strokeWidth float32 = 2
var chartHeight = 400

// Layout draws the Board and accepts input for adding alive cells.
func (chart DepthChart) Layout(gtx layout.Context) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	return layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.layoutFilled(gtx, chart.buys, false)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.layoutStroke(gtx, chart.buys, false)
				}),
			)
		}),
		layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.layoutFilled(gtx, chart.sells, true)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.layoutStroke(gtx, chart.sells, true)
				}),
			)
		}),
	)
}

func (chart DepthChart) layoutStroke(gtx layout.Context, orderBooks []*core.MiniOrder, isSell bool) layout.Dimensions {
	sizeX := gtx.Constraints.Max.X
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
		var depth float32 = float32(orderBooks[i].Qty) / scaleDepth // quanitty
		if i < len(orderBooks)-1 {
			if isSell {
				lineLength = (float32(orderBooks[i+1].Rate) - float32(orderBooks[i].Rate)) * scaleWidth // rate
			} else {
				lineLength = (float32(orderBooks[i].Rate) - float32(orderBooks[i+1].Rate)) * scaleWidth // rate
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

	clip.Stroke{
		Path:  stroke.End(),
		Style: clip.StrokeStyle{Width: strokeWidth},
	}.Op().Add(gtx.Ops)

	paint.Fill(gtx.Ops, color)

	return layout.Dimensions{Size: size}
}

func (chart DepthChart) layoutFilled(gtx layout.Context, orderBooks []*core.MiniOrder, isSell bool) layout.Dimensions {
	var sizeX int
	var color color.NRGBA
	var filled clip.Path

	if isSell {
		sizeX = 0
		color = chart.FillSellColor
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
		var depth float32 = float32(orderBooks[i].Qty) / scaleDepth // quanitty

		if i < len(orderBooks)-1 {
			if isSell {
				filledLength = (float32(orderBooks[i+1].Rate) - float32(orderBooks[i].Rate)) * scaleWidth // rate
			} else {
				filledLength = (float32(orderBooks[i].Rate) - float32(orderBooks[i+1].Rate)) * scaleWidth // rate
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
		filled.Close()
		if isSell {
			x += filledLength
		} else {
			x -= filledLength
		}
	}

	clip.Outline{Path: filled.End()}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
