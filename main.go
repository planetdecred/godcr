package main

import (
	"github.com/decred/dcrd/dcrutil"
)

const (
	appDisplayName = "GoDCR"
	netType        = "testnet3"
)

var defaultAppDataDir = dcrutil.AppDataDir("godcr", false)

func main() {
	launchUserInterface(appDisplayName, defaultAppDataDir, netType)
}
