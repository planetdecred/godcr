package helper

import (
	"image"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const (
	logoPath              = "../../gio/assets/decred.png"
	logoSymbolPath        = "../../gio/assets/decred_symbol.png"
	overviewImagePath     = "../../gio/assets/overview.png"
	transactionsImagePath = "../../gio/assets/history.png"
	walletsImagePath      = "../../gio/assets/account.png"
	moreImagePath         = "../../gio/assets/more.png"
	sendImagePath         = "../../gio/assets/send.png"
	receiveImagePath      = "../../gio/assets/receive.png"
	infoImagePath         = "../../gio/assets/info.png"

	StandaloneScreenPadding = 20
)

var (
	logo material.Image

	LogoSymbol        material.Image
	OverviewImage     material.Image
	TransactionsImage material.Image
	WalletsImage      material.Image
	MoreImage         material.Image
	SendImage         material.Image
	ReceiveImage      material.Image
	InfoImage         material.Image
)

func LoadImage(theme *Theme, path string, scale float32) (material.Image, error) {
	b, err := os.Open(path)
	if err != nil {
		return material.Image{}, err
	}

	src, _, err := image.Decode(b)
	if err != nil {
		return material.Image{}, err
	}

	img := theme.Image(paint.NewImageOp(src))
	img.Scale = scale

	return img, nil
}

func InitImages(theme *Theme) error {
	var err error
	logo, err = LoadImage(theme, logoPath, 1.3)
	if err != nil {
		return err
	}

	LogoSymbol, err = LoadImage(theme, logoSymbolPath, 0.11)
	if err != nil {
		return err
	}

	OverviewImage, err = LoadImage(theme, overviewImagePath, 0.06)
	if err != nil {
		return err
	}

	TransactionsImage, err = LoadImage(theme, transactionsImagePath, 0.06)
	if err != nil {
		return err
	}

	WalletsImage, err = LoadImage(theme, walletsImagePath, 0.25)
	if err != nil {
		return err
	}

	MoreImage, err = LoadImage(theme, moreImagePath, 0.4)
	if err != nil {
		return err
	}

	SendImage, err = LoadImage(theme, sendImagePath, 0.06)
	if err != nil {
		return err
	}

	ReceiveImage, err = LoadImage(theme, receiveImagePath, 0.06)
	if err != nil {
		return err
	}

	InfoImage, err = LoadImage(theme, infoImagePath, 0.06)
	if err != nil {
		return err
	}

	return nil
}

func DrawLogo(ctx *layout.Context) {
	inset := layout.Inset{
		Left: unit.Dp(StandaloneScreenPadding),
	}
	inset.Layout(ctx, func() {
		logo.Layout(ctx)
	})
}
