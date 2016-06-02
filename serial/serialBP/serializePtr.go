package serialBP

import (
	"fmt"
	"github.com/adamcolton/gothic/blueprint"
)

const (
	marshalPtrTemplate = `
func %s(s %s) []byte {
  return %s
}`

	unmarshalPtrTemplate = `
func %s(b *[]byte) %s {
  s := %s
  return &s
}`
)

func serializePtrFunc(t blueprint.Type) SerializeFuncs {
	ts := t.String()

	fName := getName(ts)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	sf := SerializeFuncs{
		MarshalStr:   SerialHelperPackage + "." + marshalFuncName + "(%s)",
		UnmarshalStr: SerialHelperPackage + "." + unmarshalFuncName + "(%s)",
		Marshaler:    PrependPkgMarshal,
		Unmarshaler:  PrependPkgUnmarshal,
	}
	serializers[ts] = sf

	subSerial := Serialize(t.Elem())
	sf.Package = subSerial.Package

	addFuncs(sf.Package,
		fmt.Sprintf(marshalPtrTemplate, marshalFuncName, ts, subSerial.Marshal("*s", sf.Package)),
		fmt.Sprintf(unmarshalPtrTemplate, unmarshalFuncName, ts, subSerial.Unmarshal("b", sf.Package)),
	)

	return sf
}
