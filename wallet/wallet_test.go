package wallet_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/raedahgroup/godcr-gio/wallet"
)

var _ = Describe("LoadWallets(a, b)", func() {
	When(`a is "" and b is ""`, func() {
		It("should return an error", func() {
			wal, onewallet, err := LoadWallets("", "")
			Expect(err).ToNot(BeNil())
			Expect(wal).To(BeNil())
			Expect(onewallet).To(Equal(false))
		})
	})
})
