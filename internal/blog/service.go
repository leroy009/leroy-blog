package blog

import "errors"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetPostBySlug(slug string) (*Post, error) {
	if slug == "hello-world" {
		return &Post{
			Slug:  "hello-world",
			Title: "Hello World",
			Body:  "This is your first post",
		}, nil
	}

	return nil, errors.New("post not found")
}
