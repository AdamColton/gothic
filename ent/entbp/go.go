package entbp

import (
	"bytes"
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
	"github.com/adamcolton/gothic/serial/serialbp"
)

// GS is a Go Struct generated from an entity. It includes references to the
// EntBP, the structs reference type and it's unmarshal function.
type GS struct {
	*gothicgo.Struct
	ent            *EntBP
	ref            *gothicgo.Struct
	unmarshalEntFn *gothicgo.Func
}

// Ref returns an entity reference for the underlying struct
func (gs *GS) Ref() gothicgo.Type { return gs.ref.Ptr() }

const (
	refGetBody       = "return Get%s(r.ID())"
	unmarshalEntName = "Unmarshal%sEntity"
	unmarshalEntBody = "return ent.Entity(Unmarshal%s(b))"
	getName          = "Get%s"
	getBody          = `var e ent.Entity
ent.Get(id, &e, Unmarshal%sEntity)
if e == nil{
	return nil
}
return e.(*%s)`
	unmarshalRefName = "Unmarshal%sRef"
	unmarshalRefBody = "return &%sRef{ent.Unmarshal(b)}"
	unmarshalName    = "Unmarshal%s"
)

var btSlc = gothicgo.SliceOf(gothicgo.ByteType)
var btArg = gothicgo.Arg("b", gothicgo.PointerTo(btSlc))
var entRet = gothicgo.Ret(gothicgo.DefStruct("ent.Entity"))
var idArg = gothicgo.Arg("id", btSlc)

// GoStruct generates a Go Struct from an Entity as well as a reference struct.
// The Entity should be fully defined before invoking this. A number of methods
// and function will also be generated.
func (e *EntBP) GoStruct() *GS {
	pkg := e.Package().GoPackage()
	strct := pkg.NewStruct(e.name)
	strct.Embed(EntType)
	for _, fieldName := range e.fieldOrder {
		entField := e.fields[fieldName]
		strct.AddField(entField.name, entField.typ)
	}

	saveMth := strct.NewMethod("Save")
	saveMth.ReceiverName = "o"

	ref := strct.File().NewStruct(e.name + "Ref")
	ref.Embed(EntType)
	refGet := ref.NewMethod("Get")
	refGet.ReceiverName = "r"
	refGet.Returns(strct.AsRet())
	refGet.Body = fmt.Sprintf(refGetBody, e.name)

	unmarshalEntFn := strct.File().NewFunc(fmt.Sprintf(unmarshalEntName, e.name), btArg)
	unmarshalEntFn.Body = fmt.Sprintf(unmarshalEntBody, e.name)
	unmarshalEntFn.Returns(entRet)

	getStructFn := strct.File().NewFunc(fmt.Sprintf(getName, e.name), idArg)
	getStructFn.Returns(strct.AsRet())
	getStructFn.Body = fmt.Sprintf(getBody, e.name, e.name)

	unmarshalRefFn := strct.File().NewFunc(fmt.Sprintf(unmarshalRefName, e.name), btArg)
	unmarshalRefFn.Returns(ref.AsRet())
	unmarshalRefFn.Body = fmt.Sprintf(unmarshalRefBody, e.name)

	refMth := strct.NewMethod("Ref")
	refMth.Returns(ref.AsRet())
	refMth.ReceiverName = "o"

	marshalMth := strct.NewMethod("Marshal")
	marshalMth.Returns(btRet)
	marshalMth.ReceiverName = "o"

	gs := &GS{
		ent:            e,
		Struct:         strct,
		ref:            ref,
		unmarshalEntFn: unmarshalEntFn,
	}

	file := gs.File()
	ptr := gs.Ptr()

	gothicserial.RegisterSerializeDef(ptr.String(), gothicserial.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "Unmarshal" + ptr.Elem().Name() + "(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.PrependPkgUnmarshal,
		F:            file,
	}.Def())

	gothicserial.RegisterSerializeDef(gs.ref.Ptr().String(), gothicserial.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "Unmarshal" + gs.ref.Name() + "(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.PrependPkgUnmarshal,
		F:            file,
	}.Def())

	file.AddFragGen(gs)

	return gs
}

const (
	refBody           = "return &%s{ent.Def(o.ID())}"
	marshalBodyHeader = `	var b []byte
	if o != nil {
		b = []byte{1}
		b = append(b, o.Ent.Marshal()...)
`
	marshalBodyRow      = "\t\tb = append(b, %s...)\n"
	marshalBodyFooter   = "\t} else {\n\t\tb = []byte{0}\n\t}\n\treturn b"
	unmarshalBodyHeader = `	isNil := (*b)[0]
*b = (*b)[1:]
if isNil == 0 {
	return nil
}
return &%s{
	Ent:  ent.Unmarshal(b),
`
	unmarshalBodyRow = "\t\t%s: %s,\n"
)

var btRet = gothicgo.Ret(btSlc)

func (gs *GS) Prepare() {
	mrshBody := &bytes.Buffer{}
	fmt.Fprint(mrshBody, marshalBodyHeader)

	umrshBody := &bytes.Buffer{}
	fmt.Fprintf(umrshBody, unmarshalBodyHeader, gs.Name())

	for _, fieldName := range gs.ent.fieldOrder {
		entField := gs.ent.fields[fieldName]
		structField, _ := gs.Field(entField.name)
		serialDef := serialbp.Serialize(structField.Type())
		gs.File().AddPackageImport(serialDef.PackageName())
		fmt.Fprintf(mrshBody, marshalBodyRow, serialDef.Marshal("o."+fieldName, gs.PackageName()))
		fmt.Fprintf(umrshBody, unmarshalBodyRow, fieldName, serialDef.Unmarshal("b", gs.PackageName()))
	}

	fmt.Fprint(mrshBody, marshalBodyFooter)
	fmt.Fprint(umrshBody, "}")

	saveMth, _ := gs.Method("Save")
	saveMth.Body = "ent.Store(o)"

	refMth, _ := gs.Method("Ref")
	refMth.Body = fmt.Sprintf(refBody, gs.ref.Name())

	marshalMth, _ := gs.Method("Marshal")
	marshalMth.Body = mrshBody.String()

	unmarshalFn := gs.File().NewFunc(fmt.Sprintf(unmarshalName, gs.Name()), btArg)
	unmarshalFn.Returns(gs.AsRet())
	unmarshalFn.Body = umrshBody.String()

}

func (gs *GS) Generate() (empty []string) { return }
