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
func (m *MapT) String() string           { return m.RelStr(nil) }
func (m *MapT) RelStr(i *Imports) string { return "map[" + m.key.RelStr(i) + "]" + m.elem.RelStr(i) }

// PackageRef will always return PkgBuiltin() on Map. Packages of the key and element can
// be inspected independantly
func (m *MapT) PackageRef() PackageRef { return pkgBuiltin }
func (m *MapT) File() *File            { return nil }
func (m *MapT) Kind() Kind             { return MapKind }
func (m *MapT) Elem() Type             { return m.elem }
func (m *MapT) Key() Type              { return m.key }

func MapOf(key, elem Type) MapType {
	return &MapT{
		key:  key,
		elem: elem,
	}
}
