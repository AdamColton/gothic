package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// SliceType extends Type with pointer specific information
type SliceType interface {
	Type
	Elem() Type
}

type sliceT struct {
	Type
}

func (s *sliceT) String() string { return typeToString(s, DefaultPrefixer) }
func (s *sliceT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString("[]")
	s.Type.PrefixWriteTo(sw, p)
	return sw.Rets()
}
func (s *sliceT) Kind() Kind { return SliceKind }
func (s *sliceT) Elem() Type { return s.Type }

// SliceOf returns a SliceType around t.
func SliceOf(t Type) SliceType {
	return &sliceT{t}
}
