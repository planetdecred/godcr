// Copyright (c) 2017, The dcrdata developers
// See LICENSE for details.

package dex

import (
	"fmt"
	"io"
	"os"

	"decred.org/dcrdex/dex"
	"github.com/decred/slog"
)

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var apiLog = slog.Disabled
var log = slog.Disabled

// DisableLog disables all library log output.  Logging output is disabled
// by default until UseLogger is called.
func DisableLog() {
	apiLog = slog.Disabled
	log = slog.Disabled
}

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger slog.Logger) {
	apiLog = logger
	log = logger
}

// initLogging initializes the logging rotater to write logs to logFile and
// create roll files in the same directory. initLogging must be called before
// the package-global log rotator variables are used.
func initLogging(lvl string, utc bool, w io.Writer) *dex.LoggerMaker {
	lm, err := dex.NewLoggerMaker(w, lvl, utc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create custom logger: %v\n", err)
		os.Exit(1)
	}
	log = lm.Logger("APP")
	return lm
}
