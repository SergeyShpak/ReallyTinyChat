package routing

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRouteHandler(t *testing.T) {
	normalFn := func(http.ResponseWriter, *http.Request) {}
	cases := []struct {
		handlerFn func(http.ResponseWriter, *http.Request)
	}{
		{handlerFn: normalFn},
		{handlerFn: nil},
	}
	t.Parallel()
	for i, tc := range cases {
		i := i
		tc := tc
		t.Run(fmt.Sprintf("TestRouteHandler #%d", i), func(t *testing.T) {
			r := &Route{}
			r.Handler(tc.handlerFn)
		})
	}
}
