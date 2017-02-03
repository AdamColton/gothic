package gothicserial

// Registers a Serialization definition by name. Attempting to redefine will
// cause a panic
func RegisterSerializeDef(name string, sf SerializeDef) {
	if _, ok := serializers[name]; ok {
		panic("Attempt to redefine SerializeDef: " + name)
	}
	serializers[name] = sf
}

// Gets a Serialization definition by name, return the SerializeDef and a bool
// if the definition was found. Just a wrapper around a map lookup.
func GetSerializeDef(name string) (SerializeDef, bool) {
	sf, ok := serializers[name]
	return sf, ok
}

var serializers = map[string]SerializeDef{}
