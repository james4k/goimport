package goimport

// PackageFinder is used to find new package paths.
type PackageFinder interface {
	FindPackages(PackageSetter)
}

// PackageSetter is used by PackageFinder to update a set of packages.
type PackageSetter interface {
	SetPackages(PackageFinder, []Package)
}

// Package represents a package's set of paths to locate its code
// repository.
type Package struct {
	VCS       string // git, hg, bzr, or svn
	Path      string // http path
	RootPath  string // root import path. optional; defaults to Path
	TargetURL string // target import path we forward people to
}

// Packages is a slice of Package's that implements PackageFinder
type Packages []Package

// FindPackages is called by goimport.Wrap when setting up a list of
// package paths.
func (p Packages) FindPackages(dest PackageSetter) {
	p = p.defaults()
	dest.SetPackages(p, p)
}

func (p Packages) defaults() Packages {
	for i := range p {
		if p[i].RootPath == "" {
			p[i].RootPath = p[i].Path
		}
	}
	return p
}
