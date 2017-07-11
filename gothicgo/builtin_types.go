package gothicgo

type builtin string

func (b builtin) Name() string             { return string(b) }
func (b builtin) String() string           { return string(b) }
func (b builtin) RelStr(pkg string) string { return string(b) }
func (b builtin) PackageName() string      { return "" }
func (b builtin) File() *File              { return nil }
func (b builtin) Kind() Kind               { return BuiltinKind }

var BoolType = Type(builtin("bool"))
var ByteType = Type(builtin("byte"))
var IntType = Type(builtin("int"))
var Int8Type = Type(builtin("int8"))
var Int16Type = Type(builtin("int16"))
var Int32Type = Type(builtin("int32"))
var Int64Type = Type(builtin("int64"))
var Complex128Type = Type(builtin("complex128"))
var Complex64Type = Type(builtin("complex64"))
var Float32Type = Type(builtin("float32"))
var Float64Type = Type(builtin("float64"))
var RuneType = Type(builtin("rune"))
var StringType = Type(builtin("string"))
var UintType = Type(builtin("uint"))
var Uint8Type = Type(builtin("uint8"))
var Uint16Type = Type(builtin("uint16"))
var Uint32Type = Type(builtin("uint32"))
var Uint64Type = Type(builtin("uint64"))
var UintptrType = Type(builtin("uintptr"))
var ErrorType = Type(builtin("error"))
var EmptyInterfaceType = Type(builtin("interface{}"))
