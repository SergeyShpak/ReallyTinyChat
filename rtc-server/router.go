package main

import (
	"net/http"
)

type Handle func(http.ResponseWriter, *http.Request)

type Router struct {
	routes map[string]Handle
}

func newRouter() *Router {
	return &Router{
		routes: make(map[string]Handle),
	}
}

func (r *Router) Add(path string, handle Handle) {
	r.routes[path] = handle
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head := r.URL.Path
	h, ok := rt.routes[head]
	if ok {
		h(w, r)
		return
	}
	http.NotFound(w, r)
}
