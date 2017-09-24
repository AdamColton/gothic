package gothicgo

type builtin string

func (b builtin) Name() string           { return string(b) }
func (b builtin) String() string         { return string(b) }
func (b builtin) RelStr(Prefixer) string { return string(b) }
func (b builtin) PackageRef() PackageRef { return pkgBuiltin }
func (b builtin) File() *File            { return nil }
func (b builtin) Kind() Kind             { return BuiltinKind }

// Built in Go types
var (
	BoolType           = Type(builtin("bool"))
	ByteType           = Type(builtin("byte"))
	IntType            = Type(builtin("int"))
	Int8Type           = Type(builtin("int8"))
	Int16Type          = Type(builtin("int16"))
	Int32Type          = Type(builtin("int32"))
	Int64Type          = Type(builtin("int64"))
	Complex128Type     = Type(builtin("complex128"))
	Complex64Type      = Type(builtin("complex64"))
	Float32Type        = Type(builtin("float32"))
	Float64Type        = Type(builtin("float64"))
	RuneType           = Type(builtin("rune"))
	StringType         = Type(builtin("string"))
	UintType           = Type(builtin("uint"))
	Uint8Type          = Type(builtin("uint8"))
	Uint16Type         = Type(builtin("uint16"))
	Uint32Type         = Type(builtin("uint32"))
	Uint64Type         = Type(builtin("uint64"))
	UintptrType        = Type(builtin("uintptr"))
	ErrorType          = Type(builtin("error"))
	EmptyInterfaceType = Type(builtin("interface{}"))
)
