package gothicserial

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
)

// Implements SerializeDef. If Pkg is defined and PkgName is left blank,
// Pkg.Name will be used
type SerializeFuncs struct {
	MarshalStr   string
	UnmarshalStr string
	PkgName      string
	F            *gothicgo.File
	Marshaler    func(*SerializeFuncs, string, string) string
	Unmarshaler  func(*SerializeFuncs, string, string) string
}

// Returns the string to call Marshal on v relative to the given package.
func (sf SerializeFuncs) Def() SerializeDef {
	return SerializeDef(&sf)
}

// Returns the string to call Marshal on v relative to the given package.
func (sf *SerializeFuncs) Marshal(v, pkg string) string {
	return sf.Marshaler(sf, v, pkg)
}

// Returns the string to call Marshal on v relative to the given package.
func (sf *SerializeFuncs) Unmarshal(v, pkg string) string {
	return sf.Unmarshaler(sf, v, pkg)
}

// Returns the package name the Marshal and Unmarshal functions are in
func (s *SerializeFuncs) PackageName() string {
	if s.PkgName == "" && s.F != nil {
		return s.F.Package().Name
	}
	return s.PkgName
}

// Returns the package the Marshal and Unmarshal functions are in
func (s *SerializeFuncs) File() *gothicgo.File {
	return s.F
}

// SimpleMarhsal ignores the package and just inserts the variable name into
// the function call
func SimpleMarhsal(sf *SerializeFuncs, v, pkg string) string {
	return fmt.Sprintf(sf.MarshalStr, v)
}

// SimpleUnmarhsal ignores the package and just inserts the variable name into
// the function call
func SimpleUnmarhsal(sf *SerializeFuncs, v, pkg string) string {
	return fmt.Sprintf(sf.UnmarshalStr, v)
}

// PrependPkgMarshal checks if the packages match and if not, prepends
// the package name
func PrependPkgMarshal(sf *SerializeFuncs, v, pkg string) string {
	p := ""
	if pkgName := sf.PackageName(); pkg != pkgName {
		p = pkgName + "."
	}
	return p + fmt.Sprintf(sf.MarshalStr, v)
}

// PrependPkgUnmarshal checks if the packages match and if not, prepends
// the package name
func PrependPkgUnmarshal(sf *SerializeFuncs, v, pkg string) string {
	p := ""
	if pkgName := sf.PackageName(); pkg != pkgName {
		p = pkgName + "."
	}
	return p + fmt.Sprintf(sf.UnmarshalStr, v)
}
