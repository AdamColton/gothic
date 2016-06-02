package blueprint

import (
	"os"
  "os/exec"
)

type BP struct {
	Importer
	generators []Generator
	pkg        string
	name       string
}

func New(pkg, name string) *BP {
  return &BP{
    Importer:   Importer{},
    generators: []Generator{},
    pkg:        pkg,
    name:       name,
  }
}

func (bp *BP) Package() string { return bp.pkg }
func (bp *BP) Name() string    { return bp.name }
func (bp *BP) String() string  { return bp.pkg + "." + bp.name }

func (bp *BP) Prepare() {
  // use a classic for loop instead of range so that we can add generators
  // during generator prep, and their prep method will still run
  for i:=0;i<len(bp.generators);i++{
		bp.generators[i].Prepare()
  }
}

func (bp *BP) Export() {
	dir := AppPath + bp.pkg
	os.Mkdir(dir, 0777)
  filename := dir + "/" + bp.name + ".gothic.go"
	file, _ := os.Create(filename)

	file.WriteString("package ")
	file.WriteString(bp.pkg)
	file.WriteString("\n\n")
	file.WriteString(bp.ExportImports())

	for _, gen := range bp.generators {
		file.WriteString("\n\n")
		file.WriteString(gen.Export())
	}

  file.Close()
  cmd := exec.Command("go", "fmt", filename)
  cmd.Start()
  cmd.Wait()
}

func (bp *BP) AddGenerator(gen Generator) {
	bp.generators = append(bp.generators, gen)
}
