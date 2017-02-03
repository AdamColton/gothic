package gothicgo

import (
	"fmt"
	"sort"
	"strings"
)

// Imports is a tool for managing imports. Imports can be defined by path or
// package and either may include an alias. The ResolvePackages method must be
// called to resolve any packages to paths.
type Imports struct {
	paths map[string]string
	pkgs  map[string]string
}

func NewImports() *Imports {
	return &Imports{
		paths: map[string]string{},
		pkgs:  map[string]string{},
	}
}

func (i *Imports) AddPathImport(path string) {
	if path != "" {
		i.paths[path] = ""
	}
}

func (i *Imports) AddPackageImport(pkg string) {
	if pkg != "" {
		i.pkgs[pkg] = ""
	}
}
func (i *Imports) AddPathAliasImport(path, alias string) {
	if path != "" {
		i.paths[path] = alias
	}
}
func (i *Imports) AddPackageAliasImport(pkg, alias string) {
	if pkg != "" {
		i.pkgs[pkg] = alias
	}
}
func (i *Imports) AddImports(imports *Imports) {
	for pkg, alias := range imports.pkgs {
		i.pkgs[pkg] = alias
	}
	for path, alias := range imports.paths {
		i.paths[path] = alias
	}
}

func (i *Imports) RemovePath(path string) {
	delete(i.paths, path)
}

func (i *Imports) ResolvePackages(resolver ImportResolver) {
	for pkg, alias := range i.pkgs {
		if path := resolver.Resolve(pkg); path != "" {
			i.paths[path] = alias
		}
	}
}

func (i *Imports) String() string {
	ln := len(i.paths)
	if ln == 0 {
		return ""
	}
	l := make([]string, ln+3)
	l[0] = "import ("
	j := 1
	for path, alias := range i.paths {
		if alias == "" {
			l[j] = fmt.Sprintf("\t\"%s\"", path)
		} else {
			l[j] = fmt.Sprintf("\t%s \"%s\"", alias, path)
		}
		j++
	}
	l[j] = ")"
	sort.Strings(l[1 : ln+1])
	return strings.Join(l, "\n")
}
