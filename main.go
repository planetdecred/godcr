package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/text"

	_ "net/http/pprof"

	"gioui.org/app"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui"
	"github.com/planetdecred/godcr/wallet"
)

func getAbsoultePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %s", err.Error())
	}

	exSym, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("error getting filepath after evaluating sym links")
	}

	return path.Dir(exSym), nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if cfg.Profile > 0 {
		go func() {
			log.Info(fmt.Sprintf("Starting profiling server on port %d", cfg.Profile))
			log.Error(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", cfg.Profile), nil))
		}()
	}

	dcrlibwallet.SetLogLevels(cfg.DebugLevel)

	absoluteWdPath, err := getAbsoultePath()
	if err != nil {
		panic(err)
	}

	decredIcons := make(map[string]image.Image)
	err = filepath.Walk(filepath.Join(absoluteWdPath, "ui/assets/decredicons"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := os.Open(path)
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		split := strings.Split(info.Name(), ".")
		decredIcons[split[0]] = img
		return nil
	})
	if err != nil {
		log.Warn(err)
	}

	var confirms int32 = dcrlibwallet.DefaultRequiredConfirmations

	if cfg.SpendUnconfirmed {
		confirms = 0
	}

	wal, err := wallet.NewWallet(cfg.HomeDir, cfg.Network, make(chan wallet.Response, 3), confirms)
	if err != nil {
		log.Error(err)
		return
	}

	wal.LoadWallets()
	shutdown := make(chan int)
	go func() {
		<-shutdown
		wal.Shutdown()
		os.Exit(0)
	}()

	var collection []text.FontFace
	source, err := os.Open(filepath.Join(absoluteWdPath, "ui/assets/fonts/source_sans_pro_regular.otf"))
	if err != nil {
		fmt.Println("Failed to load font")
		collection = gofont.Collection()
	} else {
		stat, err := source.Stat()
		if err != nil {
			fmt.Println(err)
		}
		bytes := make([]byte, stat.Size())
		source.Read(bytes)
		fnt, err := opentype.Parse(bytes)
		if err != nil {
			fmt.Println(err)
		}
		collection = append(collection, text.FontFace{Font: text.Font{}, Face: fnt})
	}

	win, err := ui.CreateWindow(wal, decredIcons, collection, internalLog)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}

	// Start the ui frontend
	go win.Loop(shutdown)
	app.Main()
}
