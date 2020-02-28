package window

import (
	"gioui.org/layout"

	"github.com/raedahgroup/godcr-gio/ui/page"
)

const (
	navWidth = 200
)

func (win *Window) layoutPage(gtx *layout.Context, handler page.Handler) interface{} {
	if !handler.IsNavPage {
		return handler.Page.Draw(gtx)
	}

	var evt interface{}

	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			gtx.Constraints.Width.Min = navWidth
			win.layoutNavSection(gtx)
		}),
		layout.Rigid(func() {
			evt = handler.Page.Draw(gtx)
		}),
	)

	return evt
}

func (win *Window) layoutNavSection(gtx *layout.Context) {

}
