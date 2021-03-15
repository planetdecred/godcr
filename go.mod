module github.com/planetdecred/godcr

go 1.13

require (
	gioui.org v0.0.0-20210228180843-e1248651c871
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/dcrd/dcrutil/v3 v3.0.0
	github.com/decred/slog v1.1.0
	github.com/gen2brain/beeep v0.0.0-20200526185328-e9c15c258e28
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/jrick/logrotate v1.0.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/planetdecred/dcrlibwallet v1.5.3-0.20210224132742-5d0bc5e13370
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20200618115811-c13761719519
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)

replace (
	decred.org/dcrwallet => decred.org/dcrwallet v1.6.0-rc4
	github.com/decred/dcrdata/txhelpers/v4 => github.com/decred/dcrdata/txhelpers/v4 v4.0.0-20200108145420-f82113e7e212
)
