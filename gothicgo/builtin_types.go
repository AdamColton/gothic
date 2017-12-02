package gothicgo

type builtin string

func (b builtin) Name() string                { return string(b) }
func (b builtin) String() string              { return string(b) }
func (b builtin) RelStr(Prefixer) string      { return string(b) }
func (b builtin) PackageRef() PackageRef      { return pkgBuiltin }
func (b builtin) File() *File                 { return nil }
func (b builtin) Kind() Kind                  { return BuiltinKind }
func (b builtin) AsRet() *NameType            { return Ret(b) }
func (b builtin) AsArg(name string) *NameType { return Arg(name, b) }
func (b builtin) Ptr() Type                   { return PointerTo(b) }
func (b builtin) Slice() Type                 { return SliceOf(b) }

// HelpfulType fulfils Type and can also be returned as a pointer, slice, arg
// or return.
type HelpfulType interface {
	Type
	AsRet() *NameType
	AsArg(name string) *NameType
	Ptr() Type
	Slice() Type
}

// Built in Go types
var (
	BoolType           = HelpfulType(builtin("bool"))
	ByteType           = HelpfulType(builtin("byte"))
	IntType            = HelpfulType(builtin("int"))
	Int8Type           = HelpfulType(builtin("int8"))
	Int16Type          = HelpfulType(builtin("int16"))
	Int32Type          = HelpfulType(builtin("int32"))
	Int64Type          = HelpfulType(builtin("int64"))
	Complex128Type     = HelpfulType(builtin("complex128"))
	Complex64Type      = HelpfulType(builtin("complex64"))
	Float32Type        = HelpfulType(builtin("float32"))
	Float64Type        = HelpfulType(builtin("float64"))
	RuneType           = HelpfulType(builtin("rune"))
	StringType         = HelpfulType(builtin("string"))
	UintType           = HelpfulType(builtin("uint"))
	Uint8Type          = HelpfulType(builtin("uint8"))
	Uint16Type         = HelpfulType(builtin("uint16"))
	Uint32Type         = HelpfulType(builtin("uint32"))
	Uint64Type         = HelpfulType(builtin("uint64"))
	UintptrType        = HelpfulType(builtin("uintptr"))
	ErrorType          = HelpfulType(builtin("error"))
	EmptyInterfaceType = HelpfulType(builtin("interface{}"))
)
