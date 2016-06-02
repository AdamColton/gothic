package serialBP

func RegisterSerializeFuncs(name string, sf SerializeFuncs) {
	if _, ok := serializers[name]; ok {
		panic("Attempt to redefine SerializeFuncs: " + name)
	}
	serializers[name] = sf
}

func GetSerializeFuncs(name string) (SerializeFuncs, bool) {
	sf, ok := serializers[name]
	return sf, ok
}

var serializers = map[string]SerializeFuncs{
	"string": SerializeFuncs{
		MarshalStr:   "serial.MarshalString(%s)",
		UnmarshalStr: "serial.UnmarshalString(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"uint": SerializeFuncs{
		MarshalStr:   "serial.MarshalUint(%s)",
		UnmarshalStr: "serial.UnmarshalUint(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"byte": SerializeFuncs{
		MarshalStr:   "serial.MarshalByte(%s)",
		UnmarshalStr: "serial.UnmarshalByte(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"uint8": SerializeFuncs{
		MarshalStr:   "serial.MarshalUint8(%s)",
		UnmarshalStr: "serial.UnmarshalUint8(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"uint16": SerializeFuncs{
		MarshalStr:   "serial.MarshalUint16(%s)",
		UnmarshalStr: "serial.UnmarshalUint16(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"uint32": SerializeFuncs{
		MarshalStr:   "serial.MarshalUint32(%s)",
		UnmarshalStr: "serial.UnmarshalUint32(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"uint64": SerializeFuncs{
		MarshalStr:   "serial.MarshalUint64(%s)",
		UnmarshalStr: "serial.UnmarshalUint64(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"int": SerializeFuncs{
		MarshalStr:   "serial.MarshalInt(%s)",
		UnmarshalStr: "serial.UnmarshalInt(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"int8": SerializeFuncs{
		MarshalStr:   "serial.MarshalInt8(%s)",
		UnmarshalStr: "serial.UnmarshalInt8(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"int16": SerializeFuncs{
		MarshalStr:   "serial.MarshalInt16(%s)",
		UnmarshalStr: "serial.UnmarshalInt16(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"int32": SerializeFuncs{
		MarshalStr:   "serial.MarshalInt32(%s)",
		UnmarshalStr: "serial.UnmarshalInt32(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
	"int64": SerializeFuncs{
		MarshalStr:   "serial.MarshalInt64(%s)",
		UnmarshalStr: "serial.UnmarshalInt64(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  PrependPkgUnmarshal,
	},
	"[]byte": {
		MarshalStr:   "serial.MarshalByteSlice(%s)",
		UnmarshalStr: "serial.UnmarshalByteSlice(%s)",
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	},
}
