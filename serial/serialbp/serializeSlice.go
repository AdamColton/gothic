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

func (ctx *Context) serializeSliceFunc(t gothicgo.SliceType) (gothicserial.SerializeDef, error) {
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
		fmt.Sprintf(marshalSliceTemplate, marshalFuncName, ts, subSerial.Marshal("i", pkg.Name)),
		fmt.Sprintf(unmarshalSliceTemplate, unmarshalFuncName, ts, ts, subSerial.Unmarshal("b", pkg.Name)),
	)

	file.AddPackageImport(subSerial.PackageName())

	return sf, nil
}
