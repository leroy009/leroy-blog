package blog

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var layout = []string{
	"templates/layout/base.html",
	"templates/layout/header.html",
	"templates/layout/footer.html",
}

type Handler struct {
	service *Service
	tmpls   *Tmpls
	logger  *slog.Logger
}

type Tmpls struct {
	PostIndex *template.Template
	Post      *template.Template
}

var tmplFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"pageURL": func(tag, search string, page int) string {
		v := url.Values{}
		if tag != "" {
			v.Set("tag", tag)
		}
		if search != "" {
			v.Set("q", search)
		}
		v.Set("page", strconv.Itoa(page))
		return "/posts?" + v.Encode()
	},
	"year": func() int { return time.Now().Year() },
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	logger = logger.With("component", "blog-handler")

	parse := func(page string) *template.Template {
		t, err := template.New("").Funcs(tmplFuncs).ParseFiles(append(layout, page)...)
		if err != nil {
			logger.Error("error loading template", "page", page, "error", err)
			os.Exit(1)
		}
		return t
	}

	return &Handler{
		service: service,
		logger:  logger,
		tmpls: &Tmpls{
			PostIndex: parse("templates/blog/index.html"),
			Post:      parse("templates/blog/post.html"),
		},
	}
}

func (h *Handler) render(w http.ResponseWriter, t *template.Template, data any) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		h.logger.Error("template execution error", "error", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (h *Handler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.service.GetPostBySlugWithMarkdown(slug)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	h.render(w, h.tmpls.Post, post)
}

const pageSize = 5

func (h *Handler) PostIndexHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	view, err := h.service.QueryFiltered(tag, search, page, pageSize)
	if err != nil {
		http.Error(w, "error fetching posts", http.StatusInternalServerError)
		return
	}

	h.render(w, h.tmpls.PostIndex, view)
}
