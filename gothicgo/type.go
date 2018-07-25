package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"io"
)

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
	TypeDefKind
)

type PrefixWriterTo interface {
	PrefixWriteTo(io.Writer, Prefixer) (int64, error)
}

type RegisterImports interface {
	RegisterImports(*Imports)
}

type Namer interface {
	ScopeName() string
}

// The Type interface represents a type in Go. Name is the type without the
// package, String is the type with the package and PrefixString takes a package name
// and return the string representing the type with the package included.
//
// PackageName returns a string representing the package. Package will return
// a *gothicgo.Package if the Type is part of the Gothic generation.
type Type interface {
	PrefixWriterTo
	RegisterImports
	String() string
	PackageRef() PackageRef
	File() *File // TODO: Remove this
	Kind() Kind
}

// HelpfulType fulfils Type and can also be returned as a pointer, slice, arg
// or return.
type HelpfulType interface {
	Type
	AsRet() NameType
	AsArg(name string) NameType
	Ptr() Type
	Slice() Type
}

type HelpfulTypeWrapper struct{ Type }

func (h HelpfulTypeWrapper) AsRet() NameType            { return Ret(h) }
func (h HelpfulTypeWrapper) AsArg(name string) NameType { return Arg(name, h) }
func (h HelpfulTypeWrapper) Ptr() Type                  { return PointerTo(h) }
func (h HelpfulTypeWrapper) Slice() Type                { return SliceOf(h) }

// temporary until String method is removed from type
func typeToString(t PrefixWriterTo, p Prefixer) string {
	b := bufpool.Get()
	t.PrefixWriteTo(b, p)
	return bufpool.PutStr(b)
}

type StructEmbeddable interface {
	StructEmbedName() string
}
