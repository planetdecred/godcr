package assets

import (
	"image/color"
	"log"

	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	NavigationMoreIcon *material.Icon
	ContentCopyIcon    *material.Icon
	ActionInfoIcon     *material.Icon
	DropDownIcon              *material.Icon
)

func init() {
	var err error

	NavigationMoreIcon, err = material.NewIcon(icons.NavigationMoreVert)
	if err != nil {
		log.Fatal(err)
	}

	ContentCopyIcon, err = material.NewIcon(icons.ContentContentCopy)
	if err != nil {
		log.Fatal(err)
	}

	ActionInfoIcon, err = material.NewIcon(icons.ActionInfo)
	if err != nil {
		log.Fatal(err)
	}
	ActionInfoIcon.Color = color.RGBA{44, 114, 255, 255}

	DropDownIcon, err = material.NewIcon(icons.NavigationArrowDropDown)
	if err != nil {
		log.Fatal(err)
	}
}
