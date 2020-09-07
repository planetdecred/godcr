package cli

// import (
// 	"fmt"
// 	"strconv"

// 	rpcclient "github.com/raedahgroup/dcrcli/walletrpcclient"
// 	qrcode "github.com/skip2/go-qrcode"
// )

// func balance(walletrpcclient *rpcclient.Client, commandArgs []string) (*response, error) {
// 	balances, err := walletrpcclient.Balance()
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := &response{
// 		columns: []string{
// 			"Account",
// 			"Total",
// 			"Spendable",
// 			"Locked By Tickets",
// 			"Voting Authority",
// 			"Unconfirmed",
// 		},
// 		result: make([][]interface{}, len(balances)),
// 	}
// 	for i, v := range balances {
// 		res.result[i] = []interface{}{
// 			v.AccountName,
// 			v.Total,
// 			v.Spendable,
// 			v.LockedByTickets,
// 			v.VotingAuthority,
// 			v.Unconfirmed,
// 		}
// 	}

// 	return res, nil
// }

// func normalSend(walletrpcclient *rpcclient.Client, _ []string) (*response, error) {
// 	return send(walletrpcclient, false)
// }

// func customSend(walletrpcclient *rpcclient.Client, _ []string) (*response, error) {
// 	return send(walletrpcclient, true)
// }

// func send(walletrpcclient *rpcclient.Client, custom bool) (*response, error) {
// 	var err error

// 	sourceAccount, err := getSendSourceAccount(walletrpcclient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// check if account has positive non-zero balance before proceeding
// 	// if balance is zero, there'd be no unspent outputs to use
// 	accountBalance, err := walletrpcclient.SingleAccountBalance(sourceAccount, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if accountBalance.Total == 0 {
// 		return nil, fmt.Errorf("Selected account has 0 balance. Cannot proceed")
// 	}

// 	destinationAddress, err := getSendDestinationAddress(walletrpcclient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sendAmount, err := getSendAmount()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var utxoSelection []string
// 	if custom {
// 		// get all utxos in account, pass 0 amount to get all
// 		utxos, err := walletrpcclient.UnspentOutputs(sourceAccount, 0)
// 		if err != nil {
// 			return nil, err
// 		}

// 		utxoSelection, err = getUtxosForNewTransaction(utxos, sendAmount)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	passphrase, err := getWalletPassphrase()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result *rpcclient.SendResult
// 	if custom {
// 		result, err = walletrpcclient.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount,
// 			destinationAddress, passphrase)
// 	} else {
// 		result, err = walletrpcclient.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	res := &response{
// 		columns: []string{
// 			"Result",
// 			"Hash",
// 		},
// 		result: [][]interface{}{
// 			[]interface{}{
// 				"The transaction was published successfully",
// 				result.TransactionHash,
// 			},
// 		},
// 	}

// 	return res, nil
// }

// func receive(walletrpcclient *rpcclient.Client, commandArgs []string) (*response, error) {
// 	var recieveAddress uint32

// 	// if no address passed in
// 	if len(commandArgs) == 0 {

// 		// display menu options to select account
// 		var err error
// 		recieveAddress, err = getSendSourceAccount(walletrpcclient)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		// if an address was passed in eg. ./dcrcli receive 0 use that address
// 		x, err := strconv.ParseUint(commandArgs[0], 10, 32)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error parsing account number: %s", err.Error())
// 		}

// 		recieveAddress = uint32(x)
// 	}

// 	r, err := walletrpcclient.Receive(recieveAddress)
// 	if err != nil {
// 		return nil, err
// 	}

// 	qr, err := qrcode.New(r.Address, qrcode.Medium)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error generating QR Code: %s", err.Error())
// 	}

// 	res := &response{
// 		columns: []string{
// 			"Address",
// 			"QR Code",
// 		},
// 		result: [][]interface{}{
// 			[]interface{}{
// 				r.Address,
// 				qr.ToString(true),
// 			},
// 		},
// 	}
// 	return res, nil
// }
