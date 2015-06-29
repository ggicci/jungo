package handler

import (
	"net/http"
	"strings"
)

type RESTful map[string]http.Handler

func (rest RESTful) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Match a handler.
	if h := rest.handler(r.Method); h != nil {
		h.ServeHTTP(rw, r)
		return
	}

	// Serve 405 error.
	rw.Header().Set("Allow", strings.Join(rest.allowedMethods(), ","))
	http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rest RESTful) handler(method string) http.Handler {
	handler := rest[method]
	if handler != nil {
		return handler
	}
	if method == "HEAD" && rest["GET"] != nil {
		return rest["GET"]
	}
	return nil
}

func (rest RESTful) allowedMethods() []string {
	allowedMethodsMap := map[string]struct{}{
		"HEAD": struct{}{},
	}
	for method, handler := range rest {
		if handler != nil {
			allowedMethodsMap[method] = struct{}{}
		}
	}
	allowedMethods := []string{}
	for method, _ := range allowedMethodsMap {
		allowedMethods = append(allowedMethods, method)
	}
	return allowedMethods
}
