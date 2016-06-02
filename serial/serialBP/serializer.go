package serialBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"strings"
)

var SerialHelperPackage = "gothicHelpers"

func Serialize(t blueprint.Type) SerializeFuncs {
	if serializeFuncs, ok := serializers[t.String()]; ok {
		return serializeFuncs
	}
	switch t.Kind() {
	case "slice":
		return serializeSliceFunc(t)
	case "ptr":
		return serializePtrFunc(t)
	case "map":
		return serializeMapFunc(t)
	case "struct":
		return serializeStructFunc(t)
	}
	return SerializeFuncs{}
}

func getName(typeString string) string {
	fName := strings.Replace(typeString, "[]", "Slice", -1)
	fName = strings.Replace(fName, "[", "", -1)
	fName = strings.Replace(fName, "]", "", -1)
	fName = strings.Replace(fName, ".", "", -1)
	return strings.Replace(fName, "*", "Ptr", -1)
}

func serializeStructFunc(t blueprint.Type) SerializeFuncs {
	ts := t.String()
	sf := SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "Unmarshal" + t.Name() + "(%s)",
		Marshaler:    PrependPkgMarshal,
		Unmarshaler:  PrependPkgUnmarshal,
	}
	serializers[ts] = sf
	return sf
}
