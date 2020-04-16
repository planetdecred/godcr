package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type Tab struct {
	Button  *widget.Button
	Content func() layout.Widget
	Icon    *paint.ImageOp
	TabName string
}

func NewTab(text string, content func() layout.Widget) Tab {
	return Tab{
		Content: content,
		Button:  new(widget.Button),
		TabName: text,
	}
}

func NewTabWithIcon(text string, icon paint.ImageOp, content func() layout.Widget) Tab {
	return Tab{
		Content: content,
		Button:  new(widget.Button),
		TabName: text,
		Icon:    &icon,
	}
}
