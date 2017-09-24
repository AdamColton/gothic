package gothicgo

// MapType extends the Type interface with Map specific information
type MapType interface {
	Type
	Elem() Type
	Key() Type
}

type mapT struct {
	key  Type
	elem Type
}

func (m *mapT) Name() string             { return "map[" + m.key.Name() + "]" + m.elem.Name() }
func (m *mapT) String() string           { return m.RelStr(DefaultPrefixer) }
func (m *mapT) RelStr(p Prefixer) string { return "map[" + m.key.RelStr(p) + "]" + m.elem.RelStr(p) }

// PackageRef will always return PkgBuiltin() on Map. Packages of the key and element can
// be inspected independantly
func (m *mapT) PackageRef() PackageRef { return pkgBuiltin }
func (m *mapT) File() *File            { return nil }
func (m *mapT) Kind() Kind             { return MapKind }
func (m *mapT) Elem() Type             { return m.elem }
func (m *mapT) Key() Type              { return m.key }

// MapOf returns a MapType
func MapOf(key, elem Type) MapType {
	return &mapT{
		key:  key,
		elem: elem,
	}
}
