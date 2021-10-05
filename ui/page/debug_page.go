package page

import (
	"os"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DebugPageID = "Debug"

type debugItem struct {
	text   string
	page   string
	action func()
}

type DebugPage struct {
	*load.Load
	debugItems []debugItem
	list       *decredmaterial.ClickableList

	backButton decredmaterial.IconButton
}

func NewDebugPage(l *load.Load) *DebugPage {
	debugItems := []debugItem{
		{
			text: "Check wallet logs",
			page: LogPageID,
			action: func() {
				l.ChangeFragment(NewLogPage(l))
			},
		},
		{
			text: "Check statistics",
			page: StatisticsPageID,
			action: func() {
				l.ChangeFragment(NewStatPage(l))
			},
		},
	}

	pg := &DebugPage{
		Load:       l,
		debugItems: debugItems,
		list:       l.Theme.NewClickableList(layout.Vertical),
	}
	pg.list.Radius = decredmaterial.Radius(14)

	// Add a "Reset DEX Client" option.
	pg.debugItems = append(pg.debugItems, debugItem{
		text: "Reset DEX Client",
		action: func() {
			pg.resetDexData()
		},
	})

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *DebugPage) ID() string {
	return DebugPageID
}

func (pg *DebugPage) OnResume() {

}

func (pg *DebugPage) Handle() {
	if clicked, item := pg.list.ItemClicked(); clicked {
		pg.debugItems[item].action()
	}
}

func (pg *DebugPage) OnClose() {}

func (pg *DebugPage) debugItem(gtx C, i int) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.Theme.Body1(pg.debugItems[i].text).Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Px(values.MarginPadding22)
					return pg.Icons.ChevronRight.Layout(gtx, pg.Theme.Color.Gray)
				})
			})
		}),
	)
}

func (pg *DebugPage) layoutDebugItems(gtx C) {
	background := pg.Theme.Color.Surface
	card := pg.Theme.Card()
	card.Color = background
	card.Layout(gtx, func(gtx C) D {
		return pg.list.Layout(gtx, len(pg.debugItems), func(gtx C, i int) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.debugItem(gtx, i)
				}),
				layout.Rigid(func(gtx C) D {
					if i == len(pg.debugItems)-1 {
						return layout.Dimensions{}
					}
					return layout.Inset{
						Left: values.MarginPadding16,
					}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
			)
		})
	})
}

func (pg *DebugPage) Layout(gtx C) D {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Debug",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				pg.layoutDebugItems(gtx)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			},
		}
		return sp.Layout(gtx)

	}
	return components.UniformPadding(gtx, container)
}

func (pg *DebugPage) resetDexData() {
	// Show confirm modal and delete dexc db and dex-related settings in the
	// multiwallet db.
	confirmModal := modal.NewInfoModal(pg.Load).
		Title("Confirm DEX Client Reset").
		Body("You'll need to restart godcr before you can use the DEX again. Proceed?").
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Reset DEX Client", func() {
			// Attempt to first shutdown the dex client. This will fail if there
			// are active orders.
			// TODO: Since this is a debug feature, consider allowing dex shutdown
			// even if there are active orders.
			if pg.DL.Dexc.Shutdown() {
				// Dexc shutdown was successful, perform other cleanup here
				// including deleting the dexc db.
				pg.WL.MultiWallet.DeleteUserConfigValueForKey(dexc.ConnectedDcrWalletIDConfigKey)
				err := os.RemoveAll(pg.DL.Dexc.DbPath)
				if err != nil {
					log.Warnf("DEX client data reset but failed to delete DEX db: %v", err)
				}
				pg.Toast.Notify("DEX client data reset complete.")
			} else {
				pg.Toast.NotifyError("Cannot reset DEX client data because the DEX client could not be shut down. Check the logs.")
			}
		})
	pg.ShowModal(confirmModal)
}
