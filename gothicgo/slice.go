package gothicgo

// SliceType extends Type with pointer specific information
type SliceType interface {
	Type
	Elem() Type
}

type sliceT struct {
	Type
}

func (s *sliceT) Name() string             { return "[]" + s.Type.Name() }
func (s *sliceT) String() string           { return s.RelStr(DefaultPrefixer) }
func (s *sliceT) RelStr(p Prefixer) string { return "[]" + s.Type.RelStr(p) }
func (s *sliceT) Kind() Kind               { return SliceKind }
func (s *sliceT) Elem() Type               { return s.Type }

// SliceOf returns a SliceType around t.
func SliceOf(t Type) SliceType {
	return &sliceT{t}
}
