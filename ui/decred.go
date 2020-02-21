package ui

import (
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func decredTheme() *materialplus.Theme {
	theme := materialplus.NewTheme(decredPalette)
	if theme == nil {
		return nil
	}
	theme.Icon.Cancel = mustIcon(material.NewIcon(icons.NavigationCancel))
	theme.Icon.Check = mustIcon(material.NewIcon(icons.NavigationCheck))
	theme.Icon.Logo = mustIcon(material.NewIcon(icons.ActionAccountCircle))
	return theme
}

func mustIcon(ic *material.Icon, err error) *material.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
