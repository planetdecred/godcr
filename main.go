package main

import (
	"github.com/decred/dcrd/dcrutil"
)

const appDisplayName = "GoDCR"

var defaultAppDataDir = dcrutil.AppDataDir("godcr", false)

func main() {
	LaunchUserInterface(appDisplayName, defaultAppDataDir, "testnet3")
}
