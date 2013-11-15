package goimport_test

import (
	"j4k.co/goimport"
	"j4k.co/goimport/github"
	"net/http"
)

func ExampleHandle() {
	// root could be your router from gorilla/mux (or your router of
	// choice), or a subhandler if you want to use a path prefix like
	// /pkg/ or /go/.
	root := http.NotFoundHandler()
	http.Handle("/", goimport.Handle(root, goimport.Packages{
		{
			VCS:       "git",
			Path:      "repo1",
			TargetURL: "github.com/username/repo1",
		},
		{
			VCS:       "bzr",
			Path:      "repo2",
			TargetURL: "launchpad.net/repo2",
		},
	}, github.Packages{
		User:             "james4k",
		FilterByHomepage: "(http://)?j4k.co/",
	}))
}
