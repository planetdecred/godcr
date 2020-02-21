package wallet_test

import (
	"fmt"
	"os"
	"time"

	"github.com/decred/dcrd/dcrutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

const (
	testnet = "testnet3"
)

func getTestDir() string {
	now := time.Now().UTC().Unix()
	testDir := fmt.Sprintf(".godcr_test_%d", now)
	_, err := os.Stat(testDir)
	i := 1
	for err == nil {
		testDir = fmt.Sprintf(".godcr_test_%d_%d", now, i)
		_, err = os.Stat(testDir)
		i++
	}
	err = os.Mkdir(testDir, os.ModePerm)
	Expect(err).To(BeNil())
	return testDir
}

var (
	wal *Wallet
)
var _ = BeforeSuite(func() {
	var err error
	wal, err = NewWallet(getTestDir(), testnet, make(chan Response), 2)
	Expect(err).To(BeNil())
	wal.LoadWallets()
	resp := <-wal.Send
	Expect(resp.Resp).To(BeAssignableToTypeOf(LoadedWallets{}))
	wal.CreateWallet("password")
	resp = <-wal.Send
	Expect(resp.Resp).To(BeAssignableToTypeOf(CreatedSeed{}))
})

var _ = AfterSuite(func() {
	wal.Shutdown()
})

var _ = Describe("Wallet", func() {
	It("can get the multi wallet info", func() {
		wal.GetMultiWalletInfo()
		info := <-wal.Send
		Expect(info.Resp).To(BeAssignableToTypeOf(&MultiWalletInfo{}))
		inf := info.Resp.(*MultiWalletInfo)
		Expect(inf.LoadedWallets).To(BeEquivalentTo(1))
		Expect(inf.TotalBalance).To(BeEquivalentTo(dcrutil.Amount(0).String()))
		Expect(inf.Synced).To(Equal(false))
	})
	It("can rename a wallet", func() {
		err := wal.RenameWallet(1, "random")
		Expect(err).To(BeNil())
	})
	It("can get the current address", func() {
		addr, err := wal.CurrentAddress(1, 0)
		Expect(err).To(BeNil())
		Expect(wal.IsAddressValid(addr)).To(Equal(true))
		addr2, err := wal.CurrentAddress(1, 0)
		Expect(err).To(BeNil())
		Expect(addr).To(Equal(addr2))
	})
	It("can create a new address", func() {
		addr, err := wal.NextAddress(1, 0)
		Expect(err).To(BeNil())
		Expect(wal.IsAddressValid(addr)).To(Equal(true))
	})
})
