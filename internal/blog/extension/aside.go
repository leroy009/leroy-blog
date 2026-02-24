// Package extension provides custom goldmark extensions for the blog.
// Aside blocks are fenced with ::: on their own lines:
//
//	:::note
//	This is a note.
//	:::
//
// Supported kinds: note, warning, tip, info (or any arbitrary string).
package extension

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ─── AST node ────────────────────────────────────────────────────────────────

// KindAside is the goldmark NodeKind for an aside container block.
var KindAside = ast.NewNodeKind("Aside")

// AsideBlock is a container block node that renders as <aside>.
type AsideBlock struct {
	ast.BaseBlock
	// AsideKind holds the type word after :::, e.g. "note", "warning".
	AsideKind string
}

func NewAside(kind string) *AsideBlock {
	return &AsideBlock{AsideKind: kind}
}

func (a *AsideBlock) Kind() ast.NodeKind { return KindAside }

func (a *AsideBlock) Dump(source []byte, level int) {
	ast.DumpHelper(a, source, level, map[string]string{"AsideKind": a.AsideKind}, nil)
}

// ─── Block parser ─────────────────────────────────────────────────────────────

type asideParser struct{}

var defaultAsideParser = &asideParser{}

func (p *asideParser) Trigger() []byte { return []byte{':'} }

// Open is called when the first `:` of a line is seen.
// We only proceed if the line is `:::kind` with a non-empty kind.
func (p *asideParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, seg := reader.PeekLine()
	line = bytes.TrimRight(line, "\r\n")

	if !bytes.HasPrefix(line, []byte(":::")) {
		return nil, parser.NoChildren
	}

	kind := string(bytes.TrimSpace(line[3:]))
	if kind == "" {
		return nil, parser.NoChildren
	}

	reader.Advance(seg.Len()) // consume the opening fence line
	return NewAside(kind), parser.HasChildren
}

// Continue is called for every subsequent line.
// A lone `:::` closes the block.
func (p *asideParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, seg := reader.PeekLine()
	if bytes.TrimSpace(line) != nil && bytes.Equal(bytes.TrimSpace(line), []byte(":::")) {
		reader.Advance(seg.Len()) // consume the closing fence line
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (p *asideParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}

func (p *asideParser) CanInterruptParagraph() bool { return true }
func (p *asideParser) CanAcceptIndentedLine() bool { return false }

// ─── Renderer ────────────────────────────────────────────────────────────────

type asideRenderer struct{}

func (r *asideRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindAside, r.render)
}

// render outputs an <aside> element with Tailwind classes based on kind.
func (r *asideRenderer) render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*AsideBlock)
	if entering {
		cls := asideClass(n.AsideKind)
		_, _ = w.WriteString(`<aside class="` + cls + `">`)
	} else {
		_, _ = w.WriteString(`</aside>`)
	}
	return ast.WalkContinue, nil
}

// asideClass returns Tailwind utility classes for a given aside kind.
func asideClass(kind string) string {
	switch kind {
	case "warning":
		return "my-6 border-l-4 border-yellow-400 bg-yellow-50 p-4 rounded-r-lg text-yellow-900"
	case "tip":
		return "my-6 border-l-4 border-green-400 bg-green-50 p-4 rounded-r-lg text-green-900"
	case "info":
		return "my-6 border-l-4 border-sky-400 bg-sky-50 p-4 rounded-r-lg text-sky-900"
	default: // "note" and everything else
		return "my-6 border-l-4 border-blue-400 bg-blue-50 p-4 rounded-r-lg text-blue-900"
	}
}

// ─── Extension ───────────────────────────────────────────────────────────────

// Aside is the goldmark extension that registers the aside block parser and renderer.
type Aside struct{}

func (e *Aside) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(defaultAsideParser, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&asideRenderer{}, 500),
		),
	)
}
