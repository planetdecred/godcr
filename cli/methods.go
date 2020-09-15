package main

import (
	"bufio"
	"fmt"
	"os"
)

type handler func(*Cli, interface{}) (interface{}, error)

var handlers = map[string]handler{
	"balance":      (*Cli).getbalance,
	"createwallet": (*Cli).createWallet,
}

func (c *Cli) createWallet(i interface{}) (interface{}, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	fmt.Print("Enter password: ")
	scanner.Scan()
	pwd := scanner.Text()
	c.wall.CreateWallet(pwd, c.errChan)
	return nil, nil
}

func (c *Cli) getbalance(i interface{}) (interface{}, error) {
	c.wall.GetMultiWalletInfo()
	return nil, nil
}
