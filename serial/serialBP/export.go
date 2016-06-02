package serialBP

import (
	"fmt"
	"github.com/adamcolton/gothic/blueprint"
	"os"
)

type serializeGenerator string

func (sg serializeGenerator) Export() string {
	return string(sg)
}

func (sg serializeGenerator) Prepare() {}

type serializeBlueprint struct {
	pkg        string
	generators []blueprint.Generator
}

var serializeBlueprints = map[string]*serializeBlueprint{}

func addSerializeGenerator(pkg, fn string) {
	sb, ok := serializeBlueprints[pkg]
	if !ok {
		sb = &serializeBlueprint{
			pkg:        pkg,
			generators: []blueprint.Generator{},
		}
		serializeBlueprints[pkg] = sb
		blueprint.Register(sb)
	}
	sb.AddGenerator(serializeGenerator(fn))
}

var header = ` package %s

import(
  "github.com/adamcolton/gothic/serial"
)
`

func (sb *serializeBlueprint) Export() {
	for _, bp := range serializeBlueprints {
		for _, g := range bp.generators {
			g.Prepare()
		}
	}

	for pkg, bp := range serializeBlueprints {
		dir := blueprint.AppPath + pkg
		os.Mkdir(dir, 0777)
		file, _ := os.Create(dir + "/serializers.gothic.go")
		file.WriteString(fmt.Sprintf(header, pkg))
		for _, g := range bp.generators {
			file.WriteString(g.Export())
		}
	}
}

func (sb *serializeBlueprint) Prepare()        {}
func (sb *serializeBlueprint) Package() string { return sb.pkg }
func (sb *serializeBlueprint) AddGenerator(g blueprint.Generator) {
	sb.generators = append(sb.generators, g)
}

func addFuncs(pkg, marshal, unmarshal string) {
	addSerializeGenerator(pkg, marshal)
	addSerializeGenerator(pkg, unmarshal)
}
