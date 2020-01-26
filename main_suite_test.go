package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGodcrGio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GodcrGio Suite")
}
