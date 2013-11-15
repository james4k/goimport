# goimport

Package goimport wraps an http.Handler and provides redirection to
godoc.org or to a remote import path for `?go-get=1` queries. This allows
you to use your domain for import paths while hosting your packages
elsewhere. Read more: http://golang.org/cmd/go/#hdr-Remote_import_path_syntax

Read the [docs](http://j4k.co/goimport), and also note the
[github](http://j4k.co/goimport/github) package.

## Installation

    go get j4k.co/goimport

New to Go? Have a look at how [import paths work](http://golang.org/doc/code.html#remote).
