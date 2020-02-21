package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func ListInsetedWidgets(gtx *layout.Context, list layout.List, inset layout.Inset, widgets []func()) {
	list.Layout(gtx, len(widgets), func(i int) {
		inset.Layout(gtx, widgets[i])
	})
}

func ListUniformWidgets(gtx *layout.Context, list layout.List, spacing unit.Value, widgets []func()) {
	ListInsetedWidgets(gtx, layout.List{Axis: layout.Vertical}, layout.UniformInset(spacing), widgets)
}

func VerticalInsetedList(gtx *layout.Context, inset layout.Inset, widgets []func()) {
	ListInsetedWidgets(gtx, layout.List{Axis: layout.Vertical}, inset, widgets)
}

func VerticalUniformList(gtx *layout.Context, spacing unit.Value, widgets []func()) {
	ListUniformWidgets(gtx, layout.List{Axis: layout.Vertical}, spacing, widgets)
}
