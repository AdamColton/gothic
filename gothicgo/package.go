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
	resolver   ImportResolver
	Comment    string
}

func NewPackage(name string) *Package {
	pkg := &Package{
		Name:       name,
		ImportPath: path.Join(ImportPath, name),
		OutputPath: path.Join(OutputPath, name),
		files:      map[string]*File{},
		Comment:    DefaultComment,
	}
	packages.AddGenerator(pkg)
	return pkg
}

func (p *Package) Prepare() {
	if p.Name != "main" {
		p.ImportResolver().Add(p.Name, p.ImportPath)
	}
	for _, f := range p.files {
		f.Prepare()
	}
}

func (p *Package) Generate() {
	path, _ := filepath.Abs(p.OutputPath)
	e := os.MkdirAll(path, 0777)
	if e != nil {
		panic(e)
	}
	for _, f := range p.files {
		f.Generate()
	}
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
