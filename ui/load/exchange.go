package load

import (
	"encoding/json"
	"net/http"

	"golang.org/x/text/message"
)

func GetUSDExchangeValue(target interface{}) error {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	res, err := http.Get(url)
	// TODO: include user agent in req header
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func FormatUSDBalance(p *message.Printer, balance float64) string {
	return p.Sprintf("$%.2f", balance)
}

func DCRToUSD(exchangeRate, dcr float64) float64 {
	return dcr * exchangeRate
}

func USDToDCR(exchangeRate, usd float64) float64 {
	return usd / exchangeRate
}
