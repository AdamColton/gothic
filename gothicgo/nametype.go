package gothicgo

// NameType is used for arguments and returns for function
type NameType struct {
	N string
	T Type
}

// Name value
func (n *NameType) Name() string { return n.N }

// Type value
func (n *NameType) Type() Type { return n.T }

// Arg is a short way to create a NameType as an Arg for a Func
func Arg(name string, typ Type) *NameType {
	return &NameType{
		N: name,
		T: typ,
	}
}

// Ret is a short way to make an unnamed NameType as a Return for a Func
func Ret(typ Type) *NameType {
	return &NameType{
		T: typ,
	}
}

// NmRet is a short way to make a named NameType as a return for a Func
func NmRet(name string, typ Type) *NameType {
	return &NameType{
		N: name,
		T: typ,
	}
}
