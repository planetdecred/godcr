package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	_ "github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/wallet"
)

var (
	Version   string = "1.0-beta1"
	BuildDate string
	BuildEnv  = wallet.DevBuild
)

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

	// var buildDate time.Time
	// if BuildEnv == wallet.ProdBuild {
	// 	buildDate, err = time.Parse(time.RFC3339, BuildDate)
	// 	if err != nil {
	// 		fmt.Printf("Error: %s\n", err.Error())
	// 		return
	// 	}
	// } else {
	// 	buildDate = time.Now()
	// }

	// logFilePath := filepath.Join(cfg.LogDir, defaultLogFilename)
	appInstance, err := app.Init(cfg.HomeDir, cfg.Network, Version)
	if err != nil {
		log.Error(err)
		return
	}

	// Start the GUI frontend.
	appInstance.Run(page.NewStartPage(appInstance))
}
