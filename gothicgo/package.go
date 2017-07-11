package gothicgo

import (
	"os"
	"path"
	"path/filepath"
)

// represents a directory containing Go code.
type Package struct {
	Name       string
	ImportPath string
	OutputPath string
	files      map[string]*File
	structs    map[string]*Struct
	interfaces map[string]*Interface
	resolver   ImportResolver
	Comment    string
}

func NewPackage(name string) *Package {
	pkg := &Package{
		Name:       name,
		ImportPath: path.Join(ImportPath, name),
		OutputPath: path.Join(OutputPath, name),
		files:      make(map[string]*File),
		structs:    make(map[string]*Struct),
		interfaces: make(map[string]*Interface),
		Comment:    DefaultComment,
	}
	packages.AddGenerators(pkg)
	return pkg
}

func (p *Package) Prepare() error {
	if p.Name != "main" {
		p.ImportResolver().Add(p.Name, p.ImportPath)
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
