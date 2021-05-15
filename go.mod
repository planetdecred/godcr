module github.com/planetdecred/godcr

go 1.13

require (
	decred.org/dcrdex v1.0.0
	gioui.org v0.0.0-20210418151603-3b69b5ed0512
	github.com/ararog/timeago v0.0.0-20160328174124-e9969cf18b8d
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
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
	github.com/planetdecred/dcrlibwallet v1.5.3-0.20210224132742-5d0bc5e13370
	github.com/sqweek/dialog v0.0.0-20200911184034-8a3d98e8211d
	github.com/yeqown/go-qrcode v1.5.1
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/text v0.3.4
)

replace (
	decred.org/dcrdex => github.com/decred/dcrdex v0.0.0-20210504142407-8c1920c33811
	decred.org/dcrwallet => decred.org/dcrwallet v1.7.0
	github.com/decred/dcrdata/txhelpers/v4 => github.com/decred/dcrdata/txhelpers/v4 v4.0.0-20200108145420-f82113e7e212
)
