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

func (ctx *Context) serializePtrFunc(t gothicgo.PointerType) (gothicserial.SerializeDef, error) {
	ts := t.String()

	fName := ctx.GetName(t)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	subSerial, err := ctx.Serialize(t.Elem())
	if err != nil {
		return nil, err
	}
	file := subSerial.File()
	if file == nil {
		file = ctx.GetPkg().File("serial.gothic")
	}
	pkg := file.Package()

	sf := &gothicserial.SerializeFuncs{
		MarshalStr:   pkg.Name + "." + marshalFuncName + "(%s)",
		UnmarshalStr: pkg.Name + "." + unmarshalFuncName + "(%s)",
		F:            file,
		Marshaler:    gothicserial.PrependPkgMarshal,
		Unmarshaler:  gothicserial.PrependPkgUnmarshal,
	}
	ctx.Register(t, sf)

	ts = t.RelStr(pkg.Name)

	file.AddCode(
		fmt.Sprintf(marshalPtrTemplate, marshalFuncName, ts, subSerial.Marshal("*s", pkg.Name)),
		fmt.Sprintf(unmarshalPtrTemplate, unmarshalFuncName, ts, subSerial.Unmarshal("b", pkg.Name)),
	)

	file.AddPackageImport(subSerial.PackageName())

	return sf, nil
}
