package components

import (
	"image/color"
	"regexp"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const defaultScale = .7

var (
	doubleOrMoreDecimalPlaces = regexp.MustCompile(`(([0-9]{1,3},*)*\.)\d{2,}`)
	oneDecimalPlace           = regexp.MustCompile(`(([0-9]{1,3},*)*\.)\d`)
	noDecimal                 = regexp.MustCompile(`([0-9]{1,3},*)+`)
)

func formatBalance(gtx layout.Context, theme *decredmaterial.Theme, amount string, mainTextSize unit.Value, scale float32, col color.NRGBA) D {

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

	mainText, subText := amount[:startIndex], amount[startIndex:]

	subTextSize := unit.Sp(mainTextSize.V * scale)

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := theme.Label(mainTextSize, mainText)
			txt.Color = col
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := theme.Label(subTextSize, subText)
			txt.Color = col
			return txt.Layout(gtx)
		}),
	)
}

// LayoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func LayoutBalance(gtx layout.Context, theme *decredmaterial.Theme, amount string) layout.Dimensions {
	return formatBalance(gtx, theme, amount, values.TextSize20, defaultScale, theme.Color.Text)
}

func LayoutBalanceSize(gtx layout.Context, theme *decredmaterial.Theme, amount string, mainTextSize unit.Value) layout.Dimensions {
	return formatBalance(gtx, theme, amount, mainTextSize, defaultScale, theme.Color.Text)
}

func LayoutBalanceSizeScale(gtx layout.Context, theme *decredmaterial.Theme, amount string, mainTextSize unit.Value, scale float32) layout.Dimensions {
	return formatBalance(gtx, theme, amount, mainTextSize, scale, theme.Color.Text)
}

func LayoutBalanceColor(gtx layout.Context, theme *decredmaterial.Theme, amount string, color color.NRGBA) layout.Dimensions {
	return formatBalance(gtx, theme, amount, values.TextSize20, defaultScale, color)
}
