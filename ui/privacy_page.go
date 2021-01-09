package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PagePrivacy = "Privacy"

type privacyPage struct {
	theme         *decredmaterial.Theme
	pageContainer layout.List
	backButton    decredmaterial.IconButton
	toggleMixer   *widget.Bool
	infoBtn       decredmaterial.IconButton
	line          *decredmaterial.Line
}

func (win *Window) PrivacyPage(common pageCommon) layout.Widget {
	pg := &privacyPage{
		theme:         common.theme,
		pageContainer: layout.List{Axis: layout.Vertical},
		backButton:    common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		toggleMixer:   new(widget.Bool),
		line:          common.theme.Line(),
	}
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)
	pg.infoBtn = common.theme.IconButton(new(widget.Clickable), common.icons.actionInfo)
	pg.infoBtn.Color = common.theme.Color.Gray
	pg.infoBtn.Background = common.theme.Color.Surface
	pg.infoBtn.Inset = layout.UniformInset(values.MarginPadding0)
	pg.line.Color = common.theme.Color.Background
	pg.line.Height = 1

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *privacyPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	widgets := []func(gtx C) D{
		func(gtx C) D {
			return pg.header(gtx, &common)
		},
		pg.gutter,
		func(gtx C) D {
			return pg.mixerInfoLayout(gtx, &common)
		},
		pg.gutter,
		func(gtx C) D {
			return pg.mixerSettingsLayout(gtx, &common)
		},
		pg.gutter,
		func(gtx C) D {
			return pg.dangerZoneLayout(gtx, &common)
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return pg.pageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
			return widgets[i](gtx)
		})
	})
}

func (pg *privacyPage) header(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return c.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
							return pg.backButton.Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.H6("Privacy")
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

func (pg *privacyPage) mixerInfoLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return c.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							c.icons.mixer.Scale = 0.06
							return c.icons.mixer.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := pg.theme.H6("Mixer")
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										txt := pg.theme.Body2("Ready to mix")
										txt.Color = c.theme.Color.Gray
										return txt.Layout(gtx)
									}),
								)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return material.Switch(pg.theme.Base, pg.toggleMixer).Layout(gtx)
						}),
					)
				}),
				layout.Rigid(pg.gutter),
				layout.Rigid(func(gtx C) D {
					content := c.theme.Card()
					content.Color = c.theme.Color.Background
					return content.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txt := c.theme.Label(values.TextSize14, "Unmixed balance")
											txt.Color = c.theme.Color.Gray
											return txt.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return c.theme.Body2("200 DCR").Layout(gtx)
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											t := c.theme.Label(values.TextSize14, "Mixed balance")
											t.Color = c.theme.Color.Gray
											return t.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return c.theme.Body2("0 DCR").Layout(gtx)
										}),
									)
								}),
							)
						})
					})
				}),
			)
		})
	})
}

func (pg *privacyPage) mixerSettingsLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return c.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		pg.line.Width = gtx.Constraints.Max.X

		row := func(txt1, txt2 string) D {
			return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return c.theme.Label(values.TextSize18, txt1).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return c.theme.Body2(txt2).Layout(gtx)
					}),
				)
			})
		}

		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
					return c.theme.Body2("Mixer Settings").Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D { return row("Mixed account", "mixed") }),
			layout.Rigid(func(gtx C) D { return pg.line.Layout(gtx) }),
			layout.Rigid(func(gtx C) D { return row("Change account", "unmixed") }),
			layout.Rigid(func(gtx C) D { return pg.line.Layout(gtx) }),
			layout.Rigid(func(gtx C) D { return row("Account branch", "0") }),
			layout.Rigid(func(gtx C) D { return pg.line.Layout(gtx) }),
			layout.Rigid(func(gtx C) D { return row("Shuffle server", "cspp.decred.org") }),
			layout.Rigid(func(gtx C) D { return pg.line.Layout(gtx) }),
			layout.Rigid(func(gtx C) D { return row("Shuffle port", "15760") }),
		)
	})
}

func (pg *privacyPage) dangerZoneLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return c.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.TextSize16).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.Label(values.MarginPadding15, "Danger Zone")
					txt.Color = c.theme.Color.Gray
					return txt.Layout(gtx)
				}),
			)
		})
	})
}

func (pg *privacyPage) gutter(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Dimensions{}
	})
}

func (pg *privacyPage) Handler(common pageCommon) {
	if pg.backButton.Button.Clicked() {
		*common.page = PageWallet
	}
}
