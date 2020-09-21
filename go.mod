module github.com/planetdecred/godcr

go 1.13

require (
	gioui.org v0.0.0-20200722191435-e381ff40d66b
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/markbates/pkger v0.14.1
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.8.1
	github.com/planetdecred/dcrlibwallet v1.5.2
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20200618115811-c13761719519
)

// TODO: Remove and use an actual release of dcrlibwallet
replace (
	github.com/decred/dcrwallet/wallet/v3 => github.com/raedahgroup/dcrwallet/wallet/v3 v3.2.1-badger
	github.com/planetdecred/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
)
