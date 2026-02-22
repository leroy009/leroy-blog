package blog

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	service *Service
	tmpls   *Tmpls // Added Tmpls field to store preloaded templates
}

type Tmpls struct {
	PostIndex *template.Template
	Post      *template.Template
}

func NewHandler(service *Service) *Handler {
	postTmpl, err := template.ParseFiles("templates/blog/post.html")
	if err != nil {
		log.Fatalf("error loading template: %v", err)
	}

	indexTmpl, err := template.ParseFiles("templates/blog/index.html")
	if err != nil {
		log.Fatalf("error loading template: %v", err)
	}

	return &Handler{
		service: service,
		tmpls: &Tmpls{
			PostIndex: indexTmpl,
			Post:      postTmpl,
		},
	}
}

func (h *Handler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.service.GetPostBySlugWithMarkdown(slug)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	err = h.tmpls.Post.Execute(&buf, post)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (h *Handler) PostIndexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.QueryMetadata()
	if err != nil {
		http.Error(w, "error fetching posts", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = h.tmpls.PostIndex.Execute(&buf, posts)
	if err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
