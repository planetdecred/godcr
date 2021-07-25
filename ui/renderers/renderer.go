package renderers

import (
	//"fmt"
	//"image"
	"io"
	//"reflect"
	//"regexp"
	//"strings"
	//"unicode"

	//"gioui.org/layout"
	//"gioui.org/text"
	//"gioui.org/unit"
	//"gioui.org/widget"
	//"gioui.org/widget/material"

	//md "github.com/JohannesKaufmann/html-to-markdown"
	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	//"github.com/planetdecred/godcr/ui/decredmaterial"
)

type renderer interface {
	prepareText(node *ast.Text, entering bool)
	prepareBlockQuote(node *ast.BlockQuote, entering bool)
	prepareList(node *ast.List, entering bool)
	prepareListItem(node *ast.ListItem, entering bool)
	prepareParagraph(node *ast.Paragraph, entering bool)
	prepareHeading(node *ast.Heading, entering bool)
	prepareStrong(node *ast.Strong, entering bool)
	prepareDel(node *ast.Del, entering bool)
	prepareEmph(node *ast.Emph, entering bool)
	prepareLink(node *ast.Link, entering bool)
	prepareTable(node *ast.Table, entering bool)
	prepareTableCell(node *ast.TableCell, entering bool)
	prepareTableRow(node *ast.TableRow, entering bool)
	prepareHorizontalRule(node *ast.HorizontalRule, entering bool)
	//layout() ([]layout.Widget, map[string]*widget.Clickable)
}

type nodeWalker struct {
	rootNode ast.Node
	renderer renderer
}

func newNodeWalker(doc string, renderer renderer) *nodeWalker {
	extensions := parser.NoIntraEmphasis // Ignore emphasis markers inside words
	extensions |= parser.Tables          // Parse tables
	extensions |= parser.FencedCode      // Parse fenced code blocks
	extensions |= parser.Autolink        // Detect embedded URLs that are not explicitly marked
	extensions |= parser.Strikethrough   // Strikethrough text using ~~test~~
	extensions |= parser.SpaceHeadings   // Be strict about prefix heading rules
	//extensions |= parser.HeadingIDs             // specify heading IDs  with {#id}
	extensions |= parser.BackslashLineBreak // Translate trailing backslashes into line breaks
	extensions |= parser.DefinitionLists    // Parse definition lists
	extensions |= parser.LaxHTMLBlocks      // more in HTMLBlock, less in HTMLSpan
	//extensions |= parser.NoEmptyLineBeforeBlock // no need for new line before a list
	extensions |= parser.Attributes
	//extensions |= parser.EmptyLinesBreakList
	extensions |= parser.Mmark
	extensions |= parser.LaxHTMLBlocks

	p := parser.NewWithExtensions(extensions)

	return &nodeWalker{
		rootNode: md.Parse([]byte(doc), p),
		renderer: renderer,
	}
}

func (nw *nodeWalker) walk() {
	md.Render(nw.rootNode, nw)
}

func (nw *nodeWalker) walkerFunc() {

}

func (nw *nodeWalker) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Document:
		//fmt.Println(string(node.Literal))
	case *ast.BlockQuote:
		nw.renderer.prepareBlockQuote(node, entering)
	case *ast.List:
		nw.renderer.prepareList(node, entering)
	case *ast.ListItem:
		nw.renderer.prepareListItem(node, entering)
	case *ast.Paragraph:
		nw.renderer.prepareParagraph(node, entering)
	case *ast.Heading:
		nw.renderer.prepareHeading(node, entering)
	case *ast.Strong:
		nw.renderer.prepareStrong(node, entering)
	case *ast.Del:
		nw.renderer.prepareDel(node, entering)
	case *ast.Emph:
		nw.renderer.prepareEmph(node, entering)
	case *ast.Link:
		if !entering {
			nw.renderer.prepareLink(node, entering)
			return ast.SkipChildren
		}
	case *ast.Text:
		nw.renderer.prepareText(node, entering)
	case *ast.HorizontalRule:
		nw.renderer.prepareHorizontalRule(node, entering)
	case *ast.Table:
		nw.renderer.prepareTable(node, entering)
	case *ast.TableRow:
		nw.renderer.prepareTableRow(node, entering)
	case *ast.TableCell:
		if !entering {
			nw.renderer.prepareTableCell(node, entering)
		}
	}
	return ast.GoToNext
}

