package blueprint

import (
	"testing"
)

func TestImporter(t *testing.T) {
	expect := `import (
  "A"
  "B"
)`
	imp := Importer{}
	imp.Import("B")
	imp.Import("A")
	imp.Import("B")
	if imp.ExportImports() != expect {
		t.Error("Importer don't work")
	}
}
