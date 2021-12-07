package renderers

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown/ast"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type HTMLProvider struct {
	containers    []layout.Widget
	theme         *decredmaterial.Theme
	stringBuilder strings.Builder
	styleGroups   []map[string]string
	links         map[string]*widget.Clickable
	table         *table
	isList        bool
	prefix        string
}

var (
	blockEls = []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li"}
)

const (
	openStyleTag      = "{@@"
	halfCloseStyleTag = "@}"
	closeStyleTag     = "{/@}"
	linkTag           = "[[link"
	linkSpacer        = "@@@@"
)

func RenderHTML(html string, theme *decredmaterial.Theme) *HTMLProvider {
	htmlProvider := &HTMLProvider{
		theme: theme,
	}

	converter := md.NewConverter("", true, nil)
	docStr, err := converter.ConvertString(htmlProvider.prepare(html))
	if err != nil {
		fmt.Println(err)
		return &HTMLProvider{}
	}

	newNodeWalker(docStr, htmlProvider).walk()
	return htmlProvider
}

func (p *HTMLProvider) renderSoftBreak() {
	p.renderEmptyLine()
}

func (p *HTMLProvider) renderHardBreak() {
	p.renderEmptyLine()
}

func (p *HTMLProvider) prepareBlockQuote(node *ast.BlockQuote, entering bool) {}

func (p *HTMLProvider) prepareCode(node *ast.Code, entering bool) {}

func (p *HTMLProvider) prepareCodeBlock(node *ast.CodeBlock, entering bool) {}

func (p *HTMLProvider) prepareList(node *ast.List, entering bool) {
	if next := ast.GetNextNode(node); !entering && next != nil {
		_, parentIsListItem := node.GetParent().(*ast.ListItem)
		_, nextIsList := next.(*ast.List)
		if !nextIsList && !parentIsListItem {
			p.renderEmptyLine()
		}
	}
}

func (p *HTMLProvider) prepareListItem(node *ast.ListItem, entering bool) {
	if entering {
		p.isList = true
		switch {
		// numbered list
		case node.ListFlags&ast.ListTypeOrdered != 0:
			itemNumber := 1
			siblings := node.GetParent().GetChildren()
			for _, sibling := range siblings {
				if sibling == node {
					break
				}
				itemNumber++
			}
			p.prefix += fmt.Sprintf("%d. ", itemNumber)

		// content of a definition
		case node.ListFlags&ast.ListTypeDefinition != 0:
			p.prefix += " "

		// no flags means it's the normal bullet point list
		default:
			p.prefix += " " + bulletUnicode + " "
		}
	}
}

func (p *HTMLProvider) prepareParagraph(node *ast.Paragraph, entering bool) {
	if !entering {
		p.render(p.theme.Body1(""))
		p.renderEmptyLine()
	}
}

func (p *HTMLProvider) prepareHeading(node *ast.Heading, entering bool) {
	lblFunc := p.theme.H6

	switch node.Level {
	case 1:
		lblFunc = p.theme.H4
	case 2:
		lblFunc = p.theme.H5
	case 3:
		lblFunc = p.theme.H6
	}

	p.render(lblFunc(""))
	p.renderEmptyLine()
}

func (p *HTMLProvider) prepareStrong(node *ast.Strong, entering bool) {
	label := p.theme.Body1("")
	label.Font.Weight = text.Bold
	p.render(label)
}
func (p *HTMLProvider) prepareDel(node *ast.Del, entering bool) {}

// Will be taken care off by the render function
func (p *HTMLProvider) prepareEmph(node *ast.Emph, entering bool) {}

func (p *HTMLProvider) prepareLink(node *ast.Link, entering bool) {
	dest := string(node.Destination)
	text := string(ast.GetFirstChild(node).AsLeaf().Literal)

	if p.links == nil {
		p.links = map[string]*widget.Clickable{}
	}

	if _, ok := p.links[dest]; !ok {
		p.links[dest] = new(widget.Clickable)
	}

	// fix a bug that causes the link to be written to the builder before this is called
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()
	parts := strings.Split(content, " ")
	parts = parts[:len(parts)-1]
	for i := range parts {
		p.stringBuilder.WriteString(parts[i] + " ")
	}

	word := linkTag + linkSpacer + dest + linkSpacer + text
	p.stringBuilder.WriteString(word)
}

func (p *HTMLProvider) prepareHorizontalRule(node *ast.HorizontalRule, entering bool) {}

