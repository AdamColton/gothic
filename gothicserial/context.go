package gothicserial

import (
	"github.com/adamcolton/gothic/gothicgo"
)

type Context struct {
	serializers map[string]SerializeDef
}

func New() *Context {
	return &Context{
		serializers: make(map[string]SerializeDef),
	}
}

func (c *Context) Register(t gothicgo.Type, serializer SerializeDef) {
	c.serializers[t.String()] = serializer
}

func (c *Context) Get(t gothicgo.Type) SerializeDef {
	return c.serializers[t.String()]
}
