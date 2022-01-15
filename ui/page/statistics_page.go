package page

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const StatisticsPageID = "Statistics"

type StatPage struct {
	*load.Load
	txs           *wallet.Transactions
	l             layout.List
	scrollbarList *widget.List
	startupTime   string
	netType       string

	backButton decredmaterial.IconButton
}

func NewStatPage(l *load.Load) *StatPage {
	pg := &StatPage{
		Load: l,
		txs:  l.WL.Transactions,
		l:    layout.List{Axis: layout.Vertical},
		scrollbarList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		netType: l.WL.Wallet.Net,
	}
	if pg.netType == dcrlibwallet.Testnet3 {
		pg.netType = "Testnet"
	} else {
		pg.netType = strings.Title(pg.netType)
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *StatPage) ID() string {
	return StatisticsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *StatPage) OnNavigatedTo() {
	pg.appStartTime()
}

func (pg *StatPage) layoutStats(gtx C) D {
	background := pg.Theme.Color.Surface
	card := pg.Theme.Card()
	card.Color = background
	inset := layout.Inset{
		Top:    values.MarginPadding12,
		Bottom: values.MarginPadding12,
		Right:  values.MarginPadding16,
		Left:   values.MarginPadding16,
	}

	item := func(t, v string) layout.Widget {
		return func(gtx C) D {
			l := pg.Theme.Body2(t)
			r := pg.Theme.Body2(v)
			r.Color = pg.Theme.Color.GrayText2
			return inset.Layout(gtx, func(gtx C) D {
				return components.EndToEndRow(gtx, l.Layout, r.Layout)
			})
		}
	}

	bestBlock := pg.WL.MultiWallet.GetBestBlock()
	bestBlockTime := time.Unix(bestBlock.Timestamp, 0)
	secondsSinceBestBlock := int64(time.Since(bestBlockTime).Seconds())

	items := []layout.Widget{
		item("Build", pg.netType+", "+time.Now().Format("2006-01-02")),
		pg.Theme.Separator().Layout,
		item("Peers connected", strconv.Itoa(int(pg.WL.MultiWallet.ConnectedPeers()))),
		pg.Theme.Separator().Layout,
		item("Uptime", pg.startupTime),
		pg.Theme.Separator().Layout,
		item("Network", pg.netType),
		pg.Theme.Separator().Layout,
		item("Best block", fmt.Sprintf("%d", bestBlock.Height)),
		pg.Theme.Separator().Layout,
		item("Best block timestamp", bestBlockTime.Format("2006-01-02 03:04:05 -0700")),
		pg.Theme.Separator().Layout,
		item("Best block age", wallet.SecondsToDays(secondsSinceBestBlock)),
		pg.Theme.Separator().Layout,
		item("Wallet data directory", pg.WL.WalletDirectory()),
		pg.Theme.Separator().Layout,
		item("Wallet data", pg.WL.DataSize()),
		pg.Theme.Separator().Layout,
		item("Transactions", fmt.Sprintf("%d", (*pg.txs).Total)),
		pg.Theme.Separator().Layout,
		item("Wallets", fmt.Sprintf("%d", pg.WL.MultiWallet.LoadedWalletsCount())),
	}

	return pg.Theme.List(pg.scrollbarList).Layout(gtx, 1, func(gtx C, i int) D {
		return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					return pg.l.Layout(gtx, len(items), func(gtx C, i int) D {
						return items[i](gtx)
					})
				})
			})
		})
	})
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *StatPage) Layout(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Statistics",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: pg.layoutStats,
		}
		return sp.Layout(gtx)
	}

	// Refresh frames every 1 second
	op.InvalidateOp{At: time.Now().Add(time.Second * 1)}.Add(gtx.Ops)
	return components.UniformPadding(gtx, container)
}

func (pg *StatPage) appStartTime() {
	pg.startupTime = func(t time.Time) string {
		v := int(time.Since(t).Seconds())
		h := v / 3600
		m := (v - h*3600) / 60
		s := v - h*3600 - m*60
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}(pg.WL.Wallet.StartupTime())
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *StatPage) HandleUserInteractions() {
	pg.appStartTime()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *StatPage) OnNavigatedFrom() {}
