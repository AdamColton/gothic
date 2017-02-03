package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
)

// Package applies the Go concept of a package to Entities.
type Package interface {
	Name() string
	GoPackage() *gothicgo.Package
	Ent(name string) *EntBP
}

type pkg struct {
	name  string
	goPkg *gothicgo.Package
}

// NewPackage creates a new Package by name
func NewPackage(name string) Package {
	return &pkg{
		name: name,
	}
}

// Name gets the name of the package
func (p *pkg) Name() string { return p.name }

// GoPackage returns a reference to the Go package that corresponds to the
// Entity package.
func (p *pkg) GoPackage() *gothicgo.Package {
	if p.goPkg == nil {
		p.goPkg = gothicgo.NewPackage(p.name)
	}

	return p.goPkg
}

// Ent creates an Entity in the package by name.
func (p *pkg) Ent(name string) *EntBP {
	return &EntBP{
		name:   name,
		fields: map[string]*Field{},
		pkg:    p,
	}
}
