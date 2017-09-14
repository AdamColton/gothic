package gothicgo

type SliceType interface {
	Type
	Elem() Type
}

type SliceT struct {
	Type
}

func (s *SliceT) Name() string             { return "[]" + s.Type.Name() }
func (s *SliceT) String() string           { return s.RelStr(nil) }
func (s *SliceT) RelStr(i *Imports) string { return "[]" + s.Type.RelStr(i) }
func (s *SliceT) Kind() Kind               { return SliceKind }
func (s *SliceT) Elem() Type               { return s.Type }

func SliceOf(t Type) SliceType {
	return &SliceT{t}
}
