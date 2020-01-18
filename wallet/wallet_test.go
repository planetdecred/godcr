package wallet_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/raedahgroup/godcr-gio/event"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

var _ = Describe("Wallet", func() {
	Context(`works with an almost zero wallet`, func() {
		It(`passes an error to SendChan when Sync is called`, func() {
			send := make(chan event.Event)
			wal := &Wallet{
				SendChan: send,
			}
			go wal.Sync()

			// TODO: check that it sends the correct error
			e := <-send
			_, ok := e.(error)
			Expect(ok).To(Equal(true))

		})
	})
})
