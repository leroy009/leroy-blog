package router

import (
	"net/http"

	"github.com/leroy009/leroy-blog/internal/blog"
	"github.com/leroy009/leroy-blog/internal/config"
	"github.com/leroy009/leroy-blog/internal/middleware"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []Middleware
}

func New(cfg config.Config) http.Handler {
	r := &Router{
		mux: http.NewServeMux(),
	}

	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	r.Get("/health", healthHandler)

	// Blog wiring
	fr := blog.NewFileReader(cfg.PostsDir)
	blogService := blog.NewService(fr)
	blogHandler := blog.NewHandler(blogService)

	r.Group("/posts", func(posts *Router) {
		posts.Get("/", blogHandler.PostIndexHandler)
		posts.Get("/{slug}", blogHandler.GetPostHandler)
	})

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Group creates a sub-router with a shared prefix, inheriting current middleware.
func (r *Router) Group(prefix string, fn func(*Router)) {
	sub := &Router{
		mux:         r.mux,
		prefix:      r.prefix + prefix,
		middlewares: append([]Middleware{}, r.middlewares...),
	}
	fn(sub)
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) handle(method, pattern string, handler http.Handler) {
	fullPattern := method + " " + r.prefix + pattern

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	r.mux.Handle(fullPattern, handler)
}

func (r *Router) Get(pattern string, h http.HandlerFunc) {
	r.handle("GET", pattern, h)
}

func (r *Router) Post(pattern string, h http.HandlerFunc) {
	r.handle("POST", pattern, h)
}

func (r *Router) Put(pattern string, h http.HandlerFunc) {
	r.handle("PUT", pattern, h)
}

func (r *Router) Delete(pattern string, h http.HandlerFunc) {
	r.handle("DELETE", pattern, h)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
