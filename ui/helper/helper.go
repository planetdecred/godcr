
package helper

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/wallet"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

// PaintArea paints an area with the given color and dimensions
func PaintArea(gtx *layout.Context, col color.RGBA, x, y int) {
	dim := image.Point{
		X: x,
		Y: y,
	}

	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(dim.X),
			Y: float32(dim.Y),
		},
	}

	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: rect}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: dim}
}

// Fill paints an area completely
func Fill(gtx *layout.Context, col color.RGBA) {
	PaintArea(gtx, col, gtx.Constraints.Width.Min, gtx.Constraints.Height.Min)
}
// breakBalance takes the balance string and returns it in two slices
func BreakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
	b1 = balanceParts[0]
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// RemainingSyncTime takes time on int64 and returns its string equivalent.
func RemainingSyncTime(totalTimeLeft int64) string {
	var days, hours, minutes, seconds int64

	q, r := divMod(totalTimeLeft, 24*60*60)
	days = q
	totalTimeLeft = r
	q, r = divMod(totalTimeLeft, 60*60)
	hours = q
	totalTimeLeft = r
	q, r = divMod(totalTimeLeft, 60)
	minutes = q
	totalTimeLeft = r
	seconds = totalTimeLeft
	if days > 0 {
		return fmt.Sprintf("%d"+"d"+"%d"+"h"+"%d"+"m"+"%d"+"s", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%d"+"h"+"%d"+"m"+"%d"+"s", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d"+"m"+"%d"+"s", minutes, seconds)
	}
	return fmt.Sprintf("%d"+"s", seconds)
}

func WalletSyncStatus(info wallet.InfoShort, bestBlockHeight int32) string {
	if info.IsWaiting {
		return "waiting for other wallets"
	}
	if info.BestBlockHeight < bestBlockHeight {
		return "syncing.."
	}

	return "synced"
}

func WalletSyncProgressTime(timestamp int64) string {
	return fmt.Sprintf("%s behind", dcrlibwallet.CalculateDaysBehind(timestamp))
}

func LastBlockSync(timestamp int64) string {
	return truncateTime(time.Since(time.Unix(timestamp, 0)).String(), 0)
}

func truncateTime (duration string, place int) string {
	var durationCharacter string
	durationSlice := strings.Split(duration, ".")
	if len(durationSlice) == 1 {
		return duration
	}

	secondsDecimals := durationSlice[1]
	if place > len(secondsDecimals) {
		return duration
	}

	secondLastCharacter := secondsDecimals[len(secondsDecimals) - 2:len(secondsDecimals) - 1]
	_, err := strconv.Atoi(secondLastCharacter)
	if err != nil {
		durationCharacter = secondsDecimals[len(secondsDecimals) - 2:]
	} else {
		durationCharacter = secondsDecimals[len(secondsDecimals) - 1:]
	}
	if place == 0 {
		return durationSlice[0] + durationCharacter
	}
	return durationSlice[0] + "." + secondsDecimals[0:place] + durationCharacter
}