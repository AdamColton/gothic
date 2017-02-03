package gothicgo

type Kind uint8

const (
	NoneKind = Kind(iota)
	SliceKind
	PointerKind
	MapKind
	FuncKind
	StructKind
	UnknownKind
	BuiltinKind
	InterfaceKind
)

// The Type interface represents a type in Go. Name is the type without the
// package, String is the type with the package and RelStr takes a package name
// and return the string representing the type with the package included.
//
// PackageName returns a string representing the package. Package will return
// a *gothicgo.Package if the Type is part of the Gothic generation.
type Type interface {
	Name() string
	String() string
	RelStr(pkg string) string
	PackageName() string
	File() *File
	Kind() Kind
}
