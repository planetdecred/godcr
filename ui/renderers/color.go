package renderers

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

const (
	hexRegexString = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	hexFormat      = "#%02x%02x%02x"
	hexShortFormat = "#%1x%1x%1x"
	hexToRGBFactor = 17

	rgbString       = "rgb(%d,%d,%d)"
	rgbRegexString  = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	rgbaString      = "rgba(%d,%d,%d,%g)"
	rgbaRegexString = "^rgba\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0\\.[0-9]*|[01])\\s*\\)$"
)

var (
	hexRegex  = regexp.MustCompile(hexRegexString)
	rgbRegex  = regexp.MustCompile(rgbRegexString)
	rgbaRegex = regexp.MustCompile(rgbaRegexString)

	errBadHexCode = errors.New("Bad hex color code")
	errBadRGBCode = errors.New("Bad rgb color code")
)

func parseColorCode(colorCode string) (color.NRGBA, bool) {
	colorCode = strings.ToLower(colorCode)

	if col, ok := parseHex(colorCode); ok {
		return col, ok
	}

	if col, ok := parseRGB(colorCode); ok {
		return col, ok
	}

	if col, ok := parseRGBA(colorCode); ok {
		return col, ok
	}

	return color.NRGBA{}, false
}

func parseHex(colorCode string) (color.NRGBA, bool) {
	if !hexRegex.MatchString(colorCode) {
		return color.NRGBA{}, false
	}

	var r, g, b uint8
	if len(colorCode) == 4 {
		fmt.Sscanf(colorCode, hexShortFormat, &r, &g, &b)
		r *= hexToRGBFactor
		g *= hexToRGBFactor
		b *= hexToRGBFactor
	} else {
		fmt.Sscanf(colorCode, hexFormat, &r, &g, &b)
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}, true
}

func parseRGB(colorCode string) (color.NRGBA, bool) {
	parts := rgbRegex.FindAllStringSubmatch(colorCode, -1)
	if len(parts) == 0 || len(parts[0]) == 0 {
		return color.NRGBA{}, false
	}

	r, _ := strconv.ParseUint(parts[0][1], 10, 8)
	g, _ := strconv.ParseUint(parts[0][2], 10, 8)
	b, _ := strconv.ParseUint(parts[0][3], 10, 8)

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, true
}

func parseRGBA(colorCode string) (color.NRGBA, bool) {
	parts := rgbaRegex.FindAllStringSubmatch(colorCode, -1)
	if len(parts) == 0 || len(parts[0]) == 0 {
		return color.NRGBA{}, false
	}

	r, _ := strconv.ParseUint(parts[0][1], 10, 8)
	g, _ := strconv.ParseUint(parts[0][2], 10, 8)
	b, _ := strconv.ParseUint(parts[0][3], 10, 8)
	a, _ := strconv.ParseFloat(parts[0][4], 64)

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, true
}
