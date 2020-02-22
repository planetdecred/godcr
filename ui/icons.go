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
	IconNavigationArrowForward      *material.Icon
)

func init() {
	var err error
	IconContentAdd, err = material.NewIcon(icons.ContentAdd)
	if err != nil {
		panic(err)
	}

	IconNavigationRefresh, err = material.NewIcon(icons.NavigationRefresh)
	if err != nil {
		panic(err)
	}
	IconNavigationCheck, err = material.NewIcon(icons.NavigationCheck)
	if err != nil {
		panic(err)
	}
	IconToggleIndeterminateCheckBox, err = material.NewIcon(icons.ToggleIndeterminateCheckBox)
	if err != nil {
		panic(err)
	}
	IconNavigationClose, err = material.NewIcon(icons.NavigationClose)
	if err != nil {
		panic(err)
	}
	IconNavigationArrowForward, err = material.NewIcon(icons.NavigationArrowForward)
	if err != nil {
		panic(err)
	}
}
