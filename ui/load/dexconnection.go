package load

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"sync"
	"time"

	"decred.org/dcrdex/client/comms"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/msgjson"
)

// DexConnection represents a connection to websocket server.
type dexConnection struct {
	comms.WsConn
	host string
}

type DexConnections struct {
	connMtx sync.RWMutex
	conns   map[string]*dexConnection
}

// CreateDexesConnection after login or registration, this will initialize websocket connection
// use to listen the messages, if add new dex call init again.
func (d *Load) CreateDexConnection(host string, password []byte) {
	exchanges := d.Dexc().DEXServers()
	_, ok := exchanges[host]
	if !ok {
		log.Errorf("Host %s not found", host)
		return
	}

	acct, err := d.Dexc().Core().AccountExport(password, host)
	if err != nil {
		return
	}

	cert, err := hex.DecodeString(acct.Cert)
	if err != nil {
		return
	}

	wsAddr := "wss://" + host + "/ws"
	wsURL, err := url.Parse(wsAddr)
	if err != nil {
		log.Errorf("[ERROR] " + wsAddr + " " + err.Error())
		return
	}

	cfg := &comms.WsCfg{
		URL:      wsURL.String(),
		Cert:     cert,
		PingWait: 20 * time.Second, // larger than server's pingPeriod (server/comms/server.go)
		Logger:   dex.StdOutLogger("DEX_CNN", dex.LevelInfo),
	}
	wsconn, err := comms.NewWsConn(cfg)
	if err != nil {
		log.Errorf("Failure to create new socket connection %v", err)
		return
	}

	dc := &dexConnection{
		WsConn: wsconn,
		host:   host,
	}

	d.Dexcnn.connMtx.Lock()
	d.Dexcnn.conns[host] = dc
	d.Dexcnn.connMtx.Unlock()
}

// SubscribeMarket create connection to websocket server and listen for messages.
func (d *Load) SubscribeMarket(host string, baseID, quoteID uint32) {
	ctx := context.TODO()
	exchanges := d.Dexc().DEXServers()
	_, ok := exchanges[host]
	if !ok {
		log.Errorf("Host %s not found", host)
		return
	}

	d.Dexcnn.connMtx.Lock()
	dc := d.Dexcnn.conns[host]
	d.Dexcnn.connMtx.Unlock()

	err := dex.NewConnectionMaster(dc.WsConn).Connect(ctx)
	if err != nil {
		log.Errorf(">>> Connect to websocket failure %v", err)
	}

	_, err = dc.subscribe(baseID, quoteID)
	if err != nil {
		log.Error(err)
	}
}

func (d *Load) MessageSource(host string) <-chan *msgjson.Message {
	dc, ok := d.Dexcnn.conns[host]
	if !ok || dc == nil {
		log.Errorf("Connection: %s not exist ", host)
		return nil
	}

	return dc.MessageSource()
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
