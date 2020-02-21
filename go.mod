module github.com/raedahgroup/godcr-gio

go 1.13

require (
	gioui.org v0.0.0-20200223075147-73b99a80e29a
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/markbates/pkger v0.14.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20200130094829-d902833f4f05
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
)

// TODO: Remove and use an actual release of dcrlibwallet
replace github.com/raedahgroup/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
