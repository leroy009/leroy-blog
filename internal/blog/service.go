package blog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"

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

	return &Post{
		Author:  "Leroy",
		Date:    time.Now().Format("January 2, 2006"),
		Slug:    slug,
		Title:   "Title from file",
		Content: content,
	}, nil
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

	post.Content = buf.String() // Update the post content with rendered markdown
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
