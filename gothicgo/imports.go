package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"sort"
)

// Prefixer takes a PackageRef and returns the correct prefix for it. If the
// reference is to the same pacakge we are in, it will return an empty string.
// If it's a package imported normally. it will return the package name followed
// by a period. If it is an aliased package, it will return the alias followed
// by a period.
type Prefixer interface {
	Prefix(ref PackageRef) string
}

// DefaultPrefixer always returns the package prefix.
var DefaultPrefixer = defaultPrefixer{}

type defaultPrefixer struct{}

func (defaultPrefixer) Prefix(ref PackageRef) string {
	if ref.Name() == "" {
		return ""
	}
	return ref.Name() + "."
}

// Imports is a tool for managing imports. Imports can be defined by path or
// package and either may include an alias. The ResolvePackages method must be
// called to resolve any packages to refs.
type Imports struct {
	self  PackageRef
	refs  map[string]string
	names map[string]string
}

// NewImports sets up an instance of Imports.
func NewImports(self PackageRef) *Imports {
	return &Imports{
		self:  self,
		refs:  make(map[string]string),
		names: make(map[string]string),
	}
}

// Prefix returns the name or alias of the package reference if it is
// different from Imports.self. The name will either be a blank string or will
// end with a period.
func (i *Imports) Prefix(ref PackageRef) string {
	if (i != nil && i.self != nil && ref.String() == i.self.String()) || ref.Name() == "" {
		return ""
	}
	return i.GetRefName(ref) + "."
}

// AddRefImports takes PackageRefs and adds them as imports without aliases.
func (i *Imports) AddRefImports(refs ...PackageRef) {
	for _, ref := range refs {
		if ref != nil {
			rs := ref.String()
			if rs != "" && (i.self == nil || rs != i.self.String()) {
				if _, exists := i.refs[rs]; !exists {
					i.refs[rs] = ""
				}
			}
		}
	}
}

// AddNameImports takes package names. They will be resolved to full PackageRefs
// when ResolvePackages is called.
func (i *Imports) AddNameImports(names ...string) {
	//TODO: check that name is well formed
	for _, name := range names {
		if name != "" {
			i.names[name] = ""
		}
	}
}

// AddRefAliasImport adds a PackageRef as an alias
func (i *Imports) AddRefAliasImport(ref PackageRef, alias string) {
	if ref != nil {
		rs := ref.String()
		if rs != "" && rs != i.self.String() {
			i.refs[rs] = alias
		}
	}
}

// AddNameAliasImport adds a package by name with an alias. The name will be
// resolved when ResolvePackages is called.
func (i *Imports) AddNameAliasImport(name, alias string) {
	//TODO: check that pkg is well formed
	if name != "" && alias != "" {
		i.names[name] = alias
	}
}

// AddImports takes another instance of Imports and adds all it's imports. This
// runs the risk of clobbering aliases.
func (i *Imports) AddImports(imports *Imports) {
	// TODO: handle alias collision
	for pkg, alias := range imports.names {
		i.names[pkg] = alias
	}
	for path, alias := range imports.refs {
		i.refs[path] = alias
	}
}

// RemoveRef removes a reference.
func (i *Imports) RemoveRef(ref PackageRef) {
	if ref != nil {
		delete(i.refs, ref.String())
	}
}

// ResolvePackages uses a resolver to find all the packages that were imported
// by name.
func (i *Imports) ResolvePackages(resolver ImportResolver) {
	// TODO: handle alias collision
	for pkg, alias := range i.names {
		path := resolver.Resolve(pkg)
		if ps := path.String(); ps != "" {
			i.refs[ps] = alias
		}
	}
}

// GetRefName takes a package ref and returns the name it will be referenced by
// in the Import context. If the package is aliased it will return the alias,
// otherwise it will return the package name. If there is an unresolved name
// matching the PackageRef, it will be treated as resolving to the ref.
func (i *Imports) GetRefName(ref PackageRef) string {
	if i == nil {
		return ref.Name()
	}
	name, ok := i.refs[ref.String()]
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
	i.refs[ref.String()] = name
	if name == "" {
		return rn
	}
	return name
}

// String returns the imports as Go code.
func (i *Imports) String() string {
	buf := bufpool.Get()
	i.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the Go code to a writer
func (i *Imports) WriteTo(w io.Writer) (int64, error) {
	ln := len(i.refs)
	if ln == 0 {
		return 0, nil
	}
	sum := gothicio.NewSumWriter(w)
	sum.WriteString("import (")

	refs := make([]string, 0, len(i.refs))
	for path := range i.refs {
		refs = append(refs, path)
	}
	sort.Strings(refs)

	for _, path := range refs {
		sum.WriteString("\n\t")
		if alias := i.refs[path]; alias != "" {
			sum.WriteString(alias)
			sum.WriteString(" ")
		}
		sum.WriteString("\"")
		sum.WriteString(path)
		sum.WriteString("\"")
	}
	sum.WriteString("\n)\n")
	if sum.Err != nil {
		sum.Err = errCtx(sum.Err, "While writing imports:")
	}
	return sum.Rets()
}
