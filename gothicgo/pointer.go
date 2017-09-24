package gothicgo

// PointerType extends the Type interface with pointer specific information
type PointerType interface {
	Type
	Elem() Type
}

type pointerT struct {
	Type
}

func (p *pointerT) Name() string               { return "*" + p.Type.Name() }
func (p *pointerT) String() string             { return p.RelStr(DefaultPrefixer) }
func (p *pointerT) RelStr(pre Prefixer) string { return "*" + p.Type.RelStr(pre) }
func (p *pointerT) Kind() Kind                 { return PointerKind }
func (p *pointerT) Elem() Type                 { return p.Type }

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) PointerType {
	return &pointerT{t}
}
