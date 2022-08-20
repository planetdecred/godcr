package values

import "github.com/planetdecred/godcr/ui/values/localizable"

var (
	ArrLanguages          map[string]string
	ArrExchangeCurrencies map[string]string
	ArrMixerAccounts      map[string]string
)

const (
	DefaultExchangeValue = "none"
	USDExchangeValue     = "USD (Bittrex)"

	DefaultAccount = StrDefault
	MixedAcc       = StrMixed
	UnmixedAcc     = StrUnmixed
)

func init() {
	ArrLanguages = make(map[string]string)
	ArrLanguages[localizable.ENGLISH] = StrEnglish
	ArrLanguages[localizable.FRENCH] = StrFrench
	ArrLanguages[localizable.SPANISH] = StrSpanish

	ArrExchangeCurrencies = make(map[string]string)
	ArrExchangeCurrencies[DefaultExchangeValue] = StrNone
	ArrExchangeCurrencies[USDExchangeValue] = StrUsdBittrex

	ArrMixerAccounts = make(map[string]string)
	ArrMixerAccounts[DefaultAccount] = StrDefault
	ArrMixerAccounts[MixedAcc] = StrMixed
	ArrMixerAccounts[UnmixedAcc] = StrUnmixed
}
