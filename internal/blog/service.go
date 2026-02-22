package blog

import (
	"errors"
	"io"
	"os"
)

type Service struct {
	reader SlugReader
}

type SlugReader interface {
	Read(slug string) (*Post, error)
}

type FileReader struct{}

func NewService(reader SlugReader) *Service {
	return &Service{reader: reader}
}

func (s *Service) GetPostBySlug(slug string) (*Post, error) {
	post, err := s.reader.Read(slug)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (fr *FileReader) Read(slug string) (*Post, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Here you would read the file and parse it into a Post struct
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &Post{
		Slug:  slug,
		Title: "Title from file",
		Body:  string(b),
	}, nil
}
