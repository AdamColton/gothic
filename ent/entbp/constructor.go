package entbp

import (
	"bytes"
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"strings"
)

// Adds a Entity Constructor to a Struct

func Construct(strct *gothicgo.Struct) *Constructor {
	c := &Constructor{
		strct:    strct,
		argNames: map[string]string{},
	}
	strct.File().AddFragGen(c)
	return c
}

type Constructor struct {
	strct    *gothicgo.Struct
	name     string
	fields   []string
	argNames map[string]string
}

func (c *Constructor) Prepare() {}

func (c *Constructor) Name(name string) *Constructor {
	c.name = name
	return c
}
func (c *Constructor) Fields(fields []string) *Constructor {
	c.fields = fields
	return c
}
func (c *Constructor) ArgName(field, arg string) *Constructor {
	c.argNames[field] = arg
	return c
}

const (
	constructor = "func %s(%s) *%s {\n  o := &%s{\n  Ent:  ent.New(),\n%s  }\n  ent.Register(o)\n  return o}\n"
	argSeg      = "%s %s"
	bodySeg     = "    %s: %s,\n"
)

func (c *Constructor) Generate() []string {
	if len(c.fields) == 0 {
		c.fields = c.strct.Fields()
	}
	sNm := c.strct.Name()
	if c.name == "" {
		c.name = "New" + sNm
	}

	body := &bytes.Buffer{}
	args := []string{}

	cPkg := c.strct.PackageName()

	for _, fieldName := range c.fields {
		if fieldName == "Ent" {
			continue
		}
		field, ok := c.strct.Field(fieldName)
		if !ok {
			panic("Field " + fieldName + " does not exist in " + sNm)
		}
		typeString := field.Type().RelStr(cPkg)
		argName := fieldName
		if arg, ok := c.argNames[fieldName]; ok {
			argName = arg
		}
		args = append(args, fmt.Sprintf(argSeg, argName, typeString))
		fmt.Fprintf(body, bodySeg, fieldName, argName)
	}

	return []string{
		fmt.Sprintf(constructor, c.name, strings.Join(args, ", "), sNm, sNm, body.String()),
	}
}
