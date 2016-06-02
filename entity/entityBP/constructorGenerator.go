package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
	"strings"
)

type Constructor struct {
	name   string
	ent    *Entity
	fields []string
	bin    bool
}

func (s *Entity) Construct(name string) *Constructor {
	return constr(s, name, false)
}

func (s *Entity) BinaryConstruct(name string) *Constructor {
	return constr(s, name, true)
}

func constr(s *Entity, name string, bin bool) *Constructor {
	method := &Constructor{
		name:   name,
		ent:    s,
		fields: []string{},
		bin:    bin,
	}
	s.AddGenerator(method)
	return method
}

func (constructor *Constructor) Fields(fields ...string) *Constructor {
	constructor.fields = append(constructor.fields, fields...)
	return constructor
}

func (constructor *Constructor) Prepare() {
	constructor.ent.ImportGothic("validation")
}

var constructorHeader = "func %s(%s) (*%s, validation.Errs) {\n  valErrs := validation.Errors()"
var constructorCheckErrs = "  if valErrs.HasErrs(){\n    return nil, valErrs\n  }"
var constructorReturn = "  return &%s{\n    Ent: entity.New(),"
var constructorSetRow = "    %s: %s,"

func (constructor *Constructor) Export() string {
	if constructor.bin {
		return exportBinaryConstructor(constructor)
	}
	return traditionalConstructor(constructor)
}

func traditionalConstructor(constructor *Constructor) string {
	valStr := &blueprint.StrBuilder{}
	setStr := blueprint.StringBuilderLn(constructorReturn, constructor.ent.Name())
	argsArr := []string{}
	for _, fn := range constructor.fields {
		f, _ := constructor.ent.ByName(fn)
		validators, _ := constructor.ent.validators[fn]
		for _, v := range validators {
			valStr.AddLn(v.Export(fn, constructor.ent))
		}
		setStr.AddLn(constructorSetRow, fn, fn)
		argsArr = append(argsArr, fn+" "+f.Type().RelStr(constructor.ent.Package()))
	}
	argsStr := strings.Join(argsArr, ", ")
	str := blueprint.StringBuilderLn(constructorHeader, constructor.name, argsStr, constructor.ent.Name())
	str.AddLn(valStr.String())
	str.AddLn(constructorCheckErrs)
	str.AddLn(setStr.String())
	str.AddLn("  }, valErrs\n}")
	return str.String()
}

var binaryConstructorExtractRow = "  %s := %s\n"

func exportBinaryConstructor(constructor *Constructor) string {
	str := blueprint.StringBuilderLn(constructorHeader, constructor.name, "b *[]byte", constructor.ent.Name())
	for _, fn := range constructor.fields {
		f, _ := constructor.ent.ByName(fn)
		sf := serialBP.Serialize(f.Type())
		str.Add(binaryConstructorExtractRow, fn, sf.Unmarshal("b", constructor.ent.Package()))
		validators, _ := constructor.ent.validators[fn]
		for _, v := range validators {
			str.AddLn(v.Export(fn, constructor.ent))
		}
	}
	str.AddLn(constructorCheckErrs)
	str.AddLn(constructorReturn, constructor.ent.Name())
	for _, fn := range constructor.fields {
		str.AddLn(constructorSetRow, fn, fn)
	}
	str.AddLn("  }, valErrs\n}")
	return str.String()
}
