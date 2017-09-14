package gothicgo

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// represents a directory containing Go code.
type Package struct {
	Name       string
	Ref        PackageRef
	OutputPath string
	files      map[string]*File
	structs    map[string]*Struct
	interfaces map[string]*Interface
	resolver   ImportResolver
	Comment    string
}

var nameRe = regexp.MustCompile(`^[\w\-]+$`)
var ErrBadPackageName = errStr("Bad package name")

func NewPackage(name string) (*Package, error) {
	if !nameRe.MatchString(name) {
		return nil, ErrBadPackageName
	}
	pkgRef, err := NewPackageRef(path.Join(importPath, name))
	if err != nil {
		return nil, err
	}
	pkg := &Package{
		Name:       name,
		Ref:        pkgRef,
		OutputPath: path.Join(OutputPath, name),
		files:      make(map[string]*File),
		structs:    make(map[string]*Struct),
		interfaces: make(map[string]*Interface),
		Comment:    DefaultComment,
	}
	packages.AddGenerators(pkg)
	return pkg, nil
}

func (p *Package) Prepare() error {
	if p.Name != "main" {
		p.ImportResolver().Add(p.Ref)
	}
	for _, f := range p.files {
		err := f.Prepare()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Package) Generate() error {
	path, _ := filepath.Abs(p.OutputPath)
	e := os.MkdirAll(path, 0777)
	if e != nil {
		return e
	}
	for _, f := range p.files {
		err := f.Generate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Package) ImportResolver() ImportResolver {
	if p.resolver == nil {
		return AutoResolver()
	}
	return p.resolver
}

func (p *Package) SetResolver(r ImportResolver) { p.resolver = r }

func (p *Package) Export() {
	p.Prepare()
	p.Generate()
}

type packageRef string

func (p packageRef) String() string {
	return string(p)
}

func (p packageRef) Name() string {
	last := strings.LastIndex(string(p), "/")
	if last == -1 {
		return string(p)
	}
	return string(p[last+1:])
}

func (packageRef) private() {}

var packageRefRegex = regexp.MustCompile(`^([\w\-\.]+\/)*[\w\-]+$`)

type PackageRef interface {
	String() string
	Name() string
	// PackageRef is not meant to be implemented, it's meant as an accessor to the
	// underlying packageRef. All instances should be created with NewPackageRef
	// to guarentee that the reference is well formed.
	private()
}

const ErrBadPackageRef = errStr("Bad Package Ref")

func NewPackageRef(ref string) (PackageRef, error) {
	if !packageRefRegex.MatchString(ref) {
		return nil, ErrBadPackageRef
	}
	return packageRef(ref), nil
}

func MustPackageRef(ref string) PackageRef {
	p, err := NewPackageRef(ref)
	if err != nil {
		panic(err)
	}
	return p
}

var pkgBuiltin = packageRef("")

func PkgBuiltin() PackageRef { return pkgBuiltin }
