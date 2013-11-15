/*
Package github finds package paths on GitHub for the goimport package.

	goimport.Handle(otherHandler, github.Packages{
		User:             "someuser",
		FilterByHomepage: "^(http://)?example.com/",
	})

With the above, goimport.Handle will redirect `go get` requests of
example.com/somerepo to github.com/someuser/somerepo for all GitHub
repos that are owned by user someuser and with a homepage at
example.com. For any paths that do not match a package, requests fall
back to otherHandler.

For each repository, any directories that have .go files is considered a
package.
*/
package github
