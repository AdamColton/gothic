package blueprint

import (
	"sort"
	"strings"
)

type Importer map[string]bool

func (imp Importer) Import(pkgs ...string) {
	for _, pkg := range pkgs {
		imp[pkg] = true
	}
}

func (imp Importer) ImportGothic(pkgs ...string) {
	for _, pkg := range pkgs {
		imp[GothicIncludePath+pkg] = true
	}
}

func (imp Importer) ImportApp(pkgs ...string) {
	for _, pkg := range pkgs {
		imp[ImportPath+pkg] = true
	}
}

func (imp Importer) ExportImports() string {
  if len(imp) == 0{
    return ""
  }
	packages := make([]string, len(imp))
	i := 0
	for pkg, _ := range imp {
		packages[i] = pkg
		i++
	}
	sort.Strings(packages)
	return "import (\n  \"" + strings.Join(packages, "\"\n  \"") + "\"\n)"
}
