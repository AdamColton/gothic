package gothicgo

type MapType interface {
	Type
	Elem() Type
	Key() Type
}

type MapT struct {
	key  Type
	elem Type
}

func (m *MapT) Name() string             { return "map[" + m.key.Name() + "]" + m.elem.Name() }
func (m *MapT) String() string           { return m.RelStr("") }
func (m *MapT) RelStr(pkg string) string { return "map[" + m.key.RelStr(pkg) + "]" + m.elem.RelStr(pkg) }

// PackageName will always return "" on Map. Packages of the key and element can
// be inspected independantly
func (m *MapT) PackageName() string { return "" }
func (m *MapT) File() *File         { return nil }
func (m *MapT) Kind() Kind          { return MapKind }
func (m *MapT) Elem() Type          { return m.elem }
func (m *MapT) Key() Type           { return m.key }

func MapOf(key, elem Type) MapType {
	return &MapT{
		key:  key,
		elem: elem,
	}
}
