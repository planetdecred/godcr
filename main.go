package main

import (
	"github.com/decred/dcrd/dcrutil"
)

const appDisplayName = "GoDCR"

var defaultAppDataDir = dcrutil.AppDataDir("godcr", false)

func main() {
	launchUserInterface(appDisplayName, defaultAppDataDir, "testnet3")
}
