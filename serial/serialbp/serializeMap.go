package serialbp

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
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

func serializeMapFunc(t gothicgo.MapType) gothicserial.SerializeDef {
	ts := t.String()

	fName := getName(ts)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	keySerial := Serialize(t.Key())
	valSerial := Serialize(t.Elem())
	file := keySerial.File()
	if file == nil {
		if p := valSerial.File(); p != nil {
			file = p
		} else {
			file = serialHelperPackage().File("serial.gothic")
		}
	}

	sf := &gothicserial.SerializeFuncs{
		MarshalStr:   SerialHelperPackage + "." + marshalFuncName + "(%s)",
		UnmarshalStr: SerialHelperPackage + "." + unmarshalFuncName + "(%s)",
		F:            file,
		Marshaler:    gothicserial.PrependPkgMarshal,
		Unmarshaler:  gothicserial.PrependPkgUnmarshal,
	}
	gothicserial.RegisterSerializeDef(ts, sf)

	pkgName := file.Package().Name
	ts = t.RelStr(pkgName)

	file.AddCode(
		fmt.Sprintf(marshalMapTemplate, marshalFuncName, ts, keySerial.Marshal("k", pkgName), valSerial.Marshal("v", pkgName)),
		fmt.Sprintf(unmarshalMapTemplate, unmarshalFuncName, ts, t.Key().Name(), t.Elem().Name(), keySerial.Unmarshal("b", pkgName), valSerial.Unmarshal("b", pkgName)),
	)

	file.AddPackageImport(keySerial.PackageName())
	file.AddPackageImport(valSerial.PackageName())

	return sf
}
