package router

import (
	"net/http"
	"slices"
)

type Router struct {
	globalChain []func(http.Handler) http.Handler
	routeChain  []func(http.Handler) http.Handler
	isSubRouter bool
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	if r.isSubRouter {
		r.routeChain = append(r.routeChain, mw...)
	} else {
		r.globalChain = append(r.globalChain, mw...)
	}
}

func (r *Router) Group(fn func(r *Router)) {
	subRouter := &Router{routeChain: slices.Clone(r.routeChain), isSubRouter: true, ServeMux: r.ServeMux}
	fn(subRouter)
}

func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}
	r.ServeMux.Handle(pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pattern := r.ServeMux.Handler(req)

	var h http.Handler = handler
	for _, mw := range slices.Backward(r.globalChain) {
		h = mw(h)
	}

	if pattern == "" {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not found", http.StatusNotFound)
		})
	}

	h.ServeHTTP(w, req)
}
