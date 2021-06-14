package renderers

import (
	"fmt"
	"image"
	"image/color"
	"os/exec"
	"runtime"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/PuerkitoBio/goquery"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type tagItem struct {
	spaceBelow int
	isLink     bool
	label      decredmaterial.Label
}

type HTMLRenderer struct {
	theme *decredmaterial.Theme
	doc   *goquery.Document
	items []tagItem
	//containers []layout.Widget
	links map[string]*widget.Clickable
}

const (
	pSpacing = 20
	divSpacing = 5
)

func RenderHTML(text string, theme *decredmaterial.Theme) *HTMLRenderer {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		panic(err)
	}

	renderer := &HTMLRenderer{
		doc:   doc,
		theme: theme,
	}
	renderer.render()
	return renderer
}

func (r *HTMLRenderer) render() {
	r.doc.Find("*").Each(func(_ int, node *goquery.Selection) {
		nodeName := goquery.NodeName(node)
		switch nodeName {
		case "html", "head", "body":
		case "a":
			r.renderLink(node)
		case "p":
			r.renderParagraph(node)
		case "strong":
			r.renderStrong(node)
		case "h1", "h2", "h3":
			r.renderHeading(node, nodeName)
		case "span":
			r.renderSpan(node)
		default:
			r.renderText(node)
		}
	})
}

func (r *HTMLRenderer) getEmptyLine() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		dims := r.theme.Body2("").Layout(gtx)
		dims.Size.Y -= 40

		return dims
	}
}

func (r *HTMLRenderer) getStyleMap(node *goquery.Selection) map[string]string {
	if styleStr, ok := node.Attr("style"); ok {
		spl := strings.Split(styleStr, ";")
		styleMap := map[string]string{}

		for _, v := range spl {
			items := strings.Split(v, ":")
			styleMap[strings.Trim(items[0], " ")] = strings.Trim(items[1], " ")
		}

		if labelType, ok := node.Attr("label-type"); ok {
			styleMap["label-type"] = labelType
		} else {
			styleMap["label-type"] = "body1"
		}

		return styleMap
	}

	return map[string]string{}
}

func (r *HTMLRenderer) renderParagraph(node *goquery.Selection) {
	//r.renderWords(node)
	//r.containers = append(r.containers, r.getEmptyLine())
	r.addTag(node, pSpacing)
}

func (r *HTMLRenderer) renderHeading(node *goquery.Selection, nodeName string) {
	node.SetAttr("label-type", nodeName)
	r.renderWords(node)
}

func (r *HTMLRenderer) renderStrong(node *goquery.Selection) {
	r.renderWords(node)
}

func (r *HTMLRenderer) renderSpan(node *goquery.Selection) {
	r.addTag(node, 0)
}

func (r *HTMLRenderer) renderText(node *goquery.Selection) {
	r.addTag(node, 0)
}

func (r *HTMLRenderer) renderLink(node *goquery.Selection) {
	/**href, ok := node.Attr("href")
	if ok {
		if _, ok := r.links[href]; !ok {
			r.links[href] = new(widget.Clickable)
		}
	}


	words, label := r.getWordsAndLabel(node)
	wdgt := func(gtx C) D {
		return decredmaterial.GridWrap{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		}.Layout(gtx, len(words), func(gtx C, i int) D {
			return decredmaterial.Clickable(gtx, r.links[href], func(gtx C) D {
				label.Text = words[i] + " "
				return label.Layout(gtx)
			})
		})
	}
	r.containers = append(r.containers, wdgt)**/
}

func (r *HTMLRenderer) getLabel(lblType string) decredmaterial.Label {
	labelMap := map[string]func(string) decredmaterial.Label{
		"body1":   r.theme.Body1,
		"body2":   r.theme.Body2,
		"caption": r.theme.Caption,
		"h1":      r.theme.H1,
		"h2":      r.theme.H2,
		"h3":      r.theme.H3,
		"h4":      r.theme.H4,
		"h5":      r.theme.H5,
		"h6":      r.theme.H6,
	}

	if label, ok := labelMap[lblType]; ok {
		return label("")
	}

	return labelMap["body1"]("")
}

func (r *HTMLRenderer) getColor(col string) color.NRGBA {
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

func (r *HTMLRenderer) getLabelAndStyle(style map[string]string) decredmaterial.Label {
	label := r.getLabel(style["label-type"])
	label.Color = r.getColor(style["text-color"])

	return label
}

func (r *HTMLRenderer) getWordsAndLabel(node *goquery.Selection) ([]string, decredmaterial.Label) {
	style := r.getStyleMap(node)
	content := strings.TrimSpace(node.Text())
	words := strings.Split(content, " ")
	label := r.getLabelAndStyle(style)

	return words, label
}

func (r *HTMLRenderer) addTag(node *goquery.Selection, spaceBelow int) {
	words, label := r.getWordsAndLabel(node)
	for i := range words {
		label.Text += words[i] + " "
	}

	item := tagItem{
		spaceBelow: spaceBelow,
		label:      label,
	}
	r.items = append(r.items, item)
}

func (r *HTMLRenderer) renderWords(node *goquery.Selection) {

}

func (r *HTMLRenderer) Layout(gtx C) D {
	max := gtx.Constraints.Max.X

	return decredmaterial.GridWrap{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}.Layout(gtx, len(r.items), func(gtx C, i int) D {
		if r.items[i].spaceBelow == 0 {
			return r.items[i].label.Layout(gtx)
		}

		gtx.Constraints.Min.X = max
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, 
			layout.Rigid(func(gtx C) D {
				return D{
					Size: image.Point{
						Y: r.items[i].spaceBelow,
					},
				}
			}),
			layout.Rigid(r.items[i].label.Layout),
		)
	})
}

/**func (r *HTMLRenderer) renderWords(node *goquery.Selection) {
	words, label := r.getWordsAndLabel(node)

	wdgt := func(gtx C) D {
		return decredmaterial.GridWrap{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		}.Layout(gtx, len(words), func(gtx C, i int) D {
			label.Text = words[i] + " "
			return label.Layout(gtx)
		})
	}
	r.containers = append(r.containers, wdgt)
}**/

/**func (r *HTMLRenderer) Layout() []layout.Widget {
	for href, clickable := range r.links {
		for clickable.Clicked() {
			goToURL(href)
		}
	}

	return r.containers
}**/

func goToURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}