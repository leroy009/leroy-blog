package blog

type Post struct {
	Author  string `json:"author"`
	Date    string `json:"date"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
