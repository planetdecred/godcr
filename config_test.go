package main

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("No config file", func() {
		BeforeEach(func() {
			os.RemoveAll(defaultHomeDir)
		})
		It("should return the default config", func() {

			cfg, err := loadConfig()
			if err != nil {
				Fail(err.Error())
			}
			Expect(*cfg).To(Equal(defaultConfig))
		})
		It("should create the default config file", func() {
			cfg, err := loadConfig()
			if err != nil {
				Fail(err.Error())
			}

			_, err = os.Stat(cfg.ConfigFile)

			Expect(err).To(BeNil())
		})
	})

})
