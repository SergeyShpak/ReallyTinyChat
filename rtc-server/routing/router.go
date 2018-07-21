package routing

import (
	"log"
	"net/http"
)

var defaultHandler http.Handler = http.HandlerFunc(handleUnknown)

type Router struct {
	routes         map[string]*Route
	defaultHandler http.Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

func (r *Router) SetDefaultHandler(h http.Handler) {
	r.defaultHandler = h
}

func (r *Router) newRoute(path string) *Route {
	route := &Route{
		path: path,
	}
	r.routes[path] = route
	return route
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.newRoute(path).Handler(f)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := r.getHandler(req.URL.Path)
	if handler == nil {
		defaultHandler.ServeHTTP(w, req)
		return
	}
	handler.ServeHTTP(w, req)
}

func (r *Router) getHandler(path string) http.Handler {
	route := r.routes[path]
	if route == nil {
		return nil
	}
	return route.handler
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	log.Println("Default handler")
}
