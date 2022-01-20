package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/decred/slog"
	flags "github.com/jessevdk/go-flags"
	"github.com/planetdecred/godcr/version"
)

const (
	defaultNetwork        = "mainnet"
	defaultConfigFileName = "godcr.conf"
	defaultLogFilename    = "godcr.log"
	defaultLogLevel       = "info"
	defaultLogDirname     = "logs"
)

var (
	defaultHomeDir        = dcrutil.AppDataDir("godcr", false)
	defaultConfigFilename = filepath.Join(defaultHomeDir, defaultConfigFileName)
	defaultLogDir         = filepath.Join(defaultHomeDir, defaultLogDirname)
)

type config struct {
	Network          string `long:"network" description:"Network to use"`
	HomeDir          string `long:"appdata" description:"Directory where the app configuration file and wallet data is stored"`
	ConfigFile       string `long:"configfile" description:"Filename of the config file in the app directory"`
	ShowVersion      bool   `short:"V" long:"version" description:"Display version information and exit"`
	MaxLogZips       int    `long:"max-log-zips" description:"The number of zipped log files created by the log rotator to be retained. Setting to 0 will keep all."`
	LogDir           string `long:"logdir" description:"Directory to log output."`
	DebugLevel       string `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	Quiet            bool   `short:"q" long:"quiet" description:"Easy way to set debuglevel to error"`
	SpendUnconfirmed bool   `long:"spendunconfirmed" description:"Allow the multiwallet to use transactions that have not been confirmed"`
	Profile          int    `long:"profile" description:"Runs local web server for profiling"`
}

var defaultConfig = config{
	Network:    defaultNetwork,
	HomeDir:    defaultHomeDir,
	ConfigFile: defaultConfigFilename,
	LogDir:     defaultLogDir,
	DebugLevel: defaultLogLevel,
}

// validLogLevel returns whether or not logLevel is a valid debug log level.
func validLogLevel(logLevel string) bool {
	_, ok := slog.LevelFromString(logLevel)
	return ok
}

// supportedSubsystems returns a sorted slice of the supported subsystems for
// logging purposes.
func supportedSubsystems() []string {
	// Convert the subsystemLoggers map keys to a slice.
	subsystems := make([]string, 0, len(subsystemLoggers))
	for subsysID := range subsystemLoggers {
		subsystems = append(subsystems, subsysID)
	}

	// Sort the subsystems for stable display.
	sort.Strings(subsystems)
	return subsystems
}

// parseAndSetDebugLevels attempts to parse the specified debug level and set
// the levels accordingly.  An appropriate error is returned if anything is
// invalid.
func parseAndSetDebugLevels(debugLevel string) error {
	// When the specified string doesn't have any delimiters, treat it as
	// the log level for all subsystems.
	if !strings.Contains(debugLevel, ",") && !strings.Contains(debugLevel, "=") {
		// Validate debug log level.
		if !validLogLevel(debugLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, debugLevel)
		}

		// Change the logging level for all subsystems.
		setLogLevels(debugLevel)

		return nil
	}

	// Split the specified string into subsystem/level pairs while detecting
	// issues and update the log levels accordingly.
	for _, logLevelPair := range strings.Split(debugLevel, ",") {
		if !strings.Contains(logLevelPair, "=") {
			str := "The specified debug level contains an invalid " +
				"subsystem/level pair [%v]"
			return fmt.Errorf(str, logLevelPair)
		}

		// Extract the specified subsystem and log level.
		fields := strings.Split(logLevelPair, "=")
		subsysID, logLevel := fields[0], fields[1]

		// Validate subsystem.
		if _, exists := subsystemLoggers[subsysID]; !exists {
			str := "The specified subsystem [%v] is invalid -- " +
				"supported subsystems %v"
			return fmt.Errorf(str, subsysID, supportedSubsystems())
		}

		// Validate log level.
		if !validLogLevel(logLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, logLevel)
		}

		setLogLevel(subsysID, logLevel)
	}

	return nil
}

// loadConfig initializes and parses the config using a config file and command
// line options.
func loadConfig() (*config, error) {
	loadConfigError := func(err error) (*config, error) {
		return nil, err
	}

	// // Default config
	cfg := defaultConfig
	defaultConfigNow := defaultConfig

	// Pre-parse the command line options to see if an alternative config file
	// or the version flag was specified. Override any environment variables
	// with parsed command line flags.
	preParser := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash)
	_, err := preParser.Parse()

	if err != nil {
		e, ok := err.(*flags.Error)
		if !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		if ok && e.Type == flags.ErrHelp {
			preParser.WriteHelp(os.Stdout)
			os.Exit(0)
		}
		return loadConfigError(err)
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	if cfg.ShowVersion {
		fmt.Printf("%s version %s (Go version %s)\n", appName,
			version.Version(), runtime.Version())
		os.Exit(0)
	}

	// If a non-default appdata folder is specified on the command line, it may
	// be necessary adjust the config file location. If the the config file
	// location was not specified on the command line, the default location
	// should be under the non-default appdata directory. However, if the config
	// file was specified on the command line, it should be used regardless of
	// the appdata directory.
	if defaultHomeDir != cfg.HomeDir && defaultConfigNow.ConfigFile == cfg.ConfigFile {
		cfg.ConfigFile = filepath.Join(cfg.HomeDir, defaultConfigFilename)
		// Update the defaultConfig to avoid an error if the config file in this
		// "new default" location does not exist.
		defaultConfigNow.ConfigFile = cfg.ConfigFile
	}

	// Load additional config from file.
	var configFileError error
	// Config file name for logging.
	configFile := "NONE (defaults)"
	parser := flags.NewParser(&cfg, flags.Default)

	// Do not error default config file is missing.
	if _, err := os.Stat(cfg.ConfigFile); os.IsNotExist(err) {
		// Non-default config file must exist
		if defaultConfigNow.ConfigFile != cfg.ConfigFile {
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		// Warn about missing default config file, but continue
		fmt.Printf("Config file (%s) does not exist. Using defaults.\n",
			cfg.ConfigFile)
	} else {
		// The config file exists, so attempt to parse it.
		err = flags.NewIniParser(parser).ParseFile(cfg.ConfigFile)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				fmt.Fprintln(os.Stderr, err)
				parser.WriteHelp(os.Stderr)
				return loadConfigError(err)
			}
			configFileError = err
		}
		configFile = cfg.ConfigFile
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return loadConfigError(err)
	}

	// Create the home directory if it doesn't already exist.
	funcName := "loadConfig"
	err = os.MkdirAll(cfg.HomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is linked to a
		// directory that does not exist (probably because it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "%s: failed to create home directory: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// If a non-default appdata folder is specified, it may be necessary to
	// adjust the LogDir.
	if defaultHomeDir != cfg.HomeDir {
		if defaultLogDir == cfg.LogDir {
			cfg.LogDir = filepath.Join(cfg.HomeDir, defaultLogDirname)
		}
	}

	// Warn about missing config file after the final command line parse
	// succeeds.  This prevents the warning on help messages and invalid
	// options.
	if configFileError != nil {
		fmt.Printf("%v\n", configFileError)
		return loadConfigError(configFileError)
	}

	logRotator = nil
	cfg.LogDir = cleanAndExpandPath(cfg.LogDir)

	// Initialize log rotation. After log rotation has been initialized, the
	// logger variables may be used. This creates the LogDir if needed.
	if cfg.MaxLogZips < 0 {
		cfg.MaxLogZips = 0
	}
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename), cfg.MaxLogZips)

	// Special show command to list supported subsystems and exit.
	if cfg.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Parse, validate, and set debug log level(s).
	if cfg.Quiet {
		cfg.DebugLevel = "error"
	}

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		err = fmt.Errorf("%s: %v", funcName, err.Error())
		fmt.Fprintln(os.Stderr, err)
		parser.WriteHelp(os.Stderr)
		return loadConfigError(err)
	}

	log.Debugf("Log folder: %s", cfg.LogDir)
	log.Debugf("Config file: %s", configFile)

	return &cfg, nil
}

// cleanAndExpandPath expands environment variables and leading ~ in the passed
// path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but the variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser to
	// otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}
