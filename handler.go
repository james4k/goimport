package goimport

import (
	"fmt"
	"net/http"
	"path"
)

type Path struct {
	VCS  string // git, hg, bzr, or svn
	Path string // remote import path
}

type Paths []Path

type GoImports struct {
	root   http.Handler
	prefix string
	paths  map[string]Path
}

// Wrap takes your root http.Handler and returns a wrapper which serves out
// go-import meta tags for the go tool, or redirects to godoc.org when there is
// no ?go-get=1 query. For anything that does not match, your root handler is
// called
func Wrap(root http.Handler, prefix string, paths Paths) *GoImports {
	prefix = path.Clean(prefix)
	if len(prefix) > 0 && prefix[0] == '/' {
		prefix = prefix[1:]
	}
	h := &GoImports{
		root:   root,
		prefix: prefix,
		paths:  make(map[string]Path),
	}
	for _, p := range paths {
		switch p.VCS {
		case "git", "hg", "bzr", "svn":
			s := path.Join(h.prefix, path.Base(p.Path))
			h.paths[s] = p
		default:
			panic(fmt.Errorf("goimport: unknown vcs %s", p.VCS))
		}
	}
	return h
}

func (h *GoImports) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !h.tryServe(w, req) {
		h.root.ServeHTTP(w, req)
	}
}

func (h *GoImports) tryServe(w http.ResponseWriter, req *http.Request) bool {
	var data struct {
		Host    string
		Imports []Path
	}
	query := req.URL.Query()
	_, goget := query["go-get"]
	path := req.URL.Path
	if path[0] == '/' {
		path = path[1:]
	}
	data.Host = req.Host
	if goget {
		if path == "" {
			for _, p := range h.paths {
				data.Imports = append(data.Imports, p)
			}
			render(w, data)
			return true
		} else if p, ok := h.paths[path]; ok {
			data.Imports = append(data.Imports, p)
			render(w, data)
			return true
		}
	} else {
		if _, ok := h.paths[path]; ok {
			http.Redirect(w, req, "http://godoc.org/"+req.Host+"/"+path, 302)
			return true
		}
	}
	return false
}
