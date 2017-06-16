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

func (ctx *Context) serializeMapFunc(t gothicgo.MapType) (gothicserial.SerializeDef, error) {
	ts := t.String()

	fName := ctx.GetName(t)
	marshalFuncName := "Marshal" + fName
	unmarshalFuncName := "Unmarshal" + fName

	keySerial, err := ctx.Serialize(t.Key())
	if err != nil {
		return nil, err
	}
	valSerial, err := ctx.Serialize(t.Elem())
	if err != nil {
		return nil, err
	}
	file := keySerial.File()
	if file == nil {
		if p := valSerial.File(); p != nil {
			file = p
		} else {
			file = ctx.GetPkg().File("serial.gothic")
		}
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
		fmt.Sprintf(marshalMapTemplate, marshalFuncName, ts, keySerial.Marshal("k", pkg.Name), valSerial.Marshal("v", pkg.Name)),
		fmt.Sprintf(unmarshalMapTemplate, unmarshalFuncName, ts, t.Key().Name(), t.Elem().Name(), keySerial.Unmarshal("b", pkg.Name), valSerial.Unmarshal("b", pkg.Name)),
	)

	file.AddPackageImport(keySerial.PackageName())
	file.AddPackageImport(valSerial.PackageName())

	return sf, nil
}
