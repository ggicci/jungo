package mux

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/context"
)

func RouteVars(r *http.Request) map[string]string {
	if rv := context.Get(r, 0); rv != nil {
		return map[string]string(rv.(RouteVariables))
	}

	return nil
}

func setRouteVars(r *http.Request, rvs RouteVariables) {
	fmt.Printf("setVars(%p, %v)\n", r, rvs)
	context.Set(r, 0, rvs)
}

type Mux struct {
	router          *Router
	NotFoundHandler http.Handler
}

func NewMux() *Mux {
	return &Mux{
		router: NewRouter(),
	}
}

func (m *Mux) Handle(pattern string, handler http.Handler) {
	err := m.router.Handle(pattern, handler)
	if err != nil {
		panic(err)
	}
}

func (m *Mux) HandleFunc(pattern string, fn http.HandlerFunc) {
	m.Handle(pattern, http.HandlerFunc(fn))
}

func (m *Mux) HandleStaticFile(pattern, filename string) {
	m.HandleFunc(pattern, func(rw http.ResponseWriter, r *http.Request) {
		http.ServeFile(rw, r, filename)
	})
}

// Serve a directory as a static file server.
// e.g. mux.HandleStaticDir("/assets/", "./assets")
func (m *Mux) HandleStaticDir(prefix, dir string) {
	m.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h, rvs := m.Handler(r)

	if rvs != nil && len(rvs) > 0 {
		setRouteVars(r, rvs)
		h = context.ClearHandler(h)
	}

	h.ServeHTTP(rw, r)
}

func (m *Mux) Handler(r *http.Request) (http.Handler, RouteVariables) {
	if r.RequestURI == "*" {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.ProtoAtLeast(1, 1) {
				rw.Header().Set("Connection", "close")
			}
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}), nil
	}

	if r.Method == "CONNECT" {
		np := cleanPath(r.URL.Path)
		if np != r.URL.Path {
			url := *r.URL
			url.Path = np
			return http.RedirectHandler(url.String(), http.StatusMovedPermanently), nil
		}
	}

	np := r.URL.Path // normalized path

	mr := m.router.Match(np)

	if mr.Handler != nil {
		// Found a matched handler.
		return mr.Handler.(http.Handler), mr.RouteVars
	}

	notFoundHandler := m.NotFoundHandler
	if notFoundHandler == nil {
		notFoundHandler = http.NotFoundHandler()
	}

	// `ssp`, strict slash path.
	// "/a/b" redirects to "/a/b/". "/a/b/" never redirects to "/a/b".
	ssp := ""
	if len(mr.HandlersOnTheWay) > 0 {
		ssp = mr.HandlersOnTheWay[len(mr.HandlersOnTheWay)-1].Path
		if ssp == np+"/" {
			return http.RedirectHandler(ssp, 302), nil
		}

		// Fallback to the most right handler (has "/" suffix) matched on the way.
		for i := len(mr.HandlersOnTheWay) - 1; i >= 0; i-- {
			if !strings.HasSuffix(mr.HandlersOnTheWay[i].Path, "/") {
				continue
			}
			return mr.HandlersOnTheWay[i].Handler.(http.Handler), mr.RouteVars
		}

		// 404
		return notFoundHandler, nil
	}

	// 404
	return notFoundHandler, nil
}

func (m *Mux) GetInternalRouter() *Router { return m.router }

func (m *Mux) DumpRouter() string { return m.router.DumpTree() }

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// Need the "/" postfix? I don't want it :)
	// Need it!!! Shit, 2014.10.08.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}
