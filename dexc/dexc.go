package dexc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/asset/bch"
	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/asset/ltc"
	"decred.org/dcrdex/client/comms"
	"decred.org/dcrdex/dex/msgjson"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/dcrlibwallet/dexdcr"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
)

// Register the supported asset drivers here since importing
// the packages alone no longer auto-registers them. The dcr
// asset driver isn't registered here, because it requires the
// user to select a wallet from the multiwallet instance. The
// dcr asset driver can be registered by calling
// *Dexc.RegisterDcrAssetDriver after the wallet to connect is
// gotten, either from saved data (on restarts) or from user
// selection when connecting the dcr wallet for the first time.
func init() {
	bch.RegisterDriver()
	btc.RegisterDriver()
	ltc.RegisterDriver()
}

// Dexc represents of the Decred exchange client.
type Dexc struct {
	*core.Core
	ctx    context.Context
	cancel context.CancelFunc

	DbPath             string
	IsRunning          bool
	connectedDcrWallet *dcrlibwallet.Wallet

	connMtx sync.RWMutex
	conns   map[string]*dexConnection
}

const (
	DefaultAsset            = "dcr"
	DefaultAssetID   uint32 = 42
	ConversionFactor        = 1e8

	// ConnectedDcrWalletIDConfigKey is used as the key in a simple key-value
	// database to store the ID of the dcr wallet that is connected to the DEX
	// client, to facilitate reconnecting the wallet when godcr is restarted.
	ConnectedDcrWalletIDConfigKey = "dex_dcr_wallet_id"
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
		Core:   clientCore,
		DbPath: dbPath,
		conns:  make(map[string]*dexConnection),
	}, nil
}

// Start calls the Run method of the DEX client. Provide the app-wide context
// to ensure that the DEX client is shut down when the app-wide context is
// canceled.
func (d *Dexc) Start(appCtx context.Context) {
	// Create a new cancelFunc so that the app-wide ctx isn't canceled
	// when Core stops.
	d.ctx, d.cancel = context.WithCancel(appCtx)
	go func() {
		d.IsRunning = true
		d.Core.Run(d.ctx)
		d.cancel()
		d.IsRunning = false
	}()
	<-d.Core.Ready()
}

// Shutdown causes the dex client to shutdown. The shutdown attempt will be
// prevented if there are active orders or if some other error occurs. The
// bool return indicates if shutdown was successful. If successful, dexc will
// need to be restarted before it can be used again.
func (d *Dexc) Shutdown() bool {
	return d.shutdown(false)
}

// Shutdown causes the dex client to shutdown regardless of whether or not there
// are active orders. Dexc will need to be restarted before it can be used again.
func (d *Dexc) ForceShutdown() {
	d.shutdown(true)
}

// shutdown causes the dex client to shutdown. If there are active orders,
// this shutdown attempt will fail unless `forceShutdown` is true. If shutdown
// succeeds, dexc will need to be restarted before it can be used.
func (d *Dexc) shutdown(forceShutdown bool) bool {
	err := d.Logout()
	if err != nil {
		log.Errorf("Unable to logout of the dex client: %v", err)
		if !forceShutdown { // abort shutdown because of the error since forceShutdown != true
			return false
		}
	}

	// Cancel the ctx used to run Core.
	if d.cancel != nil { // in case dexc was never actually started
		// TODO: This doesn't really trigger Core's context cancellation.
		// Investigate.
		d.cancel()
	}
	d.IsRunning = false
	return true
}

// RegisterDcrAssetDriver uses the provided *dcrlibwallet.Wallet to register
// a dcr asset driver that performs wallet operations using the provided
// wallet.
func (d *Dexc) RegisterDcrAssetDriver(wallet *dcrlibwallet.Wallet) error {
	if d.connectedDcrWallet != nil {
		// This prevents attempting to register the driver multiple times
		// as that would lead to a panic.
		return fmt.Errorf("cannot change dcr wallet after initial setting, try restarting godcr")
	}

	walletDesc := fmt.Sprintf("%d (%s)", wallet.ID, wallet.Name)
	dexdcr.RegisterDriver(wallet.Internal(), walletDesc)
	d.connectedDcrWallet = wallet
	log.Infof("Registered decred wallet driver for DEX using wallet %s", walletDesc)
	return nil
}

// AddWallet attempts to connect the wallet with the provided details to the
// DEX client.
// NOTE: Before connecting a decred wallet, first register the decred asset
// driver by calling *Dexc.RegisterDcrAssetDriver.
func (d *Dexc) AddWallet(assetID uint32, settings map[string]string, appPW, walletPW []byte) error {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		// Only possible error is unknown asset. Should never happen because
		// we register most supported assets when this package is init'ed.
		// If connecting a decred wallet, caller should first register the
		// dcr asset driver by calling (*Dexc).RegisterDcrAssetDriver.
		return fmt.Errorf("asset driver not registered for asset with BIP ID %d", assetID)
	}

	// Start building the wallet config with default values.
	config := map[string]string{}
	for _, option := range assetInfo.ConfigOpts {
		config[strings.ToLower(option.Key)] = fmt.Sprintf("%v", option.DefaultValue)
	}

	// Attempt to load additional config values from the asset's default
	// config file path. Not necessary for dcr wallets.
	if assetID != dcr.BipID {
		autoConfig, err := d.AutoWalletConfig(assetID)
		if err != nil {
			return err
		}
		for k, v := range autoConfig {
			config[k] = v
		}
	}

	// User-provided settings should override any previously
	// set config value.
	for k, v := range settings {
		config[k] = v
	}

	return d.CreateWallet(appPW, walletPW, &core.WalletForm{
		AssetID: dcr.BipID,
		Config:  config,
	})
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
	acct, err := d.Core.AccountExport(password, host)
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
