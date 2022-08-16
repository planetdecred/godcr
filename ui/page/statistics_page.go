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
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const StatisticsPageID = "Statistics"

type StatPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	txs      []dcrlibwallet.Transaction
	accounts *dcrlibwallet.Accounts

	l             layout.List
	scrollbarList *widget.List
	startupTime   string
	netType       string

	backButton decredmaterial.IconButton
}

func NewStatPage(l *load.Load) *StatPage {
	pg := &StatPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(StatisticsPageID),
		l:                layout.List{Axis: layout.Vertical},
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

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *StatPage) OnNavigatedTo() {
	txs, err := pg.WL.SelectedWallet.Wallet.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterAll, true)
	if err != nil {
		log.Errorf("Error getting txs: %s", err.Error())
	} else {
		pg.txs = txs
	}

	acc, err := pg.WL.SelectedWallet.Wallet.GetAccountsRaw()
	if err != nil {
		log.Errorf("Error getting wallet accounts: %s", err.Error())
	} else {
		pg.accounts = acc
	}

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

	bestBlock := pg.WL.SelectedWallet.Wallet.GetBestBlock()
	bestBlockTime := time.Unix(pg.WL.SelectedWallet.Wallet.GetBestBlockTimeStamp(), 0)
	secondsSinceBestBlock := int64(time.Since(bestBlockTime).Seconds())

	line := pg.Theme.Separator()
	line.Color = pg.Theme.Color.Gray2

	items := []layout.Widget{
		item(values.String(values.StrBuild), pg.netType+", "+time.Now().Format("2006-01-02")),
		line.Layout,
		item(values.String(values.StrPeersConnected), strconv.Itoa(int(pg.WL.MultiWallet.ConnectedPeers()))),
		line.Layout,
		item(values.String(values.StrUptime), pg.startupTime),
		line.Layout,
		item(values.String(values.StrNetwork), pg.netType),
		line.Layout,
		item(values.String(values.StrBestBlocks), fmt.Sprintf("%d", bestBlock)),
		line.Layout,
		item(values.String(values.StrBestBlockTimestamp), bestBlockTime.Format("2006-01-02 03:04:05 -0700")),
		line.Layout,
		item(values.String(values.StrBestBlockAge), components.SecondsToDays(secondsSinceBestBlock)),
		line.Layout,
		item(values.String(values.StrWalletDirectory), pg.WL.WalletDirectory()),
		line.Layout,
		item(values.String(values.StrDateSize), pg.WL.DataSize()),
		line.Layout,
		item(values.String(values.StrTransactions), fmt.Sprintf("%d", len(pg.txs))),
		line.Layout,
		item(values.String(values.StrAccount)+"s", fmt.Sprintf("%d", pg.accounts.Count)),
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

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *StatPage) Layout(gtx C) D {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *StatPage) layoutDesktop(gtx C) D {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrStatistics),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: pg.layoutStats,
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}

	// Refresh frames every 1 second
	op.InvalidateOp{At: time.Now().Add(time.Second * 1)}.Add(gtx.Ops)
	return components.UniformPadding(gtx, container)
}

func (pg *StatPage) layoutMobile(gtx C) D {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrStatistics),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: pg.layoutStats,
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}

	// Refresh frames every 1 second
	op.InvalidateOp{At: time.Now().Add(time.Second * 1)}.Add(gtx.Ops)
	return components.UniformMobile(gtx, false, true, container)
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
