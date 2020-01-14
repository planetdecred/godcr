package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	main "github.com/raedahgroup/godcr-gio"
)

var _ = Describe("Version", func() {
	It("should return 0.0.0", func() {
		Expect(main.Version()).To(Equal("0.0.0"))
	})
})
