package blog

import "html/template"

type Post struct {
	Author  Author        `json:"author"`
	Date    string        `json:"date"`
	Slug    string        `json:"slug"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"` // Updated to template.HTML for safe HTML rendering
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
