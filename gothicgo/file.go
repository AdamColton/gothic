package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic"
	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// File represents a Go file. Writer is intended as a hook for testing. If it is
// nil, the code will be written to the file normally, if it's set to an
// io.WriteCloser, it will write to that instead
type File struct {
	*Imports
	generators *gothic.Project
	name       string
	code       []string
	pkg        *Package
	Writer     io.WriteCloser
	Comment    string
}

// Prepare runs prepare on all the generators in the file
func (f *File) Prepare() error { return f.generators.Prepare() }

// AddGenerators to the file
func (f *File) AddGenerators(generators ...gothic.Generator) {
	f.generators.AddGenerators(generators...)
}

// AddCode will add raw strings to the file
func (f *File) AddCode(code ...string) { f.code = append(f.code, code...) }

// Generate the file
func (f *File) Generate() error {
	f.Imports.ResolvePackages(f.Package().ImportResolver())
	f.Imports.RemoveRef(f.pkg)

	err := f.generators.Generate()
	if err != nil {
		return err
	}

	s := append([]string{
		BuildComment(f.Comment, CommentWidth),
		"package " + f.pkg.name + "\n",
		f.Imports.String(),
	}, f.code...)

	code := []byte(strings.Join(s, "\n"))
	fmtCode, fmtErr := format.Source(code)

	wc := f.Writer
	if wc == nil {
		wc, err = f.open()
		if err != nil {
			return err
		}
	}

	if fmtErr == nil {
		_, err = wc.Write(fmtCode)
	} else {
		_, err = wc.Write(code)
		if err == nil {
			err = fmt.Errorf("Failed to format", f.pkg.name+"/"+f.name+".go ", fmtErr)
		}
	}
	if err != nil {
		wc.Close()
		return err
	}
	return wc.Close()
}

// File creates a file within the package. The name should not include ".go"
// which will be automatically appended.
func (p *Package) File(name string) *File {
	if file, exists := p.files[name]; exists {
		return file
	}
	f := &File{
		Imports:    NewImports(p),
		generators: gothic.New(),
		name:       name,
		pkg:        p,
		Comment:    p.Comment,
	}
	p.files[name] = f
	return f
}

func (f *File) open() (io.WriteCloser, error) {
	pth := path.Join(f.pkg.OutputPath, f.name+".go")
	pth, err := filepath.Abs(pth)
	if err != nil {
		return nil, err
	}
	return os.Create(pth)
}

// Package the file is in
func (f *File) Package() *Package { return f.pkg }
