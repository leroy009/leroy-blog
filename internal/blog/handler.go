package blog

import (
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	service *Service
	tmpl    *template.Template // Preloaded template
}

func NewHandler(service *Service) *Handler {
	tmpl, err := template.ParseFiles("templates/blog/post.html")
	if err != nil {
		log.Fatalf("error loading template: %v", err)
	}

	return &Handler{
		service: service,
		tmpl:    tmpl, // Store the preloaded template
	}
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.service.GetPostBySlugWithMarkdown(slug)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	err = h.tmpl.Execute(w, post) // Reuse the preloaded template
	if err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
