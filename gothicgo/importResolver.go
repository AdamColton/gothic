package gothicgo

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
)

// TODO: add a function to check if an alias is required and add it if so,
// so if foo is mapped to bar/glorp, return `foo "bar/glorp"

type ImportResolver interface {
	Resolve(pkgName string) string
	Add(pkgName, imprt string)
}

type ManualResolver map[string]string

func (m ManualResolver) Resolve(pkg string) string {
	return m[pkg]
}

func (m ManualResolver) Add(pkg, path string) {
	m[pkg] = path
}

type autoResolver struct {
	packages map[string]string
}

var arSingleton *autoResolver

func AutoResolver() ImportResolver {
	if arSingleton != nil {
		return arSingleton
	}
	arSingleton = &autoResolver{
		packages: map[string]string{},
	}
	for _, path := range build.Default.SrcDirs() {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		children, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			continue
		}
		for _, child := range children {
			if child.IsDir() {
				arSingleton.loadPkg(path, child.Name())
			}
		}
	}
	return arSingleton
}

func (a *autoResolver) Resolve(pkg string) string {
	return a.packages[pkg]
}

func (a *autoResolver) Add(pkg, path string) {
	a.packages[pkg] = path
}

func (a *autoResolver) loadPkg(root, importpath string) {
	shortName := path.Base(importpath)
	if shortName == "testdata" {
		return
	}

	dir := filepath.Join(root, importpath)

	_, pkgName := filepath.Split(importpath)
	if _, ok := a.packages[pkgName]; !ok {
		a.packages[pkgName] = importpath
	}

	pkgDir, err := os.Open(dir)
	if err != nil {
		return
	}
	children, err := pkgDir.Readdir(-1)
	pkgDir.Close()
	if err != nil {
		return
	}
	for _, child := range children {
		name := child.Name()
		if name == "" {
			continue
		}
		if c := name[0]; c == '.' || ('0' <= c && c <= '9') {
			continue
		}
		if child.IsDir() {
			a.loadPkg(root, filepath.Join(importpath, name))
		}
	}
}
