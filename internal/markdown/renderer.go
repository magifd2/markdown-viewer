package markdown

import (
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// SafeLinkRenderer is a custom renderer for links that prevents creating links
// to non-Markdown local files and all external websites.
type SafeLinkRenderer struct {
	html.Config
}

// NewSafeLinkRenderer creates a new instance of the SafeLinkRenderer.
func NewSafeLinkRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SafeLinkRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers the renderer for ast.Link nodes.
func (r *SafeLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
}

// isSafeLink checks if a given link destination is considered safe to render as a hyperlink.
// In this application, only relative paths to Markdown files are considered safe.
func isSafeLink(destination []byte) bool {
	destStr := string(destination)

	// Disallow all external links.
	if strings.HasPrefix(destStr, "http://") || strings.HasPrefix(destStr, "https://") {
		return false
	}

	// Allow only links to local Markdown files.
	ext := strings.ToLower(filepath.Ext(destStr))
	if ext == ".md" || ext == ".markdown" {
		return true
	}

	return false
}

// renderLink is the rendering function for links.
func (r *SafeLinkRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)

	// If the link is not safe, we don't render the <a> tag at all.
	// The link's text (children) will be rendered as plain text.
	if !isSafeLink(n.Destination) {
		return ast.WalkContinue, nil
	}

	// It's a safe link (local Markdown file), so render the <a> tag.
	if entering {
		_, _ = w.WriteString("<a href=\"")
		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_, _ = w.WriteString("\"")
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			_, _ = w.Write(n.Title)
			_, _ = w.WriteString(`"`)
		}
		_ = w.WriteByte('>')
	} else { // exiting
		_, _ = w.WriteString("</a>")
	}

	return ast.WalkContinue, nil
}
