package blog

import (
	"bytes"
	"errors"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Service struct {
	reader   PostReader
	markdown goldmark.Markdown
	logger   *slog.Logger
}

type PostReader interface {
	Read(slug string) ([]byte, error)
	Query() (*PostMetadataCollection, error)
}

type FileReader struct {
	dir    string
	logger *slog.Logger
}

func NewFileReader(dir string, logger *slog.Logger) *FileReader {
	return &FileReader{dir: dir, logger: logger.With("component", "file-reader")}
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
	}
}

func (s *Service) GetPostBySlug(slug string) (*Post, error) {
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
	post.Content = template.HTML(remaining)
	return &post, nil
}

func (s *Service) GetPostBySlugWithMarkdown(slug string) (*Post, error) {
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
	return &post, nil
}

func (fr *FileReader) Read(slug string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fr.dir, slug+".md"))
}

func (fr *FileReader) Query() (*PostMetadataCollection, error) {
	filenames, err := filepath.Glob(filepath.Join(fr.dir, "*.md"))
	if err != nil {
		return nil, err
	}

	var collection PostMetadataCollection
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			fr.logger.Error("error opening post file", "file", filename, "error", err)
			continue
		}

		var meta PostMetadata
		_, err = frontmatter.Parse(f, &meta)
		f.Close()
		if err != nil {
			fr.logger.Error("error parsing frontmatter", "file", filename, "error", err)
			continue
		}

		collection.Posts = append(collection.Posts, meta)
	}

	return &collection, nil
}

func (s *Service) QueryMetadata() (*PostMetadataCollection, error) {
	return s.reader.Query()
}
