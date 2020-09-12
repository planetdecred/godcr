module github.com/raedahgroup/cli

go 1.13

require (
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20200130094829-d902833f4f05
	github.com/raedahgroup/godcr v1.0.0
)

// TODO: Remove and use an actual release of dcrlibwallet
replace (
	github.com/raedahgroup/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
	github.com/raedahgroup/godcr => ../
)
