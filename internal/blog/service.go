package blog

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Service struct {
	reader   PostReader
	markdown goldmark.Markdown
}

type PostReader interface {
	Read(slug string) ([]byte, error)
	Query() (*PostMetadataCollection, error)
}

type FileReader struct {
	dir string
}

func NewFileReader(dir string) *FileReader {
	return &FileReader{dir: dir}
}

func NewService(reader PostReader) *Service {
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
		log.Printf("error parsing frontmatter: %v", err)
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
		log.Printf("error parsing frontmatter: %v", err)
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
			log.Printf("error opening file %s: %v", filename, err)
			continue
		}

		var meta PostMetadata
		_, err = frontmatter.Parse(f, &meta)
		f.Close()
		if err != nil {
			log.Printf("error parsing frontmatter in file %s: %v", filename, err)
			continue
		}

		collection.Posts = append(collection.Posts, meta)
	}

	return &collection, nil
}

func (s *Service) QueryMetadata() (*PostMetadataCollection, error) {
	return s.reader.Query()
}
