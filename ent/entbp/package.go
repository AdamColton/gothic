package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
)

type Package interface {
	Name() string
	GoPackage() *gothicgo.Package
	Ent(name string) *EntBP
}

type pkg struct {
	name  string
	goPkg *gothicgo.Package
}

func NewPackage(name string) Package {
	return &pkg{
		name: name,
	}
}

func (p *pkg) Name() string { return p.name }

func (p *pkg) GoPackage() *gothicgo.Package {
	if p.goPkg == nil {
		p.goPkg = gothicgo.NewPackage(p.name)
	}

	return p.goPkg
}

func (p *pkg) Ent(name string) *EntBP {
	return &EntBP{
		name:   name,
		fields: map[string]*Field{},
		pkg:    p,
	}
}
