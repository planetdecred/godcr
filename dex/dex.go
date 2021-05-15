package dex

import (
	"context"
	"fmt"
	"io"
	"sync"

	// _ "decred.org/dcrdex/client/asset/bch" // register btc asset
	_ "decred.org/dcrdex/client/asset/btc" // register btc asset
	_ "decred.org/dcrdex/client/asset/dcr" // register dcr asset
	_ "decred.org/dcrdex/client/asset/ltc" // register ltc asset

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
)

type Dex struct {
	core *core.Core
}

const DefaultAssert = 42

func NewDex(debugLevel string, dbPath, net string, w io.Writer) (*Dex, error) {
	logMaker := initLogging(debugLevel, true, w)
	log = logMaker.Logger("DEXC")

	// Prepare the Core.
	clientCore, err := core.New(&core.Config{
		DBPath: dbPath, // global set in config.go
		Net:    dex.Network(1),
		Logger: logMaker.Logger("CORE"),
		// TorProxy:     TorProxy,
		// TorIsolation: TorIsolation,
	})

	if err != nil {
		fmt.Printf("error creating client core: %s", err)
		return nil, err
	}

	return &Dex{clientCore}, nil
}

func (d *Dex) Run(appCtx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.core.Run(appCtx)
		cancel() // in the event that Run returns prematurely prior to context cancellation
	}()
	<-d.core.Ready()
}
