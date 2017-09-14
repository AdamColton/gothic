package gothicgo

import (
	"fmt"
	"sort"
	"strings"
)

// Imports is a tool for managing imports. Imports can be defined by path or
// package and either may include an alias. The ResolvePackages method must be
// called to resolve any packages to refs.
type Imports struct {
	self  PackageRef
	refs  map[PackageRef]string
	names map[string]string
}

func NewImports(self PackageRef) *Imports {
	return &Imports{
		self:  self,
		refs:  make(map[PackageRef]string),
		names: make(map[string]string),
	}
}

// Prefix returns the name or alias of the package reference if it is
// different from Imports.self. The name will either be a blank string or will
// end with a period.
func (i *Imports) Prefix(ref PackageRef) string {
	if i != nil && i.self != nil && ref.String() == i.self.String() {
		return ""
	}
	return i.GetRefName(ref) + "."
}

func (i *Imports) AddRefImports(refs ...PackageRef) {
	for _, ref := range refs {
		if ref != nil && ref.String() != "" && (i.self == nil || ref.String() != i.self.String()) {
			if _, exists := i.refs[ref]; !exists {
				i.refs[ref] = ""
			}
		}
	}
}

func (i *Imports) AddNameImports(names ...string) {
	//TODO: check that name is well formed
	for _, name := range names {
		if name != "" {
			i.names[name] = ""
		}
	}
}
func (i *Imports) AddRefAliasImport(ref PackageRef, alias string) {
	if ref != nil && ref.String() != "" && ref.String() != i.self.String() {
		i.refs[ref] = alias
	}
}
func (i *Imports) AddNameAliasImport(name, alias string) {
	//TODO: check that pkg is well formed
	if name != "" && alias != "" {
		i.names[name] = alias
	}
}
func (i *Imports) AddImports(imports *Imports) {
	// TODO: handle alias collision
	for pkg, alias := range imports.names {
		i.names[pkg] = alias
	}
	for path, alias := range imports.refs {
		i.refs[path] = alias
	}
}

func (i *Imports) RemoveRef(ref PackageRef) {
	delete(i.refs, ref)
}

func (i *Imports) ResolvePackages(resolver ImportResolver) {
	// TODO: handle alias collision
	for pkg, alias := range i.names {
		if path := resolver.Resolve(pkg); path.String() != "" {
			i.refs[path] = alias
		}
	}
}

func (i *Imports) GetRefName(ref PackageRef) string {
	if i == nil {
		return ref.Name()
	}
	name, ok := i.refs[ref]
	if ok {
		if name != "" {
			return name
		}
		return ref.Name()
	}
	rn := ref.Name()
	name, ok = i.names[rn]
	if !ok {
		return rn
	}
	delete(i.names, rn)
	i.refs[ref] = name
	if name == "" {
		return rn
	}
	return name
}

func (i *Imports) String() string {
	ln := len(i.refs)
	if ln == 0 {
		return ""
	}
	l := make([]string, ln+3)
	l[0] = "import ("
	j := 1
	for path, alias := range i.refs {
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
