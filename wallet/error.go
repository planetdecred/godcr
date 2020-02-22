package wallet

import (
	"errors"

	"github.com/raedahgroup/dcrlibwallet"
)

var (
	// ErrIDNotExist is returned when a given ID does not exist
	ErrIDNotExist = errors.New("ID does not exist")

	ErrBadPass = errors.New(dcrlibwallet.ErrInvalidPassphrase)
)

// InternalWalletError represents errors generated during the handling of the multiwallet
// and connected wallets
type InternalWalletError struct {
	Message  string
	Affected []int
	Err      error
}

func (err InternalWalletError) Unwrap() error {
	return err.Err
}

func (err InternalWalletError) Error() string {
	return err.Message
}

type MultiWalletError struct {
	Message string
	Err     error
}

func (err MultiWalletError) Error() string {
	return err.Message
}

func (err MultiWalletError) Unwrap() error {
	return err.Err
}

func ResponseError(err error) Response {
	return Response{
		Err: err,
	}
}

func ResponseResp(resp interface{}) Response {
	return Response{
		Resp: resp,
	}
}
