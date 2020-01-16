module github.com/raedahgroup/godcr-gio

go 1.13

require (
	gioui.org v0.0.0-20200113125233-3f6a1c34d37d
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrwallet/walletseed v1.0.1
	github.com/decred/slog v1.0.0
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20200112073453-40a403ab6f3a
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
)

// TODO: Remove and use an actual release of dcrlibwallet
replace github.com/raedahgroup/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
