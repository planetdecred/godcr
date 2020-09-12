package main

import "context"

type Wallet struct {
}

type handler struct {
	fn     func(*Wallet, context.Context, interface{}) (interface{}, error)
	noHelp bool
}

var handlers = map[string]handler{
	"balance": {fn: (*Wallet).getbalance},
}

func NewWallet() (*Wallet, error) {

	return &Wallet{}, nil
}

func (w *Wallet) getbalance(ctx context.Context, i interface{}) (interface{}, error) {
	wal := &Wallet{}
	return wal, nil
}
