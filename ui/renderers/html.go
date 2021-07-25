package renderers

import (
	"bytes"
	"fmt"
	"strings"

	"gioui.org/layout"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown/ast"
	//"github.com/gomarkdown/markdown/parser"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type HTMLProvider struct {
	container *layout.List
}

var (
	blockEls = []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li"}
)

const (
	openStyleTag      = "{@@"
	halfCloseStyleTag = "@}"
	closeStyleTag     = "{/@}"
)

func RenderHTML(html string, theme *decredmaterial.Theme) *HTMLProvider {
	htmlProvider := &HTMLProvider{}

	converter := md.NewConverter("", true, nil)
	docStr, err := converter.ConvertString(htmlProvider.prepare(html))
	if err != nil {
		fmt.Println(err)
		return &HTMLProvider{}
	}

	/**htmlProvider := &HTMLProvider{}

	converter := md.NewConverter("", true, nil)

	docStr, err := converter.ConvertString(docStr)
	if err != nil {
		fmt.Println(err)
		return r
	}

	return newRenderer(theme,  &HTMLProvider{})**/
	newNodeWalker(docStr, htmlProvider)

	return htmlProvider
}

func (p *HTMLProvider) prepareBlockQuote(node *ast.BlockQuote, entering bool) {

}

func (p *HTMLProvider) prepareList(node *ast.List, entering bool) {

}

func (p *HTMLProvider) prepareListItem(node *ast.ListItem, entering bool) {

}

func (p *HTMLProvider) prepareParagraph(node *ast.Paragraph, entering bool) {

}

func (p *HTMLProvider) prepareHeading(node *ast.Heading, entering bool) {

}

func (p *HTMLProvider) prepareStrong(node *ast.Strong, entering bool) {

}
func (p *HTMLProvider) prepareDel(node *ast.Del, entering bool) {

}
func (p *HTMLProvider) prepareEmph(node *ast.Emph, entering bool) {

}
func (p *HTMLProvider) prepareLink(node *ast.Link, entering bool) {

}

func (p *HTMLProvider) prepareHorizontalRule(node *ast.HorizontalRule, entering bool) {

}

func (p *HTMLProvider) prepareText(node *ast.Text, entering bool) {

}
func (p *HTMLProvider) prepareTable(node *ast.Table, entering bool) {

}
func (p *HTMLProvider) prepareTableCell(node *ast.TableCell, entering bool) {

}
func (p *HTMLProvider) prepareTableRow(node *ast.TableRow, entering bool) {

}

/**func RenderHTML(html string, theme *decredmaterial.Theme) *HTMLRenderer {
	converter := md.NewConverter("", true, nil)

	r := &HTMLRenderer{
		container: &layout.List{Axis: layout.Vertical},
		Renderer:  newRenderer(theme, true),
	}

	docStr := r.prepareHTML(html)

	docStr, err := converter.ConvertString(docStr)
	if err != nil {
		fmt.Println(err)
		return r
	}

	extensions := parser.NoIntraEmphasis        // Ignore emphasis markers inside words
	extensions |= parser.Tables                 // Parse tables
	extensions |= parser.FencedCode             // Parse fenced code blocks
	extensions |= parser.Autolink               // Detect embedded URLs that are not explicitly marked
	extensions |= parser.Strikethrough          // Strikethrough text using ~~test~~
	extensions |= parser.SpaceHeadings          // Be strict about prefix heading rules
	extensions |= parser.HeadingIDs             // specify heading IDs  with {#id}
	extensions |= parser.BackslashLineBreak     // Translate trailing backslashes into line breaks
	extensions |= parser.DefinitionLists        // Parse definition lists
	extensions |= parser.LaxHTMLBlocks          // more in HTMLBlock, less in HTMLSpan
	extensions |= parser.NoEmptyLineBeforeBlock // no need for new line before a list
	extensions |= parser.Attributes

	p := parser.NewWithExtensions(extensions)

	r.doc = p.Parse([]byte(docStr))
	r.parse()

	return r
}
**/

