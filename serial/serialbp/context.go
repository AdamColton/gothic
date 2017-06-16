package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
	"strings"
)

type Context struct {
	*gothicserial.Context
	pkg     *gothicgo.Package
	GetName func(gothicgo.Type) string
}

var namer = strings.NewReplacer("[]", "SliceOf", "map[", "Map", "]", "To", ".", "_", "*", "Ptr")

func getName(t gothicgo.Type) string {
	return namer.Replace(t.String())
}

func (ctx *Context) GetPkg() *gothicgo.Package {
	if ctx.pkg == nil {
		ctx.SetPkgString("serialHelpers")
	}
	return ctx.pkg
}

func (ctx *Context) SetPkg(pkg *gothicgo.Package) {
	ctx.pkg = pkg
}

func (ctx *Context) SetPkgString(pkgName string) {
	ctx.pkg = gothicgo.NewPackage(pkgName)
	gothicgo.AutoResolver().Add(pkgName, ctx.pkg.ImportPath)
}

func New() *Context {
	ctx := &Context{
		Context: gothicserial.New(),
		GetName: getName,
	}

	ctx.Register(gothicgo.StringType, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalString(%s)",
		UnmarshalStr: "serial.UnmarshalString(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.UintType, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint(%s)",
		UnmarshalStr: "serial.UnmarshalUint(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.ByteType, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalByte(%s)",
		UnmarshalStr: "serial.UnmarshalByte(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Uint8Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint8(%s)",
		UnmarshalStr: "serial.UnmarshalUint8(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Uint16Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint16(%s)",
		UnmarshalStr: "serial.UnmarshalUint16(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Uint32Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint32(%s)",
		UnmarshalStr: "serial.UnmarshalUint32(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Uint32Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalUint64(%s)",
		UnmarshalStr: "serial.UnmarshalUint64(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.IntType, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt(%s)",
		UnmarshalStr: "serial.UnmarshalInt(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Int8Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt8(%s)",
		UnmarshalStr: "serial.UnmarshalInt8(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Int16Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt16(%s)",
		UnmarshalStr: "serial.UnmarshalInt16(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Int32Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt32(%s)",
		UnmarshalStr: "serial.UnmarshalInt32(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.Int64Type, gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalInt64(%s)",
		UnmarshalStr: "serial.UnmarshalInt64(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))
	ctx.Register(gothicgo.SliceOf(gothicgo.ByteType), gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "serial.MarshalByteSlice(%s)",
		UnmarshalStr: "serial.UnmarshalByteSlice(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "serial",
	}))

	return ctx
}
