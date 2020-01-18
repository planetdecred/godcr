package main

import (
	"os"
	"path/filepath"

	"github.com/decred/dcrd/dcrutil/v2"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultNetwork        = "testnet3"
	defaultConfigFileName = "godcr.conf"
)

var (
	defaultHomeDir    = dcrutil.AppDataDir("godcr", false)
	defaultConfigFile = filepath.Join(defaultHomeDir, defaultConfigFileName)
)

type config struct {
	Network    string `long:"network" description:"Network to use"`
	HomeDir    string `long:"appdata" description:"Directory where the app configuration file and wallet data is stored"`
	ConfigFile string `long:"configfile" description:"Filename of the config file in the app directory"`
}

var defaultConfig = config{
	Network:    defaultNetwork,
	HomeDir:    defaultHomeDir,
	ConfigFile: defaultConfigFile,
}

// loadConfig loads the configration file stored in the defaultHomeDir
func loadConfig() (*config, error) {
	cfg := defaultConfig
	_ = flags.NewParser(&cfg, flags.Default)

	var configFile *os.File
	// if the config file does not exist, create it
	if _, err := os.Stat(cfg.ConfigFile); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.HomeDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
		configFile, err = os.Create(cfg.ConfigFile)
		defer configFile.Close()
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
