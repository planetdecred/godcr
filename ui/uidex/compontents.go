// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package uidex

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
	"github.com/planetdecred/godcr/ui/values"
)

// layoutTopBar is the top horizontal bar on every page of the app. It lays out the wallet balance, receive and send
// buttons.
func (page pageCommon) layoutTopBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding16
							v := values.MarginPadding10
							return Container{padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := page.icons.logo
											img.Scale = .3
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.theme.H5("DCRDEX").Layout(gtx)
											})
										}),
									)
								})
						})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
									// header buttons container
									return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
														func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																img := page.appBarNavItems[i].image
																img.Scale = 1.0
																return page.appBarNavItems[i].image.Layout(gtx)
															})
														})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding0,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return page.theme.Body1(page.appBarNavItems[i].page).Layout(gtx)
														})
													})
												}),
											)
										})
									})
								})
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return page.theme.Separator().Layout(gtx)
			}),
		)
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return leftWidget(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightWidget(gtx)
			})
		}),
	)
}

type depthChartStyle struct {
	strokeBuyColor  color.NRGBA
	strokeSellColor color.NRGBA
	fillBuyColor    color.NRGBA
	fillSellColor   color.NRGBA
}

// DepthChart implements depth chart sell or buy logic.
type DepthChart struct {
	buys  []*core.MiniOrder
	sells []*core.MiniOrder
	depthChartStyle
}

func NewDepthChart(buys, sells []*core.MiniOrder, style depthChartStyle) DepthChart {
	return DepthChart{
		buys,
		sells,
		style,
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
