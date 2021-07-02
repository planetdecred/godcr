package renderers

import (
	"image/color"
	"strings"

	"gioui.org/text"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

func (r *Renderer) getLabelWeight(weight string) text.Weight {
	switch weight {
	case "normal":
		return text.Normal
	case "medium":
		return text.Medium
	case "bold":
		return text.Bold
	}

	return text.Normal
}

func (r *Renderer) getColorFromMap(col string) color.NRGBA {
	colorMap := map[string]color.NRGBA{
		"primary":       r.theme.Color.Primary,
		"secondary":     r.theme.Color.Secondary,
		"text":          r.theme.Color.Text,
		"hint":          r.theme.Color.Hint,
		"overlay":       r.theme.Color.Overlay,
		"inv-text":      r.theme.Color.InvText,
		"success":       r.theme.Color.Success,
		"success2":      r.theme.Color.Success2,
		"danger":        r.theme.Color.Danger,
		"background":    r.theme.Color.Background,
		"surface":       r.theme.Color.Surface,
		"gray":          r.theme.Color.Gray,
		"black":         r.theme.Color.Black,
		"deep-blue":     r.theme.Color.DeepBlue,
		"light-blue":    r.theme.Color.LightBlue,
		"light-gray":    r.theme.Color.LightGray,
		"inactive-gray": r.theme.Color.InactiveGray,
		"active-gray":   r.theme.Color.ActiveGray,
		"gray1":         r.theme.Color.Gray1,
		"gray2":         r.theme.Color.Gray2,
		"gray3":         r.theme.Color.Gray3,
		"gray4":         r.theme.Color.Gray4,
		"gray5":         r.theme.Color.Gray5,
		"gray6":         r.theme.Color.Gray6,
		"orange":        r.theme.Color.Orange,
		"orange2":       r.theme.Color.Orange2,
	}

	if color, ok := colorMap[col]; ok {
		return color
	}

	return colorMap["text"]
}

func (r *Renderer) styleLabel(label decredmaterial.Label) decredmaterial.Label {
	if len(r.styleGroups) == 0 {
		return label
	}

	style := r.styleGroups[len(r.styleGroups)-1]
	label.Font.Weight = r.getLabelWeight(style["font-weight"])

	colStr := style["text-color"]
	if colStr == "" {
		colStr = style["color"]
	}

	if col, ok := parseColorCode(colStr); ok {
		label.Color = col
	} else {
		label.Color = r.getColorFromMap(colStr)
	}

	if fontStyle, ok := style["font-style"]; ok {
		if fontStyle == "italic" {
			label.Font.Style = text.Italic
		}
	}

	return label
}

func (r *Renderer) addStyleGroup(str string) {
	parts := strings.Split(str, "##")
	styleMap := map[string]string{}

	for i := range parts {
		if parts[i] != " " && parts[i] != "{" {
			styleParts := strings.Split(parts[i], "--")

			if len(styleParts) == 2 {
				styleMap[styleParts[0]] = styleParts[1]
			}
		}
	}

	if len(styleMap) > 0 {
		r.styleGroups = append(r.styleGroups, styleMap)
	}
}

func (r *Renderer) removeLastStyleGroup() {
	if len(r.styleGroups) > 0 {
		r.styleGroups = r.styleGroups[:len(r.styleGroups)-1]
	}
}