func (p *HTMLProvider) prepareText(node *ast.Text, entering bool) {
	if string(node.Literal) == "\n" {
		return
	}

	content := string(node.Literal)
	if shouldCleanText(node) {
		content = removeLineBreak(content)
	}
	p.stringBuilder.WriteString(content)
}
func (p *HTMLProvider) prepareTable(node *ast.Table, entering bool) {
	if entering {
		p.table = newTable(p.theme)
	} else {
		p.containers = append(p.containers, p.table.render())
		p.table = nil
	}
}
func (p *HTMLProvider) prepareTableCell(node *ast.TableCell, entering bool) {
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()

	align := cellAlignLeft
	switch node.Align {
	case ast.TableAlignmentRight:
		align = cellAlignRight
	case ast.TableAlignmentCenter:
		align = cellAlignCenter
	}

	if node.IsHeader {
		p.table.addCell(content, align, true)
	} else {
		p.table.addCell(content, cellAlignCopyHeader, false)
	}
}
func (p *HTMLProvider) prepareTableRow(node *ast.TableRow, entering bool) {
	if _, ok := node.Parent.(*ast.TableBody); ok && entering {
		p.table.startNextRow()
	}
	if _, ok := node.Parent.(*ast.TableFooter); ok && entering {
		p.table.startNextRow()
	}
}

func (p *HTMLProvider) render(lbl decredmaterial.Label) {
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()

	if p.prefix != "" {
		content = p.prefix + " " + content
	}

	var labels []decredmaterial.Label
	var inStyleBlock bool
	var isClosingStyle bool
	var isClosingBlock bool
	var currStyle string
	var currText string
	for i := range content {
		curr := content[i]

		if curr == openStyleTag[0] && getNextChar(content, i) == openStyleTag[1] {
			inStyleBlock = true
			labels = append(labels, p.getLabel(lbl, currText))
			currText = ""
		}

		if curr == halfCloseStyleTag[0] && getNextChar(content, i) == halfCloseStyleTag[1] {
			isClosingStyle = true
		}

		if curr == closeStyleTag[0] && getNextChar(content, i) == closeStyleTag[1] {
			isClosingBlock = true
		}

		if !inStyleBlock && !isClosingBlock {
			currStr := string(curr)
			currText += currStr

			if i+1 == len(content) || currStr == "" || currStr == " " {
				labels = append(labels, p.getLabel(lbl, currText))
				currText = ""
			}
		}

		if isClosingBlock && curr == closeStyleTag[3] {
			labels = append(labels, p.getLabel(lbl, currText))
			currText = ""
			p.removeLastStyleGroup()
			isClosingBlock = false

		}

		if inStyleBlock && !isClosingStyle {
			currStyle += string(curr)
		}

		if isClosingStyle && curr == halfCloseStyleTag[1] {
			isClosingStyle = false
			inStyleBlock = false
			p.addStyleGroup(currStyle)
			currStyle = ""
		}
	}

	wdgt := func(gtx C) D {
		return decredmaterial.GridWrap{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		}.Layout(gtx, len(labels), func(gtx C, i int) D {
			if strings.Trim(labels[i].Text, " ") == "" {
				return D{}
			}
			return labels[i].Layout(gtx)
		})
	}
	p.containers = append(p.containers, wdgt)
}

func (p *HTMLProvider) getLabel(lbl decredmaterial.Label, text string) decredmaterial.Label {
	l := lbl
	l.Text = text
	l = p.styleLabel(l)
	return l
}

func (p *HTMLProvider) removeLastStyleGroup() {
	if len(p.styleGroups) > 0 {
		p.styleGroups = p.styleGroups[:len(p.styleGroups)-1]
	}
}

func (p *HTMLProvider) addStyleGroup(str string) {
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
		p.styleGroups = append(p.styleGroups, styleMap)
	}
}

func (p *HTMLProvider) styleLabel(label decredmaterial.Label) decredmaterial.Label {
	if len(p.styleGroups) == 0 {
		return label
	}

	style := p.styleGroups[len(p.styleGroups)-1]
	label.Font.Weight = p.getLabelWeight(style["font-weight"])

	colStr := style["text-color"]
	if colStr == "" {
		colStr = style["color"]
	}

	if col, ok := parseColorCode(colStr); ok {
		label.Color = col
	} else {
		label.Color = p.getColorFromMap(colStr)
	}

	if fontStyle, ok := style["font-style"]; ok {
		if fontStyle == "italic" {
			label.Font.Style = text.Italic
		}
	}

	return label
}

