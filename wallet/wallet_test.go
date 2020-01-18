package wallet_test

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

var _ = Describe("Sync", func() {
	It(`works with a zero wallet`, func() {
		wal := new(Wallet)
		wal.Sync()
		// just checking that this doesn't panic
	})
})
