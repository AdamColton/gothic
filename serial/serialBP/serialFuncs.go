package serialBP

import (
	"fmt"
)

type SerializeFuncs struct {
	MarshalStr   string
	UnmarshalStr string
	Package      string
	Marshaler    func(SerializeFuncs, string, string) string
	Unmarshaler  func(SerializeFuncs, string, string) string
}

func (sf SerializeFuncs) Marshal(v, pkg string) string {
	return sf.Marshaler(sf, v, pkg)
}

func (sf SerializeFuncs) Unmarshal(v, pkg string) string {
	return sf.Unmarshaler(sf, v, pkg)
}

func SimpleMarhsal(sf SerializeFuncs, v, pkg string) string {
	return fmt.Sprintf(sf.MarshalStr, v)
}

func SimpleUnmarhsal(sf SerializeFuncs, v, pkg string) string {
	return fmt.Sprintf(sf.UnmarshalStr, v)
}

func PrependPkgMarshal(sf SerializeFuncs, v, pkg string) string {
	p := ""
	if pkg != sf.Package {
		p = sf.Package + "."
	}
	return p + fmt.Sprintf(sf.MarshalStr, v)
}

func PrependPkgUnmarshal(sf SerializeFuncs, v, pkg string) string {
	p := ""
	if pkg != sf.Package {
		p = sf.Package + "."
	}
	return p + fmt.Sprintf(sf.UnmarshalStr, v)
}
