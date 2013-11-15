/*
Package goimport wraps an http.Handler and provides redirection to
godoc.org or to a remote import path for ?go-get=1 queries. This allows
you to use your domain for import paths while hosting your packages
elsewhere. Read more:
http://golang.org/cmd/go/#hdr-Remote_import_path_syntax

In the following code, goimport.Handle will redirect `go get` requests
of example.com/somerepo to github.com/someuser/somerepo.

	goimport.Handle(router, goimport.Packages{
		{
			VCS:       "git",
			Path:      "somerepo1",
			TargetURL: "github.com/someuser/somerepo",
		},
	})

Note that the go tool always passes the ?go-get=1 query in the URL. When the
query is not found, the HTTP client is given a 302 redirect to
godoc.org/{host}/{prefix}/basepath. For all paths that do not match the repo
paths, the passed in http.Handler is used.
*/
package goimport
