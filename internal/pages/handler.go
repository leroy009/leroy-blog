package pages

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var layout = []string{
	"templates/layout/base.html",
	"templates/layout/header.html",
	"templates/layout/footer.html",
}

type Handler struct {
	logger *slog.Logger
	tmpls  *Tmpls
}

type Tmpls struct {
	Home    *template.Template
	About   *template.Template
	Contact *template.Template
}

func NewHandler(logger *slog.Logger) *Handler {
	logger = logger.With("component", "pages-handler")

	parse := func(page string) *template.Template {
		t, err := template.ParseFiles(append(layout, page)...)
		if err != nil {
			logger.Error("error loading template", "page", page, "error", err)
			os.Exit(1)
		}
		return t
	}

	return &Handler{
		logger: logger,
		tmpls: &Tmpls{
			Home:    parse("templates/pages/home.html"),
			About:   parse("templates/pages/about.html"),
			Contact: parse("templates/pages/contact.html"),
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

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.render(w, h.tmpls.Home, nil)
}

func (h *Handler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	h.render(w, h.tmpls.About, nil)
}

func (h *Handler) ContactHandler(w http.ResponseWriter, r *http.Request) {
	h.render(w, h.tmpls.Contact, nil)
}
