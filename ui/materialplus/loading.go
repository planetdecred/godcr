package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/icons"
)

type LoadingIcon struct {
	*material.Icon
}

func (t *Theme) Loading() LoadingIcon {
	return LoadingIcon{icons.ActionCached}
}

func (loading LoadingIcon) Layout(gtx *layout.Context, size unit.Value) {
	loading.Icon.Layout(gtx, size)
}
