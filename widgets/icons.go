package widgets

import (
	"log"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	AddIcon                 *Icon
	ReturnIcon              *Icon
	NavigationCheckIcon     *Icon
	NavigationArrowBackIcon *Icon
	CancelIcon              *Icon
)

func init() {
	var err error

	AddIcon, err = NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}

	ReturnIcon, err = NewIcon(icons.ContentReply)
	if err != nil {
		log.Fatal(err)
	}

	NavigationCheckIcon, err = NewIcon(icons.NavigationCheck)
	if err != nil {
		log.Fatal(err)
	}

	NavigationArrowBackIcon, err = NewIcon(icons.NavigationArrowBack)
	if err != nil {
		log.Fatal(err)
	}

	CancelIcon, err = NewIcon(icons.NavigationCancel)
	if err != nil {
		log.Fatal(err)
	}
}
