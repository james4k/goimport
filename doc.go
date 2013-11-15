/*
Package goimport provides redirection to godoc.org or to a remote import
path for `go get` requests. This allows you to use your domain for
import paths while hosting your packages elsewhere. Read more:
http://golang.org/cmd/go/#hdr-Remote_import_path_syntax

	goimport.Handle(otherHandler, goimport.Packages{
		{
			VCS:       "git",
			Path:      "somerepo1",
			TargetURL: "github.com/someuser/somerepo",
		},
	})

With the above, goimport.Handle will redirect `go get` requests of
example.com/somerepo to github.com/someuser/somerepo. For any paths that
do not match a package, requests fall back to otherHandler.

The go tool always passes the `?go-get=1` query in its request. When the
query is not found, the HTTP client is given a 302 redirect to
godoc.org/{host}/{prefix}/basepath.
*/
package goimport
