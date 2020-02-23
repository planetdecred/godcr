package layouts

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/styles"
)

func FlexedWithStyle(gtx *layout.Context, style styles.Style, weight float32, widget func()) layout.FlexChild {
	return layout.Flexed(weight, styles.WithStyle(gtx, style, widget))
}

func RigidWithStyle(gtx *layout.Context, style styles.Style, widget func()) layout.FlexChild {
	return layout.Rigid(styles.WithStyle(gtx, style, widget))
}
