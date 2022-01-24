package load

import (
	"errors"
	"fmt"
	"strings"

	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

func validateVSPServerSignature(resp *http.Response, pubKey, body []byte) error {
	sigStr := resp.Header.Get("VSP-Server-Signature")
	sig, err := base64.StdEncoding.DecodeString(sigStr)
	if err != nil {
		return fmt.Errorf("error validating VSP signature: %v", err)
	}

	if !ed25519.Verify(pubKey, body, sig) {
		return errors.New("bad signature from VSP")
	}

	return nil
}

func getVSPInfo(url string) (*dcrlibwallet.VspInfoResponse, error) {
	rq := new(http.Client)
	resp, err := rq.Get((url + "/api/v3/vspinfo"))

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 response from server: %v", string(b))
	}

	var vspInfoResponse dcrlibwallet.VspInfoResponse
	err = json.Unmarshal(b, &vspInfoResponse)
	if err != nil {
		return nil, err
	}

	err = validateVSPServerSignature(resp, vspInfoResponse.PubKey, b)
	if err != nil {
		return nil, err
	}
	return &vspInfoResponse, nil
}

// getInitVSPInfo returns the list information of the VSP
func getInitVSPInfo(url string) (map[string]*dcrlibwallet.VspInfoResponse, error) {
	rq := new(http.Client)
	resp, err := rq.Get((url))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 response from server: %v", string(b))
	}

	var vspInfoResponse map[string]*dcrlibwallet.VspInfoResponse
	err = json.Unmarshal(b, &vspInfoResponse)
	if err != nil {
		return nil, err
	}

	return vspInfoResponse, nil
}

func (wl *WalletLoad) GetVSPList() {
	var valueOut struct {
		Remember string
		List     []string
	}

	wl.MultiWallet.ReadUserConfigValue(dcrlibwallet.VSPHostConfigKey, &valueOut)
	var loadedVSP []*wallet.VSPInfo

	for _, host := range valueOut.List {
		v, err := getVSPInfo(host)
		if err == nil {
			loadedVSP = append(loadedVSP, &wallet.VSPInfo{
				Host: host,
				Info: v,
			})
		}
	}

	l, _ := getInitVSPInfo("https://api.decred.org/?c=vsp")
	var err error
	for h, v := range l {
		if strings.Contains(wl.Wallet.Net, v.Network) {
			v, err = getVSPInfo(fmt.Sprintf("https://%s", h))
			if err != nil {
				log.Error(err.Error())
				return
			}
			loadedVSP = append(loadedVSP, &wallet.VSPInfo{
				Host: fmt.Sprintf("https://%s", h),
				Info: v,
			})
		}
	}

	wl.VspInfo.List = loadedVSP
}

// TicketPrice get ticket price
func (wl *WalletLoad) TicketPrice() int64 {
	pr, err := wl.MultiWallet.WalletsIterator().Next().TicketPrice()
	if err != nil {
		log.Error(err)
		return 0
	}
	return pr.TicketPrice
}

func (wl *WalletLoad) AddVSP(host string) (err error) {
	var valueOut struct {
		Remember string
		List     []string
	}

	// check if host already exists
	_ = wl.MultiWallet.ReadUserConfigValue(dcrlibwallet.VSPHostConfigKey, &valueOut)
	for _, v := range valueOut.List {
		if v == host {
			return fmt.Errorf("existing host %s", host)
		}
	}

	// validate host network
	info, err := getVSPInfo(host)
	if err != nil {
		return err
	}

	if info.Network != wl.Wallet.Net {
		return fmt.Errorf("invalid net %s", info.Network)
	}

	valueOut.List = append(valueOut.List, host)
	wl.MultiWallet.SaveUserConfigValue(dcrlibwallet.VSPHostConfigKey, valueOut)
	(*wl.VspInfo).List = append((*wl.VspInfo).List, &wallet.VSPInfo{
		Host: host,
		Info: info,
	})

	return
}

func (wl *WalletLoad) GetRememberVSP() string {
	var valueOut struct {
		Remember string
	}
	wl.MultiWallet.ReadUserConfigValue(dcrlibwallet.VSPHostConfigKey, &valueOut)

	return valueOut.Remember
}

func (wl *WalletLoad) RememberVSP(host string) {
	var valueOut struct {
		Remember string
		List     []string
	}
	err := wl.MultiWallet.ReadUserConfigValue(dcrlibwallet.VSPHostConfigKey, &valueOut)
	if err != nil {
		log.Error(err.Error())
	}

	valueOut.Remember = host
	wl.MultiWallet.SaveUserConfigValue(dcrlibwallet.VSPHostConfigKey, valueOut)
}
