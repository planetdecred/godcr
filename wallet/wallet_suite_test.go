package wallet_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWallet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wallet Suite")
}
