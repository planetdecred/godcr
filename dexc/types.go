package dexc

import (
	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/encode"
)

// Response represents a discriminated union for wallet responses.
// Either Resp or Err must be nil.
type Response struct {
	Resp interface{}
	Err  error
}

type NewWalletForm struct {
	AssetID uint32
	Config  map[string]string
	Pass    encode.PassBytes
	AppPW   encode.PassBytes
}

type User struct {
	Info core.User
}

type TradeForm struct {
	Pass  encode.PassBytes
	Order *core.TradeForm
}

// MaxOrderEstimate is sent when the dex core is done getting max order estimate.
type MaxOrderEstimate struct {
	MaxOrderEstimate *core.MaxOrderEstimate
}
