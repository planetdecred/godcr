package main

import (
	"fmt"

	app "gioui.org/app"

	"gioui.org/font/gofont"
)

func main() {
	gofont.Register()
	win, err := createWindow(landingPage)
	if err != nil {
		fmt.Printf("Could not initialize window: %s", err)
	}
	go func(win *window) {
		win.loop()
	}(win)

	app.Main()
}
