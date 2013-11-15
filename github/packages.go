package github

import (
	"github.com/google/go-github/github"
	"j4k.co/goimport"
	"log"
	"path"
	"regexp"
	"time"
)

// Packages defines a GitHub package finder for j4k.co/goimport. GitHub
// is polled every hour for changes. API request errors are logged to
// stderr.
type Packages struct {
	User             string
	PollInterval     time.Duration // default is 1 hour
	FilterByName     string
	FilterByHomepage string

	once          bool
	matchName     *regexp.Regexp
	matchHomepage *regexp.Regexp
}

// FindPackages is called by goimport.Handler when setting up a list of
// package paths.
func (p Packages) FindPackages(dest goimport.PackageSetter) {
	if !p.once {
		p.init()
	}
	c := github.NewClient(nil)
	repos, _, err := c.Repositories.List(p.User, nil)
	if err != nil {
		// TODO: user-specified logger?
		log.Println("j4k.co/goimport/github:", err)
		rate, _, err := c.RateLimit()
		if err != nil {
			log.Println("j4k.co/goimport/github:", err)
		} else {
			log.Println("j4k.co/goimport/github:", rate)
		}
		delay := p.PollInterval / 2
		if delay == 0 {
			delay = 10 * time.Minute
		}
		time.AfterFunc(delay, func() { p.FindPackages(dest) })
		return
	}
	var pkgs []goimport.Package
	for _, r := range repos {
		if r.Name == nil || r.HTMLURL == nil {
			continue
		}
		if !p.filterRepo(r) {
			continue
		}
		pkg := goimport.Package{
			VCS:       "git",
			Path:      *r.Name,
			RootPath:  *r.Name,
			TargetURL: *r.HTMLURL,
		}
		pkgs = append(pkgs, pkg)
		pkgs = findSubPackages(pkgs, c, p.User, *r.Name, *r.HTMLURL)
	}
	dest.SetPackages(p, pkgs)
	time.AfterFunc(p.PollInterval, func() { p.FindPackages(dest) })
}

func (p *Packages) init() {
	p.once = true
	if p.PollInterval == 0 {
		p.PollInterval = 1 * time.Hour
	}
	p.initMatchers()
}

func (p *Packages) initMatchers() {
	for _, s := range []struct {
		match  **regexp.Regexp
		filter string
	}{
		{&p.matchName, p.FilterByName},
		{&p.matchHomepage, p.FilterByHomepage},
	} {
		if s.filter != "" {
			*s.match = regexp.MustCompile(s.filter)
		}
	}
}

// filterRepo returns true if we want to use the repository
func (p *Packages) filterRepo(r github.Repository) bool {
	for _, s := range []struct {
		match *regexp.Regexp
		pstr  *string
	}{
		{p.matchName, r.Name},
		{p.matchHomepage, r.Homepage},
	} {
		if s.match == nil {
			continue
		}
		if s.pstr == nil {
			return false
		}
		if !s.match.MatchString(*s.pstr) {
			return false
		}
	}
	return true
}

func findSubPackages(pkgs []goimport.Package, c *github.Client, user, name, url string) []goimport.Package {
	tree, _, err := c.Git.GetTree(user, name, "master", true)
	if err != nil {
		log.Println("j4k.co/goimport/github:", err)
		return pkgs
	}
	set := make(map[string]struct{})
	for _, entry := range tree.Entries {
		if entry.Path == nil {
			continue
		}
		s := *entry.Path
		if path.Ext(s) != ".go" {
			continue
		}
		dir := path.Dir(s)
		if dir == "." || dir == "/" {
			continue
		}
		if _, ok := set[dir]; ok {
			continue
		}
		set[dir] = struct{}{}
		pkg := goimport.Package{
			VCS:       "git",
			Path:      path.Join(name, dir),
			RootPath:  name,
			TargetURL: url,
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
