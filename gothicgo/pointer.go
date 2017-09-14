package gothicgo

type PointerType interface {
	Type
	Elem() Type
}

type PointerT struct {
	Type
}

func (p *PointerT) Name() string             { return "*" + p.Type.Name() }
func (p *PointerT) String() string           { return p.RelStr(nil) }
func (p *PointerT) RelStr(i *Imports) string { return "*" + p.Type.RelStr(i) }
func (p *PointerT) Kind() Kind               { return PointerKind }
func (p *PointerT) Elem() Type               { return p.Type }

func PointerTo(t Type) PointerType {
	return &PointerT{t}
}
