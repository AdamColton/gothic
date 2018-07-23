package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// Package represents a directory containing Go code. Package also fulfills the
// PackageRef interface.
type Package struct {
	name       string
	importPath string
	OutputPath string
	files      map[string]*File
	structs    map[string]*Struct
	interfaces map[string]*Interface
	resolver   ImportResolver
	Comment    string
}

var nameRe = regexp.MustCompile(`^[\w\-]+$`)

// ErrBadPackageName is returned when a package name is not allowed
const ErrBadPackageName = errStr("Bad package name")

// NewPackage creates a new Package. The import path will use the ImportPath
// set on the project.
func NewPackage(name string) (*Package, error) {
	if !nameRe.MatchString(name) {
		return nil, ErrBadPackageName
	}
	pkg := &Package{
		name:       name,
		importPath: importPath,
		OutputPath: path.Join(OutputPath, name),
		files:      make(map[string]*File),
		structs:    make(map[string]*Struct),
		interfaces: make(map[string]*Interface),
		Comment:    DefaultComment,
	}
	packages.AddGenerators(pkg)
	return pkg, nil
}

// MustPackage calls NewPackage and panics if there is an error
func MustPackage(name string) *Package {
	pkg, err := NewPackage(name)
	panicOnErr(err)
	return pkg
}

// Prepare calls prepare on all files
func (p *Package) Prepare() error {
	if p.name != "main" {
		p.ImportResolver().Add(p)
	}
	for _, f := range p.files {
		err := f.Prepare()
		if err != nil {
			return errCtx(err, "Prepare package %s", p.name)
		}
	}
	return nil
}

// Generate calls generate on all file
func (p *Package) Generate() error {
	path, _ := filepath.Abs(p.OutputPath)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return errCtx(err, "Generate package %s", p.name)
	}
	for _, f := range p.files {
		err := f.Generate()
		if err != nil {
			return errCtx(err, "Generate package %s", p.name)
		}
	}
	return nil
}

// ImportResolver gets the resolver being used for the package. If no resolver
// is set, the AutoResolver is used.
func (p *Package) ImportResolver() ImportResolver {
	if p.resolver == nil {
		return AutoResolver()
	}
	return p.resolver
}

// SetResolver used for the package.
func (p *Package) SetResolver(r ImportResolver) { p.resolver = r }

// SetImportPath sets the import path for the package not including the name
func (p *Package) SetImportPath(path string) error {
	if !importPathRe.MatchString(path) {
		return ErrBadImportPath
	}
	importPath = path
	return nil
}

// String returns the package import and fulfills the PackageRef and Type
// interfaces.
func (p *Package) String() string {
	return p.importPath + p.name
}

// Name returns the package name and fulfills the PackageRef and Type
// interfaces.
func (p *Package) Name() string {
	return p.name
}

func (*Package) privatePkgRef() {}

type packageRef struct {
	path, name string
}

func (p *packageRef) String() string {
	return p.path
}

func (p *packageRef) Name() string {
	return p.name
}

func (*packageRef) privatePkgRef() {}

// TODO: this regex is only mostly right
var packageRefRegex = regexp.MustCompile(`^(?:[\w\-\.]+\/)*([\w\-]+)$`)

// PackageRef represents a reference to a package.
type PackageRef interface {
	String() string
	Name() string
	// PackageRef is not meant to be implemented, it's meant as an accessor to the
	// underlying packageRef. All instances should be created with NewPackageRef
	// to guarentee that the reference is well formed.
	privatePkgRef()
}

// ErrBadPackageRef indicates a poorly formatted package ref string.
const ErrBadPackageRef = errStr("Bad Package Ref")

// NewPackageRef takes the string used to import a pacakge and returns a
// PackageRef.
func NewPackageRef(ref string) (PackageRef, error) {
	m := packageRefRegex.FindStringSubmatch(ref)
	if len(m) == 0 {
		return nil, ErrBadPackageRef
	}
	return &packageRef{
		path: m[0],
		name: m[1],
	}, nil
}

// MustPackageRef returns a new PackageRef and panics if there is an error
func MustPackageRef(ref string) PackageRef {
	p, err := NewPackageRef(ref)
	panicOnErr(err)
	return p
}

var pkgBuiltin = &packageRef{}

// PkgBuiltin is the PackageRef for the builtin types
func PkgBuiltin() PackageRef { return pkgBuiltin }

// PackageVarRef represents a package variable
type PackageVarRef interface {
	Name() string
	String() string
	PrefixWriteTo(io.Writer, Prefixer) (int64, error)
	PackageRef() PackageRef
	Type() Type
}

// NewPackageVarRef returns a reference to a package variable. Kind is not
// used internally, so if it is not needed, it is safe for it to be nil.
func NewPackageVarRef(pkg PackageRef, name string, kind Type) PackageVarRef {
	return &packageVarRef{
		pkg:  pkg,
		name: name,
		kind: kind,
	}
}

type packageVarRef struct {
	pkg  PackageRef
	name string
	kind Type
}

func (pv *packageVarRef) Name() string   { return pv.name }
func (pv *packageVarRef) String() string { return DefaultPrefixer.Prefix(pv.pkg) + pv.name }
func (pv *packageVarRef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(pre.Prefix(pv.pkg))
	sw.WriteString(pv.name)
	return sw.Rets()
}
func (pv *packageVarRef) PackageRef() PackageRef { return pv.pkg }
func (pv *packageVarRef) Type() Type             { return pv.kind }
