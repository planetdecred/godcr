module github.com/planetdecred/godcr

go 1.16

require (
	gioui.org v0.0.0-20210418151603-3b69b5ed0512
	github.com/JohannesKaufmann/html-to-markdown v1.2.1
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/ararog/timeago v0.0.0-20160328174124-e9969cf18b8d
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/dcrutil/v2 v2.0.1
	github.com/decred/dcrd/dcrutil/v3 v3.0.0
	github.com/decred/slog v1.1.0
	github.com/gen2brain/beeep v0.0.0-20210529141713-5586760f0cc1
	github.com/godbus/dbus/v5 v5.0.5 // indirect
	github.com/gomarkdown/markdown v0.0.0-20210208175418-bda154fe17d8
	github.com/gopherjs/gopherjs v0.0.0-20210901121439-eee08aaf2717 // indirect
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/jrick/logrotate v1.0.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/planetdecred/dcrlibwallet v1.6.1-0.20210816165030-bb3af17a746a
	github.com/yeqown/go-qrcode v1.5.1
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/text v0.3.3
)

replace (
	github.com/decred/dcrdata/txhelpers/v4 => github.com/decred/dcrdata/txhelpers/v4 v4.0.0-20200108145420-f82113e7e212
	github.com/planetdecred/dcrlibwallet => github.com/C-ollins/mobilewallet v1.0.0-rc1.0.20210912175524-041481a23c8b
)
