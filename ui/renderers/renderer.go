package renderers

import (
	"io"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
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
	renderSoftBreak()
	renderHardBreak()
}

type nodeWalker struct {
	rootNode ast.Node
	renderer renderer
}

func newNodeWalker(doc string, renderer renderer) *nodeWalker {
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
	extensions |= parser.EmptyLinesBreakList
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
	case *ast.Softbreak:
		nw.renderer.renderSoftBreak()
	case *ast.Hardbreak:
		nw.renderer.renderSoftBreak()
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

func (*nodeWalker) RenderHeader(w io.Writer, node ast.Node) {}

func (*nodeWalker) RenderFooter(w io.Writer, node ast.Node) {}
