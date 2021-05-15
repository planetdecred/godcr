package dex

import (
	"errors"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/encode"
)

type NewWalletForm struct {
	AssetID uint32
	Config  map[string]string
	Pass    encode.PassBytes
	AppPW   encode.PassBytes
}

func (d *Dex) InitializeClient(apppasswd string, errChan chan error) {
	go func() {
		err := d.core.InitializeClient([]byte(apppasswd))
		if err != nil {
			go func() {
				errChan <- err
			}()

			return
		}

		go func() {
			errChan <- nil
		}()
	}()
}

func (d *Dex) SupportedAsset() map[uint32]*core.SupportedAsset {
	return d.core.SupportedAssets()
}

func (d *Dex) IsInitialized() bool {
	ok, err := d.core.IsInitialized()
	if err != nil {
		log.Error(err)
	}
	return ok
}

func (d *Dex) GetUser() *core.User {
	u := d.core.User()
	return u
}

func (d *Dex) GetDefaultWalletConfig() map[string]string {
	cfg, err := d.core.AutoWalletConfig(DefaultAssert)
	if err != nil {
		return nil
	}
	return cfg
}

func (d *Dex) UnlockWallet(assetID uint32, appPW []byte) error {
	status := d.core.WalletState(assetID)
	if status == nil {
		return errors.New(fmt.Sprintf("No wallet for %d", assetID))
	}

	err := d.core.OpenWallet(assetID, appPW)
	if err != nil {
		return errors.New(fmt.Sprintf("error unlocking %s wallet: %v", assetID, err))
	}

	return nil
}

func (d *Dex) GetDEXConfig(addr string, cert string, errChan chan error, responseChan chan *core.Exchange) {
	go func() {
		cx, err := d.core.GetDEXConfig(addr, []byte(cert))
		if err != nil {
			go func() {
				errChan <- err
			}()

			return
		}

		go func() {
			responseChan <- cx
		}()

		return
	}()
}

func (d *Dex) Register(apppasswd string, addr string, fee uint64, cert string, errChan chan error) {
	go func() {
		form := &core.RegisterForm{
			AppPass: []byte(apppasswd),
			Addr:    addr,
			Fee:     fee,
			Cert:    []byte(cert),
		}
		_, err := d.core.Register(form)

		if err != nil {
			go func() {
				errChan <- err
			}()

			return
		}

		go func() {
			errChan <- nil
		}()

		return
	}()
}

func (d *Dex) Login(apppasswd string, errChan chan error) {
	go func() {
		_, err := d.core.Login([]byte(apppasswd))
		if err != nil {
			go func() {
				errChan <- err
			}()

			return
		}

		go func() {
			errChan <- nil
		}()
	}()
}

func (d *Dex) AddNewWallet(form *NewWalletForm, errChan chan error) {
	go func() {
		has := d.core.WalletState(form.AssetID) != nil
		if has {
			go func() {
				errChan <- errors.New(fmt.Sprintf("already have a wallet for %d", form.AssetID))
			}()

			return
		}

		// Wallet does not exist yet. Try to create it.
		err := d.core.CreateWallet(form.AppPW, form.Pass, &core.WalletForm{
			AssetID: form.AssetID,
			Config:  form.Config,
		})
		if err != nil {
			go func() {
				errChan <- err
			}()

			return
		}

		go func() {
			errChan <- nil
		}()
	}()
}
