package dexc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	_ "decred.org/dcrdex/client/asset/bch" // register bch asset
	_ "decred.org/dcrdex/client/asset/btc" // register btc asset
	_ "decred.org/dcrdex/client/asset/dcr" // register dcr asset
	_ "decred.org/dcrdex/client/asset/ltc" // register ltc asset
	"decred.org/dcrdex/client/comms"
	"decred.org/dcrdex/dex/msgjson"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
)

// Dexc represents of the Decred exchange client.
type Dexc struct {
	*core.Core
	wg  sync.WaitGroup
	ctx context.Context

	connMtx sync.RWMutex
	conns   map[string]*dexConnection
}

const (
	DefaultAsset            = "dcr"
	DefaultAssetID   uint32 = 42
	ConversionFactor        = 1e8
)

func NewDex(debugLevel string, dbPath, net string, w io.Writer) (*Dexc, error) {
	logMaker := initLogging(debugLevel, true, w)
	log = logMaker.Logger("DEXC")
	if net == "testnet3" {
		net = "testnet"
	}
	n, err := dex.NetFromString(net)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Prepare the Core.
	clientCore, err := core.New(&core.Config{
		DBPath: dbPath, // global set in config.go
		Net:    n,
		Logger: logMaker.Logger("CORE"),
		// TorProxy:     TorProxy,
		// TorIsolation: TorIsolation,
	})

	if err != nil {
		log.Errorf("error creating client core: %v", err)
		return nil, err
	}

	return &Dexc{
		Core:  clientCore,
		conns: make(map[string]*dexConnection),
	}, nil
}

func (d *Dexc) Run(appCtx context.Context, cancel context.CancelFunc) {
	d.ctx = appCtx
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.Core.Run(d.ctx)
		cancel()
	}()
	<-d.Core.Ready()
}

// DexConnection represents a connection to websocket server.
type dexConnection struct {
	comms.WsConn
	host string
}

func (d *Dexc) ConnectDexes(host string, password []byte) {
	exchanges := d.Exchanges()
	exchange, ok := exchanges[host]
	if !ok {
		log.Errorf("Host %s not found", host)
		return
	}

	dc, _ := d.connectDex(host, password)
	for _, mkt := range exchange.Markets {
		_, err := dc.subscribe(mkt.BaseID, mkt.QuoteID)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	d.connMtx.Lock()
	d.conns[host] = dc
	d.connMtx.Unlock()
}

func (d *Dexc) connectDex(host string, password []byte) (*dexConnection, error) {
	ctx := d.ctx
	acct, err := d.Core.AccountExport([]byte(password), host)
	if err != nil {
		return nil, err
	}

	cert, err := hex.DecodeString(acct.Cert)
	if err != nil {
		return nil, err
	}

	wsAddr := "wss://" + host + "/ws"
	wsURL, err := url.Parse(wsAddr)
	if err != nil {
		return nil, fmt.Errorf("error parsing ws address %s: %w", wsAddr, err)
	}

	dc := &dexConnection{}
	cfg := &comms.WsCfg{
		URL:      wsURL.String(),
		Cert:     cert,
		PingWait: 20 * time.Second, // larger than server's pingPeriod (server/comms/server.go)
		Logger:   dex.StdOutLogger("DEX_CNN", dex.LevelInfo),
	}
	wsconn, err := comms.NewWsConn(cfg)
	if err != nil {
		log.Errorf("Failure to create new socket connection %v", err)
		return nil, err
	}

	cmaster := dex.NewConnectionMaster(wsconn)
	dc.WsConn = wsconn
	dc.host = host

	err = cmaster.Connect(ctx)
	if err != nil {
		log.Errorf(">>> Connect to websocket failure %v", err)
	}

	return dc, nil
}

func (dc *dexConnection) subscribe(b, q uint32) (*msgjson.OrderBook, error) {
	mkt, _ := dex.MarketName(b, q)
	req, err := msgjson.NewRequest(dc.NextID(), msgjson.OrderBookRoute, &msgjson.OrderBookSubscription{
		Base:  b,
		Quote: q,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding 'orderbook' request: %w", err)
	}
	errChan := make(chan error, 1)
	result := new(msgjson.OrderBook)
	err = dc.RequestWithTimeout(req, func(msg *msgjson.Message) {
		errChan <- msg.UnmarshalResult(result)
	}, comms.DefaultResponseTimeout, func() {
		errChan <- fmt.Errorf("timed out waiting for '%s' response", msgjson.OrderBookRoute)
	})
	if err != nil {
		return nil, fmt.Errorf("error subscribing to %s orderbook: %w", mkt, err)
	}
	err = <-errChan
	if err != nil {
		return nil, fmt.Errorf("error subscribing to %s orderbook: %w", mkt, err)
	}
	return result, nil
}

func (d *Dexc) MessageSource(host string) <-chan *msgjson.Message {
	dc, ok := d.conns[host]
	if !ok || dc == nil {
		log.Errorf("Connection: %s not exist ", host)
		return nil
	}
	return dc.MessageSource()
}
