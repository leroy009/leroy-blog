package blog

import (
	"html/template"
	"time"
)

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PostMetadata struct {
	Author      Author    `json:"author"`
	Date        time.Time `json:"date"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
}

// TOCItem represents a single heading in the table of contents.
type TOCItem struct {
	Level int
	Text  string
	ID    string
}

type Post struct {
	PostMetadata
	Content     template.HTML `json:"content"`
	TOC         []TOCItem
	ReadingTime int // estimated minutes to read
}

type PostMetadataCollection struct {
	Posts []PostMetadata
}

// PostIndexView is the data passed to the blog index template.
type PostIndexView struct {
	Posts      []PostMetadata
	Tags       []string // all unique tags across all posts
	ActiveTag  string   // currently filtered tag, empty means all
	Search     string   // current search query
	Page       int
	TotalPages int
}