/**func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Document:
		// Nothing to do
	case *ast.BlockQuote:
		r.renderBlockQuote(entering)
	case *ast.List:
		// extra new line at the end of a list *if* next is not a list
		if next := ast.GetNextNode(node); !entering && next != nil {
			_, parentIsListItem := node.GetParent().(*ast.ListItem)
			_, nextIsList := next.(*ast.List)
			if !nextIsList && !parentIsListItem {
				r.renderEmptyLine()
			}
		}
	case *ast.ListItem:
		r.renderList(node, entering)
	case *ast.Paragraph:
		if !entering {
			r.renderParagraph()
		}
	case *ast.Heading:
		if !entering {
			r.renderHeading(node.Level, true)
		}
	case *ast.Strong:
		r.renderStrong(entering)
	case *ast.Del:
		r.renderDel(entering)
	case *ast.Emph:
		r.renderEmph(entering)
	case *ast.Link:
		if !entering {
			r.renderLink(node)
			return ast.SkipChildren
		}
	case *ast.Text:
		r.renderText(node)
	case *ast.Table:
		r.renderTable(entering)
	case *ast.TableCell:
		if !entering {
			r.renderTableCell(node)
		}
	case *ast.TableRow:
		r.renderTableRow(node, entering)
	}

	return ast.GoToNext
}
**/

func (*nodeWalker) RenderHeader(w io.Writer, node ast.Node) {}

func (*nodeWalker) RenderFooter(w io.Writer, node ast.Node) {}

