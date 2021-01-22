package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PagePrivacy = "Privacy"

type privacyPage struct {
	theme                                *decredmaterial.Theme
	pageContainer                        layout.List
	toggleMixer, allowUnspendUnmixedAcct *widget.Bool
	infoBtn                              decredmaterial.IconButton
	line                                 *decredmaterial.Line
	toPrivacySetup                       decredmaterial.Button
	privacyPageSetupVisibility           bool
	dangerZoneCollapsible                *decredmaterial.Collapsible
	errChann                             chan error
}

func (win *Window) PrivacyPage(common pageCommon) layout.Widget {
	pg := &privacyPage{
		theme:                   common.theme,
		pageContainer:           layout.List{Axis: layout.Vertical},
		toggleMixer:             new(widget.Bool),
		allowUnspendUnmixedAcct: new(widget.Bool),
		line:                    common.theme.Line(),
		toPrivacySetup:          common.theme.Button(new(widget.Clickable), "Set up mixer for this wallet"),
		dangerZoneCollapsible:   common.theme.Collapsible(),
		errChann:                common.errorChannels[PagePrivacy],
	}
	pg.toPrivacySetup.Background = pg.theme.Color.Primary
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

func (pg *privacyPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	d := func(gtx C) D {
		load := SubPage{
			title:      "Privacy",
			walletName: c.info.Wallets[*c.selectedWallet].Name,
			back: func() {
				*c.page = PageWallet
			},
			infoTemplateTitle: "How to use the mixer?",
			infoTemplate:      PrivacyInfoTemplate,
			body: func(gtx layout.Context) layout.Dimensions {
				if c.wallet.ReadyToMix(c.info.Wallets[*c.selectedWallet].ID) {
					widgets := []func(gtx C) D{
						func(gtx C) D {
							return pg.mixerInfoLayout(gtx, &c)
						},
						pg.gutter,
						func(gtx C) D {
							return pg.mixerSettingsLayout(gtx, &c)
						},
						pg.gutter,
						func(gtx C) D {
							return pg.dangerZoneLayout(gtx, &c)
						},
					}
					return pg.pageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return widgets[i](gtx)
					})
				}
				return pg.privacyIntroLayout(gtx, &c)
			},
		}
		return c.SubPageLayout(gtx, load)
	}
	return c.Layout(gtx, d)
}

func (pg *privacyPage) privacyIntroLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
										c.icons.transactionFingerPrintIcon.Scale = 0.09
										return c.icons.transactionFingerPrintIcon.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									c.icons.arrowFowardIcon.Scale = 0.18
									return c.icons.arrowFowardIcon.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
										c.icons.mixer.Scale = 0.25
										return c.icons.mixer.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									c.icons.arrowFowardIcon.Scale = 0.18
									return c.icons.arrowFowardIcon.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
										c.icons.transactionIcon.Scale = 0.09
										return c.icons.transactionIcon.Layout(gtx)
									})
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.theme.H6("How does CoinShuffle++ mixer enhance your privacy?")
							txt2 := pg.theme.Body1("CoinShuffle++ mixer can mix your DCRs through CoinJoin transactions.")
							txt3 := pg.theme.Body1("Using mixed DCRs protects you from exposing your financial activities to the public (e.g. how much you own, who pays you).")
							txt.Alignment, txt2.Alignment, txt3.Alignment = text.Middle, text.Middle, text.Middle

							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(txt.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt2.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt3.Layout)
								}),
							)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return pg.toPrivacySetup.Layout(gtx)
				})
			}),
		)
	})
}

func (pg *privacyPage) mixerInfoLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return c.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							c.icons.mixer.Scale = 0.05
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
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
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
			return layout.Inset{
				Left:   values.MarginPadding15,
				Right:  values.MarginPadding15,
				Top:    values.MarginPadding10,
				Bottom: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return c.theme.Label(values.TextSize16, txt1).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return c.theme.Body2(txt2).Layout(gtx)
					}),
				)
			})
		}

		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
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
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return pg.dangerZoneCollapsible.Layout(gtx,
				func(gtx C) D {
					txt := pg.theme.Label(values.MarginPadding15, "Danger Zone")
					txt.Color = c.theme.Color.Gray
					return txt.Layout(gtx)
				},
				func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, func(gtx C) D {
								return c.theme.Label(values.TextSize16, "Allow spending from unmixed accounts").Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return material.Switch(pg.theme.Base, pg.allowUnspendUnmixedAcct).Layout(gtx)
							}),
						)
					})
				},
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
	if pg.toPrivacySetup.Button.Clicked() {
		go pg.showModalSetupMixerInfo(&common)
	}

	select {
	case err := <-pg.errChann:
		common.Notify(err.Error(), false)
	default:
	}
}

func (pg *privacyPage) showModalSetupMixerInfo(common *pageCommon) {
	common.modalReceiver <- &modalLoad{
		template: ConfirmSetupMixerTemplate,
		title:    "Set up mixer by creating two needed accounts",
		confirm: func() {
			go pg.showModalSetupMixerAcct(common)
		},
		confirmText: "Begin setup",
		cancel:      common.closeModal,
		cancelText:  "Cancel",
	}
}

func (pg *privacyPage) showModalSetupMixerAcct(common *pageCommon) {
	common.modalReceiver <- &modalLoad{
		template: ConfirmSetupMixerAcctTemplate,
		title:    "Confirm to create needed accounts",
		confirm: func(p string) {
			for _, acct := range common.info.Wallets[*common.selectedWallet].Accounts {
				if acct.Name == "mixed" || acct.Name == "unmixed" {
					go pg.showModalSetupExistAcct(common, acct.Name)
					return
				}
			}
			common.wallet.SetupAccountMixer(common.info.Wallets[*common.selectedWallet].ID, p, pg.errChann)
		},
		confirmText: "Confirm",
		cancel:      common.closeModal,
		cancelText:  "Cancel",
	}
}

func (pg *privacyPage) showModalSetupExistAcct(common *pageCommon, acctName string) {
	common.modalReceiver <- &modalLoad{
		template: ConfirmMixerAcctExistTemplate,
		title:    fmt.Sprintf("Account “%s” is taken", acctName),
		confirm: func() {
			common.closeModal()
			*common.page = PageWallet
		},
		confirmText: "Go back & rename",
		cancel:      common.closeModal,
	}
}
