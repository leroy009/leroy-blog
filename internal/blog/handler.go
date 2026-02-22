package blog

import (
	"bytes"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

type Handler struct {
	service *Service
	tmpls   *Tmpls
	logger  *slog.Logger
}

type Tmpls struct {
	PostIndex *template.Template
	Post      *template.Template
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	logger = logger.With("component", "handler")

	postTmpl, err := template.ParseFiles("templates/blog/post.html")
	if err != nil {
		logger.Error("error loading template", "template", "post", "error", err)
		os.Exit(1)
	}

	indexTmpl, err := template.ParseFiles("templates/blog/index.html")
	if err != nil {
		logger.Error("error loading template", "template", "index", "error", err)
		os.Exit(1)
	}

	return &Handler{
		service: service,
		logger:  logger,
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
	if err = h.tmpls.Post.Execute(&buf, post); err != nil {
		h.logger.Error("template execution error", "template", "post", "slug", slug, "error", err)
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
	if err = h.tmpls.PostIndex.Execute(&buf, posts); err != nil {
		h.logger.Error("template execution error", "template", "post-index", "error", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
