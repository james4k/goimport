package goimport_test

import (
	"j4k.co/goimport"
	"net/http"
)

func ExampleWrap() {
	// root could be your router from gorilla/mux (or your router of choice)
	root := http.NotFoundHandler()
	http.Handle("/", goimport.Wrap(root, "/", goimport.Paths{
		{"git", "github.com/username/repo1"},
		{"git", "github.com/username/repo2"},
		{"bzr", "launchpad.net/lpadrepo"},
	}))
}
