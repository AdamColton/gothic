package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
	"strings"
)

func Serialize(t gothicgo.Type) gothicserial.SerializeDef {
	if t == nil {
		panic("Cannot serialize nil Type")
	}
	if serializeDef, ok := gothicserial.GetSerializeDef(t.String()); ok {
		return serializeDef
	}
	switch t.Kind() {
	case gothicgo.SliceKind:
		st, ok := t.(gothicgo.SliceType)
		if !ok {
			panic("Cast to SliceType failed on SliceKind")
		}
		return serializeSliceFunc(st)
	case gothicgo.PointerKind:
		pt, ok := t.(gothicgo.PointerType)
		if !ok {
			panic("Cast to PointerType failed on PointerKind")
		}
		return serializePtrFunc(pt)
	case gothicgo.MapKind:
		mt, ok := t.(gothicgo.MapType)
		if !ok {
			panic("Cast to MapType failed on MapKind")
		}
		return serializeMapFunc(mt)
	case gothicgo.FuncKind:
		panic("Cannot serialize func")
	case gothicgo.StructKind:
		panic("Cannot auto-serialize struct: " + t.String())
	case gothicgo.InterfaceKind:
		panic("Cannot serialize interface")
	default:
		panic("Cannot serialize unknown kind")
	}
	return nil
}

func getName(typeString string) string {
	//TODO: write tests for getName
	fName := strings.Replace(typeString, "[]", "Slice", -1)
	fName = strings.Replace(fName, "[", "", -1)
	fName = strings.Replace(fName, "]", "", -1)
	fName = strings.Replace(fName, ".", "", -1)
	return strings.Replace(fName, "*", "Ptr", -1)
}
