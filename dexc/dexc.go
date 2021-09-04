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
	"decred.org/dcrdex/dex/order"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
)

type Dexc struct {
	core *core.Core
	Send chan Response
	Net  string

	connMtx sync.RWMutex
	conns   map[string]*dexConnection

	wg  sync.WaitGroup
	ctx context.Context
}

const (
	DefaultAsset            = "dcr"
	DefaultAssetID   uint32 = 42
	conversionFactor        = 1e8
)

func NewDex(debugLevel string, dbPath, net string, send chan Response, w io.Writer) (*Dexc, error) {
	logMaker := initLogging(debugLevel, true, w)
	log = logMaker.Logger("DEXC")

	_, err := dex.NetFromString(net)
	if err != nil {
		log.Error(err)
		// return nil, err
	}

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

	return &Dexc{
		core:  clientCore,
		Send:  send,
		Net:   net,
		conns: make(map[string]*dexConnection),
	}, nil
}

func (d *Dexc) Run(appCtx context.Context, cancel context.CancelFunc) {
	d.ctx = appCtx

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.core.Run(d.ctx)
		cancel()
	}()
	<-d.core.Ready()

	d.GetUser()

	host := "127.0.0.1:7232"
	dc, _ := d.connectDex(host, "123")

	d.connMtx.Lock()
	d.conns[host] = dc
	d.connMtx.Unlock()
}

// DexConnection represents a connection to websocket server
type dexConnection struct {
	comms.WsConn
	host string
}

func (d *Dexc) connectDex(host string, password string) (*dexConnection, error) {
	ctx := d.ctx
	acct, err := d.core.AccountExport([]byte(password), host)
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

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		go d.listen(dc)
	}()

	err = cmaster.Connect(ctx)
	if err != nil {
		log.Errorf(">>>>>>>>>>>>>>>>>>>>> Connect to websocket failer %v", err)
	}

	{
		// TODO: will create function for subscribing

		req, err := msgjson.NewRequest(wsconn.NextID(), msgjson.OrderBookRoute, &msgjson.OrderBookSubscription{
			Base:  42,
			Quote: 0,
		})
		if err != nil {
			log.Errorf("error encoding 'orderbook' request: %w", err)
		}
		errChan := make(chan error, 1)
		result := new(msgjson.OrderBook)
		err = wsconn.RequestWithTimeout(req, func(msg *msgjson.Message) {
			errChan <- msg.UnmarshalResult(result)
		}, comms.DefaultResponseTimeout, func() {
			errChan <- fmt.Errorf("timed out waiting for '%s' response", msgjson.OrderBookRoute)
		})
		if err != nil {
			log.Errorf("error subscribing to %s orderbook: %w", "test mkt", err)
		}
		err = <-errChan
		if err != nil {
			log.Errorf("error subscribing to %s orderbook: %w", err)
		}
	}

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.readNotifications(ctx)
	}()

	return dc, nil
}

// listen monitors the DEX websocket connection for server requests and
// notifications. This should be run as a goroutine. listen will return when
// either c.ctx is canceled or the Message channel from the dexConnection's
// MessageSource method is closed. The latter would be the case when the
// dexConnection's WsConn is shut down / ConnectionMaster stopped.
func (d *Dexc) listen(dc *dexConnection) {
	msgs := (*dc).MessageSource()

out:
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Errorf("listen(wc): Connection terminated for %wsc.", "test")
				return
			}

			var handler routeHandler
			var found bool
			switch msg.Type {
			case msgjson.Request:
				// log.Infof(">>>>>>>>>>>>>>>>>>>>> Receive message source Request: %s ", msg)
				handler, found = reqHandlers[msg.Route]
			case msgjson.Notification:
				// log.Infof(">>>>>>>>>>>>>>>>>>>>> Receive message source: Notification %s ", msg)
				handler, found = noteHandlers[msg.Route]
			case msgjson.Response:
				// client/comms.wsConn handles responses to requests we sent.
				log.Errorf("A response was received in the message queue: %s", msg)
				continue
			default:
				log.Errorf("Invalid message type %d from MessageSource", msg.Type)
				continue
			}
			// Until all the routes have handlers, check for nil too.
			if !found || handler == nil {
				log.Errorf("No handler found for route '%s'", msg.Route)
				continue
			}

			// handling of this message.
			handler(d, dc, msg)
		case <-d.ctx.Done():
			break out
		}
	}
}

