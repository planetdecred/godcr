package components

import (
	"image/color"
	"regexp"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const defaultScale = .7

var (
	doubleOrMoreDecimalPlaces = regexp.MustCompile(`(([0-9]{1,3},*)*\.)\d{2,}`)
	oneDecimalPlace           = regexp.MustCompile(`(([0-9]{1,3},*)*\.)\d`)
	noDecimal                 = regexp.MustCompile(`([0-9]{1,3},*)+`)
)

func formatBalance(gtx layout.Context, l *load.Load, amount string, mainTextSize unit.Sp, scale float32, col color.NRGBA, withUnit bool) D {

	startIndex := 0

	if doubleOrMoreDecimalPlaces.MatchString(amount) {
		decimalIndex := strings.Index(amount, ".")
		startIndex = decimalIndex + 3
	} else if oneDecimalPlace.MatchString(amount) {
		decimalIndex := strings.Index(amount, ".")
		startIndex = decimalIndex + 2
	} else if noDecimal.MatchString(amount) {
		loc := noDecimal.FindStringIndex(amount)
		startIndex = loc[1] // start scaling from the end
	}

	indexUnit := len(amount) - 4
	if !withUnit {
		indexUnit = len(amount) - 1
	}

	mainText, subText, unitValue := amount[:startIndex], amount[startIndex:indexUnit], amount[indexUnit:]

	subTextSize := unit.Sp(float32(mainTextSize) * scale)

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := l.Theme.Label(mainTextSize, mainText)
			txt.Color = col
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := l.Theme.Label(subTextSize, subText)
			txt.Color = col
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := l.Theme.Label(mainTextSize, unitValue)
			txt.Color = col
			if withUnit {
				return txt.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

// LayoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func LayoutBalance(gtx layout.Context, l *load.Load, amount string) layout.Dimensions {
	return formatBalance(gtx, l, amount, values.TextSize20, defaultScale, l.Theme.Color.Text, false)
}

func LayoutBalanceWithUnit(gtx layout.Context, l *load.Load, amount string) layout.Dimensions {
	return formatBalance(gtx, l, amount, values.TextSize20, defaultScale, l.Theme.Color.Text, true)
}

func LayoutBalanceSize(gtx layout.Context, l *load.Load, amount string, mainTextSize unit.Sp) layout.Dimensions {
	return formatBalance(gtx, l, amount, mainTextSize, defaultScale, l.Theme.Color.Text, false)
}

func LayoutBalanceSizeScale(gtx layout.Context, l *load.Load, amount string, mainTextSize unit.Sp, scale float32) layout.Dimensions {
	return formatBalance(gtx, l, amount, mainTextSize, scale, l.Theme.Color.Text, false)
}

func LayoutBalanceColor(gtx layout.Context, l *load.Load, amount string, color color.NRGBA) layout.Dimensions {
	return formatBalance(gtx, l, amount, values.TextSize20, defaultScale, color, false)
}
