package serialBP

import (
	"fmt"
	"github.com/adamcolton/gothic/blueprint"
)

const (
	marshalMapTemplate = `
func %s(s %s) []byte {
  b := serial.MarshalInt(len(s))
  for k, v := range s {
    b = append(b, %s...)
    b = append(b, %s...)
  }
  return b
}`

	unmarshalMapTemplate = `
func %s(b *[]byte) %s {
  l := serial.UnmarshalInt(b)
  m := make(map[%s]%s)
  for i := 0; i < l; i++ {
    m[%s] = %s
  }
  return m
}`
)

func serializeMapFunc(t blueprint.Type) SerializeFuncs {
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

	keySerial := Serialize(t.Key())
	valSerial := Serialize(t.Elem())
	sf.Package = valSerial.Package

	addFuncs(sf.Package,
		fmt.Sprintf(marshalMapTemplate, marshalFuncName, ts, keySerial.Marshal("k", sf.Package), valSerial.Marshal("v", sf.Package)),
		fmt.Sprintf(unmarshalMapTemplate, unmarshalFuncName, ts, t.Key().Name(), t.Elem().Name(), keySerial.Unmarshal("b", sf.Package), valSerial.Unmarshal("b", sf.Package)),
	)

	return sf
}