func (r *HTMLProvider) prepare(html string) string {
	//html = strings.Replace(html, "<br/>", " \n\n ", -1)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}

	doc.Find("*").Each(func(_ int, node *goquery.Selection) {
		nodeName := goquery.NodeName(node)
		switch nodeName {
		case "i":
			r.prepareItalic(node)
		case "em":
			r.prepareItalic(node)
		case "b", "strong":
			r.prepareBold(node)
		case "font":
			r.prepareFont(node)
		case "br":
			r.prepareBreak(node)
		}
	})

	doc.Find("body > *").Each(func(_ int, node *goquery.Selection) {
		styleMap := r.getStyleMap(node)
		newStyleMap := r.setNodeStyle(node, styleMap)
		r.traverse(node, newStyleMap)
	})

	return doc.Text()
}

func (r *HTMLProvider) prepareItalic(node *goquery.Selection) {
	style, ok := node.Attr("style")
	if ok {
		style += "; font-style: italic"
	} else {
		style = "font-style: italic"
	}

	node.ReplaceWithHtml(fmt.Sprintf(`<span style="%s">%s</span>`, style, node.Text()))
}

func (r *HTMLProvider) prepareBold(node *goquery.Selection) {
	style, ok := node.Attr("style")
	if ok {
		style += "; font-weight: bold"
	} else {
		style = "font-weight: bold"
	}

	node.ReplaceWithHtml(fmt.Sprintf(`<span style="%s">%s</span>`, style, node.Text()))
}

func (r *HTMLProvider) prepareFont(node *goquery.Selection) {
	style, _ := node.Attr("style")
	if style != "" {
		style += "; "
	}

	if color, ok := node.Attr("color"); ok {
		style += "text-color: " + color + "; "
	}

	if weight, ok := node.Attr("weight"); ok {
		style += "font-weight: " + weight + "; "
	}

	node.ReplaceWithHtml(fmt.Sprintf(`<span style="%s">%s</span>`, style, node.Text()))
}

func (r *HTMLProvider) prepareBreak(node *goquery.Selection) {
	node.ReplaceWithHtml("\n\n")
}

func (r *HTMLProvider) mapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func (r *HTMLProvider) getStyleMap(node *goquery.Selection) map[string]string {
	if styleStr, ok := node.Attr("style"); ok {
		spl := strings.Split(styleStr, ";")
		styleMap := map[string]string{}

		for _, v := range spl {
			items := strings.Split(v, ":")
			if len(items) == 2 {
				styleMap[strings.Trim(items[0], " ")] = strings.Trim(items[1], " ")
			}
		}

		return styleMap
	}

	return map[string]string{}
}

func (r *HTMLProvider) styleMapToString(m map[string]string) string {
	str := ""
	for k, v := range m {
		str += "##" + k + "--" + v
	}

	return str
}

func (r *HTMLProvider) traverse(node *goquery.Selection, parentStyle map[string]string) {
	node.Children().Each(func(_ int, s *goquery.Selection) {
		newStyle := r.setNodeStyle(s, parentStyle)
		r.traverse(s, newStyle)
	})
}

func (r *HTMLProvider) isBlockElement(element string) bool {
	for i := range blockEls {
		if element == blockEls[i] {
			return true
		}
	}

	return false
}

func (r *HTMLProvider) setNodeStyle(node *goquery.Selection, parentStyle map[string]string) map[string]string {
	styleMap := r.getStyleMap(node)
	for key, val := range parentStyle {
		if _, ok := styleMap[key]; !ok {
			styleMap[key] = val
		}
	}

	styleTag := openStyleTag + r.styleMapToString(styleMap) + halfCloseStyleTag
	endTag := closeStyleTag
	node = node.PrependHtml(styleTag)
	if r.isBlockElement(goquery.NodeName(node)) {
		endTag += " \n "
	}
	node.AppendHtml(endTag)

	return styleMap
}

func (r *HTMLProvider) parse() []byte {
	var buf bytes.Buffer
	/**ast.WalkFunc(r.doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return r.RenderNode(&buf, node, entering)
	})**/

	return buf.Bytes()
}

func (r *HTMLProvider) Layout(gtx C) D {
	/**w, _ := r.Renderer.Layout()

	return r.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return w[i](gtx)
	})**/
	return D{}
}
