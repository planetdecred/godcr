package wallet

import (
	"github.com/planetdecred/dcrlibwallet"
)

// TODO: responses.go file to be deprecated with future code clean up

type UnspentOutput struct {
	UTXO     dcrlibwallet.UnspentOutput
	Amount   string
	DateTime string
}

// UnspentOutputs wraps the dcrlibwallet UTXO type and adds processed data
type UnspentOutputs struct {
	List []*UnspentOutput
}
