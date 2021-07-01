package renderers

import (
	"strings"

	"gioui.org/layout"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type MarkdownRenderer struct {
	*Renderer
}

func RenderMarkdown(gtx layout.Context, theme *decredmaterial.Theme, source string) *MarkdownRenderer {
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

	p := parser.NewWithExtensions(extensions)

	r := &MarkdownRenderer{
		newRenderer(theme, false),
	}

	source = r.prepareDocForTable(source)
	nodes := md.Parse([]byte(source), p)

	md.Render(nodes, r.Renderer)

	return r
}

func (r *MarkdownRenderer) prepareDocForTable(doc string) string {
	d := strings.Replace(doc, ":|", "------:|", -1)
	d = strings.Replace(d, "-|", "------|", -1)
	d = strings.Replace(d, "|-", "|------", -1)
	d = strings.Replace(d, "|:-", "|:------", -1)

	return d
}
