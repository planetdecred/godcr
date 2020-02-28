module github.com/raedahgroup/godcr-gio

go 1.13

require (
	gioui.org v0.0.0-20200116122050-18cddc030077
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/dcrwallet/wallet/v3 v3.0.5
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/markbates/pkger v0.14.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20200130094829-d902833f4f05
)

// TODO: Remove and use an actual release of dcrlibwallet
replace github.com/raedahgroup/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
