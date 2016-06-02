package serialBP

import (
	"fmt"
	"github.com/adamcolton/gothic/blueprint"
)

const (
	marshalSliceTemplate = `
func %s(s %s) []byte {
  b := serial.MarshalInt(len(s))
  for _, i := range s {
    b = append(b, %s...)
  }
  return b
}`

	unmarshalSliceTemplate = `
func %s(b *[]byte) %s {
  l := serial.UnmarshalInt(b)
  s := make(%s, l)
  for i := 0; i < l; i++ {
    s[i] = %s
  }
  return s
}`
)

func serializeSliceFunc(t blueprint.Type) SerializeFuncs {
	ts := t.String()

	fName := getName(ts)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	subSerial := Serialize(t.Elem())
	pkg := subSerial.Package
	if pkg == "" {
		pkg = SerialHelperPackage
	}

	sf := SerializeFuncs{
		MarshalStr:   marshalFuncName + "(%s)",
		UnmarshalStr: unmarshalFuncName + "(%s)",
		Package:      pkg,
		Marshaler:    PrependPkgMarshal,
		Unmarshaler:  PrependPkgUnmarshal,
	}
	serializers[ts] = sf

	ts = t.RelStr(pkg)

	addFuncs(sf.Package,
		fmt.Sprintf(marshalSliceTemplate, marshalFuncName, ts, subSerial.Marshal("i", pkg)),
		fmt.Sprintf(unmarshalSliceTemplate, unmarshalFuncName, ts, ts, subSerial.Unmarshal("b", pkg)),
	)

	return sf
}
