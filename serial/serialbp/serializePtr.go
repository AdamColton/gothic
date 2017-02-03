package serialbp

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
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

func serializePtrFunc(t gothicgo.PointerType) gothicserial.SerializeDef {
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
		fmt.Sprintf(marshalPtrTemplate, marshalFuncName, ts, subSerial.Marshal("*s", pkgName)),
		fmt.Sprintf(unmarshalPtrTemplate, unmarshalFuncName, ts, subSerial.Unmarshal("b", pkgName)),
	)

	file.AddPackageImport(subSerial.PackageName())

	return sf
}
