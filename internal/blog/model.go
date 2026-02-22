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

type Post struct {
	PostMetadata
	Content template.HTML `json:"content"`
}

type PostMetadataCollection struct {
	Posts []PostMetadata
}
