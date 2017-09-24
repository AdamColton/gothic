package gothicgo

// Kind represents the different kinds of types. Two different structs will have
// different types, but the same kind (StructKind)
type Kind uint8

// Defined Kinds
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

// TODO: this concept sort of doesn't make sense. A struct has a name, but a
// func doesn't really have a name without package reference.

// The Type interface represents a type in Go. Name is the type without the
// package, String is the type with the package and RelStr takes a package name
// and return the string representing the type with the package included.
//
// PackageName returns a string representing the package. Package will return
// a *gothicgo.Package if the Type is part of the Gothic generation.
type Type interface {
	Name() string
	String() string
	RelStr(Prefixer) string
	PackageRef() PackageRef
	File() *File
	Kind() Kind
}
