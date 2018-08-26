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
	InterfaceTypeDefKind
)

type PrefixWriterTo interface {
	PrefixWriteTo(io.Writer, Prefixer) (int64, error)
}

type IgnorePrefixer struct{ io.WriterTo }

func (ip IgnorePrefixer) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	n, err := ip.WriteTo(w)
	return int64(n), err
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
	Named(string) NameType
}

type HelpfulTypeWrapper struct{ Type }

func (h HelpfulTypeWrapper) AsRet() NameType            { return Ret(h) }
func (h HelpfulTypeWrapper) AsArg(name string) NameType { return Arg(name, h) }
func (h HelpfulTypeWrapper) Ptr() Type                  { return PointerTo(h) }
func (h HelpfulTypeWrapper) Slice() Type                { return SliceOf(h) }
func (h HelpfulTypeWrapper) Named(name string) NameType {
	return NameType{name, h}
}

// temporary until String method is removed from type
func typeToString(t PrefixWriterTo, p Prefixer) string {
	b := bufpool.Get()
	t.PrefixWriteTo(b, p)
	return bufpool.PutStr(b)
}
