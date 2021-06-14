module github.com/planetdecred/godcr

go 1.13

require (
	gioui.org v0.0.0-20210418151603-3b69b5ed0512
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/ararog/timeago v0.0.0-20160328174124-e9969cf18b8d
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.3-0.20200921185235-6d75c7ec1199
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/dcrd/dcrutil/v3 v3.0.0
	github.com/decred/slog v1.1.0
	github.com/gen2brain/beeep v0.0.0-20200526185328-e9c15c258e28
	github.com/gomarkdown/markdown v0.0.0-20210208175418-bda154fe17d8
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/jrick/logrotate v1.0.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/planetdecred/dcrlibwallet v1.6.0
	github.com/yeqown/go-qrcode v1.5.1
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/text v0.3.3
)

replace (
	github.com/decred/dcrdata/txhelpers/v4 => github.com/decred/dcrdata/txhelpers/v4 v4.0.0-20200108145420-f82113e7e212
	github.com/planetdecred/dcrlibwallet => github.com/planetdecred/dcrlibwallet v1.6.1-0.20210618124747-f776ef63322f
)
