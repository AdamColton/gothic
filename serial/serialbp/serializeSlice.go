package serialbp

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
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

func serializeSliceFunc(t gothicgo.SliceType) gothicserial.SerializeDef {
	ts := t.String()

	fName := getName(ts)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	subSerial := Serialize(t.Elem())
	file := subSerial.File()
	if file == nil {
		file = serialHelperPackage().File("serial.gothic")
	}

	sf := &gothicserial.SerializeFuncs{
		MarshalStr:   marshalFuncName + "(%s)",
		UnmarshalStr: unmarshalFuncName + "(%s)",
		F:            file,
		Marshaler:    gothicserial.PrependPkgMarshal,
		Unmarshaler:  gothicserial.PrependPkgUnmarshal,
	}
	gothicserial.RegisterSerializeDef(ts, sf)

	pkgName := file.Package().Name
	ts = t.RelStr(pkgName)

	file.AddCode(
		fmt.Sprintf(marshalSliceTemplate, marshalFuncName, ts, subSerial.Marshal("i", pkgName)),
		fmt.Sprintf(unmarshalSliceTemplate, unmarshalFuncName, ts, ts, subSerial.Unmarshal("b", pkgName)),
	)

	file.AddPackageImport(subSerial.PackageName())

	return sf
}
