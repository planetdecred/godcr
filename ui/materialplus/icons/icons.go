package icons

import (
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	ContentAdd                  = mustIcon(material.NewIcon(icons.ContentAdd))
	NavigationRefresh           *material.Icon
	NavigationCheck             *material.Icon
	ToggleIndeterminateCheckBox *material.Icon
	NavigationClose             *material.Icon
	NavigationArrowForward      *material.Icon
	ActionCached                *material.Icon
)

func init() {
	var err error
	ContentAdd, err = material.NewIcon(icons.ContentAdd)
	if err != nil {
		panic(err)
	}

	NavigationRefresh, err = material.NewIcon(icons.NavigationRefresh)
	if err != nil {
		panic(err)
	}
	NavigationCheck, err = material.NewIcon(icons.NavigationCheck)
	if err != nil {
		panic(err)
	}
	ToggleIndeterminateCheckBox, err = material.NewIcon(icons.ToggleIndeterminateCheckBox)
	if err != nil {
		panic(err)
	}
	NavigationClose, err = material.NewIcon(icons.NavigationClose)
	if err != nil {
		panic(err)
	}
	NavigationArrowForward, err = material.NewIcon(icons.NavigationArrowForward)
	if err != nil {
		panic(err)
	}

	ActionCached, err = material.NewIcon(icons.ActionCached)
	if err != nil {
		panic(err)
	}
}

func mustIcon(ic *material.Icon, err error) *material.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
