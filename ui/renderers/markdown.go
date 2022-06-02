package renderers

import (
	"strings"
	"unicode"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/gomarkdown/markdown/ast"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	bulletUnicode = "\u2022"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type layoutRow struct {
	widgets []layout.Widget
}

type MarkdownProvider struct {
	containers     []layoutRow
	theme          *decredmaterial.Theme
	listItemNumber int // should be negative when not rendering a list
	links          map[string]*widget.Clickable
	table          *table
	label          *decredmaterial.Label
	prefix         string

	stringBuilder strings.Builder
	tagStack      []string

	shouldRemoveBold bool
}

func RenderMarkdown(gtx C, theme *decredmaterial.Theme, source string) *MarkdownProvider {
	lbl := theme.Body1("")
	source = strings.Replace(source, " \n*", " \n\n *", -1)

	mdProvider := &MarkdownProvider{
		theme:          theme,
		listItemNumber: -1,
		label:          &lbl,
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

func (p *MarkdownProvider) Layout() ([]layout.Widget, map[string]*widget.Clickable) {
	w := func(gtx C) D {
		max := gtx.Constraints.Max.X
		rows := layout.List{Axis: layout.Vertical}
		return rows.Layout(gtx, len(p.containers), func(gtx C, i int) D {
			return decredmaterial.GridWrap{
				Axis:      layout.Horizontal,
				Alignment: layout.Start,
			}.Layout(gtx, len(p.containers[i].widgets), func(gtx C, j int) D {
				gtx.Constraints.Max.X = max
				return p.containers[i].widgets[j](gtx)
			})
		})
	}

	return []layout.Widget{w}, p.links
}

func (p *MarkdownProvider) prepareBlockQuote(node *ast.BlockQuote, entering bool) {
	p.openOrCloseTag(blockQuoteTagName, entering)
}

func (p *MarkdownProvider) prepareCode(node *ast.Code, entering bool) {
	content := string(node.Literal)
	p.stringBuilder.WriteString(content)
}

func (p *MarkdownProvider) prepareCodeBlock(node *ast.CodeBlock, entering bool) {
	content := string(node.Literal)
	p.createNewRow()
	wdg := func(gtx C) D {
		return decredmaterial.LinearLayout{
			Orientation: layout.Vertical,
			Width:       decredmaterial.WrapContent,
			Height:      decredmaterial.WrapContent,
			Background:  p.theme.Color.Background,
			Padding:     layout.UniformInset(values.MarginPadding16),
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return p.theme.Body1(content).Layout(gtx)
			}),
		)
	}

	p.appendToLastRow(wdg)
	p.addVerticalSpacing(15)
}

func (p *MarkdownProvider) renderSoftBreak() {
	p.createNewRow()
}

func (p *MarkdownProvider) renderHardBreak() {
	p.createNewRow()
}

func (p *MarkdownProvider) prepareStrong(node *ast.Strong, entering bool) {
	p.openOrCloseTag(strongTagName, entering)
}

func (p *MarkdownProvider) prepareDel(node *ast.Del, entering bool) {
	p.openOrCloseTag(blockQuoteTagName, entering)
}

func (p *MarkdownProvider) prepareEmph(node *ast.Emph, entering bool) {
	p.openOrCloseTag(emphTagName, entering)
}

func (p *MarkdownProvider) prepareHorizontalRule(node *ast.HorizontalRule, entering bool) {
	p.drawLineRow(layout.Horizontal)
}

func (p *MarkdownProvider) prepareList(node *ast.List, entering bool) {
	if entering {
		p.listItemNumber = 1
		p.prefix = bulletUnicode + " "
	} else {
		p.listItemNumber = -1
	}

	if next := ast.GetNextNode(node); !entering && next != nil {
		_, parentIsListItem := node.GetParent().(*ast.ListItem)
		_, nextIsList := next.(*ast.List)
		if !nextIsList && !parentIsListItem {
			p.listItemNumber = -1
		}
	}
}

func (p *MarkdownProvider) renderListItem(content string) {
	if strings.Trim(content, " ") == "" {
		return
	}

	w := func(gtx C) D {
		lbl := p.getLabel()
		strongLabel := p.getLabel()
		strongLabel.Font.Weight = text.Bold

		return layout.Flex{}.Layout(gtx,
			layout.Flexed(0.02, func(gtx C) D {
				strongLabel.Text = ""
				return strongLabel.Layout(gtx)
			}),
			layout.Flexed(0.05, func(gtx C) D {
				strongLabel.Text = p.prefix
				return strongLabel.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				lbl.Text = content
				return lbl.Layout(gtx)
			}),
		)
	}

	p.createNewRow()
	p.appendToLastRow(w)
}

