package main

import ( // app "gioui.org/app"
	// "gioui.org/font/gofont"
	"fmt"

	"github.com/raedahgroup/godcr-gio/wallet"
)

func main() {
	// gofont.Register()
	// win, err := createWindow(landingPage)
	// if err != nil {
	// 	fmt.Printf("Could not initialize window: %s\ns", err)
	// 	return
	// }
	// go func(win *window) {
	// 	win.loop()
	// }(win)

	// app.Main()

	cfg, err := loadConfig()

	if err != nil {
		return
	}
	_, atleastone, err := wallet.LoadWallets(cfg.HomeDir, cfg.Network)
	if err != nil {
		fmt.Println(err)
		return
	}
	if atleastone {
		fmt.Println("Wallet loaded")
	} else {
		fmt.Println("No Wallets")
	}
}
