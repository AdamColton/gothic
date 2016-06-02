package blueprint

import (
	"strings"
)

type Type interface {
	Name() string
	String() string
	RelStr(pkg string) string
	Package() string
	Kind() string
	Elem() Type
	Key() Type
}

type T struct {
	name string
	pkg  string
	kind string
	elem Type
	key  Type
}

func (t T) Name() string {
	return t.name
}

func (t T) Package() string {
	return t.pkg
}

func (t T) RelStr(pkg string) string {
	if t.kind == "slice" {
		return "[]" + t.elem.RelStr(pkg)
	}
	if t.kind == "ptr" {
		return "*" + t.elem.RelStr(pkg)
	}
	if t.kind == "map" {
		return "map[" + t.key.RelStr(pkg) + "]" + t.elem.RelStr(pkg)
	}
	if t.pkg != "" && t.pkg != pkg {
		return t.pkg + "." + t.name
	}
	return t.name
}

func (t T) String() string {
	return t.RelStr("")
}

func (t T) Kind() string {
	return t.kind
}

func (t T) Elem() Type {
	return t.elem
}

func (t T) Key() Type {
	return t.key
}

func Slice(t Type) Type {
	return Type(T{
		kind: "slice",
		elem: t,
	})
}

func Ptr(t Type) Type {
	return Type(T{
		kind: "ptr",
		elem: t,
	})
}

func Map(k, v Type) Type {
	return Type(T{
		kind: "map",
		elem: v,
		key:  k,
	})
}

func TypeString(s string) Type {
	t, _ := typeString(s)
	return t
}

func typeString(s string) (Type, string) {
	if strings.HasPrefix(s, "[]") {
		t, c := typeString(s[2:])
		return T{
			kind: "slice",
			elem: t,
		}, c
	} else if strings.HasPrefix(s, "*") {
		t, c := typeString(s[1:])
		return T{
			kind: "ptr",
			elem: t,
		}, c
	} else if strings.HasPrefix(s, "map[") {
		key, c := typeString(s[4:])
		val, x := typeString(c[1:])
		return T{
			kind: "map",
			elem: val,
			key:  key,
		}, x
	} else {
		i := strings.Index(s, "]")
		name, rest, pkg := "", "", ""
		if i == -1 {
			name = s
		} else {
			name = s[:i]
			rest = s[i:]
		}
		i = strings.Index(name, ".")
		if i > -1 {
			pkg = name[:i]
			name = name[i+1:]
		}
		if pkg == "" {
			return T{
				name: name,
				kind: name,
				pkg:  pkg,
			}, rest
		}
		return T{
			name: name,
			kind: "struct",
			pkg:  pkg,
		}, rest
	}
}
