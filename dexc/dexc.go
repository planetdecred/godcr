package dexc

import (
	"context"
	"fmt"
	"os"
	"strings"

	"decred.org/dcrdex/client/asset"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/dcrlibwallet/dexdcr"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
)

// Dexc represents the Decred Dex client.
type Dexc struct {
	*core.Core
	ctx    context.Context
	cancel context.CancelFunc

	coreConfig *core.Config
	IsRunning  bool
	IsLoggedIn bool // Keep user logged in state
}

const (
	// ConnectedDcrWalletIDConfigKey is used as the key in a simple key-value
	// database to store the ID of the dcr wallet that is connected to the DEX
	// client, to facilitate reconnecting the wallet when godcr is restarted.
	ConnectedDcrWalletIDConfigKey = "dex_dcr_wallet_id"
)

func NewDex(dbPath, net string, logMaker *dex.LoggerMaker) (*Dexc, error) {
	if net == "testnet3" {
		net = "testnet"
	}
	n, err := dex.NetFromString(net)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	coreConfig := &core.Config{
		DBPath: dbPath,
		Net:    n,
		Logger: logMaker.Logger("CORE"),
	}
	clientCore, err := core.New(coreConfig)
	if err != nil {
		log.Errorf("error creating client core: %v", err)
		return nil, err
	}

	return &Dexc{
		Core:       clientCore,
		coreConfig: coreConfig,
	}, nil
}

// Start calls the Run method of the DEX client. Provide the app-wide context
// to ensure that the DEX client is shut down when the app-wide context is
// canceled.
func (d *Dexc) Start(appCtx context.Context) error {
	// Re-setup Core if it was previously shutdown.
	if d.Core == nil {
		var err error
		d.Core, err = core.New(d.coreConfig)
		if err != nil {
			log.Errorf("error re-creating client core: %v", err)
			return err
		}
	}

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

	return nil
}

// Reset attempts to shutdown Core if it is running and if successful, deletes
// the DEX client database.
func (d *Dexc) Reset() bool {
	shutdownOk := d.shutdown(false)
	if shutdownOk {
		err := os.RemoveAll(d.coreConfig.DBPath)
		if err != nil {
			log.Warnf("DEX client reset failed: erroring deleting DEX db: %v", err)
			return false
		}
	}
	return shutdownOk
}

// Shutdown causes the dex client to shutdown. The shutdown attempt will be
// prevented if there are active orders or if some other error occurs. The
// bool return indicates if shutdown was successful. If successful, dexc will
// need to be restarted before it can be used again.
func (d *Dexc) Shutdown() bool {
	return d.shutdown(false)
}

// ForceShutdown causes the dex client to shutdown regardless of whether or not there
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
		d.cancel()
	}
	d.IsRunning = false
	d.IsLoggedIn = false
	d.Core = nil // Clear this to prevent panic in d.Core.Run if (*Dexc).Start() is re-called.
	return true
}

// SetWalletForDcrAsset configures the DEX client to use the provided wallet
// for dcr wallet operations. This method should be called before attempting
// to connect a dcr wallet to the DEX client for the first time. After this
// initial dcr wallet connection, this method should always be called with the
// same *dcrlibwallet.Wallet before starting the DEX client to ensure that Core
// is able to load the previously connected dcr wallet on startup.
func (d *Dexc) SetWalletForDcrAsset(wallet *dcrlibwallet.Wallet) {
	walletDesc := fmt.Sprintf("%d (%s)", wallet.ID, wallet.Name)
	dexdcr.UseSpvWalletForDexClient(wallet.Internal(), walletDesc)
}

// AddWallet attempts to connect the wallet with the provided details to the
// DEX client.
// NOTE: Before connecting a dcr wallet, first call *Dexc.SetWalletForDcrAsset
// to configure the dcr ExchangeWallet to use a custom wallet instead of the
// default rpc wallet.
func (d *Dexc) AddWallet(assetID uint32, walletType string, settings map[string]string, appPW, walletPW []byte) error {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Errorf("asset driver not registered for asset with BIP ID %d", assetID)
	}
	var walletDef *asset.WalletDefinition
	for _, def := range assetInfo.AvailableWallets {
		if def.Type == walletType {
			walletDef = def
		}
	}
	if walletDef == nil {
		return fmt.Errorf("cannot add %s wallet of type %q", assetInfo.Name, walletType)
	}

	// Start building the wallet config with default values.
	config := map[string]string{}
	for _, option := range walletDef.ConfigOpts {
		config[strings.ToLower(option.Key)] = fmt.Sprintf("%v", option.DefaultValue)
	}

	// User-provided settings should override any previously
	// set config value.
	for k, v := range settings {
		config[k] = v
	}

	return d.CreateWallet(appPW, walletPW, &core.WalletForm{
		AssetID: assetID,
		Config:  config,
		Type:    walletType,
	})
}
