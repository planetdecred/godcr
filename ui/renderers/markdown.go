package renderers

import (
	"fmt"
	"strings"
	"unicode"

	"gioui.org/layout"
	"gioui.org/widget"

	//md "github.com/gomarkdown/markdown"
	//"github.com/gomarkdown/markdown/parser"
	"github.com/gomarkdown/markdown/ast"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const (
	bulletUnicode = "\u2022"
	linkTag       = "[[link"
	linkSpacer    = "@@@@"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type MarkdownProvider struct {
	containers     []layout.Widget
	theme          *decredmaterial.Theme
	listItemNumber int // should be negative when not rendering a list
	links          map[string]*widget.Clickable
	table          *table

	stringBuilder strings.Builder
	tagStack      []string
}

func RenderMarkdown(gtx C, theme *decredmaterial.Theme, source string) *MarkdownProvider {
	mdProvider := &MarkdownProvider{
		theme:          theme,
		listItemNumber: -1,
	}
	source = mdProvider.prepare(source)

	newNodeWalker(source, mdProvider).walk()
	return mdProvider
}

func (*MarkdownProvider) prepare(doc string) string {
	d := strings.Replace(doc, ":|", "------:|", -1)
	d = strings.Replace(d, "-|", "------|", -1)
	d = strings.Replace(d, "|-", "|------", -1)
	d = strings.Replace(d, "|:-", "|:------", -1)

	return d
}

func (m *MarkdownProvider) Layout() ([]layout.Widget, map[string]*widget.Clickable) {
	return m.containers, map[string]*widget.Clickable{}
}

func (p *MarkdownProvider) prepareBlockQuote(node *ast.BlockQuote, entering bool) {
	p.openOrCloseTag(blockQuoteTagName, entering)
}

func (p *MarkdownProvider) prepareStrong(node *ast.Strong, entering bool) {
	p.openOrCloseTag(strongTagName, entering)
}

func (p *MarkdownProvider) prepareDel(node *ast.Del, entering bool) {
	p.openOrCloseTag(strikeTagName, entering)
}

func (p *MarkdownProvider) prepareEmph(node *ast.Emph, entering bool) {
	p.openOrCloseTag(emphTagName, entering)
}

func (p *MarkdownProvider) prepareHorizontalRule(node *ast.HorizontalRule, entering bool) {
	p.containers = append(p.containers, renderHorizontalLine(p.theme))
}

func (p *MarkdownProvider) prepareList(node *ast.List, entering bool) {
	if next := ast.GetNextNode(node); !entering && next != nil {
		_, parentIsListItem := node.GetParent().(*ast.ListItem)
		_, nextIsList := next.(*ast.List)
		if !nextIsList && !parentIsListItem {
			p.renderEmptyLine(true)
			p.listItemNumber = -1

		}
	}
}

func (p *MarkdownProvider) prepareListItem(node *ast.ListItem, entering bool) {
	var prefix string

	if entering {
		p.listItemNumber++
		if node.ListFlags&ast.ListTypeOrdered != 0 {
			// numbered list
			prefix = fmt.Sprintf("%d. ", p.listItemNumber+1)
		} else if node.ListFlags&ast.ListTypeDefinition != 0 {
			prefix = ""
		} else {
			//prefix = bulletUnicode + " "
		}

		p.openTag(listItemTagName)
		p.stringBuilder.WriteString(prefix)
	}
}

func (p *MarkdownProvider) prepareParagraph(node *ast.Paragraph, entering bool) {
	if !entering {
		p.render()
		p.renderEmptyLine(false)
	} else {
		p.stringBuilder.WriteString(string(node.Literal))
	}
}

func (p *MarkdownProvider) prepareHeading(node *ast.Heading, entering bool) {
	if entering {
		var tag string
		switch node.Level {
		case 1:
			tag = h1TagName
		case 2:
			tag = h2TagName
		case 3:
			tag = h3TagName
		case 4:
			tag = h4TagName
		case 5:
			tag = h5TagName
		case 6:
			tag = h6TagName
		}
		p.openTag(tag)
	} else {
		p.closeTag()
		p.render()
		p.renderEmptyLine(false)
	}
}

func (p *MarkdownProvider) prepareLink(node *ast.Link, entering bool) {
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

	word := linkTag + linkSpacer + dest + linkSpacer + strings.Replace(text, " ", "---", -1)
	p.stringBuilder.WriteString(word)
}

func (p *MarkdownProvider) prepareText(node *ast.Text, entering bool) {
	if string(node.Literal) == "\n" {
		return
	}

	content := string(node.Literal)
	if shouldCleanText(node) || p.listItemNumber > -1 {
		content = removeLineBreak(content)
	}
	p.stringBuilder.WriteString(content)
	if p.listItemNumber > -1 {
		p.closeTag()
	}
}

func (p *MarkdownProvider) prepareTable(node *ast.Table, entering bool) {
	if entering {
		p.table = newTable(p.theme)
	} else {
		p.containers = append(p.containers, p.table.render())
		p.table = nil
	}
}

func (p *MarkdownProvider) prepareTableCell(node *ast.TableCell, entering bool) {
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()

	align := cellAlignLeft
	switch node.Align {
	case ast.TableAlignmentRight:
		align = cellAlignRight
	case ast.TableAlignmentCenter:
		align = cellAlignCenter
	}

	var alignment cellAlign
	if node.IsHeader {
		alignment = align
	} else {
		alignment = cellAlignCopyHeader
	}
	p.table.addCell(content, alignment, node.IsHeader)
}

func (p *MarkdownProvider) prepareTableRow(node *ast.TableRow, entering bool) {
	if _, ok := node.Parent.(*ast.TableBody); ok && entering {
		p.table.startNextRow()
	}
	if _, ok := node.Parent.(*ast.TableFooter); ok && entering {
		p.table.startNextRow()
	}
}

func (p *MarkdownProvider) pushTag(tagName string) {
	p.tagStack = append(p.tagStack, tagName)
}

func (p *MarkdownProvider) popTag() {
	if len(p.tagStack) > 0 {
		p.tagStack = p.tagStack[:len(p.tagStack)-1]
	}
}

func (p *MarkdownProvider) renderEmptyLine(isList bool) {
	p.containers = append(p.containers, renderEmptyLine(p.theme, isList))
}

func (p *MarkdownProvider) renderLineBreak() layout.Widget {
	return func(gtx C) D {
		dims := p.theme.Body2("").Layout(gtx)
		dims.Size.Y = dims.Size.Y + 5
		return dims
	}
}

func (p *MarkdownProvider) renderCurrentText(txt string) layout.Widget {
	lbl := p.theme.Body1(txt)
	var container layout.Widget

	for index := range p.tagStack {
		i := index
		switch p.tagStack[i] {
		case listItemTagName:
			container = renderListItem(lbl, p.theme)
		case italicsTagName, emphTagName:
			lbl = setStyle(lbl, italicsTagName)
		case strongTagName:
			lbl = setWeight(lbl, strongTagName)
		case strikeTagName:
			container = renderStrike(lbl, p.theme)
		case blockQuoteTagName:
			container = renderBlockQuote(lbl, p.theme)
		case h1TagName, h2TagName, h3TagName, h4TagName, h5TagName, h6TagName:
			lbl = getHeading(txt, p.tagStack[i], p.theme)
		default:

		}
	}

	if container == nil {
		return lbl.Layout
	}

	return container
}

func (p *MarkdownProvider) render() {
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()

	var wdgts []layout.Widget
	var isGettingTagName bool
	var isClosingBlock bool
	var isInBlock bool
	var currentTag string
	var currText string

	//fmt.Println(content)

	for index := range content {
		i := index
		curr := content[i]

		if curr == openTagPrefix[0] && getNextChar(content, i) == openTagPrefix[1] {
			wdgts = append(wdgts, p.renderCurrentText(currText))
			currText = ""

			isGettingTagName = true
			isInBlock = true
			currentTag = string(curr)
			continue
		} else if !isInBlock {
			currStr := string(curr)
			currText += currStr

			if i+1 == len(content) || currStr == "" || currStr == " " {
				if strings.HasPrefix(currText, linkTag) {
					wdgts = append(wdgts,  p.getLinkWidget(currText))
				} else {
					lbl := p.theme.Body1(currText)
					wdgts = append(wdgts, lbl.Layout)
				}

				currText = ""
			}
			continue
		}

		if isGettingTagName {
			currentTag += string(curr)
			if curr == openTagSuffix[1] {
				isGettingTagName = false
				currentTag = strings.Replace(currentTag, openTagPrefix, "", -1)
				currentTag = strings.Replace(currentTag, openTagSuffix, "", -1)
				p.pushTag(currentTag)
				currentTag = ""
			}
			continue
		}

		if curr == closeTag[0] && getNextChar(content, i) == closeTag[1] {
			isClosingBlock = true
			continue
		}

		if isClosingBlock {
			if curr == closeTag[2] {
				wdgts = append(wdgts, p.renderCurrentText(currText))
				p.popTag()
				currText = ""
				isClosingBlock = false
				isInBlock = false
			}
			continue
		}
		currText += string(curr)
	}

	wdgt := func(gtx C) D {
		return decredmaterial.GridWrap{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		}.Layout(gtx, len(wdgts), func(gtx C, i int) D {
			if wdgts[i] == nil {
				return D{}
			}

			return wdgts[i](gtx)
		})
	}
	p.containers = append(p.containers, wdgt)
}

func (p *MarkdownProvider) openOrCloseTag(tagName string, entering bool) {
	if entering {
		p.openTag(tagName)
	} else {
		p.closeTag()
	}
}
func (p *MarkdownProvider) openTag(tagName string) {
	tag := openTagPrefix + tagName + openTagSuffix
	p.stringBuilder.WriteString(tag)
}

func (p *MarkdownProvider) closeTag() {
	p.stringBuilder.WriteString(closeTag)
}

func shouldCleanText(node ast.Node) bool {
	for node != nil {
		switch node.(type) {
		case *ast.BlockQuote:
			return false

		case *ast.Heading, *ast.Image, *ast.Link,
			*ast.TableCell, *ast.Document, *ast.ListItem:
			return true
		}
		node = node.GetParent()
	}

	return false
}

func removeLineBreak(text string) string {
	lines := strings.Split(text, "\n")

	if len(lines) <= 1 {
		return text
	}

	for i, l := range lines {
		switch i {
		case 0:
			lines[i] = strings.TrimRightFunc(l, unicode.IsSpace)
		case len(lines) - 1:
			lines[i] = strings.TrimLeftFunc(l, unicode.IsSpace)
		default:
			lines[i] = strings.TrimFunc(l, unicode.IsSpace)
		}
	}

	return strings.Join(lines, " ")
}

func getNextChar(content string, currIndex int) byte {
	if currIndex+2 <= len(content) {
		return content[currIndex+1]
	}

	return 0
}

/**func RenderMarkdown(gtx layout.Context, theme *decredmaterial.Theme, source string) interface{} {
	/**extensions := parser.NoIntraEmphasis        // Ignore emphasis markers inside words
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

	p := parser.NewWithExtensions(extensions)

	r := &MarkdownRenderer{
		//newRenderer(theme, false),
	}

	source = r.prepareDocForTable(source)
	nodes := md.Parse([]byte(source), p)

	md.Render(nodes, r.Renderer)

	return r
	return nil
}

func (r *MarkdownRenderer) prepareDocForTable(doc string) string {
	d := strings.Replace(doc, ":|", "------:|", -1)
	d = strings.Replace(d, "-|", "------|", -1)
	d = strings.Replace(d, "|-", "|------", -1)
	d = strings.Replace(d, "|:-", "|:------", -1)

	return d
}**/
