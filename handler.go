package goimport

import (
	"fmt"
	"net/http"
	"path"
	"sync"
)

// Handler is the http.Handler
type Handler struct {
	mu    sync.RWMutex
	root  http.Handler
	paths map[string]Package
	srcs  map[PackageFinder]Packages
}

// Handle takes a root http.Handler and returns a wrapper which serves out
// go-import meta tags for the go tool, or redirects to godoc.org when there is
// no ?go-get=1 query. For anything that does not match, your root handler is
// called. If root is nil, we use http.NotFoundHandler.
func Handle(root http.Handler, sources ...PackageFinder) *Handler {
	if root == nil {
		root = http.NotFoundHandler()
	}
	h := &Handler{
		root:  root,
		paths: make(map[string]Package),
		srcs:  make(map[PackageFinder]Packages),
	}
	for _, src := range sources {
		go src.FindPackages(h)
	}
	return h
}

// ServeHTTP serves our redirect, html meta tag, or the root
// http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !h.tryServe(w, req) {
		h.root.ServeHTTP(w, req)
	}
}

func validVCS(vcs string) bool {
	switch vcs {
	case "git", "hg", "bzr", "svn":
		return true
	default:
		return false
	}
}

// SetPackages is called by PackageFinders for updating our set of
// packages.
func (h *Handler) SetPackages(src PackageFinder, pkgs []Package) {
	h.mu.Lock()
	defer h.mu.Unlock()
	prevPkgs, ok := h.srcs[src]
	if ok {
		for _, p := range prevPkgs {
			delete(h.paths, p.Path)
		}
	}
	for i := range pkgs {
		p := &pkgs[i]
		if validVCS(p.VCS) {
			s := path.Clean(p.Path)
			p.Path = s
			h.paths[s] = *p
		} else {
			panic(fmt.Errorf("goimport: unknown vcs %s", p.VCS))
		}
	}
	h.srcs[src] = pkgs
}

type tmplData struct {
	Host string
	Path string
	VCS  string
	Repo string
}

func (h *Handler) tryServe(w http.ResponseWriter, req *http.Request) bool {
	var data tmplData
	query := req.URL.Query()
	_, goget := query["go-get"]
	path := req.URL.Path
	if path[0] == '/' {
		path = path[1:]
	}
	data.Host = req.Host
	if goget && h.goget(&data, path) {
		render(w, data)
		return true
	} else if url, ok := h.godocRedirect(path, req.Host); ok {
		http.Redirect(w, req, url, 302)
		return true
	}
	return false
}

func (h *Handler) goget(data *tmplData, path string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if p, ok := h.paths[path]; ok {
		data.Path = p.RootPath
		data.VCS = p.VCS
		data.Repo = p.TargetURL
		return true
	}
	return false
}

func (h *Handler) godocRedirect(path, host string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if _, ok := h.paths[path]; ok {
		return "http://godoc.org/" + host + "/" + path, true
	}
	return "", false
}
