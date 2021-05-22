package dex

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
