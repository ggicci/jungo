package handler

import (
	"net/http"
	"strings"
)

type RESTful map[string]http.Handler

func (rest RESTful) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	method := r.Method

	// Match a handler.
	if h := rest[method]; h != nil {
		h.ServeHTTP(rw, r)
		return
	}

	// Serve 405 error.
	allowed := []string{}
	for k, h := range rest {
		if h != nil {
			allowed = append(allowed, k)
		}
	}
	rw.Header().Set("Allow", strings.Join(allowed, ","))
	http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
