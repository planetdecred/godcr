package renderers

import (
	//"fmt"
	//"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

func setStyle(lbl decredmaterial.Label, style string) decredmaterial.Label {
	switch style {
	case italicsTagName, emphTagName:
		lbl.Font.Style = text.Italic
	}

	return lbl
}

func setWeight(lbl decredmaterial.Label, weight string) decredmaterial.Label {
	var w text.Weight

	switch weight {
	case "normal":
		w = text.Normal
	case "medium":
		w = text.Medium
	case "bold", "strong":
		w = text.Bold
	default:
		w = lbl.Font.Weight
	}

	lbl.Font.Weight = w
	return lbl
}

func getHeading(txt string, tagName string, theme *decredmaterial.Theme) decredmaterial.Label {
	var lblWdgt func(string) decredmaterial.Label

	switch tagName {
	case h5TagName:
		lblWdgt = theme.H5
	case h6TagName:
		lblWdgt = theme.H6
	default:
		lblWdgt = theme.H4
	}

	return lblWdgt(txt)
}

func renderStrike(lbl decredmaterial.Label, theme *decredmaterial.Theme) layout.Widget {
	return func(gtx C) D {
		var dims D
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				dims = lbl.Layout(gtx)
				return dims
			}),
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: unit.Dp((float32(dims.Size.Y) / float32(2))),
				}.Layout(gtx, func(gtx C) D {
					l := theme.Separator()
					l.Color = lbl.Color
					l.Width = dims.Size.X
					return l.Layout(gtx)
				})
			}),
		)
	}
}

func renderBlockQuote(lbl decredmaterial.Label, theme *decredmaterial.Theme) layout.Widget {
	words := strings.Fields(lbl.Text)

	return func(gtx C) D {
		var dims D

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				l := theme.SeparatorVertical(dims.Size.Y, 10)
				l.Color = theme.Color.Gray
				return l.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				dims = layout.Inset{
					Left: unit.Dp(4),
				}.Layout(gtx, func(gtx C) D {
					return decredmaterial.GridWrap{
						Axis:      layout.Horizontal,
						Alignment: layout.Start,
					}.Layout(gtx, len(words), func(gtx C, i int) D {
						lbl.Text = words[i] + " "
						return lbl.Layout(gtx)
					})
				})

				return dims
			}),
		)
	}
}

func (p *MarkdownProvider) getLinkWidget(linkWord string) layout.Widget {
	parts := strings.Split(linkWord, linkSpacer)

	return func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Constraints.Max.X - 200
		return material.Clickable(gtx, p.links[parts[1]], func(gtx C) D {
			lbl := p.theme.Body2(strings.Replace(parts[2], "---", " ", -1) + " ")
			lbl.Color = p.theme.Color.Primary
			return lbl.Layout(gtx)
		})
	}
}

/**func (p *MarkdownProvider) getLinkWidget(gtx layout.Context, linkWord string) D {
	parts := strings.Split(linkWord, linkSpacer)

	gtx.Constraints.Max.X = gtx.Constraints.Max.X - 200
	return material.Clickable(gtx, p.links[parts[1]], func(gtx C) D {
		lbl := p.theme.Body2(parts[2] + " ")
		lbl.Color = p.theme.Color.Primary
		return lbl.Layout(gtx)
	})
}
**/

func renderHorizontalLine(theme *decredmaterial.Theme) layout.Widget {
	l := theme.Separator()
	l.Width = 1
	return l.Layout
}

func renderEmptyLine(theme *decredmaterial.Theme, isList bool) layout.Widget {
	var padding = -5

	if isList {
		padding = -10
	}

	return func(gtx C) D {
		dims := theme.Body2("").Layout(gtx)
		dims.Size.Y = dims.Size.Y + padding
		return dims
	}
}

func renderListItem(lbl decredmaterial.Label, theme *decredmaterial.Theme) layout.Widget {
	words := strings.Fields(lbl.Text)
	if len(words) == 0 {
		return func(gtx C) D { return D{} }
	}

	return func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lbl.Text = words[0]
				if len(words) > 1 {
					words = words[1:]
				}
				return lbl.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				return decredmaterial.GridWrap{
					Axis:      layout.Horizontal,
					Alignment: layout.Start,
				}.Layout(gtx, len(words), func(gtx C, i int) D {
					lbl.Text = words[i] + " "
					return lbl.Layout(gtx)
				})
			}),
		)
	}
}

/**
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

func (r *Renderer) styleHTMLLabel(label decredmaterial.Label) decredmaterial.Label {
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

func (r *Renderer) setLabelStyle(label decredmaterial.Label, value string) decredmaterial.Label {
	if value == "italic" {
		label.Font.Style = text.Italic
	}

	return label
}

func (r *Renderer) setLabelWeight(label decredmaterial.Label, value string) decredmaterial.Label {
	switch value {
	case "bold":
		label.Font.Weight = text.Bold
	case "normal":
		label.Font.Weight = text.Normal
	}

	return label
}

func (r *Renderer) styleMarkdownLabel(label decredmaterial.Label) layout.Widget {
	var wdgt layout.Widget

	 for i := range r.styleGroups {
		for style, value := range r.styleGroups[i] {
			switch style {
			case "font-style":
				wdgt = r.setLabelStyle(label, value).Layout
			case "font-weight":
				wdgt = r.setLabelWeight(label, value).Layout
			case "font-decoration":
				if value == strikeTagName {
					wdgt = r.strikeLabel(label)
				}
			}
		}
	}

	if wdgt == nil {
		return label.Layout
	}
	return wdgt
}

func (r *Renderer) getMarkdownWidgetAndStyle(label decredmaterial.Label) layout.Widget {
	return r.styleMarkdownLabel(label)
}

func (r *Renderer) addHTMLStyleGroup(str string) {
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

func (r *Renderer) addStyleItem(style, value string) {
	styleMap := map[string]string{
		style: value,
	}

	r.styleGroups = append(r.styleGroups, styleMap)
}

func (r *Renderer) addStyleGroupFromTagName(tagName string) {
	var key, val string

	switch tagName {
	case italicsTagName:
		key, val = "font-style", "italic"
	case emphTagName:
		key, val = "font-style", "italic"
	case strongTagName:
		key, val = "font-weight", "bold"
	case strikeTagName:
		key, val = "font-decoration", "strike"
	}

	if key != "" && val != "" {
		r.addStyleItem(key, val)
	}
}

func (r *Renderer) removeLastStyleGroup() {
	if len(r.styleGroups) > 0 {
		r.styleGroups = r.styleGroups[:len(r.styleGroups)-1]
	}
}
**/
