module github.com/planetdecred/godcr

go 1.13

require (
	gioui.org v0.0.0-20200722191435-e381ff40d66b
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/ararog/timeago v0.0.0-20160328174124-e9969cf18b8d
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/slog v1.1.0
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/jrick/logrotate v1.0.0
	github.com/markbates/pkger v0.17.1
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/planetdecred/dcrlibwallet v1.5.3-0.20210127151106-5c5bfcb04d6d
	github.com/russross/blackfriday v2.0.0+incompatible
	github.com/shurcooL/go v0.0.0-20200502201357-93f07166e636 // indirect
	github.com/shurcooL/markdownfmt v0.0.0-20200725144734-77d690767c81
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	gitlab.com/golang-commonmark/markdown v0.0.0-20191127184510-91b5b3c99c19
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20200618115811-c13761719519
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
)

// TODO: Remove and use an actual release of dcrlibwallet
replace (
	decred.org/dcrwallet => decred.org/dcrwallet v1.6.0-rc4
	//github.com/planetdecred/dcrlibwallet/spv => github.com/raedahgroup/dcrlibwallet/spv v0.0.0-20200113081741-39a55988f78c
	github.com/decred/dcrdata/txhelpers/v4 => github.com/decred/dcrdata/txhelpers/v4 v4.0.0-20200108145420-f82113e7e212
	github.com/decred/dcrwallet/wallet/v3 => github.com/raedahgroup/dcrwallet/wallet/v3 v3.2.1-badger
)
