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

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
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

	var post Post
	remaining, err := frontmatter.Parse(bytes.NewReader(raw), &post)
	if err != nil {
		s.logger.Error("error parsing frontmatter", "slug", slug, "error", err)
		return nil, errors.New("error parsing frontmatter")
	}

	var buf bytes.Buffer
	if err = s.markdown.Convert(remaining, &buf); err != nil {
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