// routeHandler is a handler for a message from the DEX.
type routeHandler func(*Dexc, *dexConnection, *msgjson.Message) error

var reqHandlers = map[string]routeHandler{
	msgjson.PreimageRoute: handlePreimageRequest,
	// msgjson.MatchRoute:      handleMatchRoute,
	// msgjson.AuditRoute:      handleAuditRoute,
	// msgjson.RedemptionRoute: handleRedemptionRoute, // TODO: to ntfn
}

var noteHandlers = map[string]routeHandler{
	// msgjson.MatchProofRoute:      handleMatchProofMsg,
	msgjson.BookOrderRoute:  handleBookOrderMsg,
	msgjson.EpochOrderRoute: handleEpochOrderMsg,
	// msgjson.UnbookOrderRoute:     handleUnbookOrderMsg,
	// msgjson.UpdateRemainingRoute: handleUpdateRemainingMsg,
	// msgjson.EpochReportRoute:     handleEpochReportMsg,
	// msgjson.SuspensionRoute:      handleTradeSuspensionMsg,
	// msgjson.ResumptionRoute:      handleTradeResumptionMsg,
	// msgjson.NotifyRoute:          handleNotifyMsg,
	// msgjson.PenaltyRoute:         handlePenaltyMsg,
	// msgjson.NoMatchRoute:         handleNoMatchRoute,
	// msgjson.RevokeOrderRoute:     handleRevokeOrderMsg,
	// msgjson.RevokeMatchRoute:     handleRevokeMatchMsg,
}

// readNotifications reads from the Core notification channel
func (d *Dexc) readNotifications(ctx context.Context) {
	ch := d.core.NotificationFeed()
	for {
		select {
		case n := <-ch:
			log.Info("Recv notification", n)
		case <-ctx.Done():
			return
		}
	}
}

// handleEpochOrderMsg is called when an epoch_order notification is
// received.
func handleEpochOrderMsg(d *Dexc, dc *dexConnection, msg *msgjson.Message) error {
	note := new(msgjson.EpochOrderNote)
	err := msg.Unmarshal(note)
	if err != nil {
		return fmt.Errorf("epoch order note unmarshal error: %w", err)
	}

	var resp Response
	resp.Resp = BookUpdate{
		Action:   msg.Route,
		Host:     dc.host,
		MarketID: note.MarketID,
		Payload:  minifyOrder(note.OrderID, &note.TradeNote, note.Epoch),
	}

	d.Send <- resp

	return nil
}

// handleBookOrderMsg is called when a book_order notification is received.
func handleBookOrderMsg(_ *Dexc, dc *dexConnection, msg *msgjson.Message) error {
	note := new(msgjson.BookOrderNote)
	err := msg.Unmarshal(note)
	if err != nil {
		return fmt.Errorf("book order note unmarshal error: %w", err)
	}

	ord := minifyOrder(note.OrderID, &note.TradeNote, 0)
	log.Info(">>>>>> handleBookOrderMsg...", ord)

	return nil
}

// handlePreimageRequest handles a DEX-originating request for an order
// preimage. If the order id in the request is not known, it may launch a
// goroutine to wait for a market/limit/cancel request to finish processing.
func handlePreimageRequest(c *Dexc, dc *dexConnection, msg *msgjson.Message) error {
	req := new(msgjson.PreimageRequest)
	err := msg.Unmarshal(req)
	if err != nil {
		return fmt.Errorf("preimage request parsing error: %w", err)
	}

	if len(req.OrderID) != order.OrderIDSize {
		return fmt.Errorf("invalid order ID in preimage request")
	}

	var oid order.OrderID
	copy(oid[:], req.OrderID)

	// TODO: implement logic here
	return nil
}

// minifyOrder creates a MiniOrder from a TradeNote. The epoch and order ID must
// be supplied.
func minifyOrder(oid dex.Bytes, trade *msgjson.TradeNote, epoch uint64) *core.MiniOrder {
	return &core.MiniOrder{
		Qty:   float64(trade.Quantity) / conversionFactor,
		Rate:  float64(trade.Rate) / conversionFactor,
		Sell:  trade.Side == msgjson.SellOrderNum,
		Token: token(oid),
		Epoch: epoch,
	}
}

// token is a short representation of a byte-slice-like ID, such as a match ID
// or an order ID. The token is meant for display where the 64-character
// hexadecimal IDs are untenable.
func token(id []byte) string {
	if len(id) < 4 {
		return ""
	}
	return hex.EncodeToString(id[:4])
}
