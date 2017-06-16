package serialbp

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
)

func (ctx *Context) Serialize(t gothicgo.Type) (gothicserial.SerializeDef, error) {
	if t == nil {
		panic("Cannot serialize nil Type")
	}
	if serializeDef := ctx.Get(t); serializeDef != nil {
		return serializeDef, nil
	}
	switch t.Kind() {
	case gothicgo.SliceKind:
		st, ok := t.(gothicgo.SliceType)
		if !ok {
			return nil, fmt.Errorf("Cast to SliceType failed on SliceKind")
		}
		return ctx.serializeSliceFunc(st)
	case gothicgo.PointerKind:
		pt, ok := t.(gothicgo.PointerType)
		if !ok {
			return nil, fmt.Errorf("Cast to PointerType failed on PointerKind")
		}
		return ctx.serializePtrFunc(pt)
	case gothicgo.MapKind:
		mt, ok := t.(gothicgo.MapType)
		if !ok {
			return nil, fmt.Errorf("Cast to MapType failed on MapKind")
		}
		return ctx.serializeMapFunc(mt)
	case gothicgo.FuncKind:
		return nil, fmt.Errorf("Cannot serialize func")
	case gothicgo.StructKind:
		return nil, fmt.Errorf("Cannot auto-serialize struct: " + t.String())
	case gothicgo.InterfaceKind:
		return nil, fmt.Errorf("Cannot serialize interface")
	}
	return nil, fmt.Errorf("Cannot serialize unknown kind")
}
