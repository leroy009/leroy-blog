package router

import (
	"net/http"

	"github.com/leroy009/leroy-blog/internal/config"
	"github.com/leroy009/leroy-blog/internal/middleware"
	"mgithub.com/leroy009/leroy-blog/internal/blog"
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

	// Blog wiring
	blogService := blog.NewService()
	blogHandler := blog.NewHandler(blogService)

	r.Get("/health", healthHandler)
	r.Get("/posts/{slug}", blogHandler.GetPost)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
