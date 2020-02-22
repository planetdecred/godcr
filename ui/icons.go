package ui

import (
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	IconContentAdd                  *material.Icon
	IconNavigationRefresh           *material.Icon
	IconNavigationCheck             *material.Icon
	IconToggleIndeterminateCheckBox *material.Icon
	IconNavigationClose             *material.Icon
)

func init() {
	var err error
	IconContentAdd, err = material.NewIcon(icons.ContentAdd)

	IconNavigationRefresh, err = material.NewIcon(icons.NavigationRefresh)
	IconNavigationCheck, err = material.NewIcon(icons.NavigationCheck)
	IconToggleIndeterminateCheckBox, err = material.NewIcon(icons.ToggleIndeterminateCheckBox)
	IconNavigationClose, err = material.NewIcon(icons.NavigationClose)
	if err != nil {
		panic(err)
	}
}
