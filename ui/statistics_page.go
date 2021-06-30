package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageStat = "Stat"

type statPage struct {
	common      *pageCommon
	txs         **wallet.Transactions
	theme       *decredmaterial.Theme
	l           layout.List
	startupTime time.Time
	syncStatus  *wallet.SyncStatus
	netType     string

	backButton decredmaterial.IconButton
}

func StatPage(common *pageCommon) Page {
	pg := &statPage{
		txs:    common.walletTransactions,
		common: common,
		theme:  common.theme,
		l: layout.List{
			Axis: layout.Vertical,
		},
	}
	pg.startupTime = time.Now()
	pg.syncStatus = common.walletSyncStatus
	if common.wallet.Net == "testnet3" {
		pg.netType = "Testnet"
	} else {
		pg.netType = strings.Title(common.wallet.Net)
	}

	pg.backButton, _ = common.SubPageHeaderButtons()

	return pg
}

func (pg *statPage) OnResume() {

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
			l := pg.theme.Body2(t)
			r := pg.theme.Body2(v)
			r.Color = pg.theme.Color.Gray
			return inset.Layout(gtx, func(gtx C) D {
				return endToEndRow(gtx, l.Layout, r.Layout)
			})
		}
	}
	uptime := func(t time.Time) string {
		v := int(time.Since(t).Seconds())
		h := v / 3600
		m := (v - h*3600) / 60
		s := v - h*3600 - m*60
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}(pg.startupTime)

	items := []layout.Widget{
		item("Build", pg.netType+", "+time.Now().Format("2006-01-02")),
		pg.theme.Separator().Layout,
		item("Peers connected", strconv.Itoa(int(pg.syncStatus.ConnectedPeers))),
		pg.theme.Separator().Layout,
		item("Uptime", uptime),
		pg.theme.Separator().Layout,
		item("Network", pg.netType),
		pg.theme.Separator().Layout,
		item("Best block", fmt.Sprintf("%d", pg.common.info.BestBlockHeight)),
		pg.theme.Separator().Layout,
		item("Best block timestamp", time.Unix(pg.common.info.BestBlockTime, 0).Format("2006-01-02 03:04:05 -0700")),
		pg.theme.Separator().Layout,
		item("Best block age", pg.common.info.LastSyncTime),
		pg.theme.Separator().Layout,
		item("Wallet data directory", pg.common.wallet.WalletDirectory()),
		pg.theme.Separator().Layout,
		item("Wallet data", pg.common.wallet.DataSize()),
		pg.theme.Separator().Layout,
		item("Transactions", fmt.Sprintf("%d", (*pg.txs).Total)),
		pg.theme.Separator().Layout,
		item("Wallets", fmt.Sprintf("%d", len(pg.common.info.Wallets))),
	}

	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
			return pg.l.Layout(gtx, len(items), func(gtx C, i int) D {
				return items[i](gtx)
			})
		})
	})
}

func (pg *statPage) Layout(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		page := SubPage{
			title:      "Statistics",
			backButton: pg.backButton,
			back: func() {
				pg.common.changePage(PageDebug)
			},
			body: pg.layoutStats,
		}
		return pg.common.SubPageLayout(gtx, page)
	}

	// Refresh frames every 1 second
	op.InvalidateOp{At: time.Now().Add(time.Second * 1)}.Add(gtx.Ops)
	return pg.common.UniformPadding(gtx, container)
}

func (pg *statPage) handle()  {}
func (pg *statPage) onClose() {}
