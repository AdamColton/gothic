package serialbp

import (
	"github.com/adamcolton/gothic/gothicserial"
)

func init() {
	gothicserial.RegisterSerializeDef("string", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalString(%s)",
		UnmarshalStr: "serial.UnmarshalString(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("uint", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint(%s)",
		UnmarshalStr: "serial.UnmarshalUint(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("byte", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalByte(%s)",
		UnmarshalStr: "serial.UnmarshalByte(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("uint8", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint8(%s)",
		UnmarshalStr: "serial.UnmarshalUint8(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("uint16", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint16(%s)",
		UnmarshalStr: "serial.UnmarshalUint16(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("uint32", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint32(%s)",
		UnmarshalStr: "serial.UnmarshalUint32(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("uint64", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint64(%s)",
		UnmarshalStr: "serial.UnmarshalUint64(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("int", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt(%s)",
		UnmarshalStr: "serial.UnmarshalInt(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("int8", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt8(%s)",
		UnmarshalStr: "serial.UnmarshalInt8(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("int16", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt16(%s)",
		UnmarshalStr: "serial.UnmarshalInt16(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("int32", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt32(%s)",
		UnmarshalStr: "serial.UnmarshalInt32(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("int64", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt64(%s)",
		UnmarshalStr: "serial.UnmarshalInt64(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	gothicserial.RegisterSerializeDef("[]byte", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalByteSlice(%s)",
		UnmarshalStr: "serial.UnmarshalByteSlice(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
}
