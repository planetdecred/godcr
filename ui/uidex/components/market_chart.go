package components

import (
	"fmt"
	"image"
	"image/color"

	"decred.org/dcrdex/client/core"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	scaleHeight        float32 = 0.025
	scaleWidth         float32 = 8
	strokeWidth        float32 = 2
	minimumChartHeight int     = 220
)

// DepthChart implements depth chart sell or buy logic.
type DepthChart struct {
	buys  []*core.MiniOrder
	sells []*core.MiniOrder

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

func NewDepthChart(buys, sells []*core.MiniOrder, theme *decredmaterial.Theme) DepthChart {
	increase := theme.IconButton(new(widget.Clickable), utils.MustIcon(widget.NewIcon(icons.ContentAdd)))
	decrease := theme.IconButton(new(widget.Clickable), utils.MustIcon(widget.NewIcon(icons.ContentRemove)))
	increase.Background, decrease.Background = color.NRGBA{}, color.NRGBA{}
	increase.Color, decrease.Color = theme.Color.Text, theme.Color.Text

	return DepthChart{
		buys:            buys,
		sells:           sells,
		strokeBuyColor:  theme.Color.ChartBuyLine,
		fillBuyColor:    theme.Color.ChartBuyFill,
		strokeSellColor: theme.Color.ChartSellLine,
		fillSellColor:   theme.Color.ChartSellFill,
		zoom: &zoomWdg{
			increase: increase,
			decrease: decrease,
			text:     theme.Body1("5"),
			value:    5,
		},
	}
}

// Layout draws the Board and accepts input for adding alive cells.
func (chart DepthChart) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = minimumChartHeight

	// cut layout if go outside
	clip.Rect{Max: gtx.Constraints.Max}.Add(gtx.Ops)
	chart.handler()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return chart.depthChartLayout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return chart.zoomButtonLayout(gtx)
		}),
	)
}

func (chart DepthChart) zoomButtonLayout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return chart.zoom.decrease.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Left:  values.MarginPadding8,
					Right: values.MarginPadding8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return chart.zoom.text.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return chart.zoom.increase.Layout(gtx)
			}),
		)
	})
}

func (chart DepthChart) depthChartLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.filledLayout(gtx, chart.buys, false)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.strokeLayout(gtx, chart.buys, false)
				}),
			)
		}),
		layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.filledLayout(gtx, chart.sells, true)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return chart.strokeLayout(gtx, chart.sells, true)
				}),
			)
		}),
	)
}

func (chart DepthChart) strokeLayout(gtx layout.Context, orderBooks []*core.MiniOrder, isSell bool) layout.Dimensions {
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
		var depth float32 = float32(orderBooks[i].Qty) * scaleHeight * float32(chart.zoom.value)
		if i < len(orderBooks)-1 {
			if isSell {
				lineLength = (float32(orderBooks[i+1].Rate) - float32(orderBooks[i].Rate)) * scaleWidth * float32(chart.zoom.value)
			} else {
				lineLength = (float32(orderBooks[i].Rate) - float32(orderBooks[i+1].Rate)) * scaleWidth * float32(chart.zoom.value)
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

	clip.Stroke{
		Path:  stroke.End(),
		Style: clip.StrokeStyle{Width: strokeWidth},
	}.Op().Add(gtx.Ops)

	paint.Fill(gtx.Ops, color)

	return layout.Dimensions{Size: size}
}

func (chart DepthChart) filledLayout(gtx layout.Context, orderBooks []*core.MiniOrder, isSell bool) layout.Dimensions {
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
		var depth float32 = float32(orderBooks[i].Qty) * scaleHeight * float32(chart.zoom.value) // quanitty

		if i < len(orderBooks)-1 {
			if isSell {
				filledLength = (float32(orderBooks[i+1].Rate) - float32(orderBooks[i].Rate)) * scaleWidth * float32(chart.zoom.value)
			} else {
				filledLength = (float32(orderBooks[i].Rate) - float32(orderBooks[i+1].Rate)) * scaleWidth * float32(chart.zoom.value)
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

	clip.Outline{Path: filled.End()}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (chart *DepthChart) handler() {
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
