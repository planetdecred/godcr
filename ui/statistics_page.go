package ui

import (
	"strconv"
	"time"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageStat = "Stat"

type statItems struct {
	network, buildDate string
	walletDir          string
	numWallets, numTxs string
	walletDataSize     string
	syncStatus         *wallet.SyncStatus
	startupTime        time.Time
}

type statPage struct {
	divider *decredmaterial.Line
	theme   *decredmaterial.Theme
	statItems
}

func (win *Window) StatPage(common pageCommon) layout.Widget {
	l := common.theme.Line(2, 2)
	pg := &statPage{
		theme:   common.theme,
		divider: &l,
	}

	pg.divider.Color = common.theme.Color.Background

	pg.statItems.startupTime = time.Now()

	pg.syncStatus = win.walletSyncStatus

	return func(gtx C) D {
		return pg.Layout(gtx, common)
	}
}

func (pg *statPage) statItem(title, value string, gtx C, common pageCommon) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				return common.theme.Body1(title).Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
					return common.theme.Body1(value).Layout(gtx)
				})
			})
		}),
	)
}

func (pg *statPage) lineSeparator() layout.Widget {
	m := values.MarginPadding1
	return func(gtx C) D {
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
			pg.divider.Width = gtx.Constraints.Max.X
			return pg.divider.Layout(gtx)
		})
	}
}

func (pg *statPage) layoutStats(gtx C, common pageCommon) D {
	background := common.theme.Color.Surface
	card := common.theme.Card()
	card.Color = background
	item := func(t, v string) layout.Widget {
		return func(gtx C) D {
			return pg.statItem(t, v, gtx, common)
		}
	}

	now := time.Now()
	uptime := now.Sub(pg.startupTime)

	items := []layout.Widget{
		item("Build", pg.statItems.network+","+pg.statItems.buildDate),
		pg.lineSeparator(),
		item("Peers connected", strconv.Itoa(int(pg.statItems.syncStatus.ConnectedPeers))),
		pg.lineSeparator(),
		item("Uptime", uptime.String()),
		pg.lineSeparator(),
		item("Network", pg.statItems.network),
		pg.lineSeparator(),
		item("Best block", strconv.Itoa(int(pg.statItems.syncStatus.CurrentBlockHeight))),
		pg.lineSeparator(),
		item("Best block timestamp", "value"),
		pg.lineSeparator(),
		item("Best block age", "value"),
		pg.lineSeparator(),
		item("Wallet data directory", pg.statItems.walletDir),
		pg.lineSeparator(),
		item("Wallet data", pg.statItems.walletDataSize),
		pg.lineSeparator(),
		item("Transactions", pg.numTxs),
		pg.lineSeparator(),
		item("Wallets", pg.statItems.numWallets),
	}

	return card.Layout(gtx, func(gtx C) D {
		m15 := values.MarginPadding15
		return layout.Inset{Left: m15, Right: m15}.Layout(gtx, func(gtx C) D {
			l := layout.List{
				Axis: layout.Vertical,
			}
			return l.Layout(gtx, len(items), func(gtx C, i int) D {
				return items[i](gtx)
			})
		})
	})
}

func (pg *statPage) Layout(gtx C, common pageCommon) D {
	container := func(gtx C) D {
		page := SubPage{
			title: "Statistics",
			back: func() {
				common.changePage(PageDebug)
			},
			body: func(gtx C) D {
				return pg.layoutStats(gtx, common)
			},
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.Layout(gtx, container)
}
