package ui

import (
	"fmt"
	"strconv"
	"time"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageStat = "Stat"

type statPage struct {
	common      pageCommon
	txs         **wallet.Transactions
	divider     *decredmaterial.Line
	theme       *decredmaterial.Theme
	l           layout.List
	startupTime time.Time
	syncStatus  *wallet.SyncStatus
}

func (win *Window) StatPage(common pageCommon) Page {
	l := common.theme.Line(2, 2)
	pg := &statPage{
		txs:     &win.walletTransactions,
		common:  common,
		theme:   common.theme,
		divider: &l,
		l: layout.List{
			Axis: layout.Vertical,
		},
	}

	pg.divider.Color = common.theme.Color.Background

	pg.startupTime = time.Now()
	pg.syncStatus = win.walletSyncStatus

	return pg
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

func (pg *statPage) layoutStats(gtx C) D {
	background := pg.common.theme.Color.Surface
	card := pg.common.theme.Card()
	card.Color = background
	inset := layout.Inset{
		Top:    values.MarginPadding12,
		Bottom: values.MarginPadding12,
		Right:  values.MarginPadding16,
		Left:   values.MarginPadding16,
	}

	item := func(t, v string) layout.Widget {
		return func(gtx C) D {
			l := pg.theme.Label(values.TextSize14, t)
			r := pg.theme.Label(values.TextSize14, v)
			r.Color = pg.theme.Color.Gray
			return inset.Layout(gtx, func(gtx C) D {
				return endToEndRow(gtx, l.Layout, r.Layout)
			})
		}
	}
	items := []layout.Widget{
		item("Build", pg.common.wallet.Net+", "+time.Now().Format("2006-01-02")),
		pg.lineSeparator(),
		item("Peers connected", strconv.Itoa(int(pg.syncStatus.ConnectedPeers))),
		pg.lineSeparator(),
		item("Uptime", time.Since(pg.startupTime).String()),
		pg.lineSeparator(),
		item("Network", pg.common.wallet.Net),
		pg.lineSeparator(),
		item("Best block", fmt.Sprintf("%d", pg.common.info.BestBlockHeight)),
		pg.lineSeparator(),
		item("Best block timestamp", time.Unix(pg.common.info.BestBlockTime, 0).Format("2006-01-02 03:04:05 -0700")),
		pg.lineSeparator(),
		item("Best block age", pg.common.info.LastSyncTime),
		pg.lineSeparator(),
		item("Wallet data directory", pg.common.wallet.WalletDirectory()),
		pg.lineSeparator(),
		item("Wallet data", pg.common.wallet.DataSize()),
		pg.lineSeparator(),
		item("Transactions", fmt.Sprintf("%d", (*pg.txs).Total)),
		pg.lineSeparator(),
		item("Wallets", fmt.Sprintf("%d", len(pg.common.info.Wallets))),
	}

	return card.Layout(gtx, func(gtx C) D {
		m15 := values.MarginPadding15
		return layout.Inset{Left: m15, Right: m15}.Layout(gtx, func(gtx C) D {
			return pg.l.Layout(gtx, len(items), func(gtx C, i int) D {
				return items[i](gtx)
			})
		})
	})
}

func (pg *statPage) Layout(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		page := SubPage{
			title: "Statistics",
			back: func() {
				pg.common.changePage(PageDebug)
			},
			body: func(gtx C) D {
				return pg.layoutStats(gtx)
			},
		}
		return pg.common.SubPageLayout(gtx, page)
	}
	return pg.common.Layout(gtx, func(gtx C) D {
		return pg.common.UniformPadding(gtx, container)
	})
}

func (pg *statPage) handle()  {}
func (pg *statPage) onClose() {}
