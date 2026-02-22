package blog

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Service struct {
	reader   SlugReader
	markdown goldmark.Markdown
}

type SlugReader interface {
	Read(slug string) (string, error) // Updated to return raw content as a string
}

type FileReader struct{}

func NewService(reader SlugReader) *Service {
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
	content, err := s.reader.Read(slug)
	if err != nil {
		return nil, errors.New("post not found")
	}

	var post Post
	remainingContent, err := frontmatter.Parse(bytes.NewBufferString(content), &post)
	if err != nil {
		log.Printf("error parsing frontmatter: %v", err)
		return nil, errors.New("error parsing frontmatter")
	}
	post.Content = template.HTML(remainingContent) // Store the remaining content as template.HTML

	return &post, nil
}

func (s *Service) GetPostBySlugWithMarkdown(slug string) (*Post, error) {
	post, err := s.GetPostBySlug(slug) // Reuse GetPostBySlug to fetch the post
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = s.markdown.Convert([]byte(post.Content), &buf)
	if err != nil {
		return nil, err
	}

	post.Content = template.HTML(buf.String()) // Update the post content with rendered markdown as template.HTML
	return post, nil
}

func (fr *FileReader) Read(slug string) (string, error) {
	f, err := os.Open("posts/" + slug + ".md") // Use relative path to the 'md' folder
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil // Return raw content as a string
}