func (p *MarkdownProvider) prepareListItem(node *ast.ListItem, entering bool) {

}

func (p *MarkdownProvider) prepareParagraph(node *ast.Paragraph, entering bool) {
	if !entering {
		p.renderBlock()
		p.createNewRow()
		p.addVerticalSpacing(15)
	}
}

func (p *MarkdownProvider) prepareHeading(node *ast.Heading, entering bool) {
	if !entering {
		content := p.stringBuilder.String()
		p.stringBuilder.Reset()
		p.createNewRow()
		p.appendToLastRow(getHeading(content, node.Level, p.theme).Layout)
		p.addVerticalSpacing(8)
		if node.Level == 1 {
			p.drawLineRow(layout.Horizontal)
			p.addVerticalSpacing(14)
		}
	}
}

func (p *MarkdownProvider) prepareLink(node *ast.Link, entering bool) {
	p.renderBlock()
}

func (p *MarkdownProvider) renderBlock() {
	content := p.stringBuilder.String()
	p.stringBuilder.Reset()

	var inBlock bool
	var isGettingTagName bool
	var isClosingBlock bool
	var currentTag string
	currText := new(strings.Builder)
	for i := range content {
		curr := content[i]

		if curr == openTagPrefix[0] && getNextChar(content, i) == openTagPrefix[1] {
			p.render(currText)
			isGettingTagName = true
			inBlock = true
			currentTag = string(curr)
			continue
		} else if !inBlock {
			currText.WriteByte(curr)
			if i+1 == len(content) || curr == openTagPrefix[0] || curr == closeTag[0] {
				p.render(currText)
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
				p.render(currText)
				p.popTag()
				isClosingBlock = false
				inBlock = false
			}
			continue
		}
		currText.WriteByte(curr)
	}
}

func (p *MarkdownProvider) getLabel() decredmaterial.Label {
	lbl := p.theme.Body1("")
	if len(p.tagStack) > 0 {
		for i := range p.tagStack {
			switch p.tagStack[i] {
			case strongTagName:
				setWeight(&lbl, "bold")
			case emphTagName:
				setStyle(&lbl, "italic")
			case listItemTagName:
				//p.createNewRow()
			}
		}
	}

	return lbl
}

func (p *MarkdownProvider) render(content *strings.Builder) {
	lbl := p.getLabel()
	words := strings.Fields(content.String())
	content.Reset()

	for index := range words {
		lbl.Text = words[index] + " "
		p.appendToLastRow(lbl.Layout)
	}
}

func (p *MarkdownProvider) addVerticalSpacing(height int) {
	p.appendToLastRow(func(gtx C) D {
		dims := p.theme.Caption(" ").Layout(gtx)
		dims.Size.X = gtx.Constraints.Max.X
		dims.Size.Y = height
		return dims
	})
}

func (p *MarkdownProvider) createNewRow() {
	row := layoutRow{}
	p.containers = append(p.containers, row)
}

func (p *MarkdownProvider) appendToLastRow(wdgt layout.Widget) {
	if len(p.containers) == 0 {
		p.createNewRow()
	}

	l := len(p.containers)
	lastRow := p.containers[l-1]
	lastRow.widgets = append(lastRow.widgets, wdgt)
	p.containers[l-1] = lastRow
}

func (p *MarkdownProvider) drawLineRow(axis layout.Axis) {
	var l decredmaterial.Line

	if axis == layout.Vertical {
		l = p.theme.SeparatorVertical(1, 10)
	} else {
		l = p.theme.Separator()
	}

	p.createNewRow()
	p.appendToLastRow(l.Layout)
}

func (p *MarkdownProvider) prepareText(node *ast.Text, entering bool) {
	if string(node.Literal) == "\n" {
		return
	}

	content := string(node.Literal)
	if shouldCleanText(node) {
		content = removeLineBreak(content)
	}

	if p.listItemNumber > -1 {
		p.renderListItem(content)
		return
	}

	p.stringBuilder.WriteString(content)
}

func (p *MarkdownProvider) prepareTable(node *ast.Table, entering bool) {
	if entering {
		p.table = newTable(p.theme)
	} else {
		p.createNewRow()
		p.appendToLastRow(p.table.render())
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
