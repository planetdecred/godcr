package materialplus

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
	DropDownIcon       *material.Icon
<<<<<<< HEAD
	CancelIcon         *material.Icon
=======
>>>>>>> implemented account selection modal, added default live data to receive page on load and fixed minor bugs
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
<<<<<<< HEAD

	CancelIcon, err = material.NewIcon(icons.NavigationClose)
	if err != nil {
		log.Fatal(err)
	}
=======
>>>>>>> implemented account selection modal, added default live data to receive page on load and fixed minor bugs
}
