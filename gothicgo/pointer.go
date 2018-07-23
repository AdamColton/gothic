package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// PointerType extends the Type interface with pointer specific information
type PointerType interface {
	Type
	Elem() Type
}

type pointerT struct {
	Type
}

func (p *pointerT) Name() string   { return "*" + p.Type.Name() }
func (p *pointerT) String() string { return typeToString(p, DefaultPrefixer) }
func (p *pointerT) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteRune('*')
	p.Type.PrefixWriteTo(sw, pre)
	return sw.Rets()
}
func (p *pointerT) Kind() Kind { return PointerKind }
func (p *pointerT) Elem() Type { return p.Type }

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) PointerType {
	return &pointerT{t}
}
