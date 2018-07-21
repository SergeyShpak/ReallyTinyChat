package routing

import (
	"log"
	"net/http"
)

var defaultHandler http.Handler = http.HandlerFunc(handleUnknown)

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = "GET"
	HttpMethodPost HttpMethod = "POST"
)

type Router struct {
	routes         map[string](map[HttpMethod]*Route)
	defaultHandler http.Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string](map[HttpMethod]*Route)),
	}
}

func (r *Router) SetDefaultHandler(h http.Handler) {
	r.defaultHandler = h
}

func (r *Router) newRoute(path string, method HttpMethod) *Route {
	route := &Route{
		path: path,
	}
	routesMap, ok := r.routes[path]
	if !ok {
		r.routes[path] = make(map[HttpMethod]*Route)
		routesMap = r.routes[path]
	}
	routesMap[method] = route
	return route
}

func (r *Router) HandleFunc(path string, method HttpMethod, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.newRoute(path, method).Handler(f)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := r.getHandler(req.URL.Path, HttpMethod(req.Method))
	if handler == nil {
		defaultHandler.ServeHTTP(w, req)
		return
	}
	handler.ServeHTTP(w, req)
}

func (r *Router) getHandler(path string, method HttpMethod) http.Handler {
	routesMap, ok := r.routes[path]
	if !ok || routesMap == nil {
		return nil
	}
	route, ok := routesMap[method]
	if !ok || route == nil {
		return nil
	}
	return route.handler
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	log.Println("Default handler")
}
