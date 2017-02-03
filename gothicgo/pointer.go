package gothicgo

type PointerType interface {
	Type
	Elem() Type
}

type PointerT struct {
	Type
}

func (p *PointerT) Name() string             { return "*" + p.Type.Name() }
func (p *PointerT) String() string           { return p.RelStr("") }
func (p *PointerT) RelStr(pkg string) string { return "*" + p.Type.RelStr(pkg) }
func (p *PointerT) Kind() Kind               { return PointerKind }
func (p *PointerT) Elem() Type               { return p.Type }

func PointerTo(t Type) PointerType {
	return &PointerT{t}
}
