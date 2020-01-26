package wallet_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/raedahgroup/godcr-gio/event"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

const (
	testroot = ".godcr"
	testnet  = "testnet3"
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

var _ = Describe("Completely new Wallet", func() {
	var (
		wal     *Wallet
		duplex  event.Duplex
		testDir string
		info    *os.File
		writer  *bufio.Writer
	)
	BeforeEach(func() {
		Expect(os.RemoveAll(testroot)).To(BeNil())
		dup := event.NewDuplexBase()
		testDir = getTestDir()
		walb := NewWallet(testDir, testnet, dup.Duplex())
		wal = walb
		duplex = dup.Reverse()

		testDesc := CurrentGinkgoTestDescription()

		info, err := os.Create(filepath.Join(testDir, "info.txt"))
		Expect(err).To(BeNil())

		writer = bufio.NewWriter(info)

		_, err = writer.WriteString("Test: " + testDesc.FullTestText + ".\n")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		writer.Flush()
		info.Close()
	})

	It("shuts down properly", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go wal.Sync(&wg)

		<-duplex.Receive
		duplex.Send <- event.WalletCmd{
			Cmd: event.ShutdownCmd,
		}
		wg.Wait()
	})

	It("can create a new wallet", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go wal.Sync(&wg)

		<-duplex.Receive

		duplex.Send <- event.WalletCmd{
			Cmd: event.CreateCmd,
			Arguments: &event.ArgumentQueue{
				Queue: []interface{}{"password", 1},
			},
		}

		e := <-duplex.Receive
		resp, ok := e.(event.WalletResponse)
		Expect(ok).To(Equal(true))

		seed, err := resp.Results.PopString()
		Expect(err).To(BeNil())
		Expect(len(strings.Split(seed, " "))).To(Equal(33))

		duplex.Send <- event.WalletCmd{
			Cmd: event.ShutdownCmd,
		}
		wg.Wait()
	})
})

var _ = Describe("Wallet with one newly created wallet", func() {
	var (
		wal     *Wallet
		duplex  event.Duplex
		wg      sync.WaitGroup
		testDir string
		info    *os.File
		writer  *bufio.Writer
	)
	BeforeEach(func() {
		dup := event.NewDuplexBase()
		testDir = getTestDir()
		wal = NewWallet(testDir, testnet, dup.Duplex())
		duplex = dup.Reverse()

		wg.Add(1)
		go wal.Sync(&wg)
		<-duplex.Receive

		duplex.Send <- event.WalletCmd{
			Cmd: event.CreateCmd,
			Arguments: &event.ArgumentQueue{
				Queue: []interface{}{"password", 1},
			},
		}

		e := <-duplex.Receive
		resp, ok := e.(event.WalletResponse)
		Expect(ok).To(Equal(true))

		seed, err := resp.Results.PopString()
		Expect(err).To(BeNil())
		Expect(len(strings.Split(seed, " "))).To(Equal(33))

		testDesc := CurrentGinkgoTestDescription()

		info, err = os.Create(filepath.Join(testDir, "info.txt"))
		Expect(err).To(BeNil())

		writer = bufio.NewWriter(info)

		_, err = writer.WriteString("Test: " + testDesc.FullTestText + ".\n")
		Expect(err).To(BeNil())

		_, err = writer.WriteString("Seed words: " + seed + "\n")
		Expect(err).To(BeNil())

	})

	AfterEach(func() {
		duplex.Send <- event.WalletCmd{
			Cmd: event.ShutdownCmd,
		}
		wg.Wait()
		writer.Flush()
		info.Close()
	})

	It("returns 0 for total balance", func() {
		duplex.Send <- event.WalletCmd{
			Cmd: event.InfoCmd,
		}

		e := <-duplex.Receive
		info, ok := e.(event.WalletInfo)
		By("returning a WalletInfo")
		Expect(ok).To(Equal(true), fmt.Sprintf("Actual val: %+v", e))

		By("LoadedWallets == 1")
		Expect(info.LoadedWallets).To(Equal(1))

		By("TotalBalance == 0")
		Expect(info.TotalBalance).To(Equal(int64(0)))
	})
	// FIt("syncs properly", func(done Done) {
	// 	duplex.Send <- event.WalletCmd{
	// 		Cmd: event.StartSyncCmd,
	// 	}
	// 	e := <-duplex.Receive
	// 	evt, ok := e.(event.Sync)
	// 	By("returning a SyncStart event")
	// 	Expect(ok).To(Equal(true), fmt.Sprintf("Actual val: %+v", e))
	// 	Expect(evt.Event).To(Equal(event.SyncStart))

	// 	e = <-duplex.Receive
	// 	evt, ok = e.(event.Sync)
	// 	By("returning a SyncEnd event")
	// 	Expect(ok).To(Equal(true), fmt.Sprintf("Actual val: %+v", e))
	// 	Expect(evt.Event).To(Equal(event.SyncEnd))
	// 	close(done)
	// }, 30*60*60)
})