func (p *HTMLProvider) getLabelWeight(weight string) text.Weight {
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

func (p *HTMLProvider) getColorFromMap(col string) color.NRGBA {
	colorMap := map[string]color.NRGBA{
		"primary":    p.theme.Color.Primary,
		"text":       p.theme.Color.Text,
		"grayText1":  p.theme.Color.GrayText1,
		"grayText2":  p.theme.Color.GrayText2,
		"grayText3":  p.theme.Color.GrayText3,
		"grayText4":  p.theme.Color.GrayText4,
		"greenText":  p.theme.Color.GreenText,
		"inv-text":   p.theme.Color.InvText,
		"success":    p.theme.Color.Success,
		"success2":   p.theme.Color.Success2,
		"danger":     p.theme.Color.Danger,
		"surface":    p.theme.Color.Surface,
		"black":      p.theme.Color.Black,
		"light-blue": p.theme.Color.LightBlue,
		"orange":     p.theme.Color.Orange,
		"orange2":    p.theme.Color.Orange2,
	}

	if color, ok := colorMap[col]; ok {
		return color
	}

	return colorMap["text"]
}

func (p *HTMLProvider) renderEmptyLine() {
	var padding = -5

	if p.isList {
		padding = -10
		p.isList = false
	}

	p.containers = append(p.containers, func(gtx C) D {
		dims := p.theme.Body2("").Layout(gtx)
		dims.Size.Y = dims.Size.Y + padding
		return dims
	})
}

func (p *HTMLProvider) prepare(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}

	doc.Find("*").Each(func(_ int, node *goquery.Selection) {
		nodeName := goquery.NodeName(node)
		switch nodeName {
		case "i":
			p.prepareItalic(node)
		case "em":
			p.prepareItalic(node)
		case "b", "strong":
			p.prepareBold(node)
		case "font":
			p.prepareFont(node)
		case "br":
			p.prepareBreak(node)
		}
	})

	doc.Find("body > *").Each(func(_ int, node *goquery.Selection) {
		styleMap := p.getStyleMap(node)
		newStyleMap := p.setNodeStyle(node, styleMap)
		p.traverse(node, newStyleMap)
	})

	return doc.Text()
}

func (p *HTMLProvider) prepareItalic(node *goquery.Selection) {
	style, ok := node.Attr("style")
	if ok {
		style += "; font-style: italic"
	} else {
		style = "font-style: italic"
	}

	node.ReplaceWithHtml(fmt.Sprintf(`<span style="%s">%s</span>`, style, node.Text()))
}

func (p *HTMLProvider) prepareBold(node *goquery.Selection) {
	style, ok := node.Attr("style")
	if ok {
		style += "; font-weight: bold"
	} else {
		style = "font-weight: bold"
	}

	node.ReplaceWithHtml(fmt.Sprintf(`<span style="%s">%s</span>`, style, node.Text()))
}

func (p *HTMLProvider) prepareFont(node *goquery.Selection) {
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

func (p *HTMLProvider) prepareBreak(node *goquery.Selection) {
	node.ReplaceWithHtml("\n\n")
}

func (p *HTMLProvider) mapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func (p *HTMLProvider) getStyleMap(node *goquery.Selection) map[string]string {
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

func (p *HTMLProvider) styleMapToString(m map[string]string) string {
	str := ""
	for k, v := range m {
		str += "##" + k + "--" + v
	}

	return str
}

func (p *HTMLProvider) traverse(node *goquery.Selection, parentStyle map[string]string) {
	node.Children().Each(func(_ int, s *goquery.Selection) {
		newStyle := p.setNodeStyle(s, parentStyle)
		p.traverse(s, newStyle)
	})
}

func (p *HTMLProvider) isBlockElement(element string) bool {
	for i := range blockEls {
		if element == blockEls[i] {
			return true
		}
	}

	return false
}

func (p *HTMLProvider) setNodeStyle(node *goquery.Selection, parentStyle map[string]string) map[string]string {
	styleMap := p.getStyleMap(node)
	for key, val := range parentStyle {
		if _, ok := styleMap[key]; !ok {
			styleMap[key] = val
		}
	}

	styleTag := openStyleTag + p.styleMapToString(styleMap) + halfCloseStyleTag
	endTag := closeStyleTag
	node = node.PrependHtml(styleTag)
	if p.isBlockElement(goquery.NodeName(node)) {
		endTag += " \n "
	}
	node.AppendHtml(endTag)

	return styleMap
}

func (p *HTMLProvider) Layout(gtx C) D {
	return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(p.containers), func(gtx C, i int) D {
		return p.containers[i](gtx)
	})
}
