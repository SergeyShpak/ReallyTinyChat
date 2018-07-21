package routing

import "net/http"

type Route struct {
	path    string
	handler http.Handler
}

func (r *Route) Handler(f func(http.ResponseWriter, *http.Request)) *Route {
	r.handler = http.HandlerFunc(f)
	return r
}
