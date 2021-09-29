package values

import "github.com/planetdecred/godcr/ui/values/localizable"

var (
	ArrLanguages          map[string]string
	ArrExchangeCurrencies map[string]string
)

func init() {
	ArrLanguages = make(map[string]string)
	ArrLanguages[localizable.ENGLISH] = StrEnglish
	ArrLanguages[localizable.FRENCH] = StrFrench
	ArrLanguages[localizable.SPANISH] = StrSpanish
}
