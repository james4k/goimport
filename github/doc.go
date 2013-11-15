/*
Package github finds package paths on GitHub for the goimport package.

In the following code, goimport.Handle will redirect `go get` requests
of j4k.co/somerepo to github.com/james4k/somerepo for all GitHub repos
that are owned by user james4k and with a homepage that starts with
j4k.co.

	goimport.Handle(router, github.Packages{
		User:             "james4k",
		FilterByHomepage: "^(http://)?j4k.co/",
	})

For each repository, any directories that have .go files is considered a
package.
*/
package github
