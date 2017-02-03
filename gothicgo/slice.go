package gothicgo

type SliceType interface {
	Type
	Elem() Type
}

type SliceT struct {
	Type
}

func (s *SliceT) Name() string             { return "[]" + s.Type.Name() }
func (s *SliceT) String() string           { return s.RelStr("") }
func (s *SliceT) RelStr(pkg string) string { return "[]" + s.Type.RelStr(pkg) }
func (s *SliceT) Kind() Kind               { return SliceKind }
func (s *SliceT) Elem() Type               { return s.Type }

func SliceOf(t Type) SliceType {
	return &SliceT{t}
}
