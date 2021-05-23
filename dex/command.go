package dex

import (
	"errors"
	"fmt"

	"decred.org/dcrdex/client/core"
)

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

func (d *Dex) GetUser() {
	go func() {
		var resp Response
		resp.Resp = User{
			Info: *d.core.User(),
		}
		d.Send <- resp
	}()
}

func (d *Dex) AutoWalletConfig(assetID uint32) map[string]string {
	cfg, err := d.core.AutoWalletConfig(assetID)
	if err != nil {
		return nil
	}

	return cfg
}

func (d *Dex) UnlockWallet(assetID uint32, appPW []byte, errChan chan error) {
	go func() {
		status := d.core.WalletState(assetID)

		if status == nil {
			go func() {
				errChan <- errors.New(fmt.Sprintf("No wallet for %d", assetID))
			}()

			return
		}

		err := d.core.OpenWallet(assetID, appPW)
		if err != nil {
			log.Errorf("UnlockWallet error: %v", err)

			go func() {
				errChan <- errors.New(fmt.Sprintf("error unlocking %s wallet: %v", assetID, err))
			}()

			return
		}

		go func() {
			errChan <- nil
		}()
	}()
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
