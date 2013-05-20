/*
Package goimport wraps your root http.Handler to provide redirection to
godoc.org or to a remote import path for ?go-get=1 queries. This allows you to
use your domain for import paths while hosting your packages elsewhere. Read
more: http://golang.org/cmd/go/#hdr-Remote_import_path_syntax

Incomplete, API could change, not tested well, etc.
TODO: Some kind of versioning would be swell...

In the following code, the "pkg" prefix is used to direct the go tool from
example.com/pkg/somerepo to github.com/someuser/somerepo. You can also use an
empty string as the prefix, or put multiple levels of paths.

	goimport.Wrap(router, "pkg", goimport.Paths{
		{"git", "github.com/someuser/somerepo"},
		{"bzr", "launchpad.net/anotherepo"},
	})

Note that the go tool always passes the ?go-get=1 query in the URL. When the
query is not found, the HTTP client is given a 302 redirect to
godoc.org/{host}/{prefix}/basepath. For all paths that do not match the repo
paths, the passed in http.Handler is called.
*/
package goimport
