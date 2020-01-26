package wallet_test

import (
	"fmt"
	"os"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/raedahgroup/godcr-gio/event"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

const (
	testroot = ".godcr"
	testnet  = "testnet3"
)

var _ = Describe("Wallet", func() {
	Context(`with new multi wallet`, func() {
		var wal *Wallet
		var duplex event.Duplex
		BeforeEach(func() {
			Expect(os.RemoveAll(testroot)).To(BeNil())
			dup := event.NewDuplexBase()
			walb := NewWallet(testroot, testnet, dup.Duplex())
			wal = walb
			duplex = dup.Reverse()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(testroot)).To(BeNil())
		})

		It("shuts down properly", func(done Done) {
			var wg sync.WaitGroup
			wg.Add(1)
			go wal.Sync(&wg)

			<-duplex.Receive
			duplex.Send <- event.WalletCmd{
				Cmd: event.ShutdownCmd,
			}
			wg.Wait()
			close(done)
		}, 5.0)

		It("can create a new wallet", func(done Done) {
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

			words := strings.Split(seed, " ")
			fmt.Printf("Seed words: %+v\n", words)
			Expect(len(words)).To(Equal(33))

			duplex.Send <- event.WalletCmd{
				Cmd: event.ShutdownCmd,
			}
			wg.Wait()
			close(done)
		}, 5.0)

	})
	Context("with multiwallet with one new wallet", func() {
		var wal *Wallet
		var duplex event.Duplex
		var wg sync.WaitGroup
		BeforeEach(func(done Done) {
			Expect(os.RemoveAll(testroot)).To(BeNil())
			dup := event.NewDuplexBase()
			walb := NewWallet(testroot, testnet, dup.Duplex())
			wal = walb
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
			_, ok := e.(event.WalletResponse)
			Expect(ok).To(Equal(true))
			close(done)

			// seed, err := resp.Results.PopString()

			// Expect(err).To(BeNil())

			// words := strings.Split(seed, " ")
			// fmt.Printf("Seed words: %+v\n", words)
			// Expect(len(words)).To(Equal(33))

		}, 5.0)

		AfterEach(func(done Done) {
			duplex.Send <- event.WalletCmd{
				Cmd: event.ShutdownCmd,
			}
			wg.Wait()
			Expect(os.RemoveAll(testroot)).To(BeNil())
			close(done)
		}, 5.0)

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
})