/**type labelFunc func(string) decredmaterial.Label

type Renderer struct {
	theme      *decredmaterial.Theme
	provider   provider
	isList     bool
	isListItem bool

	prefix         string
	padAccumulator []string
	leftPad        int
	links          map[string]*widget.Clickable
	stringBuilder  strings.Builder
	containers     []layout.Widget
	styleGroups    []map[string]string
	isHTML         bool

	table *table
}

const (
	bulletUnicode = "\u2022"
	linkTag       = "[[link"
	linkSpacer    = "@@@@"
)

var (
	textBeforeCloseRegexp = regexp.MustCompile("(.*){/#}")
)

/**func newRenderer(theme *decredmaterial.Theme, isHTML bool) *Renderer {
	return &Renderer{
		theme:  theme,
		isHTML: isHTML,
	}
}
/**
func (r *Renderer) pad() string {
	return strings.Repeat(" ", r.leftPad) + strings.Join(r.padAccumulator, "")
}

func (r *Renderer) addPad(pad string) {
	r.padAccumulator = append(r.padAccumulator, pad)
}

func (r *Renderer) popPad() {
	r.padAccumulator = r.padAccumulator[:len(r.padAccumulator)-1]
}

func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Document:
		// Nothing to do
	case *ast.BlockQuote:
		r.renderBlockQuote(entering)
	case *ast.List:
		// extra new line at the end of a list *if* next is not a list
		if next := ast.GetNextNode(node); !entering && next != nil {
			_, parentIsListItem := node.GetParent().(*ast.ListItem)
			_, nextIsList := next.(*ast.List)
			if !nextIsList && !parentIsListItem {
				r.renderEmptyLine()
			}
		}
	case *ast.ListItem:
		r.renderList(node, entering)
	case *ast.Paragraph:
		if !entering {
			r.renderParagraph()
		}
	case *ast.Heading:
		if !entering {
			r.renderHeading(node.Level, true)
		}
	case *ast.Strong:
		r.renderStrong(entering)
	case *ast.Del:
		r.renderDel(entering)
	case *ast.Emph:
		r.renderEmph(entering)
	case *ast.Link:
		if !entering {
			r.renderLink(node)
			return ast.SkipChildren
		}
	case *ast.Text:
		r.renderText(node)
	case *ast.Table:
		r.renderTable(entering)
	case *ast.TableCell:
		if !entering {
			r.renderTableCell(node)
		}
	case *ast.TableRow:
		r.renderTableRow(node, entering)
	}

	return ast.GoToNext
}

func (r *Renderer) openMarkdownTag(tagName string) {
	tag := openTagPrefix + tagName + openTagSuffix
	r.stringBuilder.WriteString(tag)
}

func (r *Renderer) closeMarkdownTag() {
	r.stringBuilder.WriteString(closeTag)
}

func (r *Renderer) renderBlockQuote(entering bool) {
	if r.isHTML {
		return
	} else if entering {
		r.openMarkdownTag(blockQuoteTagName)
	} else {
		r.closeMarkdownTag()
	}
}

func (r *Renderer) renderStrong(entering bool) {
	if r.isHTML {
		label := r.theme.Body1("")
		label.Font.Weight = text.Bold
		r.renderWords(label)
	} else if entering {
		r.openMarkdownTag(strongTagName)
	} else {
		r.closeMarkdownTag()
	}
}

func (r *Renderer) renderEmph(entering bool) {
	if r.isHTML {
		return
	} else if entering {
		r.openMarkdownTag(emphTagName)
	} else {
		r.closeMarkdownTag()
	}
}

func (r *Renderer) renderDel(entering bool) {
	if r.isHTML {
		return
	} else if entering {
		r.openMarkdownTag(strikeTagName)
	} else {
		r.closeMarkdownTag()
	}
}

func (r *Renderer) renderParagraph() {
	r.renderWords(r.theme.Body1(""))
	r.renderEmptyLine()
}

func (r *Renderer) renderHeading(level int, block bool) {
	lblFunc := r.theme.H6

	switch level {
	case 1:
		lblFunc = r.theme.H4
	case 2:
		lblFunc = r.theme.H5
	case 3:
		lblFunc = r.theme.H6
	}

	r.renderWords(lblFunc(""))
	if block {
		// add dummy widget for new line
		r.renderEmptyLine()
	}
}

func (r *Renderer) renderLink(node *ast.Link) {
	dest := string(node.Destination)
	text := string(ast.GetFirstChild(node).AsLeaf().Literal)

	if r.links == nil {
		r.links = map[string]*widget.Clickable{}
	}

	if _, ok := r.links[dest]; !ok {
		r.links[dest] = new(widget.Clickable)
	}

	// fix a bug that causes the link to be written to the builder before this is called
	content := r.stringBuilder.String()
	r.stringBuilder.Reset()
	parts := strings.Split(content, " ")
	parts = parts[:len(parts)-1]
	for i := range parts {
		r.stringBuilder.WriteString(parts[i] + " ")
	}

	word := linkTag + linkSpacer + dest + linkSpacer + text
	r.stringBuilder.WriteString(word)
}

func (r *Renderer) renderText(node *ast.Text) {
	if string(node.Literal) == "\n" {
		return
	}

	content := string(node.Literal)
	if shouldCleanText(node) {
		content = removeLineBreak(content)
	}
	r.stringBuilder.WriteString(content)
}

func (r *Renderer) getNextChar(content string, currIndex int) byte {
	if currIndex+2 <= len(content) {
		return content[currIndex+1]
	}

	return 0
}

func (r *Renderer) renderWords(lbl decredmaterial.Label) {
	if r.isHTML {
		r.renderHTML(lbl)
		return
	}

	r.renderMarkdown(lbl)
}

func (r *Renderer) getLabel(lbl decredmaterial.Label, text string) decredmaterial.Label {
	l := lbl
	l.Text = text
	if r.isHTML {
		l = r.styleHTMLLabel(l)
	} else {

	}
	return l
}

func (r *Renderer) getMarkdownWidget(lbl decredmaterial.Label, text string) layout.Widget {
	l := lbl
	l.Text = text

	return r.getMarkdownWidgetAndStyle(l)
}

func (r *Renderer) renderMarkdown(lbl decredmaterial.Label) {
	content := r.stringBuilder.String()
	r.stringBuilder.Reset()

	var wdgts []layout.Widget
	var isGettingTagName bool
	var isClosingBlock bool
	var currentTag string
	var currText string

	for i := range content {
		curr := content[i]

		if curr == openTagPrefix[0] && r.getNextChar(content, i) == openTagPrefix[1] {
			wdgts = append(wdgts,r.getMarkdownWidget(lbl, currText))
			currText = ""

			isGettingTagName = true
			currentTag = string(curr)
			continue
		}

		if isGettingTagName {
			currentTag += string(curr)
			if curr == openTagSuffix[1] {
				isGettingTagName = false
				currentTag = strings.Replace(currentTag, openTagPrefix, "", -1)
				currentTag = strings.Replace(currentTag, openTagSuffix, "", -1)
				r.addStyleGroupFromTagName(currentTag)
				currentTag = ""
			}
			continue
		}

		if curr == closeTag[0] && r.getNextChar(content, i) == closeTag[1] {
			isClosingBlock = true
			continue
		}

		if isClosingBlock {
			if curr == closeTag[2] {
				wdgts = append(wdgts,r.getMarkdownWidget(lbl, currText))
				r.removeLastStyleGroup()
				currText = ""
				isClosingBlock = false
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
	r.containers = append(r.containers, wdgt)

}

func (r *Renderer) strikeLabel(label decredmaterial.Label) layout.Widget {
	return func(gtx C) D {
		var dims D
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				dims = label.Layout(gtx)
				return dims
			}),
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: unit.Dp((float32(dims.Size.Y) / float32(2))),
				}.Layout(gtx, func(gtx C) D {
					l := r.theme.Separator()
					l.Color = label.Color
					l.Width = dims.Size.X
					return l.Layout(gtx)
				})
			}),
		)
	}
}

func (r *Renderer) renderHTML(lbl decredmaterial.Label) {
	content := r.stringBuilder.String()
	r.stringBuilder.Reset()

	var labels []decredmaterial.Label
	var inStyleBlock bool
	var isClosingStyle bool
	var isClosingBlock bool
	var currStyle string
	var currText string
	for i := range content {
		curr := content[i]

		if curr == openStyleTag[0] && r.getNextChar(content, i) == openStyleTag[1] {
			inStyleBlock = true
			labels = append(labels, r.getLabel(lbl, currText))
			currText = ""
		}

		if curr == halfCloseStyleTag[0] && r.getNextChar(content, i) == halfCloseStyleTag[1] {
			isClosingStyle = true
		}

		if curr == closeStyleTag[0] && r.getNextChar(content, i) == closeStyleTag[1] {
			isClosingBlock = true
		}

		if !inStyleBlock && !isClosingBlock {
			currStr := string(curr)
			currText += currStr

			if i+1 == len(content) || currStr == "" || currStr == " " {
				labels = append(labels, r.getLabel(lbl, currText))
				currText = ""
			}
		}

		if isClosingBlock && curr == closeStyleTag[3] {
			labels = append(labels, r.getLabel(lbl, currText))
			currText = ""
			r.removeLastStyleGroup()
			isClosingBlock = false

		}

		if inStyleBlock && !isClosingStyle {
			currStyle += string(curr)
		}

		if isClosingStyle && curr == halfCloseStyleTag[1] {
			isClosingStyle = false
			inStyleBlock = false
			r.addHTMLStyleGroup(currStyle)
			currStyle = ""
		}
	}

	wdgt := func(gtx C) D {
		return decredmaterial.GridWrap{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		}.Layout(gtx, len(labels), func(gtx C, i int) D {
			return labels[i].Layout(gtx)
		})
	}
	r.containers = append(r.containers, wdgt)
}

func (r *Renderer) renderEmptyLine() {
	var padding = -5

	if r.isList {
		padding = -10
		r.isList = false
	}

	r.containers = append(r.containers, func(gtx C) D {
		dims := r.theme.Body2("").Layout(gtx)
		dims.Size.Y = dims.Size.Y + padding
		return dims
	})
}

func (r *Renderer) renderList(node *ast.ListItem, entering bool) {
	if entering {
		r.isList = true
		r.isListItem = true
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
			r.prefix += fmt.Sprintf("%d. ", itemNumber)

		// content of a definition
		case node.ListFlags&ast.ListTypeDefinition != 0:
			r.prefix += " "

		// no flags means it's the normal bullet point list
		default:
			r.prefix += bulletUnicode + " "
		}
	} else {
		r.isListItem = false
	}
}

func (r *Renderer) renderTable(entering bool) {
	if entering {
		r.table = newTable(r.theme)
	} else {
		r.containers = append(r.containers, r.table.render())
		r.table = nil
	}
}

func (r *Renderer) renderTableCell(node *ast.TableCell) {
	content := r.stringBuilder.String()
	r.stringBuilder.Reset()

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
	r.table.addCell(content, alignment, node.IsHeader)
}

func (r *Renderer) renderTableRow(node *ast.TableRow, entering bool) {
	if entering {
		switch node.Parent.(type) {
		case *ast.TableHeader, *ast.TableBody, *ast.TableFooter:
			r.table.startNextRow()
		}
	}
}

func (*Renderer) RenderHeader(w io.Writer, node ast.Node) {}

func (*Renderer) RenderFooter(w io.Writer, node ast.Node) {}

func (r *Renderer) Layout() ([]layout.Widget, map[string]*widget.Clickable) {
	return r.containers, r.links
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

func (r *Renderer) getLinkWidget(gtx layout.Context, linkWord string) D {
	parts := strings.Split(linkWord, linkSpacer)

	gtx.Constraints.Max.X = gtx.Constraints.Max.X - 200
	return material.Clickable(gtx, r.links[parts[1]], func(gtx C) D {
		lbl := r.theme.Body2(parts[2] + " ")
		lbl.Color = r.theme.Color.Primary
		return lbl.Layout(gtx)
	})
}**/
