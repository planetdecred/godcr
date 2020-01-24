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
	Context(`New`, func() {
		It(`fails with invalid arguments`, func() {
			_, _, err := New("", "")
			Expect(err).ToNot(BeNil())
		})
	})
	Context(`with valid wallet`, func() {
		var wal *Wallet
		var duplex event.Duplex
		BeforeEach(func() {
			Expect(os.RemoveAll(testroot)).To(BeNil())
			walb, dupb, err := New(testroot, testnet)
			Expect(err).To(BeNil())
			wal = walb
			duplex = dupb.Reverse()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(testroot)).To(BeNil())
		})

		It("shuts down properly", func(done Done) {
			var wg sync.WaitGroup
			wg.Add(1)
			go wal.Sync(&wg)

			duplex.Send <- event.WalletCmd{
				Cmd: event.ShutdownCmd,
			}
			wg.Wait()
			close(done)
		}, 3.0)

		FIt("can create a new wallet", func(done Done) {
			var wg sync.WaitGroup
			wg.Add(1)
			go wal.Sync(&wg)

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
		}, 10.0)
	})
})
