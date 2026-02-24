package blog

import (
	"bytes"
	"errors"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/leroy009/leroy-blog/internal/blog/extension"
)

type Service struct {
	reader   PostReader
	markdown goldmark.Markdown
	logger   *slog.Logger
	mu       sync.RWMutex
	cache    map[string]*Post
}

type PostReader interface {
	Read(slug string) ([]byte, error)
	Query() (*PostMetadataCollection, error)
}

// FileReader loads all posts from disk into memory at startup.
type FileReader struct {
	logger   *slog.Logger
	posts    map[string][]byte
	metadata PostMetadataCollection
}

func NewFileReader(dir string, logger *slog.Logger) *FileReader {
	logger = logger.With("component", "file-reader")
	fr := &FileReader{
		logger: logger,
		posts:  make(map[string][]byte),
	}
	fr.load(dir)
	return fr
}

// load reads all .md files from dir into memory and pre-parses metadata.
func (fr *FileReader) load(dir string) {
	filenames, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		fr.logger.Error("failed to glob posts directory", "dir", dir, "error", err)
		return
	}

	for _, filename := range filenames {
		raw, err := os.ReadFile(filename)
		if err != nil {
			fr.logger.Error("failed to read post file", "file", filename, "error", err)
			continue
		}

		slug := strings.TrimSuffix(filepath.Base(filename), ".md")
		fr.posts[slug] = raw

		var meta PostMetadata
		if _, err = frontmatter.Parse(bytes.NewReader(raw), &meta); err != nil {
			fr.logger.Error("failed to parse frontmatter", "file", filename, "error", err)
			continue
		}
		fr.metadata.Posts = append(fr.metadata.Posts, meta)
	}

	fr.logger.Info("posts loaded into memory", "count", len(fr.posts))
}

func (fr *FileReader) Read(slug string) ([]byte, error) {
	raw, ok := fr.posts[slug]
	if !ok {
		return nil, errors.New("post not found")
	}
	return raw, nil
}

func (fr *FileReader) Query() (*PostMetadataCollection, error) {
	return &fr.metadata, nil
}

func NewService(reader PostReader, logger *slog.Logger) *Service {
	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
			&extension.Aside{},
		),
	)
	return &Service{
		reader:   reader,
		markdown: markdown,
		logger:   logger.With("component", "service"),
		cache:    make(map[string]*Post),
	}
}

func (s *Service) GetPostBySlugWithMarkdown(slug string) (*Post, error) {
	s.mu.RLock()
	if post, ok := s.cache[slug]; ok {
		s.mu.RUnlock()
		return post, nil
	}
	s.mu.RUnlock()

	raw, err := s.reader.Read(slug)
	if err != nil {
		return nil, errors.New("post not found")
	}

	var meta PostMetadata
	remaining, err := frontmatter.Parse(bytes.NewReader(raw), &meta)
	if err != nil {
		s.logger.Error("error parsing frontmatter", "slug", slug, "error", err)
		return nil, errors.New("error parsing frontmatter")
	}
	post := Post{PostMetadata: meta}

	// Parse the markdown into an AST so we can walk it before rendering.
	reader := text.NewReader(remaining)
	doc := s.markdown.Parser().Parse(reader)

	post.TOC = extractTOC(doc, remaining)
	post.ReadingTime = readingTime(remaining)

	var buf bytes.Buffer
	if err = s.markdown.Renderer().Render(&buf, remaining, doc); err != nil {
		return nil, err
	}
	post.Content = template.HTML(buf.String())

	s.mu.Lock()
	s.cache[slug] = &post
	s.mu.Unlock()

	return &post, nil
}

func (s *Service) QueryMetadata() (*PostMetadataCollection, error) {
	return s.reader.Query()
}

// QueryFiltered returns a paginated, tag-filtered, and text-searched PostIndexView.
func (s *Service) QueryFiltered(tag, search string, page, pageSize int) (*PostIndexView, error) {
	all, err := s.reader.Query()
	if err != nil {
		return nil, err
	}

	// collect unique tags from ALL posts (not just filtered)
	seen := map[string]bool{}
	var tags []string
	for _, p := range all.Posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}

	// filter by tag
	filtered := all.Posts
	if tag != "" {
		filtered = filtered[:0:0]
		for _, p := range all.Posts {
			for _, t := range p.Tags {
				if t == tag {
					filtered = append(filtered, p)
					break
				}
			}
		}
	}

	// filter by search (case-insensitive title + description match)
	if search != "" {
		q := strings.ToLower(search)
		var matched []PostMetadata
		for _, p := range filtered {
			if strings.Contains(strings.ToLower(p.Title), q) ||
				strings.Contains(strings.ToLower(p.Description), q) {
				matched = append(matched, p)
			}
		}
		filtered = matched
	}

	total := len(filtered)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	return &PostIndexView{
		Posts:      filtered[start:end],
		Tags:       tags,
		ActiveTag:  tag,
		Search:     search,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

// extractTOC walks the goldmark AST and collects heading nodes into a TOC slice.
func extractTOC(doc gast.Node, source []byte) []TOCItem {
	var toc []TOCItem
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		h, ok := n.(*gast.Heading)
		if !ok {
			return gast.WalkContinue, nil
		}
		id := ""
		if v, ok := h.AttributeString("id"); ok {
			if b, ok := v.([]byte); ok {
				id = string(b)
			}
		}
		toc = append(toc, TOCItem{
			Level: h.Level,
			Text:  headingText(h, source),
			ID:    id,
		})
		return gast.WalkSkipChildren, nil
	})
	return toc
}

// headingText extracts the plain text content from a heading node.
func headingText(h *gast.Heading, source []byte) string {
	var buf bytes.Buffer
	for c := h.FirstChild(); c != nil; c = c.NextSibling() {
		switch t := c.(type) {
		case *gast.Text:
			buf.Write(t.Segment.Value(source))
		case *gast.String:
			buf.Write(t.Value)
		}
	}
	return buf.String()
}

// readingTime estimates minutes to read based on an average of 200 wpm.
func readingTime(src []byte) int {
	words := 0
	inWord := false
	for _, r := range string(src) {
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	mins := words / 200
	if mins < 1 {
		mins = 1
	}
	return mins
}
