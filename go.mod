module github.com/planetdecred/godcr

go 1.16

require (
	gioui.org v0.0.0-20211011183043-05f0f5c20f45
	github.com/JohannesKaufmann/html-to-markdown v1.2.1
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/ararog/timeago v0.0.0-20160328174124-e9969cf18b8d
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/dcrd/dcrutil/v3 v3.0.0
	github.com/decred/slog v1.2.0
	github.com/gen2brain/beeep v0.0.0-20210529141713-5586760f0cc1
	github.com/godbus/dbus/v5 v5.0.5 // indirect
	github.com/gomarkdown/markdown v0.0.0-20210208175418-bda154fe17d8
	github.com/gopherjs/gopherjs v0.0.0-20210901121439-eee08aaf2717 // indirect
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/jrick/logrotate v1.0.0
	github.com/nxadm/tail v1.4.4
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/planetdecred/dcrlibwallet v1.6.1-rc1.0.20210915175038-31878a61e002
	github.com/planetdecred/dcrlibwallet/dexdcr v0.0.0-00010101000000-000000000000
	github.com/yeqown/go-qrcode v1.5.1
	golang.org/x/exp v0.0.0-20210722180016-6781d3edade3
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/text v0.3.6
)

replace (
	decred.org/dcrdex => github.com/itswisdomagain/dcrdex v0.0.0-20211004141752-92a02cc7352a
	github.com/planetdecred/dcrlibwallet => github.com/itswisdomagain/dcrlibwallet v1.0.0-rc1.0.20211004143133-f2d643fe4c9b
	github.com/planetdecred/dcrlibwallet/dexdcr => github.com/itswisdomagain/dcrlibwallet/dexdcr v0.0.0-20211004143133-f2d643fe4c9b
)
